package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/commander/tracker"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const lowGasLimit = 40_000

type TxsTrackingTestSuite struct {
	*require.Assertions
	suite.Suite
	cmd     *Commander
	storage *st.TestStorage
	client  *eth.TestClient
	cfg     *config.Config
	wallets []bls.Wallet
}

func (s *TxsTrackingTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = config.GetTestConfig()
	s.cfg.Rollup.MinCommitmentsPerBatch = 1
	s.cfg.Rollup.MaxCommitmentsPerBatch = 32
	s.cfg.Rollup.MinTxsPerCommitment = 1
	s.cfg.Rollup.MaxTxsPerCommitment = 1
	s.cfg.Rollup.DisableSignatures = true
}

func (s *TxsTrackingTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
}

func (s *TxsTrackingTestSuite) setupTestWithClientConfig(cfg *eth.ClientConfig) {
	s.cmd = NewCommander(s.cfg, nil)
	// pass txs channels to testClient to use commander tracking worker
	clientCfg := &eth.TestClientConfig{
		ClientConfig: *cfg,
		TxsChannels:  s.cmd.txsTrackingChannels,
	}

	s.client = newClientWithGenesisStateWithClientConfig(s.T(), s.storage, clientCfg)

	setStateLeaves(s.T(), s.storage.Storage)
	s.cmd.client = s.client.Client
	s.cmd.blockchain = s.client.Blockchain
	s.cmd.storage = s.storage.Storage
	s.cmd.txsTracker = tracker.NewTracker(s.client.Client, clientCfg.TxsChannels.SentTxs, clientCfg.TxsChannels.Requests)

	err := s.cmd.addGenesisBatch()
	s.NoError(err)

	domain, err := s.client.GetDomain()
	s.NoError(err)
	s.wallets = testutils.GenerateWallets(s.Assertions, domain, 2)
	s.setAccountsAndChainState()

	s.startWorkers()
}

func (s *TxsTrackingTestSuite) TearDownTest() {
	s.cmd.workersWaitGroup.Wait()
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *TxsTrackingTestSuite) TestTrackSentTxs_TransferTransaction() {
	s.setupTestWithClientConfig(&eth.ClientConfig{TransferBatchSubmissionGasLimit: ref.Uint64(lowGasLimit)})
	transfer := testutils.MakeTransfer(0, 1, 0, 400)
	s.submitBatch(&transfer, batchtype.Transfer)
}

func (s *TxsTrackingTestSuite) TestTrackSentTxs_Create2TransfersTransaction() {
	s.setupTestWithClientConfig(&eth.ClientConfig{C2TBatchSubmissionGasLimit: ref.Uint64(lowGasLimit)})
	transfer := testutils.MakeCreate2Transfer(
		0,
		ref.Uint32(1),
		0,
		50,
		&models.PublicKey{2, 3, 4},
	)
	s.submitBatch(&transfer, batchtype.Create2Transfer)
}

func (s *TxsTrackingTestSuite) TestTrackSentTxs_MassMigrationTransaction() {
	s.setupTestWithClientConfig(&eth.ClientConfig{MMBatchSubmissionGasLimit: ref.Uint64(lowGasLimit)})
	massMigration := testutils.MakeMassMigration(0, 2, 0, 50)
	s.submitBatch(&massMigration, batchtype.MassMigration)
}

func (s *TxsTrackingTestSuite) TestTrackSentTxs_BatchAccountRegistrationTransaction() {
	s.setupTestWithClientConfig(&eth.ClientConfig{BatchAccountRegistrationGasLimit: ref.Uint64(lowGasLimit)})
	publicKeys := make([]models.PublicKey, st.AccountBatchSize)
	_, err := s.client.Client.RegisterBatchAccount(publicKeys)
	s.NoError(err)
}

func (s *TxsTrackingTestSuite) TestTrackSentTxs_WithdrawStake() {
	s.setupTestWithClientConfig(&eth.ClientConfig{StakeWithdrawalGasLimit: ref.Uint64(lowGasLimit)})
	transfer := testutils.MakeTransfer(0, 1, 0, 400)
	batch := s.submitBatch(&transfer, batchtype.Transfer)

	_, err := s.client.Client.WithdrawStake(&batch.ID)
	s.NoError(err)
}

func (s *TxsTrackingTestSuite) TestTrackSentTxs_SubmitDepositBatch() {
	s.setupTestWithClientConfig(&eth.ClientConfig{DepositBatchSubmissionGasLimit: ref.Uint64(lowGasLimit)})
	err := s.storage.AddPendingDepositSubtree(&models.PendingDepositSubtree{
		ID:       models.MakeUint256(1),
		Root:     utils.RandomHash(),
		Deposits: testutils.GetFourDeposits(),
	})
	s.NoError(err)
	executionCtx := executor.NewTestExecutionContext(s.storage.Storage, s.client.Client, nil)
	depositsCtx := executor.NewTestDepositsContext(executionCtx)

	_, _, err = depositsCtx.CreateAndSubmitBatch()
	s.NoError(err)
}

func (s *TxsTrackingTestSuite) TestTrackSentTxs_ClosesTxsChannelOnEthTxError() {
	s.setupTestWithClientConfig(&eth.ClientConfig{
		BatchAccountRegistrationGasLimit: ref.Uint64(lowGasLimit),
	})
	transfer := testutils.MakeCreate2Transfer(
		0,
		ref.Uint32(1),
		0,
		50,
		&models.PublicKey{2, 3, 4},
	)
	s.submitBatch(&transfer, batchtype.Create2Transfer)
}

func (s *TxsTrackingTestSuite) submitBatch(tx models.GenericTransaction, batchType batchtype.BatchType) *models.Batch {
	executionCtx := executor.NewTestExecutionContext(s.storage.Storage, s.client.Client, s.cfg.Rollup)
	txsCtx := executor.NewTestTxsContext(executionCtx, batchType)

	err := s.storage.AddTransaction(tx)
	s.NoError(err)

	batch, _, err := txsCtx.CreateAndSubmitBatch()
	s.NoError(err)
	s.client.Backend.Commit()
	return batch
}

func (s *TxsTrackingTestSuite) startWorkers() {
	s.cmd.startWorker("Test Sending Requested Txs", func() error {
		err := s.cmd.txsTracker.SendRequestedTxs(s.cmd.workersContext)
		s.NoError(err)
		return err
	})
	s.cmd.startWorker("Test Tracking Sent Txs", func() error {
		err := s.cmd.txsTracker.TrackSentTxs(s.cmd.workersContext)
		s.Error(err)
		return err
	})
}

func (s *TxsTrackingTestSuite) setAccountsAndChainState() {
	setChainState(s.T(), s.storage)
	setAccountLeaves(s.T(), s.storage.Storage, s.wallets)
}

func newClientWithGenesisStateWithClientConfig(t *testing.T, storage *st.TestStorage, conf *eth.TestClientConfig) *eth.TestClient {
	setStateLeaves(t, storage.Storage)
	genesisRoot, err := storage.StateTree.Root()
	require.NoError(t, err)

	client, err := eth.NewConfiguredTestClient(&rollup.DeploymentConfig{
		Params: rollup.Params{
			GenesisStateRoot: genesisRoot,
		},
	}, conf)

	require.NoError(t, err)

	return client
}

func TestTxsTrackingTestSuite(t *testing.T) {
	suite.Run(t, new(TxsTrackingTestSuite))
}
