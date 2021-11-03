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

func (s *BenchmarkSuite) sendC2T(wallet bls.Wallet, from uint32, to *models.PublicKey, nonce models.Uint256) common.Hash {
	transfer, err := api.SignCreate2Transfer(&wallet, dto.Create2Transfer{
		FromStateID: &from,
		ToPublicKey: to,
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
