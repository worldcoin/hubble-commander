package admin

import (
	"context"

	"github.com/Worldcoin/hubble-commander/models"
)

// reverts this batch and all batches following it
// does so by jumping the state tree back to a previous state
// and inserting all batched transactions back into the mempool
func (a *API) RevertBatches(ctx context.Context, id models.Uint256) error {
	// see commander/executor/revert_batches.go

	// TODO: open a transaction, and probably take out a big lock to exclude the RollupLoop

	// TODO: return some kind of indication of what was performed?

	batch, err := a.storage.GetBatch(id)
	if err != nil {
		return err // TODO: sanitize error
	}

	// TODO: wrap this in a transaction, it performs a lot of
	//       operations non-transactionally, and should probably
	//       be called unsafeRevertBatches.
	return a.storage.RevertBatches(batch) // TODO: sanitize error
}
