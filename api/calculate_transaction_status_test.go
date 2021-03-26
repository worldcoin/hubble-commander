package api

import (
	"math"
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	commitment = models.Commitment{
		Transactions:      utils.RandomBytes(24),
		FeeReceiver:       1,
		CombinedSignature: models.MakeSignature(1, 2),
		PostStateRoot:     utils.RandomHash(),
	}
)

type CalculateTransactionStatusTestSuite struct {
	*require.Assertions
	suite.Suite
	db      *db.TestDB
	storage *st.Storage
	sim     *simulator.Simulator
	tx      *models.Transaction
}

func (s *CalculateTransactionStatusTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *CalculateTransactionStatusTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)

	s.storage = st.NewTestStorage(testDB.DB)
	s.db = testDB

	sim, err := simulator.NewSimulator()
	s.NoError(err)
	s.sim = sim

	userState := models.UserState{
		AccountIndex: 1,
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(420),
		Nonce:        models.MakeUint256(0),
	}

	tree := st.NewStateTree(s.storage)
	err = tree.Set(1, &userState)
	s.NoError(err)

	tx := &models.Transaction{
		FromIndex: 1,
		ToIndex:   2,
		Amount:    *models.NewUint256(50),
		Fee:       *models.NewUint256(10),
		Nonce:     *models.NewUint256(0),
		Signature: []byte{1, 2, 3, 4},
	}

	s.tx = tx
}

func (s *CalculateTransactionStatusTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *CalculateTransactionStatusTestSuite) TestApi_CalculateTransactionStatus_Pending() {
	status, err := CalculateTransactionStatus(s.storage, s.tx, 0)
	s.NoError(err)

	s.Equal(models.Pending.Message(), status.Message())
}

func (s *CalculateTransactionStatusTestSuite) TestApi_CalculateTransactionStatus_Committed() {
	commitmentID, err := s.storage.AddCommitment(&commitment)
	s.NoError(err)

	s.tx.IncludedInCommitment = commitmentID

	status, err := CalculateTransactionStatus(s.storage, s.tx, 0)
	s.NoError(err)

	s.Equal(models.Committed.Message(), status.Message())
}

func (s *CalculateTransactionStatusTestSuite) TestApi_CalculateTransactionStatus_InBatch() {
	batch := models.Batch{
		Hash:              utils.RandomHash(),
		FinalisationBlock: math.MaxUint32,
	}
	err := s.storage.AddBatch(&batch)
	s.NoError(err)

	includedCommitment := commitment
	includedCommitment.IncludedInBatch = &batch.Hash
	commitmentID, err := s.storage.AddCommitment(&includedCommitment)
	s.NoError(err)

	s.tx.IncludedInCommitment = commitmentID

	status, err := CalculateTransactionStatus(s.storage, s.tx, 0)
	s.NoError(err)

	s.Equal(models.InBatch.Message(), status.Message())
}

// nolint:misspell
func (s *CalculateTransactionStatusTestSuite) TestApi_CalculateTransactionStatus_Finalised() {
	currentBlockNumber, err := s.sim.GetLatestBlockNumber()
	s.NoError(err)
	batch := models.Batch{
		Hash:              utils.RandomHash(),
		FinalisationBlock: *currentBlockNumber + 1,
	}
	err = s.storage.AddBatch(&batch)
	s.NoError(err)

	includedCommitment := commitment
	includedCommitment.IncludedInBatch = &batch.Hash
	commitmentID, err := s.storage.AddCommitment(&includedCommitment)
	s.NoError(err)

	s.tx.IncludedInCommitment = commitmentID

	s.sim.Commit()
	latestBlockNumber, err := s.sim.GetLatestBlockNumber()
	s.NoError(err)

	status, err := CalculateTransactionStatus(s.storage, s.tx, *latestBlockNumber)
	s.NoError(err)

	s.Equal(models.Finalised.Message(), status.Message())
}

func (s *CalculateTransactionStatusTestSuite) TestApi_CalculateTransactionStatus_Error() {
	s.tx.ErrorMessage = ref.String("Gold Duck Error")
	status, err := CalculateTransactionStatus(s.storage, s.tx, 0)
	s.NoError(err)

	s.Equal(models.Error.Message(), status.Message())
}

func TestCalculateTransactionStatusTestSuite(t *testing.T) {
	suite.Run(t, new(CalculateTransactionStatusTestSuite))
}
