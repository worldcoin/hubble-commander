package api

import (
	"math"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	commitment = models.TxCommitment{
		CommitmentBase: models.CommitmentBase{
			Type:          batchtype.Transfer,
			PostStateRoot: utils.RandomHash(),
		},
		FeeReceiver:       1,
		CombinedSignature: models.MakeRandomSignature(),
	}
)

type CalculateTransactionStatusTestSuite struct {
	*require.Assertions
	suite.Suite
	storage  *st.TestStorage
	sim      *simulator.Simulator
	transfer *models.Transfer
}

func (s *CalculateTransactionStatusTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *CalculateTransactionStatusTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	sim, err := simulator.NewSimulator()
	s.NoError(err)
	s.sim = sim

	userState := models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	}

	_, err = s.storage.StateTree.Set(1, &userState)
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
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *CalculateTransactionStatusTestSuite) TestCalculateTransactionStatus_TxInMempool() {
	status, err := CalculateTransactionStatus(s.storage.Storage, &s.transfer.TransactionBase, 0)
	s.NoError(err)

	s.Equal(txstatus.Pending, *status)
}

func (s *CalculateTransactionStatusTestSuite) TestCalculateTransactionStatus_TxInPendingBatch() {
	batch := models.Batch{
		ID: models.MakeUint256(1),
	}
	err := s.storage.AddBatch(&batch)
	s.NoError(err)

	s.transfer.CommitmentID = &models.CommitmentID{
		BatchID:      batch.ID,
		IndexInBatch: 0,
	}

	status, err := CalculateTransactionStatus(s.storage.Storage, &s.transfer.TransactionBase, 0)
	s.NoError(err)

	s.Equal(txstatus.Pending, *status)
}

func (s *CalculateTransactionStatusTestSuite) TestCalculateTransactionStatus_InBatch() {
	batch := models.Batch{
		ID:                models.MakeUint256(1),
		FinalisationBlock: ref.Uint32(math.MaxUint32),
	}
	err := s.storage.AddBatch(&batch)
	s.NoError(err)

	s.transfer.CommitmentID = &models.CommitmentID{
		BatchID: batch.ID,
	}

	status, err := CalculateTransactionStatus(s.storage.Storage, &s.transfer.TransactionBase, 0)
	s.NoError(err)

	s.Equal(txstatus.Mined, *status)
}

// nolint:misspell
func (s *CalculateTransactionStatusTestSuite) TestCalculateTransactionStatus_Finalised() {
	currentBlockNumber, err := s.sim.GetLatestBlockNumber()
	s.NoError(err)
	batch := models.Batch{
		ID:                models.MakeUint256(1),
		FinalisationBlock: ref.Uint32(uint32(*currentBlockNumber) + 1),
	}
	err = s.storage.AddBatch(&batch)
	s.NoError(err)

	s.transfer.CommitmentID = &models.CommitmentID{
		BatchID: batch.ID,
	}

	s.sim.GetBackend().Commit()
	latestBlockNumber, err := s.sim.GetLatestBlockNumber()
	s.NoError(err)

	status, err := CalculateTransactionStatus(s.storage.Storage, &s.transfer.TransactionBase, uint32(*latestBlockNumber))
	s.NoError(err)

	s.Equal(txstatus.Finalised, *status)
}

func (s *CalculateTransactionStatusTestSuite) TestCalculateTransactionStatus_Error() {
	s.transfer.ErrorMessage = ref.String("Gold Duck Error")
	status, err := CalculateTransactionStatus(s.storage.Storage, &s.transfer.TransactionBase, 0)
	s.NoError(err)

	s.Equal(txstatus.Error, *status)
}

func TestCalculateTransactionStatusTestSuite(t *testing.T) {
	suite.Run(t, new(CalculateTransactionStatusTestSuite))
}
