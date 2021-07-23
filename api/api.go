package api

import (
	"fmt"
	"net/http"

	"github.com/Worldcoin/hubble-commander/api/middleware"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pkg/errors"
)

type API struct {
	cfg           *config.APIConfig
	storage       *st.StorageBase // TODO-INTERNAL use Storage here?
	client        *eth.Client
	mockSignature models.Signature
	devMode       bool
}

func NewAPIServer(cfg *config.APIConfig, storage *st.StorageBase, client *eth.Client, devMode bool) (*http.Server, error) {
	server, err := getAPIServer(cfg, storage, client, devMode)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	if devMode {
		mux.Handle("/", middleware.Logger(server))
	} else {
		mux.HandleFunc("/", server.ServeHTTP)
	}

	addr := fmt.Sprintf(":%s", cfg.Port)
	return &http.Server{Addr: addr, Handler: mux}, nil
}

func getAPIServer(cfg *config.APIConfig, storage *st.StorageBase, client *eth.Client, devMode bool) (*rpc.Server, error) {
	api := API{
		cfg:     cfg,
		storage: storage,
		client:  client,
		devMode: devMode,
	}
	if err := api.initSignature(); err != nil {
		return nil, errors.WithMessage(err, "failed to create mock signature")
	}
	server := rpc.NewServer()

	if err := server.RegisterName("hubble", &api); err != nil {
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
