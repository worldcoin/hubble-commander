package client

import (
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/utils/consts"
	"github.com/pkg/errors"
	"github.com/ybbus/jsonrpc/v2"
)

type Hubble interface {
	GetPendingBatches() ([]dto.PendingBatch, error)
}

type hubble struct {
	client jsonrpc.RPCClient
}

func NewHubble(url, authenticationKey string) Hubble {
	client := jsonrpc.NewClientWithOpts(url, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			consts.AuthKeyHeader: authenticationKey,
		},
	})

	return &hubble{
		client: client,
	}
}

func (h *hubble) GetPendingBatches() ([]dto.PendingBatch, error) {
	var pendingBatches []dto.PendingBatch
	err := h.client.CallFor(&pendingBatches, "admin_getPendingBatches")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return pendingBatches, nil
}
