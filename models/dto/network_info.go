package dto

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
)

type NetworkInfo struct {
	models.ChainState
	BlockNumber          uint32          `json:"blockNumber"`
	TransactionCount     int             `json:"transactionCount"`
	AccountCount         uint32          `json:"accountCount"`
	LatestBatch          *models.Uint256 `json:"latestBatch"`
	LatestFinalisedBatch *models.Uint256 `json:"latestFinalisedBatch"`
	SignatureDomain      bls.Domain      `json:"signatureDomain"`
}
