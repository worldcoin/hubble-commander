// +build e2e

package e2e

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
)

func Test_Commander(t *testing.T) {
	commander, err := StartCommander(StartOptions{
		Image: "ghcr.io/worldcoin/hubble-commander:latest",
	})
	require.NoError(t, err)
	defer func() {
		err = commander.Stop()
		require.NoError(t, err)
	}()

	var version string
	err = commander.Client.CallFor(&version, "hubble_getVersion")
	require.NoError(t, err)
	require.Equal(t, "dev-0.1.0", version)

	var userStates []models.ReturnUserState
	err = commander.Client.CallFor(&userStates, "hubble_getUserStates", []interface{}{models.PublicKey{1, 2, 3}})
	require.NoError(t, err)
	require.Equal(t, 1, len(userStates))
	require.Equal(t, 0, userStates[0].Nonce.Cmp(big.NewInt(0)))
}
