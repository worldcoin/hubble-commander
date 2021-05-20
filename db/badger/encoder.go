package badger

import (
	"github.com/Worldcoin/hubble-commander/models"
	bh "github.com/timshannon/badgerhold/v3"
)

func Encode(value interface{}) ([]byte, error) {
	switch v := value.(type) {
	case *models.MerklePath:
		return v.Bytes(), nil
	default:
		return bh.DefaultEncode(value)
	}
}

func Decode(data []byte, value interface{}) error {
	switch v := value.(type) {
	case *models.MerklePath:
		return v.SetBytes(data)
	default:
		return bh.DefaultDecode(data, value)
	}
}
