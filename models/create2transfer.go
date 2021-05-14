package models

type Create2Transfer struct {
	TransactionBase
	ToStateID   uint32    `db:"to_state_id"`
	ToPublicKey PublicKey `db:"to_public_key"`
}

type Create2TransferForCommitment struct {
	TransactionBaseForCommitment
	ToStateID   uint32    `db:"to_state_id"`
	ToPublicKey PublicKey `db:"to_public_key"`
}
