package executor

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Create2TransferCommitmentsTestSuite struct {
	*require.Assertions
	suite.Suite
	teardown               func() error
	storage                *st.Storage
	client                 *eth.TestClient
	cfg                    *config.RollupConfig
	transactionExecutor    *TransactionExecutor
	maxTxBytesInCommitment int
}

func (s *Create2TransferCommitmentsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *Create2TransferCommitmentsTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.client, err = eth.NewTestClient()
	s.NoError(err)
	s.storage = testStorage.Storage
	s.teardown = testStorage.Teardown
	s.cfg = &config.RollupConfig{
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    4,
		FeeReceiverPubKeyID:    2,
		MaxCommitmentsPerBatch: 1,
	}
	s.maxTxBytesInCommitment = encoder.Create2TransferLength * int(s.cfg.MaxTxsPerCommitment)

	err = populateAccounts(s.storage, genesisBalances)
	s.NoError(err)

	s.transactionExecutor = NewTestTransactionExecutor(s.storage, s.client.Client, s.cfg, context.Background())
}

func (s *Create2TransferCommitmentsTestSuite) TearDownTest() {
	s.client.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_WithMinTxsPerCommitment() {
	transfers := generateValidCreate2Transfers(1)
	s.addCreate2Transfers(transfers)

	preRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.Create2TransferLength * len(transfers)
	commitments, err := s.transactionExecutor.CreateCreate2TransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, expectedTxsLength)

	postRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_WithMoreThanMinTxsPerCommitment() {
	transfers := generateValidCreate2Transfers(3)
	s.addCreate2Transfers(transfers)

	preRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.Create2TransferLength * len(transfers)
	commitments, err := s.transactionExecutor.CreateCreate2TransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, expectedTxsLength)

	postRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_QueriesForMorePendingTransfersUntilSatisfied() {
	addAccountWithHighNonce(s.Assertions, s.storage, 123)

	transfers := generateValidCreate2Transfers(6)
	s.invalidateCreate2Transfers(transfers[3:6])

	highNonceTransfer := createC2T(123, nil, 10, 1, &models.PublicKey{5, 4, 3, 2, 1})
	transfers = append(transfers, highNonceTransfer)

	s.addCreate2Transfers(transfers)

	preRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	commitments, err := s.transactionExecutor.CreateCreate2TransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, s.maxTxBytesInCommitment)

	postRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_ForMultipleCommitmentsInBatch() {
	s.cfg = &config.RollupConfig{
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    4,
		FeeReceiverPubKeyID:    2,
		MaxCommitmentsPerBatch: 3,
	}

	// TODO-MIN Why the client is different than the one in transfers tests
	s.transactionExecutor = NewTestTransactionExecutor(s.storage, s.client.Client, s.cfg, context.Background())

	addAccountWithHighNonce(s.Assertions, s.storage, 123)

	transfers := generateValidCreate2Transfers(9)
	s.invalidateCreate2Transfers(transfers[7:9])

	highNonceTransfers := []models.Create2Transfer{
		createC2T(123, nil, 10, 1, &models.PublicKey{5, 4, 3, 2, 1}),
		createC2T(123, nil, 11, 1, &models.PublicKey{5, 4, 3, 2, 1}),
	}

	transfers = append(transfers, highNonceTransfers...)
	s.addCreate2Transfers(transfers)

	preRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	commitments, err := s.transactionExecutor.CreateCreate2TransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 3)
	s.Len(commitments[0].Transactions, s.maxTxBytesInCommitment)
	s.Len(commitments[1].Transactions, s.maxTxBytesInCommitment)
	s.Len(commitments[2].Transactions, encoder.Create2TransferLength)

	postRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[2].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) invalidateCreate2Transfers(transfers []models.Create2Transfer) {
	for i := range transfers {
		tx := &transfers[i]
		tx.Amount = *genesisBalances[tx.FromStateID].MulN(10)
	}
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_ReturnsErrorWhenThereAreNotEnoughPendingTransfers() {
	preRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	commitments, err := s.transactionExecutor.CreateCreate2TransferCommitments(testDomain)
	s.Nil(commitments)
	s.Equal(ErrNotEnoughC2Transfers, err)

	postRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_ReturnsErrorWhenThereAreNotEnoughValidTransfers() {
	s.cfg = &config.RollupConfig{
		MinTxsPerCommitment:    32,
		MaxTxsPerCommitment:    32,
		FeeReceiverPubKeyID:    2,
		MaxCommitmentsPerBatch: 1,
	}

	// TODO-MIN validate the client used here. Check TestCreateCreate2TransferCommitments_ForMultipleCommitmentsInBatch
	s.transactionExecutor = NewTestTransactionExecutor(s.storage, &eth.Client{}, s.cfg, context.Background())

	transfers := generateValidCreate2Transfers(2)
	s.addCreate2Transfers(transfers)

	pendingTransfers, err := s.storage.GetPendingCreate2Transfers(32)
	s.NoError(err)
	s.Len(pendingTransfers, 2)

	preRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	commitments, err := s.transactionExecutor.CreateCreate2TransferCommitments(testDomain)
	s.Nil(commitments)
	s.Equal(ErrNotEnoughC2Transfers, err)

	postRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_StoresCorrectCommitment() {
	transfersCount := uint32(4)
	s.preparePendingCreate2Transfers(transfersCount)

	preRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.Create2TransferLength * int(transfersCount)
	commitments, err := s.transactionExecutor.CreateCreate2TransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, expectedTxsLength)
	s.Equal(commitments[0].FeeReceiver, uint32(2))
	s.Nil(commitments[0].IncludedInBatch)

	postRoot, err := s.transactionExecutor.stateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_CreatesMaximallyAsManyCommitmentsAsSpecifiedInConfig() {
	s.preparePendingCreate2Transfers(5)

	commitments, err := s.transactionExecutor.CreateCreate2TransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCreate2TransferCommitments_UpdateTransfers() {
	s.preparePendingCreate2Transfers(2)

	pendingTransfers, err := s.storage.GetPendingCreate2Transfers(2)
	s.NoError(err)
	s.Len(pendingTransfers, 2)

	commitments, err := s.transactionExecutor.CreateCreate2TransferCommitments(testDomain)
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

	s.Equal([]models.Create2Transfer{transfer1, transfer3}, removeC2Ts(transfers, toRemove))
}

func TestCreate2TransferCommitmentsTestSuite(t *testing.T) {
	suite.Run(t, new(Create2TransferCommitmentsTestSuite))
}

func (s *Create2TransferCommitmentsTestSuite) addCreate2Transfers(transfers []models.Create2Transfer) {
	for i := range transfers {
		_, err := s.storage.AddCreate2Transfer(&transfers[i])
		s.NoError(err)
	}
}

func (s *Create2TransferCommitmentsTestSuite) preparePendingCreate2Transfers(transfersAmount uint32) {
	transfers := generateValidCreate2Transfers(transfersAmount)
	s.addCreate2Transfers(transfers)
}
