package tracker

import (
	"github.com/Worldcoin/hubble-commander/models"
)

type txRequest struct {
	params          interface{}
	responseChannel chan txRequestResponse
}

type submitTransfersBatchRequestParams struct {
	batchID     *models.Uint256
	commitments []models.CommitmentWithTxs
}

type submitCreate2TransfersBatchRequestParams struct {
	batchID     *models.Uint256
	commitments []models.CommitmentWithTxs
}

type submitMassMigrationsBatchRequestParams struct {
	batchID     *models.Uint256
	commitments []models.CommitmentWithTxs
}

type withdrawStakeRequestParams struct {
	batchID *models.Uint256
}
