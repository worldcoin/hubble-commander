package api

import (
	"fmt"
	"net/http"

	"github.com/Worldcoin/hubble-commander/config"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/rpc"
)

type API struct {
	cfg     *config.Config
	storage *st.Storage
}

func StartAPIServer(cfg *config.Config) error {
	storage, err := st.NewStorage(cfg)
	if err != nil {
		return err
	}

	server, err := getAPIServer(cfg, storage)
	if err != nil {
		return err
	}

	http.HandleFunc("/", server.ServeHTTP)
	addr := fmt.Sprintf(":%s", cfg.Port)
	return http.ListenAndServe(addr, nil)
}

func getAPIServer(cfg *config.Config, storage *st.Storage) (*rpc.Server, error) {
	api := API{cfg, storage}
	server := rpc.NewServer()

	err := server.RegisterName("hubble", &api)
	if err != nil {
		return nil, err
	}

	return server, nil
}
