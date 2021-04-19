package models

type Create2Transfer struct {
	TransactionBase
	ToStateID  uint32 `db:"to_state_id"`
	ToPubkeyID uint32 `db:"to_pubkey_id"`
}
