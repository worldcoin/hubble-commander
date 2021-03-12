package models

type Account struct {
	AccountIndex uint32 `db:"account_index"`
	PublicKey    []byte `db:"public_key"`
}
