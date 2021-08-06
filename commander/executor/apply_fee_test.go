package executor

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ApplyFeeTestSuite struct {
	*require.Assertions
	suite.Suite
	storage             *storage.TestStorage
	transactionExecutor *TransactionExecutor
}

func (s *ApplyFeeTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ApplyFeeTestSuite) SetupTest() {
	var err error
	s.storage, err = storage.NewTestStorageWithBadger()
	s.NoError(err)
	s.transactionExecutor = NewTestTransactionExecutor(s.storage.Storage, nil, nil, context.Background())
}

func (s *ApplyFeeTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ApplyFeeTestSuite) TestApplyFee() {
	feeReceiverStateID := receiverState.PubKeyID
	_, err := s.storage.StateTree.Set(feeReceiverStateID, &receiverState)
	s.NoError(err)

	stateProof, err := s.transactionExecutor.ApplyFee(feeReceiverStateID, models.MakeUint256(555))
	s.NoError(err)
	s.Equal(receiverState, *stateProof.UserState)

	feeReceiverState, err := s.storage.StateTree.Leaf(feeReceiverStateID)
	s.NoError(err)

	s.Equal(uint64(555), feeReceiverState.Balance.Uint64())
}

func (s *ApplyFeeTestSuite) TestApplyFeeForSync_InvalidTokenID() {
	feeReceiver := &FeeReceiver{
		StateID: receiverState.PubKeyID,
		TokenID: receiverState.TokenID,
	}
	_, err := s.storage.StateTree.Set(feeReceiver.StateID, &receiverState)
	s.NoError(err)

	stateProof, transferError, appError := s.transactionExecutor.ApplyFeeForSync(feeReceiver, models.NewUint256(2), models.NewUint256(555))
	s.NoError(appError)
	s.ErrorIs(transferError, ErrInvalidFeeReceiverTokenID)
	s.Equal(receiverState, *stateProof.UserState)

	feeReceiverState, err := s.storage.StateTree.Leaf(feeReceiver.StateID)
	s.NoError(err)

	s.Equal(uint64(555), feeReceiverState.Balance.Uint64())
}

func TestApplyFeeTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyFeeTestSuite))
}
