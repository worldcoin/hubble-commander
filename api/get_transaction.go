package api

import (
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/ethereum/go-ethereum/common"
)

func (a *API) GetTransaction(hash common.Hash) (*dto.TransferReceipt, error) {
	transfer, err := a.storage.GetTransfer(hash)
	if err != nil {
		return nil, err
	}

	status, err := CalculateTransferStatus(a.storage, transfer, a.storage.GetLatestBlockNumber())
	if err != nil {
		return nil, err
	}

	returnTx := &dto.TransferReceipt{
		Transfer: *transfer,
		Status:   *status,
	}

	return returnTx, nil
}
