package p2p

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestRpc struct {
}

type IntParam struct {
	value int
}

func (t *TestRpc) Double(arg IntParam) IntParam {
	log.Println("RPC Called with parameter", arg.value)
	return IntParam{arg.value*2 + 1}
}

func TestPeer(t *testing.T) {
	fmt.Println("# Test start")

	// Create node Alice serving test_double
	alice, err := NewPeerWithRandomKey(0)
	require.NoError(t, err)
	err = alice.server.RegisterName("test", &TestRpc{})
	require.NoError(t, err)
	fmt.Println("Alice id:", alice.host.ID())
	fmt.Println("Alice addr:", alice.host.Addrs())

	// Create node Bob
	bob, err := NewPeerWithRandomKey(0)
	require.NoError(t, err)
	fmt.Println("Bob id:", bob.host.ID())
	fmt.Println("Bob addr:", bob.host.Addrs())

	// Have Bob call Alice
	var res IntParam
	addr := alice.ListenAddr()
	client, err := bob.Dial(addr)
	require.NoError(t, err)
	err = client.Call(&res, "test_double", IntParam{3})
	require.NoError(t, err)
	fmt.Println(res.value)
	err = client.Call(&res, "test_double", IntParam{5})
	require.NoError(t, err)
	fmt.Println(res.value)

	// Tear down
	require.NoError(t, alice.Close())
	require.NoError(t, bob.Close())
}
