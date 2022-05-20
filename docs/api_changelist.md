API Changes since v0.4
==================

This list was created retroactively and is unlikely to be exhaustive. Our architecture
makes it difficult to list all the places our behavior under errors has changed but the
error handling was significantly revamped and it is very likely that code which relies on
specific error codes or messages will need to be changed.

1. `hubble_getNetworkInfo`
  - Capitalization of all keys has changed; all keys used to be in `camelCase` but
    they are now in `TitleCase`. e.g. `latestBatch` is now `LatestBatch`.
  - The `DeploymentBlock` key has been renamed to `AccountRegistryDeploymentBlock`
  - The `TokenRegistry`, `SpokeRegistry`, `DepositManager`, `WithdrawManager`, and
    `TransactionCount` fields have been added.

2. `hubble_getUserStates`
  - If something goes wrong we now return a generic error: 99003

3. `hubble_getBatchByID`
  - `SubmissionBlock` and `SubmissionTime` were renamed to `MinedBlock` and `MinedTime`.
  - There is a new option for the `Type` field: `DEPOSIT`.

4. `hubble_getCommitment`
  - This used to accept a single parameter, the commitment ID as an integer. Commitment
    IDs are now objects with a `BatchID` key and an `IndexInBatch` key.

    Here is an example request:
    ```json
    {
        "jsonrpc": "2.0",
        "method": "hubble_getCommitment",
        "params": [
            {
                "BatchID": "1",
                "IndexInBatch": 0
            }
        ],
        "id": 1
    }
    ```

  - errors now return a generic error: 20000
  - The `Commitment` field has been removed and its fields have been inlined.
    - As an example: `.Commitment.ID` is now just `ID.
    - Specifically: these fields have been added: `ID`, `Type`, `PostStateRoot`, `LeafHash`, `TokenID`, `FeeReceivedStateID`, and `CombinedSignature`.
  - `BatchTime` was renamed to `MinedTime`.

5. `hubble_getTransaction`
  - The `IN_BATCH` status no longer exists, it has been split into `SUBMITTED` and `MINED` statuses.
  - `BatchTime` has been renamed to `MinedTime`.
  - `ToPublicKey` and `SpokeID` fields have been added.
