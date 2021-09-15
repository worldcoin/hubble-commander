# Main application struct

`Commander` in `commander.go` is the main App struct.

`rollupLoop` periodically creates `rollupLoopIteration`. Batches have min and max size.


`accounts.go`: `syncAccounts` Synchronizes on-chain account state to internal state.

`batches.go`: `syncBatches` Synchronizes on-chain batch state to internal state.

`new_block.go`: `newBlockLoop` Watches for new blocks.

`syncForward`: Sync a batch of blocks (batched for internal efficiency, not visible from the outside).

`syncTokens`: Synchronizes list of tokens.


TODO:

* Currently, alternates between Transfer and Create2 batches. In the future allocate
  based on demand.
