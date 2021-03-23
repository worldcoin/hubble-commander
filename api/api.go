package api

import (
	"fmt"
	"net/http"

	"github.com/Worldcoin/hubble-commander/config"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/rpc"
)

type API struct {
	cfg     *config.APIConfig
	storage *st.Storage
}

func StartAPIServer(cfg *config.Config) error {
	storage, err := st.NewStorage(&cfg.DB)
	if err != nil {
		return err
	}

	server, err := getAPIServer(&cfg.API, storage)
	if err != nil {
		return err
	}

	http.HandleFunc("/", server.ServeHTTP)
	addr := fmt.Sprintf(":%s", cfg.API.Port)
	return http.ListenAndServe(addr, nil)
}

func getAPIServer(cfg *config.APIConfig, storage *st.Storage) (*rpc.Server, error) {
	api := API{cfg, storage}
	server := rpc.NewServer()

	err := server.RegisterName("hubble", &api)
	if err != nil {
		return nil, err
	}

	return server, nil
}
