package models

import (
	"reflect"
)

type ByteEncoder interface {
	Bytes() []byte
	SetBytes(data []byte) error
}

func GetBadgerHoldPrefix(dataType interface{}) []byte {
	return []byte("bh_" + reflect.TypeOf(dataType).Name() + ":")
}

func GetTypeName(dataType interface{}) []byte {
	return []byte(reflect.TypeOf(dataType).Name())
}

func EncodeCommitmentIDPointer(id *CommitmentID) []byte {
	b := make([]byte, commitmentIDDataLength+1)
	if id == nil {
		return b
	}
	b[0] = 1
	copy(b[1:], id.Bytes())
	return b
}

func DecodeCommitmentIDPointer(data []byte) (*CommitmentID, error) {
	if data[0] == 0 {
		return nil, nil
	}

	var commitmentID CommitmentID
	err := commitmentID.SetBytes(data[1:])
	if err != nil {
		return nil, err
	}
	return &commitmentID, nil
}
