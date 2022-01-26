package api

import (
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAPIConfigure_BatchCreation(t *testing.T) {
	var enabled bool
	api := API{enableBatchCreation: func(enable bool) {
		enabled = enable
	}}

	err := api.Configure(ConfigureParams{
		CreateBatches: ref.Bool(true),
	})
	require.NoError(t, err)
	require.Equal(t, true, enabled)

	err = api.Configure(ConfigureParams{
		CreateBatches: ref.Bool(false),
	})
	require.NoError(t, err)
	require.Equal(t, false, enabled)
}

func TestAPIConfigure_AcceptingTransactions(t *testing.T) {
	api := API{}

	err := api.Configure(ConfigureParams{
		AcceptTransactions: ref.Bool(true),
	})
	require.NoError(t, err)
	require.Equal(t, true, api.isAcceptingTransactions)

	err = api.Configure(ConfigureParams{
		AcceptTransactions: ref.Bool(false),
	})
	require.NoError(t, err)
	require.Equal(t, false, api.isAcceptingTransactions)
}
