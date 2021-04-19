package models

type Account struct {
	PubkeyID  uint32    `db:"pubkey_id"`
	PublicKey PublicKey `db:"public_key"`
}
