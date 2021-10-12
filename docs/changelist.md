# Changelist 0.4.0 -> master (b43d41eca3598515a687253beb08fd2cc5e8fca5)

## Changes

- New standardized API errors (refer to the table at the bottom of the changelist)
- Data models changes:
    - `Commitment`
        - `ID` type changed from `uint32` to `{BatchID: string, IndexInBatch: uint8}` object
        - `IncludedInBatch` field was removed
    - `NetworkInfo`
        - Renamed fields:
            - `chainID` -> `ChainID`
            - `accountRegistry` -> `AccountRegistry`
            - `deploymentBlock` -> `AccountRegistryDeploymentBlock`
            - `rollup` -> `Rollup`
            - `blockNumber` -> `BlockNumber`
            - `transactionCount` -> `TransactionCount`
            - `accountCount` -> `AccountCount`
            - `latestBatch` -> `LatestBatch`
            - `latestFinalisedBatch` -> `LatestFinalisedBatch`
            - `signatureDomain` -> `SignatureDomain`
        - New fields:
            - `TokenRegistry` string (address)
            - `DepositManager` string (address)
    - `Transfer`/`Create2Transfer`:
        - `IncludedInCommitment [int32]` was replaced with `CommitmentID [object]`
            - `CommitmentID` type is the `{BatchID: string, IndexInBatch: uint8}` object
            - refer to the API change made to `hubble_getTransactions` endpoint for an example

- API changes:
    - `hubble_getPublicKeyByID` endpoint renamed to `hubble_getPublicKeyByPubKeyID`
    - `hubble_getNetworkInfo`
        - New output - refer to the changes made to the `NetworkInfo` data model
            - Before:
              ```json
              {
                "chainId": "1337",
                "accountRegistry": "0x10bd6732fe3908b8a816f6a1b271e0864de78ca1",
                "deploymentBlock": 74,
                "rollup": "0xf2a409ccf78e6e32e02d5e3a3ac274ca6880d9ac",
                "blockNumber": 2146,
                "transactionCount": 2,
                "accountCount": 6,
                "latestBatch": "2",
                "latestFinalisedBatch": "0",
                "signatureDomain": "0x47b39cc40c04341a600ee0941a8231bf3a04725da5c65ac93286ef9147d23bbc"
              }
              ```
            - After:
              ```json
              {
                "ChainID": "1337",
                "AccountRegistry": "0x10bd6732fe3908b8a816f6a1b271e0864de78ca1",
                "AccountRegistryDeploymentBlock": 74,
                "TokenRegistry": "0x07389715ae1f0a891fba82e65099f6a3fa7da593",
                "DepositManager": "0xa3accd1cfabc8b09aea4d0e25f21f25c526c9be8",
                "Rollup": "0xf2a409ccf78e6e32e02d5e3a3ac274ca6880d9ac",
                "BlockNumber": 2146,
                "TransactionCount": 2,
                "AccountCount": 6,
                "LatestBatch": "2",
                "LatestFinalisedBatch": "0",
                "SignatureDomain": "0x47b39cc40c04341a600ee0941a8231bf3a04725da5c65ac93286ef9147d23bbc"
              }
              ```
    - `hubble_getTransactions`
        - New output - refer to the changes made to the `Transaction`/`Create2Transfer` data models
            - Before:
              ```json
              [
                {
                  "Hash": "0x03b15bc97adb5e86fffbffd8b049629b80eb696499fe3aa62813fcbed87f4023",
                  "TxType": "TRANSFER",
                  "FromStateID": 1,
                  "Amount": "50",
                  "Fee": "1",
                  "Nonce": "1",
                  "Signature": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
                  "ReceiveTime": null,
                  "IncludedInCommitment": 2, // removed
                  "ErrorMessage": null,
                  "ToStateID": 2,
                  "BatchHash": "0x87e7380690ae69c1d06796828f12113eadc436942a3dd3aa9182eb1c9b164f90",
                  "BatchTime": 1633693199,
                  "Status": "IN_BATCH"
                },
                {
                  "Hash": "0x9b442316136f46247a399169aff5b9931060331f4b66971766a81b77765cfb36",
                  "TxType": "TRANSFER",
                  "FromStateID": 1,
                  "Amount": "50",
                  "Fee": "1",
                  "Nonce": "0",
                  "Signature": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
                  "ReceiveTime": null,
                  "IncludedInCommitment": 1, // removed
                  "ErrorMessage": null,
                  "ToStateID": 2,
                  "BatchHash": "0xeb590ba0ce14d821caebc56514fe867521da78b46b7b78ce4810a353e619f315",
                  "BatchTime": 1633692591,
                  "Status": "IN_BATCH"
                },
              ]
              ```
            - After:
              ```json
              [
                {
                  "Hash": "0x03b15bc97adb5e86fffbffd8b049629b80eb696499fe3aa62813fcbed87f4023",
                  "TxType": "TRANSFER",
                  "FromStateID": 1,
                  "Amount": "50",
                  "Fee": "1",
                  "Nonce": "1",
                  "Signature": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
                  "ReceiveTime": null,
                  "CommitmentID": { // new
                      "BatchID": "2",
                      "IndexInBatch": 0
                  },
                  "ErrorMessage": null,
                  "ToStateID": 2,
                  "BatchHash": "0x87e7380690ae69c1d06796828f12113eadc436942a3dd3aa9182eb1c9b164f90",
                  "BatchTime": 1633693199,
                  "Status": "IN_BATCH"
                },
                {
                  "Hash": "0x9b442316136f46247a399169aff5b9931060331f4b66971766a81b77765cfb36",
                  "TxType": "TRANSFER",
                  "FromStateID": 1,
                  "Amount": "50",
                  "Fee": "1",
                  "Nonce": "0",
                  "Signature": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
                  "ReceiveTime": null,
                  "CommitmentID": { // new
                    "BatchID": "1",
                    "IndexInBatch": 0
                  },
                  "ErrorMessage": null,
                  "ToStateID": 2,
                  "BatchHash": "0xeb590ba0ce14d821caebc56514fe867521da78b46b7b78ce4810a353e619f315",
                  "BatchTime": 1633692591,
                  "Status": "IN_BATCH"
                }
              ]
              ```
    - `hubble_getTransaction`
        - New output - refer to the changes made to the `Transaction`/`Create2Transfer` data models
            - Before:
              ```json
              {
                "Hash": "0x9b442316136f46247a399169aff5b9931060331f4b66971766a81b77765cfb36",
                "TxType": "TRANSFER",
                "FromStateID": 1,
                "Amount": "50",
                "Fee": "1",
                "Nonce": "0",
                "Signature": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
                "ReceiveTime": null,
                "IncludedInCommitment": 1, // removed
                "ErrorMessage": null,
                "ToStateID": 2,
                "BatchHash": "0xeb590ba0ce14d821caebc56514fe867521da78b46b7b78ce4810a353e619f315",
                "BatchTime": 1633692591,
                "Status": "IN_BATCH"
              }
              ```
            - After:
              ```json
              {
                "Hash": "0x9b442316136f46247a399169aff5b9931060331f4b66971766a81b77765cfb36",
                "TxType": "TRANSFER",
                "FromStateID": 1,
                "Amount": "50",
                "Fee": "1",
                "Nonce": "0",
                "Signature": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
                "ReceiveTime": null,
                "CommitmentID": { // new
                    "BatchID": "1",
                    "IndexInBatch": 0
                },
                "ErrorMessage": null,
                "ToStateID": 2,
                "BatchHash": "0xeb590ba0ce14d821caebc56514fe867521da78b46b7b78ce4810a353e619f315",
                "BatchTime": 1633692591,
                "Status": "IN_BATCH"
              }
              ```
    - `hubble_getCommitment`
        - The parameter for `hubble_getCommitment` endpoint was changed:
            - `[uint32]` -> `[{"BatchID": string, "IndexInBatch": uint8}]`
                - Example: `[1]` -> `[{"BatchID": "1", "IndexInBatch": 0}]`
        - New output - refer to the changes made to the `Commitment` data model
            - Before:
              ```json
              {
                "ID": 1, // changed
                "Type": "TRANSFER",
                "FeeReceiver": 0,
                "CombinedSignature": "0x1152450e7da64c68023921d3a37ea750df4158bb17203317bf7af9ac7d8c6a3216d982a417c204593c82dc1f64851cad49361a4a4175636e0c062497c7ef2f9c",
                "PostStateRoot": "0x81cf78ec55d3393ff2e9c0e081dc6ced3cd4a7e9e42f3c6e441b035035a6839a",
                "IncludedInBatch": "1", // removed
                "Status": "IN_BATCH",
                "BatchTime": 1633692591,
                "Transactions": [
                  {
                    "Hash": "0x9b442316136f46247a399169aff5b9931060331f4b66971766a81b77765cfb36",
                    "FromStateID": 1,
                    "Amount": "50",
                    "Fee": "1",
                    "Nonce": "0",
                    "Signature": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
                    "ReceiveTime": 1633692588,
                    "ToStateID": 2
                  }
                ]
              }
              ```
            - After:
              ```json
              {
                "ID": { // new
                  "BatchID": "1",
                  "IndexInBatch": 0
                },
                "Type": "TRANSFER",
                "FeeReceiver": 0,
                "CombinedSignature": "0x1152450e7da64c68023921d3a37ea750df4158bb17203317bf7af9ac7d8c6a3216d982a417c204593c82dc1f64851cad49361a4a4175636e0c062497c7ef2f9c",
                "PostStateRoot": "0x81cf78ec55d3393ff2e9c0e081dc6ced3cd4a7e9e42f3c6e441b035035a6839a",
                "Status": "IN_BATCH",
                "BatchTime": 1633692591,
                "Transactions": [
                  {
                    "Hash": "0x9b442316136f46247a399169aff5b9931060331f4b66971766a81b77765cfb36",
                    "FromStateID": 1,
                    "Amount": "50",
                    "Fee": "1",
                    "Nonce": "0",
                    "Signature": "0x1152450e7da64c68023921d3a37ea750df4158bb17203317bf7af9ac7d8c6a3216d982a417c204593c82dc1f64851cad49361a4a4175636e0c062497c7ef2f9c",
                    "ReceiveTime": 1633692588,
                    "ToStateID": 2
                  }
                ]
              }
              ```
    - `hubble_getBatchByID`/`hubble_getBatchByHash`
        - New output - refer to the changes made to the `Commitment` data model
            - Before:
              ```json
              {
                "ID": "2",
                "Hash": "0x87e7380690ae69c1d06796828f12113eadc436942a3dd3aa9182eb1c9b164f90",
                "Type": "TRANSFER",
                "TransactionHash": "0xd80f934de37a5400ea434c7349baa01923a5d04b2817adb5e66ee6b543d368f6",
                "SubmissionBlock": 851,
                "SubmissionTime": 1633693199,
                "FinalisationBlock": 41171,
                "AccountTreeRoot": "0xb261c40259ad5dbaf32efb2256225bbf03dcda8e84cffdfe67e68b958e3c7a95",
                "Commitments": [
                  {
                    "ID": 2, // changed
                    "LeafHash": "0xfa61444076068795954727098250f88da93d9e75fd76b476ad41da5bd4f6cba1",
                    "TokenID": "0",
                    "FeeReceiverStateID": 0,
                    "CombinedSignature": "0x233fb912ea37ec1e18e1ce5e24ebfe2d80c35c6d9431ba74017a884c635569442c708e2b1f5163d99c81c7fd2cd5406d1831f8fe2e2da8003a85b7441ef0403f",
                    "PostStateRoot": "0xd4f3de11d8b3035b163da99dffcac576975b8cac358d9471b8a1e9ef4d6fbb30"
                  }
                ]
              }
              ```
            - After:
              ```json
              {
                "ID": "2",
                "Hash": "0x87e7380690ae69c1d06796828f12113eadc436942a3dd3aa9182eb1c9b164f90",
                "Type": "TRANSFER",
                "TransactionHash": "0xd80f934de37a5400ea434c7349baa01923a5d04b2817adb5e66ee6b543d368f6",
                "SubmissionBlock": 851,
                "SubmissionTime": 1633693199,
                "FinalisationBlock": 41171,
                "AccountTreeRoot": "0xb261c40259ad5dbaf32efb2256225bbf03dcda8e84cffdfe67e68b958e3c7a95",
                "Commitments": [
                  {
                    "ID": { // changed
                      "BatchID": "2",
                      "IndexInBatch": 0
                    },
                    "LeafHash": "0xfa61444076068795954727098250f88da93d9e75fd76b476ad41da5bd4f6cba1",
                    "TokenID": "0",
                    "FeeReceiverStateID": 0,
                    "CombinedSignature": "0x233fb912ea37ec1e18e1ce5e24ebfe2d80c35c6d9431ba74017a884c635569442c708e2b1f5163d99c81c7fd2cd5406d1831f8fe2e2da8003a85b7441ef0403f",
                    "PostStateRoot": "0xd4f3de11d8b3035b163da99dffcac576975b8cac358d9471b8a1e9ef4d6fbb30"
                  }
                ]
              }
              ```

## API Error code ranges

- `999` - Unknown Errors
- `10XXX` - Transaction Errors
- `20XXX` - Commitment Errors
- `30XXX` - Batch Errors
- `40XXX` - Badger Errors
- `99XXX` - Uncategorized Errors like NetworkInfo, BLS, UserStates etc.

## API Errors

|  Error code  |                                Message                              |
| -------------| ------------------------------------------------------------------- |
| `999`        | `unknown error: COMMANDER_ERROR`                                    |
| `10000`      | `transaction not found`                                             |
| `10001`      | `transactions not found`                                            |
| `10002`      | `some field is missing, verify the transfer/create2transfer object` |
| `10003`      | `invalid recipient, cannot send funds to yourself`                  |
| `10004`      | `nonce too low`                                                     |
| `10005`      | `nonce too high`                                                    |
| `10006`      | `not enough balance`                                                |
| `10007`      | `amount must be greater than 0`                                     |
| `10008`      | `fee too low`                                                       |
| `10009`      | `invalid signature`                                                 |
| `10010`      | `amount is not encodable as multi-precission decimal`               |
| `10011`      | `fee is not encodable as multi-precission decimal`                  |
| `20000`      | `commitment not found`                                              |
| `30000`      | `batch not found`                                                   |
| `30001`      | `batches not found`                                                 |
| `40000`      | `an error occurred while saving data to the Badger database`        |
| `40001`      | `an error occurred while iterating over Badger database`            |
| `99000`      | `an error occurred while fetching the account count`                |
| `99001`      | `public key not found`                                              |
| `99002`      | `user state not found`                                              |
| `99003`      | `user states not found`                                             |
| `99004`      | `an error occurred while fetching the domain for signing`           |
