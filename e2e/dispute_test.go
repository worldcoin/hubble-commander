package e2e

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/ybbus/jsonrpc/v2"
)

func TestCommanderDispute(t *testing.T) {
	cmd := setup.CreateInProcessCommander()
	err := cmd.Start()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, cmd.Stop())
	}()

	domain := getDomain(t, cmd.Client())

	wallets, err := setup.CreateWallets(domain)
	require.NoError(t, err)

	senderWallet := wallets[1]

	firstTransferHash := testSendTransfer(t, cmd.Client(), senderWallet, models.NewUint256(0))
	testGetTransaction(t, cmd.Client(), firstTransferHash)
	send31MoreTransfers(t, cmd.Client(), senderWallet)

	waitForTxToBeIncludedInBatch(t, cmd.Client(), firstTransferHash)

	sendInvalidTransfer(t, cmd.Client(), senderWallet, models.NewUint256(32))
}

func sendInvalidTransfer(t *testing.T, client jsonrpc.RPCClient, senderWallet bls.Wallet, nonce *models.Uint256) {
	transferDTO, err := api.SignTransfer(&senderWallet, dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(2_000_000_000_000_000_000),
		Fee:         models.NewUint256(10),
		Nonce:       nonce,
	})
	require.NoError(t, err)

	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: *transferDTO.FromStateID,
			Amount:      *transferDTO.Amount,
			Fee:         *transferDTO.Fee,
		},
		ToStateID: *transferDTO.ToStateID,
	}

	encodedTransfer, err := encoder.EncodeTransferForCommitment(&transfer)
	require.NoError(t, err)

	ethClient := newEthClient(t, client)

	commitment := models.Commitment{
		Transactions:      encodedTransfer,
		FeeReceiver:       0,
		CombinedSignature: models.Signature{},
		PostStateRoot:     utils.RandomHash(),
	}
	transaction, err := ethClient.SubmitTransfersBatch([]models.Commitment{commitment})
	require.NoError(t, err)

	_, err = deployer.WaitToBeMined(ethClient.ChainConnection.GetBackend(), transaction)
	require.NoError(t, err)

	_, err = ethClient.GetBatch(models.NewUint256(2))
	require.NoError(t, err)
}

func newEthClient(t *testing.T, client jsonrpc.RPCClient) *eth.Client {
	var info dto.NetworkInfo
	err := client.CallFor(&info, "hubble_getNetworkInfo")
	require.NoError(t, err)

	chainState := models.ChainState{
		ChainID:         info.ChainID,
		AccountRegistry: info.AccountRegistry,
		DeploymentBlock: info.DeploymentBlock,
		Rollup:          info.Rollup,
	}

	cfg := config.GetConfig()
	chain, err := deployer.NewRPCChainConnection(cfg.Ethereum)
	require.NoError(t, err)

	accountRegistry, err := accountregistry.NewAccountRegistry(chainState.AccountRegistry, chain.GetBackend())
	require.NoError(t, err)

	rollupContract, err := rollup.NewRollup(chainState.Rollup, chain.GetBackend())
	require.NoError(t, err)

	ethClient, err := eth.NewClient(chain, &eth.NewClientParams{
		ChainState:      chainState,
		Rollup:          rollupContract,
		AccountRegistry: accountRegistry,
	})
	require.NoError(t, err)
	return ethClient
}
