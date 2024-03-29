# 🦡 Badger data structures

## Stored structures
<!-- toc -->

## Notes and design rationale
- Some indices are specified using badgerhold tags on struct fields.
  Others by implementing the `Indexes()` method of `bh.Storer` interface.
- Watch out when using indices on fields of pointer types.
  By default, badgerhold adds IDs of all structs that have the indexed field set to `nil` to a single index entry, for instance: 

  `_bhIndex:Batch:Hash:nil -> bh.KeyList{ batchID1, batchID2, ... }`  

  Such index entry can quickly grow in size.
  For structs that have indices on fields of pointer type we can implement `Indexes()` method and specify our own `IndexFunc`.
  Returning `nil` from such `IndexFunc` for `nil` field values prevents creation of the `nil` index entry.
  In case you go down this path, make sure to initialise the index using `initializeIndex` function (see `storage/initialize_index.go` for more info).
- All transaction details were previously held in a single **Transaction** structure. 
  We had to split transaction details between **Stored Transaction** and **Stored Transaction Receipt** because of conflict 
  between API `hubble_sendTransaction` method and Rollup loop iterations. If we kept all data in the same structure the API method would be 
  adding new keys updating the indices. At the same time Rollup loop would be modifying the stored transactions updating the same indices.
  As a result some DB transactions would error with `bh.ErrConflict`.

## State Tree

### State Leaf

- Holds `UserState` data

Key: state ID `uint32`

Value: `stored.StateLeaf`

```go
type StateLeaf struct {
    StateID  uint32
    DataHash common.Hash
    PubKeyID uint32
    TokenID  Uint256
    Balance  Uint256
    Nonce    Uint256
}
```

#### Index on `PubKeyID`

Key: pub key ID `uint32`

Value: list of state IDs `bh.KeyList`


### State Tree Node

- Holds state tree nodes (hashes). The leaf level nodes are hashes of `UserState` structs.

Key: node path `models.NamespacedMerklePath`

Prefix: `bh_MerkleTreeNode:state`

Value: node `common.Hash` (through clever encoding of `models.MerkleTreeNode`)


### State Update

- Holds state tree updates including leaves data that were previously in the tree

Key: auto-incremented ID `uint64`

Value: `models.StateUpdate`

```go
type StateUpdate struct {
    ID            uint64 `badgerhold:"key"`
    CurrentRoot   common.Hash
    PrevRoot      common.Hash
    PrevStateLeaf StateLeaf
}

type StateLeaf struct {
    StateID  uint32
    DataHash common.Hash
    UserState
}

type UserState struct {
    PubKeyID uint32
    TokenID  Uint256
    Balance  Uint256
    Nonce    Uint256
}
```

## Account Tree

### Account Leaf

- Holds pubKeyID - publicKey mapping

Key: pubKeyID `uint32`

Value: {pubKeyID, publicKey} `models.AccountLeaf`

#### Index on `PublicKey`

Key: publicKey `models.PublicKey`

Prefix: `_bhIndex:AccountLeaf:PublicKey:`

Value: list of pubKeyIDs `bh.KeyList`


### Account Tree Node

- Holds account tree nodes (hashes). The leaf level nodes are hashes of public keys.

Key: node path `models.NamespacedMerklePath`

Prefix: `bh_MerkleTreeNode:account`

Value: node `common.Hash` (through clever encoding of `models.MerkleTreeNode`)

## Transactions

### Pending Transaction

Key: tx hash `common.Hash`

Value: `stored.PendingTx`

```go
type PendingTx struct {
    Hash        common.Hash
    TxType      txtype.TransactionType
    FromStateID uint32
    Amount      models.Uint256
    Fee         models.Uint256
    Nonce       models.Uint256
    Signature   models.Signature
    ReceiveTime *models.Timestamp

    Body TxBody // interface
}
```

#### Transfer
Body: `stored.TxTransferBody`

```go
type TxTransferBody struct {
    ToStateID uint32
}
```

#### Create2Transfer
Body: `stored.TxCreate2TransferBody`

```go
type TxCreate2TransferBody struct {
    ToPublicKey models.PublicKey
    ToStateID   *uint32
}
```

#### MassMigration
Body: `stored.TxMassMigrationBody`

```go
type TxMassMigrationBody struct {
    SpokeID uint32
}
```

### Batched Transaction

- Stores transactions which have been added to a batch and submitted to the chain

Key: (BatchID, IndexInBatch, IndexInCommitment) `models.CommitmentSlot`

Value: `stored.BatchedTx`

```go
type BatchedTx struct {
	PendingTx
	ID stored.CommitmentSlot
}
```

### Failed Transaction

- Stores transactions which the rollup loop attempted to add to a batch, but failed

Key: tx hash `common.Hash`

Value: `stored.FailedTx`

```go
type FailedTx struct {
	PendingTx

	ErrorMessage *string
}
```

#### Index on `Hash`

Key: Hash `common.Hash`

Prefix: `_bhIndex:BatchedTx:Hash:`

Value: list `bh.KeyList` with single element, the tx

## Deposits

### Pending Deposit

- Holds individual pending deposits until they are moved into **Pending Deposit Subtrees**. Individual deposits correspond to `DepositQueued` events emitted by the SCs.  

Key: <subtreeID, depositIndex> `models.DepositID`

Value: `models.PendingDeposit`

```go
type DepositID struct {
    SubtreeID    Uint256 // the subtree which contains this deposit
    DepositIndex Uint256 // deposit number in the subtree
}
```

```go
type PendingDeposit struct {
    ID         DepositID
    ToPubKeyID uint32
    TokenID    Uint256
    L2Amount   Uint256
}
```

### Pending Deposit Subtree

- FIFO queue for ready Deposit Subtrees that can be submitted in a new Batch

Key: SubtreeID `models.Uint256`

Value: `models.PendingDepositSubtree`

```go
type PendingDepositSubtree struct {
    ID       models.Uint256          // assigned in SC
    Root     common.Hash             // subtree root that will be inserted into the state tree
    Deposits []models.PendingDeposit // deposits included in the subtree
}
```

## Commitments and Batches

### Stored Commitment
- Hold commitment details for both transaction and deposit commitmetns


Key: {BatchID, IndexInBatch} `models.CommitmentID`

Value: `stored.Commitment`

```go
type Commitment struct {
    models.CommitmentBase
    Body CommitmentBody // interface
}

type CommitmentBase struct {
    ID            CommitmentID
    Type          batchtype.BatchType
    PostStateRoot common.Hash
}

type CommitmentID struct {
    BatchID      Uint256
    IndexInBatch uint8
}
```

#### Transfer / Create2Transfer Commitment

Body: `stored.TxCommitmentBody`

```go
type TxCommitmentBody struct {
    FeeReceiver       uint32
    CombinedSignature models.Signature
    BodyHash          *common.Hash
}
```

#### MassMigration Commitment

Body: `stored.MMCommitmentBody`

```go
type MMCommitmentBody struct {
    CombinedSignature models.Signature
    BodyHash          *common.Hash
    Meta              models.MassMigrationMeta
    WithdrawRoot      common.Hash
}

type MassMigrationMeta struct {
    SpokeID     uint32
    TokenID     Uint256
    Amount      Uint256
    FeeReceiver uint32
}
```

#### Deposit Commitment

- When Deposit batch is created data is moved from **Pending Deposit Subtree** to **Stored Commitment**

Body: `stored.DepositCommitmentBody`

```go
type DepositCommitmentBody struct {
    SubtreeID   Uint256
    SubtreeRoot common.Hash
    Deposits    []models.PendingDeposit
}
```


### Stored Batch

- Hold details of both mined and pending batches. Pending batches have some fields left `nil`.

Key: batch ID `models.Uint256`

Value: `stored.Batch`

```go
type Batch struct {
    ID                Uint256
    BType             batchtype.BatchType // not named `Type` to avoid collision with Type() method needed to implement bh.Storer interface
    TransactionHash   common.Hash
    Hash              *common.Hash // root of merkle tree of all commitments included in this batch
    FinalisationBlock *uint32
    AccountTreeRoot   *common.Hash
    PrevStateRoot     *common.Hash
    MinedTime         *models.Timestamp
}
```

#### Index on `Hash`
Key: batch hash `*common.Hash`

Prefix: `_bhIndex:Batch:Hash:`

Value: list of batch IDs `bh.KeyList`

## Other

### Chain State

- Holds a single value: the current chain state

Key: "ChainState" `string`

Value: `models.ChainState`

```go
type ChainState struct {
    ChainID                        Uint256
    AccountRegistry                common.Address
    AccountRegistryDeploymentBlock uint64
    TokenRegistry                  common.Address
    DepositManager                 common.Address
    Rollup                         common.Address
    GenesisAccounts                []PopulatedGenesisAccount
}

type PopulatedGenesisAccount struct {
    PublicKey [128]byte
    PubKeyID  uint32
    StateID   uint32
    Balance   Uint256
}
```

### Registered Token

- Holds tokenID - token contract address mapping

Key: token ID `models.Uint256`

Value: token contract address `common.Address` (through clever encoding of `models.RegisteredToken`)

```go
type RegisteredToken struct {
    ID       models.Uint256
    Contract common.Address
}
```

### Registered Spoke

- Holds spokeID - spoke contract address mapping

Key: spoke ID `models.Uint256`

Value: spoke contract address `common.Address` (through clever encoding of `models.RegisteredSpoke`)

```go
type RegisteredSpoke struct {
    ID       models.Uint256
    Contract common.Address
}
```

### Pending Stake Withdrawal

- Holds batchID - finalisation block mapping

Key: batchID `models.Uint256`

Value: finalisation block `uint32`

```go
type PendingStakeWithdrawal struct {
	BatchID           Uint256
	FinalisationBlock uint32
}
```
