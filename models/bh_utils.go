package models

import (
	"reflect"
)

func GetTypeName(dataType interface{}) []byte {
	return []byte(reflect.TypeOf(dataType).Name())
}

func GetBadgerHoldPrefix(dataType interface{}) []byte {
	return []byte("bh_" + reflect.TypeOf(dataType).Name() + ":")
}
