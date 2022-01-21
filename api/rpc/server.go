package rpc

import (
	"context"
	"net/http"

	"github.com/ethereum/go-ethereum/rpc"
)

type ContextKey int

const (
	authKeyHeader = "Auth-Key"

	AuthKey ContextKey = iota
)

// Server is an RPC server wrapper that pass additional auth header value to context.
type Server struct {
	*rpc.Server
}

func NewServer() *Server {
	return &Server{
		Server: rpc.NewServer(),
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	authHeaderValue := r.Header.Get(authKeyHeader)
	if authHeaderValue == "" {
		s.Server.ServeHTTP(w, r)
		return
	}

	ctx := context.WithValue(r.Context(), AuthKey, authHeaderValue)
	r = r.WithContext(ctx)
	s.Server.ServeHTTP(w, r)
}
