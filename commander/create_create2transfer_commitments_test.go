package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Create2TransferCommitmentsTestSuite struct {
	*require.Assertions
	suite.Suite
	teardown func() error
	storage  *storage.Storage
	cfg      *config.RollupConfig
}

func (s *Create2TransferCommitmentsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *Create2TransferCommitmentsTestSuite) SetupTest() {
	testStorage, err := storage.NewTestStorageWithBadger()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.teardown = testStorage.Teardown
	s.cfg = &config.RollupConfig{
		TxsPerCommitment:       2,
		FeeReceiverPubKeyID:    2,
		MaxCommitmentsPerBatch: 1,
	}
	err = PopulateGenesisAccounts(s.storage, genesisAccounts)
	s.NoError(err)
}

func (s *Create2TransferCommitmentsTestSuite) TearDownTest() {
	err := s.teardown()
	s.NoError(err)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_DoesNothingWhenThereAreNotEnoughPendingTransfers() {
	preRoot, err := storage.NewStateTree(s.storage).Root()
	s.NoError(err)

	commitments, err := createTransferCommitments([]models.Transfer{}, s.storage, s.cfg)
	s.NoError(err)
	s.Len(commitments, 0)

	postRoot, err := storage.NewStateTree(s.storage).Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_DoesNothingWhenThereAreNotEnoughValidTransfers() {
	transfers := generateValidCreate2Transfers(2)
	transfers[1].Amount = models.MakeUint256(99999999999)
	s.addCreate2Transfers(transfers)

	pendingTransfers, err := s.storage.GetPendingCreate2Transfers()
	s.NoError(err)
	s.Len(pendingTransfers, 2)

	preRoot, err := storage.NewStateTree(s.storage).Root()
	s.NoError(err)

	commitments, err := createCreate2TransferCommitments(pendingTransfers, s.storage, s.cfg)
	s.NoError(err)
	s.Len(commitments, 0)

	postRoot, err := storage.NewStateTree(s.storage).Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_StoresCorrectCommitment() {
	pendingTransfers := s.prepareAndReturnPendingCreate2Transfers(2)

	preRoot, err := storage.NewStateTree(s.storage).Root()
	s.NoError(err)

	commitments, err := createCreate2TransferCommitments(pendingTransfers, s.storage, s.cfg)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, 32)
	s.Equal(commitments[0].FeeReceiver, uint32(2))
	s.Nil(commitments[0].AccountTreeRoot)
	s.Nil(commitments[0].IncludedInBatch)

	postRoot, err := storage.NewStateTree(s.storage).Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_CreatesMaximallyAsManyCommitmentsAsSpecifiedInConfig() {
	pendingTransfers := s.prepareAndReturnPendingCreate2Transfers(2)

	commitments, err := createCreate2TransferCommitments(pendingTransfers, s.storage, s.cfg)
	s.NoError(err)
	s.Len(commitments, 1)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_MarksTransfersAsIncludedInCommitment() {
	pendingTransfers := s.prepareAndReturnPendingCreate2Transfers(2)

	commitments, err := createCreate2TransferCommitments(pendingTransfers, s.storage, s.cfg)
	s.NoError(err)
	s.Len(commitments, 1)

	for i := range pendingTransfers {
		tx, err := s.storage.GetCreate2Transfer(pendingTransfers[i].Hash)
		s.NoError(err)
		s.Equal(*tx.IncludedInCommitment, int32(1))
	}
}

func (s *Create2TransferCommitmentsTestSuite) TestRemoveCreate2Transfer() {
	transfer1 := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash: utils.RandomHash(),
		},
	}
	transfer2 := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash: utils.RandomHash(),
		},
	}
	transfer3 := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash: utils.RandomHash(),
		},
	}

	transfers := []models.Create2Transfer{transfer1, transfer2, transfer3}
	toRemove := []models.Create2Transfer{transfer2}

	s.Equal([]models.Create2Transfer{transfer1, transfer3}, removeCreate2Transfer(transfers, toRemove))
}

func TestCreate2TransferCommitmentsTestSuite(t *testing.T) {
	suite.Run(t, new(Create2TransferCommitmentsTestSuite))
}

func (s *Create2TransferCommitmentsTestSuite) addCreate2Transfers(transfers []models.Create2Transfer) {
	for i := range transfers {
		err := s.storage.AddCreate2Transfer(&transfers[i])
		s.NoError(err)
	}
}

func (s *Create2TransferCommitmentsTestSuite) prepareAndReturnPendingCreate2Transfers(transfersAmount int) []models.Create2Transfer {
	transfers := generateValidCreate2Transfers(transfersAmount)
	s.addCreate2Transfers(transfers)

	pendingTransfers, err := s.storage.GetPendingCreate2Transfers()
	s.NoError(err)
	s.Len(pendingTransfers, transfersAmount)

	return pendingTransfers
}
