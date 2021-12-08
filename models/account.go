package models

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
)

var AccountLeafPrefix = GetBadgerHoldPrefix(AccountLeaf{})

type AccountLeaf struct {
	PubKeyID  uint32
	PublicKey PublicKey `badgerhold:"index"`
}

type AccountNode struct {
	MerklePath MerklePath
	DataHash   common.Hash
}

func (a *AccountLeaf) Bytes() []byte {
	b := make([]byte, 132)
	binary.BigEndian.PutUint32(b[0:4], a.PubKeyID)
	copy(b[4:132], a.PublicKey.Bytes())
	return b
}

func (a *AccountLeaf) SetBytes(data []byte) error {
	a.PubKeyID = binary.BigEndian.Uint32(data[0:4])
	return a.PublicKey.SetBytes(data[4:132])
}
