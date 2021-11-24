package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTCamelCaseToSnakeCase(t *testing.T) {
	res := CamelCaseToSnakeCase("DepositSubTreeReady")
	require.Equal(t, "deposit_sub_tree_ready", res)
}

func TestTCamelCaseToSnakeCase_FirstWordLowerCase(t *testing.T) {
	res := CamelCaseToSnakeCase("oneTwoThree")
	require.Equal(t, "one_two_three", res)
}
