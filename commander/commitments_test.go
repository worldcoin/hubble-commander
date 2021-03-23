package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	sender = RegisteredGenesisAccount{
		GenesisAccount: GenesisAccount{
			PublicKey: models.PublicKey{1, 2, 3},
			Balance:   models.MakeUint256(1000),
		},
		AccountIndex: 0,
	}
	receiver = RegisteredGenesisAccount{
		GenesisAccount: GenesisAccount{
			PublicKey: models.PublicKey{2, 3, 4},
			Balance:   models.MakeUint256(1000),
		},
		AccountIndex: 1,
	}
	feeReceiver = RegisteredGenesisAccount{
		GenesisAccount: GenesisAccount{
			PublicKey: models.PublicKey{3, 4, 5},
			Balance:   models.MakeUint256(1000),
		},
		AccountIndex: 2,
	}
	genesisAccounts = []RegisteredGenesisAccount{sender, receiver, feeReceiver}
)

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
		FeeReceiverIndex: 2,
	}
	err = PopulateGenesisAccounts(storage.NewStateTree(s.storage), genesisAccounts)
	s.NoError(err)
}

func (s *CommitmentsLoopTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *CommitmentsLoopTestSuite) TestCommitTransactions_ReturnsErrorWhenThereAreNotEnoughTxs() {
	err := CommitTransactions(s.storage, s.cfg)
	s.ErrorIs(err, ErrNotEnoughTransactions)
}

func (s *CommitmentsLoopTestSuite) addTransactions(txs []models.Transaction) {
	for i := range txs {
		err := s.storage.AddTransaction(&txs[i])
		s.NoError(err)
	}
}

func (s *CommitmentsLoopTestSuite) TestCommitTransactions_ReturnsErrorWhenThereAreNotEnoughPendingTxs() {
	txs := generateValidTransactions(2)
	txs[1].ErrorMessage = ref.String("some error")
	s.addTransactions(txs)

	err := CommitTransactions(s.storage, s.cfg)
	s.ErrorIs(err, ErrNotEnoughTransactions)
}

func (s *CommitmentsLoopTestSuite) TestCommitTransactions_ReturnsErrorWhenThereAreNotEnoughValidTxs() {
	txs := generateInvalidTransactions(2)
	s.addTransactions(txs)

	err := CommitTransactions(s.storage, s.cfg)
	s.ErrorIs(err, ErrNotEnoughTransactions)
}

func TestCommitmentsLoopTestSuite(t *testing.T) {
	suite.Run(t, new(CommitmentsLoopTestSuite))
}
