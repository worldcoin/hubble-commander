package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	genesisBalances = []models.Uint256{
		models.MakeUint256(1000),
		models.MakeUint256(1000),
		models.MakeUint256(1000),
	}
)

type TransferCommitmentsTestSuite struct {
	*require.Assertions
	suite.Suite
	storage                *st.TestStorage
	cfg                    *config.RollupConfig
	txsCtx                 *TxsContext
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
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 1,
	}
	s.maxTxBytesInCommitment = encoder.TransferLength * int(s.cfg.MaxTxsPerCommitment)

	err = populateAccounts(s.storage.Storage, genesisBalances)
	s.NoError(err)

	executionCtx := NewTestExecutionContext(s.storage.Storage, eth.DomainOnlyTestClient, s.cfg)
	s.txsCtx = NewTestTxsContext(executionCtx, batchtype.Transfer)
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

func (s *TransferCommitmentsTestSuite) TestCreateCommitments_WithMinTxsPerCommitment() {
	transfers := testutils.GenerateValidTransfers(1)
	s.addTransfers(transfers)

	preRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.TransferLength * len(transfers)
	batchData, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(batchData.Commitments(), 1)
	s.Len(batchData.Commitments()[0].Transactions, expectedTxsLength)

	postRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(batchData.Commitments()[0].PostStateRoot, *postRoot)
}

func (s *TransferCommitmentsTestSuite) TestCreateCommitments_WithMoreThanMinTxsPerCommitment() {
	transfers := testutils.GenerateValidTransfers(3)
	s.addTransfers(transfers)

	preRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.TransferLength * len(transfers)
	batchData, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(batchData.Commitments(), 1)
	s.Len(batchData.Commitments()[0].Transactions, expectedTxsLength)

	postRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(batchData.Commitments()[0].PostStateRoot, *postRoot)
}

func (s *TransferCommitmentsTestSuite) TestCreateCommitments_ForMultipleCommitmentsInBatch() {
	s.txsCtx.cfg = &config.RollupConfig{
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    4,
		FeeReceiverPubKeyID:    2,
		MaxCommitmentsPerBatch: 3,
	}

	addAccountWithHighNonce(s.Assertions, s.storage.Storage, 123)

	transfers := testutils.GenerateValidTransfers(9)
	s.invalidateTransfers(transfers[7:9])

	highNonceTransfers := []models.Transfer{
		testutils.MakeTransfer(123, 1, 10, 1),
		testutils.MakeTransfer(123, 1, 11, 1),
	}

	transfers = append(transfers, highNonceTransfers...)
	s.addTransfers(transfers)

	preRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)

	batchData, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(batchData.Commitments(), 3)
	s.Len(batchData.Commitments()[0].Transactions, s.maxTxBytesInCommitment)
	s.Len(batchData.Commitments()[1].Transactions, s.maxTxBytesInCommitment)
	s.Len(batchData.Commitments()[2].Transactions, encoder.TransferLength)

	postRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(batchData.Commitments()[2].PostStateRoot, *postRoot)
}

func (s *TransferCommitmentsTestSuite) invalidateTransfers(transfers []models.Transfer) {
	for i := range transfers {
		tx := &transfers[i]
		tx.Amount = *genesisBalances[tx.FromStateID].MulN(10)
	}
}

func (s *TransferCommitmentsTestSuite) TestCreateCommitments_ReturnsErrorWhenThereAreNotEnoughPendingTransfers() {
	preRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)

	batchData, err := s.txsCtx.CreateCommitments()
	s.Nil(batchData)
	s.ErrorIs(err, ErrNotEnoughTxs)

	postRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *TransferCommitmentsTestSuite) TestCreateCommitments_ReturnsErrorWhenThereAreNotEnoughValidTransfers() {
	s.txsCtx.cfg = &config.RollupConfig{
		MinTxsPerCommitment:    32,
		MaxTxsPerCommitment:    32,
		FeeReceiverPubKeyID:    2,
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 1,
	}

	transfers := testutils.GenerateValidTransfers(2)
	s.addTransfers(transfers)

	preRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)

	batchData, err := s.txsCtx.CreateCommitments()
	s.Nil(batchData)
	s.ErrorIs(err, ErrNotEnoughTxs)

	postRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *TransferCommitmentsTestSuite) TestCreateCommitments_StoresCorrectCommitment() {
	transfersCount := uint32(4)
	s.preparePendingTransfers(transfersCount)

	preRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.TransferLength * int(transfersCount)
	batchData, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(batchData.Commitments(), 1)
	s.Len(batchData.Commitments()[0].Transactions, expectedTxsLength)
	s.Equal(batchData.Commitments()[0].FeeReceiver, uint32(2))

	postRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(batchData.Commitments()[0].PostStateRoot, *postRoot)
}

func (s *TransferCommitmentsTestSuite) TestCreateCommitments_CreatesMaximallyAsManyCommitmentsAsSpecifiedInConfig() {
	s.preparePendingTransfers(5)

	batchData, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(batchData.Commitments(), 1)
}

func (s *TransferCommitmentsTestSuite) TestCreateCommitments_MarksTransfersAsIncludedInCommitment() {
	transfers := testutils.GenerateValidTransfers(4)
	s.addTransfers(transfers)

	batchData, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(batchData.Commitments(), 1)
	commitment := &batchData.Commitments()[0]

	for i := range transfers {
		tx, err := s.storage.GetTransfer(transfers[i].Hash)
		s.NoError(err)
		s.Equal(commitment.ID, *tx.CommitmentID)
		s.Nil(tx.ErrorMessage)
	}
}

func (s *TransferCommitmentsTestSuite) TestCreateCommitments_SkipsNonceTooHighTx() {
	txs := testutils.GenerateValidTransfers(5)
	validTxs := txs[:4]
	nonceTooHighTx := &txs[4]
	nonceTooHighTx.Nonce = models.MakeUint256(21)

	s.addTransfers(validTxs)
	err := s.storage.AddTransfer(nonceTooHighTx)
	s.NoError(err)

	batchData, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(batchData.Commitments(), 1)
	commitment := &batchData.Commitments()[0]

	for i := range validTxs {
		var tx *models.Transfer
		tx, err = s.storage.GetTransfer(validTxs[i].Hash)
		s.NoError(err)
		s.Equal(commitment.ID, *tx.CommitmentID)
	}

	tx, err := s.storage.GetTransfer(nonceTooHighTx.Hash)
	s.NoError(err)
	s.Nil(tx.CommitmentID)
	s.Nil(tx.ErrorMessage)
}

func TestTransferCommitmentsTestSuite(t *testing.T) {
	suite.Run(t, new(TransferCommitmentsTestSuite))
}

func (s *TransferCommitmentsTestSuite) addTransfers(transfers []models.Transfer) {
	err := s.storage.BatchAddTransfer(transfers)
	s.NoError(err)
}

func (s *TransferCommitmentsTestSuite) preparePendingTransfers(transfersAmount uint32) {
	transfers := testutils.GenerateValidTransfers(transfersAmount)
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
