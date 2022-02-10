package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models/enums/healthstatus"
	"github.com/stretchr/testify/require"
)

func TestAPI_GetStatus_Ready(t *testing.T) {
	api := &API{isMigrating: func() bool {
		return false
	}}

	status := api.GetStatus()
	require.Equal(t, healthstatus.Ready, status)
}

func TestAPI_GetStatus_Migrating(t *testing.T) {
	api := &API{isMigrating: func() bool {
		return true
	}}

	status := api.GetStatus()
	require.Equal(t, healthstatus.Migrating, status)
}
