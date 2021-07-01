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
	storage             *storage.TestStorage
	tree                *storage.StateTree
	transactionExecutor *TransactionExecutor
}

func (s *ApplyFeeTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ApplyFeeTestSuite) SetupTest() {
	var err error
	s.storage, err = storage.NewTestStorageWithBadger()
	s.NoError(err)
	s.tree = storage.NewStateTree(s.storage.Storage)
	s.transactionExecutor = NewTestTransactionExecutor(s.storage.Storage, nil, nil, TransactionExecutorOpts{})
}

func (s *ApplyFeeTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ApplyFeeTestSuite) TestApplyFee() {
	feeReceiverStateID := receiverState.PubKeyID
	_, err := s.tree.Set(feeReceiverStateID, &receiverState)
	s.NoError(err)

	err = s.transactionExecutor.ApplyFee(feeReceiverStateID, models.MakeUint256(555))
	s.NoError(err)

	feeReceiverState, err := s.storage.GetStateLeaf(feeReceiverStateID)
	s.NoError(err)

	s.Equal(uint64(555), feeReceiverState.Balance.Uint64())
}

func TestApplyFeeTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyFeeTestSuite))
}
