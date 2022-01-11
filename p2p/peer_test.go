package p2p

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestRPC struct {
}

type IntParam struct {
	// Note that Value needs to be uppercase so `json.Marshall` can see it
	Value int
}

func (t *TestRPC) Double(arg IntParam) IntParam {
	return IntParam{arg.Value * 2}
}

func TestPeer(t *testing.T) {
	// Create node Alice serving test_double
	alice, err := NewPeerWithRandomKey(0)
	require.NoError(t, err)
	err = alice.server.RegisterName("test", &TestRPC{})
	require.NoError(t, err)

	// Create node Bob
	bob, err := NewPeerWithRandomKey(0)
	require.NoError(t, err)

	// Have Bob call Alice
	var res IntParam
	addr := alice.ListenAddr()
	client, err := bob.Dial(addr)
	require.NoError(t, err)
	err = client.Call(&res, "test_double", IntParam{3})
	require.NoError(t, err)
	require.Equal(t, 6, res.Value)
	err = client.Call(&res, "test_double", IntParam{5})
	require.NoError(t, err)
	require.Equal(t, 10, res.Value)

	// Tear down
	require.NoError(t, alice.Close())
	require.NoError(t, bob.Close())
}
