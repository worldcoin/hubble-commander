# 🛰 JSON RPC

The default namespace would be `hubble`. We could also have other ones like `admin`.

### `hubble_getVersion()`

Returns the version number

Example result:

```json
"dev-0.0.1"
```

### `hubble_getNetworkInfo()`

This returns a number of datapoints about the current state of the system:

- ethereum network chain id
- AccountRegistry, TokenRegistry and Rollup contract addresses
- Block at which contracts were deployed (for new instance of commander to know where to start syncing events from)
- Current ethereum block number
- Number of transactions and accounts
- ID of the latest batch
- ID of the latest finalised batch
- Domain used for signatures

Example result:

```json
{
    "ChainID": "1616067554748",
    "AccountRegistry": "0x61e655cac5ac4f224501a8e578852c1919c119c5",
		"AccountRegistryDeploymentBlock": 20,
		"DepositManager": "0x2d72f696d6b1fecdfa1e55f5f8778953509e5329",
    "Rollup": "0x568fcd21edc59900314a944175805e9fc61d1bdd",
    "BlockNumber": 48,
		"TransactionCount": 65,
		"AccountCount": 38,
    "LatestBatch": "12",
    "LatestFinalisedBatch": "10",
    "SignatureDomain": "0x123123abc..."
}
```

### `hubble_getGenesisAccounts()`

Returns a list of genesis accounts added at batch #0. Example result:

```json
[
    {
        "PublicKey": "0x0df68cb87856229b0bc3f158fff8b82b04deb1a4c23dadbf3ed2da4ec6f6efcb1c165c6b47d8c89ab2ddb0831c182237b27a4b3d9701775ad6c180303f87ef260566cb2f0bcc7b89c2260de2fee8ec29d7b5e575a1e36eb4bcead52a74a511b7188d7df7c9d08f94b9daa9d89105fbdf22bf14e30b84f8adefb3695ebff00e88",
        "PubKeyID": 0,
        "StateID": 0,
        "Balance": "1000000000000000000"
    },
    {
        "PublicKey": "0x0097f465fe827ce4dad751988f6ce5ec747458075992180ca11b0776b9ea3a910c3ee4dca4a03d06c3863778affe91ce38d502138356a35ae12695c565b24ea6151b83eabd41a6090b8ac3bb25e173c84c3b080a5545260b1327495920c342c02d51cac4418228db1a3d98aa12e6fd7b3267c703475f5999b2ec7a197ad7d8bc",
        "PubKeyID": 1,
        "StateID": 1,
        "Balance": "1000000000000000000"
    }
]
```

### `hubble_sendTransaction(IncomingTransaction)`

Adds a transaction to the pending list of transactions. The transaction needs a valid signature and nonce.

Example result:

```json
"0x9b442316136f46247a399169aff5b9931060331f4b66971766a81b77765cfb36"
```

### `hubble_getTransaction(Hash)`

Returns transaction object including its status:

- PENDING
- IN_BATCH
- FINALISED
- ERROR

Example result (TRANSFER):

```json
{
    "Hash": "0x9b442316136f46247a399169aff5b9931060331f4b66971766a81b77765cfb36",
    "TxType": "TRANSFER",
    "FromStateID": 1,
    "Amount": "50",
    "Fee": "1",
    "Nonce": "0",
    "Signature": "0x19721d261934684e37c582a1cbd69eb406eb430e56c1865e0bd071d350728c5225fc666aeb79165b6b29616a1018aabf0d317d60b0feb5668e372d7810081d96",
    "IncludedInCommitment": null,
    "ErrorMessage": null,
    "ToStateID": 2,
    "ReceiveTime": 1625153276, // timestamp at which the tx was received by the coordinator for inclusion in batch. Can be null, when the tx was synced from blockchain.
		"BatchTime": 1645153276, // timestamp at which the tx was included in a batch submitted on chain. Can be null, when the tx hasn't been included yet.
    "Status": "PENDING"
}
```

Example result (CREATE2TRANSFER):

```json
{
    "Hash": "0x6cbc0e7428308e2f0397e5f35d6f5eb8922cd67bdd2eda39d55bfd379c2f2f1a",
    "TxType": "CREATE2TRANSFER",
    "FromStateID": 2,
    "Amount": "50",
    "Fee": "1",
    "Nonce": "0",
    "Signature": "0x06e788fb9494da058be2c6d3e078d4833401fcf992f410ffe98ff88a02b833d828f399d33b5368ccd3f03c6d6d1aa6b7b36f398210abf056434b4b15b3462c59",
    "IncludedInCommitment": null,
    "ErrorMessage": null,
    "ToStateID": 0,
    "ToPublicKey": "0x0097f465fe827ce4dad751988f6ce5ec747458075992180ca11b0776b9ea3a910c3ee4dca4a03d06c3863778affe91ce38d502138356a35ae12695c565b24ea6151b83eabd41a6090b8ac3bb25e173c84c3b080a5545260b1327495920c342c02d51cac4418228db1a3d98aa12e6fd7b3267c703475f5999b2ec7a197ad7d8bc",
    "ReceiveTime": 1625153276, // timestamp at which the tx was received by the coordinator for inclusion in batch. Can be null, when the tx was synced from blockchain.
		"BatchTime": 1645153276, // timestamp at which the tx was included in a batch submitted on chain. Can be null, when the tx hasn't been included yet.
    "Status": "PENDING"
}
```

### `hubble_getTransactions(pubKey)`

Returns an array of transactions (TRANSFER and CREATE2TRANSFER type) for given public key

### `hubble_getUserState(stateId)`

```json
{
      "StateID": 3,
      "PubKeyID": 3,
      "TokenIndex": "0",
      "Balance": "1000000000000000000",
      "Nonce": "0"
  }
```

### `hubble_getUserStates(pubKey)`

Return all UserState objects related to a public key

```json
[
    {
        "StateID": 3,
        "PubKeyID": 3,
        "TokenIndex": "0",
        "Balance": "1000000000000000000",
        "Nonce": "0"
    },
    {
        "StateID": 1,
        "PubKeyID": 1,
        "TokenIndex": "0",
        "Balance": "999999999999996800",
        "Nonce": "32"
    }
]
```

### `hubble_getPublicKeyByID(pubkeyId)`

```json
{
    "ID": 34,
    "PublicKey": "0x123123123"
}
```

### `hubble_getPublicKeyByStateID(stateID)`

```json
{
    "ID": 34,
    "PublicKey": "0x123123123"
}
```

### `hubble_getBatches(from, to)`

Returns an array of batches in given ID range

Example result:

```json
[
   {
	    "ID": "1",
	    "Hash": "0xb1786be90de852376032956f0ed165011bf4ef3a6f6a0753a1f9f1f59e99441f",
	    "Type": "TRANSFER",
	    "TransactionHash": "0xa680c859cbe25b6a26f0cd017318fe51d7a195d32cd004f6a8862dbe0a348f94",
	    "SubmissionBlock": 59,
			"SubmissionTime": 1645153276,
	    "FinalisationBlock": 40379
   }
]
```

### `hubble_getBatchByHash(hash)`

Returns batch information and list of included commitments in batch

Example result:

```json
{
		"ID": "1",
    "Hash": "0x05dc2502909e415e5f6c6b3d71e94192c277274816f65cd4fba000fda4706372",
    "Type": "TRANSFER",
		"TransactionHash": "0xa680c859cbe25b6a26f0cd017318fe51d7a195d32cd004f6a8862dbe0a348f94"
    "FinalisationBlock": 40363,
    "SubmissionBlock": 268,
		"SubmissionTime": 1645153276,
    "AccountTreeRoot": "0xce31f2b7478a42b23e5661eda2c7665488e0b51b55e24e60a129110b28db4d68",
    "Commitments": [
        {
            "ID": 1,
            "LeafHash": "0x8d394d5d9f1223e729571f8727b46a6e3dd8e4799bb5dfcdd7d44a1a7391c126",
            "TokenID": "0",
            "FeeReceiverStateID": 0,
            "CombinedSignature": "0x1b111c3fdeddc3326a326fb180c8a825b8919959dfbcccf5cf51cb76fb2d9a9c066bcbabb65792f9c9b7cfce8cc484de8e0a06c2703b023d69fa93f8bdfc8f44",
            "PostStateRoot": "0xfb2f86a3c7d96a560f88f3e22522195a5c96fd5df74690ecd0256c4a3c4a0836"
        }
    ]
}
```

### `hubble_getBatchByID(id)`

Same as hubble_getBatchByHash(hash)

### `hubble_getCommitment(id)`

Returns commitment information and list of included transactions

Example result:

```json
{
    "ID": 1,
    "Type": "TRANSFER",
    "FeeReceiver": 0,
    "CombinedSignature": "0x1b111c3fdeddc3326a326fb180c8a825b8919959dfbcccf5cf51cb76fb2d9a9c066bcbabb65792f9c9b7cfce8cc484de8e0a06c2703b023d69fa93f8bdfc8f44",
    "PostStateRoot": "0xfb2f86a3c7d96a560f88f3e22522195a5c96fd5df74690ecd0256c4a3c4a0836",
    "AccountTreeRoot": "0xce31f2b7478a42b23e5661eda2c7665488e0b51b55e24e60a129110b28db4d68",
    "IncludedInBatch": "0x05dc2502909e415e5f6c6b3d71e94192c277274816f65cd4fba000fda4706372",
    "Status": "IN_BATCH",
		"BatchTime": 1645153276,
    "Transactions": [
        {
            "Hash": "0x4cabb7314775aae932fb3995315298b5241f1793753fd4fd4b5c5405f5dc18e6",
            "FromStateID": 1,
            "Amount": "90",
            "Fee": "10",
            "Nonce": "0",
            "Signature": "0x28c71cc24191b4fd335bc5b2045d27e723820dc1071ad20882f2ec4347cccfff302d354a7adb570b1a53aebd6a59537271648391eb20e72c8fe4d83a2e6e4df6",
            "ToStateID": 2
            "ReceiveTime": 1625153276 // timestamp at which the tx was received by the coordinator for inclusion in batch. Can be null, when the tx was synced from blockchain.
        }
    ]
}
```

# API usage

## Sending a transaction

1. Call `hubble_getUserStates(senderPubkey)` to list sender's accounts. Pick one with an appropriate token index and balance.

1. Call `hubble_getUserStates(recipientPubkey)` to list recipients accounts pick one with an appropriate token index.
2. Call `hubble_sendTransaction` with `txType=TRANSFER` using the state indexes from steps 1 & 2 and nonce from step 1.
3. Call `hubble_getTransfer(Hash)` with the hash from step 3 to monitor transaction progress.

## Alternative: Recipient doesn't have a state leaf for a given token

1. Call `hubble_getUserStates(senderPubkey)` to list sender's accounts. Pick one with an appropriate token index and balance.
2. Call `hubble_sendTransaction` with `txType=CREATE2TRANSFER` with the recipient's pubkey. The commander will decide the appropriate account index and state index for the recipient.
3. Call `hubble_getTransaction(Hash)` with the hash from step 2 to monitor transaction progress.

DRAFT:

### `hubble_getBatchByNumber(batchNumber)`

Returns the information about the batch (hash, id) as well as the included commitments. Should this include transactions? Should hash and id be the same thing?

### `hubble_getAccounts(PublicKey)`

Returns the indices of an account in the account tree. Used for constructing transactions.

### `hubble_getTokenIndex(Token)`

Returns the index of the token. Figure out if the tokens actually do have indices. Used for constructing transactions.

### `hubble_getBalanceByAddress(Address, Token)`

Returns the user's balance in the specified token.
Potentially in the future extend this with the ability to see the current balance, committed balance and finalized (non-disputable) balance.

### `hubble_getTransactionCount(Address)`

Returns the user's transaction count. Used for setting the nonce. Should this be per token?

### `hubble_getTransactions(Address)`

Returns the user's list of transactions. Should this be per token? Should this have pagination?

### `hubble_getLatestCommitment()`

Returns the latest commitment, see below.

### `hubble_getCommitment(id)`

Returns the information about the commitment (hash, id) as well as the included transactions. Should hash and id be the same thing?

### `hubble_getLatestBatch()`

Returns the latest batch, see below.

```json
{
    "ID": "1",
    "Hash": "0xb1786be90de852376032956f0ed165011bf4ef3a6f6a0753a1f9f1f59e99441f",
    "Type": "TRANSFER",
    "TransactionHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "SubmissionBlock": 59,
    "FinalisationBlock": 40379
},
```
