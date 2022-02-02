package commander

import (
	"context"
	"fmt"
	"sort"

	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/pkg/errors"
	"github.com/ybbus/jsonrpc/v2"
)

//TODO-mig: move to const
const authKeyHeader = "Auth-Key"

func (c *Commander) migrate() error {
	if c.cfg.Bootstrap.BootstrapNodeURL == nil {
		return fmt.Errorf("missing bootstram node url for migration mode")
	}

	client := jsonrpc.NewClientWithOpts(*c.cfg.Bootstrap.BootstrapNodeURL, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			authKeyHeader: c.cfg.API.AuthenticationKey,
		},
	})

	//TODO-mig: fetch pending txs
	//TODO-mig: fetch failed txs

	err := c.fetchPendingBatches(client)
	if err != nil {
		return err
	}

	return nil
}

func (c *Commander) fetchPendingBatches(client jsonrpc.RPCClient) error {
	var pendingBatches []dto.PendingBatch
	err := client.CallFor(&pendingBatches, "admin_getPendingBatches")
	if err != nil {
		return errors.WithStack(err)
	}

	sort.Slice(pendingBatches, func(i, j int) bool {
		return pendingBatches[i].ID.Cmp(&pendingBatches[j].ID) < 0
	})

	for i := range pendingBatches {
		err = c.syncPendingBatch(&pendingBatches[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Commander) syncPendingBatch(batch *dto.PendingBatch) error {
	switch batch.Type {
	case batchtype.Transfer, batchtype.Create2Transfer, batchtype.MassMigration:
		return c.updateTxsState(batch)
	case batchtype.Deposit:
		panic("invalid batch type")
	case batchtype.Genesis:
		panic("invalid batch type")
	}
	panic("invalid batch type")
}

func (c *Commander) updateTxsState(batch *dto.PendingBatch) (err error) {
	ctx := executor.NewTxsContext(c.storage, c.client, c.cfg.Rollup, c.metrics, context.Background(), batch.Type)
	defer ctx.Rollback(&err)

	err = c.storage.AddBatch(&models.Batch{
		ID:              batch.ID,
		Type:            batch.Type,
		TransactionHash: batch.TransactionHash,
	})
	if err != nil {
		return err
	}

	for i := range batch.Commitments {
		err = c.storage.AddCommitment(batch.Commitments[i].Commitment)
		if err != nil {
			return err
		}

		err = c.storage.BatchAddTransaction(batch.Commitments[i].Transactions)
		if err != nil {
			return err
		}

		var feeReceiver *executor.FeeReceiver
		feeReceiver, err = c.getFeeReceiver(batch.Commitments[i].Commitment)
		if err != nil {
			return err
		}
		//TODO-mig: at least check invalid and skipped txs
		_, err = ctx.ExecuteTxs(batch.Commitments[i].Transactions, feeReceiver)
		if err != nil {
			return err
		}
	}

	return ctx.Commit()
}

func (c *Commander) getFeeReceiver(commitment models.Commitment) (*executor.FeeReceiver, error) {
	var stateID uint32
	switch commitment.GetCommitmentBase().Type {
	case batchtype.Transfer, batchtype.Create2Transfer:
		stateID = commitment.ToTxCommitment().FeeReceiver
	case batchtype.Deposit:
		stateID = commitment.ToMMCommitment().Meta.FeeReceiver
	}

	feeReceiver, err := c.storage.StateTree.Leaf(stateID)
	if err != nil {
		return nil, err
	}

	return &executor.FeeReceiver{
		StateID: feeReceiver.StateID,
		TokenID: feeReceiver.TokenID,
	}, nil
}
