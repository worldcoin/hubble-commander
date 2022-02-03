package commander

import (
	"context"
	"fmt"
	"sort"

	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/pkg/errors"
	"github.com/ybbus/jsonrpc/v2"
)

//TODO-mig: move to const
const authKeyHeader = "Auth-Key"

func (c *Commander) migrate() error {
	if c.cfg.Bootstrap.BootstrapNodeURL == nil {
		return fmt.Errorf("bootstrap node is required for migration mode")
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

	c.setMigrate(false)
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

func (c *Commander) syncPendingBatch(batch *dto.PendingBatch) (err error) {
	ctx := executor.NewRollupLoopContext(c.storage, c.client, c.cfg.Rollup, c.metrics, context.Background(), batch.Type)
	defer ctx.Rollback(&err)

	err = ctx.ExecutePendingBatch(batch)
	if err != nil {
		return err
	}

	return ctx.Commit()
}
