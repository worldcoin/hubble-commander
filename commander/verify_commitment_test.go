package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type VerifyCommitmentTestSuite struct {
	*require.Assertions
	suite.Suite
	storage  *st.Storage
	teardown func() error
	client   *eth.TestClient
}

func (s *VerifyCommitmentTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *VerifyCommitmentTestSuite) SetupTest() {
	storage, err := st.NewTestStorage()
	s.NoError(err)
	s.storage = storage.Storage
	s.teardown = storage.Teardown
	s.client, err = eth.NewTestClient()
	s.NoError(err)
}

func (s *VerifyCommitmentTestSuite) TearDownTest() {
	s.client.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *VerifyCommitmentTestSuite) TestVerifyCommitment_ValidCommitmentRoot() {
	err := PopulateGenesisAccounts(s.storage, s.client.ChainState.GenesisAccounts)
	s.NoError(err)

	err = verifyCommitmentRoot(s.storage, s.client.Client)
	s.NoError(err)
}

func (s *VerifyCommitmentTestSuite) TestVerifyCommitment_InvalidCommitmentRoot() {
	s.client.ChainState.GenesisAccounts = append(s.client.ChainState.GenesisAccounts, models.PopulatedGenesisAccount{
		PublicKey: models.PublicKey{5, 6, 7},
		StateID:   1,
		State: models.UserState{
			PubKeyID: 10,
			TokenID:  models.MakeUint256(0),
			Balance:  models.MakeUint256(500),
			Nonce:    models.MakeUint256(0),
		},
	})
	err := PopulateGenesisAccounts(s.storage, s.client.ChainState.GenesisAccounts)
	s.NoError(err)

	err = verifyCommitmentRoot(s.storage, s.client.Client)
	s.NotNil(err)
	s.Equal(ErrInvalidCommitmentRoot, err)
}

func TestVerifyCommitmentTestSuite(t *testing.T) {
	suite.Run(t, new(VerifyCommitmentTestSuite))
}
