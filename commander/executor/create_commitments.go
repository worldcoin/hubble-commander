package executor

import (
	"context"
	"time"

	"github.com/Worldcoin/hubble-commander/mempool"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/o11y"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

var ErrNotEnoughTxs = NewRollupError("not enough transactions")

type FeeReceiver struct {
	StateID uint32
	TokenID models.Uint256
}

func (c *TxsContext) CreateCommitments(ctx context.Context) ([]models.CommitmentWithTxs, error) {
	spanCtx, span := otel.Tracer("txsContext").Start(ctx, "CreateCommitments")
	defer span.End()

	err := c.newHeap()
	if err != nil {
		return nil, err
	}

	log.WithFields(o11y.TraceFields(ctx)).Info("Pre transaction")
	txController, batchMempool := c.Mempool.BeginTransaction()
	defer txController.Rollback()

	log.WithFields(o11y.TraceFields(ctx)).Info("Started transaction")

	commitmentID, err := c.NextCommitmentID()
	if err != nil {
		return nil, err
	}

	commitments := make([]models.CommitmentWithTxs, 0, c.cfg.MaxCommitmentsPerBatch)
	pendingAccounts := make([]models.AccountLeaf, 0)

	for i := uint8(0); len(commitments) < int(c.cfg.MaxCommitmentsPerBatch); i++ {
		var result CreateCommitmentResult
		commitmentID.IndexInBatch = i

		result, err = c.createCommitment(ctx, batchMempool, commitmentID)
		log.WithFields(o11y.TraceFields(ctx)).Info("CommitmentID: %v, %s", commitmentID)

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
	log.WithFields(o11y.TraceFields(ctx)).Info("Created all commitments")
	if len(commitments) < int(c.minCommitmentsPerBatch) {
		log.WithFields(o11y.TraceFields(ctx)).Infof("Not enough commitments per batch, %s", len(commitments))
		return nil, errors.WithStack(ErrNotEnoughTxs)
	}

	log.WithFields(o11y.TraceFields(ctx)).Info("Enough commitments per batch")

	select {
	case <-c.ctx.Done():
		log.WithFields(o11y.TraceFields(ctx)).Info("Rollup context cancelled")
		return nil, errors.WithStack(ErrRollupContextCanceled)
	default:
	}

	log.WithFields(o11y.TraceFields(ctx)).Info("Before registering accounts")
	err = c.registerPendingAccounts(spanCtx, pendingAccounts)
	if err != nil {
		return nil, err
	}

	log.WithFields(o11y.TraceFields(ctx)).Info("Registered accounts")

	log.WithFields(o11y.TraceFields(ctx)).Info("Before commit")
	txController.Commit()
	log.WithFields(o11y.TraceFields(ctx)).Info("Committed transaction")
	return commitments, nil
}

func (c *TxsContext) createCommitment(ctx context.Context, batchMempool *mempool.TxMempool, commitmentID *models.CommitmentID) (CreateCommitmentResult, error) {
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

		log.WithFields(o11y.TraceFields(ctx)).Info("XXX:before executing txs for commitment")
		executeResult, err = c.executeTxsForCommitment(ctx, batchMempool, feeReceiver)
		log.WithFields(o11y.TraceFields(ctx)).Info("XXX:executed txs for commitment, error: %s", err)
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

func (c *TxsContext) executeTxsForCommitment(ctx context.Context, batchMempool *mempool.TxMempool, feeReceiver *FeeReceiver) (result ExecuteTxsForCommitmentResult, err error) {
	count := c.Mempool.TxCount(txtype.TransactionType(c.BatchType))

	log.WithFields(o11y.TraceFields(ctx)).Infof("XXX: transaction count: %d, batch: %s", c.Mempool.TxCount(txtype.TransactionType(c.BatchType)), c.BatchType.String())

	if c.Mempool.TxCount(txtype.TransactionType(c.BatchType)) < int(c.minTxsPerCommitment) {
		log.WithFields(o11y.TraceFields(ctx)).Infof("XXX:not enough tx point 1, count: %d", count)
		return nil, errors.WithStack(ErrNotEnoughTxs)
	}

	txController, commitmentMempool := batchMempool.BeginTransaction()
	defer txController.Rollback()

	executeTxsResult, err := c.ExecuteTxs(commitmentMempool, feeReceiver)
	if err != nil {
		return nil, err
	}

	if executeTxsResult.AppliedTxs().Len() < int(c.minTxsPerCommitment) {
		log.WithFields(o11y.TraceFields(ctx)).Info("XXX:not enough tx point 2")
		return nil, ErrNotEnoughTxs
	}

	txController.Commit()
	return c.Executor.NewExecuteTxsForCommitmentResult(executeTxsResult), nil
}

func (c *TxsContext) verifyTxsCount() error {
	if c.Mempool.TxCount(txtype.TransactionType(c.BatchType)) >= int(c.cfg.MinCommitmentsPerBatch*c.cfg.MinTxsPerCommitment) {
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

func (c *TxsContext) registerPendingAccounts(ctx context.Context, accounts []models.AccountLeaf) error {
	spanCtx, span := otel.Tracer("txsContext").Start(ctx, "registerPendingAccounts")
	defer span.End()

	accounts, err := c.fillMissingAccounts(accounts)
	if err != nil {
		return err
	}
	publicKeys := make([]models.PublicKey, 0, st.AccountBatchSize)
	for i := range accounts {
		publicKeys = append(publicKeys, accounts[i].PublicKey)
		if len(publicKeys) == st.AccountBatchSize {
			tx, err := c.client.RegisterBatchAccount(spanCtx, publicKeys)
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
