package models

type Transfer struct {
	TransactionBase
	ToStateID uint32 `db:"to_state_id"`
}

type TransferForCommitment struct {
	TransactionBaseForCommitment
	ToStateID uint32 `db:"to_state_id"`
}
