package dto

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

type NetworkInfo struct {
	ChainID                        models.Uint256
	AccountRegistry                common.Address
	AccountRegistryDeploymentBlock uint64
	TokenRegistry                  common.Address
	SpokeRegistry                  common.Address
	DepositManager                 common.Address
	WithdrawManager                common.Address
	Rollup                         common.Address
	BlockNumber                    uint32
	TransactionCount               int
	AccountCount                   uint32
	LatestBatch                    *models.Uint256
	LatestFinalisedBatch           *models.Uint256
	SignatureDomain                bls.Domain
}
