package dto

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

type NetworkInfo struct {
	ChainID              models.Uint256 `json:"chainId"`
	AccountRegistry      common.Address `json:"accountRegistry"`
	DeploymentBlock      uint64         `json:"deploymentBlock"`
	Rollup               common.Address `json:"rollup"`
	BlockNumber          uint32         `json:"blockNumber"`
	TransactionCount     int            `json:"transactionCount"`
	AccountCount         uint32         `json:"accountCount"`
	LatestBatch          *string        `json:"latestBatch"`
	LatestFinalisedBatch *string        `json:"latestFinalisedBatch"`
	SignatureDomain      bls.Domain     `json:"signatureDomain"`
}
