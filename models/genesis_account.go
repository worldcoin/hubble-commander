package models

import (
	"encoding/binary"
)

const populatedGenesisAccountByteSize = 232

type GenesisAccount struct {
	PublicKey PublicKey `yaml:"public_key"`
	StateID   uint32    `yaml:"state_id"`
	State     UserState `yaml:"state"`
}

func (a *GenesisAccount) Bytes() []byte {
	b := make([]byte, populatedGenesisAccountByteSize)

	copy(b[:128], a.PublicKey.Bytes())
	binary.BigEndian.PutUint32(b[128:132], a.StateID)
	binary.BigEndian.PutUint32(b[132:136], a.State.PubKeyID)
	copy(b[136:168], a.State.TokenID.Bytes())
	copy(b[168:200], a.State.Balance.Bytes())
	copy(b[200:232], a.State.Nonce.Bytes())

	return b
}

func (a *GenesisAccount) SetBytes(data []byte) error {
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
