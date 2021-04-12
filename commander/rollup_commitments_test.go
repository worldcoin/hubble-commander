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

type RollupLoopCommitmentsTestSuite struct {
	*require.Assertions
	suite.Suite
	db        *db.TestDB
	storage   *storage.Storage
	cfg       *config.RollupConfig
}

func (s *RollupLoopCommitmentsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *RollupLoopCommitmentsTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.db = testDB
	s.storage = storage.NewTestStorage(testDB.DB)
	s.cfg = &config.RollupConfig{
		TxsPerCommitment:       2,
		FeeReceiverIndex:       2,
		MaxCommitmentsPerBatch: 1,
	}
	err = PopulateGenesisAccounts(storage.NewStateTree(s.storage), genesisAccounts)

	s.NoError(err)
}

func (s *RollupLoopCommitmentsTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *RollupLoopCommitmentsTestSuite) Test_CommitmentsLoop_DoNothingWhenThereAreNotEnoughPendingTxs() {
	commitments, err := CommitmentsLoop([]models.Transaction{}, s.storage, s.cfg)
	s.NoError(err)
	s.Len(commitments, 0)
}

func (s *RollupLoopCommitmentsTestSuite) Test_CommitmentsLoop_DoNothingWhenThereAreNotEnoughValidTxs() {
	txs := generateValidTransactions(2)
	txs[1].Amount = models.MakeUint256(99999999999)
	s.addTransactions(txs)

	pendingTransactions, err := s.storage.GetPendingTransactions()
	s.NoError(err)
	s.Len(pendingTransactions, 2)
	
	commitments, err := CommitmentsLoop(pendingTransactions, s.storage, s.cfg)
	s.NoError(err)
	s.Len(commitments, 0)
}


func (s *RollupLoopCommitmentsTestSuite) Test_CommitmentsLoop_StoresCorrectCommitment() {
	pendingTransactions := s.prepareAndReturnPendingTransactions(3)

	commitments, err := CommitmentsLoop(pendingTransactions, s.storage, s.cfg)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, 24)
	s.Equal(commitments[0].FeeReceiver, uint32(2))
	s.Nil(commitments[0].AccountTreeRoot)
	s.Nil(commitments[0].IncludedInBatch)

	root, err := storage.NewStateTree(s.storage).Root()
	s.NoError(err)
	s.Equal(commitments[0].PostStateRoot, *root)
}

func (s *RollupLoopCommitmentsTestSuite) Test_CommitmentsLoop_MarksTransactionsAsIncludedInCommitment() {
	pendingTransactions := s.prepareAndReturnPendingTransactions(2)

	commitments, err := CommitmentsLoop(pendingTransactions, s.storage, s.cfg)
	s.NoError(err)
	s.Len(commitments, 1)

	for i := range pendingTransactions {
		tx, err := s.storage.GetTransaction(pendingTransactions[i].Hash)
		s.NoError(err)
		s.Equal(*tx.IncludedInCommitment, int32(1))
	}
}

func TestRollupLoopCommitmentsTestSuite(t *testing.T) {
	suite.Run(t, new(RollupLoopCommitmentsTestSuite))
}

func (s *RollupLoopCommitmentsTestSuite) addTransactions(txs []models.Transaction) {
	for i := range txs {
		err := s.storage.AddTransaction(&txs[i])
		s.NoError(err)
	}
}

func (s *RollupLoopCommitmentsTestSuite) prepareAndReturnPendingTransactions(txAmount int) []models.Transaction {
	txs := generateValidTransactions(txAmount)
	s.addTransactions(txs)

	pendingTransactions, err := s.storage.GetPendingTransactions()
	s.NoError(err)
	s.Len(pendingTransactions, txAmount)

	return pendingTransactions
}

