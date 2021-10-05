package syncer

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
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

	seedDB(s.Assertions, s.storage.Storage, s.wallets)
}

func seedDB(s *require.Assertions, storage *st.Storage, wallets []bls.Wallet) {
	err := storage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  0,
		PublicKey: *wallets[0].PublicKey(),
	})
	s.NoError(err)

	err = storage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: *wallets[1].PublicKey(),
	})
	s.NoError(err)

	_, err = storage.StateTree.Set(0, &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	_, err = storage.StateTree.Set(1, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)
}

func (s *syncTestSuite) createCommitmentWithEmptyTransactions(commitmentType batchtype.BatchType) models.Commitment {
	stateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	return models.Commitment{
		Type:              commitmentType,
		Transactions:      []byte{},
		FeeReceiver:       0,
		CombinedSignature: models.Signature{},
		PostStateRoot:     *stateRoot,
	}
}

func (s *syncTestSuite) syncAllBatches() {
	newRemoteBatches, err := s.client.GetAllBatches()
	s.NoError(err)

	for i := range newRemoteBatches {
		remoteBatch := &newRemoteBatches[i]
		err = s.syncCtx.SyncBatch(remoteBatch)
		s.NoError(err)
	}
}

func (s *syncTestSuite) recreateDatabase() {
	err := s.storage.Teardown()
	s.NoError(err)

	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	executionCtx := executor.NewTestExecutionContext(s.storage.Storage, s.client.Client, s.cfg)
	s.rollupCtx = executor.NewTestRollupContext(executionCtx, s.rollupCtx.BatchType)
	s.syncCtx = NewTestContext(s.storage.Storage, s.client.Client, s.cfg, s.syncCtx.BatchType)

	seedDB(s.Assertions, s.storage.Storage, s.wallets)
}

func (s *syncTestSuite) getAccountTreeRoot() common.Hash {
	rawAccountRoot, err := s.client.AccountRegistry.Root(nil)
	s.NoError(err)
	return common.BytesToHash(rawAccountRoot[:])
}

func (s *syncTestSuite) submitBatch(tx models.GenericTransaction) *models.Batch {
	pendingBatch, commitments := s.createBatch(tx)

	err := s.rollupCtx.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	s.client.Commit()
	return pendingBatch
}

func (s *syncTestSuite) createBatch(tx models.GenericTransaction) (*models.Batch, []models.Commitment) {
	if tx.Type() == txtype.Transfer {
		err := s.storage.AddTransfer(tx.ToTransfer())
		s.NoError(err)
	} else {
		err := s.storage.AddCreate2Transfer(tx.ToCreate2Transfer())
		s.NoError(err)
	}

	pendingBatch, err := s.rollupCtx.NewPendingBatch(s.rollupCtx.BatchType)
	s.NoError(err)

	commitments, err := s.rollupCtx.CreateCommitments()
	s.NoError(err)
	s.Len(commitments, 1)

	return pendingBatch, commitments
}
