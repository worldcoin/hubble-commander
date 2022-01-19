package api

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAPIEnableBatchCreation(t *testing.T) {
	var enabled bool
	api := API{enableBatchCreation: func(enable bool) {
		enabled = enable
	}}

	err := api.EnableBatchCreation(true)
	require.NoError(t, err)
	require.Equal(t, true, enabled)

	err = api.EnableBatchCreation(false)
	require.NoError(t, err)
	require.Equal(t, false, enabled)
}
