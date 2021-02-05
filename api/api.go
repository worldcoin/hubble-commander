package api

import (
	"net/http"

	"github.com/ethereum/go-ethereum/rpc"
)

type Api struct{}

func StartApiServer(addr string) error {
	server, err := getApiServer()
	if err != nil {
		return err
	}
	http.HandleFunc("/", server.ServeHTTP)
	return http.ListenAndServe(addr, nil)
}

func getApiServer() (*rpc.Server, error) {
	server := rpc.NewServer()
	err := server.RegisterName("hubble", new(Api))
	if err != nil {
		return nil, err
	}
	return server, nil
}
