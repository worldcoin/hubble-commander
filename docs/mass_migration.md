# 🎒 Mass Migrations

## Withdrawal flow

1. User signs and sends a L2 MM tx with spokeId=WithdrawManager
2. Commander builds and submits a batch with that tx
    1. Wait for the batch to be finalized
3. Someone (anyone) calls `WithdrawManager.processWithdrawCommitment`
4. User (from their ethereum wallet) calls `WithdrawManager.claimTokens`
    1. Q: Does the Worldcoin mobile app have user's ethereum wallet or do we build a dapp so you can make this call from metamask (or other wallets)?

```go
struct MassMigrationBody {
    bytes32 accountRoot;
    uint256[2] signature;
    uint256 spokeID;
    bytes32 withdrawRoot; // root of a merkle tree made up of withdrawLeaves
    uint256 tokenID;
    uint256 amount; // total amount of withdrawals without fees
    uint256 feeReceiver;
    bytes txs;
}
```

The withdrawLeaves are new created UserStates, each with nonce=0 and balances corresponding to the amounts deducted from the "migrated" accounts.

```go
struct MassMigration {
    uint256 fromIndex;
    uint256 amount;
    uint256 fee;
}
```

TODO:

- What is SpokeID?