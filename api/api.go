package api

import (
	"fmt"
	"net"
	"net/http"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/middleware"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/etclabscore/go-openrpc-reflect/examples"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"

	go_openrpc_reflect "github.com/etclabscore/go-openrpc-reflect"
	meta_schema "github.com/open-rpc/meta-schema"
)

type API struct {
	cfg           *config.APIConfig
	storage       *st.Storage
	client        *eth.Client
	mockSignature models.Signature
}

func NewAPIServer(cfg *config.APIConfig, storage *st.Storage, client *eth.Client) (*http.Server, error) {
	server, err := getAPIServer(cfg, storage, client)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	if cfg.DevMode {
		mux.Handle("/", middleware.Logger(server))
	} else {
		mux.HandleFunc("/", server.ServeHTTP)
	}

	addr := fmt.Sprintf(":%s", cfg.Port)
	return &http.Server{Addr: addr, Handler: mux}, nil
}

func getAPIServer(cfg *config.APIConfig, storage *st.Storage, client *eth.Client) (*rpc.Server, error) {
	api := API{
		cfg:     cfg,
		storage: storage,
		client:  client,
	}
	if err := api.initSignature(); err != nil {
		return nil, errors.WithMessage(err, "failed to create mock signature")
	}
	server := rpc.NewServer()

	if err := server.RegisterName("hubble", &api); err != nil {
		return nil, err
	}

	rpcDiscoverService := createOpenRPCDiscoverService(&api)

	err := server.RegisterName("rpc", rpcDiscoverService)
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (a *API) initSignature() error {
	domain, err := a.storage.GetDomain(a.client.ChainState.ChainID)
	if err != nil {
		return err
	}
	wallet, err := bls.NewRandomWallet(*domain)
	if err != nil {
		return err
	}
	signature, err := wallet.Sign([]byte{1, 2, 3, 4})
	if err != nil {
		return err
	}
	a.mockSignature = *signature.ModelsSignature()
	return nil
}

func createOpenRPCDiscoverService(api *API) *examples.RPCEthereum {
	// Instantiate our document with sane defaults.
	doc := &go_openrpc_reflect.Document{}

	// Set up some minimum-viable application-specific information.
	// These are 3 fields grouped as 'Meta' in the case are server and application-specific data
	// that depend entirely on application context.
	// These fields are filled functionally, and Servers uses a lambda.
	// The fields are:
	// - Servers: describes server information like address, protocol, etc.
	// - Info: describes title, license, links, etc.
	// - ExternalDocs: links to document-level external docs.
	// This is the only place you really have to get your hands dirty.
	// Note that Servers and Info fields aren't strictly-speaking allowed to be nil for
	// an OpenRPC document to be 'valid' by spec (they're *required*), but this is just to
	// show that these are the only things that you have to actually think about
	// and we don't really care about meeting spec in a simple example.
	doc.WithMeta(&go_openrpc_reflect.MetaT{
		GetServersFn: func() func(listeners []net.Listener) (*meta_schema.Servers, error) {
			return func([]net.Listener) (*meta_schema.Servers, error) { return nil, nil }
		},
		GetInfoFn: func() (info *meta_schema.InfoObject) {
			return nil
		},
		GetExternalDocsFn: func() (exdocs *meta_schema.ExternalDocumentationObject) {
			return nil
		},
	})

	// Use a Standard reflector pattern.
	// This is a sane default supplied by the library which fits Go's net/rpc reflection conventions.
	// If you want, you can also roll your own, or edit pretty much any part of this standard object you want.
	// Highly tweakable.
	doc.WithReflector(go_openrpc_reflect.EthereumReflector)

	// Register our calculator service to the rpc.Server and rpc.Doc
	// I've grouped these together because in larger applications
	// multiple receivers may be registered on a single server,
	// and receiver registration is often done in a loop.
	// NOTE that net/rpc will log warnings like:
	//   > rpc.Register: method "BrokenReset" has 1 input parameters; needs exactly three'
	// This is because internal/fakearithmetic has spurious methods for testing this package.

	doc.RegisterReceiverName("hubble", api) // <- Register the receiver to the doc.

	// Wrap the document in a very simple default 'RPC' service, which provides one method: Discover.
	// This meets the OpenRPC specification for the service discovery endpoint to be at the reserved
	// rpc.discover endpoint.
	// You can easily roll your own Discover service if you'd like to do anything tweakable or fancy or different
	// with the document endpoint.
	rpcDiscoverService := &examples.RPCEthereum{doc}
	// (For the curious, here's what the whole of this RPC service looks like behind the scenes.)
	/*
		type RPC struct {
			Doc *Document
		}

		type RPCArg int // noop

		func (d *RPC) Discover(rpcArg *RPCArg, document *meta_schema.OpenrpcDocument) error {
			doc, err := d.Doc.Discover()
			if err != nil {
				return err
			}
			*document = *doc
			return err
		}
	*/
	return rpcDiscoverService
}
