package models

import (
	"encoding/binary"

	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v3"
)

// Encode Remember to provide cases for both value and pointer types when adding new encoders
func Encode(value interface{}) ([]byte, error) {
	switch v := value.(type) {
	case MerklePath:
		return v.Bytes(), nil
	case *MerklePath:
		return nil, errors.Errorf("pass by value")
	case StateNode:
		return EncodeDataHash(&v)
	case *StateNode:
		return nil, errors.Errorf("pass by value")
	case FlatStateLeaf:
		return v.Bytes(), nil
	case *FlatStateLeaf:
		return nil, errors.Errorf("pass by value")
	case StateUpdate:
		return v.Bytes(), nil
	case *StateUpdate:
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
	case *MerklePath:
		return v.SetBytes(data)
	case *StateNode:
		return DecodeDataHash(data, v)
	case *FlatStateLeaf:
		return v.SetBytes(data)
	case *StateUpdate:
		return v.SetBytes(data)
	case *uint32:
		return DecodeUint32(data, v)
	case *uint64:
		return DecodeUint64(data, v)
	default:
		return bh.DefaultDecode(data, value)
	}
}

func EncodeDataHash(node *StateNode) ([]byte, error) {
	return node.DataHash.Bytes(), nil
}

func EncodeUint32(number *uint32) ([]byte, error) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b[0:4], *number)
	return b, nil
}

func DecodeDataHash(data []byte, node *StateNode) error {
	node.DataHash.SetBytes(data)
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
