//go:build e2e
// +build e2e

package e2e

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/e2e/setup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type CoreCommanderE2ETestSuite struct {
	setup.E2ETestSuite
}

func (s *CoreCommanderE2ETestSuite) SetupTest() {
	cfg := config.GetConfig()
	cfg.Rollup.MinTxsPerCommitment = 32
	cfg.Rollup.MaxTxsPerCommitment = 32
	cfg.Rollup.MinCommitmentsPerBatch = 1
	cfg.Rollup.MaxTxnDelay = 2 * time.Second

	s.SetupTestEnvironment(cfg, nil)
}

func (s *CoreCommanderE2ETestSuite) TestCommander() {
	s.testGetVersion()

	// Public key of the wallet that has two different states registered to the same PubKeyID
	testWalletPublicKey := *s.Wallets[1].PublicKey()

	testUserState := s.testGetUserStates(testWalletPublicKey)
	s.testGetPublicKeyMethods(testUserState, testWalletPublicKey)

	s.testSubmitTxBatchesAndWait()

	s.testMaxBatchDelay(96)

	s.testSubmitDepositBatchAndWait()

	s.testUserStateAfterTransfers(
		*s.Wallets[1].PublicKey(),
		1,
		32*3+1,
		-1*(32*100*3+100),
	)

	feeReceiverStateID := uint64(s.Cfg.Rollup.FeeReceiverPubKeyID)
	s.testUserStateAfterTransfers(
		*s.Wallets[feeReceiverStateID].PublicKey(),
		feeReceiverStateID,
		0,
		32*10*3+10,
	)

	s.testBatches([]batchtype.BatchType{
		batchtype.Genesis,
		batchtype.Transfer,
		batchtype.Create2Transfer,
		batchtype.MassMigration,
		batchtype.Transfer,
		batchtype.Deposit,
	})
}

func (s *CoreCommanderE2ETestSuite) TestSyncing() {
	s.testSubmitTxBatchesAndWait()
	s.testSubmitDepositBatchAndWait()

	passiveCommander := s.preparePassiveCommander()

	err := passiveCommander.Start()
	s.NoError(err)
	defer func() {
		s.NoError(passiveCommander.Stop())
	}()

	activeCommanderNetworkInfo := s.GetNetworkInfo()

	// All API calls afterwards will be done to Passive Commander
	s.RPCClient = passiveCommander.Client()
	s.waitForFullSync(activeCommanderNetworkInfo.LatestBatch)

	s.testUserStateAfterTransfers(
		*s.Wallets[1].PublicKey(),
		1,
		32*3,
		-1*(32*100*3),
	)

	feeReceiverStateID := uint64(s.Cfg.Rollup.FeeReceiverPubKeyID)
	s.testUserStateAfterTransfers(
		*s.Wallets[feeReceiverStateID].PublicKey(),
		feeReceiverStateID,
		0,
		32*10*3,
	)

	s.testBatches([]batchtype.BatchType{
		batchtype.Genesis,
		batchtype.Transfer,
		batchtype.Create2Transfer,
		batchtype.MassMigration,
		batchtype.Deposit,
	})
}

func (s *CoreCommanderE2ETestSuite) testGetVersion() {
	var version string
	err := s.RPCClient.CallFor(&version, "hubble_getVersion")
	s.NoError(err)
	s.Equal(config.GetConfig().API.Version, version)
}

func (s *CoreCommanderE2ETestSuite) testGetUserStates(targetPublicKey models.PublicKey) *dto.UserStateWithID {
	userStates := s.GetUserStates(targetPublicKey)
	s.Len(userStates, 2)

	s.EqualValues(1, userStates[0].StateID)
	s.Equal(models.MakeUint256(0), userStates[0].Nonce)
	s.EqualValues(3, userStates[1].StateID)
	s.Equal(models.MakeUint256(0), userStates[1].Nonce)

	return &userStates[0]
}

func (s *CoreCommanderE2ETestSuite) testGetPublicKeyMethods(state *dto.UserStateWithID, targetPublicKey models.PublicKey) {
	var publicKey models.PublicKey
	err := s.RPCClient.CallFor(&publicKey, "hubble_getPublicKeyByPubKeyID", []interface{}{state.PubKeyID})
	s.NoError(err)
	s.Equal(targetPublicKey, publicKey)

	err = s.RPCClient.CallFor(&publicKey, "hubble_getPublicKeyByStateID", []interface{}{state.StateID})
	s.NoError(err)
	s.Equal(targetPublicKey, publicKey)
}

func (s *CoreCommanderE2ETestSuite) testSubmitTxBatchesAndWait() {
	s.SubmitTxBatchAndWait(func() common.Hash {
		return s.SendNTransactions(32, dto.Transfer{
			FromStateID: ref.Uint32(1),
			ToStateID:   ref.Uint32(2),
			Amount:      models.NewUint256(90),
			Fee:         models.NewUint256(10),
			Nonce:       models.NewUint256(0),
		})
	})

	s.SubmitTxBatchAndWait(func() common.Hash {
		firstC2TWallet := s.Wallets[len(s.Wallets)-32]
		return s.SendNTransactions(32, dto.Create2Transfer{
			FromStateID: ref.Uint32(1),
			ToPublicKey: firstC2TWallet.PublicKey(),
			Amount:      models.NewUint256(90),
			Fee:         models.NewUint256(10),
			Nonce:       models.NewUint256(32),
		})
	})

	s.SubmitTxBatchAndWait(func() common.Hash {
		return s.SendNTransactions(32, dto.MassMigration{
			FromStateID: ref.Uint32(1),
			SpokeID:     ref.Uint32(1),
			Amount:      models.NewUint256(90),
			Fee:         models.NewUint256(10),
			Nonce:       models.NewUint256(64),
		})
	})
}

func (s *CoreCommanderE2ETestSuite) testSubmitDepositBatchAndWait() {
	registeredToken, _ := s.GetDeployedToken(0)
	s.ApproveToken(registeredToken.Contract, "100")
	depositTargetPubKeyID := models.NewUint256(2)
	s.SubmitDepositBatchAndWait(depositTargetPubKeyID, &registeredToken.ID, "10")
}

// Confirms that batches smaller than the minimum will be submitted if any tx is left pending for too long
func (s *CoreCommanderE2ETestSuite) testMaxBatchDelay(startNonce uint64) {
	txHash := s.SendTransaction(dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(90),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(startNonce),
	})

	time.Sleep(1 * time.Second)

	txReceipt := s.GetTransaction(txHash)
	s.NotEqual(txReceipt.Status, txstatus.Mined)

	log.Warn("Delayed tx is not yet in batch, waiting..")

	s.Eventually(func() bool {
		txReceipt := s.GetTransaction(txHash)
		return txReceipt.Status == txstatus.Mined
	}, 10*time.Second, testutils.TryInterval)
}

func (s *CoreCommanderE2ETestSuite) testUserStateAfterTransfers(targetPublicKey models.PublicKey, expectedStateID, expectedNonce uint64, expectedBalanceDifference int64) {
	userStates := s.GetUserStates(targetPublicKey)
	senderState := userStates[0]

	s.EqualValues(expectedStateID, senderState.StateID)
	s.Equal(models.MakeUint256(expectedNonce), senderState.Nonce)
	initialBalance := models.MakeUint256(setup.InitialGenesisBalance)

	if expectedBalanceDifference > 0 {
		s.Equal(*initialBalance.AddN(uint64(expectedBalanceDifference)), senderState.Balance)
	} else {
		s.Equal(*initialBalance.SubN(uint64(-1 * expectedBalanceDifference)), senderState.Balance)
	}
}

func (s *CoreCommanderE2ETestSuite) testBatches(expectedBatches []batchtype.BatchType) {
	batches := s.GetAllBatches()
	s.Len(batches, len(expectedBatches))

	for i := range batches {
		s.Equal(expectedBatches[i], batches[i].Type)
	}
}

func (s *CoreCommanderE2ETestSuite) testCommanderRestart(startNonce uint64) {
	err := s.Commander.Restart()
	s.NoError(err)

	txHash := s.SendTransaction(dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(90),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(startNonce),
	})

	txReceipt := s.GetTransaction(txHash)
	s.Equal(txReceipt.Status, txstatus.Pending)
}

func (s *CoreCommanderE2ETestSuite) preparePassiveCommander() *setup.InProcessCommander {
	newCommanderCfg := *s.Cfg
	newCommanderCfg.Bootstrap.Prune = true
	newCommanderCfg.API.Port = "5555"
	newCommanderCfg.Metrics.Port = "2222"
	newCommanderCfg.Badger.Path += "_passive"
	newCommanderCfg.Ethereum.PrivateKey = "ab6919fd6ac00246bb78657e0696cf72058a4cb395133d074eabaddb83d8b00c"

	passiveCommander, err := setup.CreateInProcessCommander(&newCommanderCfg, nil)
	s.NoError(err)

	return passiveCommander
}

func (s *CoreCommanderE2ETestSuite) waitForFullSync(latestBatch *models.Uint256) {
	s.Eventually(func() bool {
		networkInfo := s.GetNetworkInfo()
		return networkInfo.LatestBatch != nil && networkInfo.LatestBatch.Cmp(latestBatch) >= 0
	}, 30*time.Second, testutils.TryInterval)
}

func TestCoreCommanderE2ETestSuite(t *testing.T) {
	suite.Run(t, new(CoreCommanderE2ETestSuite))
}
