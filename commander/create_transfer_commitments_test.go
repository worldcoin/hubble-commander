package commander

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
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
	testDomain      = &bls.Domain{1, 2, 3, 4}
)

type TransferCommitmentsTestSuite struct {
	*require.Assertions
	suite.Suite
	teardown            func() error
	storage             *storage.Storage
	cfg                 *config.RollupConfig
	transactionExecutor *transactionExecutor
}

func (s *TransferCommitmentsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TransferCommitmentsTestSuite) SetupTest() {
	testStorage, err := storage.NewTestStorageWithBadger()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.teardown = testStorage.Teardown
	s.cfg = &config.RollupConfig{
		TxsPerCommitment:          2,
		PendingTxsCountMultiplier: 2,
		FeeReceiverPubKeyID:       2,
		MaxCommitmentsPerBatch:    1,
	}

	err = PopulateGenesisAccounts(s.storage, AssignStateIDs(genesisAccounts))
	s.NoError(err)

	s.transactionExecutor = newTestTransactionExecutor(s.storage, &eth.Client{}, s.cfg, transactionExecutorOpts{})
}

func (s *TransferCommitmentsTestSuite) TearDownTest() {
	err := s.teardown()
	s.NoError(err)
}

func (s *TransferCommitmentsTestSuite) TestCreateTransferCommitments_QueriesForMorePendingTransfersUntilSatisfied() {
	transfers := generateValidTransfers(6)

	for i := 1; i < 6; i++ {
		transfers[i].Amount = models.MakeUint256(99999999999)
	}

	transfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.Transfer,
			FromStateID: 24,
			Amount:      models.MakeUint256(1),
			Fee:         models.MakeUint256(1),
			Nonce:       models.MakeUint256(10),
		},
		ToStateID: 1,
	}
	transfers = append(transfers, transfer)

	s.addTransfers(transfers)

	err := addNewDummyState(s.storage, s.transactionExecutor.stateTree, 24)
	s.NoError(err)

	pendingTransfers, err := s.storage.GetPendingTransfers(s.cfg.PendingTxsCountMultiplier*s.cfg.TxsPerCommitment, nil)
	s.NoError(err)
	s.Len(pendingTransfers, 4)

	preRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	commitments, err := s.transactionExecutor.createTransferCommitments(pendingTransfers, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	postRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *TransferCommitmentsTestSuite) TestCreateTransferCommitments_DoesNothingWhenThereAreNotEnoughPendingTransfers() {
	preRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	commitments, err := s.transactionExecutor.createTransferCommitments([]models.Transfer{}, testDomain)
	s.NoError(err)
	s.Len(commitments, 0)

	postRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *TransferCommitmentsTestSuite) TestCreateTransferCommitments_DoesNothingWhenThereAreNotEnoughValidTransfers() {
	transfers := generateValidTransfers(2)
	transfers[1].Amount = models.MakeUint256(99999999999)
	s.addTransfers(transfers)

	pendingTransfers, err := s.storage.GetPendingTransfers(s.cfg.TxsPerCommitment, nil)
	s.NoError(err)
	s.Len(pendingTransfers, 2)

	preRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	commitments, err := s.transactionExecutor.createTransferCommitments(pendingTransfers, testDomain)
	s.NoError(err)
	s.Len(commitments, 0)

	postRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *TransferCommitmentsTestSuite) TestCreateTransferCommitments_StoresCorrectCommitment() {
	pendingTransfers := s.prepareAndReturnPendingTransfers(3)

	preRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	commitments, err := s.transactionExecutor.createTransferCommitments(pendingTransfers, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, 24)
	s.Equal(commitments[0].FeeReceiver, uint32(2))
	s.Nil(commitments[0].IncludedInBatch)

	postRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *TransferCommitmentsTestSuite) TestCreateTransferCommitments_CreatesMaximallyAsManyCommitmentsAsSpecifiedInConfig() {
	pendingTransfers := s.prepareAndReturnPendingTransfers(2)

	commitments, err := s.transactionExecutor.createTransferCommitments(pendingTransfers, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
}

func (s *TransferCommitmentsTestSuite) TestCreateTransferCommitments_MarksTransfersAsIncludedInCommitment() {
	pendingTransfers := s.prepareAndReturnPendingTransfers(2)

	commitments, err := s.transactionExecutor.createTransferCommitments(pendingTransfers, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	for i := range pendingTransfers {
		tx, err := s.storage.GetTransfer(pendingTransfers[i].Hash)
		s.NoError(err)
		s.Equal(*tx.IncludedInCommitment, int32(1))
	}
}

func (s *TransferCommitmentsTestSuite) TestRemoveTransfer() {
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

func TestTransferCommitmentsTestSuite(t *testing.T) {
	suite.Run(t, new(TransferCommitmentsTestSuite))
}

func (s *TransferCommitmentsTestSuite) addTransfers(transfers []models.Transfer) {
	for i := range transfers {
		err := s.storage.AddTransfer(&transfers[i])
		s.NoError(err)
	}
}

func (s *TransferCommitmentsTestSuite) prepareAndReturnPendingTransfers(transfersAmount uint64) []models.Transfer {
	transfers := generateValidTransfers(transfersAmount)
	s.addTransfers(transfers)

	pendingTransfers, err := s.storage.GetPendingTransfers(transfersAmount, nil)
	s.NoError(err)
	s.Len(pendingTransfers, int(transfersAmount))

	return pendingTransfers
}

func addNewDummyState(storage *storage.Storage, stateTree *storage.StateTree, stateID uint32) error {
	dummyAccount := models.Account{
		PubKeyID: 500,
		PublicKey: models.MakePublicKeyFromInts([4]*big.Int{
			big.NewInt(9),
			big.NewInt(10),
			big.NewInt(3),
			big.NewInt(2),
		}),
	}

	err := storage.AddAccountIfNotExists(&dummyAccount)
	if err != nil {
		return err
	}

	err = stateTree.Set(stateID, &models.UserState{
		PubKeyID:   dummyAccount.PubKeyID,
		TokenIndex: models.MakeUint256(0),
		Balance:    models.MakeUint256(1000),
		Nonce:      models.MakeUint256(10),
	})
	if err != nil {
		return err
	}

	return nil
}
