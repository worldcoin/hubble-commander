package commander

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	baseCommitment = models.Commitment{
		Transactions:      utils.RandomBytes(24),
		FeeReceiver:       1,
		CombinedSignature: models.MakeSignature(1, 2),
		PostStateRoot:     utils.RandomHash(),
	}
)

type SubmitBatchTestSuite struct {
	*require.Assertions
	suite.Suite
	db         *db.TestDB
	storage    *st.Storage
	tree       *st.StateTree
	cfg        *config.RollupConfig
	testClient *eth.TestClient
}

func (s *SubmitBatchTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SubmitBatchTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
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

	userState := models.UserState{
		PubKeyID:   1,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(1000),
		Nonce:      models.MakeUint256(0),
	}

	err = s.tree.Set(1, &userState)
	s.NoError(err)
}

func (s *SubmitBatchTestSuite) TearDownTest() {
	s.testClient.Close()
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *SubmitBatchTestSuite) Test_SubmitBatch_ErrorsIfNotEnoughCommitments() {
	err := submitBatch([]models.Commitment{}, s.storage, s.testClient.Client, s.cfg)
	s.Equal(ErrNotEnoughCommitments, err)
}

func (s *SubmitBatchTestSuite) Test_SubmitBatch_SubmitsCommitmentsOnChain() {
	commitmentID, err := s.storage.AddCommitment(&baseCommitment)
	s.NoError(err)

	commitment, err := s.storage.GetCommitment(*commitmentID)
	s.NoError(err)

	err = submitBatch([]models.Commitment{*commitment}, s.storage, s.testClient.Client, s.cfg)
	s.NoError(err)

	nextBatchID, err := s.testClient.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(2), nextBatchID)
}

func (s *SubmitBatchTestSuite) Test_SubmitBatch_StoresBatchRecord() {
	commitmentID, err := s.storage.AddCommitment(&baseCommitment)
	s.NoError(err)

	commitment, err := s.storage.GetCommitment(*commitmentID)
	s.NoError(err)

	err = submitBatch([]models.Commitment{*commitment}, s.storage, s.testClient.Client, s.cfg)
	s.NoError(err)

	batch, err := s.storage.GetBatchByID(models.MakeUint256(1))
	s.NoError(err)
	s.NotNil(batch)
}

func (s *SubmitBatchTestSuite) addCommitments(count int) ([]int32, []models.Commitment) {
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

func (s *SubmitBatchTestSuite) Test_SubmitBatch_MarksCommitmentsAsIncluded() {
	ids, commitments := s.addCommitments(2)

	err := submitBatch(commitments, s.storage, s.testClient.Client, s.cfg)
	s.NoError(err)

	batch, err := s.storage.GetBatchByID(models.MakeUint256(1))
	s.NoError(err)

	for _, id := range ids {
		commit, err := s.storage.GetCommitment(id)
		s.NoError(err)
		s.Equal(batch.Hash, *commit.IncludedInBatch)
	}
}

func (s *SubmitBatchTestSuite) Test_SubmitBatch_MarksCommitmentsAsIncluded_UnsavedCommitment() {
	err := submitBatch([]models.Commitment{baseCommitment}, s.storage, s.testClient.Client, s.cfg)
	s.EqualError(err, "no rows were affected by the update")
}

func (s *SubmitBatchTestSuite) Test_SubmitBatch_UpdatesCommitmentsAccountRoot() {
	ids, commitments := s.addCommitments(2)

	err := submitBatch(commitments, s.storage, s.testClient.Client, s.cfg)
	s.NoError(err)

	accountRoot, err := s.testClient.AccountRegistry.Root(nil)
	s.NoError(err)

	for _, id := range ids {
		commit, err := s.storage.GetCommitment(id)
		s.NoError(err)
		s.Equal(common.BytesToHash(accountRoot[:]), *commit.AccountTreeRoot)
	}
}

func TestSubmitBatchTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitBatchTestSuite))
}
