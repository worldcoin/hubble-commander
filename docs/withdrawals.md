# Withdrawals

## Flow

1. User sends a mass migration transaction with a selected spoke.
2. User waits a certain amount of time until the batch with sent transactions is finalized (finalization status can be checked
   with `hubble_getBatch`, `hubble_getCommitment` or `hubble_getTransaction`).
3. User gathers the mass migration commitment inclusion proof with `hubble_getMassMigrationCommitmentProof`.
4. User or anyone else calls `WithdrawManager.processWithdrawCommitment`.
    1. `WithdrawManager` verifies the entire commitment and marks it as processed (`WithdrawManager.claimTokens` can be now called for all
       mass migration transactions from said commitment).
    2. ERC20 tokens are transferred from the `Vault` to the `WithdrawManager` contract.
5. User gathers the withdrawal proof with `hubble_getWithdrawProof`.
6. User gathers the public key proof with `hubble_getPublicKeyProofByPubKeyID`.
7. User signs their ethereum address with their BLS private key.
8. User calls `WithdrawManager.claimTokens`.
    1. `WithdrawManager` verifies the request.
    2. ERC20 tokens are transferred from the `WithdrawManager` to the user.
