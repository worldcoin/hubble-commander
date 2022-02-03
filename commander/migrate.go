package commander

import (
	"context"
	"fmt"
	"sort"

	"github.com/Worldcoin/hubble-commander/client"
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/pkg/errors"
)

var errMissingBootstrapNodeURL = fmt.Errorf("bootstrap node URL is required for migration mode")

func (c *Commander) migrate() error {
	if c.cfg.Bootstrap.BootstrapNodeURL == nil {
		return errors.WithStack(errMissingBootstrapNodeURL)
	}

	hubbleClient := client.NewHubble(*c.cfg.Bootstrap.BootstrapNodeURL, c.cfg.API.AuthenticationKey)
	return c.migrateWithClient(hubbleClient)
}

func (c *Commander) migrateWithClient(hubble client.Hubble) error {
	//TODO: fetch pending txs
	//TODO: fetch failed txs

	err := c.syncPendingBatches(hubble)
	if err != nil {
		return err
	}

	c.setMigrate(false)
	return nil
}

func (c *Commander) syncPendingBatches(hubble client.Hubble) error {
	pendingBatches, err := hubble.GetPendingBatches()
	if err != nil {
		return err
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
