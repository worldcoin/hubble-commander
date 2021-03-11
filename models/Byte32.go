package models

import (
	"encoding/hex"
	"errors"
)

type Bytes32 [32]byte

func MakeBytes32(hexString string) (Bytes32, error) {
	var bytes [32]byte
	bytesSlice, err := hex.DecodeString(hexString)
	if err != nil {
		return Bytes32{}, err
	}
	if len(bytesSlice) != 32 {
		return Bytes32{}, errors.New("invalid hex string length")
	}
	copy(bytes[:], bytesSlice)
	return bytes, nil
}
