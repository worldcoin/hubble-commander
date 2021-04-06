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

func StartAPIServer(cfg *config.Config, eth *eth.Client) error {
	storage, err := st.NewStorage(&cfg.DB)
	if err != nil {
		return err
	}

	server, err := getAPIServer(&cfg.API, storage, eth)
	if err != nil {
		return err
	}

	http.HandleFunc("/", server.ServeHTTP)
	addr := fmt.Sprintf(":%s", cfg.API.Port)
	return http.ListenAndServe(addr, nil)
}

func getAPIServer(cfg *config.APIConfig, storage *st.Storage, eth *eth.Client) (*rpc.Server, error) {
	api := API{
		cfg:     cfg,
		storage: storage,
		client:  eth,
	}
	server := rpc.NewServer()

	err := server.RegisterName("hubble", &api)
	if err != nil {
		return nil, err
	}

	return server, nil
}
