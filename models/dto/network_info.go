package dto

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

type NetworkInfo struct {
	ChainID              models.Uint256  `json:"chainId"`
	AccountRegistry      common.Address  `json:"accountRegistry"`
	TokenRegistry        common.Address  `json:"tokenRegistry"`
	DeploymentBlock      uint64          `json:"deploymentBlock"`
	DepositManager       common.Address  `json:"depositManager"`
	Rollup               common.Address  `json:"rollup"`
	BlockNumber          uint32          `json:"blockNumber"`
	TransactionCount     int             `json:"transactionCount"`
	AccountCount         uint32          `json:"accountCount"`
	LatestBatch          *models.Uint256 `json:"latestBatch"`
	LatestFinalisedBatch *models.Uint256 `json:"latestFinalisedBatch"`
	SignatureDomain      bls.Domain      `json:"signatureDomain"`
}
