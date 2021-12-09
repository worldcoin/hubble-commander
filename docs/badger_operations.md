# Badger Operations

Statistics for large operations done on Badger database.
This is mostly for tracking whether we don't hit Badger txn limits. 

## Rollup loop

### Transfer batch
```ignore
min_txs_per_commitment: 32
min_commitments_per_batch: 32
```

```ignore
Key: bh_MerkleTreeNode, Count: 68640, Size: 4942080
Key: _bhIndex:TxReceipt:CommitmentID, Count: 1024, Size: 954880
Key: bh_StateUpdate, Count: 2080, Size: 505440
Key: bh_StateLeaf, Count: 2080, Size: 351520
Key: _bhIndex:FlatStateLeaf:PubKeyID, Count: 4160, Size: 260000
Key: bh_TxReceipt, Count: 1024, Size: 138240
Key: bh_Commitment, Count: 32, Size: 7424
Key: bh_Batch, Count: 1, Size: 244
SUM: Count: 79042, Size: 7159849
```

### Create2Transfer batch (realistic)
Measured with `TestBenchPubKeysRegistration`. Batch with 1024 public key registrations.
```ignore
min_txs_per_commitment: 32
min_commitments_per_batch: 32
```
```ignore
Key: bh_MerkleTreeNode, Count: 102432, Size: 7442688
Key: _bhIndex:TxReceipt:CommitmentID, Count: 1024, Size: 954880
Key: bh_StateUpdate, Count: 2080, Size: 505440
Key: bh_StateLeaf, Count: 2080, Size: 351520
Key: _bhIndex:StateLeaf:PubKeyID, Count: 3136, Size: 210932
Key: _bhIndex:AccountLeaf:PublicKey, Count: 1024, Size: 203094
Key: bh_AccountLeaf, Count: 1024, Size: 166912
Key: bh_TxReceipt, Count: 1024, Size: 138240
Key: _bhIndex:TxReceipt:ToStateID, Count: 1024, Size: 112640
Key: bh_Commitment, Count: 32, Size: 7424
Key: bh_Batch, Count: 1, Size: 244
SUM: Count: 114882, Size: 10094035
```

### Create2Transfer batch (unrealistic)
Measured with `TestBenchCreate2TransfersCommander`. 
```ignore
min_txs_per_commitment: 32
min_commitments_per_batch: 32
```


This test sends an enormous number of C2Ts to a bunch of registered accounts. 
A lot of new User States are created for a small set of pub key IDs.
As a result `_bhIndex:FlatStateLeaf:PubKeyID` index grows with every batch.


First batch:
```ignore
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
```ignore
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
```ignore
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
Measured with `TestBenchSyncCommander` set to send and sync only Transfer batches.
```ignore
min_txs_per_commitment: 32
min_commitments_per_batch: 32
```
The test sends an enormous number of Transfers using just a bunch of User States. 
As a result `_bhIndex:Tx:FromStateID` and `_bhIndex:Tx:ToStateID` grow with every batch.

First batch:
```ignore
Key: bh_MerkleTreeNode, Count: 68640, Size: 4942080
Key: _bhIndex:Tx:ToStateID, Count: 1024, Size: 3938404
Key: _bhIndex:Tx:FromStateID, Count: 1024, Size: 3922632
Key: _bhIndex:TxReceipt:CommitmentID, Count: 1024, Size: 954880
Key: bh_StateUpdate, Count: 2080, Size: 505440
Key: bh_FlatStateLeaf, Count: 2080, Size: 351520
Key: bh_Tx, Count: 1024, Size: 279552
Key: _bhIndex:FlatStateLeaf:PubKeyID, Count: 4160, Size: 260000
Key: bh_TxReceipt, Count: 1024, Size: 138240
Key: bh_Commitment, Count: 32, Size: 7424
Key: bh_Batch, Count: 1, Size: 244
Key: _bhIndex:Batch:Hash, Count: 1, Size: 125
SUM: Count: 82115, Size: 15300562
```
Second batch:
```ignore
Key: _bhIndex:Tx:ToStateID, Count: 1024, Size: 11643552
Key: _bhIndex:Tx:FromStateID, Count: 1024, Size: 11612160
Key: bh_MerkleTreeNode, Count: 68640, Size: 4942080
Key: _bhIndex:TxReceipt:CommitmentID, Count: 1024, Size: 954880
Key: bh_StateUpdate, Count: 2080, Size: 505440
Key: bh_FlatStateLeaf, Count: 2080, Size: 351520
Key: bh_Tx, Count: 1024, Size: 279552
Key: _bhIndex:FlatStateLeaf:PubKeyID, Count: 4160, Size: 260000
Key: bh_TxReceipt, Count: 1024, Size: 138240
Key: bh_Commitment, Count: 32, Size: 7424
Key: bh_Batch, Count: 1, Size: 244
Key: _bhIndex:Batch:Hash, Count: 1, Size: 125
SUM: Count: 82115, Size: 30695238
```
Third batch:
```ignore
Key: _bhIndex:Tx:FromStateID, Count: 1024, Size: 19301688
Key: _bhIndex:Tx:ToStateID, Count: 1024, Size: 19294712
Key: bh_MerkleTreeNode, Count: 68640, Size: 4942080
Key: _bhIndex:TxReceipt:CommitmentID, Count: 1024, Size: 954880
Key: bh_StateUpdate, Count: 2080, Size: 505440
Key: bh_FlatStateLeaf, Count: 2080, Size: 351520
Key: bh_Tx, Count: 1024, Size: 279552
Key: _bhIndex:FlatStateLeaf:PubKeyID, Count: 4160, Size: 260000
Key: bh_TxReceipt, Count: 1024, Size: 138240
Key: bh_Commitment, Count: 32, Size: 7424
Key: bh_Batch, Count: 1, Size: 244
Key: _bhIndex:Batch:Hash, Count: 1, Size: 125
SUM: Count: 82115, Size: 46035926
```
