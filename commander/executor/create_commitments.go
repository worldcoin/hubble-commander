package executor

import (
	"context"
	"time"

	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
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

	// TODO: here we create a heap which is populated from all the executable txs in
	//       our mempool. We then call c.createCommitment a few times, which is just
	//       a wrapper for c.executeTxsForCommitment, which is just a wrapper for
	//       c.ExecuteTxs, which is the only place we ever interact with the heap.
	//       Stashing the heap on this struct is very convenient but obscures the
	//       dataflow, and means the heap lives longer than it should. The heap is
	//       intimately related to the mempool, right? Why isn't it attached to that?

	err := c.verifyTxsCount()
	if err != nil {
		return nil, err
	}

	// TODO: so here we open a transaction and do not update any of the buckets until
	//       we have registered the pending accounts and are about to return a
	//       commitment.
	//       (I) What happens if there is an error in SubmitBatch?
	//           we are caled by commander/executor/create_batch.go:23
	//           if the call to SubmitBatch fails then we will have dropped our txns
	//           from the mempool without saving the batch to the chain or the db.
	//           - we should commit to the db first, then commit the mempool changes
	//       (II) What happens if the api adds some transactions while the txn is
	//            open? The rollup loop will overwrite any changes made to these
	//            buckets since it opened its txn. With the right data layout that
	//            could be okay. Locking will not save us, because the rollup loop
	//            will keep the txn open for much longer than we should block API
	//            requests.

	// TODO: we are already inside a transaction, right?
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

		result, err = c.createCommitment(mempoolHeap, commitmentID)
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

	err = c.registerPendingAccounts(spanCtx, pendingAccounts)
	if err != nil {
		return nil, err
	}

	return commitments, nil
}

func (c *TxsContext) createCommitment(mempoolHeap *st.MempoolHeap, commitmentID *models.CommitmentID) (
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

	log.Printf(
		"Created a %s commitment from %d transactions in %s",
		c.BatchType,
		executeResult.AppliedTxs().Len(),
		duration,
	)

	return c.Executor.NewCreateCommitmentResult(executeResult, commitment), nil
}

func (c *TxsContext) executeTxsForCommitment(mempoolHeap *st.MempoolHeap, feeReceiver *FeeReceiver) (
	result ExecuteTxsForCommitmentResult,
	err error,
) {
	// TODO: replace this w our own?
	/*
	if c.Mempool.TxCount(txtype.TransactionType(c.BatchType)) < int(c.minTxsPerCommitment) {
		return nil, errors.WithStack(ErrNotEnoughTxs)
	}
	*/

	// TODO: we open this transaction so that our fetched txs will be not be taken
	// from the mempool in the event that we fail to build a complete commitment.
	// (I) this appears to be redundant w the transaction opened by
	//     commander/executor/create_commitments.go:48
	// (II) during rollbacks we don't appear to put the txns back onto c.heap? I think
	//      this might be a bug.

	// I think both are fine?
	// - if we return an exceptional error then everything is getting unwound anyway,
	//   the heap may be in a bad state but it's never used again
	// - if we fail because we run out of txs then this commitment will not be added
	//   to the batch and the heap will be in a bad state, but we never read from the
	//   heap again because we've moved on to bundling the successful commitments into
	//   a batch.

	// Actually, here's why it's not redundant: we actually want the nested txn! If we
	// fail by running out of txns then the incomplete commitment will be thrown away,
	// we do not want those txns to be pulled from the mempool when the batchMempool
	// commits.

	executeTxsResult, err := c.ExecuteTxs(mempoolHeap, feeReceiver)
	if err != nil {
		return nil, err
	}

	if executeTxsResult.AppliedTxs().Len() < int(c.minTxsPerCommitment) {
		return nil, ErrNotEnoughTxs
	}

	// there were enough txs so we were able to successfully build this commitment.
	// Savepoint() removes these pendingTxs from the badger transaction, when
	// badger commits they will be removed from the mempool.
	mempoolHeap.Savepoint()

	return c.Executor.NewExecuteTxsForCommitmentResult(executeTxsResult), nil
}

func (c *TxsContext) verifyTxsCount() error {
	txType := txtype.TransactionType(c.BatchType)

	txCount, err := c.storage.CountPendingTxsOfType(txType)
	if err != nil {
		return err
	}

	// TODO: add a test which exercises this branch
	if txCount >= (c.cfg.MinCommitmentsPerBatch * c.cfg.MinTxsPerCommitment) {
		return nil
	}

	oldestTxn, err := c.storage.FindOldestMempoolTransaction(txType)
	if err != nil {
		return err
	}
	if oldestTxn == nil {
		/*
		log.WithFields(log.Fields{
			"txType": txType,
		}).Debug("quitting loop, no txs of desired type")
		*/
		return errors.WithStack(ErrNotEnoughTxs)
	}

	if time.Since(oldestTxn.ReceiveTime.Time) > c.cfg.MaxTxnDelay {
		log.WithFields(log.Fields{
			"hash": oldestTxn.Hash,
			"from": oldestTxn.FromStateID,
			"nonce": oldestTxn.Nonce.Uint64(),
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
