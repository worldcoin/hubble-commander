package p2p

import (
	"bytes"
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

type Request struct {
	Data []byte
}
type Response struct {
	Data []byte
}

type Service struct{}

type Peer struct {
	host    host.Host
	rpc     *rpc.Server
	client  *gorpc.Client
	server  *gorpc.Server
	service Service
}

type Connection struct {
	server *rpc.Server
	client *netRpc.Client
}

// Name to advertise our service on the P2P network. For connecting and service discovery.
var protocolID = protocol.ID("/worldcoin/rpc/1.0.0")
var serviceName = "RPC"
var serviceMethod = "Call" // Same as function name below

func (t *Service) Call(ctx context.Context, req Request, res *Response) error {
	sender, err := gorpc.GetRequestSender(ctx)
	if err != nil {
		return err
	}
	log.Println("Received a request from", sender)
	res.Data = req.Data
	return nil
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

	// Create a gorpc client protocol handler
	clt := gorpc.NewClient(h, protocolID)

	// Create a gorpc server protocol handler and add the RPC service
	svr := gorpc.NewServer(h, protocolID)
	svc := Service{}
	err = svr.RegisterName(serviceName, &svc)
	if err != nil {
		return nil, err
	}

	// Create a geth RPC server
	rpc := rpc.NewServer()

	p := &Peer{
		host:    h,
		rpc:     rpc,
		client:  clt,
		server:  svr,
		service: svc,
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
	var req Request
	b := make([]byte, 64)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	req.Data = b

	// Allocate reply
	var res Response

	// Execute RPC request
	log.Println("Send ping request")
	err = p.client.Call(peer, serviceName, "Call", req, &res)
	if err != nil {
		panic(err)
	}
	log.Println("Received ping reply")

	println("handleDial end")
}

func (p *Peer) Call(destination string, method string, args ...interface{}) error {

	// Register remote peer
	dest, err := p.PeerID(destination)
	if err != nil {
		return err
	}

	// Encode request as JSON-RPC
	var in bytes.Buffer
	var out bytes.Buffer
	ctx := context.Background()
	client, err := rpc.DialIO(ctx, in, out)

	// Construct requrest
	var req Request
	// TODO

	// Call
	var res Response
	err = p.client.Call(*dest, serviceName, serviceMethod, req, &res)
	if err != nil {
		return err
	}

	// TODO

	return nil
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
	return p.host.Close()
}
