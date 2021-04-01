// +build e2e

package e2e

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Commander(t *testing.T) {
	commander, err := StartCommander(StartOptions{
		Image:             "ghcr.io/worldcoin/hubble-commander:latest",
		UseHostNetworking: true, // TODO: This needs to be on when running on linux & off when running on mac.
	})
	require.NoError(t, err)

	var version string
	err = commander.Client.CallFor(&version, "hubble_getVersion", []interface{}{})
	require.NoError(t, err)

	require.Equal(t, "dev-0.1.0", version)

	err = commander.Process.Signal(os.Interrupt)
	require.NoError(t, err)
}
