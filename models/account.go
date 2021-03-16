package models

type Account struct {
	AccountIndex uint32    `db:"account_index"`
	PublicKey    PublicKey `db:"public_key"`
}
