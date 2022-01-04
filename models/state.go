package models

import (
	"encoding/binary"
	"reflect"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
)

var FlatStateLeafPrefix = []byte("bh_" + reflect.TypeOf(FlatStateLeaf{}).Name() + ":")

type UserState struct {
	PubKeyID uint32
	TokenID  Uint256
	Balance  Uint256
	Nonce    Uint256
}

type StateLeaf struct {
	StateID  uint32
	DataHash common.Hash
	UserState
}

type StateUpdate struct {
	ID            uint64 `badgerhold:"key"`
	CurrentRoot   common.Hash
	PrevRoot      common.Hash
	PrevStateLeaf StateLeaf
}

// nolint:gocritic
func (s UserState) Copy() *UserState {
	return &s
}

func (u *StateUpdate) Bytes() []byte {
	b := make([]byte, 208)
	binary.BigEndian.PutUint64(b[0:8], u.ID)
	copy(b[8:40], u.CurrentRoot[:])
	copy(b[40:72], u.PrevRoot[:])

	leaf := &u.PrevStateLeaf
	binary.BigEndian.PutUint32(b[72:76], leaf.StateID)
	copy(b[76:108], leaf.DataHash[:])
	binary.BigEndian.PutUint32(b[108:112], leaf.PubKeyID)
	copy(b[112:144], utils.PadLeft(leaf.TokenID.Bytes(), 32))
	copy(b[144:176], utils.PadLeft(leaf.Balance.Bytes(), 32))
	copy(b[176:208], utils.PadLeft(leaf.Nonce.Bytes(), 32))
	return b
}

func (u *StateUpdate) SetBytes(data []byte) error {
	u.ID = binary.BigEndian.Uint64(data[0:8])
	u.CurrentRoot.SetBytes(data[8:40])
	u.PrevRoot.SetBytes(data[40:72])

	u.PrevStateLeaf.StateID = binary.BigEndian.Uint32(data[72:76])
	u.PrevStateLeaf.DataHash.SetBytes(data[76:108])
	u.PrevStateLeaf.PubKeyID = binary.BigEndian.Uint32(data[108:112])
	u.PrevStateLeaf.TokenID.SetBytes(data[112:144])
	u.PrevStateLeaf.Balance.SetBytes(data[144:176])
	u.PrevStateLeaf.Nonce.SetBytes(data[176:208])
	return nil
}
