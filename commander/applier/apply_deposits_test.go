package applier

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ApplyDepositsTestSuite struct {
	*require.Assertions
	suite.Suite
	storage *st.TestStorage
	applier *Applier
}

func (s *ApplyDepositsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *ApplyDepositsTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.applier = NewApplier(s.storage.Storage)
}

func (s *ApplyDepositsTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *ApplyDepositsTestSuite) TestApplyDeposits() {
	deposits := make([]models.PendingDeposit, 4)
	for i := range deposits {
		deposits[i] = models.PendingDeposit{
			ID:         models.DepositID{BlockNumber: 1, LogIndex: uint32(i)},
			ToPubKeyID: uint32(i),
			TokenID:    models.MakeUint256(uint64(i)),
			L2Amount:   models.MakeUint256(uint64(100 + i)),
		}
	}

	startStateID := uint32(1)
	err := s.applier.ApplyDeposits(startStateID, deposits)
	s.NoError(err)

	for i := range deposits {
		leaf, err := s.applier.storage.StateTree.Leaf(startStateID + uint32(i))
		s.NoError(err)
		s.Equal(deposits[i].ToPubKeyID, leaf.PubKeyID)
		s.Equal(deposits[i].TokenID, leaf.TokenID)
		s.Equal(deposits[i].L2Amount, leaf.Balance)
		s.Equal(models.MakeUint256(0), leaf.Nonce)
	}
}

func TestApplyDepositsTestSuite(t *testing.T) {
	suite.Run(t, new(ApplyDepositsTestSuite))
}
