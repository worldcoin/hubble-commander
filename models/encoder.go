package models

import (
	"encoding/binary"
	"reflect"

	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
)

type ByteEncoder interface {
	Bytes() []byte
	SetBytes(data []byte) error
}

// encoderPointer encodes given value
// if value is nil it sets first byte to 0
// if value isn't nil it sets first byte to 1
// and encodes value with Bytes() function
func encodePointer(length int, value ByteEncoder) []byte {
	b := make([]byte, length+1)
	if value == nil || reflect.ValueOf(value).IsNil() {
		return b
	}
	b[0] = 1
	copy(b[1:], value.Bytes())
	return b
}

func EncodeHashPointer(value *common.Hash) []byte {
	b := make([]byte, 33)
	if value == nil {
		return b
	}
	b[0] = 1
	copy(b[1:], value[:])
	return b
}

func DecodeHashPointer(data []byte) *common.Hash {
	if data[0] == 0 {
		return nil
	}
	return ref.Hash(common.BytesToHash(data[1:]))
}

func encodeUint32Pointer(value *uint32) []byte {
	b := make([]byte, 5)
	if value == nil {
		return b
	}
	b[0] = 1
	binary.BigEndian.PutUint32(b[1:], *value)
	return b
}

func decodeUint32Pointer(data []byte) *uint32 {
	if data[0] == 0 {
		return nil
	}
	return ref.Uint32(binary.BigEndian.Uint32(data[1:]))
}

func decodeTimestampPointer(data []byte) (*Timestamp, error) {
	if data[0] == 0 {
		return nil, nil
	}

	var timestamp Timestamp
	err := timestamp.SetBytes(data[1:])
	if err != nil {
		return nil, err
	}
	return &timestamp, nil
}
