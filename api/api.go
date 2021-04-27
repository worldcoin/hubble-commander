package api

import (
	"fmt"
	"log"
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

func StartAPIServer(cfg *config.Config, client *eth.Client) (*http.Server, error) {
	storage, err := st.NewStorage(&cfg.DB)
	if err != nil {
		return nil, err
	}

	server, err := getAPIServer(&cfg.API, storage, client)
	if err != nil {
		return nil, err
	}

	http.HandleFunc("/", server.ServeHTTP)
	addr := fmt.Sprintf(":%s", cfg.API.Port)
	httpServer := &http.Server{Addr: addr}
	go func() {
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Printf("%+v", err)
		}
	}()

	return httpServer, nil
}

func getAPIServer(cfg *config.APIConfig, storage *st.Storage, client *eth.Client) (*rpc.Server, error) {
	api := API{
		cfg:     cfg,
		storage: storage,
		client:  client,
	}
	server := rpc.NewServer()

	err := server.RegisterName("hubble", &api)
	if err != nil {
		return nil, err
	}

	return server, nil
}
