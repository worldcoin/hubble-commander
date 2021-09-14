# Claim procedure design

## Phillip Rooftop Version

1. User generates proof on mobile
2. Aggregate proofs on server
3. Send aggregate proof to a SC
4. SC verifies the aggregate proof, blacklists the proofs (no more proofs for the same purpose from the same user)
5. Transfers the tokens to rollup contract (batch deposit)

## C2T Version

1. User generates proof on mobile
2. Aggregate proofs on server and verify them against RegistrationTree
3. Submit a C2T Batch
4. Wait 2 weeks
5. Withdraw reward from WorldcoinVault passing inclusion proofs

Inclusion proofs would need additional

- fraud proof on uniquness of userState, to avoid race conditions with other CA/Proposers

IS IT CHEAPER?

TODO: can you deposit to existing account?

## Making C2T Version optimistic

![MiscCharts-Page-3.png](Claim%20procedure%20design%208a313a102e4f4acbb281c2fb0fdefb5a/MiscCharts-Page-3.png)