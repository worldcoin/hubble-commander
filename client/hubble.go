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
	var pendingBatches []Batch
	err := h.client.CallFor(&pendingBatches, "admin_getPendingBatches")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	batches := make([]dto.PendingBatch, 0, len(pendingBatches))
	for i := range pendingBatches {
		batches = append(batches, pendingBatches[i].ToDTO())
	}
	return batches, nil
}

func (h *hubble) GetPendingTransactions() (models.GenericTransactionArray, error) {
	pendingTxs := make([]Transaction, 0)
	err := h.client.CallFor(&pendingTxs, "admin_getPendingTransactions")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return txsToTransactionArray(pendingTxs), nil
}

func (h *hubble) GetFailedTransactions() (models.GenericTransactionArray, error) {
	failedTxs := make([]Transaction, 0)
	err := h.client.CallFor(&failedTxs, "admin_getFailedTransactions")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return txsToTransactionArray(failedTxs), nil
}
