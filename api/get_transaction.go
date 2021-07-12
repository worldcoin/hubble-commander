package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
)

func (a *API) GetTransaction(hash common.Hash) (interface{}, error) {
	transfer, err := a.storage.GetTransferWithBatchDetails(hash)
	if err != nil && !st.IsNotFoundError(err) {
		return nil, err
	}
	if transfer != nil {
		return a.returnTransferReceipt(transfer)
	}

	transaction, err := a.storage.GetCreate2TransferWithBatchDetails(hash)
	if err != nil {
		return nil, err
	}
	return a.returnCreate2TransferReceipt(transaction)
}

func (a *API) returnTransferReceipt(transfer *models.TransferWithBatchDetails) (*dto.TransferReceipt, error) {
	status, err := CalculateTransactionStatus(a.storage, &transfer.TransactionBase, a.storage.GetLatestBlockNumber())
	if err != nil {
		return nil, err
	}

	return &dto.TransferReceipt{
		TransferWithBatchDetails: *transfer,
		ReceiveTime:              models.NewTimestamp(transfer.ReceiveTime),
		Status:                   *status,
	}, nil
}

func (a *API) returnCreate2TransferReceipt(transfer *models.Create2TransferWithBatchDetails) (*dto.Create2TransferReceipt, error) {
	status, err := CalculateTransactionStatus(a.storage, &transfer.TransactionBase, a.storage.GetLatestBlockNumber())
	if err != nil {
		return nil, err
	}

	return &dto.Create2TransferReceipt{
		Create2TransferWithBatchDetails: *transfer,
		ReceiveTime:                     models.NewTimestamp(transfer.ReceiveTime),
		Status:                          *status,
	}, nil
}
