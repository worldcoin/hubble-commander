# 🛰 JSON RPC

## Hubble API

### `hubble_getVersion()`

Returns the version number

Example result:

```json
"0.5.0"
```

### `hubble_getNetworkInfo()`

This returns a number of datapoints about the current state of the system:

- Ethereum network chain ID
- AccountRegistry, TokenRegistry, SpokeRegistry, DepositManager, WithdrawManager and Rollup contract addresses
- Block at which contracts were deployed (for new instance of commander to know where to start syncing events from)
- Current ethereum block number
- Number of transactions and accounts
- ID of the latest batch
- ID of the latest finalised batch
- Domain used for signatures

Example result:

```json
{
    "ChainID": "1337",
    "AccountRegistry": "0x10bd6732fe3908b8a816f6a1b271e0864de78ca1",
    "AccountRegistryDeploymentBlock": 74,
    "TokenRegistry": "0x07389715ae1f0a891fba82e65099f6a3fa7da593",
    "SpokeRegistry": "0x535ca2e7dc31afce3dde4d78ded91aedf55b04b8",
    "DepositManager": "0xa3accd1cfabc8b09aea4d0e25f21f25c526c9be8",
    "WithdrawManager": "0x7eaa005432a4602044ae2242c79234650304f290",
    "Rollup": "0xf2a409ccf78e6e32e02d5e3a3ac274ca6880d9ac",
    "BlockNumber": 2146,
    "TransactionCount": 2,
    "AccountCount": 6,
    "LatestBatch": "2",
    "LatestFinalisedBatch": "0",
    "SignatureDomain": "0x123123abc..."
}
```

### `hubble_getGenesisAccounts()`

Returns a list of genesis accounts added at batch #0.

Example result:

```json
[
    {
        "PublicKey": "0x0df68cb87856229b0bc3f158fff8b82b04deb1a4c23dadbf3ed2da4ec6f6efcb1c165c6b47d8c89ab2ddb0831c182237b27a4b3d9701775ad6c180303f87ef260566cb2f0bcc7b89c2260de2fee8ec29d7b5e575a1e36eb4bcead52a74a511b7188d7df7c9d08f94b9daa9d89105fbdf22bf14e30b84f8adefb3695ebff00e88",
        "StateID": 0,
        "State": {
            "PubKeyID": 0,
            "TokenID": "0",
            "Balance": "1000000000000000000",
            "Nonce": "0"
        }
    },
    {
        "PublicKey": "0x0097f465fe827ce4dad751988f6ce5ec747458075992180ca11b0776b9ea3a910c3ee4dca4a03d06c3863778affe91ce38d502138356a35ae12695c565b24ea6151b83eabd41a6090b8ac3bb25e173c84c3b080a5545260b1327495920c342c02d51cac4418228db1a3d98aa12e6fd7b3267c703475f5999b2ec7a197ad7d8bc",
        "StateID": 1,
        "State": {
            "PubKeyID": 1,
            "TokenID": "0",
            "Balance": "1000000000000000000",
            "Nonce": "0"
        }
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

- `PENDING`
- `SUBMITTED`
- `MINED`
- `FINALISED`
- `ERROR`

Example result (`TRANSFER`):

```json
{
    "Hash": "0x9b442316136f46247a399169aff5b9931060331f4b66971766a81b77765cfb36",
    "TxType": "TRANSFER",
    "FromStateID": 1,
    "Amount": "50",
    "Fee": "1",
    "Nonce": "0",
    "Signature": "0x19721d261934684e37c582a1cbd69eb406eb430e56c1865e0bd071d350728c5225fc666aeb79165b6b29616a1018aabf0d317d60b0feb5668e372d7810081d96",
    // timestamp at which the tx was received by the coordinator for inclusion in batch. Can be null, when the tx was synced from blockchain.
    "ReceiveTime": 1625153276,
    "CommitmentID": {
        "BatchID": "1",
        "IndexInBatch": 0
    },
    "ErrorMessage": null,
    "ToStateID": 2,
    "BatchHash": "0xeb590ba0ce14d821caebc56514fe867521da78b46b7b78ce4810a353e619f315",
    // timestamp at which the tx was included in a batch submitted on chain. Can be null, when the tx hasn't been included yet.
    "MinedTime": 1633692591,
    "Status": "FINALISED"
}
```

Example result (`CREATE2TRANSFER`):

```json
{
    "Hash": "0x6cbc0e7428308e2f0397e5f35d6f5eb8922cd67bdd2eda39d55bfd379c2f2f1a",
    "TxType": "CREATE2TRANSFER",
    "FromStateID": 1,
    "Amount": "50",
    "Fee": "1",
    "Nonce": "0",
    "Signature": "0x06e788fb9494da058be2c6d3e078d4833401fcf992f410ffe98ff88a02b833d828f399d33b5368ccd3f03c6d6d1aa6b7b36f398210abf056434b4b15b3462c59",
    // timestamp at which the tx was received by the coordinator for inclusion in batch. Can be null, when the tx was synced from blockchain.
    "ReceiveTime": 1625153276,
    "CommitmentID": {
        "BatchID": "1",
        "IndexInBatch": 0
    },
    "ErrorMessage": null,
    "ToStateID": 0,
    "ToPublicKey": "0x0097f465fe827ce4dad751988f6ce5ec747458075992180ca11b0776b9ea3a910c3ee4dca4a03d06c3863778affe91ce38d502138356a35ae12695c565b24ea6151b83eabd41a6090b8ac3bb25e173c84c3b080a5545260b1327495920c342c02d51cac4418228db1a3d98aa12e6fd7b3267c703475f5999b2ec7a197ad7d8bc",
    "BatchHash": "0xeb590ba0ce14d821caebc56514fe867521da78b46b7b78ce4810a353e619f315",
    // timestamp at which the tx was included in a batch submitted on chain. Can be null, when the tx hasn't been included yet.
    "MinedTime": 1633692591,
    "Status": "FINALISED"
}
```

Example result (`MASS_MIGRATION`):

```json
{
    "Hash": "0x5e40559e87254b436245fa9723f1489e4328869d8cf1a960066c74e74750c2ae",
    "TxType": "MASS_MIGRATION",
    "FromStateID": 1,
    "Amount": "50",
    "Fee": "1",
    "Nonce": "0",
    "Signature": "0x1f49179e0dc59a3d7ce425f7379d80d0d8ffac2a9d9d2a915d5f21fec8a59e3816b5ec10ce3bd608116c0c18e4bd07a8206e2f6a66696982f0ff441a357e8300",
    // timestamp at which the tx was received by the coordinator for inclusion in batch. Can be null, when the tx was synced from blockchain.
    "ReceiveTime": 1625153276,
    "CommitmentID": {
        "BatchID": "1",
        "IndexInBatch": 0
    },
    "ErrorMessage": null,
    "SpokeID": 2,
    "BatchHash": "0x9fc863e718defd9506764f4f623f2b0a63fa51dc1a2f85ca314846aaf5cf422c",
    // timestamp at which the tx was included in a batch submitted on chain. Can be null, when the tx hasn't been included yet.
    "MinedTime": 1633692591,
    "Status": "FINALISED"
}
```

### `hubble_getUserState(stateId)`

Example result:

```json
{
    "StateID": 3,
    "PubKeyID": 3,
    "TokenID": "0",
    "Balance": "1000000000000000000",
    "Nonce": "0"
}
```

### `hubble_getUserStates(pubKey)`

Return all UserState objects related to a public key.

Example result:

```json
[
    {
        "StateID": 3,
        "PubKeyID": 3,
        "TokenID": "0",
        "Balance": "1000000000000000000",
        "Nonce": "0"
    },
    {
        "StateID": 1,
        "PubKeyID": 1,
        "TokenID": "0",
        "Balance": "999999999999996800",
        "Nonce": "32"
    }
]
```

### `hubble_getPublicKeyByPubKeyID(pubKeyId)`

Example result:

```json
"0x0097f465fe827ce4dad751988f6ce5ec747458075992180ca11b0776b9ea3a910c3ee4dca4a03d06c3863778affe91ce38d502138356a35ae12695c565b24ea6151b83eabd41a6090b8ac3bb25e173c84c3b080a5545260b1327495920c342c02d51cac4418228db1a3d98aa12e6fd7b3267c703475f5999b2ec7a197ad7d8bc"
```

### `hubble_getPublicKeyByStateID(stateID)`

Example result:

```json
"0x0097f465fe827ce4dad751988f6ce5ec747458075992180ca11b0776b9ea3a910c3ee4dca4a03d06c3863778affe91ce38d502138356a35ae12695c565b24ea6151b83eabd41a6090b8ac3bb25e173c84c3b080a5545260b1327495920c342c02d51cac4418228db1a3d98aa12e6fd7b3267c703475f5999b2ec7a197ad7d8bc"
```

### `hubble_getBatches(from, to)`

Returns an array of batches with is statuses in given ID range. Batch statuses:

- `SUBMITTED`
- `MINED`
- `FINALISED`

Example result:

```json
[
    {
        "ID": "0",
        "Hash": "0x97222c9001ff4b2e8c328d15c02fc508259f2f9bec5af8b286a932813fe920d6",
        "Type": "GENESIS",
        "TransactionHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "MinedBlock": 44009,
        "MinedTime": null,
        "Status": "FINALISED",
        "FinalisationBlock": 44009
    },
    {
        "ID": "1",
        "Hash": "0xb286d51f69b2f2ff52d387663b449c78eb7f4e6aa86e81ac1efac5ef2020d505",
        "Type": "TRANSFER",
        "TransactionHash": "0x40aec9dddb89b73e09727a15d2f0bb2dd59cc9ae76da86b5b079f5c870df1033",
        "MinedBlock": 46332,
        "MinedTime": 1642591575,
        "Status": "MINED",
        "FinalisationBlock": 86652
    },
    {
        "ID": "2",
        "Hash": "0xedd8e1b9ef9e3279f80df1434bd3a0b46b073a95d9b8be5aa4a645b6ef0c7a25",
        "Type": "CREATE2TRANSFER",
        "TransactionHash": "0x9da0fcf69cbeec799f4a41a70431dc803c28ef32f80ec40b28b5f013e6749cfe",
        "MinedBlock": 46412,
        "MinedTime": 1642591655,
        "Status": "MINED",
        "FinalisationBlock": 86732
    },
    {
        "ID": "3",
        "Hash": "0x612d7e555b246b82a93ef412377631bc0edeb566a0d6fb57f23f14d705ce21e8",
        "Type": "MASS_MIGRATION",
        "TransactionHash": "0xd0b056121219000d69da10501470ea69163c35139abcf1b4f1774f7ca7b4a07b",
        "MinedBlock": 46398,
        "MinedTime": 1642591641,
        "Status": "MINED",
        "FinalisationBlock": 86718
    },
    {
        "ID": "4",
        "Hash": "0xe3df23dfda31fc9498bdf990f05f65ca47a18a3c1247709bf9ae634a27970677",
        "Type": "DEPOSIT",
        "TransactionHash": "0xcd03405862fdc876bd6f1ccd2b6a843835324021e09ad8f86ebee2de67144cc9",
        "MinedBlock": 44027,
        "MinedTime": 1642589195,
        "Status": "MINED",
        "FinalisationBlock": 84347
    },
    {
        "ID": "5",
        "Hash": null,
        "Type": "TRANSFER",
        "TransactionHash": "0xa428c43175e6156c933c436455a675c1e521e3841d9d44c6482ee26f861e49c1",
        "MinedBlock": null,
        "MinedTime": null,
        "Status": "SUBMITTED",
        "FinalisationBlock": null
    }
]
```

### `hubble_getBatchByHash(hash)`

Returns mined or finalised batch information and list of included commitments in batch.

Example result (`TRANSFER`):

```json
{
    "ID": "1",
    "Hash": "0xb286d51f69b2f2ff52d387663b449c78eb7f4e6aa86e81ac1efac5ef2020d505",
    "Type": "TRANSFER",
    "TransactionHash": "0x40aec9dddb89b73e09727a15d2f0bb2dd59cc9ae76da86b5b079f5c870df1033",
    "MinedBlock": 46332,
    "MinedTime": 1642591575,
    "Status": "MINED",
    "FinalisationBlock": 86652,
    "AccountTreeRoot": "0xb261c40259ad5dbaf32efb2256225bbf03dcda8e84cffdfe67e68b958e3c7a95",
    "Commitments": [
        {
            "ID": {
                "BatchID": "1",
                "IndexInBatch": 0
            },
            "PostStateRoot": "0x1be2feecb81d0338bf049a6e3324a254cd5b5a93f481d941a7f9e31dabc7cc9c",
            "LeafHash": "0xcb4df4741da627b1a21f85c793cde72c86d06b1402d968cab6ba6c2635348b01",
            "TokenID": "0",
            "FeeReceiverStateID": 0,
            "CombinedSignature": "0x046961ddc008eb0a8e91a9a73a82ab3197741be327b210b7a443b7c763af47970548c262ed1c6039163f075071a36dbf93a7c230d25bfc869d88ded50dd15d36"
        }
    ]
}
```

Example result (`CREATE2TRANSFER`):

```json
{
    "ID": "2",
    "Hash": "0xedd8e1b9ef9e3279f80df1434bd3a0b46b073a95d9b8be5aa4a645b6ef0c7a25",
    "Type": "CREATE2TRANSFER",
    "TransactionHash": "0x9da0fcf69cbeec799f4a41a70431dc803c28ef32f80ec40b28b5f013e6749cfe",
    "MinedBlock": 46412,
    "MinedTime": 1642591655,
    "Status": "MINED",
    "FinalisationBlock": 86732,
    "AccountTreeRoot": "0xb261c40259ad5dbaf32efb2256225bbf03dcda8e84cffdfe67e68b958e3c7a95",
    "Commitments": [
        {
            "ID": {
                "BatchID": "2",
                "IndexInBatch": 0
            },
            "PostStateRoot": "0x795a9e9e4bd52ac321c67eca39f5bb120f3c22181b0c825339df9be67a988ae6",
            "LeafHash": "0x03cfbdcb60a9372138d84a33738d3992b9d7c0e969f0ebcdad6f82259dfb9a3a",
            "TokenID": "0",
            "FeeReceiverStateID": 0,
            "CombinedSignature": "0x046961ddc008eb0a8e91a9a73a82ab3197741be327b210b7a443b7c763af47970548c262ed1c6039163f075071a36dbf93a7c230d25bfc869d88ded50dd15d36"
        }
    ]
}
```

Example result (`MASS_MIGRATION`):

```json
{
    "ID": "3",
    "Hash": "0x612d7e555b246b82a93ef412377631bc0edeb566a0d6fb57f23f14d705ce21e8",
    "Type": "MASS_MIGRATION",
    "TransactionHash": "0xd0b056121219000d69da10501470ea69163c35139abcf1b4f1774f7ca7b4a07b",
    "MinedBlock": 46398,
    "MinedTime": 1642591641,
    "Status": "MINED",
    "FinalisationBlock": 86718,
    "AccountTreeRoot": "0xb261c40259ad5dbaf32efb2256225bbf03dcda8e84cffdfe67e68b958e3c7a95",
    "Commitments": [
        {
            "ID": {
                "BatchID": "3",
                "IndexInBatch": 0
            },
            "PostStateRoot": "0xe2afe7bb8d8f0310cf68636777794b97327f3ce3633127b58c4f5f566cbf8ac7",
            "LeafHash": "0x44f03af79586fd8ab74b442be1bb920e40eefd467f50d85cea6b29fc98ebd7a8",
            "CombinedSignature": "0x046961ddc008eb0a8e91a9a73a82ab3197741be327b210b7a443b7c763af47970548c262ed1c6039163f075071a36dbf93a7c230d25bfc869d88ded50dd15d36",
            "WithdrawRoot": "0x7b76f0f62d3774ee059f48632072d284a0cd421abaf0415ad031efc5f3e22866",
            "Meta": {
                "SpokeID": 1,
                "TokenID": "0",
                "Amount": "50",
                "FeeReceiverStateID": 0
            }
        }
    ]
}
```

Example result (`DEPOSIT`):

```json
{
    "ID": "4",
    "Hash": "0xe3df23dfda31fc9498bdf990f05f65ca47a18a3c1247709bf9ae634a27970677",
    "Type": "DEPOSIT",
    "TransactionHash": "0xcd03405862fdc876bd6f1ccd2b6a843835324021e09ad8f86ebee2de67144cc9",
    "MinedBlock": 44027,
    "MinedTime": 1642589195,
    "Status": "MINED",
    "FinalisationBlock": 84347,
    "AccountTreeRoot": null,
    "Commitments": [
        {
            "ID": {
                "BatchID": "4",
                "IndexInBatch": 0
            },
            "PostStateRoot": "0xec16268cef171bc3ca07c1144d1d8a83d8e09f3ba4f43d7bba83c87729ffc62f",
            "LeafHash": "0x920330642c6cc8464fa50ef9ce46e0988b3f2013bc0d8a33fcfc00efec96f04d",
            "SubtreeID": "1",
            "SubtreeRoot": "0x819d6b845c2454d916cc180de8c9e29da8bbb89d17e6374d48d8e9f2de300f19"
        }
    ]
}
```

Example result (`GENESIS`):

```json
{
    "ID": "0",
    "Hash": "0x97222c9001ff4b2e8c328d15c02fc508259f2f9bec5af8b286a932813fe920d6",
    "Type": "GENESIS",
    "TransactionHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
    "MinedBlock": 44009,
    "MinedTime": null,
    "Status": "FINALISED",
    "FinalisationBlock": 44009,
    "AccountTreeRoot": null,
    "Commitments": null
}
```

### `hubble_getBatchByID(id)`

Same as hubble_getBatchByHash(hash), however this endpoint can also return commitments with `SUBMITTED` status.

Example result (`TRANSFER`) with `SUBMITTED` status:

```json
{
    "ID": "5",
    "Hash": null,
    "Type": "TRANSFER",
    "TransactionHash": "0xa428c43175e6156c933c436455a675c1e521e3841d9d44c6482ee26f861e49c1",
    "MinedBlock": null,
    "MinedTime": null,
    "Status": "SUBMITTED",
    "FinalisationBlock": null,
    "AccountTreeRoot": null,
    "Commitments": [
        {
            "ID": {
                "BatchID": "5",
                "IndexInBatch": 0
            },
            "PostStateRoot": "0x165d9644ba7aced9ed8cd69f2e06b6f5d619bb7e4ec9a0cffd2ff232f308cf4e",
            "LeafHash": null,
            "TokenID": "0",
            "FeeReceiverStateID": 0,
            "CombinedSignature": "0x2a2bd763f5d61e5ab2c8a96cb9961b1d18c4856546c61a28d4d4a2368e31a4d503e33bc6c29cdbb025c33091f41d016e377e27fee634f7b799da476889166a95"
        }
    ]
}
```

### `hubble_getCommitment(commitmentID: {BatchID: string, IndexInBatch: uint8})`

Returns commitment object including list of included transactions and commitment status. Commitment statuses:

- `SUBMITTED`
- `MINED`
- `FINALISED`

Example result (`TRANSFER`):

```json
{
    "ID": {
        "BatchID": "1",
        "IndexInBatch": 0
    },
    "Type": "TRANSFER",
    "PostStateRoot": "0x1be2feecb81d0338bf049a6e3324a254cd5b5a93f481d941a7f9e31dabc7cc9c",
    "LeafHash": "0xcb4df4741da627b1a21f85c793cde72c86d06b1402d968cab6ba6c2635348b01",
    "TokenID": "0",
    "FeeReceiverStateID": 0,
    "CombinedSignature": "0x046961ddc008eb0a8e91a9a73a82ab3197741be327b210b7a443b7c763af47970548c262ed1c6039163f075071a36dbf93a7c230d25bfc869d88ded50dd15d36",
    "Status": "MINED",
    "MinedTime": 1642591575,
    "Transactions": [
        {
            "Hash": "0x9b442316136f46247a399169aff5b9931060331f4b66971766a81b77765cfb36",
            "FromStateID": 1,
            "Amount": "50",
            "Fee": "1",
            "Nonce": "0",
            "Signature": "0x046961ddc008eb0a8e91a9a73a82ab3197741be327b210b7a443b7c763af47970548c262ed1c6039163f075071a36dbf93a7c230d25bfc869d88ded50dd15d36",
            // timestamp at which the tx was received by the coordinator for inclusion in batch. Can be null, when the tx was synced from blockchain.
            "ReceiveTime": 1642591573,
            "ToStateID": 2
        }
    ]
}
```

Example result (`CREATE2TRANSFER`):

```json
{
    "ID": {
        "BatchID": "2",
        "IndexInBatch": 0
    },
    "Type": "CREATE2TRANSFER",
    "PostStateRoot": "0x795a9e9e4bd52ac321c67eca39f5bb120f3c22181b0c825339df9be67a988ae6",
    "LeafHash": "0x03cfbdcb60a9372138d84a33738d3992b9d7c0e969f0ebcdad6f82259dfb9a3a",
    "TokenID": "0",
    "FeeReceiverStateID": 0,
    "CombinedSignature": "0x046961ddc008eb0a8e91a9a73a82ab3197741be327b210b7a443b7c763af47970548c262ed1c6039163f075071a36dbf93a7c230d25bfc869d88ded50dd15d36",
    "Status": "MINED",
    "MinedTime": 1642591655,
    "Transactions": [
        {
            "Hash": "0x3b98cb3f5c10647a80c753108ef14f44b4edebf096fe21fa49c7b578f4069ae0",
            "FromStateID": 2,
            "Amount": "50",
            "Fee": "1",
            "Nonce": "0",
            "Signature": "0x046961ddc008eb0a8e91a9a73a82ab3197741be327b210b7a443b7c763af47970548c262ed1c6039163f075071a36dbf93a7c230d25bfc869d88ded50dd15d36",
            // timestamp at which the tx was received by the coordinator for inclusion in batch. Can be null, when the tx was synced from blockchain.
            "ReceiveTime": 1642591653,
            "ToStateID": 6,
            "ToPublicKey": "0x0097f465fe827ce4dad751988f6ce5ec747458075992180ca11b0776b9ea3a910c3ee4dca4a03d06c3863778affe91ce38d502138356a35ae12695c565b24ea6151b83eabd41a6090b8ac3bb25e173c84c3b080a5545260b1327495920c342c02d51cac4418228db1a3d98aa12e6fd7b3267c703475f5999b2ec7a197ad7d8bc"
        }
    ]
}
```

Example result (`MASS_MIGRATION`):

```json
{
    "ID": {
        "BatchID": "3",
        "IndexInBatch": 0
    },
    "Type": "MASS_MIGRATION",
    "PostStateRoot": "0xe2afe7bb8d8f0310cf68636777794b97327f3ce3633127b58c4f5f566cbf8ac7",
    "LeafHash": "0x44f03af79586fd8ab74b442be1bb920e40eefd467f50d85cea6b29fc98ebd7a8",
    "CombinedSignature": "0x046961ddc008eb0a8e91a9a73a82ab3197741be327b210b7a443b7c763af47970548c262ed1c6039163f075071a36dbf93a7c230d25bfc869d88ded50dd15d36",
    "Status": "MINED",
    "MinedTime": 1642591641,
    "WithdrawRoot": "0x7b76f0f62d3774ee059f48632072d284a0cd421abaf0415ad031efc5f3e22866",
    "Meta": {
        "SpokeID": 1,
        "TokenID": "0",
        "Amount": "50",
        "FeeReceiverStateID": 0
    },
    "Transactions": [
        {
            "Hash": "0x716e1a2788c96e446e0248c3011b24c792afe7860fc674c9ac40b584a0b4953c",
            "FromStateID": 1,
            "Amount": "50",
            "Fee": "1",
            "Nonce": "1",
            "Signature": "0x028aafebc011dd95f2138740f66a05f7e8986111416227164a0e88a9e745ac06112d75e267d4be5aaa9e5509949f298529a0da1a6e86de0c9d5fc98325a8af4d",
            // timestamp at which the tx was received by the coordinator for inclusion in batch. Can be null, when the tx was synced from blockchain.
            "ReceiveTime": 1642591627,
            "SpokeID": 1
        }
    ]
}
```

Example result (`DEPOSIT`):

```json
{
    "ID": {
        "BatchID": "4",
        "IndexInBatch": 0
    },
    "Type": "DEPOSIT",
    "PostStateRoot": "0xec16268cef171bc3ca07c1144d1d8a83d8e09f3ba4f43d7bba83c87729ffc62f",
    "LeafHash": "0x920330642c6cc8464fa50ef9ce46e0988b3f2013bc0d8a33fcfc00efec96f04d",
    "Status": "MINED",
    "MinedTime": 1642589195,
    "SubtreeID": "1",
    "SubtreeRoot": "0x819d6b845c2454d916cc180de8c9e29da8bbb89d17e6374d48d8e9f2de300f19",
    "Deposits": [
        {
            "ID": {
                "SubtreeID": "1",
                "DepositIndex": "0"
            },
            "ToPubKeyID": 1,
            "TokenID": "0",
            "L2Amount": "10000000000"
        },
        {
            "ID": {
                "SubtreeID": "1",
                "DepositIndex": "1"
            },
            "ToPubKeyID": 2,
            "TokenID": "0",
            "L2Amount": "10000000000"
        },
        {
            "ID": {
                "SubtreeID": "1",
                "DepositIndex": "2"
            },
            "ToPubKeyID": 3,
            "TokenID": "0",
            "L2Amount": "10000000000"
        },
        {
            "ID": {
                "SubtreeID": "1",
                "DepositIndex": "3"
            },
            "ToPubKeyID": 4,
            "TokenID": "0",
            "L2Amount": "10000000000"
        }
    ]
}
```

Example result (`TRANSFER`) with `SUBMITTED` status:

```json
{
    "ID": {
        "BatchID": "5",
        "IndexInBatch": 0
    },
    "Type": "TRANSFER",
    "PostStateRoot": "0x165d9644ba7aced9ed8cd69f2e06b6f5d619bb7e4ec9a0cffd2ff232f308cf4e",
    "LeafHash": null,
    "TokenID": "0",
    "FeeReceiverStateID": 0,
    "CombinedSignature": "0x2a2bd763f5d61e5ab2c8a96cb9961b1d18c4856546c61a28d4d4a2368e31a4d503e33bc6c29cdbb025c33091f41d016e377e27fee634f7b799da476889166a95",
    "Status": "SUBMITTED",
    "MinedTime": null,
    "Transactions": [
        {
            "Hash": "0x9b442316136f46247a399169aff5b9931060331f4b66971766a81b77765cfb36",
            "FromStateID": 1,
            "Amount": "50",
            "Fee": "1",
            "Nonce": "0",
            "Signature": "0x2a2bd763f5d61e5ab2c8a96cb9961b1d18c4856546c61a28d4d4a2368e31a4d503e33bc6c29cdbb025c33091f41d016e377e27fee634f7b799da476889166a95",
            // timestamp at which the tx was received by the coordinator for inclusion in batch. Can be null, when the tx was synced from blockchain.
            "ReceiveTime": 1643379987,
            "ToStateID": 2
        }
    ]
}
```

### `hubble_getUserStateProof(stateID)`

Returns the merkle proof of the state tree and associated user state for the requested ID, see below.

```json
{
    "UserState": {
        "PubKeyID": 0,
        "TokenID": "0",
        "Balance": "1000000000000000000",
        "Nonce": "0"
    },
    "Witness": [
        "0x93081a8c3a12cc2c99211b299b84340f62bfe1c8c49678ed0873a2c69f233161",
        "0x772f16b497b46e7495658a7be1ab7b6502a13f041f9c8b97f67884719a23161e",
        "0x372034de54d25cd1250fd0fe608dc74903bd747572dc799dde9a8be89ba02fcf",
        "0x3b8ec09e026fdc305365dfc94e189a81b38c7597b3d941c279f042e8206e0bd8",
        "0xecd50eee38e386bd62be9bedb990706951b65fe053bd9d8a521af753d139e2da",
        "0xdefff6d330bb5403f63b14f33b578274160de3a50df4efecf0e0db73bcdd3da5",
        "0x617bdd11f7c0a11f49db22f629387a12da7596f9d1704d7465177c63d88ec7d7",
        "0x292c23a9aa1d8bea7e2435e555a4a60e379a5a35f3f452bae60121073fb6eead",
        "0xe1cea92ed99acdcb045a6726b2f87107e8a61620a232cf4d7d5b5766b3952e10",
        "0x7ad66c0a68c72cb89e4fb4303841966e4062a76ab97451e3b9fb526a5ceb7f82",
        "0xe026cc5a4aed3c22a58cbd3d2ac754c9352c5436f638042dca99034e83636516",
        "0x3d04cffd8b46a874edf5cfae63077de85f849a660426697b06a829c70dd1409c",
        "0xad676aa337a485e4728a0b240d92b3ef7b3c372d06d189322bfd5f61f1e7203e",
        "0xa2fca4a49658f9fab7aa63289c91b7c7b6c832a6d0e69334ff5b0a3483d09dab",
        "0x4ebfd9cd7bca2505f7bef59cc1c12ecc708fff26ae4af19abe852afe9e20c862",
        "0x2def10d13dd169f550f578bda343d9717a138562e0093b380a1120789d53cf10",
        "0x776a31db34a1a0a7caaf862cffdfff1789297ffadc380bd3d39281d340abd3ad",
        "0xe2e7610b87a5fdf3a72ebe271287d923ab990eefac64b6e59d79f8b7e08c46e3",
        "0x504364a5c6858bf98fff714ab5be9de19ed31a976860efbd0e772a2efe23e2e0",
        "0x4f05f4acb83f5b65168d9fef89d56d4d77b8944015e6b1eed81b0238e2d0dba3",
        "0x44a6d974c75b07423e1d6d33f481916fdd45830aea11b6347e700cd8b9f0767c",
        "0xedf260291f734ddac396a956127dde4c34c0cfb8d8052f88ac139658ccf2d507",
        "0x6075c657a105351e7f0fce53bc320113324a522e8fd52dc878c762551e01a46e",
        "0x6ca6a3f763a9395f7da16014725ca7ee17e4815c0ff8119bf33f273dee11833b",
        "0x1c25ef10ffeb3c7d08aa707d17286e0b0d3cbcb50f1bd3b6523b63ba3b52dd0f",
        "0xfffc43bd08273ccf135fd3cacbeef055418e09eb728d727c4d5d5c556cdea7e3",
        "0xc5ab8111456b1f28f3c7a0a604b4553ce905cb019c463ee159137af83c350b22",
        "0x0ff273fcbf4ae0f2bd88d6cf319ff4004f8d7dca70d4ced4e74d2c74139739e6",
        "0x7fa06ba11241ddd5efdc65d4e39c9f6991b74fd4b81b62230808216c876f827c",
        "0x7e275adf313a996c7e2950cac67caba02a5ff925ebf9906b58949f3e77aec5b9",
        "0x8f6162fa308d2b3a15dc33cffac85f13ab349173121645aedf00f471663108be",
        "0x78ccaaab73373552f207a63599de54d7d8d0c1805f86ce7da15818d09f4cff62"
    ]
}
```

### `hubble_getPublicKeyProofByPubKeyID(pubKeyID)`

Returns the merkle path and associated public key for the requested public key ID, see below.

```json
{
    "PublicKey": "0x0df68cb87856229b0bc3f158fff8b82b04deb1a4c23dadbf3ed2da4ec6f6efcb1c165c6b47d8c89ab2ddb0831c182237b27a4b3d9701775ad6c180303f87ef260566cb2f0bcc7b89c2260de2fee8ec29d7b5e575a1e36eb4bcead52a74a511b7188d7df7c9d08f94b9daa9d89105fbdf22bf14e30b84f8adefb3695ebff00e88",
    "Witness": [
        "0x0219c7e21708fe950f8c3f2107150aafb47cb9c181a932994905a5bd64b04170",
        "0x388aaaa9155b428f04a5fa0075f7243100156545205902faf78d0f0afd888469",
        "0x172658ed181de68736fb9dfdda8a4313a9396cdbbdbdaf7573036e7174607e88",
        "0x8eb2750868f95e70d25aba7070e4126bf756cf7bd0d172903c9808940aeaa129",
        "0xf3868669437495e7e5f8726bb9ac7e0b40fe3a7b15dca40caf3654ac4ebf64be",
        "0x153a5b39fad8ab5bf7c3789f4fc47a36497f1f6cd3cbbeb34b29c4a805c37708",
        "0x617bdd11f7c0a11f49db22f629387a12da7596f9d1704d7465177c63d88ec7d7",
        "0x292c23a9aa1d8bea7e2435e555a4a60e379a5a35f3f452bae60121073fb6eead",
        "0xe1cea92ed99acdcb045a6726b2f87107e8a61620a232cf4d7d5b5766b3952e10",
        "0x7ad66c0a68c72cb89e4fb4303841966e4062a76ab97451e3b9fb526a5ceb7f82",
        "0xe026cc5a4aed3c22a58cbd3d2ac754c9352c5436f638042dca99034e83636516",
        "0x3d04cffd8b46a874edf5cfae63077de85f849a660426697b06a829c70dd1409c",
        "0xad676aa337a485e4728a0b240d92b3ef7b3c372d06d189322bfd5f61f1e7203e",
        "0xa2fca4a49658f9fab7aa63289c91b7c7b6c832a6d0e69334ff5b0a3483d09dab",
        "0x4ebfd9cd7bca2505f7bef59cc1c12ecc708fff26ae4af19abe852afe9e20c862",
        "0x2def10d13dd169f550f578bda343d9717a138562e0093b380a1120789d53cf10",
        "0x776a31db34a1a0a7caaf862cffdfff1789297ffadc380bd3d39281d340abd3ad",
        "0xe2e7610b87a5fdf3a72ebe271287d923ab990eefac64b6e59d79f8b7e08c46e3",
        "0x504364a5c6858bf98fff714ab5be9de19ed31a976860efbd0e772a2efe23e2e0",
        "0x4f05f4acb83f5b65168d9fef89d56d4d77b8944015e6b1eed81b0238e2d0dba3",
        "0x44a6d974c75b07423e1d6d33f481916fdd45830aea11b6347e700cd8b9f0767c",
        "0xedf260291f734ddac396a956127dde4c34c0cfb8d8052f88ac139658ccf2d507",
        "0x6075c657a105351e7f0fce53bc320113324a522e8fd52dc878c762551e01a46e",
        "0x6ca6a3f763a9395f7da16014725ca7ee17e4815c0ff8119bf33f273dee11833b",
        "0x1c25ef10ffeb3c7d08aa707d17286e0b0d3cbcb50f1bd3b6523b63ba3b52dd0f",
        "0xfffc43bd08273ccf135fd3cacbeef055418e09eb728d727c4d5d5c556cdea7e3",
        "0xc5ab8111456b1f28f3c7a0a604b4553ce905cb019c463ee159137af83c350b22",
        "0x0ff273fcbf4ae0f2bd88d6cf319ff4004f8d7dca70d4ced4e74d2c74139739e6",
        "0x7fa06ba11241ddd5efdc65d4e39c9f6991b74fd4b81b62230808216c876f827c",
        "0x7e275adf313a996c7e2950cac67caba02a5ff925ebf9906b58949f3e77aec5b9",
        "0x8f6162fa308d2b3a15dc33cffac85f13ab349173121645aedf00f471663108be",
        "0x78ccaaab73373552f207a63599de54d7d8d0c1805f86ce7da15818d09f4cff62"
    ]
}
```

### `hubble_getCommitmentProof(commitmentID)`

Returns the transfer/create2transfer commitment inclusion proof for the given commitment ID, see below.

```json
{
    "StateRoot": "0x81cf78ec55d3393ff2e9c0e081dc6ced3cd4a7e9e42f3c6e441b035035a6839a",
    "Body": {
        "AccountRoot": "0xb261c40259ad5dbaf32efb2256225bbf03dcda8e84cffdfe67e68b958e3c7a95",
        "Signature": "0x2e03d13aec0f8bad52b3045d245c1318ecbeacbc2549642bcc11a9aef1fabdc51fc104a88321c5093897a679e244d32e306c69187fbbc5f81bbd50cd71f2bbc9",
        "FeeReceiver": 0,
        "Transactions": [
            {
                "Hash": "0x9b442316136f46247a399169aff5b9931060331f4b66971766a81b77765cfb36",
                "FromStateID": 1,
                "Amount": "50",
                "Fee": "1",
                "Nonce": "0",
                "Signature": "0x2e03d13aec0f8bad52b3045d245c1318ecbeacbc2549642bcc11a9aef1fabdc51fc104a88321c5093897a679e244d32e306c69187fbbc5f81bbd50cd71f2bbc9",
                "ReceiveTime": 1634691270,
                "ToStateID": 2
            }
        ]
    },
    "Path": {
        "Path": 0,
        "Depth": 2
    },
    "Witness": [
        "0x290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563"
    ]
}
```

### `hubble_getMassMigrationCommitmentProof(commitmentID)`

Returns the mass migration commitment inclusion proof for the given commitment ID, see below.

```json
{
    "StateRoot": "0xe2afe7bb8d8f0310cf68636777794b97327f3ce3633127b58c4f5f566cbf8ac7",
    "Path": {
        "Path": 0,
        "Depth": 2
    },
    "Witness": [
        "0x290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563"
    ],
    "Body": {
        "AccountRoot": "0xb261c40259ad5dbaf32efb2256225bbf03dcda8e84cffdfe67e68b958e3c7a95",
        "Signature": "0x046961ddc008eb0a8e91a9a73a82ab3197741be327b210b7a443b7c763af47970548c262ed1c6039163f075071a36dbf93a7c230d25bfc869d88ded50dd15d36",
        "Meta": {
            "SpokeID": 1,
            "TokenID": "0",
            "Amount": "50",
            "FeeReceiverStateID": 0
        },
        "WithdrawRoot": "0x7b76f0f62d3774ee059f48632072d284a0cd421abaf0415ad031efc5f3e22866",
        // transactions are serialized and encoded using Base64 encoding format
        "Transactions": "AAAAARAFAAE="
    }
}
```

### `hubble_getWithdrawProof(commitmentID, transactionHash)`

Returns the withdrawal proof for the given mass migration transaction, see below.

```json
{
    "UserState": {
        "PubKeyID": 1,
        "TokenID": "0",
        "Balance": "50",
        "Nonce": "0"
    },
    "Path": {
        "Path": 0,
        "Depth": 2
    },
    "Witness": [
        "0x290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563"
    ],
    "Root": "0x7b76f0f62d3774ee059f48632072d284a0cd421abaf0415ad031efc5f3e22866"
}
```

## Admin API

Admin API endpoints requires authentication via authentication key specified in config.

### `admin_getPendingBatches()`

Returns pending batches that haven't been mined yet, see below.

```json
[
    {
        "ID": "1",
        "Type": "TRANSFER",
        "TransactionHash": "0x7c24f82171937322d49b382ec79d6fb40bc10086694d805a399b92c061b17e0d",
        "Commitments": [
            {
                "Commitment": {
                    "ID": {
                        "BatchID": "1",
                        "IndexInBatch": 0
                    },
                    "Type": "TRANSFER",
                    "PostStateRoot": "0xae6b8d73ccfca01f6d893c59b9353b9a81bd7d0b4f6d602fc35da9fd0161ee64",
                    "FeeReceiver": 0,
                    "CombinedSignature": "0x2041452aa916e776ded08c4bd9cba2614e4632e6a2c4b18116ae7ce69ac0d71903a31b7fc04c2dce4601ec0a9d301e8f0de3a7cb0b7ed5683219e6403c58eaf9",
                    "BodyHash": null
                },
                "Transactions": [
                    {
                        "Hash": "0x9b442316136f46247a399169aff5b9931060331f4b66971766a81b77765cfb36",
                        "TxType": "TRANSFER",
                        "FromStateID": 1,
                        "Amount": "50",
                        "Fee": "1",
                        "Nonce": "0",
                        "Signature": "0x2041452aa916e776ded08c4bd9cba2614e4632e6a2c4b18116ae7ce69ac0d71903a31b7fc04c2dce4601ec0a9d301e8f0de3a7cb0b7ed5683219e6403c58eaf9",
                        "ReceiveTime": 1642845254,
                        "CommitmentID": {
                            "BatchID": "1",
                            "IndexInBatch": 0
                        },
                        "ErrorMessage": null,
                        "ToStateID": 2
                    }
                ]
            }
        ]
    },
    {
        "ID": "2",
        "Type": "CREATE2TRANSFER",
        "TransactionHash": "0x5326e17a37dcf9358c4a969cfda8fc34c598031c8dbe6effa56e9e15caff32d7",
        "Commitments": [
            {
                "Commitment": {
                    "ID": {
                        "BatchID": "2",
                        "IndexInBatch": 0
                    },
                    "Type": "CREATE2TRANSFER",
                    "PostStateRoot": "0x5d9cdfe6e2516b669d841b0006dc664580abd436b0fa216692bd45c05d24a80e",
                    "FeeReceiver": 0,
                    "CombinedSignature": "0x2041452aa916e776ded08c4bd9cba2614e4632e6a2c4b18116ae7ce69ac0d71903a31b7fc04c2dce4601ec0a9d301e8f0de3a7cb0b7ed5683219e6403c58eaf9",
                    "BodyHash": null
                },
                "Transactions": [
                    {
                        "Hash": "0x3b98cb3f5c10647a80c753108ef14f44b4edebf096fe21fa49c7b578f4069ae0",
                        "TxType": "CREATE2TRANSFER",
                        "FromStateID": 2,
                        "Amount": "50",
                        "Fee": "1",
                        "Nonce": "0",
                        "Signature": "0x2041452aa916e776ded08c4bd9cba2614e4632e6a2c4b18116ae7ce69ac0d71903a31b7fc04c2dce4601ec0a9d301e8f0de3a7cb0b7ed5683219e6403c58eaf9",
                        "ReceiveTime": 1642845266,
                        "CommitmentID": {
                            "BatchID": "2",
                            "IndexInBatch": 0
                        },
                        "ErrorMessage": null,
                        "ToStateID": 6,
                        "ToPublicKey": "0x0097f465fe827ce4dad751988f6ce5ec747458075992180ca11b0776b9ea3a910c3ee4dca4a03d06c3863778affe91ce38d502138356a35ae12695c565b24ea6151b83eabd41a6090b8ac3bb25e173c84c3b080a5545260b1327495920c342c02d51cac4418228db1a3d98aa12e6fd7b3267c703475f5999b2ec7a197ad7d8bc"
                    }
                ]
            }
        ]
    },
    {
        "ID": "3",
        "Type": "MASS_MIGRATION",
        "TransactionHash": "0x26f1f6307fa8be04d5f1b56292af3fc0cda1cd4cd6d3cbed71ec9d52e2df24db",
        "Commitments": [
            {
                "Commitment": {
                    "ID": {
                        "BatchID": "3",
                        "IndexInBatch": 0
                    },
                    "Type": "MASS_MIGRATION",
                    "PostStateRoot": "0x2212a3961a54c730d796a674ae10155a9d573de2418a0ccdd228ca92650178aa",
                    "FeeReceiver": 0,
                    "CombinedSignature": "0x2041452aa916e776ded08c4bd9cba2614e4632e6a2c4b18116ae7ce69ac0d71903a31b7fc04c2dce4601ec0a9d301e8f0de3a7cb0b7ed5683219e6403c58eaf9",
                    "BodyHash": null,
                    "Meta": {
                        "SpokeID": 1,
                        "TokenID": "0",
                        "Amount": "50",
                        "FeeReceiver": 0
                    },
                    "WithdrawRoot": "0x8dbe126663cf1d1c53e6c6b52d0fbc180c928af61a5cb2b2744d5708234d51cb"
                },
                "Transactions": [
                    {
                        "Hash": "0x7a43485af7a557de8ea4e00b9e11c77d24cab35a8f7a3d8b87f3815e1b4ece31",
                        "TxType": "MASS_MIGRATION",
                        "FromStateID": 0,
                        "Amount": "50",
                        "Fee": "1",
                        "Nonce": "0",
                        "Signature": "0x2041452aa916e776ded08c4bd9cba2614e4632e6a2c4b18116ae7ce69ac0d71903a31b7fc04c2dce4601ec0a9d301e8f0de3a7cb0b7ed5683219e6403c58eaf9",
                        "ReceiveTime": 1642845275,
                        "CommitmentID": {
                            "BatchID": "3",
                            "IndexInBatch": 0
                        },
                        "ErrorMessage": null,
                        "SpokeID": 1
                    }
                ]
            }
        ]
    },
    {
        "ID": "4",
        "Type": "DEPOSIT",
        "TransactionHash": "0x6b862c990d0c3b7d1f8d87b303fc43a0ce1f926f7c04f75498f661696d5ecf98",
        "Commitments": [
            {
                "Commitment": {
                    "ID": {
                        "BatchID": "4",
                        "IndexInBatch": 0
                    },
                    "Type": "DEPOSIT",
                    "PostStateRoot": "0xec16268cef171bc3ca07c1144d1d8a83d8e09f3ba4f43d7bba83c87729ffc62f",
                    "SubtreeID": "1",
                    "SubtreeRoot": "0x819d6b845c2454d916cc180de8c9e29da8bbb89d17e6374d48d8e9f2de300f19",
                    "Deposits": [
                        {
                            "ID": {
                                "SubtreeID": "1",
                                "DepositIndex": "0"
                            },
                            "ToPubKeyID": 1,
                            "TokenID": "0",
                            "L2Amount": "10000000000"
                        },
                        {
                            "ID": {
                                "SubtreeID": "1",
                                "DepositIndex": "1"
                            },
                            "ToPubKeyID": 1,
                            "TokenID": "0",
                            "L2Amount": "10000000000"
                        },
                        {
                            "ID": {
                                "SubtreeID": "1",
                                "DepositIndex": "2"
                            },
                            "ToPubKeyID": 1,
                            "TokenID": "0",
                            "L2Amount": "10000000000"
                        },
                        {
                            "ID": {
                                "SubtreeID": "1",
                                "DepositIndex": "3"
                            },
                            "ToPubKeyID": 1,
                            "TokenID": "0",
                            "L2Amount": "10000000000"
                        }
                    ]
                },
                "Transactions": null
            }
        ]
    }
]
```

### `admin_getPendingTransactions()`

Returns pending transactions that haven't been included in batch yet, see below.

```json
[
    {
        "Hash": "0x3b98cb3f5c10647a80c753108ef14f44b4edebf096fe21fa49c7b578f4069ae0",
        "TxType": "CREATE2TRANSFER",
        "FromStateID": 2,
        "Amount": "50",
        "Fee": "1",
        "Nonce": "0",
        "Signature": "0x0f94ee2c5207b2f0ce4627a7083c1ea9619faeb83d4ae458677b479520cdc7b70032531aca0e42cf1ddc5b2f2e40c15bd78d7fbbb5b385d59b5c9855bb83e1cc",
        "ReceiveTime": 1643970131,
        "CommitmentID": null,
        "ErrorMessage": null,
        "ToStateID": null,
        "ToPublicKey": "0x0097f465fe827ce4dad751988f6ce5ec747458075992180ca11b0776b9ea3a910c3ee4dca4a03d06c3863778affe91ce38d502138356a35ae12695c565b24ea6151b83eabd41a6090b8ac3bb25e173c84c3b080a5545260b1327495920c342c02d51cac4418228db1a3d98aa12e6fd7b3267c703475f5999b2ec7a197ad7d8bc"
    }
]
```

### `admin_getFailedTransactions()`

Returns all failed transactions.

```json
[
    {
        "Hash": "0x25c3b159f42a82dc8277a662b90a079a61c569bff8784e8ebd51bfd610792901",
        "TxType": "TRANSFER",
        "FromStateID": 1,
        "Amount": "50",
        "Fee": "1",
        "Nonce": "0",
        "Signature": "0x191e0510f1657c01c22b65c3548806ba2bd19faa721c61516bf49cbb787e120a2b416731ac80dd38a561bba25aac89c6474a37bce42f3b213d135c920221d810",
        "ReceiveTime": 1643977662,
        "CommitmentID": null,
        "ErrorMessage": "receiver state ID does not exist",
        "ToStateID": 100
    }
]
```

### `admin_configure(configureParams)`

Can be used to enable/disable:

- Accepting new transactions by `hubble_sendTransaction`.
- Creating and submitting new transaction batches.

# API usage

## Sending a transaction

1. Call `hubble_getUserStates(senderPubKey)` to list sender's accounts. Pick one with an appropriate token index and balance.
2. Call `hubble_getUserStates(recipientPubKey)` to list recipients accounts pick one with an appropriate token index.
3. Call `hubble_sendTransaction` with `txType=TRANSFER` using the state indexes from steps 1 & 2 and nonce from step 1.
4. Call `hubble_getTransfer(Hash)` with the hash from step 3 to monitor transaction progress.

## Alternative: Recipient doesn't have a state leaf for a given token

1. Call `hubble_getUserStates(senderPubKey)` to list sender's accounts. Pick one with an appropriate token index and balance.
2. Call `hubble_sendTransaction` with `txType=CREATE2TRANSFER` with the recipient's pubKey. The commander will decide the appropriate
   account index and state index for the recipient.
3. Call `hubble_getTransaction(Hash)` with the hash from step 2 to monitor transaction progress.
