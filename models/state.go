package models

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
)

const (
	userStateLength   = 100 // 4 + 32 + 32 + 32
	stateUpdateLength = 208 // 4 + 32 + 32 + 100
)

var StateUpdatePrefix = GetBadgerHoldPrefix(StateUpdate{})

type UserState struct {
	PubKeyID uint32  `yaml:"pub_key_id"`
	TokenID  Uint256 `yaml:"token_id"`
	Balance  Uint256 `yaml:"balance"`
	Nonce    Uint256 `yaml:"nonce"`
}

func (s *UserState) Bytes() []byte {
	b := make([]byte, 100)

	binary.BigEndian.PutUint32(b[:4], s.PubKeyID)
	copy(b[4:36], s.TokenID.Bytes())
	copy(b[36:68], s.Balance.Bytes())
	copy(b[68:100], s.Nonce.Bytes())

	return b
}

func (s *UserState) SetBytes(data []byte) error {
	if len(data) != userStateLength {
		return ErrInvalidLength
	}

	s.PubKeyID = binary.BigEndian.Uint32(data[:4])
	s.TokenID.SetBytes(data[4:36])
	s.Balance.SetBytes(data[36:68])
	s.Nonce.SetBytes(data[68:100])

	return nil
}

type StateLeaf struct {
	StateID  uint32
	DataHash common.Hash `json:"-"`
	UserState
}

type StateUpdate struct {
	ID            uint64 `badgerhold:"key"`
	CurrentRoot   common.Hash
	PrevRoot      common.Hash
	PrevStateLeaf StateLeaf
}

//nolint:gocritic
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
	copy(b[108:208], leaf.UserState.Bytes())
	return b
}

func (u *StateUpdate) SetBytes(data []byte) error {
	if len(data) != stateUpdateLength {
		return ErrInvalidLength
	}
	err := u.PrevStateLeaf.UserState.SetBytes(data[108:208])
	if err != nil {
		return err
	}

	u.ID = binary.BigEndian.Uint64(data[0:8])
	u.CurrentRoot.SetBytes(data[8:40])
	u.PrevRoot.SetBytes(data[40:72])

	u.PrevStateLeaf.StateID = binary.BigEndian.Uint32(data[72:76])
	u.PrevStateLeaf.DataHash.SetBytes(data[76:108])
	return nil
}
