# Main application struct

`Commander` in `commander.go` is the main App struct.

`rollupLoop` periodically creates `rollupLoopIteration`. Batches have min and max size.


`accounts.go`: `syncAccounts` Synchronizes on-chain account state to internal state.

`batches.go`: `syncBatches` Synchronizes on-chain batch state to internal state.

`registered_tokens.go`: `syncTokens` Synchronizes on-chain registered tokens to internal state.

`new_block.go`: `newBlockLoop` Watches for new blocks. Triggers syncing methods.

`syncForward`: Sync a batch of blocks (batched for internal efficiency, not visible from the outside).


TODO:

* Currently, alternates between Transfer and Create2 batches. In the future allocate
  based on demand.
