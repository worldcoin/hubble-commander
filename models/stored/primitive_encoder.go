package stored

import (
	"encoding/binary"

	"github.com/ethereum/go-ethereum/common"
)

func EncodeHash(hash *common.Hash) ([]byte, error) {
	return hash.Bytes(), nil
}

func DecodeHash(data []byte, hash *common.Hash) error {
	hash.SetBytes(data)
	return nil
}

func EncodeUint32(number uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b[0:4], number)
	return b
}

func DecodeUint32(data []byte, number *uint32) error {
	newUint32 := binary.BigEndian.Uint32(data)
	*number = newUint32
	return nil
}

func EncodeUint64(value uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b[0:8], value)
	return b
}

func DecodeUint64(data []byte, value *uint64) error {
	newUint64 := binary.BigEndian.Uint64(data)
	*value = newUint64
	return nil
}

func EncodeString(value string) []byte {
	return []byte(value)
}

func DecodeString(data []byte, value *string) error {
	*value = string(data)
	return nil
}
