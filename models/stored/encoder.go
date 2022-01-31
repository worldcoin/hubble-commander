package stored

import (
	"encoding/binary"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
)

type ByteEncoder interface {
	Bytes() []byte
	SetBytes(data []byte) error
}

func EncodeHashPointer(value *common.Hash) []byte {
	b := make([]byte, 33)
	if value == nil {
		return b
	}
	b[0] = 1
	copy(b[1:], value.Bytes())
	return b
}

func decodeHashPointer(data []byte) *common.Hash {
	if data[0] == 0 {
		return nil
	}
	return ref.Hash(common.BytesToHash(data[1:]))
}

func EncodeUint32Pointer(value *uint32) []byte {
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

func encodeStringPointer(value *string) []byte {
	if value == nil {
		return []byte{0}
	}
	b := make([]byte, len(*value)+1)
	b[0] = 1
	copy(b[1:], *value)
	return b
}

func decodeStringPointer(data []byte) *string {
	if data[0] == 0 {
		return nil
	}
	return ref.String(string(data[1:]))
}

func encodeTimestampPointer(timestamp *models.Timestamp) []byte {
	b := make([]byte, 16)
	if timestamp == nil {
		return b
	}
	b[0] = 1
	copy(b[1:], timestamp.Bytes())
	return b
}

func decodeTimestampPointer(data []byte) (*models.Timestamp, error) {
	if data[0] == 0 {
		return nil, nil
	}

	var timestamp models.Timestamp
	err := timestamp.SetBytes(data[1:])
	if err != nil {
		return nil, err
	}
	return &timestamp, nil
}

func EncodeCommitmentIDPointer(id *models.CommitmentID) []byte {
	b := make([]byte, models.CommitmentIDDataLength+1)
	if id == nil {
		return b
	}
	b[0] = 1
	copy(b[1:], id.Bytes())
	return b
}

func decodeCommitmentIDPointer(data []byte) (*models.CommitmentID, error) {
	if data[0] == 0 {
		return nil, nil
	}

	var commitmentID models.CommitmentID
	err := commitmentID.SetBytes(data[1:])
	if err != nil {
		return nil, err
	}
	return &commitmentID, nil
}
