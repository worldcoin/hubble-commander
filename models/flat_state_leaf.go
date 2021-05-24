package models

import (
	"encoding/binary"
	"reflect"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/timshannon/badgerhold/v3"
)

var flatStateLeafT = reflect.TypeOf(FlatStateLeaf{})

type FlatStateLeaf struct {
	StateID    uint32
	DataHash   common.Hash
	PubKeyID   uint32
	TokenIndex Uint256 // TODO: Consider adding a tuple index of (Pubkey; tokenIdx)
	Balance    Uint256
	Nonce      Uint256
}

func MakeFlatStateLeaf(leaf *StateLeaf) FlatStateLeaf {
	return FlatStateLeaf{
		StateID:    leaf.StateID,
		DataHash:   leaf.DataHash,
		PubKeyID:   leaf.PubKeyID,
		TokenIndex: leaf.TokenIndex,
		Balance:    leaf.Balance,
		Nonce:      leaf.Nonce,
	}
}

func (l *FlatStateLeaf) StateLeaf() *StateLeaf {
	return &StateLeaf{
		StateID:  l.StateID,
		DataHash: l.DataHash,
		UserState: UserState{
			PubKeyID:   l.PubKeyID,
			TokenIndex: l.TokenIndex,
			Balance:    l.Balance,
			Nonce:      l.Nonce,
		},
	}
}

func (l *FlatStateLeaf) Bytes() []byte {
	b := make([]byte, 136)
	binary.BigEndian.PutUint32(b[0:4], l.StateID)
	copy(b[4:36], l.DataHash[:])
	binary.BigEndian.PutUint32(b[36:40], l.PubKeyID)
	copy(b[40:72], utils.PadLeft(l.TokenIndex.Bytes(), 32))
	copy(b[72:104], utils.PadLeft(l.Balance.Bytes(), 32))
	copy(b[104:136], utils.PadLeft(l.Nonce.Bytes(), 32))
	return b
}

func (l *FlatStateLeaf) SetBytes(data []byte) error {
	l.StateID = binary.BigEndian.Uint32(data[0:4])
	l.DataHash.SetBytes(data[4:36])
	l.PubKeyID = binary.BigEndian.Uint32(data[36:40])
	l.TokenIndex.SetBytes(data[40:72])
	l.Balance.SetBytes(data[72:104])
	l.Nonce.SetBytes(data[104:136])
	return nil
}

// nolint:gocritic
// Type implements badgerhold.Storer
func (l FlatStateLeaf) Type() string {
	return flatStateLeafT.Name()
}

// nolint:gocritic
// Indexes implements badgerhold.Storer
func (l FlatStateLeaf) Indexes() map[string]badgerhold.Index {
	return map[string]badgerhold.Index{
		"Combined": {
			IndexFunc: PubKeyIDIndex,
			Unique:    false,
		},
	}
}

func PubKeyIDIndex(_ string, value interface{}) ([]byte, error) {
	leaf, ok := value.(FlatStateLeaf)
	if !ok {
		return nil, errors.New("invalid type for FlatStateLeaf index")
	}
	index := &StateLeafIndex{
		PubKeyID:   leaf.PubKeyID,
		TokenIndex: leaf.TokenIndex,
	}
	return Encode(index)
}

type StateLeafIndex struct {
	PubKeyID   uint32
	TokenIndex Uint256
}

func (c *StateLeafIndex) Bytes() []byte {
	b := make([]byte, 36)
	binary.BigEndian.PutUint32(b[0:4], c.PubKeyID)
	copy(b[4:36], utils.PadLeft(c.TokenIndex.Bytes(), 32))
	return b
}

func (c *StateLeafIndex) SetBytes(data []byte) error {
	c.PubKeyID = binary.BigEndian.Uint32(data[0:4])
	c.TokenIndex.SetBytes(data[4:36])
	return nil
}
