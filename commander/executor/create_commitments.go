package executor

import (
	"context"
	"time"

	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/o11y"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var ErrNotEnoughTxs = NewRollupError("not enough transactions")

type FeeReceiver struct {
	StateID uint32
	TokenID models.Uint256
}

func (c *TxsContext) CreateCommitments(ctx context.Context) ([]models.CommitmentWithTxs, error) {
	spanCtx, span := otel.Tracer("txsContext").Start(ctx, "CreateCommitments")
	defer span.End()

	err := c.verifyTxsCount()
	if err != nil {
		return nil, err
	}

	txType := txtype.TransactionType(c.BatchType)
	mempoolHeap, err := c.storage.NewMempoolHeap(txType)
	if err != nil {
		return nil, err
	}

	commitmentID, err := c.NextCommitmentID()
	if err != nil {
		return nil, err
	}

	commitments := make([]models.CommitmentWithTxs, 0, c.cfg.MaxCommitmentsPerBatch)
	pendingAccounts := make([]models.AccountLeaf, 0)

	for i := uint8(0); len(commitments) < int(c.cfg.MaxCommitmentsPerBatch); i++ {
		var result CreateCommitmentResult
		commitmentID.IndexInBatch = i

		result, err = c.createCommitment(spanCtx, mempoolHeap, commitmentID)
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

	span.SetAttributes(attribute.Int("hubble.commitmentCount", len(commitments)))

	select {
	case <-c.ctx.Done():
		return nil, errors.WithStack(ErrRollupContextCanceled)
	default:
	}

	err = c.registerPendingAccounts(spanCtx, pendingAccounts)
	if err != nil {
		return nil, err
	}

	return commitments, nil
}

func (c *TxsContext) createCommitment(ctx context.Context, mempoolHeap *st.MempoolHeap, commitmentID *models.CommitmentID) (
	CreateCommitmentResult, error,
) {
	spanCtx, span := otel.Tracer("txsContext").Start(ctx, "createCommitment")
	defer span.End()

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

		executeResult, err = c.executeTxsForCommitment(mempoolHeap, feeReceiver)
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

	log.WithFields(o11y.TraceFields(spanCtx)).Printf(
		"Created a %s commitment from %d transactions in %s",
		c.BatchType,
		executeResult.AppliedTxs().Len(),
		duration,
	)
	span.SetAttributes(attribute.Int(
		"hubble.commitmentTxCount",
		executeResult.AppliedTxs().Len(),
	))

	return c.Executor.NewCreateCommitmentResult(executeResult, commitment), nil
}

func (c *TxsContext) executeTxsForCommitment(mempoolHeap *st.MempoolHeap, feeReceiver *FeeReceiver) (
	result ExecuteTxsForCommitmentResult,
	err error,
) {
	// This block is an optimization, if we don't include it we'll instead fail inside
	// ExecuteTxs when we try and fail to pull more txns. It's okay if this gives
	// false negatives, but it should not give false positives.
	count, err := c.storage.CountPendingTxsOfType(txtype.TransactionType(c.BatchType))
	if err != nil {
		return nil, err
	}
	if count < c.minTxsPerCommitment {
		return nil, errors.WithStack(ErrNotEnoughTxs)
	}

	executeTxsResult, err := c.ExecuteTxs(mempoolHeap, feeReceiver)
	if err != nil {
		return nil, err
	}

	if executeTxsResult.AppliedTxs().Len() < int(c.minTxsPerCommitment) {
		return nil, ErrNotEnoughTxs
	}

	// there were enough txs so we were able to successfully build this commitment.
	// Savepoint() deletes these pendingTxs from badger
	err = mempoolHeap.Savepoint()
	if err != nil {
		return nil, err
	}

	return c.Executor.NewExecuteTxsForCommitmentResult(executeTxsResult), nil
}

func (c *TxsContext) updateMempoolMetrics(allTxCount int, txTypeCount uint32) {
	c.commanderMetrics.MempoolSize.Set(float64(allTxCount))

	var gauge prometheus.Gauge

	switch txtype.TransactionType(c.BatchType) {
	case txtype.Transfer:
		gauge = c.commanderMetrics.MempoolSizeTransfer
	case txtype.Create2Transfer:
		gauge = c.commanderMetrics.MempoolSizeCreate2Transfer
	case txtype.MassMigration:
		gauge = c.commanderMetrics.MempoolSizeMassMigration
	}

	gauge.Set(float64(txTypeCount))
}

func (c *TxsContext) verifyTxsCount() error {
	txType := txtype.TransactionType(c.BatchType)

	allMempoolTxs, err := c.storage.GetAllMempoolTransactions()
	if err != nil {
		return err
	}

	txCount := uint32(0)
	for i := range allMempoolTxs {
		if allMempoolTxs[i].TxType == txType {
			txCount += 1
		}
	}

	c.updateMempoolMetrics(len(allMempoolTxs), txCount)

	// TODO: add a test which exercises this branch
	if txCount >= (c.cfg.MinCommitmentsPerBatch * c.cfg.MinTxsPerCommitment) {
		return nil
	}

	oldestTxn, err := c.storage.FindOldestMempoolTransaction(txType)
	if err != nil {
		return err
	}
	if oldestTxn == nil {
		return errors.WithStack(ErrNotEnoughTxs)
	}

	if time.Since(oldestTxn.ReceiveTime.Time) > c.cfg.MaxTxnDelay {
		log.WithFields(log.Fields{
			"hash":        oldestTxn.Hash,
			"from":        oldestTxn.FromStateID,
			"nonce":       oldestTxn.Nonce.Uint64(),
			"receiveTime": oldestTxn.ReceiveTime,
		}).Debug("Forcing a batch because a transaction is older than MaxTxnDelay")
		c.minTxsPerCommitment = 1
		c.minCommitmentsPerBatch = 1
		return nil
	}
	return errors.WithStack(ErrNotEnoughTxs)
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
			log.WithFields(o11y.TraceFields(spanCtx)).Debugf(
				"Submitted a batch account registration transaction. Transaction nonce: %d, hash: %v",
				tx.Nonce(),
				tx.Hash(),
			)
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
