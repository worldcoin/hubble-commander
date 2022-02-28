package setup

import (
	"os"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/contracts/erc20"
	"github.com/Worldcoin/hubble-commander/contracts/test/customtoken"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchstatus"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/ybbus/jsonrpc/v2"
)

type E2ETestSuite struct {
	*require.Assertions
	suite.Suite

	QueueDepositGasLimit uint64

	Cfg       *config.Config
	Commander Commander
	RPCClient jsonrpc.RPCClient
	ETHClient *eth.Client

	Domain  bls.Domain
	Wallets []bls.Wallet
}

func (s *E2ETestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())

	s.QueueDepositGasLimit = 600_000
}

func (s *E2ETestSuite) SetupTestEnvironment(commanderCfg *config.Config, deployerConfig *config.DeployerConfig) {
	var err error

	if commanderCfg == nil {
		commanderCfg = config.GetConfig()
	}

	s.Cfg = commanderCfg

	s.Commander, err = NewConfiguredCommanderFromEnv(commanderCfg, deployerConfig)
	s.NoError(err)
	err = s.Commander.Start()
	s.NoError(err)

	s.RPCClient = s.Commander.Client()

	s.ETHClient = s.newEthClient()

	s.Domain = s.getDomain()

	s.Wallets, err = CreateWallets(s.Domain)
	s.NoError(err)
}

func (s *E2ETestSuite) TearDownTest() {
	s.NoError(s.Commander.Stop())
	s.NoError(os.Remove(*s.Cfg.Bootstrap.ChainSpecPath))
}

func (s *E2ETestSuite) GetNetworkInfo() dto.NetworkInfo {
	var networkInfo dto.NetworkInfo
	err := s.RPCClient.CallFor(&networkInfo, "hubble_getNetworkInfo")
	s.NoError(err)
	return networkInfo
}

func (s *E2ETestSuite) GetTransaction(txHash common.Hash) dto.TransactionReceipt {
	var txReceipt dto.TransactionReceipt
	err := s.RPCClient.CallFor(&txReceipt, "hubble_getTransaction", []interface{}{txHash})
	s.NoError(err)
	s.Equal(txHash, txReceipt.Hash)
	return txReceipt
}

func (s *E2ETestSuite) GetAllBatches() []dto.Batch {
	var batches []dto.Batch
	err := s.RPCClient.CallFor(&batches, "hubble_getBatches", []interface{}{nil, nil})
	s.NoError(err)
	return batches
}

func (s *E2ETestSuite) GetBatchByID(batchID uint64) dto.BatchWithRootAndCommitments {
	var batch dto.BatchWithRootAndCommitments
	err := s.RPCClient.CallFor(&batch, "hubble_getBatchByID", []interface{}{models.MakeUint256(batchID)})
	s.NoError(err)
	return batch
}

func (s *E2ETestSuite) GetCommitment(commitmentID models.CommitmentID) dto.TxCommitment {
	var commitment dto.TxCommitment
	err := s.RPCClient.CallFor(&commitment, "hubble_getCommitment", []interface{}{commitmentID})
	s.NoError(err)
	return commitment
}

func (s *E2ETestSuite) GetUserStates(publicKey *models.PublicKey) []dto.UserStateWithID {
	var userStates []dto.UserStateWithID
	err := s.RPCClient.CallFor(&userStates, "hubble_getUserStates", []interface{}{publicKey})
	s.NoError(err)
	return userStates
}

func (s *E2ETestSuite) SendTransaction(tx interface{}) common.Hash {
	switch parsedTx := tx.(type) {
	case dto.Transfer:
		return s.sendTransfer(parsedTx)
	case dto.Create2Transfer:
		return s.sendCreate2Transfer(parsedTx)
	case dto.MassMigration:
		return s.sendMassMigration(parsedTx)
	default:
		panic("unexpected tx type")
	}
}

func (s *E2ETestSuite) SendNTransactions(n int, baseTx interface{}) common.Hash {
	switch t := baseTx.(type) {
	case dto.Transfer:
		return s.sendNTransfers(n, t)
	case dto.Create2Transfer:
		return s.sendNCreate2Transfers(n, t)
	case dto.MassMigration:
		return s.sendNMassMigrations(n, t)
	default:
		panic("unexpected tx type")
	}
}

func (s *E2ETestSuite) WaitForBatchStatus(batchID uint64, status batchstatus.BatchStatus) *dto.BatchWithRootAndCommitments {
	var batch dto.BatchWithRootAndCommitments

	s.Eventually(func() bool {
		var rpcError *jsonrpc.RPCError
		err := s.RPCClient.CallFor(&batch, "hubble_getBatchByID", []interface{}{models.MakeUint256(batchID)})
		if err != nil && errors.As(err, &rpcError) {
			if rpcError.Code == 30000 {
				return false
			}
		}
		s.NoError(err)
		return batch.Status == status
	}, 30*time.Second, testutils.TryInterval)

	return &batch
}

func (s *E2ETestSuite) WaitForTxToBeIncludedInBatch(txHash common.Hash) {
	s.Eventually(func() bool {
		var txReceipt dto.TransactionReceipt
		err := s.RPCClient.CallFor(&txReceipt, "hubble_getTransaction", []interface{}{txHash})
		s.NoError(err)
		return txReceipt.Status == txstatus.Mined
	}, 30*time.Second, testutils.TryInterval)
}

func (s *E2ETestSuite) ApproveToken(tokenAddress common.Address, amount string) {
	approvedAmount := models.NewUint256FromBig(*utils.ParseEther(amount))

	token, err := erc20.NewERC20(tokenAddress, s.ETHClient.Blockchain.GetBackend())
	s.NoError(err)

	tx, err := token.Approve(s.ETHClient.Blockchain.GetAccount(), s.ETHClient.ChainState.DepositManager, approvedAmount.ToBig())
	s.NoError(err)

	_, err = s.ETHClient.WaitToBeMined(tx)
	s.NoError(err)
}

func (s *E2ETestSuite) GetDeployedToken(tokenID uint64) (*models.RegisteredToken, *customtoken.TestCustomToken) {
	registeredToken, err := s.ETHClient.GetRegisteredToken(models.NewUint256(tokenID))
	s.NoError(err)

	tokenContract, err := customtoken.NewTestCustomToken(registeredToken.Contract, s.ETHClient.Blockchain.GetBackend())
	s.NoError(err)

	return registeredToken, tokenContract
}

func (s *E2ETestSuite) SubmitTxBatchAndWait(submit func() common.Hash) {
	firstTxHash := submit()
	s.WaitForTxToBeIncludedInBatch(firstTxHash)
}

func (s *E2ETestSuite) SubmitDepositBatchAndWait(targetPubKeyID, tokenID *models.Uint256, depositAmount string) {
	batchID, err := s.ETHClient.Rollup.NextBatchID(nil)
	s.NoError(err)

	fullDepositBatchCount := s.CalculateDepositsCountForFullBatch()
	parsedDepositAmount := models.NewUint256FromBig(*utils.ParseEther(depositAmount))

	txs := make([]types.Transaction, 0, fullDepositBatchCount)
	for i := 0; i < fullDepositBatchCount; i++ {
		var tx *types.Transaction
		tx, err = s.ETHClient.QueueDeposit(s.QueueDepositGasLimit, targetPubKeyID, parsedDepositAmount, tokenID)
		s.NoError(err)
		txs = append(txs, *tx)
	}
	receipts, err := s.ETHClient.WaitForMultipleTxs(txs...)
	s.NoError(err)

	for i := range receipts {
		s.EqualValues(1, receipts[i].Status)
	}

	s.WaitForBatchStatus(batchID.Uint64(), batchstatus.Mined)
}
