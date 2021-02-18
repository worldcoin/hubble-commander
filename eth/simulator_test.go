package eth

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/contracts/transfer"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestSimulator(t *testing.T) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	require.NoError(t, err)

	genesisAccounts := make(core.GenesisAlloc)
	genesisAccounts[auth.From] = core.GenesisAccount{Balance: big.NewInt(10000000000), PrivateKey: key.D.Bytes()}
	sim := backends.NewSimulatedBackend(genesisAccounts, 12_500_000)
	defer sim.Close()

	ctx := context.Background()
	fmt.Println(sim.BalanceAt(ctx, auth.From, big.NewInt(0)))

	_, _, contract, err := transfer.DeployTransfer(auth, sim)
	require.NoError(t, err)

	sim.Commit()

	result, err := contract.Encode(&bind.CallOpts{Pending: false}, transfer.OffchainTransfer{
		TxType:    big.NewInt(0),
		FromIndex: big.NewInt(0),
		ToIndex:   big.NewInt(0),
		Amount:    big.NewInt(0),
		Fee:       big.NewInt(0),
		Nonce:     big.NewInt(0),
	})
	require.NoError(t, err)

	fmt.Println("result =", result)

	// _, _, token, err := DeployMyToken(auth, sim, new(big.Int), "Simulated blockchain tokens", 0, "SBT")

	// core.GenesisAccount{Address: auth.From, Balance: big.NewInt(10000000000)}
}
