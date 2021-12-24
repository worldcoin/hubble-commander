package executor

import (
	"time"

	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNotEnoughTxs = NewRollupError("not enough transactions")
	mockPublicKey   = models.PublicKey{1, 2, 3}
)

type FeeReceiver struct {
	StateID uint32
	TokenID models.Uint256
}

func (c *TxsContext) CreateCommitments() (BatchData, error) {
	// an evil hack because I'm not sure how to fix the test suite
	c.minTxsPerCommitment = c.cfg.MinTxsPerCommitment
	c.minCommitmentsPerBatch = c.cfg.MinCommitmentsPerBatch
	log.Warnf("minTxs=%d minCommit=%d", c.minTxsPerCommitment, c.minCommitmentsPerBatch)

	txQueue, err := c.queryPendingTxs()
	if err != nil {
		return nil, err
	}

	commitmentID, err := c.NextCommitmentID()
	if err != nil {
		return nil, err
	}

	batchData := c.Executor.NewBatchData(c.cfg.MaxCommitmentsPerBatch)
	pendingAccounts := make([]models.AccountLeaf, 0)

	for i := uint8(0); batchData.Len() != int(c.cfg.MaxCommitmentsPerBatch); i++ {
		var result CreateCommitmentResult
		commitmentID.IndexInBatch = i

		result, err = c.createCommitment(txQueue, commitmentID)
		if errors.Is(err, ErrNotEnoughTxs) {
			break
		}
		if err != nil {
			return nil, err
		}

		batchData.AddCommitment(result.Commitment())
		err = c.Executor.GenerateMetaAndWithdrawRoots(batchData, result)
		if err != nil {
			return nil, err
		}
		pendingAccounts = append(pendingAccounts, result.PendingAccounts()...)
	}

	if batchData.Len() < int(c.minCommitmentsPerBatch) {
		return nil, errors.WithStack(ErrNotEnoughCommitments)
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

	return batchData, nil
}

func (c *TxsContext) createCommitment(txQueue *TxQueue, commitmentID *models.CommitmentID) (
	CreateCommitmentResult, error,
) {
	var commitment *models.CommitmentWithTxs
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

		executeResult, err = c.executeTxsForCommitment(txQueue, feeReceiver)
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
		"type": metrics.BatchTypeToMetricsBatchType(commitment.Type),
	})

	log.Printf(
		"Created a %s commitment from %d transactions in %s",
		c.BatchType,
		executeResult.AppliedTxs().Len(),
		duration,
	)

	return c.Executor.NewCreateCommitmentResult(executeResult, commitment), nil
}

func (c *TxsContext) executeTxsForCommitment(txQueue *TxQueue, feeReceiver *FeeReceiver) (
	result ExecuteTxsForCommitmentResult,
	err error,
) {
	pendingTxs := txQueue.PickTxsForCommitment()

	if pendingTxs.Len() < int(c.minTxsPerCommitment) {
		return nil, ErrNotEnoughTxs
	}

	executeTxsResult, err := c.ExecuteTxs(pendingTxs, feeReceiver)
	if err != nil {
		return nil, err
	}

	if executeTxsResult.AppliedTxs().Len() < int(c.minTxsPerCommitment) {
		return nil, ErrNotEnoughTxs
	}

	txQueue.RemoveFromQueue(executeTxsResult.AllTxs())
	return c.Executor.NewExecuteTxsForCommitmentResult(executeTxsResult), nil
}

func (c *TxsContext) setBatchMinimums(pendingTxs models.GenericTransactionArray) {
	oldestTxnTime := models.FindOldestTransactionTime(pendingTxs)
	if oldestTxnTime == nil {
		return
	}

	oldestTxnDelay := time.Now().Sub(oldestTxnTime.Time)

	if oldestTxnDelay > c.cfg.MaxTxnDelay {
		log.Warn("ignoring batch minimums")
		c.minTxsPerCommitment = 1
		c.minCommitmentsPerBatch = 1
	}
}

func (c *TxsContext) queryPendingTxs() (*TxQueue, error) {
	pendingTxs, err := c.storage.GetPendingTransactions(txtype.TransactionType(c.BatchType))
	if err != nil {
		return nil, err
	}

	c.setBatchMinimums(pendingTxs)

	if pendingTxs.Len() < int(c.minTxsPerCommitment*c.minCommitmentsPerBatch) {
		return nil, errors.WithStack(ErrNotEnoughTxs)
	}
	return NewTxQueue(pendingTxs), nil
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
			PublicKey: mockPublicKey,
		})
	}

	err := c.storage.AccountTree.SetInBatch(accounts[len(accounts)-missingAccounts:]...)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
