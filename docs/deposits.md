# Deposits

## Flow

1. User calls `DepositManager.depositFor` method.
    1. ERC20 tokens are transferred to the `Vault` contract.
    2. A `DepositQueued` event is emitted including `pubkeyID`, `tokenID` and `l2Amount`.
    3. A new `UserState` is created, encoded and hashed. The leaf hash is stored in a **baby tree**.
2. Other users call the `DepositManager.depositFor` method.
    1. Once the baby tree is full the root hash is added to a FIFO queue for submission in a batch.
    2. A `DepositSubTreeReady` event is emitted including `subtreeID` and `subtreeRoot`.
3. In the meantime commander picks up individual `DepositQueued` events and adds them to **Deposits** table. Once commander
   notices `DepositSubTreeReady` it gathers corresponding deposits and stores a record in **Pending Deposit Subtrees** table.
4. Rollup loop reads from Pending Deposit Subtrees table and submits deposit batches on chain. After a deposit subtree is submitted it is
   removed from Pending Deposit Subtrees table and a corresponding commitment and pending batch are stored.
5. Once the submission transaction is mined, the commander syncs it and marks the batch as mined.
