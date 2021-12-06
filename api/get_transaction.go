package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
)

var getTransactionAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(10000, "transaction not found"),
}

func (a *API) GetTransaction(hash common.Hash) (interface{}, error) {
	transaction, err := a.unsafeGetTransaction(hash)
	if err != nil {
		return nil, sanitizeError(err, getTransactionAPIErrors)
	}

	return transaction, nil
}

func (a *API) unsafeGetTransaction(hash common.Hash) (interface{}, error) {
	transfer, err := a.storage.GetTransferWithBatchDetails(hash)
	if err != nil && !storage.IsNotFoundError(err) {
		return nil, err
	}
	if transfer != nil {
		return a.returnTransferReceipt(transfer)
	}

	create2transfer, err := a.storage.GetCreate2TransferWithBatchDetails(hash)
	if err != nil && !storage.IsNotFoundError(err) {
		return nil, err
	}
	if create2transfer != nil {
		return a.returnCreate2TransferReceipt(create2transfer)
	}

	massMigration, err := a.storage.GetMassMigrationWithBatchDetails(hash)
	if err != nil {
		return nil, err
	}
	return a.returnMassMigrationReceipt(massMigration)
}

func (a *API) returnTransferReceipt(transfer *models.TransferWithBatchDetails) (*dto.TransferReceipt, error) {
	status, err := CalculateTransactionStatus(a.storage, &transfer.TransactionBase, a.storage.GetLatestBlockNumber())
	if err != nil {
		return nil, err
	}

	return &dto.TransferReceipt{
		TransferWithBatchDetails: *transfer,
		Status:                   *status,
	}, nil
}

func (a *API) returnCreate2TransferReceipt(create2Transfer *models.Create2TransferWithBatchDetails) (*dto.Create2TransferReceipt, error) {
	status, err := CalculateTransactionStatus(a.storage, &create2Transfer.TransactionBase, a.storage.GetLatestBlockNumber())
	if err != nil {
		return nil, err
	}

	return &dto.Create2TransferReceipt{
		Create2TransferWithBatchDetails: *create2Transfer,
		Status:                          *status,
	}, nil
}

func (a *API) returnMassMigrationReceipt(massMigration *models.MassMigrationWithBatchDetails) (*dto.MassMigrationReceipt, error) {
	status, err := CalculateTransactionStatus(a.storage, &massMigration.TransactionBase, a.storage.GetLatestBlockNumber())
	if err != nil {
		return nil, err
	}

	return &dto.MassMigrationReceipt{
		MassMigrationWithBatchDetails: *massMigration,
		Status:                        *status,
	}, nil
}
