package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
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
	storage                *st.TestStorage
	cfg                    *config.RollupConfig
	executionCtx           *ExecutionContext
	maxTxBytesInCommitment int
}

func (s *TransferCommitmentsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TransferCommitmentsTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.cfg = &config.RollupConfig{
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    4,
		FeeReceiverPubKeyID:    2,
		MaxCommitmentsPerBatch: 1,
	}
	s.maxTxBytesInCommitment = encoder.TransferLength * int(s.cfg.MaxTxsPerCommitment)

	err = populateAccounts(s.storage.Storage, genesisBalances)
	s.NoError(err)

	s.executionCtx = NewTestExecutionContext(s.storage.Storage, &eth.Client{}, s.cfg)
}

func populateAccounts(storage *st.Storage, balances []models.Uint256) error {
	for i := uint32(0); i < uint32(len(balances)); i++ {
		_, err := storage.StateTree.Set(i, &models.UserState{
			PubKeyID: i,
			TokenID:  models.MakeUint256(0),
			Balance:  balances[i],
			Nonce:    models.MakeUint256(0),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *TransferCommitmentsTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *TransferCommitmentsTestSuite) TestCreateTransferCommitments_WithMinTxsPerCommitment() {
	transfers := generateValidTransfers(1)
	s.addTransfers(transfers)

	preRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.TransferLength * len(transfers)
	commitments, err := s.executionCtx.CreateTransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, expectedTxsLength)

	postRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *TransferCommitmentsTestSuite) TestCreateTransferCommitments_WithMoreThanMinTxsPerCommitment() {
	transfers := generateValidTransfers(3)
	s.addTransfers(transfers)

	preRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.TransferLength * len(transfers)
	commitments, err := s.executionCtx.CreateTransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, expectedTxsLength)

	postRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *TransferCommitmentsTestSuite) TestCreateTransferCommitments_QueriesForMorePendingTransfersUntilSatisfied() {
	addAccountWithHighNonce(s.Assertions, s.storage.Storage, 123)

	transfers := generateValidTransfers(6)
	s.invalidateTransfers(transfers[3:6])

	highNonceTransfer := testutils.MakeTransfer(123, 1, 10, 1)
	transfers = append(transfers, highNonceTransfer)

	s.addTransfers(transfers)

	preRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	commitments, err := s.executionCtx.CreateTransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, s.maxTxBytesInCommitment)

	postRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *TransferCommitmentsTestSuite) TestCreateTransferCommitments_ForMultipleCommitmentsInBatch() {
	s.cfg = &config.RollupConfig{
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    4,
		FeeReceiverPubKeyID:    2,
		MaxCommitmentsPerBatch: 3,
	}

	s.executionCtx = NewTestExecutionContext(s.storage.Storage, &eth.Client{}, s.cfg)

	addAccountWithHighNonce(s.Assertions, s.storage.Storage, 123)

	transfers := generateValidTransfers(9)
	s.invalidateTransfers(transfers[7:9])

	highNonceTransfers := []models.Transfer{
		testutils.MakeTransfer(123, 1, 10, 1),
		testutils.MakeTransfer(123, 1, 11, 1),
	}

	transfers = append(transfers, highNonceTransfers...)
	s.addTransfers(transfers)

	preRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	commitments, err := s.executionCtx.CreateTransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 3)
	s.Len(commitments[0].Transactions, s.maxTxBytesInCommitment)
	s.Len(commitments[1].Transactions, s.maxTxBytesInCommitment)
	s.Len(commitments[2].Transactions, encoder.TransferLength)

	postRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[2].PostStateRoot, *postRoot)
}

func (s *TransferCommitmentsTestSuite) invalidateTransfers(transfers []models.Transfer) {
	for i := range transfers {
		tx := &transfers[i]
		tx.Amount = *genesisBalances[tx.FromStateID].MulN(10)
	}
}

func (s *TransferCommitmentsTestSuite) TestCreateTransferCommitments_ReturnsErrorWhenThereAreNotEnoughPendingTransfers() {
	preRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	commitments, err := s.executionCtx.CreateTransferCommitments(testDomain)
	s.Nil(commitments)
	s.Equal(ErrNotEnoughTransfers, err)

	postRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *TransferCommitmentsTestSuite) TestCreateTransferCommitments_ReturnsErrorWhenThereAreNotEnoughValidTransfers() {
	s.cfg = &config.RollupConfig{
		MinTxsPerCommitment:    32,
		MaxTxsPerCommitment:    32,
		FeeReceiverPubKeyID:    2,
		MaxCommitmentsPerBatch: 1,
	}

	s.executionCtx = NewTestExecutionContext(s.storage.Storage, &eth.Client{}, s.cfg)

	transfers := generateValidTransfers(2)
	s.addTransfers(transfers)

	preRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	commitments, err := s.executionCtx.CreateTransferCommitments(testDomain)
	s.Nil(commitments)
	s.Equal(ErrNotEnoughTransfers, err)

	postRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *TransferCommitmentsTestSuite) TestCreateTransferCommitments_StoresCorrectCommitment() {
	transfersCount := uint32(4)
	s.preparePendingTransfers(transfersCount)

	preRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.TransferLength * int(transfersCount)
	commitments, err := s.executionCtx.CreateTransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].Transactions, expectedTxsLength)
	s.Equal(commitments[0].FeeReceiver, uint32(2))

	postRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].PostStateRoot, *postRoot)
}

func (s *TransferCommitmentsTestSuite) TestCreateTransferCommitments_CreatesMaximallyAsManyCommitmentsAsSpecifiedInConfig() {
	s.preparePendingTransfers(5)

	commitments, err := s.executionCtx.CreateTransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)
}

func (s *TransferCommitmentsTestSuite) TestCreateTransferCommitments_MarksTransfersAsIncludedInCommitment() {
	transfersCount := uint32(4)
	s.preparePendingTransfers(transfersCount)

	pendingTransfers, err := s.storage.GetPendingTransfers(transfersCount)
	s.NoError(err)
	s.Len(pendingTransfers, int(transfersCount))

	commitments, err := s.executionCtx.CreateTransferCommitments(testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	for i := range pendingTransfers {
		tx, err := s.storage.GetTransfer(pendingTransfers[i].Hash)
		s.NoError(err)
		s.Equal(commitments[0].ID, *tx.CommitmentID)
	}
}

func (s *TransferCommitmentsTestSuite) TestRemoveTransfer() {
	transfer1 := createRandomTransferWithHash()
	transfer2 := createRandomTransferWithHash()
	transfer3 := createRandomTransferWithHash()

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
	_, err := storage.StateTree.Set(stateID, &models.UserState{
		PubKeyID: 500,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(10),
	})
	s.NoError(err)
}

func createRandomTransferWithHash() models.Transfer {
	return models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash: utils.RandomHash(),
		},
	}
}
