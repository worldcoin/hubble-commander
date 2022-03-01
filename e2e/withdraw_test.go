//go:build e2e
// +build e2e

package e2e

import (
	"math/big"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/test/customtoken"
	"github.com/Worldcoin/hubble-commander/contracts/withdrawmanager"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchstatus"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/consts"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/suite"
)

type WithdrawalsE2ETestSuite struct {
	setup.E2ETestSuite

	transactor             *bind.TransactOpts
	withdrawManager        *withdrawmanager.WithdrawManager
	withdrawManagerAddress common.Address
	token                  *models.RegisteredToken
	tokenContract          *customtoken.TestCustomToken
	senderWallet           bls.Wallet
}

func (s *WithdrawalsE2ETestSuite) SetupTest() {
	commanderConfig := config.GetConfig()
	commanderConfig.Rollup.MinTxsPerCommitment = 32
	commanderConfig.Rollup.MaxTxsPerCommitment = 32
	commanderConfig.Rollup.MinCommitmentsPerBatch = 1
	commanderConfig.Rollup.MaxTxnDelay = 2 * time.Second

	deployerConfig := config.GetDeployerTestConfig()
	deployerConfig.Bootstrap.BlocksToFinalise = 1

	s.SetupTestEnvironment(commanderConfig, deployerConfig)

	s.transactor = s.getTransactor(1_000_000)
	s.withdrawManager, s.withdrawManagerAddress = s.getWithdrawManager()
	s.token, s.tokenContract = s.getDeployedTokenAndApprove()
	s.senderWallet = s.Wallets[1]
}

func (s *WithdrawalsE2ETestSuite) TestWithdrawals() {
	newUserStates := s.makeDeposit()

	targetMassMigrationHash := s.submitWithdrawBatch(newUserStates[0].StateID)

	s.testProcessWithdrawCommitment()

	s.testClaimTokens(targetMassMigrationHash)
}

func (s *WithdrawalsE2ETestSuite) makeDeposit() []dto.UserStateWithID {
	fullDepositBatchCount := s.CalculateDepositsCountForFullBatch()

	userStatesBeforeDeposit := s.GetUserStates(s.senderWallet.PublicKey())

	s.makeFullDepositBatch(fullDepositBatchCount)

	userStatesAfterDeposit := s.GetUserStates(s.senderWallet.PublicKey())

	newUserStates := s.userStatesDifference(userStatesAfterDeposit, userStatesBeforeDeposit)
	s.Len(newUserStates, fullDepositBatchCount)

	return newUserStates
}

func (s *WithdrawalsE2ETestSuite) testProcessWithdrawCommitment() {
	commitmentID := &models.CommitmentID{
		BatchID:      models.MakeUint256(2),
		IndexInBatch: 0,
	}

	proof := s.getMassMigrationCommitmentProof(commitmentID)

	typedProof := s.massMigrationCommitmentProofToCalldata(proof)

	expectedBalanceDifference := *proof.Body.Meta.Amount.MulN(consts.L2Unit)
	s.testDoActionAndAssertTokenBalanceDifference(s.withdrawManagerAddress, expectedBalanceDifference, func() {
		tx, err := s.withdrawManager.ProcessWithdrawCommitment(s.transactor, commitmentID.BatchID.ToBig(), typedProof)
		s.NoError(err)

		receipt, err := s.ETHClient.WaitToBeMined(tx)
		s.NoError(err)
		s.NotZero(receipt.Status)
	})
}

func (s *WithdrawalsE2ETestSuite) testClaimTokens(transactionHash common.Hash) {
	commitmentID := models.CommitmentID{
		BatchID:      models.MakeUint256(2),
		IndexInBatch: 0,
	}
	proof := s.getWithdrawProof(commitmentID, transactionHash)

	typedProof := s.withdrawProofToCalldata(proof)

	message, err := s.senderWallet.Sign(s.transactor.From.Bytes())
	s.NoError(err)

	publicKeyProof := s.getPublicKeyProof(proof.UserState.PubKeyID)

	expectedBalanceDifference := *proof.UserState.Balance.MulN(consts.L2Unit)
	s.testDoActionAndAssertTokenBalanceDifference(s.transactor.From, expectedBalanceDifference, func() {
		tx, err := s.withdrawManager.ClaimTokens(
			s.transactor,
			utils.ByteSliceTo32ByteArray(proof.Root.Bytes()),
			typedProof,
			s.senderWallet.PublicKey().BigInts(),
			message.BigInts(),
			publicKeyProof.Witness.Bytes(),
		)
		s.NoError(err)

		receipt, err := s.ETHClient.WaitToBeMined(tx)
		s.NoError(err)
		s.NotZero(receipt.Status)
	})
}

func (s *WithdrawalsE2ETestSuite) getWithdrawManager() (*withdrawmanager.WithdrawManager, common.Address) {
	networkInfo := s.GetNetworkInfo()

	blockchain, err := chain.NewRPCConnection(s.Cfg.Ethereum)
	s.NoError(err)

	backend := blockchain.GetBackend()

	withdrawManager, err := withdrawmanager.NewWithdrawManager(networkInfo.WithdrawManager, backend)
	s.NoError(err)

	return withdrawManager, networkInfo.WithdrawManager
}

func (s *WithdrawalsE2ETestSuite) getTransactor(gasLimit uint64) *bind.TransactOpts {
	account := s.ETHClient.Blockchain.GetAccount()
	account.GasLimit = gasLimit

	return account
}

func (s *WithdrawalsE2ETestSuite) getDeployedTokenAndApprove() (*models.RegisteredToken, *customtoken.TestCustomToken) {
	token, tokenContract := s.GetDeployedToken(0)
	s.ApproveToken(token.Contract, "100")

	return token, tokenContract
}

func (s *WithdrawalsE2ETestSuite) testDoActionAndAssertTokenBalanceDifference(
	address common.Address,
	expectedBalanceDifference models.Uint256,
	action func(),
) {
	balanceBeforeAction, err := s.tokenContract.BalanceOf(nil, address)
	s.NoError(err)

	action()

	balanceAfterAction, err := s.tokenContract.BalanceOf(nil, address)
	s.NoError(err)

	signedBalanceDifference := balanceAfterAction.Sub(balanceAfterAction, balanceBeforeAction)
	balanceDifference := models.MakeUint256FromBig(*big.NewInt(0).Abs(signedBalanceDifference))
	s.Equal(expectedBalanceDifference, balanceDifference)
}

func (s *WithdrawalsE2ETestSuite) makeFullDepositBatch(depositsNeeded int) {
	depositAmount := models.NewUint256FromBig(*utils.ParseEther("10"))

	expectedBalanceDifference := *depositAmount.MulN(uint64(depositsNeeded))
	targetPubKeyID := models.NewUint256(1)
	s.testDoActionAndAssertTokenBalanceDifference(s.transactor.From, expectedBalanceDifference, func() {
		txs := make([]types.Transaction, 0, depositsNeeded)
		for i := 0; i < depositsNeeded; i++ {
			tx, err := s.ETHClient.QueueDeposit(s.QueueDepositGasLimit, targetPubKeyID, depositAmount, &s.token.ID)
			s.NoError(err)
			txs = append(txs, *tx)
		}
		_, err := s.ETHClient.WaitForMultipleTxs(txs...)
		s.NoError(err)

		s.WaitForBatchStatus(1, batchstatus.Submitted)
	})
}

func (s *WithdrawalsE2ETestSuite) sendMMFromCustomWallet(wallet bls.Wallet, massMigration dto.MassMigration) common.Hash {
	signedMassMigration, err := api.SignMassMigration(&wallet, massMigration)
	s.NoError(err)

	var txHash common.Hash
	err = s.RPCClient.CallFor(&txHash, "hubble_sendTransaction", []interface{}{signedMassMigration})
	s.NoError(err)
	s.NotZero(txHash)

	return txHash
}

// userStatesDifference returns the user states in `a` that aren't in `b`.
func (s *WithdrawalsE2ETestSuite) userStatesDifference(a, b []dto.UserStateWithID) []dto.UserStateWithID {
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

func (s *WithdrawalsE2ETestSuite) submitWithdrawBatch(fromStateID uint32) common.Hash {
	massMigrationWithdrawalAmount := models.NewUint256(9 * consts.L2Unit)

	var targetMassMigrationHash common.Hash
	s.SubmitTxBatchAndWait(func() common.Hash {
		targetMassMigrationHash = s.sendMMFromCustomWallet(s.senderWallet, dto.MassMigration{
			FromStateID: ref.Uint32(fromStateID),
			SpokeID:     ref.Uint32(1),
			Amount:      massMigrationWithdrawalAmount,
			Fee:         models.NewUint256(1),
			Nonce:       models.NewUint256(0),
		})

		s.sendMMFromCustomWallet(s.senderWallet, dto.MassMigration{
			FromStateID: ref.Uint32(fromStateID),
			SpokeID:     ref.Uint32(1),
			Amount:      models.NewUint256(90),
			Fee:         models.NewUint256(1),
			Nonce:       models.NewUint256(1),
		})

		return targetMassMigrationHash
	})

	return targetMassMigrationHash
}

func (s *WithdrawalsE2ETestSuite) getMassMigrationCommitmentProof(commitmentID *models.CommitmentID) *dto.MassMigrationCommitmentProof {
	var proof dto.MassMigrationCommitmentProof
	err := s.RPCClient.CallFor(&proof, "hubble_getMassMigrationCommitmentProof", []interface{}{commitmentID})
	s.NoError(err)

	return &proof
}

func (s *WithdrawalsE2ETestSuite) massMigrationCommitmentProofToCalldata(
	proof *dto.MassMigrationCommitmentProof,
) withdrawmanager.TypesMMCommitmentInclusionProof {
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

func (s *WithdrawalsE2ETestSuite) getWithdrawProof(
	commitmentID models.CommitmentID,
	transactionHash common.Hash,
) *dto.WithdrawProof {
	var proof dto.WithdrawProof
	err := s.RPCClient.CallFor(&proof, "hubble_getWithdrawProof", []interface{}{commitmentID, transactionHash})
	s.NoError(err)

	return &proof
}

func (s *WithdrawalsE2ETestSuite) getPublicKeyProof(pubKeyID uint32) dto.PublicKeyProof {
	var proof dto.PublicKeyProof
	err := s.RPCClient.CallFor(&proof, "hubble_getPublicKeyProofByPubKeyID", []interface{}{pubKeyID})
	s.NoError(err)

	return proof
}

func (s *WithdrawalsE2ETestSuite) withdrawProofToCalldata(proof *dto.WithdrawProof) withdrawmanager.TypesStateMerkleProofWithPath {
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

func TestWithdrawalsE2ETestSuite(t *testing.T) {
	suite.Run(t, new(WithdrawalsE2ETestSuite))
}
