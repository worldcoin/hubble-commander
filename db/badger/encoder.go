package badger

import (
	bh "github.com/timshannon/badgerhold/v3"
)

func Encode(value interface{}) ([]byte, error) {
	encode, err := bh.DefaultEncode(value)
	return encode, err
}

func Decode(data []byte, value interface{}) error {
	err := bh.DefaultDecode(data, value)
	return err
}
