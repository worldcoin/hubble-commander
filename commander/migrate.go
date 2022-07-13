package commander

import (
	"context"
	"fmt"
	"sort"

	"github.com/Worldcoin/hubble-commander/client"
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var errMissingBootstrapNodeURL = fmt.Errorf("bootstrap node URL is required for migration mode")

func (c *Commander) migrate() error {
	nodeURL := c.cfg.Bootstrap.BootstrapNodeURL
	if nodeURL == nil {
		return errors.WithStack(errMissingBootstrapNodeURL)
	}
	log.Printf("Migration mode is on, syncing data from commander instance running at %s\n", *nodeURL)

	hubbleClient := client.NewHubble(*nodeURL, c.cfg.API.AuthenticationKey)
	return c.migrateCommanderData(hubbleClient)
}

func (c *Commander) migrateCommanderData(hubble client.Hubble) error {
	// TODO: we probably don't need to sync these?
	failedTxsCount, err := c.syncFailedTxs(hubble)
	if err != nil {
		return err
	}
	log.Printf("Synced %d failed transaction(s)\n", failedTxsCount)

	pendingBatchesCount, err := c.syncPendingBatches(hubble)
	if err != nil {
		return err
	}
	log.Printf("Synced %d pending batch(es)\n", pendingBatchesCount)

	// txns must be synced after pending batches, because the mempool is only
	// guaranteed to cleanly apply to the pending state

	pendingTxsCount, err := c.syncPendingTxs(hubble)
	if err != nil {
		return err
	}
	log.Printf("Synced %d pending transaction(s)\n", pendingTxsCount)

	c.setMigrate(false)
	return nil
}

func (c *Commander) syncPendingTxs(hubble client.Hubble) (int, error) {
	txs, err := hubble.GetPendingTransactions()
	if err != nil {
		return 0, err
	}

	err = c.addPendingTxs(txs)
	if err != nil {
		return 0, err
	}

	return txs.Len(), nil
}

type byNonce struct {
	array models.GenericTransactionArray
}

func (b byNonce) Len() int {
	return b.array.Len()
}

func (b byNonce) Swap(i, j int) {
	left, right := b.array.At(i), b.array.At(j)
	b.array.Set(i, right)
	b.array.Set(j, left)
}

func (b byNonce) Less(i, j int) bool {
	left, right := b.array.At(i), b.array.At(j)
	leftNonceAddressableValue := left.GetNonce()
	rightNonceAddressableValue := right.GetNonce()
	return leftNonceAddressableValue.Cmp(&rightNonceAddressableValue) <= 0
}

func (c *Commander) addPendingTxs(txs models.GenericTransactionArray) error {
	// these are given to us in a random order, but the mempool will only accept them
	// in the correct order.
	// TODO: this will fail if any transactions are funded by other transactions,
	//       we need to throw them into a non-validating mempool and then read them
	//       out??
	sort.Sort(byNonce{txs})

	for i := 0; i < txs.Len(); i += 1 {
		tx := txs.At(i)
		nonce := tx.GetNonce()
		log.WithFields(log.Fields{
			"nonce": nonce.Uint64(),
			"from": tx.GetFromStateID(),
		}).Debug("tx: ", tx)
	}

	for i := 0; i < txs.Len(); i += 1 {
		tx := txs.At(i)
		err := c.storage.AddMempoolTx(tx)
		if err != nil {
			// TODO: skip this tx if the error is not a badger error
			return err
		}
	}

	return nil
}

func (c *Commander) syncFailedTxs(hubble client.Hubble) (int, error) {
	txs, err := hubble.GetFailedTransactions()
	if err != nil {
		return 0, err
	}

	err = c.storage.AddFailedTransactions(txs)
	if err != nil {
		return 0, err
	}

	return txs.Len(), nil
}

func (c *Commander) syncPendingBatches(hubble client.Hubble) (int, error) {
	pendingBatches, err := hubble.GetPendingBatches()
	if err != nil {
		return 0, err
	}

	for i := range pendingBatches {
		err = c.syncPendingBatch(dtoToModelsBatch(&pendingBatches[i]))
		if err != nil {
			return 0, err
		}
	}

	return len(pendingBatches), nil
}

func (c *Commander) syncPendingBatch(batch *models.PendingBatch) (err error) {
	ctx := executor.NewRollupLoopContext(c.storage, c.client, c.cfg.Rollup, c.metrics, context.Background(), batch.Type)
	defer ctx.Rollback(&err)

	err = ctx.ExecutePendingBatch(batch)
	if err != nil {
		return err
	}

	return ctx.Commit()
}

func dtoToModelsBatch(dtoBatch *dto.PendingBatch) *models.PendingBatch {
	batch := models.PendingBatch{
		ID:              dtoBatch.ID,
		Type:            dtoBatch.Type,
		TransactionHash: dtoBatch.TransactionHash,
		PrevStateRoot:   dtoBatch.PrevStateRoot,
		Commitments:     make([]models.PendingCommitment, 0, len(dtoBatch.Commitments)),
	}

	for i := range dtoBatch.Commitments {
		batch.Commitments = append(batch.Commitments, models.PendingCommitment(dtoBatch.Commitments[i]))
	}

	return &batch
}
