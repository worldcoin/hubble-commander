package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAPI_EnableBatchCreation(t *testing.T) {
	var enabled bool
	api := API{enableBatchCreation: func(enable bool) {
		enabled = enable
	}}

	api.EnableBatchCreation(true)
	require.Equal(t, true, enabled)

	api.EnableBatchCreation(false)
	require.Equal(t, false, enabled)
}
