package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAPIError(t *testing.T) {
	_ = NewAPIError(1234, "something")
	require.Panics(t, func() {
		_ = NewAPIError(1234, "something")
	})
}
