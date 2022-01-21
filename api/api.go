package api

import (
	"fmt"
	"net/http"

	"github.com/Worldcoin/hubble-commander/api/middleware"
	"github.com/Worldcoin/hubble-commander/api/rpc"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type API struct {
	cfg               *config.APIConfig
	storage           *st.Storage
	client            *eth.Client
	mockSignature     models.Signature
	commanderMetrics  *metrics.CommanderMetrics
	disableSignatures bool
}

func NewServer(
	cfg *config.Config,
	storage *st.Storage,
	client *eth.Client,
	commanderMetrics *metrics.CommanderMetrics,
) (*http.Server, error) {
	server, err := getAPIServer(cfg.API, storage, client, commanderMetrics, cfg.Rollup.DisableSignatures)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	if log.IsLevelEnabled(log.DebugLevel) {
		mux.Handle("/", middleware.Logger(server, commanderMetrics))
	} else {
		mux.Handle("/", middleware.DefaultHandler(server, commanderMetrics))
	}

	addr := fmt.Sprintf(":%s", cfg.API.Port)
	return &http.Server{Addr: addr, Handler: mux}, nil
}

func getAPIServer(
	cfg *config.APIConfig,
	storage *st.Storage,
	client *eth.Client,
	commanderMetrics *metrics.CommanderMetrics,
	disableSignatures bool,
) (*rpc.Server, error) {
	api := API{
		cfg:               cfg,
		storage:           storage,
		client:            client,
		commanderMetrics:  commanderMetrics,
		disableSignatures: disableSignatures,
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
	domain, err := a.client.GetDomain()
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
