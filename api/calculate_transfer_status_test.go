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

type CalculateTransferStatusTestSuite struct {
	*require.Assertions
	suite.Suite
	db       *db.TestDB
	storage  *st.Storage
	sim      *simulator.Simulator
	transfer *models.Transfer
}

func (s *CalculateTransferStatusTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *CalculateTransferStatusTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)

	s.storage = st.NewTestStorage(testDB.DB)
	s.db = testDB

	sim, err := simulator.NewSimulator()
	s.NoError(err)
	s.sim = sim

	userState := models.UserState{
		PubkeyID:   1,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(420),
		Nonce:      models.MakeUint256(0),
	}

	tree := st.NewStateTree(s.storage)
	err = tree.Set(1, &userState)
	s.NoError(err)

	transfer := &models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      *models.NewUint256(50),
			Fee:         *models.NewUint256(10),
			Nonce:       *models.NewUint256(0),
			Signature:   []byte{1, 2, 3, 4},
		},
		ToStateID: 2,
	}

	s.transfer = transfer
}

func (s *CalculateTransferStatusTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *CalculateTransferStatusTestSuite) TestApi_CalculateTransferStatus_Pending() {
	status, err := CalculateTransferStatus(s.storage, s.transfer, 0)
	s.NoError(err)

	s.Equal(models.Pending, *status)
}

func (s *CalculateTransferStatusTestSuite) TestApi_CalculateTransferStatus_InBatch() {
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

	s.transfer.IncludedInCommitment = commitmentID

	status, err := CalculateTransferStatus(s.storage, s.transfer, 0)
	s.NoError(err)

	s.Equal(models.InBatch, *status)
}

// nolint:misspell
func (s *CalculateTransferStatusTestSuite) TestApi_CalculateTransferStatus_Finalised() {
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

	s.transfer.IncludedInCommitment = commitmentID

	s.sim.Commit()
	latestBlockNumber, err := s.sim.GetLatestBlockNumber()
	s.NoError(err)

	status, err := CalculateTransferStatus(s.storage, s.transfer, *latestBlockNumber)
	s.NoError(err)

	s.Equal(models.Finalised, *status)
}

func (s *CalculateTransferStatusTestSuite) TestApi_CalculateTransferStatus_Error() {
	s.transfer.ErrorMessage = ref.String("Gold Duck Error")
	status, err := CalculateTransferStatus(s.storage, s.transfer, 0)
	s.NoError(err)

	s.Equal(models.Error, *status)
}

func TestCalculateTransferStatusTestSuite(t *testing.T) {
	suite.Run(t, new(CalculateTransferStatusTestSuite))
}
