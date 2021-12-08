package p2p

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type TestRpc struct {

}

type IntParam struct {
	value int
}

func (t *TestRpc) Double(arg IntParam) (*IntParam, error) {
	return &IntParam{arg.value * 2}, nil
}

func TestPeer(t *testing.T) {
	alice, err := NewPeerWithRandomKey(0, func(conn Connection) {
		err := conn.server.RegisterName("test", TestRpc{})
		if err != nil {
			panic(err)
		}
	})
	require.NoError(t, err)

	bob, err := NewPeerWithRandomKey(0, func(conn Connection) {
		var res int
		err := conn.client.Call("test_Double", 3, &res)
		if err != nil {
			panic(err)
		}
	})
	require.NoError(t, err)

	addr := alice.ListenAddr()
	err = bob.Dial(addr)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	//require.NoError(t, alice.Close())
	//require.NoError(t, bob.Close())
}
