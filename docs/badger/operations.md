# Badger Operations

Statistics for large operations done on Badger database. This is mostly for tracking whether we don't hit Badger txn limits. Stats generated
using [these changes to Badger](https://github.com/msieczko/badger/commit/bf43a3a4b9dfb80019a97b10d0e7d269a7eff34e).

All tests run for full batches:

```
min_txs_per_commitment: 32
min_commitments_per_batch: 32
```

## Rollup loop

### Transfer batch

Measured with `TestBenchTransfersCommander`. Tx count: `10000`.

Both Badger tx operations count and tx size stable for consecutive batches.

```
Key: bh_MerkleTreeNode, Count: 68640, Size: 4942080
Key: _bhIndex:TxReceipt:CommitmentID, Count: 1024, Size: 847360
Key: bh_StateUpdate, Count: 2080, Size: 505440
Key: bh_StateLeaf, Count: 2080, Size: 343200
Key: _bhIndex:StateLeaf:PubKeyID, Count: 4160, Size: 235040
Key: bh_TxReceipt, Count: 1024, Size: 132096
Key: bh_Commitment, Count: 32, Size: 7232
Key: bh_Batch, Count: 1, Size: 238
SUM: Count: 79042, Size: 7012707
```

### Create2Transfer batch (realistic)

Measured with `TestBenchCreate2TransfersCommander`. Batch with 1024 public key registrations. Tx count: `10000`.

Badger tx operations count stable and tx size mostly stable for consecutive batches. Only `_bhIndex:AccountLeaf:PublicKey`
and `_bhIndex:StateLeaf:PubKeyID` index sizes grow slowly.

```
Key: bh_MerkleTreeNode, Count: 102432, Size: 7442688
Key: _bhIndex:TxReceipt:CommitmentID, Count: 1024, Size: 847360
Key: bh_StateUpdate, Count: 2080, Size: 505440
Key: bh_StateLeaf, Count: 2080, Size: 343200
Key: _bhIndex:AccountLeaf:PublicKey, Count: 1024, Size: 203436
Key: _bhIndex:StateLeaf:PubKeyID, Count: 3136, Size: 190035
Key: bh_AccountLeaf, Count: 1024, Size: 166912
Key: bh_TxReceipt, Count: 1024, Size: 132096
Key: _bhIndex:TxReceipt:ToStateID, Count: 1024, Size: 100352
Key: bh_Commitment, Count: 32, Size: 7232
Key: bh_Batch, Count: 1, Size: 238
SUM: Count: 114882, Size: 9939010
```

### Create2Transfer batch (unrealistic)

Measured with `TestBenchCreate2TransfersCommander` (
at [`6b813227`](https://github.com/worldcoin/hubble-commander/commit/6b81322780bb73f21ce25c434265062fc72a44bd)).

This test sends an enormous number of C2Ts to a bunch of registered accounts. A lot of new User States are created for a small set of pub
key IDs. As a result `_bhIndex:FlatStateLeaf:PubKeyID` index grows with every batch.

First batch:

```
Key: bh_MerkleTreeNode, Count: 68640, Size: 4942080
Key: _bhIndex:FlatStateLeaf:PubKeyID, Count: 3136, Size: 980392
Key: _bhIndex:TxReceipt:CommitmentID, Count: 1024, Size: 954880
Key: bh_StateUpdate, Count: 2080, Size: 505440
Key: bh_FlatStateLeaf, Count: 2080, Size: 351520
Key: bh_TxReceipt, Count: 1024, Size: 138240
Key: _bhIndex:TxReceipt:ToStateID, Count: 1024, Size: 112640
Key: bh_Commitment, Count: 32, Size: 7424
Key: bh_Batch, Count: 1, Size: 244
SUM: Count: 79042, Size: 7992881
```

Second batch:

```
Key: bh_MerkleTreeNode, Count: 68640, Size: 4942080
Key: _bhIndex:FlatStateLeaf:PubKeyID, Count: 3136, Size: 2549209
Key: _bhIndex:TxReceipt:CommitmentID, Count: 1024, Size: 954880
Key: bh_StateUpdate, Count: 2080, Size: 505440
Key: bh_FlatStateLeaf, Count: 2080, Size: 351520
Key: bh_TxReceipt, Count: 1024, Size: 138240
Key: _bhIndex:TxReceipt:ToStateID, Count: 1024, Size: 112640
Key: bh_Commitment, Count: 32, Size: 7424
Key: bh_Batch, Count: 1, Size: 244
SUM: Count: 79042, Size: 9561698
```

Third batch:

```
Key: bh_MerkleTreeNode, Count: 68640, Size: 4942080
Key: _bhIndex:FlatStateLeaf:PubKeyID, Count: 3136, Size: 4241675
Key: _bhIndex:TxReceipt:CommitmentID, Count: 1024, Size: 954880
Key: bh_StateUpdate, Count: 2080, Size: 505440
Key: bh_FlatStateLeaf, Count: 2080, Size: 351520
Key: bh_TxReceipt, Count: 1024, Size: 138240
Key: _bhIndex:TxReceipt:ToStateID, Count: 1024, Size: 112640
Key: bh_Commitment, Count: 32, Size: 7424
Key: bh_Batch, Count: 1, Size: 244
SUM: Count: 79042, Size: 11254164
```

## Syncing

### Transfer batch

Measured with `TestBenchSyncCommander` set to send and sync only Transfer batches. Tx count: `10000`.

The test sends an enormous number of Transfers using just a bunch of User States. Removing `_bhIndex:Tx:FromStateID`
and `_bhIndex:Tx:ToStateID` indices made the size of txs stable for consecutive batches.

```
Key: bh_MerkleTreeNode, Count: 68640, Size: 4942080
Key: _bhIndex:TxReceipt:CommitmentID, Count: 1024, Size: 847360
Key: bh_StateUpdate, Count: 2080, Size: 505440
Key: bh_StateLeaf, Count: 2080, Size: 343200
Key: bh_Tx, Count: 1024, Size: 273408
Key: _bhIndex:StateLeaf:PubKeyID, Count: 4160, Size: 235040
Key: bh_TxReceipt, Count: 1024, Size: 132096
Key: bh_Commitment, Count: 32, Size: 7232
Key: bh_Batch, Count: 1, Size: 238
Key: _bhIndex:Batch:Hash, Count: 1, Size: 113
SUM: Count: 80067, Size: 7286228
```

### Create2Transfer batch

Measured with `TestBenchSyncCommander` set to send and sync only Create2Transfer batches. Tx count: `5000`.

Badger tx size stable for consecutive batches.

```
Key: bh_MerkleTreeNode, Count: 68640, Size: 4942080
Key: _bhIndex:TxReceipt:CommitmentID, Count: 1024, Size: 847360
Key: bh_StateUpdate, Count: 2080, Size: 505440
Key: bh_Tx, Count: 1024, Size: 400384
Key: bh_StateLeaf, Count: 2080, Size: 343200
Key: _bhIndex:StateLeaf:PubKeyID, Count: 3136, Size: 190069
Key: bh_TxReceipt, Count: 1024, Size: 132096
Key: _bhIndex:TxReceipt:ToStateID, Count: 1024, Size: 100352
Key: bh_Commitment, Count: 32, Size: 7232
Key: bh_Batch, Count: 1, Size: 238
Key: _bhIndex:Batch:Hash, Count: 1, Size: 113
SUM: Count: 80067, Size: 7468585
```
