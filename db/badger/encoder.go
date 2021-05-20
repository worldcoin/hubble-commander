package badger

import (
	"github.com/Worldcoin/hubble-commander/models"
	bh "github.com/timshannon/badgerhold/v3"
)

func Encode(value interface{}) ([]byte, error) {
	switch v := value.(type) {
	case *models.MerklePath:
		return v.Bytes(), nil
	case *models.StateNode:
		return EncodeDataHash(v)
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
	default:
		return bh.DefaultDecode(data, value)
	}
}

func EncodeDataHash(node *models.StateNode) ([]byte, error) {
	return node.DataHash.Bytes(), nil
}

func DecodeDataHash(data []byte, node *models.StateNode) error {
	node.DataHash.SetBytes(data)
	return nil
}
