package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyMap(t *testing.T) {
	m1 := map[string]interface{}{
		"a": "bbb",
		"b": map[string]interface{}{
			"c": 123,
		},
	}

	m2 := CopyMap(m1)

	m1["a"] = "zzz"
	delete(m1, "b")

	require.Equal(t, map[string]interface{}{"a": "zzz"}, m1)
	require.Equal(t, map[string]interface{}{
		"a": "bbb",
		"b": map[string]interface{}{
			"c": 123,
		},
	}, m2)
}
