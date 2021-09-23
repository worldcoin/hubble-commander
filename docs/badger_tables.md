# 🦡 Badger tables

## Current tables

### State Leaf

- Holds `UserState` data

Key: stateID `uint32`

Value: `models.FlatStateLeaf`

### Merkle Node

- Holds state tree nodes (hashes). The leaf level nodes are hashes of `UserState` structs.

Key: node path `models.MerklePath`

Value: node `common.Hash` (through clever encoding of `models.StateNode`)

### State Update

- Holds state tree updates including leaves data that were previously in the tree

Key: auto-incremented ID `uint64`

Value: `models.StateUpdate`

### Account Leaf

- Holds pubKeyID - publicKey mapping

Key: pubKeyID `uint32`

Value: {pubKeyID, publicKey} `models.Account`

#### PublicKey index

- Index on `PublicKey` field makes badger auto-generate and update an index table

Key: publicKey `models.PublicKey`

Value: list of pubKeyIDs `[]uint32 ?`

### Chain State

- Holds a single value: the current chain state

Key: "ChainState" `string`

Value: `models.ChainState`

### Commitment

- Hold commitment details along with serialized transactions data

Key: {BatchID, IndexInBatch} `models.CommitmentID`

Value: `models.Commitment`

### Batch

- Hold details of both mined and pending batches. Pending batches have some fields left `nil`.

Key: BatchID `models.Uint256`

Value: `models.Batch`

## Migrating Postgres tables to Badger

- Draft docs for implementation purposes

### Chain State

```go
"ChainState" -> {
	ChainID         Uint256        
	AccountRegistry common.Address 
	DeploymentBlock uint64         
	Rollup          common.Address
	GenesisAccounts GenesisAccounts 
	SyncedBlock     uint64
}
```

### Batches

```go
BatchID -> {
  ID                Uint256
	Type              txtype.TransactionType
	TransactionHash   common.Hash  // submission tx
	Hash              *common.Hash // root of tree containing all commitments included in this batch
	FinalisationBlock *uint32      
	AccountTreeRoot   *common.Hash 
	PrevStateRoot     *common.Hash 
	SubmissionTime    *Timestamp   
}
```

### Commitments

```go
<BatchID, IndexInBatch> -> {
	BatchID           Uint256 // keep it in struct, just not serialize it
	IndexInBatch      int32   // keep it in struct, just not serialize it  
	Type              txtype.TransactionType
	Transactions      []byte
	FeeReceiver       uint32
	CombinedSignature Signature
	PostStateRoot     common.Hash
}
```

### Transactions

```go
Hash -> {
	Hash                 common.Hash // keep it in struct, just not serialize it         
	TxType               txtype.TransactionType 
	FromStateID          uint32                 
	Amount               Uint256
	Fee                  Uint256
	Nonce                Uint256
	Signature            Signature
	ReceiveTime          *Timestamp 
	IncludedInCommitment *{BatchID: Uint256, IndexInBatch: uint8}
	ErrorMessage         *string

  Body TransactionBody
}
```

## Tables and changes required for Deposits

- Draft docs for implementation purposes

### Token Registry

- Holds tokenID - token contract address mapping

Key: tokenID `models.Uint256`

Value: tokenContract `common.Address` (through clever encoding of `models.RegisteredToken`)

```go
type RegisteredToken struct {
	ID       models.Uint256
	Contract common.Address
}
```

### Pending Deposits

- Holds deposits until they are moved into **Pending Deposit SubTrees**

Key: <blockNumber, logIndex> `models.PendingDepositID`

Value: `models.PendingDeposit`

```go
type PendingDepositID struct {
	BlockNumber uint32 // block in which the deposit tx was included
	LogIndex    uint32 // `DepositQueued` log index in that block
}
```

```go
type PendingDeposit struct {
	ID                   DepositID
	ToPubKeyID           uint32
  TokenID              models.Uint256
  L2Amount             models.Uint256
}
```

### Pending Deposit SubTrees

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

### Commitment → TransactionCommitment / DepositCommitment

- The current Commitment table will need to be updated to support both transaction and deposit commitments

Stored structures (objects persisted to Badger):

```go
type StoredCommitment struct {
	ID            models.CommitmentID // keep it in struct, just not serialize it
	Type          batchtype.BatchType // type added in PR #337
	PostStateRoot common.Hash

  Body models.StoredCommitmentBody // interface type, different body depending on `Type`
}

type StoredCommitmentTxBody struct {
	FeeReceiver       uint32
	CombinedSignature models.Signature
	Transactions      []byte
}

type StoredCommitmentDepositBody struct {
	SubTreeID         models.Uint256
	SubTreeRoot       common.Hash
	Deposits          []models.Deposit
}

type Deposit struct {
	ID                   DepositID
	ToPubKeyID           uint32
	TokenID              models.Uint256
	L2Amount             models.Uint256
}

```

Types for internal use:

```go
type CommitmentBase struct {
	ID                models.CommitmentID
	Type              txtype.TransactionType
	PostStateRoot     common.Hash
}

type TransactionCommitment struct {
	CommitmentBase
	FeeReceiver       uint32
	CombinedSignature models.Signature
	Transactions      []byte
}

type DepositCommitment struct {
	CommitmentBase
	SubTreeID         models.Uint256
	SubTreeRoot       common.Hash
	Deposits          []models.Deposit	
}
```
