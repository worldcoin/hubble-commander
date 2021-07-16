package models

type Account struct {
	PubKeyID  uint32
	PublicKey PublicKey `badgerhold:"index"`
}
