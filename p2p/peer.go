package p2p

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/inconshreveable/muxado"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"io"
	"log"
	netRpc "net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"sync"
	"time"
)

type Peer struct {
	host             host.Host
	handleConnection func(conn Connection)
}

type Connection struct {
	server *rpc.Server
	client *netRpc.Client
}

// NewPeer creates a new transaction exchange with P2P capabilities.
// port - is the TCP port to listen for incoming P2P connections on. Pass 0 to let OS pick the port.
// privateKey - is the identity of the P2P instance
func NewPeer(port int, handleConnection func(conn Connection), privateKey crypto.PrivKey) (*Peer, error) {
	sourceMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
	if err != nil {
		return nil, err
	}

	h, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(privateKey),
	)
	if err != nil {
		return nil, err
	}

	p := &Peer{
		host:             h,
		handleConnection: handleConnection,
	}

	h.SetStreamHandler("/worldcoin/1.0.0", func(stream network.Stream) {
		p.handleStream(stream)
	})

	return p, nil
}

func NewPeerWithRandomKey(port int, handleConnection func(conn Connection)) (*Peer, error) {
	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		return nil, err
	}

	return NewPeer(port, handleConnection, prvKey)
}

type conn struct {
	io.Reader
	io.Writer
}

func (c conn) Close() error {
	return nil
}

func (c conn) SetWriteDeadline(time time.Time) error {
	return nil
}

func (p *Peer) handleStream(stream network.Stream) {
	fmt.Println("handleStream")

	mux := muxado.Server(stream, nil)

	p.handleMuxed(mux)

	println("handleStream end")
}

func (p *Peer) handleDial(stream network.Stream) {
	fmt.Println("handleDial")

	mux := muxado.Client(stream, nil)

	time.Sleep(100 * time.Millisecond)

	p.handleMuxed(mux)

	println("handleDial end")
}

func (p *Peer) handleMuxed(mux muxado.Session) {

	conn := Connection{}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		fmt.Println("Waiting to accept")
		serverStream, err := mux.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Accepted")

		conn.server = rpc.NewServer()

		go conn.server.ServeCodec(rpc.NewCodec(serverStream), 0)
		wg.Done()
	}()

	time.Sleep(200 * time.Millisecond)

	go func() {
		clientStream, err := mux.Open()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Opened")

		conn.client = jsonrpc.NewClient(clientStream)
		wg.Done()
	}()

	wg.Wait()

	fmt.Println("handleConnection")

	p.handleConnection(conn)

	fmt.Println("handleConnection end")
}

func (p *Peer) Dial(destination string) error {
	// Turn the destination into a multiaddr.
	maddr, err := multiaddr.NewMultiaddr(destination)
	if err != nil {
		return err
	}

	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return err
	}

	// Add the destination's peer multiaddress in the peerstore.
	// This will be used during connection and stream creation by libp2p.
	p.host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	// Start a stream with the destination.
	// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
	s, err := p.host.NewStream(context.Background(), info.ID, "/worldcoin/1.0.0")
	if err != nil {
		return err
	}

	go p.handleDial(s)

	println("Dial end")
	return nil
}

func (p *Peer) ListenAddr() string {
	port := int(0)

	for _, la := range p.host.Network().ListenAddresses() {
		if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
			i, _ := strconv.ParseInt(p, 10, 32)
			port = int(i)
		}
	}

	return fmt.Sprintf("/ip4/127.0.0.1/tcp/%v/p2p/%s", port, p.host.ID().Pretty())
}

func (p *Peer) Close() error {
	return p.host.Close()
}
