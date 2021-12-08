package disputer

import (
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/commander/syncer"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/eth/deployer/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testSuiteWithContexts struct {
	*require.Assertions
	suite.Suite
	storage      *st.TestStorage
	txController *db.TxController
	cfg          *config.RollupConfig
	client       *eth.TestClient
	txsCtx       *executor.TxsContext
	syncCtx      *syncer.TxsContext
	disputeCtx   *Context
}

func (s *testSuiteWithContexts) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *testSuiteWithContexts) SetupTest(batchType batchtype.BatchType) {
	s.SetupTestWithConfig(batchType, &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    1,
		DisableSignatures:      false,
	})
}

func (s *testSuiteWithContexts) SetupTestWithConfig(batchType batchtype.BatchType, cfg *config.RollupConfig) {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.cfg = cfg

	s.setGenesisState()
	root, err := s.storage.StateTree.Root()
	s.NoError(err)

	s.client, err = eth.NewConfiguredTestClient(rollup.DeploymentConfig{
		Params: rollup.Params{
			GenesisStateRoot: root,
		},
	}, eth.ClientConfig{})
	s.NoError(err)

	s.addGenesisBatch(root)
	s.newContexts(s.storage.Storage, s.client.Client, s.cfg, batchType)
}

func (s *testSuiteWithContexts) setGenesisState() {
	userStates := []models.UserState{
		*createUserState(0, 300, 0),
		*createUserState(1, 200, 0),
		*createUserState(2, 100, 0),
	}

	for i := range userStates {
		_, err := s.storage.StateTree.Set(uint32(i), &userStates[i])
		s.NoError(err)
	}
}

func (s *testSuiteWithContexts) addGenesisBatch(root *common.Hash) {
	batch, err := s.client.GetBatch(models.NewUint256(0))
	s.NoError(err)

	batch.PrevStateRoot = root
	err = s.storage.AddBatch(batch)
	s.NoError(err)
}

func (s *testSuiteWithContexts) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *testSuiteWithContexts) newContexts(
	storage *st.Storage, client *eth.Client, cfg *config.RollupConfig, batchType batchtype.BatchType,
) {
	executionCtx := executor.NewTestExecutionContext(storage, s.client.Client, s.cfg)
	s.txsCtx = executor.NewTestTxsContext(executionCtx, batchType)
	s.syncCtx = syncer.NewTestTxsContext(storage, client, cfg, batchType)
	s.disputeCtx = NewContext(storage, s.client.Client)
}

func (s *testSuiteWithContexts) beginTransaction() {
	txController, txStorage := s.storage.BeginTransaction(st.TxOptions{})
	s.txController = txController
	s.newContexts(txStorage, s.client.Client, s.cfg, s.txsCtx.BatchType)
}

func (s *testSuiteWithContexts) commitTransaction() {
	err := s.txController.Commit()
	s.NoError(err)
	s.newContexts(s.storage.Storage, s.client.Client, s.cfg, s.txsCtx.BatchType)
}

func (s *testSuiteWithContexts) rollback() {
	s.txController.Rollback(nil)
	s.newContexts(s.storage.Storage, s.client.Client, s.cfg, s.txsCtx.BatchType)
}

func (s *testSuiteWithContexts) submitBatch(tx models.GenericTransaction) *models.Batch {
	pendingBatch, batchData := s.createBatch(tx)

	err := s.txsCtx.SubmitBatch(pendingBatch, batchData)
	s.NoError(err)

	s.client.GetBackend().Commit()
	return pendingBatch
}

func (s *testSuiteWithContexts) createBatch(tx models.GenericTransaction) (*models.Batch, executor.BatchData) {
	var err error
	switch tx.Type() {
	case txtype.Transfer:
		err = s.disputeCtx.storage.AddTransfer(tx.ToTransfer())
	case txtype.Create2Transfer:
		err = s.disputeCtx.storage.AddCreate2Transfer(tx.ToCreate2Transfer())
	case txtype.MassMigration:
		err = s.disputeCtx.storage.AddMassMigration(tx.ToMassMigration())
	}
	s.NoError(err)

	pendingBatch, err := s.txsCtx.NewPendingBatch(s.txsCtx.BatchType)
	s.NoError(err)

	batchData, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(batchData.Commitments(), 1)

	return pendingBatch, batchData
}
