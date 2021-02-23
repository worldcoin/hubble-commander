package simulator

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/frontendtransfer"
	"github.com/stretchr/testify/require"
)

func TestSimulator(t *testing.T) {
	sim, err := NewSimulator()
	require.NoError(t, err)
	defer sim.Close()

	_, _, contract, err := frontendtransfer.DeployFrontendTransfer(sim.Account, sim.Backend)
	require.NoError(t, err)

	sim.Backend.Commit()

	_, err = contract.Encode(nil, frontendtransfer.OffchainTransfer{
		TxType:    big.NewInt(0),
		FromIndex: big.NewInt(0),
		ToIndex:   big.NewInt(0),
		Amount:    big.NewInt(0),
		Fee:       big.NewInt(0),
		Nonce:     big.NewInt(0),
	})
	require.NoError(t, err)
}
