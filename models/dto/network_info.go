package dto

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
)

type NetworkInfo struct {
	models.ChainState
	BlockNumber          uint32     `json:"blockNumber"`
	TransactionCount     int        `json:"transactionCount"`
	LatestBatch          *string    `json:"latestBatch"`
	LatestFinalisedBatch *string    `json:"latestFinalisedBatch"`
	SignatureDomain      bls.Domain `json:"signatureDomain"`
}
