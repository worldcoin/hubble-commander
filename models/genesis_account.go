package models

import (
	"encoding/binary"
)

const populatedGenesisAccountByteSize = 232 // 128 + 4 + 100

type GenesisAccount struct {
	PublicKey PublicKey `yaml:"public_key"`
	StateID   uint32    `yaml:"state_id"`
	State     UserState `yaml:"state"`
}

func (a *GenesisAccount) Bytes() []byte {
	b := make([]byte, populatedGenesisAccountByteSize)

	copy(b[:128], a.PublicKey.Bytes())
	binary.BigEndian.PutUint32(b[128:132], a.StateID)
	copy(b[132:232], a.State.Bytes())

	return b
}

func (a *GenesisAccount) SetBytes(data []byte) error {
	err := a.PublicKey.SetBytes(data[:128])
	if err != nil {
		return err
	}
	err = a.State.SetBytes(data[132:232])
	if err != nil {
		return err
	}

	a.StateID = binary.BigEndian.Uint32(data[128:132])
	return nil
}
