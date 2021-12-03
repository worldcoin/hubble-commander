package applier

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ApplyFeeTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *st.TestStorage
	applier *Applier
}

func (s *ApplyFeeTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ApplyFeeTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.applier = NewApplier(s.storage.Storage)
}

func (s *ApplyFeeTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ApplyFeeTestSuite) TestApplyFee() {
	feeReceiverStateID := receiverState.PubKeyID
	_, err := s.storage.StateTree.Set(feeReceiverStateID, &receiverState)
	s.NoError(err)

	stateProof, err := s.applier.ApplyFee(feeReceiverStateID, models.MakeUint256(555))
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

	stateProof, txError, appError := s.applier.ApplyFeeForSync(feeReceiver, models.NewUint256(2), models.NewUint256(555))
	s.NoError(appError)
	s.ErrorIs(txError, ErrInvalidFeeReceiverTokenID)
	s.Equal(receiverState, *stateProof.UserState)

	feeReceiverState, err := s.storage.StateTree.Leaf(feeReceiver)
	s.NoError(err)

	s.Equal(uint64(555), feeReceiverState.Balance.Uint64())
}

func TestApplyFeeTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyFeeTestSuite))
}
