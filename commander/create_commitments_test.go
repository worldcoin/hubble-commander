package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	sender = models.RegisteredGenesisAccount{
		GenesisAccount: models.GenesisAccount{
			Balance: models.MakeUint256(1000),
		},
		PubKeyID: 0,
	}
	receiver = models.RegisteredGenesisAccount{
		GenesisAccount: models.GenesisAccount{
			Balance: models.MakeUint256(1000),
		},
		PubKeyID: 1,
	}
	feeReceiver = models.RegisteredGenesisAccount{
		GenesisAccount: models.GenesisAccount{
			Balance: models.MakeUint256(1000),
		},
		PubKeyID: 2,
	}
	genesisAccounts = []models.RegisteredGenesisAccount{sender, receiver, feeReceiver}
)

type CreateCommitmentsTestSuite struct {
	*require.Assertions
	suite.Suite
	db      *db.TestDB
	storage *storage.Storage
	cfg     *config.RollupConfig
}

func (s *CreateCommitmentsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *CreateCommitmentsTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	s.db = testDB
	s.storage = storage.NewTestStorage(testDB.DB)
	s.cfg = &config.RollupConfig{
		TxsPerCommitment:       2,
		FeeReceiverIndex:       2,
		MaxCommitmentsPerBatch: 1,
	}
	err = PopulateGenesisAccounts(s.storage, genesisAccounts)

	s.NoError(err)
}

func (s *CreateCommitmentsTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_DoesNothingWhenThereAreNotEnoughPendingTransfers() {
	preRoot, err := storage.NewStateTree(s.storage).Root()
	s.NoError(err)

	commitments, err := createCommitments([]models.Transfer{}, s.storage, s.cfg)
	s.NoError(err)
	s.Len(commitments, 0)

	postRoot, err := storage.NewStateTree(s.storage).Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_DoesNothingWhenThereAreNotEnoughValidTransfers() {
	transfers := generateValidTransfers(2)
	transfers[1].Amount = models.MakeUint256(99999999999)
	s.addTransfers(transfers)

	pendingTransfers, err := s.storage.GetPendingTransfers()
	s.NoError(err)
	s.Len(pendingTransfers, 2)

	preRoot, err := storage.NewStateTree(s.storage).Root()
	s.NoError(err)

	commitments, err := createCommitments(pendingTransfers, s.storage, s.cfg)
	s.NoError(err)
	s.Len(commitments, 0)

	postRoot, err := storage.NewStateTree(s.storage).Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_StoresCorrectCommitment() {
	pendingTransfers := s.prepareAndReturnPendingTransfers(3)

	preRoot, err := storage.NewStateTree(s.storage).Root()
	s.NoError(err)

	commitments, err := createCommitments(pendingTransfers, s.storage, s.cfg)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, 24)
	s.Equal(commitments[0].FeeReceiver, uint32(2))
	s.Nil(commitments[0].AccountTreeRoot)
	s.Nil(commitments[0].IncludedInBatch)

	postRoot, err := storage.NewStateTree(s.storage).Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_CreatesMaximallyAsManyCommitmentsAsSpecifiedInConfig() {
	pendingTransfers := s.prepareAndReturnPendingTransfers(2)

	commitments, err := createCommitments(pendingTransfers, s.storage, s.cfg)
	s.NoError(err)
	s.Len(commitments, 1)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_MarksTransfersAsIncludedInCommitment() {
	pendingTransfers := s.prepareAndReturnPendingTransfers(2)

	commitments, err := createCommitments(pendingTransfers, s.storage, s.cfg)
	s.NoError(err)
	s.Len(commitments, 1)

	for i := range pendingTransfers {
		tx, err := s.storage.GetTransfer(pendingTransfers[i].Hash)
		s.NoError(err)
		s.Equal(*tx.IncludedInCommitment, int32(1))
	}
}

func (s *CreateCommitmentsTestSuite) TestRemoveTransactions() {
	transfer1 := models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash: utils.RandomHash(),
		},
	}
	transfer2 := models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash: utils.RandomHash(),
		},
	}
	transfer3 := models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash: utils.RandomHash(),
		},
	}

	transfers := []models.Transfer{transfer1, transfer2, transfer3}
	toRemove := []models.Transfer{transfer2}

	s.Equal([]models.Transfer{transfer1, transfer3}, removeTransfer(transfers, toRemove))
}

func TestCreateCommitmentsTestSuite(t *testing.T) {
	suite.Run(t, new(CreateCommitmentsTestSuite))
}

func (s *CreateCommitmentsTestSuite) addTransfers(transfers []models.Transfer) {
	for i := range transfers {
		err := s.storage.AddTransfer(&transfers[i])
		s.NoError(err)
	}
}

func (s *CreateCommitmentsTestSuite) prepareAndReturnPendingTransfers(transfersAmount int) []models.Transfer {
	transfers := generateValidTransfers(transfersAmount)
	s.addTransfers(transfers)

	pendingTransfers, err := s.storage.GetPendingTransfers()
	s.NoError(err)
	s.Len(pendingTransfers, transfersAmount)

	return pendingTransfers
}
