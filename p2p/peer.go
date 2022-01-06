package p2p

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	netRpc "net/rpc"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-libp2p-core/protocol"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/multiformats/go-multiaddr"
)

type Peer struct {
	host             host.Host
	client           *gorpc.Client
	server           *gorpc.Server
	service          PingService
	handleConnection func(conn Connection)
}

type Connection struct {
	server *rpc.Server
	client *netRpc.Client
}

type PingArgs struct {
	Data []byte
}
type PingReply struct {
	Data []byte
}

type PingService struct{}

// Name to advertise our service on the P2P network. For connecting and service discovery.
var protocolID = protocol.ID("/worldcoin/rpc/1.0.0")

func (t *PingService) Ping(ctx context.Context, argType PingArgs, replyType *PingReply) error {
	sender, err := gorpc.GetRequestSender(ctx)
	if err != nil {
		return err
	}
	log.Println("Received a Ping call from", sender)
	replyType.Data = argType.Data
	return nil
}

// NewPeer creates a new transaction exchange with P2P capabilities.
// port - is the TCP port to listen for incoming P2P connections on. Pass 0 to let OS pick the port.
// privateKey - is the identity of the P2P instance
func NewPeer(port int, handleConnection func(conn Connection), privateKey crypto.PrivKey) (*Peer, error) {
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

	// Create a gorpc client protocol handler
	clt := gorpc.NewClient(h, protocolID)

	// Create a gorpc server protocol handler and add the RPC service
	svr := gorpc.NewServer(h, protocolID)
	svc := PingService{}
	err = svr.Register(&svc)
	if err != nil {
		return nil, err
	}

	p := &Peer{
		host:             h,
		client:           clt,
		server:           svr,
		service:          svc,
		handleConnection: handleConnection,
	}

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

func (p *Peer) handleDial(peer peer.ID) {
	fmt.Println("handleDial")

	// Construct request
	var args PingArgs
	c := 64
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	args.Data = b

	// Allocate reply
	var reply PingReply

	// Execute RPC request
	log.Println("Send ping request")
	err = p.client.Call(peer, "PingService", "Ping", args, &reply)
	if err != nil {
		panic(err)
	}
	log.Println("Received ping reply")

	println("handleDial end")
}

func (p *Peer) Dial(destination string) error {
	fmt.Println("Dialing %0", destination)

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

	// Connect to destination peer
	go p.handleDial(info.ID)

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
