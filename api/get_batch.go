package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

func (a *API) GetBatchByHash(hash common.Hash) ([]models.Commitment, error) {
	return a.storage.GetCommitmentsByBatchHash(&hash)
}
