package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BatchesTestSuite struct {
	*require.Assertions
	suite.Suite
	cmd        *Commander
	testClient *eth.TestClient
	cfg        *config.Config
	teardown   func() error
	wallets    []bls.Wallet
}

func (s *BatchesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = config.GetTestConfig()
	s.cfg.Rollup.MinCommitmentsPerBatch = 1
	s.cfg.Rollup.MaxCommitmentsPerBatch = 32
	s.cfg.Rollup.TxsPerCommitment = 1
	s.cfg.Rollup.DevMode = false
}

func (s *BatchesTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.teardown = testStorage.Teardown
	s.testClient, err = eth.NewTestClient()
	s.NoError(err)
	err = testStorage.SetChainState(&s.testClient.ChainState)
	s.NoError(err)

	s.cmd = NewCommander(s.cfg)
	s.cmd.client = s.testClient.Client
	s.cmd.storage = testStorage.Storage
	s.cmd.stopChannel = make(chan bool)

	s.wallets = generateWallets(s.T(), s.testClient.ChainState.Rollup, 2)
	seedDB(s.T(), testStorage.Storage, st.NewStateTree(testStorage.Storage), s.wallets)
}

func (s *BatchesTestSuite) TearDownTest() {
	stopCommander(s.cmd)
	s.testClient.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *BatchesTestSuite) TestUnsafeSyncBatches_DoesNotSyncExistingBatchTwice() {
	tx := testutils.MakeTransfer(0, 1, 0, 400)
	signTransfer(s.T(), &s.wallets[tx.FromStateID], &tx)
	createAndSubmitTransferBatch(s.T(), s.cmd, &tx)
	s.testClient.Commit()

	s.syncAllBlocks()

	tx2 := testutils.MakeTransfer(1, 0, 0, 100)
	signTransfer(s.T(), &s.wallets[tx2.FromStateID], &tx2)
	createAndSubmitTransferBatch(s.T(), s.cmd, &tx2)
	s.testClient.Commit()

	batches, err := s.cmd.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 1)

	s.syncAllBlocks()

	batches, err = s.cmd.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 2)

	state0, err := s.cmd.storage.GetStateLeaf(0)
	s.NoError(err)
	s.Equal(models.MakeUint256(700), state0.Balance)

	state1, err := s.cmd.storage.GetStateLeaf(1)
	s.NoError(err)
	s.Equal(models.MakeUint256(300), state1.Balance)
}

// TODO test syncRemoteBatch:
//  - test that in case of race condition local batch(es) is reverted and remote one is synced
//  - test that when syncing fraudulent batch dispute is made
//  - test that both race condition and dispute can be handled

func (s *BatchesTestSuite) syncAllBlocks() {
	latestBlockNumber, err := s.testClient.GetLatestBlockNumber()
	s.NoError(err)

	err = s.cmd.unsafeSyncBatches(0, *latestBlockNumber)
	s.NoError(err)
}

func TestBatchesTestSuite(t *testing.T) {
	suite.Run(t, new(BatchesTestSuite))
}
