package api

import (
	"math"
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	commitment = models.Commitment{
		Type:              txtype.Transfer,
		Transactions:      utils.RandomBytes(24),
		FeeReceiver:       1,
		CombinedSignature: models.MakeRandomSignature(),
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
		PubKeyID:   1,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(420),
		Nonce:      models.MakeUint256(0),
	}

	err = s.storage.AddAccountIfNotExists(&models.Account{
		PubKeyID:  userState.PubKeyID,
		PublicKey: models.PublicKey{1, 2, 3},
	})
	s.NoError(err)

	tree := st.NewStateTree(s.storage)
	err = tree.Set(1, &userState)
	s.NoError(err)

	transfer := &models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(50),
			Fee:         models.MakeUint256(10),
			Nonce:       models.MakeUint256(0),
		},
		ToStateID: 2,
	}

	s.transfer = transfer
}

func (s *CalculateTransferStatusTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *CalculateTransferStatusTestSuite) TestCalculateTransferStatus_Pending() {
	status, err := CalculateTransferStatus(s.storage, s.transfer, 0)
	s.NoError(err)

	s.Equal(txstatus.Pending, *status)
}

func (s *CalculateTransferStatusTestSuite) TestCalculateTransferStatus_InBatch() {
	batch := models.Batch{
		Hash:              utils.RandomHash(),
		Type:              txtype.Transfer,
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

	s.Equal(txstatus.InBatch, *status)
}

// nolint:misspell
func (s *CalculateTransferStatusTestSuite) TestCalculateTransferStatus_Finalised() {
	currentBlockNumber, err := s.sim.GetLatestBlockNumber()
	s.NoError(err)
	batch := models.Batch{
		Hash:              utils.RandomHash(),
		Type:              txtype.Transfer,
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

	s.Equal(txstatus.Finalised, *status)
}

func (s *CalculateTransferStatusTestSuite) TestCalculateTransferStatus_Error() {
	s.transfer.ErrorMessage = ref.String("Gold Duck Error")
	status, err := CalculateTransferStatus(s.storage, s.transfer, 0)
	s.NoError(err)

	s.Equal(txstatus.Error, *status)
}

func TestCalculateTransferStatusTestSuite(t *testing.T) {
	suite.Run(t, new(CalculateTransferStatusTestSuite))
}
