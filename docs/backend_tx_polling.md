# Backend Transaction Polling

**Main loop:**

1. Call `hubble_getNetworkInfo` every 500ms
2. For each batch between last call and `latestBatch` :
    1. Call `hubble_getBatchByID`
    2. For each commitment in Commitments
    3. Call `hubble_getCommitment` with `ID` as parameter
    4. If Type == "TRANSFER" (see below for example)
        - For each transaction in Transactions
            - Get sender (get user from database by `FromStateID` )
            - Get receiver (get user from database by `ToStateID` )
            - Store transaction with current timestamp (TODO: needs to be changed to real timestamp from crypto side) and sender and receiver user uuid
            - Execute push notification to receiver uuid
    5. If Type == "CREATE2TRANSFER" (see below for example)
        - For each transaction in Transactions
            - Get user by pubkey (`transaction.ToPublicKey`) from postgres
            - Add state id to user (`transaction.ToStateID`)

**Endpoints for the app:**

- GetTransactions:
    - Returns all transactions for the logged in user (sender or receiver)
    - Each transaction contains sender user object and timestamp
        - To be discussed: For privacy reasons I think we should only return the worldcoin tag, not the full user object. If both people are friends we could also show the profile picture
- Let's skip the balance endpoint for now. The app can just subscribe to the transaction graphql endpoint and for every update fetch the new balance from the crypto backend

**API response examples:**

Transfer example: 

```json
{
    "jsonrpc": "2.0",
    "id": 1,
    "result": {
        "ID": 4372,
        "Type": "TRANSFER",
        "FeeReceiver": 0,
        "CombinedSignature": "0x10b9d20858b240cf00f0abc2c7de3e31e53b60fd1c42f336eb4e25909b1abcea1c36d144ccdd1509b5fd4b1b9a60d939dddc49dbc9acb907b586442c500608d1",
        "PostStateRoot": "0x401b2c83ff1bd8ffc95b61ad45cc725530383d4a7fe4aeb93062cb744f3ff6f9",
        "AccountTreeRoot": "0x01e4208fe9592b9209e46501587f895f67bc0cd641e8d99d5a288939b8b737d6",
        "IncludedInBatch": "0x0e3a4491f72d7d7f72c7a84be79ceb60b557228572a22ee7b4423e51a0842686",
        "Status": "FINALISED",
        "Transactions": [
            {
                "Hash": "0x530d7738f30b5c6af498cffe837a7b10ef09a746382a0fa0ca2153dc87fe1249",
                "FromStateID": 21,
                "Amount": "1000000000",
                "Fee": "1",
                "Nonce": "0",
                "Signature": "0x10b9d20858b240cf00f0abc2c7de3e31e53b60fd1c42f336eb4e25909b1abcea1c36d144ccdd1509b5fd4b1b9a60d939dddc49dbc9acb907b586442c500608d1",
                "ToStateID": 20
            }
        ]
    }
}
```

Create2 example:

```json
{
    "jsonrpc": "2.0",
    "id": 1,
    "result": {
        "ID": 4370,
        "Type": "CREATE2TRANSFER",
        "FeeReceiver": 0,
        "CombinedSignature": "0x301cfb8422718c8198e67bdb6649ea79f8b303c7ddb7c67ece5cc2318aadd57c037f18ee3cc874e78e188eccbbfafd6fc0297b197d971ab6cb09202797ade3da",
        "PostStateRoot": "0x00eb8c6f37b7f1ef468a4dc7bea4c4a2fb5c7db0f0c297e09fa3b50830ef8907",
        "AccountTreeRoot": "0x92e0db4ddcfb814e75fcca8305c1ae7f870ee6fcadc1fb7228e831a910bdbb43",
        "IncludedInBatch": "0xe2b07b71e2f033604ef4d6888d9ee46a6bf66e99e6ab16ed2e716d2f4724afbb",
        "Status": "FINALISED",
        "Transactions": [
            {
                "Hash": "0x5d5d4629362cdd5336f7894b8513146984c80e5b81f9be7bd6ae97915bb11e50",
                "FromStateID": 2,
                "Amount": "100000000000",
                "Fee": "1",
                "Nonce": "718",
                "Signature": "0x301cfb8422718c8198e67bdb6649ea79f8b303c7ddb7c67ece5cc2318aadd57c037f18ee3cc874e78e188eccbbfafd6fc0297b197d971ab6cb09202797ade3da",
                "ToStateID": 20,
                "ToPublicKey": "0x15763d54124aec3aaf274f1339cbf1e64068305a9730b605160b769f26855cb41fc742984744ac47ac94edad9b4dbdecb82dd59dbe97691fee0c55f8a055917e20549efb8880ab2b2a346fae952d5bada5482500e40031d3a7253501f05543a11695bb0fad5ba3a2f85f40ae1e1007d353a2bd552a669a2ae0337ae68fe2adb4"
            }
        ]
    }
}
```