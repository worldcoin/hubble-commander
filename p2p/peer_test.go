package p2p

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

type TestRpc struct {

}

type IntParam struct {
	value int
}

func (t *TestRpc) Double(arg IntParam) IntParam {
	return IntParam{arg.value * 2}
}

func TestPeer(t *testing.T) {
	alice, err := NewPeerWithRandomKey(0, func(conn Connection) {
		err := conn.server.RegisterName("test", &TestRpc{})
		if err != nil {
			panic(err)
		}
	})
	require.NoError(t, err)

	bob, err := NewPeerWithRandomKey(0, func(conn Connection) {
		var res IntParam
		err := conn.client.Call("test_double", IntParam{3}, &res)
		fmt.Println(res.value)
		if err != nil {
			log.Fatal(err)
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
