package api

import (
	"fmt"
	"net/http"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/rpc"
)

type API struct {
	cfg     *config.APIConfig
	storage *st.Storage
	client  *eth.Client
}

func NewAPIServer(cfg *config.APIConfig, storage *st.Storage, client *eth.Client) (*http.Server, error) {
	server, err := getAPIServer(cfg, storage, client)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", server.ServeHTTP)
	addr := fmt.Sprintf(":%s", cfg.Port)
	return &http.Server{Addr: addr}, nil
}

func getAPIServer(cfg *config.APIConfig, storage *st.Storage, client *eth.Client) (*rpc.Server, error) {
	api := API{
		cfg:     cfg,
		storage: storage,
		client:  client,
	}
	server := rpc.NewServer()

	if err := server.RegisterName("hubble", &api); err != nil {
		return nil, err
	}
	return server, nil
}
