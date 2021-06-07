package api

import (
	"math"
	"testing"

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

type CalculateTransactionStatusTestSuite struct {
	*require.Assertions
	suite.Suite
	teardown func() error
	storage  *st.Storage
	sim      *simulator.Simulator
	transfer *models.Transfer
}

func (s *CalculateTransactionStatusTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *CalculateTransactionStatusTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.teardown = testStorage.Teardown

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

func (s *CalculateTransactionStatusTestSuite) TearDownTest() {
	err := s.teardown()
	s.NoError(err)
}

func (s *CalculateTransactionStatusTestSuite) TestCalculateTransactionStatus_Pending() {
	status, err := CalculateTransactionStatus(s.storage, &s.transfer.TransactionBase, 0)
	s.NoError(err)

	s.Equal(txstatus.Pending, *status)
}

func (s *CalculateTransactionStatusTestSuite) TestCalculateTransactionStatus_InBatch() {
	batch := models.Batch{
		FinalisationBlock: ref.Uint32(math.MaxUint32),
	}
	batchID, err := s.storage.AddBatch(&batch)
	s.NoError(err)

	includedCommitment := commitment
	includedCommitment.IncludedInBatch = batchID
	commitmentID, err := s.storage.AddCommitment(&includedCommitment)
	s.NoError(err)

	s.transfer.IncludedInCommitment = commitmentID

	status, err := CalculateTransactionStatus(s.storage, &s.transfer.TransactionBase, 0)
	s.NoError(err)

	s.Equal(txstatus.InBatch, *status)
}

// nolint:misspell
func (s *CalculateTransactionStatusTestSuite) TestCalculateTransactionStatus_Finalised() {
	currentBlockNumber, err := s.sim.GetLatestBlockNumber()
	s.NoError(err)
	batch := models.Batch{
		FinalisationBlock: ref.Uint32(uint32(*currentBlockNumber) + 1),
	}
	batchID, err := s.storage.AddBatch(&batch)
	s.NoError(err)

	includedCommitment := commitment
	includedCommitment.IncludedInBatch = batchID
	commitmentID, err := s.storage.AddCommitment(&includedCommitment)
	s.NoError(err)

	s.transfer.IncludedInCommitment = commitmentID

	s.sim.Commit()
	latestBlockNumber, err := s.sim.GetLatestBlockNumber()
	s.NoError(err)

	status, err := CalculateTransactionStatus(s.storage, &s.transfer.TransactionBase, uint32(*latestBlockNumber))
	s.NoError(err)

	s.Equal(txstatus.Finalised, *status)
}

func (s *CalculateTransactionStatusTestSuite) TestCalculateTransactionStatus_Error() {
	s.transfer.ErrorMessage = ref.String("Gold Duck Error")
	status, err := CalculateTransactionStatus(s.storage, &s.transfer.TransactionBase, 0)
	s.NoError(err)

	s.Equal(txstatus.Error, *status)
}

func TestCalculateTransactionStatusTestSuite(t *testing.T) {
	suite.Run(t, new(CalculateTransactionStatusTestSuite))
}
