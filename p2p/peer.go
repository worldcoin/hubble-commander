package p2p

import (
	"context"
	"crypto/rand"
	"fmt"
	"net"
	netRpc "net/rpc"
	"strconv"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/protocol"
	gostream "github.com/libp2p/go-libp2p-gostream"
	p2pstream "github.com/libp2p/go-libp2p-gostream"
	"github.com/multiformats/go-multiaddr"
)

var protocolID = protocol.ID("/worldcoin/rpc/1.0.0")

type Peer struct {
	host     host.Host
	server   *rpc.Server
	listener net.Listener
}

type Connection struct {
	server *rpc.Server
	client *netRpc.Client
}

// NewPeer creates a new transaction exchange with P2P capabilities.
// port - is the TCP port to listen for incoming P2P connections on. Pass 0 to let OS pick the port.
// privateKey - is the identity of the P2P instance
func NewPeer(port int, privateKey crypto.PrivKey) (*Peer, error) {
	// Create a listening address
	sourceMultiAddr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
	if err != nil {
		return nil, err
	}

	// Create a libp2p Host
	h, err := libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(privateKey),
	)
	if err != nil {
		return nil, err
	}

	// Start a libp2p-gostream based Geth JSON-RPC server
	server := rpc.NewServer()
	listener, _ := p2pstream.Listen(h, protocolID)
	go server.ServeListener(listener)

	p := &Peer{
		host:     h,
		server:   server,
		listener: listener,
	}

	return p, nil
}

func NewPeerWithRandomKey(port int) (*Peer, error) {
	// Creates a new RSA key pair for this host.
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	if err != nil {
		return nil, err
	}

	return NewPeer(port, prvKey)
}

func (p *Peer) Dial(destination string) (*rpc.Client, error) {
	// Register remote peer and get url
	dest, err := p.PeerID(destination)
	if err != nil {
		return nil, err
	}

	// Create a socket
	ctx := context.Background()
	conn, err := gostream.Dial(ctx, p.host, *dest, protocolID)
	if err != nil {
		return nil, err
	}

	// Create libp2p-http based Geth JSON-RPC client
	return rpc.DialIO(ctx, conn, conn)
}

func (p *Peer) PeerID(destination string) (*peer.ID, error) {
	fmt.Println("Meeting", destination)

	// Turn the destination into a multiaddr.
	maddr, err := multiaddr.NewMultiaddr(destination)
	if err != nil {
		return nil, err
	}

	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return nil, err
	}

	// Add the destination's peer multiaddress in the peerstore.
	// This will be used during connection and stream creation by libp2p.
	p.host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	return &info.ID, nil
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
	p.listener.Close()
	return p.host.Close()
}
