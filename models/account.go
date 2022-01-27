package models

import (
	"encoding/binary"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

var (
	AccountLeafName                = GetTypeName(AccountLeaf{})
	AccountLeafPrefix              = GetBadgerHoldPrefix(AccountLeaf{})
	errInvalidAccountLeafIndexType = fmt.Errorf("invalid models.AccountLeaf index type")
)

type AccountLeaf struct {
	PubKeyID  uint32
	PublicKey PublicKey `badgerhold:"index"`
}

type AccountNode struct {
	MerklePath MerklePath
	DataHash   common.Hash
}

// Type implements badgerhold.Storer
func (a AccountLeaf) Type() string {
	return string(AccountLeafName)
}

// Indexes implements badgerhold.Storer
// We're adding a lot of Accounts with ZeroPublicKey to the database as a result of the way C2T processing is currently implemented.
// See the usages of ZeroPublicKey for more context.
func (a AccountLeaf) Indexes() map[string]bh.Index {
	return map[string]bh.Index{
		"PublicKey": {
			IndexFunc: func(_ string, value interface{}) ([]byte, error) {
				v, err := interfaceToAccountLeaf(value)
				if err != nil {
					return nil, err
				}
				if v.PublicKey == ZeroPublicKey {
					return nil, nil
				}
				return v.PublicKey.Bytes(), nil
			},
		},
	}
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

func interfaceToAccountLeaf(value interface{}) (*AccountLeaf, error) {
	p, ok := value.(*AccountLeaf)
	if ok {
		return p, nil
	}
	v, ok := value.(AccountLeaf)
	if ok {
		return &v, nil
	}
	return nil, errors.WithStack(errInvalidAccountLeafIndexType)
}
