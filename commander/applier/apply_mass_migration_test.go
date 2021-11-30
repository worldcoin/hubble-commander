package applier

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ApplyMassMigrationTestSuite struct {
	*require.Assertions
	suite.Suite
	storage       *st.TestStorage
	applier       *Applier
	massMigration models.MassMigration
	tokenID       models.Uint256
}

func (s *ApplyMassMigrationTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
	s.massMigration = models.MassMigration{
		TransactionBase: models.TransactionBase{
			FromStateID: 1,
			Amount:      models.MakeUint256(100),
			Fee:         models.MakeUint256(10),
			Nonce:       models.MakeUint256(0),
		},
		SpokeID: models.MakeUint256(2),
	}
}

func (s *ApplyMassMigrationTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.applier = NewApplier(s.storage.Storage, nil)

	_, err = s.storage.StateTree.Set(senderState.PubKeyID, &senderState)
	s.NoError(err)

	s.tokenID = senderState.TokenID
}

func (s *ApplyMassMigrationTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ApplyMassMigrationTestSuite) TestApplyMassMigration() {
	_, txError, appError := s.applier.ApplyMassMigration(&s.massMigration, s.tokenID)
	s.NoError(txError)
	s.NoError(appError)

	senderLeaf, err := s.storage.StateTree.Leaf(1)
	s.NoError(err)

	s.Equal(uint64(290), senderLeaf.Balance.Uint64())
}

func TestApplyMassMigrationTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyMassMigrationTestSuite))
}
