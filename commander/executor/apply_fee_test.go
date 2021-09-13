package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ApplyFeeTestSuite struct {
	*require.Assertions
	suite.Suite
	storage      *storage.TestStorage
	executionCtx *ExecutionContext
}

func (s *ApplyFeeTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ApplyFeeTestSuite) SetupTest() {
	var err error
	s.storage, err = storage.NewTestStorage()
	s.NoError(err)
	s.executionCtx = NewTestExecutionContext(s.storage.Storage, nil, nil)
}

func (s *ApplyFeeTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ApplyFeeTestSuite) TestApplyFee() {
	feeReceiverStateID := receiverState.PubKeyID
	_, err := s.storage.StateTree.Set(feeReceiverStateID, &receiverState)
	s.NoError(err)

	stateProof, err := s.executionCtx.ApplyFee(feeReceiverStateID, models.MakeUint256(555))
	s.NoError(err)
	s.Equal(receiverState, *stateProof.UserState)

	feeReceiverState, err := s.storage.StateTree.Leaf(feeReceiverStateID)
	s.NoError(err)

	s.Equal(uint64(555), feeReceiverState.Balance.Uint64())
}

func (s *ApplyFeeTestSuite) TestApplyFeeForSync_InvalidTokenID() {
	feeReceiver := receiverState.PubKeyID
	_, err := s.storage.StateTree.Set(feeReceiver, &receiverState)
	s.NoError(err)

	stateProof, transferError, appError := s.executionCtx.ApplyFeeForSync(feeReceiver, models.NewUint256(2), models.NewUint256(555))
	s.NoError(appError)
	s.ErrorIs(transferError, ErrInvalidFeeReceiverTokenID)
	s.Equal(receiverState, *stateProof.UserState)

	feeReceiverState, err := s.storage.StateTree.Leaf(feeReceiver)
	s.NoError(err)

	s.Equal(uint64(555), feeReceiverState.Balance.Uint64())
}

func TestApplyFeeTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyFeeTestSuite))
}
