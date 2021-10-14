# 💸 Transaction structure

```go
struct Transfer {
    uint256 txType;
    uint256 fromIndex;
    uint256 toIndex;
    uint256 amount;
    uint256 fee;
    uint256 nonce;
}

struct Create2Transfer {
    uint256 txType;
    uint256 fromIndex;
    uint256 toIndex;
    uint256 toPubkeyID;
    uint256 amount;
    uint256 fee;
    uint256 nonce;
}

struct MassMigration {
    uint256 txType;
    uint256 fromIndex;
    uint256 amount;
    uint256 fee;
    uint256 spokeID;
    uint256 nonce;
}

const (
	Transfer        TransactionType = 1
	Create2Transfer TransactionType = 3
	MassMigration   TransactionType = 5
)
```

Common fields:

```json
uint256 txType;
uint256 fromIndex;
uint256 amount;
uint256 fee;
uint256 nonce;
```

## Signatures

Transfer (12 bytes)

```jsx
abi.encodePacked(
    TRANSFER,      
    _tx.fromIndex, // sender state ID
    _tx.toIndex,   // recipient state ID
    nonce,
    _tx.amount,
    _tx.fee
);
```

Create2Transfer

```jsx
abi.encodePacked(
    CREATE2TRANSFER,
    _tx.fromIndex, // sender state ID
    to,            // recipient pubkey
    nonce,
    _tx.amount,
    _tx.fee
);
```