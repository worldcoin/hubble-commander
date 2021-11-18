# 🦡 Badger data structures

## Stored structures
- State Tree
    - [State Leaf](#State-Leaf)
        - [Index on `PubKeyID`](#Index-on-PubKeyID)
    - [State Tree Node](#State-Tree-Node)
    - [State Update](#State-Update)
- Account Tree
    - [Account Leaf](#Account-Leaf)
        - [Index on `PublicKey`](#Index-on-PublicKey)
    - [Account Tree Node](#Account-Tree-Node)
- Transactions
    - [Stored Transaction](#Stored-Transaction)
        - [Transfer](#Transfer)
            - [Index on `ToStateID`](#Index-on-ToStateID)
        - [Create2Transfer](#Create2Transfer)
        - [MassMigration](#MassMigration)
        - [Index on `FromStateID`](#Index-on-FromStateID)
    - [Stored Transaction Receipt](#Stored-Transaction-Receipt)
        - [Index on `CommitmentID`](#Index-on-CommitmentID)
        - [Index on `ToStateID`](#Index-on-ToStateID-1)
- Deposits
    - [Pending Deposit](#Pending-Deposit)
    - [Pending Deposit SubTree](#Pending-Deposit-SubTree)
- Commitments and Batches
    - [Stored Commitment](#Stored-Commitment)
        - [Transaction Commitment](#Transaction-Commitment)
        - [Deposit Commitment](#Deposit-Commitment)
    - [Batch](#Batch)
        - [Index on `Hash`](#Index-on-Hash)
- Other
    - [Chain State](#Chain-State)
    - [Registered Token](#Registered-Token)

## Notes and design rationale
- Some indices are specified using badgerhold tags on struct fields.
  Others by implementing the `Indexes()` method of `bh.Storer` interface.
- Badger does not support indices on fields of pointer types well.
  By default, it would add IDs of all structs that have the indexed field set to `nil` to a single index entry, for instance: 

  `_bhIndex:StoredTxReceipt:ToStateID:nil -> bh.KeyList{ txHash1, txHash2, ... }`  

  Such index entry can quickly grow in size.
  Thus, for structs that have indices on fields of pointer type we implement `Indexes()` method and specify our own `IndexFunc`.
  Returning `nil` from such `IndexFunc` for `nil` field values prevents creation of the `nil` index entry.

## State Tree

### State Leaf

- Holds `UserState` data

Key: state ID `uint32`

Value: `models.FlatStateLeaf`

```go
type FlatStateLeaf struct {
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

### Stored Transaction
- Stores pending and mined transactions data

Key: tx hash `common.Hash`

Value: `models.StoredTx`

```go
type StoredTx struct {
    Hash        common.Hash
    TxType      txtype.TransactionType
    FromStateID uint32
    Amount      Uint256
    Fee         Uint256
    Nonce       Uint256
    Signature   Signature
    ReceiveTime *Timestamp

    Body TxBody // interface
}
```

#### Transfer
Body: `models.StoredTxTransferBody`

```go
type StoredTxTransferBody struct {
    ToStateID uint32
}
```

##### Index on `ToStateID`

- This index is updated only when storing transactions with `StoredTxTransferBody`

Key: to state ID `uint32`

Prefix: `_bhIndex:StoredTx:ToStateID:`

Value: list of tx hashes `bh.KeyList`

#### Create2Transfer
Body: `models.StoredTxCreate2TransferBody`

```go
type StoredTxCreate2TransferBody struct {
    ToPublicKey PublicKey
}
```

#### MassMigration
Body: `models.StoredTxMassMigrationBody`

```go
type StoredTxMassMigrationBody struct {
    SpokeID Uint256
}
```

#### Index on `FromStateID`
- This index is updated for all stored transactions
Key: from state ID `uint32`

Prefix: `_bhIndex:StoredTx:FromStateID:`

Value: list of tx hashes `bh.KeyList`

### Stored Transaction Receipt
- Stores transactions details known only after it is mined

Key: tx hash `common.Hash`

Value: `models.StoredTxReceipt`

```go
type StoredTxReceipt struct {
    Hash         common.Hash
    CommitmentID *CommitmentID
    ToStateID    *uint32 // specified for C2Ts, nil for Transfers and MassMigrations
    ErrorMessage *string
}
```

#### Index on `CommitmentID`
Key: commitment ID `models.CommitmentID`

Prefix: `_bhIndex:StoredTxReceipt:CommitmentID:`

Value: list of tx hashes `bh.KeyList`

#### Index on `ToStateID`
Key: to state ID `uint32`

Prefix: `_bhIndex:StoredTxReceipt:ToStateID:`

Value: list of tx hashes `bh.KeyList`

## Deposits

### Pending Deposit

- Holds individual pending deposits until they are moved into **Pending Deposit SubTrees**. Individual deposits correspond to `DepositQueued` events emitted by the SCs.  

Key: <blockNumber, logIndex> `models.DepositID`

Value: `models.PendingDeposit`

```go
type DepositID struct {
    BlockNumber uint32 // block in which the deposit tx was included
    LogIndex    uint32 // `DepositQueued` log index in that block
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

### Pending Deposit SubTree

- FIFO queue for ready Deposit SubTrees that can be submitted in a new Batch

Key: SubTreeID `models.Uint256`

Value: `models.PendingDepositSubTree`

```go
type PendingDepositSubTree struct {
    ID       models.Uint256          // assigned in SC
    Root     common.Hash             // subtree root that will be inserted into the state tree
    Deposits []models.PendingDeposit // deposits included in the subtree
}
```

## Commitments and Batches

### Stored Commitment
- Hold commitment details for both transaction and deposit commitmetns


Key: {BatchID, IndexInBatch} `models.CommitmentID`

Value: `models.StoredCommitment`

```go
type StoredCommitment struct {
    CommitmentBase
    Body StoredCommitmentBody // interface
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

#### Transaction Commitment

Body: `models.StoredTxCommitmentBody`

```go
type StoredTxCommitmentBody struct {
    FeeReceiver       uint32
    CombinedSignature [64]byte
    BodyHash          *common.Hash
}
```

#### Deposit Commitment
- When Deposit batch is created data is moved from **Pending Deposit SubTree** to **Stored Commitment**

Body: `models.StoredDepositCommitmentBody`

```go
type StoredDepositCommitmentBody struct {
    SubTreeID   Uint256
    SubTreeRoot common.Hash
    Deposits    []PendingDeposit
}
```


### Batch

- Hold details of both mined and pending batches. Pending batches have some fields left `nil`.

Key: BatchID `models.Uint256`

Value: `models.Batch`

```go
type Batch struct {
    ID                Uint256
    Type              batchtype.BatchType
    TransactionHash   common.Hash
    Hash              *common.Hash // root of merkle tree of all commitments included in this batch
    FinalisationBlock *uint32
    AccountTreeRoot   *common.Hash
    PrevStateRoot     *common.Hash
    SubmissionTime    *Timestamp
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
    SyncedBlock                    uint64
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
