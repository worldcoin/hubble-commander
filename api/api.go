package api

import (
	"fmt"
	"net/http"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/ethereum/go-ethereum/rpc"
)

type Api struct {
	cfg     *config.Config
	storage *db.Storage
}

func StartApiServer(cfg *config.Config) error {
	storage, err := db.NewStorage(cfg)
	if err != nil {
		return err
	}

	server, err := getApiServer(cfg, storage)
	if err != nil {
		return err
	}

	http.HandleFunc("/", server.ServeHTTP)
	addr := fmt.Sprintf(":%s", cfg.Port)
	return http.ListenAndServe(addr, nil)
}

func getApiServer(cfg *config.Config, storage *db.Storage) (*rpc.Server, error) {
	api := Api{cfg, storage}
	server := rpc.NewServer()

	err := server.RegisterName("hubble", &api)
	if err != nil {
		return nil, err
	}

	return server, nil
}
