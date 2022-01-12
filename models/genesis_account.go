package models

import (
	"encoding/binary"
)

const populatedGenesisAccountByteSize = 168

type RawGenesisAccount struct {
	PublicKey  string       `yaml:"publicKey"`
	PrivateKey string       `yaml:"privateKey"`
	Balance    uint64       `yaml:"balance"` //TODO-chan: remove later
	State      GenesisState `yaml:"state"`
}

type GenesisState struct {
	StateID  uint32 `yaml:"stateID"`
	PubKeyID uint32 `yaml:"pubKeyID"`
	TokenID  uint64 `yaml:"tokenID"`
	Balance  uint64 `yaml:"balance"`
	Nonce    uint64 `yaml:"nonce"`
}

func (s *GenesisState) ToStateLeaf() *StateLeaf {
	return &StateLeaf{
		StateID: s.StateID,
		UserState: UserState{
			PubKeyID: s.PubKeyID,
			TokenID:  MakeUint256(s.TokenID),
			Balance:  MakeUint256(s.Balance),
			Nonce:    MakeUint256(s.Nonce),
		},
	}
}

type GenesisAccount struct {
	PublicKey  *PublicKey
	PrivateKey *[32]byte
	Balance    Uint256
	State      *StateLeaf
}

type RegisteredGenesisAccount struct {
	GenesisAccount
	PublicKey PublicKey
	PubKeyID  uint32
}

type PopulatedGenesisAccount struct {
	PublicKey PublicKey `yaml:"public_key"`
	PubKeyID  uint32    `yaml:"pub_key_id"`
	StateID   uint32    `yaml:"state_id"`
	Balance   Uint256
}

func (a *PopulatedGenesisAccount) Bytes() []byte {
	b := make([]byte, populatedGenesisAccountByteSize)

	copy(b[:128], a.PublicKey.Bytes())
	binary.BigEndian.PutUint32(b[128:132], a.PubKeyID)
	binary.BigEndian.PutUint32(b[132:136], a.StateID)
	copy(b[136:168], a.Balance.Bytes())

	return b
}

func (a *PopulatedGenesisAccount) SetBytes(data []byte) error {
	err := a.PublicKey.SetBytes(data[:128])
	if err != nil {
		return err
	}

	a.PubKeyID = binary.BigEndian.Uint32(data[128:132])
	a.StateID = binary.BigEndian.Uint32(data[132:136])
	a.Balance.SetBytes(data[136:168])

	return nil
}
