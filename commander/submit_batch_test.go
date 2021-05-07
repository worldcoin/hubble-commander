package commander

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db/postgres"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	baseCommitment = models.Commitment{
		Type:              txtype.Transfer,
		Transactions:      utils.RandomBytes(24),
		FeeReceiver:       1,
		CombinedSignature: models.MakeRandomSignature(),
		PostStateRoot:     utils.RandomHash(),
	}
)

type SubmitTransferBatchTestSuite struct {
	*require.Assertions
	suite.Suite
	db         *postgres.TestDB
	storage    *st.Storage
	tree       *st.StateTree
	cfg        *config.RollupConfig
	testClient *eth.TestClient
}

func (s *SubmitTransferBatchTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SubmitTransferBatchTestSuite) SetupTest() {
	testDB, err := postgres.NewTestDB()
	s.NoError(err)
	s.db = testDB
	s.storage = st.NewTestStorage(testDB.DB)
	s.tree = st.NewStateTree(s.storage)
	s.cfg = &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
	}

	s.testClient, err = eth.NewTestClient()
	s.NoError(err)

	err = s.storage.AddAccountIfNotExists(&models.Account{
		PubKeyID:  1,
		PublicKey: models.PublicKey{1, 2, 3},
	})
	s.NoError(err)

	userState := models.UserState{
		PubKeyID:   1,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(1000),
		Nonce:      models.MakeUint256(0),
	}

	err = s.tree.Set(1, &userState)
	s.NoError(err)
}

func (s *SubmitTransferBatchTestSuite) TearDownTest() {
	s.testClient.Close()
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_Transfers_ErrorsIfNotEnoughCommitments() {
	err := submitBatch(txtype.Transfer, []models.Commitment{}, s.storage, s.testClient.Client, s.cfg)
	s.Equal(ErrNotEnoughCommitments, err)
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_Transfers_SubmitsCommitmentsOnChain() {
	commitmentID, err := s.storage.AddCommitment(&baseCommitment)
	s.NoError(err)

	commitment, err := s.storage.GetCommitment(*commitmentID)
	s.NoError(err)

	err = submitBatch(txtype.Transfer, []models.Commitment{*commitment}, s.storage, s.testClient.Client, s.cfg)
	s.NoError(err)

	nextBatchID, err := s.testClient.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(2), nextBatchID)
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_Create2Transfers_SubmitsCommitmentsOnChain() {
	commitmentID, err := s.storage.AddCommitment(&baseCommitment)
	s.NoError(err)

	commitment, err := s.storage.GetCommitment(*commitmentID)
	s.NoError(err)

	err = submitBatch(txtype.Create2Transfer, []models.Commitment{*commitment}, s.storage, s.testClient.Client, s.cfg)
	s.NoError(err)

	nextBatchID, err := s.testClient.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(2), nextBatchID)
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_Transfers_StoresBatchRecord() {
	commitmentID, err := s.storage.AddCommitment(&baseCommitment)
	s.NoError(err)

	commitment, err := s.storage.GetCommitment(*commitmentID)
	s.NoError(err)

	err = submitBatch(txtype.Transfer, []models.Commitment{*commitment}, s.storage, s.testClient.Client, s.cfg)
	s.NoError(err)

	batch, err := s.storage.GetBatchByID(models.MakeUint256(1))
	s.NoError(err)
	s.Equal(batch.Type, txtype.Transfer)
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_Create2Transfers_StoresBatchRecord() {
	commitmentID, err := s.storage.AddCommitment(&baseCommitment)
	s.NoError(err)

	commitment, err := s.storage.GetCommitment(*commitmentID)
	s.NoError(err)

	err = submitBatch(txtype.Create2Transfer, []models.Commitment{*commitment}, s.storage, s.testClient.Client, s.cfg)
	s.NoError(err)

	batch, err := s.storage.GetBatchByID(models.MakeUint256(1))
	s.NoError(err)
	s.Equal(batch.Type, txtype.Create2Transfer)
}

// nolint:unparam
func (s *SubmitTransferBatchTestSuite) addCommitments(count int) ([]int32, []models.Commitment) {
	ids := make([]int32, 0, count)
	commitments := make([]models.Commitment, 0, count)
	for i := 0; i < count; i++ {
		id, err := s.storage.AddCommitment(&baseCommitment)
		s.NoError(err)
		ids = append(ids, *id)

		commitment, err := s.storage.GetCommitment(*id)
		s.NoError(err)
		commitments = append(commitments, *commitment)
	}
	return ids, commitments
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_Transfers_MarksCommitmentsAsIncluded() {
	ids, commitments := s.addCommitments(2)

	err := submitBatch(txtype.Transfer, commitments, s.storage, s.testClient.Client, s.cfg)
	s.NoError(err)

	batch, err := s.storage.GetBatchByID(models.MakeUint256(1))
	s.NoError(err)

	for _, id := range ids {
		commit, err := s.storage.GetCommitment(id)
		s.NoError(err)
		s.Equal(batch.Hash, *commit.IncludedInBatch)
	}
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_Create2Transfers_MarksCommitmentsAsIncluded() {
	ids, commitments := s.addCommitments(2)

	err := submitBatch(txtype.Create2Transfer, commitments, s.storage, s.testClient.Client, s.cfg)
	s.NoError(err)

	batch, err := s.storage.GetBatchByID(models.MakeUint256(1))
	s.NoError(err)

	for _, id := range ids {
		commit, err := s.storage.GetCommitment(id)
		s.NoError(err)
		s.Equal(batch.Hash, *commit.IncludedInBatch)
	}
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_Transfers_MarksCommitmentsAsIncluded_UnsavedCommitment() {
	err := submitBatch(txtype.Transfer, []models.Commitment{baseCommitment}, s.storage, s.testClient.Client, s.cfg)
	s.EqualError(err, "no rows were affected by the update")
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_Transfers_UpdatesCommitmentsAccountRoot() {
	ids, commitments := s.addCommitments(2)

	err := submitBatch(txtype.Transfer, commitments, s.storage, s.testClient.Client, s.cfg)
	s.NoError(err)

	accountRoot, err := s.testClient.AccountRegistry.Root(nil)
	s.NoError(err)

	for _, id := range ids {
		commit, err := s.storage.GetCommitment(id)
		s.NoError(err)
		s.Equal(common.BytesToHash(accountRoot[:]), *commit.AccountTreeRoot)
	}
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_Create2Transfers_UpdatesCommitmentsAccountRoot() {
	ids, commitments := s.addCommitments(2)

	err := submitBatch(txtype.Create2Transfer, commitments, s.storage, s.testClient.Client, s.cfg)
	s.NoError(err)

	accountRoot, err := s.testClient.AccountRegistry.Root(nil)
	s.NoError(err)

	for _, id := range ids {
		commit, err := s.storage.GetCommitment(id)
		s.NoError(err)
		s.Equal(common.BytesToHash(accountRoot[:]), *commit.AccountTreeRoot)
	}
}

func TestSubmitBatch_TransfersTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitTransferBatchTestSuite))
}
