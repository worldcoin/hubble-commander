package models

type Create2Transfer struct {
	TransactionBase
	ToStateID  uint32 `db:"to_state_id"`
	ToPubKeyID uint32 `db:"to_pub_key_id"`
}

type Create2TransferForCommitment struct {
	TransactionBaseForCommitment
	ToStateID  uint32 `db:"to_state_id"`
	ToPubKeyID uint32 `db:"to_pub_key_id"`
}
