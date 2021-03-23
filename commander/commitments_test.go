package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var genesisAccounts = []RegisteredGenesisAccount{
	{
		GenesisAccount: GenesisAccount{
			PublicKey: models.PublicKey{1, 2, 3},
			Balance:   models.MakeUint256(1000),
		},
		AccountIndex: 1,
	},
	{
		GenesisAccount: GenesisAccount{
			PublicKey: models.PublicKey{2, 3, 4},
			Balance:   models.MakeUint256(1000),
		},
		AccountIndex: 2,
	},
	{
		GenesisAccount: GenesisAccount{
			PublicKey: models.PublicKey{3, 4, 5},
			Balance:   models.MakeUint256(1000),
		},
		AccountIndex: 3,
	},
}

type CommitmentsLoopTestSuite struct {
	*require.Assertions
	suite.Suite
	db      *db.TestDB
	storage *storage.Storage
	cfg     *config.RollupConfig
}

func (s *CommitmentsLoopTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *CommitmentsLoopTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.db = testDB
	s.storage = storage.NewTestStorage(testDB.DB)
	s.cfg = &config.RollupConfig{
		TxsPerCommitment: 2,
	}
	err = PopulateGenesisAccounts(storage.NewStateTree(s.storage), genesisAccounts)
	s.NoError(err)
}

func (s *CommitmentsLoopTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *CommitmentsLoopTestSuite) TestCommitTransactions_ReturnsErrorWhenThereAreNotEnoughPendingTxs() {
	err := CommitTransactions(s.storage, s.cfg)
	s.ErrorIs(err, ErrNotEnoughTransactions)
}

func TestCommitmentsLoopTestSuite(t *testing.T) {
	suite.Run(t, new(CommitmentsLoopTestSuite))
}
