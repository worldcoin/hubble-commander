//go:build e2e
// +build e2e

package bench

import (
	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/ethereum/go-ethereum/common"
)

func (s *BenchmarkSuite) sendTransfer(wallet bls.Wallet, from, to uint32, nonce models.Uint256) common.Hash {
	transfer, err := api.SignTransfer(&wallet, dto.Transfer{
		FromStateID: &from,
		ToStateID:   &to,
		Amount:      models.NewUint256(1),
		Fee:         models.NewUint256(1),
		Nonce:       &nonce,
	})
	s.NoError(err)

	var transferHash common.Hash
	err = s.commander.Client().CallFor(&transferHash, "hubble_sendTransaction", []interface{}{*transfer})
	s.NoError(err)
	s.NotNil(transferHash)

	return transferHash
}
