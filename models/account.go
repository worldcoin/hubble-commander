package models

import "github.com/ethereum/go-ethereum/common"

type Account struct {
	PubKeyID  uint32
	PublicKey PublicKey `badgerhold:"index"`
}

type AccountNode struct {
	MerklePath MerklePath
	DataHash   common.Hash
}
