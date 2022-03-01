package syncer

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/ethereum/go-ethereum/common"
)

var syncTestSuiteConfig = config.RollupConfig{
	MinCommitmentsPerBatch: 1,
	MaxCommitmentsPerBatch: 32,
	MinTxsPerCommitment:    1,
	MaxTxsPerCommitment:    1,
	DisableSignatures:      false,
}

// Other test suites encapsulate syncTestSuite. Don't add any tests on syncTestSuite to avoid repeated runs.
type syncTestSuite struct {
	testSuiteWithSyncAndRollupContext

	domain  *bls.Domain
	wallets []bls.Wallet
}

func (s *syncTestSuite) setupTest() {
	s.NotNil(s.client) // make sure testSuiteWithSyncAndRollupContext.SetupTest was called before

	var err error
	s.domain, err = s.client.GetDomain()
	s.NoError(err)

	s.wallets = testutils.GenerateWallets(s.Assertions, s.domain, 2)

	s.seedDB(s.wallets)
}

func (s *syncTestSuite) seedDB(wallets []bls.Wallet) {
	err := s.storage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  0,
		PublicKey: *wallets[0].PublicKey(),
	})
	s.NoError(err)

	err = s.storage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: *wallets[1].PublicKey(),
	})
	s.NoError(err)

	s.setBatchAccounts()

	_, err = s.storage.StateTree.Set(0, &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	_, err = s.storage.StateTree.Set(1, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)
}

func (s *syncTestSuite) setBatchAccounts() {
	accounts := make([]models.AccountLeaf, st.AccountBatchSize)
	accounts[0] = models.AccountLeaf{
		PubKeyID:  st.AccountBatchOffset,
		PublicKey: *s.wallets[0].PublicKey(),
	}

	for i := 1; i < len(accounts); i++ {
		accounts[i] = models.AccountLeaf{
			PubKeyID:  uint32(st.AccountBatchOffset + i),
			PublicKey: *s.wallets[0].PublicKey(),
		}
	}
	err := s.storage.AccountTree.SetBatch(accounts)
	s.NoError(err)
}

func (s *syncTestSuite) createCommitmentWithEmptyTransactions(commitmentType batchtype.BatchType) models.TxCommitmentWithTxs {
	stateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	return models.TxCommitmentWithTxs{
		TxCommitment: models.TxCommitment{
			CommitmentBase: models.CommitmentBase{
				Type:          commitmentType,
				PostStateRoot: *stateRoot,
			},
			FeeReceiver:       0,
			CombinedSignature: models.Signature{},
		},
		Transactions: []byte{},
	}
}

func (s *syncTestSuite) syncAllBatches() {
	newRemoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)

	for i := range newRemoteBatches {
		err = s.syncCtx.SyncBatch(newRemoteBatches[i])
		s.NoError(err)
	}
}

func (s *syncTestSuite) recreateDatabase() {
	err := s.storage.Teardown()
	s.NoError(err)

	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	executionCtx := executor.NewTestExecutionContext(s.storage.Storage, s.client.Client, s.cfg)
	s.txsCtx, err = executor.NewTestTxsContext(executionCtx, s.txsCtx.BatchType)
	s.NoError(err)
	s.syncCtx, err = NewTestContext(s.storage.Storage, s.client.Client, s.cfg, s.txsCtx.BatchType)
	s.NoError(err)

	s.seedDB(s.wallets)
}

func (s *syncTestSuite) getAccountTreeRoot() common.Hash {
	rawAccountRoot, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)
	return common.BytesToHash(rawAccountRoot[:])
}

func (s *syncTestSuite) submitBatch(tx models.GenericTransaction) []models.CommitmentWithTxs {
	pendingBatch, commitments := s.createBatch(tx)

	err := s.txsCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.GetBackend().Commit()
	return commitments
}

func (s *syncTestSuite) createBatch(tx models.GenericTransaction) (*models.Batch, []models.CommitmentWithTxs) {
	s.addTx(tx)

	pendingBatch, err := s.txsCtx.NewPendingBatch(s.txsCtx.BatchType)
	s.NoError(err)

	commitments, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(commitments, 1)

	return pendingBatch, commitments
}

func (s *syncTestSuite) addTx(tx models.GenericTransaction) {
	err := s.storage.AddTransaction(tx)
	s.NoError(err)
	_, err = s.txsCtx.Mempool.AddOrReplace(s.storage.Storage, tx)
	s.NoError(err)
}
