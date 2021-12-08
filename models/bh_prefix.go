package models

import (
	"reflect"
)

func GetBadgerHoldPrefix(dataType interface{}) []byte {
	return []byte("bh_" + reflect.TypeOf(dataType).Name() + ":")
}
