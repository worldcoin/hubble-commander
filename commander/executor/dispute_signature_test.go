package executor

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DisputeSignatureTestSuite struct {
	*require.Assertions
	suite.Suite
	storage             *st.Storage
	teardown            func() error
	client              *eth.TestClient
	cfg                 *config.RollupConfig
	transactionExecutor *TransactionExecutor
}

func (s *DisputeSignatureTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.cfg = &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    1,
		DevMode:                false,
	}
}

func (s *DisputeSignatureTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.teardown = testStorage.Teardown

	s.client, err = eth.NewTestClient()
	s.NoError(err)

	s.transactionExecutor = NewTestTransactionExecutor(s.storage, s.client.Client, s.cfg, context.Background())
}

func (s *DisputeSignatureTestSuite) TearDownTest() {
	s.client.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *DisputeSignatureTestSuite) TestGetUserStateProof() {
	userState := createUserState(1, 300, 0)
	witness, err := s.transactionExecutor.stateTree.Set(1, userState)
	s.NoError(err)

	stateProof, err := s.transactionExecutor.getUserStateProof(1)
	s.NoError(err)
	s.Equal(userState, stateProof.UserState)
	s.Equal(witness, stateProof.Witness)
}

func TestDisputeSignatureTestSuite(t *testing.T) {
	suite.Run(t, new(DisputeSignatureTestSuite))
}
