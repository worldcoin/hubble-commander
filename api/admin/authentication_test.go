package admin

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/api/rpc"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/stretchr/testify/require"
)

func TestAPI_verifyAuthKey_ValidKey(t *testing.T) {
	api := API{cfg: &config.APIConfig{AuthenticationKey: authKeyValue}}

	err := api.verifyAuthKey(contextWithAuthKey(authKeyValue))
	require.NoError(t, err)
}

func TestAPI_verifyAuthKey_MissingKey(t *testing.T) {
	api := API{cfg: &config.APIConfig{AuthenticationKey: authKeyValue}}

	err := api.verifyAuthKey(context.Background())
	require.ErrorIs(t, err, errMissingAuthKey)
}

func TestAPI_verifyAuthKey_InvalidKey(t *testing.T) {
	api := API{cfg: &config.APIConfig{AuthenticationKey: authKeyValue}}

	err := api.verifyAuthKey(contextWithAuthKey("invalid key"))
	require.ErrorIs(t, err, errInvalidAuthKey)
}

func contextWithAuthKey(authKeyValue string) context.Context {
	return context.WithValue(context.Background(), rpc.AuthKey, authKeyValue)
}
