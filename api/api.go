package api

import (
	"fmt"
	"net/http"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/ethereum/go-ethereum/rpc"
)

type Api struct {
	cfg *config.Config
}

func StartApiServer(cfg *config.Config) error {
	server, err := getApiServer(cfg)
	if err != nil {
		return err
	}
	http.HandleFunc("/", server.ServeHTTP)
	addr := fmt.Sprintf(":%d", cfg.Port)
	return http.ListenAndServe(addr, nil)
}

func getApiServer(cfg *config.Config) (*rpc.Server, error) {
	api := Api{cfg}
	server := rpc.NewServer()
	err := server.RegisterName("hubble", &api)
	if err != nil {
		return nil, err
	}
	return server, nil
}
