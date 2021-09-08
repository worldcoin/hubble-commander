package executor

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DepositsTestSuite struct {
	*require.Assertions
	suite.Suite
	storage             *st.TestStorage
	transactionExecutor *TransactionExecutor
}

func (s *DepositsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *DepositsTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.transactionExecutor = NewTestTransactionExecutor(s.storage.Storage, &eth.Client{}, nil, context.Background())
}

func (s *DepositsTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *DepositsTestSuite) TestGetVacancyProof() {
	stateID, err := s.transactionExecutor.storage.StateTree.NextVacantSubtree(2)
	s.NoError(err)

	vacancyProof, err := s.transactionExecutor.GetVacancyProof(*stateID, 2)
	s.NoError(err)
	s.EqualValues(vacancyProof.PathAtDepth, 0)
	s.Len(vacancyProof.Witness, 30)
}

func TestDepositsTestSuite(t *testing.T) {
	suite.Run(t, new(DepositsTestSuite))
}
