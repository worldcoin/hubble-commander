package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	senderState = models.UserState{
		PubKeyID:   1,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(420),
		Nonce:      models.MakeUint256(0),
	}
	receiverState = models.UserState{
		PubKeyID:   2,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(0),
		Nonce:      models.MakeUint256(0),
	}
)

type ApplyTransferTestSuite struct {
	*require.Assertions
	suite.Suite
	db      *db.TestDB
	storage *storage.Storage
	tree    *storage.StateTree
}

func (s *ApplyTransferTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ApplyTransferTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.db = testDB
	s.storage = storage.NewTestStorage(testDB.DB)
	s.tree = storage.NewStateTree(s.storage)
}

func (s *ApplyTransferTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *ApplyTransferTestSuite) Test_CalculateStateAfterTransfer_UpdatesStates() {
	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(100),
			Fee:         models.MakeUint256(10),
			Nonce:       models.MakeUint256(0),
		},
		ToStateID: 2,
	}

	newSenderState, newReceiverState, err := CalculateStateAfterTransfer(
		&senderState,
		&receiverState,
		&transfer,
	)
	s.NoError(err)

	s.Equal(models.MakeUint256(1), newSenderState.Nonce)
	s.Equal(models.MakeUint256(310), newSenderState.Balance)

	s.Equal(models.MakeUint256(0), newReceiverState.Nonce)
	s.Equal(models.MakeUint256(100), newReceiverState.Balance)

	s.NotEqual(&newSenderState, &senderState)
	s.NotEqual(&newReceiverState, &receiverState)
}

func (s *ApplyTransferTestSuite) Test_CalculateStateAfterTransfer_Validation_Nonce() {
	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(100),
			Fee:         models.MakeUint256(10),
			Nonce:       models.MakeUint256(1),
		},
		ToStateID: 2,
	}

	_, _, err := CalculateStateAfterTransfer(&senderState, &receiverState, &transfer)
	s.Error(err)
}

func (s *ApplyTransferTestSuite) Test_CalculateStateAfterTransfer_Validation_Balance() {
	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(50),
			Nonce:       models.MakeUint256(0),
		},
		ToStateID: 2,
	}

	_, _, err := CalculateStateAfterTransfer(&senderState, &receiverState, &transfer)
	s.Error(err)
}

func (s *ApplyTransferTestSuite) Test_ApplyTransfer_Validation_Nil() {
	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(50),
			Nonce:       models.MakeUint256(0),
		},
		ToStateID: 2,
	}

	transferError, appError := ApplyTransfer(s.tree, nil, models.MakeUint256(1))
	s.Error(appError, ErrTransactionIsNil)
	s.NoError(transferError)
	transferError, appError = ApplyTransfer(nil, &transfer, models.MakeUint256(1))
	s.Equal(appError, ErrStateTreeIsNil)
	s.NoError(transferError)
}

func (s *ApplyTransferTestSuite) Test_ApplyTransfer_Validation_TokenIndex() {
	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(50),
			Nonce:       models.MakeUint256(0),
		},
		ToStateID: 2,
	}

	senderIndex := senderState.PubKeyID
	receiverIndex := receiverState.PubKeyID

	err := s.tree.Set(senderIndex, &senderState)
	s.NoError(err)
	err = s.tree.Set(receiverIndex, &receiverState)
	s.NoError(err)

	transferError, appError := ApplyTransfer(s.tree, &transfer, models.MakeUint256(3))
	s.Equal(appError, ErrIncorrectTokenIndices)
	s.NoError(transferError)
}

func (s *ApplyTransferTestSuite) Test_ApplyTransfer() {
	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(100),
			Fee:         models.MakeUint256(50),
			Nonce:       models.MakeUint256(0),
		},
		ToStateID: 2,
	}

	senderIndex := senderState.PubKeyID
	receiverIndex := receiverState.PubKeyID

	err := s.tree.Set(senderIndex, &senderState)
	s.NoError(err)
	err = s.tree.Set(receiverIndex, &receiverState)
	s.NoError(err)

	transferError, appError := ApplyTransfer(s.tree, &transfer, models.MakeUint256(1))
	s.NoError(appError)
	s.NoError(transferError)

	senderLeaf, err := s.tree.Leaf(senderIndex)
	s.NoError(err)
	receiverLeaf, err := s.tree.Leaf(receiverIndex)
	s.NoError(err)

	s.Equal(int64(270), senderLeaf.Balance.Int64())
	s.Equal(int64(100), receiverLeaf.Balance.Int64())
}

func (s *ApplyTransferTestSuite) Test_ApplyFee() {
	receiverIndex := receiverState.PubKeyID
	err := s.tree.Set(receiverIndex, &receiverState)
	s.NoError(err)

	err = ApplyFee(s.tree, receiverIndex, models.MakeUint256(555))
	s.NoError(err)

	receiverLeaf, err := s.tree.Leaf(receiverIndex)
	s.NoError(err)

	s.Equal(int64(555), receiverLeaf.Balance.Int64())
}

func TestApplyTransferTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyTransferTestSuite))
}
