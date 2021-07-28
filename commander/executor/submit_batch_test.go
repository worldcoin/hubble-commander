package executor

import (
	"context"
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
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
	storage             *st.TestStorage
	cfg                 *config.RollupConfig
	client              *eth.TestClient
	transactionExecutor *TransactionExecutor
}

func (s *SubmitTransferBatchTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *SubmitTransferBatchTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorageWithBadger()
	s.NoError(err)
	s.cfg = &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
	}

	s.client, err = eth.NewTestClient()
	s.NoError(err)

	err = s.storage.AddAccountLeafIfNotExists(&models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: models.PublicKey{1, 2, 3},
	})
	s.NoError(err)

	userState := models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	}

	_, err = s.storage.StateTree.Set(1, &userState)
	s.NoError(err)

	s.transactionExecutor = NewTestTransactionExecutor(s.storage.Storage, s.client.Client, s.cfg, context.Background())
}

func (s *SubmitTransferBatchTestSuite) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_Transfers_ErrorsIfNotEnoughCommitments() {
	pendingBatch, err := s.transactionExecutor.NewPendingBatch(txtype.Transfer)
	s.NoError(err)
	err = s.transactionExecutor.SubmitBatch(pendingBatch, []models.Commitment{})
	s.Equal(ErrNotEnoughCommitments, err)
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_Transfers_SubmitsCommitmentsOnChain() {
	commitmentID, err := s.storage.AddCommitment(&baseCommitment)
	s.NoError(err)

	commitment, err := s.storage.GetCommitment(*commitmentID)
	s.NoError(err)

	nextBatchID, err := s.client.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(1), nextBatchID)

	pendingBatch, err := s.transactionExecutor.NewPendingBatch(txtype.Transfer)
	s.NoError(err)
	err = s.transactionExecutor.SubmitBatch(pendingBatch, []models.Commitment{*commitment})
	s.NoError(err)

	s.client.Commit()

	nextBatchID, err = s.client.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(2), nextBatchID)
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_Create2Transfers_SubmitsCommitmentsOnChain() {
	commitmentID, err := s.storage.AddCommitment(&baseCommitment)
	s.NoError(err)

	commitment, err := s.storage.GetCommitment(*commitmentID)
	s.NoError(err)

	nextBatchID, err := s.client.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(1), nextBatchID)

	pendingBatch, err := s.transactionExecutor.NewPendingBatch(txtype.Create2Transfer)
	s.NoError(err)
	err = s.transactionExecutor.SubmitBatch(pendingBatch, []models.Commitment{*commitment})
	s.NoError(err)

	s.client.Commit()

	nextBatchID, err = s.client.Rollup.NextBatchID(nil)
	s.NoError(err)
	s.Equal(big.NewInt(2), nextBatchID)
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_Transfers_StoresPendingBatchRecord() {
	commitmentID, err := s.storage.AddCommitment(&baseCommitment)
	s.NoError(err)

	commitment, err := s.storage.GetCommitment(*commitmentID)
	s.NoError(err)

	pendingBatch, err := s.transactionExecutor.NewPendingBatch(txtype.Transfer)
	s.NoError(err)
	err = s.transactionExecutor.SubmitBatch(pendingBatch, []models.Commitment{*commitment})
	s.NoError(err)

	batch, err := s.storage.GetBatch(models.MakeUint256(1))
	s.NoError(err)
	s.Equal(txtype.Transfer, batch.Type)
	s.Equal(models.MakeUint256(1), batch.ID)
	s.NotEqual(common.Hash{}, batch.TransactionHash)
	s.Equal(pendingBatch.PrevStateRoot, batch.PrevStateRoot)
	s.Nil(batch.Hash)
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_Create2Transfers_StoresPendingBatchRecord() {
	commitmentID, err := s.storage.AddCommitment(&baseCommitment)
	s.NoError(err)

	commitment, err := s.storage.GetCommitment(*commitmentID)
	s.NoError(err)

	pendingBatch, err := s.transactionExecutor.NewPendingBatch(txtype.Create2Transfer)
	s.NoError(err)
	err = s.transactionExecutor.SubmitBatch(pendingBatch, []models.Commitment{*commitment})
	s.NoError(err)

	batch, err := s.storage.GetBatch(models.MakeUint256(1))
	s.NoError(err)
	s.Equal(pendingBatch.Type, batch.Type)
	s.NotEqual(common.Hash{}, batch.TransactionHash)
	s.Equal(models.MakeUint256(1), batch.ID)
	s.Equal(pendingBatch.PrevStateRoot, batch.PrevStateRoot)
	s.Nil(batch.Hash)
}

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

	pendingBatch, err := s.transactionExecutor.NewPendingBatch(txtype.Transfer)
	s.NoError(err)
	err = s.transactionExecutor.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	batch, err := s.storage.GetBatch(models.MakeUint256(1))
	s.NoError(err)

	for _, id := range ids {
		commit, err := s.storage.GetCommitment(id)
		s.NoError(err)
		s.Equal(batch.ID, *commit.IncludedInBatch)
	}
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_Create2Transfers_MarksCommitmentsAsIncluded() {
	ids, commitments := s.addCommitments(2)

	pendingBatch, err := s.transactionExecutor.NewPendingBatch(txtype.Create2Transfer)
	s.NoError(err)
	err = s.transactionExecutor.SubmitBatch(pendingBatch, commitments)
	s.NoError(err)

	batch, err := s.storage.GetBatch(models.MakeUint256(1))
	s.NoError(err)

	for _, id := range ids {
		commit, err := s.storage.GetCommitment(id)
		s.NoError(err)
		s.Equal(batch.ID, *commit.IncludedInBatch)
	}
}

func (s *SubmitTransferBatchTestSuite) TestSubmitBatch_Transfers_MarksCommitmentsAsIncluded_UnsavedCommitment() {
	pendingBatch, err := s.transactionExecutor.NewPendingBatch(txtype.Transfer)
	s.NoError(err)
	err = s.transactionExecutor.SubmitBatch(pendingBatch, []models.Commitment{baseCommitment})
	s.EqualError(err, "no rows were affected by the update")
}

func TestSubmitBatch_TransfersTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitTransferBatchTestSuite))
}
