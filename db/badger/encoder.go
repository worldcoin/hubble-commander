package badger

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/models"
	bh "github.com/timshannon/badgerhold/v3"
)

func Encode(value interface{}) ([]byte, error) {
	switch v := value.(type) {
	case *models.MerklePath:
		return v.Bytes(), nil
	case models.StateNode:
		return EncodeDataHash(&v)
	case models.FlatStateLeaf:
		return v.Bytes(), nil
	case models.StateUpdate:
		return v.Bytes(), nil
	case uint32:
		return EncodeUint32(&v)
	default:
		return bh.DefaultEncode(value)
	}
}

func Decode(data []byte, value interface{}) error {
	switch v := value.(type) {
	case *models.MerklePath:
		return v.SetBytes(data)
	case *models.StateNode:
		return DecodeDataHash(data, v)
	case *models.FlatStateLeaf:
		return v.SetBytes(data)
	case *models.StateUpdate:
		return v.SetBytes(data)
	case *uint32:
		return DecodeUint32(data, v)
	default:
		return bh.DefaultDecode(data, value)
	}
}

func EncodeDataHash(node *models.StateNode) ([]byte, error) {
	return node.DataHash.Bytes(), nil
}

func EncodeUint32(number *uint32) ([]byte, error) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b[0:4], *number)
	return b, nil
}

func DecodeDataHash(data []byte, node *models.StateNode) error {
	node.DataHash.SetBytes(data)
	return nil
}

func DecodeUint32(data []byte, number *uint32) error {
	newUint32 := binary.BigEndian.Uint32(data)
	*number = newUint32
	return nil
}
