package client

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/utils/consts"
	"github.com/pkg/errors"
	"github.com/ybbus/jsonrpc/v2"
)

type Hubble interface {
	GetPendingBatches() ([]dto.PendingBatch, error)
	GetPendingTransactions() (models.GenericTransactionArray, error)
	GetFailedTransactions() (models.GenericTransactionArray, error)
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

func (h *hubble) GetPendingTransactions() (models.GenericTransactionArray, error) {
	var pendingTxs models.GenericTransactionArray
	err := h.client.CallFor(&pendingTxs, "admin_getPendingTransactions")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return pendingTxs, nil
}

func (h *hubble) GetFailedTransactions() (models.GenericTransactionArray, error) {
	var failedTxs models.GenericTransactionArray
	err := h.client.CallFor(&failedTxs, "admin_getFailedTransactions")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return failedTxs, nil
}
