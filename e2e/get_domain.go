package e2e

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/stretchr/testify/require"
	"github.com/ybbus/jsonrpc/v2"
)

func GetDomain(t *testing.T, client jsonrpc.RPCClient) bls.Domain {
	var info dto.NetworkInfo
	err := client.CallFor(&info, "hubble_getNetworkInfo")
	require.NoError(t, err)

	return info.SignatureDomain
}
