package models

import (
	"encoding/binary"
)

const populatedGenesisAccountByteSize = 232

type RawGenesisAccount struct {
	PublicKey string       `yaml:"publicKey"`
	State     GenesisState `yaml:"state"`
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
	PublicKey PublicKey
	State     *StateLeaf
}

type PopulatedGenesisAccount struct {
	PublicKey PublicKey `yaml:"public_key"`
	StateID   uint32    `yaml:"state_id"`
	State     UserState `yaml:"state"`
}

func (a *PopulatedGenesisAccount) Bytes() []byte {
	b := make([]byte, populatedGenesisAccountByteSize)

	copy(b[:128], a.PublicKey.Bytes())
	binary.BigEndian.PutUint32(b[128:132], a.StateID)
	binary.BigEndian.PutUint32(b[132:136], a.State.PubKeyID)
	copy(b[136:168], a.State.TokenID.Bytes())
	copy(b[168:200], a.State.Balance.Bytes())
	copy(b[200:232], a.State.Nonce.Bytes())

	return b
}

func (a *PopulatedGenesisAccount) SetBytes(data []byte) error {
	err := a.PublicKey.SetBytes(data[:128])
	if err != nil {
		return err
	}

	a.StateID = binary.BigEndian.Uint32(data[128:132])
	a.State.PubKeyID = binary.BigEndian.Uint32(data[132:136])
	a.State.TokenID.SetBytes(data[136:168])
	a.State.Balance.SetBytes(data[168:200])
	a.State.Nonce.SetBytes(data[200:232])

	return nil
}
