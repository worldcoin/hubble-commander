package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testType struct {
	First  int
	Second string
}

func TestValueToInterfaceSlice(t *testing.T) {
	input := make([]testType, 5)
	expected := make([]interface{}, 5)
	for i := range input {
		input[i] = testType{
			First: i,
		}
		expected[i] = i
	}

	output := ValueToInterfaceSlice(input, "First")
	require.Equal(t, expected, output)
}

func TestValueToInterfaceSlice_InvalidType(t *testing.T) {
	input := testType{
		First: 1,
	}
	require.Panics(t, func() {ValueToInterfaceSlice(input, "First")})
}
