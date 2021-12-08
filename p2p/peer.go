package p2p

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/multiformats/go-multiaddr"
	netRpc "net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"time"
)

type Peer struct {
	host host.Host
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
	*bufio.Reader
	*bufio.Writer
}

func (c conn) Close() error {
	return nil
}

func (c conn) SetWriteDeadline(time time.Time) error {
	return nil
}

func (p *Peer) handleStream(stream network.Stream) {
	fmt.Println("handleStream")

	server := rpc.NewServer()

	c := conn{
		Reader: bufio.NewReader(stream),
		Writer: bufio.NewWriter(stream),
	}

	client := jsonrpc.NewClient(c)

	p.handleConnection(Connection{
		server: server,
		client: client,
	})

	go server.ServeCodec(rpc.NewCodec(c), 0)
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

	println("Dial end")

	p.handleStream(s)

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
