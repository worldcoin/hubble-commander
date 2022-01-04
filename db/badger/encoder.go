package badger

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v4"
)

// nolint:gocyclo
// Encode Remember to provide cases for both value and pointer types when adding new encoders
func Encode(value interface{}) ([]byte, error) {
	switch v := value.(type) {
	case models.AccountNode:
		return EncodeDataHash(&v.DataHash)
	case *models.AccountNode:
		return nil, errors.Errorf("pass by value")
	case models.AccountLeaf:
		return v.Bytes(), nil
	case *models.AccountLeaf:
		return nil, errors.Errorf("pass by value")
	case models.ChainState:
		return v.Bytes(), nil
	case *models.ChainState:
		return nil, errors.Errorf("pass by value")
	case models.NamespacedMerklePath:
		return v.Bytes(), nil
	case *models.NamespacedMerklePath:
		return nil, errors.Errorf("pass by value")
	case models.MerkleTreeNode:
		return EncodeDataHash(&v.DataHash)
	case *models.MerkleTreeNode:
		return nil, errors.Errorf("pass by value")
	case models.FlatStateLeaf:
		return v.Bytes(), nil
	case *models.FlatStateLeaf:
		return nil, errors.Errorf("pass by value")
	case models.StateUpdate:
		return v.Bytes(), nil
	case *models.StateUpdate:
		return nil, errors.Errorf("pass by value")
	case uint32:
		return EncodeUint32(&v)
	case *uint32:
		return nil, errors.Errorf("pass by value")
	case uint64:
		return EncodeUint64(&v)
	case *uint64:
		return nil, errors.Errorf("pass by value")
	default:
		return bh.DefaultEncode(value)
	}
}

func Decode(data []byte, value interface{}) error {
	switch v := value.(type) {
	case *models.AccountNode:
		return DecodeDataHash(data, &v.DataHash)
	case *models.AccountLeaf:
		return v.SetBytes(data)
	case *models.ChainState:
		return v.SetBytes(data)
	case *models.NamespacedMerklePath:
		return v.SetBytes(data)
	case *models.MerkleTreeNode:
		return DecodeDataHash(data, &v.DataHash)
	case *models.FlatStateLeaf:
		return v.SetBytes(data)
	case *models.StateUpdate:
		return v.SetBytes(data)
	case *uint32:
		return DecodeUint32(data, v)
	case *uint64:
		return DecodeUint64(data, v)
	default:
		return bh.DefaultDecode(data, value)
	}
}

func EncodeDataHash(dataHash *common.Hash) ([]byte, error) {
	return dataHash.Bytes(), nil
}

func EncodeUint32(number *uint32) ([]byte, error) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b[0:4], *number)
	return b, nil
}

func DecodeDataHash(data []byte, dataHash *common.Hash) error {
	dataHash.SetBytes(data)
	return nil
}

func DecodeUint32(data []byte, number *uint32) error {
	newUint32 := binary.BigEndian.Uint32(data)
	*number = newUint32
	return nil
}

func EncodeUint64(value *uint64) ([]byte, error) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b[0:8], *value)
	return b, nil
}

func DecodeUint64(data []byte, value *uint64) error {
	newUint64 := binary.BigEndian.Uint64(data)
	*value = newUint64
	return nil
}

func DecodeKey(data []byte, key interface{}, prefix []byte) error {
	return Decode(data[len(prefix):], key)
}
