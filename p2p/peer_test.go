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
	return IntParam{arg.value * 2}
}

func TestPeer(t *testing.T) {
	fmt.Println("# Test start")
	alice, err := NewPeerWithRandomKey(0, func(conn Connection) {
		fmt.Println("# Alice 1")
		err := conn.server.RegisterName("test", &TestRpc{})
		if err != nil {
			panic(err)
		}
	})
	require.NoError(t, err)
	fmt.Println("Alice id:", alice.host.ID())
	fmt.Println("Alice addr:", alice.host.Addrs())

	done := make(chan bool)

	bob, err := NewPeerWithRandomKey(0, func(conn Connection) {
		fmt.Println("# Bob 1")
		var res IntParam
		err := conn.client.Call("test_double", IntParam{3}, &res)
		fmt.Println(res.value)
		if err != nil {
			log.Fatal(err)
		}

		done <- true
	})
	require.NoError(t, err)
	fmt.Println("Bob id:", bob.host.ID())
	fmt.Println("Bob addr:", bob.host.Addrs())

	fmt.Println("# Connect")
	addr := alice.ListenAddr()
	err = bob.Dial(addr)
	require.NoError(t, err)

	fmt.Println("# Wait done")
	<-done
	fmt.Println("# Done")

	//require.NoError(t, alice.Close())
	//require.NoError(t, bob.Close())
}
