package models

import (
	"encoding/binary"
	"reflect"

	"github.com/pkg/errors"
	bh "github.com/timshannon/badgerhold/v3"
)

const keyListTypeName = "keyList"

var keyListT = reflect.TypeOf([][]byte{})

// Encode Remember to provide cases for both value and pointer types when adding new encoders
func Encode(value interface{}) ([]byte, error) {
	if reflect.TypeOf(value).Name() == "keyList" {
		return EncodeKeyList(value)
	}
	switch v := value.(type) {
	case MerklePath:
		return v.Bytes(), nil
	case *MerklePath:
		return nil, errors.Errorf("pass by value")
	case StateNode:
		return EncodeDataHash(&v)
	case *StateNode:
		return nil, errors.Errorf("pass by value")
	case FlatStateLeaf:
		return v.Bytes(), nil
	case *FlatStateLeaf:
		return nil, errors.Errorf("pass by value")
	case StateUpdate:
		return v.Bytes(), nil
	case *StateUpdate:
		return nil, errors.Errorf("pass by value")
	case uint32:
		return EncodeUint32(&v)
	case *uint32:
		return nil, errors.Errorf("pass by value")
	}

	if reflect.TypeOf(value).ConvertibleTo(keyListT) {
		return EncodeKeyList(value)
	}
	return bh.DefaultEncode(value)
}

func Decode(data []byte, value interface{}) error {
	switch v := value.(type) {
	case *MerklePath:
		return v.SetBytes(data)
	case *StateNode:
		return DecodeDataHash(data, v)
	case *FlatStateLeaf:
		return v.SetBytes(data)
	case *StateUpdate:
		return v.SetBytes(data)
	case *uint32:
		return DecodeUint32(data, v)
	}

	if rValue, ok := isKeyListPtrType(value); ok {
		return DecodeKeyList(data, *rValue)
	}
	return bh.DefaultDecode(data, value)
}

func EncodeDataHash(node *StateNode) ([]byte, error) {
	return node.DataHash.Bytes(), nil
}

func EncodeUint32(number *uint32) ([]byte, error) {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b[0:4], *number)
	return b, nil
}

func DecodeDataHash(data []byte, node *StateNode) error {
	node.DataHash.SetBytes(data)
	return nil
}

func DecodeUint32(data []byte, number *uint32) error {
	newUint32 := binary.BigEndian.Uint32(data)
	*number = newUint32
	return nil
}

func EncodeKeyList(value interface{}) ([]byte, error) {
	keyList, ok := reflect.ValueOf(value).Convert(keyListT).Interface().([][]byte)
	if !ok {
		return nil, errors.New("keyList: cannot convert to [][]byte")
	}

	itemLength := 0
	if len(keyList) > 0 {
		itemLength = len(keyList[0])
	}

	b := make([]byte, 8, len(keyList)*itemLength+8)
	binary.BigEndian.PutUint32(b[:4], uint32(len(keyList)))
	binary.BigEndian.PutUint32(b[4:8], uint32(itemLength))
	for i := range keyList {
		if len(keyList[i]) != itemLength {
			panic("keyList: different items length")
		}
		b = append(b, keyList[i]...)
	}
	return b, nil
}

func DecodeKeyList(data []byte, value reflect.Value) error {
	if len(data) < 8 {
		return errors.New("keyList: should have at least 8 bytes")
	}

	sliceLen := binary.BigEndian.Uint32(data[:4])
	itemLen := binary.BigEndian.Uint32(data[4:8])
	if uint32(len(data)) != sliceLen*itemLen+8 {
		return errors.New("keyList: invalid data length")
	}

	keyList := make([][]byte, sliceLen)
	for i := uint32(0); i < sliceLen; i++ {
		keyList[i] = data[i*itemLen+8 : (i+1)*itemLen+8]
	}
	value.Set(reflect.Indirect(reflect.ValueOf(keyList)))
	return nil
}

func isKeyListPtrType(v interface{}) (*reflect.Value, bool) {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return nil, false
	}
	value = value.Elem()
	return &value, value.Type().AssignableTo(keyListT)
}
