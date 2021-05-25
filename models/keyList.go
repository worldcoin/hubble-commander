package models

import (
	"encoding/binary"
	"reflect"

	"github.com/pkg/errors"
)

var keyListT = reflect.TypeOf([][]byte{})

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
