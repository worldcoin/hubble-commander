package db

import (
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v3"
)

const (
	keyListMetadataLength = 8
	bhIndexPrefix         = "_bhIndex"
)

var (
	errInconsistentItemsLength      = fmt.Errorf("inconsistent KeyList items length")
	errInvalidKeyListLength         = fmt.Errorf("invalid KeyList data length")
	errInvalidKeyListMetadataLength = fmt.Errorf("invalid KeyListMetadata data length")
)

// EncodeKeyList Format: [numKeys][keyLength][1st key ...][2nd key ...]...
func EncodeKeyList(value *bh.KeyList) ([]byte, error) {
	metadata := KeyListMetadata{ListLen: uint32(len(*value))}
	if metadata.ListLen == 0 {
		return metadata.Bytes(), nil
	}

	metadata.ItemLen = uint32(len((*value)[0]))
	b := make([]byte, metadata.GetKeyListByteLength())
	copy(b[0:keyListMetadataLength], metadata.Bytes())

	bp := keyListMetadataLength
	for i := range *value {
		if uint32(len((*value)[i])) != metadata.ItemLen {
			return nil, errors.WithStack(errInconsistentItemsLength)
		}
		bp += copy(b[bp:], (*value)[i])
	}
	return b, nil
}

func DecodeKeyList(data []byte, value *bh.KeyList) error {
	var metadata KeyListMetadata
	err := metadata.SetBytes(data)
	if err != nil {
		return err
	}
	if metadata.ListLen == 0 {
		return nil
	}

	if metadata.GetKeyListByteLength() != len(data) {
		return errInvalidKeyListLength
	}

	*value = make([][]byte, metadata.ListLen)
	index := keyListMetadataLength
	for i := range *value {
		(*value)[i] = make([]byte, metadata.ItemLen)
		index += copy((*value)[i], data[index:index+int(metadata.ItemLen)])
	}
	return nil
}

type KeyListMetadata struct {
	ListLen uint32
	ItemLen uint32
}

func (m *KeyListMetadata) Bytes() []byte {
	b := make([]byte, keyListMetadataLength)
	binary.BigEndian.PutUint32(b[0:4], m.ListLen)
	binary.BigEndian.PutUint32(b[4:8], m.ItemLen)
	return b
}

func (m *KeyListMetadata) SetBytes(data []byte) error {
	if len(data) < keyListMetadataLength {
		return errInvalidKeyListMetadataLength
	}
	m.ListLen = binary.BigEndian.Uint32(data[0:4])
	m.ItemLen = binary.BigEndian.Uint32(data[4:8])
	return nil
}

func (m *KeyListMetadata) GetKeyListByteLength() int {
	return int(keyListMetadataLength + m.ListLen*m.ItemLen)
}

func IndexKeyPrefix(typeName []byte, indexName string) []byte {
	return []byte(bhIndexPrefix + ":" + string(typeName) + ":" + indexName)
}

func IndexKey(typeName []byte, indexName string, value []byte) []byte {
	return append(IndexKeyPrefix(typeName, indexName), value...)
}
