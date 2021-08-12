package models

import (
	"encoding/binary"
	"reflect"

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

func encodeHashPointer(value *common.Hash) []byte {
	b := make([]byte, 33)
	if value == nil {
		return b
	}
	b[0] = 1
	copy(b[1:], value[:])
	return b
}

func encodeUint32Pointer(value *uint32) []byte {
	b := make([]byte, 33)
	if value == nil {
		return b
	}
	b[0] = 1
	binary.BigEndian.PutUint32(b[1:], *value)
	return b
}
