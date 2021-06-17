package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Create2TransferCommitmentsTestSuite struct {
	*require.Assertions
	suite.Suite
	teardown            func() error
	storage             *storage.Storage
	client              *eth.TestClient
	cfg                 *config.RollupConfig
	transactionExecutor *transactionExecutor
}

func (s *Create2TransferCommitmentsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *Create2TransferCommitmentsTestSuite) SetupTest() {
	testStorage, err := storage.NewTestStorageWithBadger()
	s.NoError(err)
	s.client, err = eth.NewTestClient()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.teardown = testStorage.Teardown
	s.cfg = &config.RollupConfig{
		TxsPerCommitment:       2,
		FeeReceiverPubKeyID:    2,
		MaxCommitmentsPerBatch: 1,
	}

	err = PopulateGenesisAccounts(s.storage, AssignStateIDs(genesisAccounts))
	s.NoError(err)

	s.transactionExecutor = newTestTransactionExecutor(s.storage, s.client.Client, s.cfg, transactionExecutorOpts{})
}

func (s *Create2TransferCommitmentsTestSuite) TearDownTest() {
	s.client.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_QueriesForMorePendingTransfersUntilSatisfied() {
	transfers := generateValidCreate2Transfers(6, &models.PublicKey{1, 2, 3})

	for i := 1; i < 6; i++ {
		transfers[i].Amount = models.MakeUint256(99999999999)
	}

	transfer := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.Create2Transfer,
			FromStateID: 24,
			Amount:      models.MakeUint256(1),
			Fee:         models.MakeUint256(1),
			Nonce:       models.MakeUint256(10),
		},
		ToStateID:   nil,
		ToPublicKey: models.PublicKey{5, 4, 3, 2, 1},
	}
	transfers = append(transfers, transfer)

	s.addCreate2Transfers(transfers)

	err := addNewDummyState(s.storage, s.transactionExecutor.stateTree, 24)
	s.NoError(err)

	pendingTransfers, err := s.storage.GetPendingCreate2Transfers(pendingTxsCountMultiplier * s.cfg.TxsPerCommitment)
	s.NoError(err)
	s.Len(pendingTransfers, 4)

	preRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	commitments, err := s.transactionExecutor.createCreate2TransferCommitments(pendingTransfers, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	postRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_DoesNothingWhenThereAreNotEnoughPendingTransfers() {
	preRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	commitments, err := s.transactionExecutor.createCreate2TransferCommitments([]models.Create2Transfer{}, testDomain)
	s.NoError(err)
	s.Len(commitments, 0)

	postRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_DoesNothingWhenThereAreNotEnoughValidTransfers() {
	transfers := generateValidCreate2Transfers(2, &models.PublicKey{1, 2, 3})
	transfers[1].Amount = models.MakeUint256(99999999999)
	s.addCreate2Transfers(transfers)

	pendingTransfers, err := s.storage.GetPendingCreate2Transfers(32)
	s.NoError(err)
	s.Len(pendingTransfers, 2)

	preRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	commitments, err := s.transactionExecutor.createCreate2TransferCommitments(pendingTransfers, testDomain)
	s.NoError(err)
	s.Len(commitments, 0)

	postRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_StoresCorrectCommitment() {
	pendingTransfers := s.prepareAndReturnPendingCreate2Transfers(2)

	preRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	commitments, err := s.transactionExecutor.createCreate2TransferCommitments(pendingTransfers, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, 32)
	s.Equal(commitments[0].FeeReceiver, uint32(2))
	s.Nil(commitments[0].IncludedInBatch)

	postRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_CreatesMaximallyAsManyCommitmentsAsSpecifiedInConfig() {
	pendingTransfers := s.prepareAndReturnPendingCreate2Transfers(2)

	commitments, err := s.transactionExecutor.createCreate2TransferCommitments(pendingTransfers, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_UpdateTransfers() {
	pendingTransfers := s.prepareAndReturnPendingCreate2Transfers(2)

	commitments, err := s.transactionExecutor.createCreate2TransferCommitments(pendingTransfers, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	for i := range pendingTransfers {
		tx, err := s.storage.GetCreate2Transfer(pendingTransfers[i].Hash)
		s.NoError(err)
		s.Equal(int32(1), *tx.IncludedInCommitment)
		s.Equal(uint32(i+3), *tx.ToStateID)
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

func (s *Create2TransferCommitmentsTestSuite) prepareAndReturnPendingCreate2Transfers(transfersAmount uint64) []models.Create2Transfer {
	transfers := generateValidCreate2Transfers(transfersAmount, &models.PublicKey{1, 2, 3})
	s.addCreate2Transfers(transfers)

	pendingTransfers, err := s.storage.GetPendingCreate2Transfers(transfersAmount)
	s.NoError(err)
	s.Len(pendingTransfers, int(transfersAmount))

	return pendingTransfers
}
