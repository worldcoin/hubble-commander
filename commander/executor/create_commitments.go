package executor

import (
	"time"

	"github.com/Worldcoin/hubble-commander/mempool"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var ErrNotEnoughTxs = NewRollupError("not enough transactions")

type FeeReceiver struct {
	StateID uint32
	TokenID models.Uint256
}

func (c *TxsContext) CreateCommitments() ([]models.CommitmentWithTxs, error) {
	err := c.newHeap()
	if err != nil {
		return nil, err
	}
	txController, batchMempool := c.Mempool.BeginTransaction()
	defer txController.Rollback()

	commitmentID, err := c.NextCommitmentID()
	if err != nil {
		return nil, err
	}

	commitments := make([]models.CommitmentWithTxs, 0, c.cfg.MaxCommitmentsPerBatch)
	pendingAccounts := make([]models.AccountLeaf, 0)

	for i := uint8(0); len(commitments) < int(c.cfg.MaxCommitmentsPerBatch); i++ {
		var result CreateCommitmentResult
		commitmentID.IndexInBatch = i

		result, err = c.createCommitment(batchMempool, commitmentID)
		if errors.Is(err, ErrNotEnoughTxs) {
			break
		}
		if err != nil {
			return nil, err
		}

		commitment := result.Commitment()
		commitments = append(commitments, commitment)
		err = c.Executor.GenerateMetaAndWithdrawRoots(commitment, result)
		if err != nil {
			return nil, err
		}
		pendingAccounts = append(pendingAccounts, result.PendingAccounts()...)
	}

	if len(commitments) < int(c.minCommitmentsPerBatch) {
		return nil, errors.WithStack(ErrNotEnoughTxs)
	}

	select {
	case <-c.ctx.Done():
		return nil, errors.WithStack(ErrRollupContextCanceled)
	default:
	}

	err = c.registerPendingAccounts(pendingAccounts)
	if err != nil {
		return nil, err
	}

	txController.Commit()
	return commitments, nil
}

func (c *TxsContext) createCommitment(batchMempool *mempool.TxMempool, commitmentID *models.CommitmentID) (
	CreateCommitmentResult, error,
) {
	var commitment models.CommitmentWithTxs
	var executeResult ExecuteTxsForCommitmentResult

	duration, err := metrics.MeasureDuration(func() error {
		feeReceiver, err := c.getCommitmentFeeReceiver()
		if err != nil {
			return err
		}

		initialStateRoot, err := c.storage.StateTree.Root()
		if err != nil {
			return err
		}

		executeResult, err = c.executeTxsForCommitment(batchMempool, feeReceiver)
		if errors.Is(err, ErrNotEnoughTxs) {
			if uint32(commitmentID.IndexInBatch+1) <= c.minCommitmentsPerBatch {
				return err // No need to revert the StateTree in this case as the DB tx will be rolled back anyway
			}
			if revertErr := c.storage.StateTree.RevertTo(*initialStateRoot); revertErr != nil {
				return revertErr
			}
			return err
		}
		if err != nil {
			return err
		}

		commitment, err = c.BuildCommitment(executeResult, commitmentID, feeReceiver.StateID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	metrics.SaveHistogramMeasurement(duration, c.commanderMetrics.CommitmentBuildDuration, prometheus.Labels{
		"type": metrics.BatchTypeToMetricsBatchType(c.BatchType),
	})

	log.Printf(
		"Created a %s commitment from %d transactions in %s",
		c.BatchType,
		executeResult.AppliedTxs().Len(),
		duration,
	)

	return c.Executor.NewCreateCommitmentResult(executeResult, commitment), nil
}

func (c *TxsContext) executeTxsForCommitment(batchMempool *mempool.TxMempool, feeReceiver *FeeReceiver) (
	result ExecuteTxsForCommitmentResult,
	err error,
) {
	if c.Mempool.TxCount() < int(c.minTxsPerCommitment) {
		return nil, errors.WithStack(ErrNotEnoughTxs)
	}

	txController, commitmentMempool := batchMempool.BeginTransaction()
	defer txController.Rollback()

	executeTxsResult, err := c.ExecuteTxs(commitmentMempool, feeReceiver)
	if err != nil {
		return nil, err
	}

	if executeTxsResult.AppliedTxs().Len() < int(c.minTxsPerCommitment) {
		return nil, ErrNotEnoughTxs
	}

	txController.Commit()
	return c.Executor.NewExecuteTxsForCommitmentResult(executeTxsResult), nil
}

func (c *TxsContext) verifyTxsCount() error {
	if c.Mempool.TxCount() >= int(c.cfg.MinCommitmentsPerBatch*c.cfg.MinTxsPerCommitment) {
		return nil
	}

	oldestTxnTime := c.findOldestTransactionTime()
	if oldestTxnTime == nil {
		return errors.WithStack(ErrNotEnoughTxs)
	}

	if time.Since(oldestTxnTime.Time) > c.cfg.MaxTxnDelay {
		log.Debug("Creating a batch because a transaction is older than MaxTxnDelay")
		c.minTxsPerCommitment = 1
		c.minCommitmentsPerBatch = 1
		return nil
	}
	return errors.WithStack(ErrNotEnoughTxs)
}

func (c *TxsContext) newHeap() error {
	err := c.verifyTxsCount()
	if err != nil {
		return err
	}

	txs := c.Mempool.GetExecutableTxs(txtype.TransactionType(c.BatchType))
	c.heap = mempool.NewTxHeap(txs...)
	return nil
}

func (c *TxsContext) getCommitmentFeeReceiver() (*FeeReceiver, error) {
	commitmentTokenID := models.MakeUint256(0) // TODO support multiple tokens
	feeReceiverState, err := c.storage.GetFeeReceiverStateLeaf(c.cfg.FeeReceiverPubKeyID, commitmentTokenID)
	if err != nil {
		return nil, err
	}
	return &FeeReceiver{
		StateID: feeReceiverState.StateID,
		TokenID: feeReceiverState.TokenID,
	}, nil
}

func (c *TxsContext) registerPendingAccounts(accounts []models.AccountLeaf) error {
	accounts, err := c.fillMissingAccounts(accounts)
	if err != nil {
		return err
	}
	publicKeys := make([]models.PublicKey, 0, st.AccountBatchSize)
	for i := range accounts {
		publicKeys = append(publicKeys, accounts[i].PublicKey)
		if len(publicKeys) == st.AccountBatchSize {
			tx, err := c.client.RegisterBatchAccount(publicKeys)
			if err != nil {
				return err
			}
			log.Debugf("Submitted a batch account registration transaction. Transaction nonce: %d, hash: %v", tx.Nonce(), tx.Hash())
			publicKeys = make([]models.PublicKey, 0, st.AccountBatchSize)
		}
	}
	return nil
}

func (c *TxsContext) fillMissingAccounts(accounts []models.AccountLeaf) ([]models.AccountLeaf, error) {
	missingAccounts := st.AccountBatchSize - len(accounts)%st.AccountBatchSize
	if missingAccounts == st.AccountBatchSize {
		return accounts, nil
	}
	for i := 0; i < missingAccounts; i++ {
		lastAccount := &accounts[len(accounts)-1]
		accounts = append(accounts, models.AccountLeaf{
			PubKeyID:  lastAccount.PubKeyID + 1,
			PublicKey: models.ZeroPublicKey,
		})
	}

	err := c.storage.AccountTree.SetInBatch(accounts[len(accounts)-missingAccounts:]...)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
