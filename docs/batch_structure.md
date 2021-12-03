# 🗞 Batch Structure and Fraud Proofs

## Batch Submission

![batch_submission.png](batch_submission.png)

## Batch types

```go
const (
    Genesis         BatchType = 0
    Transfer        BatchType = 1
    MassMigration   BatchType = 2
    Create2Transfer BatchType = 3
    Deposit         BatchType = 4
)
```
