package models

type Account struct {
	PubKeyID  uint32    `db:"pub_key_id"`
	PublicKey PublicKey `db:"public_key"`
}
