package utils

import (
	"reflect"
)

func ValueToInterfaceSlice(slice interface{}, fieldName string) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("ValueToInterfaceSlice() given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil
	}

	result := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		result[i] = s.Index(i).FieldByName(fieldName).Interface()
	}

	return result
}
