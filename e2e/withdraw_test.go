//go:build e2e
// +build e2e

package e2e

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/test/customtoken"
	"github.com/Worldcoin/hubble-commander/contracts/withdrawmanager"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/consts"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/ybbus/jsonrpc/v2"
)

func TestWithdrawProcess(t *testing.T) {
	commanderConfig := config.GetConfig()
	commanderConfig.Rollup.MinTxsPerCommitment = 2
	commanderConfig.Rollup.MaxTxsPerCommitment = 2
	commanderConfig.Rollup.MinCommitmentsPerBatch = 1
	commanderConfig.API.EnableProofMethods = true

	deployerConfig := config.GetDeployerTestConfig()
	deployerConfig.Bootstrap.BlocksToFinalise = 1

	commander, err := setup.NewConfiguredCommanderFromEnv(commanderConfig, deployerConfig)
	require.NoError(t, err)
	err = commander.Start()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, commander.Stop())
	}()

	domain := GetDomain(t, commander.Client())

	wallets, err := setup.CreateWallets(domain)
	require.NoError(t, err)

	senderWallet := wallets[1]

	commanderClient := newEthClient(t, commander.Client(), commanderConfig.Ethereum.PrivateKey)
	ethClient := newEthClient(t, commander.Client(), setup.EthClientPrivateKey)
	token, tokenContract := getDeployedToken(t, ethClient)
	transferTokens(t, tokenContract, commanderClient, ethClient.Blockchain.GetAccount().From)
	approveTokens(t, tokenContract, ethClient)

	depositAmount := models.NewUint256FromBig(*utils.ParseEther("10"))
	fullDepositBatchCount := calculateDepositsCountForFullBatch(t, ethClient)

	userStatesBeforeDeposit := getSenderUserStates(t, commander.Client(), senderWallet.PublicKey())

	makeFullDepositBatch(t, commander.Client(), ethClient, depositAmount, &token.ID, tokenContract, fullDepositBatchCount)

	userStatesAfterDeposit := getSenderUserStates(t, commander.Client(), senderWallet.PublicKey())

	newUserStates := userStatesDifference(userStatesAfterDeposit, userStatesBeforeDeposit)
	require.Len(t, newUserStates, fullDepositBatchCount)

	targetMassMigrationHash := testSubmitWithdrawBatch(t, commander.Client(), senderWallet, newUserStates[0].StateID)

	withdrawManager, withdrawManagerAddress := getWithdrawManager(t, commander.Client(), ethClient)
	testProcessWithdrawCommitment(t, commander.Client(), ethClient, withdrawManager, withdrawManagerAddress, tokenContract)

	testClaimTokens(t, commander.Client(), ethClient, withdrawManager, tokenContract, senderWallet, targetMassMigrationHash)
}

func transferTokens(t *testing.T, tokenContract *customtoken.TestCustomToken, commanderClient *eth.Client, recipient common.Address) {
	tx, err := tokenContract.Transfer(commanderClient.Blockchain.GetAccount(), recipient, utils.ParseEther("1000"))
	require.NoError(t, err)
	_, err = commanderClient.WaitToBeMined(tx)
	require.NoError(t, err)
}

func getWithdrawManager(t *testing.T, client jsonrpc.RPCClient, ethClient *eth.Client) (*withdrawmanager.WithdrawManager, common.Address) {
	var info dto.NetworkInfo
	err := client.CallFor(&info, "hubble_getNetworkInfo")
	require.NoError(t, err)

	withdrawManager, err := withdrawmanager.NewWithdrawManager(info.WithdrawManager, ethClient.Blockchain.GetBackend())
	require.NoError(t, err)

	return withdrawManager, info.WithdrawManager
}

func getSenderUserStates(t *testing.T, client jsonrpc.RPCClient, senderPublicKey *models.PublicKey) []dto.UserStateWithID {
	var userStates []dto.UserStateWithID
	err := client.CallFor(&userStates, "hubble_getUserStates", []interface{}{senderPublicKey})
	require.NoError(t, err)

	return userStates
}

func calculateDepositsCountForFullBatch(t *testing.T, ethClient *eth.Client) int {
	subtreeDepth, err := ethClient.GetMaxSubtreeDepthParam()
	require.NoError(t, err)
	depositsCount := 1 << *subtreeDepth

	return depositsCount
}

func testDoActionAndAssertTokenBalanceDifference(
	t *testing.T,
	token *customtoken.TestCustomToken,
	address common.Address,
	expectedBalanceDifference models.Uint256,
	action func(),
) {
	balanceBeforeAction, err := token.BalanceOf(nil, address)
	require.NoError(t, err)

	action()

	balanceAfterAction, err := token.BalanceOf(nil, address)
	require.NoError(t, err)

	signedBalanceDifference := balanceAfterAction.Sub(balanceAfterAction, balanceBeforeAction)
	balanceDifference := models.MakeUint256FromBig(*big.NewInt(0).Abs(signedBalanceDifference))
	require.Equal(t, expectedBalanceDifference, balanceDifference)
}

func makeFullDepositBatch(
	t *testing.T,
	client jsonrpc.RPCClient,
	ethClient *eth.Client,
	depositAmount, tokenID *models.Uint256,
	token *customtoken.TestCustomToken,
	depositsNeeded int,
) {
	expectedBalanceDifference := *depositAmount.MulN(uint64(depositsNeeded))
	testDoActionAndAssertTokenBalanceDifference(t, token, ethClient.Blockchain.GetAccount().From, expectedBalanceDifference, func() {
		txs := make([]types.Transaction, 0, depositsNeeded)
		for i := 0; i < depositsNeeded; i++ {
			tx, err := ethClient.QueueDeposit(queueDepositGasLimit, models.NewUint256(1), depositAmount, tokenID)
			require.NoError(t, err)
			txs = append(txs, *tx)
		}
		receipt, err := ethClient.WaitForMultipleTxs(txs...)
		require.NoError(t, err)

		for i := range receipt {
			require.EqualValues(t, 1, receipt[i].Status)
		}

		waitForBatch(t, client, models.MakeUint256(1))
	})
}

// userStatesDifference returns the user states in `a` that aren't in `b`.
func userStatesDifference(a, b []dto.UserStateWithID) []dto.UserStateWithID {
	mb := make(map[uint32]struct{}, len(b))
	for _, x := range b {
		mb[x.StateID] = struct{}{}
	}
	var diff []dto.UserStateWithID
	for _, x := range a {
		if _, found := mb[x.StateID]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func testSendMassMigrationsForWithdrawal(
	t *testing.T,
	client jsonrpc.RPCClient,
	senderWallet bls.Wallet,
	fromStateID uint32,
	amounts []models.Uint256,
	startingNonce int,
) common.Hash {
	var firstTxHash common.Hash

	for i := range amounts {
		massMigration, err := api.SignMassMigration(&senderWallet, dto.MassMigration{
			FromStateID: ref.Uint32(fromStateID),
			SpokeID:     ref.Uint32(1),
			Amount:      &amounts[i],
			Fee:         models.NewUint256(1),
			Nonce:       models.NewUint256(uint64(startingNonce + i)),
		})
		require.NoError(t, err)

		var txHash common.Hash
		err = client.CallFor(&txHash, "hubble_sendTransaction", []interface{}{*massMigration})
		require.NoError(t, err)
		require.NotZero(t, txHash)

		testGetTransaction(t, client, txHash)

		if i == 0 {
			firstTxHash = txHash
		}
	}

	return firstTxHash
}

func testSubmitWithdrawBatch(
	t *testing.T,
	client jsonrpc.RPCClient,
	senderWallet bls.Wallet,
	fromStateID uint32,
) common.Hash {
	massMigrationWithdrawalAmount := models.MakeUint256(9 * consts.L2Unit)

	var targetMassMigrationHash common.Hash
	submitTxBatchAndWait(t, client, func() common.Hash {
		targetMassMigrationHash = testSendMassMigrationsForWithdrawal(
			t,
			client,
			senderWallet,
			fromStateID,
			[]models.Uint256{massMigrationWithdrawalAmount, models.MakeUint256(90)},
			0,
		)

		return targetMassMigrationHash
	})

	return targetMassMigrationHash
}

func testGetMassMigrationCommitmentProof(
	t *testing.T,
	client jsonrpc.RPCClient,
	commitmentID *models.CommitmentID,
) *dto.MassMigrationCommitmentProof {
	var proof dto.MassMigrationCommitmentProof
	err := client.CallFor(&proof, "hubble_getMassMigrationCommitmentProof", []interface{}{commitmentID})
	require.NoError(t, err)

	return &proof
}

func massMigrationCommitmentProofToCalldata(proof *dto.MassMigrationCommitmentProof) withdrawmanager.TypesMMCommitmentInclusionProof {
	return withdrawmanager.TypesMMCommitmentInclusionProof{
		Commitment: withdrawmanager.TypesMassMigrationCommitment{
			StateRoot: utils.ByteSliceTo32ByteArray(proof.StateRoot.Bytes()),
			Body: withdrawmanager.TypesMassMigrationBody{
				AccountRoot:  utils.ByteSliceTo32ByteArray(proof.Body.AccountRoot.Bytes()),
				Signature:    proof.Body.Signature.BigInts(),
				SpokeID:      big.NewInt(int64(proof.Body.Meta.SpokeID)),
				WithdrawRoot: utils.ByteSliceTo32ByteArray(proof.Body.WithdrawRoot.Bytes()),
				TokenID:      proof.Body.Meta.TokenID.ToBig(),
				Amount:       proof.Body.Meta.Amount.ToBig(),
				FeeReceiver:  big.NewInt(int64(proof.Body.Meta.FeeReceiverStateID)),
				Txs:          proof.Body.Transactions,
			},
		},
		Path:    big.NewInt(int64(proof.Path.Path)),
		Witness: proof.Witness.Bytes(),
	}
}

func testProcessWithdrawCommitment(
	t *testing.T,
	client jsonrpc.RPCClient,
	ethClient *eth.Client,
	withdrawManager *withdrawmanager.WithdrawManager,
	withdrawManagerAddress common.Address,
	token *customtoken.TestCustomToken,
) {
	transactor := ethClient.Blockchain.GetAccount()
	transactor.GasLimit = 1_000_000

	commitmentID := &models.CommitmentID{
		BatchID:      models.MakeUint256(2),
		IndexInBatch: 0,
	}

	proof := testGetMassMigrationCommitmentProof(t, client, commitmentID)

	typedProof := massMigrationCommitmentProofToCalldata(proof)

	expectedBalanceDifference := *proof.Body.Meta.Amount.MulN(consts.L2Unit)
	testDoActionAndAssertTokenBalanceDifference(t, token, withdrawManagerAddress, expectedBalanceDifference, func() {
		tx, err := withdrawManager.ProcessWithdrawCommitment(transactor, commitmentID.BatchID.ToBig(), typedProof)
		require.NoError(t, err)

		receipt, err := ethClient.WaitToBeMined(tx)
		require.NoError(t, err)
		require.NotZero(t, receipt.Status)
	})
}

func testGetWithdrawProof(
	t *testing.T,
	client jsonrpc.RPCClient,
	commitmentID models.CommitmentID,
	transactionHash common.Hash,
) *dto.WithdrawProof {
	var proof dto.WithdrawProof
	err := client.CallFor(&proof, "hubble_getWithdrawProof", []interface{}{commitmentID, transactionHash})
	require.NoError(t, err)

	return &proof
}

func testGetPublicKeyProof(t *testing.T, client jsonrpc.RPCClient, pubKeyID uint32) dto.PublicKeyProof {
	var proof dto.PublicKeyProof
	err := client.CallFor(&proof, "hubble_getPublicKeyProofByPubKeyID", []interface{}{pubKeyID})
	require.NoError(t, err)

	return proof
}

func withdrawProofToCalldata(proof *dto.WithdrawProof) withdrawmanager.TypesStateMerkleProofWithPath {
	return withdrawmanager.TypesStateMerkleProofWithPath{
		State: withdrawmanager.TypesUserState{
			PubkeyID: big.NewInt(int64(proof.UserState.PubKeyID)),
			TokenID:  proof.UserState.TokenID.ToBig(),
			Balance:  proof.UserState.Balance.ToBig(),
			Nonce:    proof.UserState.Nonce.ToBig(),
		},
		Path:    big.NewInt(int64(proof.Path.Path)),
		Witness: proof.Witness.Bytes(),
	}
}

func testClaimTokens(
	t *testing.T,
	client jsonrpc.RPCClient,
	ethClient *eth.Client,
	withdrawManager *withdrawmanager.WithdrawManager,
	token *customtoken.TestCustomToken,
	sender bls.Wallet,
	transactionHash common.Hash,
) {
	transactor := ethClient.Blockchain.GetAccount()
	transactor.GasLimit = 1_000_000

	commitmentID := models.CommitmentID{
		BatchID:      models.MakeUint256(2),
		IndexInBatch: 0,
	}
	proof := testGetWithdrawProof(t, client, commitmentID, transactionHash)

	typedProof := withdrawProofToCalldata(proof)

	message, err := sender.Sign(transactor.From.Bytes())
	require.NoError(t, err)

	publicKeyProof := testGetPublicKeyProof(t, client, proof.UserState.PubKeyID)

	expectedBalanceDifference := *proof.UserState.Balance.MulN(consts.L2Unit)
	testDoActionAndAssertTokenBalanceDifference(t, token, transactor.From, expectedBalanceDifference, func() {
		tx, err := withdrawManager.ClaimTokens(
			transactor,
			utils.ByteSliceTo32ByteArray(proof.Root.Bytes()),
			typedProof,
			sender.PublicKey().BigInts(),
			message.BigInts(),
			publicKeyProof.Witness.Bytes(),
		)
		require.NoError(t, err)

		receipt, err := ethClient.WaitToBeMined(tx)
		require.NoError(t, err)
		require.NotZero(t, receipt.Status)
	})
}
