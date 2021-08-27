package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CopyStringUint32Map(t *testing.T) {
	m1 := map[string]uint32{
		"a": 123,
		"b": 456,
	}

	m2 := CopyStringUint32Map(m1)

	m1["a"] = 999
	delete(m1, "b")

	require.Equal(t, map[string]uint32{"a": 999}, m1)
	require.Equal(t, map[string]uint32{
		"a": 123,
		"b": 456,
	}, m2)
}
