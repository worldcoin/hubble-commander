package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	genesisBalances = []models.Uint256{
		models.MakeUint256(1000),
		models.MakeUint256(1000),
		models.MakeUint256(1000),
	}
	testDomain = &bls.Domain{1, 2, 3, 4}
)

type TransferCommitmentsTestSuite struct {
	*require.Assertions
	suite.Suite
	teardown            func() error
	storage             *st.Storage
	cfg                 *config.RollupConfig
	transactionExecutor *TransactionExecutor
}

func (s *TransferCommitmentsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TransferCommitmentsTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.teardown = testStorage.Teardown
	s.cfg = &config.RollupConfig{
		TxsPerCommitment:       2,
		FeeReceiverPubKeyID:    2,
		MaxCommitmentsPerBatch: 1,
	}

	err = populateAccounts(s.storage, genesisBalances)
	s.NoError(err)

	s.transactionExecutor = NewTestTransactionExecutor(s.storage, &eth.Client{}, s.cfg, TransactionExecutorOpts{})
}

func populateAccounts(storage *st.Storage, balances []models.Uint256) error {
	stateTree := st.NewStateTree(storage)
	for i := uint32(0); i < uint32(len(balances)); i++ {
		err := storage.AddAccountIfNotExists(&models.Account{
			PubKeyID:  i,
			PublicKey: models.PublicKey{},
		})
		if err != nil {
			return err
		}

		err = stateTree.Set(i, &models.UserState{
			PubKeyID:   i,
			TokenIndex: models.MakeUint256(0),
			Balance:    balances[i],
			Nonce:      models.MakeUint256(0),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *TransferCommitmentsTestSuite) TearDownTest() {
	err := s.teardown()
	s.NoError(err)
}

func (s *TransferCommitmentsTestSuite) TestCreateTransferCommitments_QueriesForMorePendingTransfersUntilSatisfied() {
	addAccountWithHighNonce(s.Assertions, s.storage, 123)

	transfers := generateValidTransfers(6)
	s.invalidateTransfers(transfers[1:6])

	highNonceTransfer := models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.Transfer,
			FromStateID: 123,
			Amount:      models.MakeUint256(1),
			Fee:         models.MakeUint256(1),
			Nonce:       models.MakeUint256(10),
		},
		ToStateID: 1,
	}
	transfers = append(transfers, highNonceTransfer)

	s.addTransfers(transfers)

	pendingTransfers, err := s.storage.GetPendingTransfers(pendingTxsCountMultiplier * s.cfg.TxsPerCommitment)
	s.NoError(err)
	s.Len(pendingTransfers, 4)

	preRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	commitments, err := s.transactionExecutor.CreateTransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	postRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *TransferCommitmentsTestSuite) invalidateTransfers(transfers []models.Transfer) {
	for i := range transfers {
		tx := &transfers[i]
		tx.Amount = *genesisBalances[tx.FromStateID].MulN(10)
	}
}

func (s *TransferCommitmentsTestSuite) TestCreateTransferCommitments_ReturnsErrorWhenThereAreNotEnoughPendingTransfers() {
	preRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	commitments, err := s.transactionExecutor.CreateTransferCommitments(testDomain)
	s.Nil(commitments)
	s.Equal(ErrNotEnoughTransfers, err)

	postRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *TransferCommitmentsTestSuite) TestCreateTransferCommitments_ReturnsErrorWhenThereAreNotEnoughValidTransfers() {
	transfers := generateValidTransfers(2)
	transfers[1].Amount = models.MakeUint256(99999999999)
	s.addTransfers(transfers)

	preRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	commitments, err := s.transactionExecutor.CreateTransferCommitments(testDomain)
	s.Nil(commitments)
	s.Equal(ErrNotEnoughTransfers, err)

	postRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *TransferCommitmentsTestSuite) TestCreateTransferCommitments_StoresCorrectCommitment() {
	s.preparePendingTransfers(3)

	preRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	commitments, err := s.transactionExecutor.CreateTransferCommitments(testDomain)
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
	s.preparePendingTransfers(2)

	commitments, err := s.transactionExecutor.CreateTransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
}

func (s *TransferCommitmentsTestSuite) TestCreateTransferCommitments_MarksTransfersAsIncludedInCommitment() {
	s.preparePendingTransfers(2)

	pendingTransfers, err := s.storage.GetPendingTransfers(2)
	s.NoError(err)
	s.Len(pendingTransfers, 2)

	commitments, err := s.transactionExecutor.CreateTransferCommitments(testDomain)
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

	s.Equal([]models.Transfer{transfer1, transfer3}, removeTransfers(transfers, toRemove))
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

func (s *TransferCommitmentsTestSuite) preparePendingTransfers(transfersAmount uint32) {
	transfers := generateValidTransfers(transfersAmount)
	s.addTransfers(transfers)
}

func addAccountWithHighNonce(s *require.Assertions, storage *st.Storage, stateID uint32) {
	dummyAccount := models.Account{
		PubKeyID:  500,
		PublicKey: models.PublicKey{1, 2, 3, 4},
	}

	err := storage.AddAccountIfNotExists(&dummyAccount)
	s.NoError(err)

	stateTree := st.NewStateTree(storage)
	err = stateTree.Set(stateID, &models.UserState{
		PubKeyID:   dummyAccount.PubKeyID,
		TokenIndex: models.MakeUint256(0),
		Balance:    models.MakeUint256(1000),
		Nonce:      models.MakeUint256(10),
	})
	s.NoError(err)
}
