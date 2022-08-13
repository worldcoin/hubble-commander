package disputer

import (
	"context"

	"github.com/Worldcoin/hubble-commander/bls"
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
	"github.com/Worldcoin/hubble-commander/testutils"
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

func (s *testSuiteWithContexts) SetupTest(batchType batchtype.BatchType, disableSignatures bool) {
	s.SetupTestWithConfig(batchType, &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    2,
		DisableSignatures:      disableSignatures,
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

	s.client, err = eth.NewConfiguredTestClient(&rollup.DeploymentConfig{
		Params: rollup.Params{GenesisStateRoot: root},
	}, &eth.TestClientConfig{})
	s.NoError(err)

	s.addGenesisBatch(root)
	s.newContexts(s.storage.Storage, s.client.Client, s.cfg, batchType)
}

func (s *testSuiteWithContexts) setGenesisState() {
	userStates := []models.UserState{
		*createUserState(0, 300),
		*createUserState(1, 200),
		*createUserState(2, 100),
	}

	for i := range userStates {
		_, err := s.storage.StateTree.Set(uint32(i), &userStates[i])
		s.NoError(err)
	}
}

func (s *testSuiteWithContexts) setAccounts(domain *bls.Domain) []bls.Wallet {
	wallets := testutils.GenerateWallets(s.Assertions, domain, 3)
	for i := range wallets {
		pubKeyID, err := s.client.RegisterAccountAndWait(wallets[i].PublicKey())
		s.NoError(err)
		s.EqualValues(i, *pubKeyID)

		err = s.storage.AccountTree.SetSingle(&models.AccountLeaf{
			PubKeyID:  uint32(i),
			PublicKey: *wallets[i].PublicKey(),
		})
		s.NoError(err)
	}
	return wallets
}

func (s *testSuiteWithContexts) addGenesisBatch(root *common.Hash) {
	contractBatch, err := s.client.GetContractBatch(models.NewUint256(0))
	s.NoError(err)

	batch := contractBatch.ToModelBatch()
	batch.PrevStateRoot = *root
	err = s.storage.AddBatch(batch)
	s.NoError(err)
}

func (s *testSuiteWithContexts) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *testSuiteWithContexts) newContexts(
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	batchType batchtype.BatchType,
) {
	var err error
	executionCtx := executor.NewTestExecutionContext(storage, s.client.Client, s.cfg)
	s.txsCtx, err = executor.NewTestTxsContext(executionCtx, batchType)
	s.NoError(err)

	s.syncCtx, err = syncer.NewTestTxsContext(storage, client, cfg, txtype.TransactionType(batchType))
	s.NoError(err)
	s.disputeCtx = NewContext(storage, s.client.Client)
}

func (s *testSuiteWithContexts) beginTransaction() {
	txController, txStorage := s.storage.BeginTransaction(st.TxOptions{})
	s.txController = txController
	s.newContexts(txStorage, s.client.Client, s.cfg, s.txsCtx.BatchType)
}

func (s *testSuiteWithContexts) rollback() {
	s.txController.Rollback(nil)
	s.newContexts(s.storage.Storage, s.client.Client, s.cfg, s.txsCtx.BatchType)
}

func (s *testSuiteWithContexts) submitBatch(tx models.GenericTransaction) *models.Batch {
	s.addTxs(models.MakeGenericArray(tx))

	pendingBatch, _, err := s.txsCtx.CreateAndSubmitBatch(context.Background())
	s.NoError(err)

	s.client.GetBackend().Commit()
	return pendingBatch
}

func (s *testSuiteWithContexts) addTxs(txs models.GenericTransactionArray) {
	err := s.disputeCtx.storage.BatchAddTransaction(txs)
	s.NoError(err)

	for i := 0; i < txs.Len(); i++ {
		err = s.disputeCtx.storage.AddMempoolTx(txs.At(i))
		s.NoError(err)
	}
}
