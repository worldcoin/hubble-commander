package executor

import (
	"context"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/ethereum/go-ethereum/common"
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

type CreateCommitmentsTestSuite struct {
	testSuiteWithTxsContext
	wallets                []bls.Wallet
	maxTxBytesInCommitment int
}

func (s *CreateCommitmentsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *CreateCommitmentsTestSuite) SetupTest() {
	s.testSuiteWithTxsContext.SetupTestWithConfig(batchtype.Transfer, &config.RollupConfig{
		MinTxsPerCommitment:    2,
		MaxTxsPerCommitment:    4,
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		FeeReceiverPubKeyID:    2,
	})
	s.maxTxBytesInCommitment = encoder.TransferLength * int(s.cfg.MaxTxsPerCommitment)

	domain, err := s.client.GetDomain()
	s.NoError(err)
	s.wallets = testutils.GenerateWallets(s.Assertions, domain, 2)

	s.addUserStates()

	err = populateAccounts(s.storage.Storage, genesisBalances)
	s.NoError(err)
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

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_WithMinTxsPerCommitment() {
	s.cfg.MinTxsPerCommitment = 1
	s.cfg.MinCommitmentsPerBatch = 1

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

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_WithMoreThanMinTxsPerCommitment() {
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

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_ForMultipleCommitmentsInBatch() {
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

func (s *CreateCommitmentsTestSuite) invalidateTransfers(transfers []models.Transfer) {
	for i := range transfers {
		tx := &transfers[i]
		tx.Amount = *genesisBalances[tx.FromStateID].MulN(10)
	}
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_ReturnsErrorWhenThereAreNotEnoughPendingTransfers() {
	preRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)

	batchData, err := s.txsCtx.CreateCommitments()
	s.Nil(batchData)
	s.ErrorIs(err, ErrNotEnoughTxs)

	postRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_ReturnsErrorWhenThereAreNotEnoughValidTransfers() {
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

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_StoresCorrectCommitment() {
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

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_CreatesMaximallyAsManyCommitmentsAsSpecifiedInConfig() {
	s.preparePendingTransfers(5)

	batchData, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(batchData.Commitments(), 1)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_MarksTransfersAsIncludedInCommitment() {
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

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_SkipsNonceTooHighTx() {
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

func (s *CreateCommitmentsTestSuite) preparePendingTransfers(transfersAmount uint32) {
	transfers := testutils.GenerateValidTransfers(transfersAmount)
	s.addTransfers(transfers)
}

func (s *CreateCommitmentsTestSuite) addTransfers(transfers []models.Transfer) {
	err := s.storage.BatchAddTransfer(transfers)
	s.NoError(err)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_DoesNotCreateCommitmentsWithLessTxsThanRequired() {
	validTransfer := testutils.MakeTransfer(1, 2, 0, 100)
	s.hashSignAndAddTransfer(&s.wallets[0], &validTransfer)
	invalidTransfer := testutils.MakeTransfer(2, 1, 1234, 100)
	s.hashSignAndAddTransfer(&s.wallets[1], &invalidTransfer)

	batchData, err := s.txsCtx.CreateCommitments()
	s.Nil(batchData)
	s.ErrorIs(err, ErrNotEnoughCommitments)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_ReadyTransactionSkipsMinCommitmentsCheck() {
	s.cfg.MinTxsPerCommitment = 1
	s.cfg.MaxTxsPerCommitment = 1
	s.cfg.MinCommitmentsPerBatch = 2
	s.cfg.MaxTxnDelay = 1 * time.Second

	validTransfer := testutils.MakeTransfer(1, 2, 0, 100)
	{
		twoSecondsAgo := time.Now().UTC().Add(time.Duration(-2) * time.Second)
		validTransfer.ReceiveTime = models.NewTimestamp(twoSecondsAgo)
	}
	s.hashSignAndAddTransfer(&s.wallets[0], &validTransfer)

	batchData, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.NotNil(batchData)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_ReturnsErrorIfCouldNotCreateEnoughCommitments() {
	s.cfg.MinTxsPerCommitment = 1
	s.cfg.MaxTxsPerCommitment = 1
	s.cfg.MinCommitmentsPerBatch = 2

	validTransfer := testutils.MakeTransfer(1, 2, 0, 100)
	s.hashSignAndAddTransfer(&s.wallets[0], &validTransfer)
	invalidTransfer := testutils.MakeTransfer(2, 1, 1234, 100)
	s.hashSignAndAddTransfer(&s.wallets[1], &invalidTransfer)

	batchData, err := s.txsCtx.CreateCommitments()
	s.Nil(batchData)
	s.ErrorIs(err, ErrNotEnoughCommitments)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_StoresErrorMessagesOfInvalidTransactions() {
	s.cfg.MinTxsPerCommitment = 1

	invalidTransfer := testutils.MakeTransfer(1, 1234, 0, 100)
	s.hashSignAndAddTransfer(&s.wallets[0], &invalidTransfer)

	batchData, err := s.txsCtx.CreateCommitments()
	s.Nil(batchData)
	s.ErrorIs(err, ErrNotEnoughCommitments)

	s.Len(s.txsCtx.txErrorsToStore, 1)
	s.Equal(invalidTransfer.Hash, s.txsCtx.txErrorsToStore[0].TxHash)
	s.Equal(applier.ErrNonexistentReceiver.Error(), s.txsCtx.txErrorsToStore[0].ErrorMessage)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_DoesNotCallRevertToWhenNotNecessary() {
	validTransfer := testutils.MakeTransfer(1, 2, 0, 100)
	s.hashSignAndAddTransfer(&s.wallets[0], &validTransfer)
	invalidTransfer := testutils.MakeTransfer(2, 1, 1234, 100)
	s.hashSignAndAddTransfer(&s.wallets[1], &invalidTransfer)

	preStateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	batchData, err := s.txsCtx.CreateCommitments()
	s.Nil(batchData)
	s.ErrorIs(err, ErrNotEnoughCommitments)

	postStateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	s.NotEqual(preStateRoot, postStateRoot)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_CallsRevertToWhenNecessary() {
	validTransfers := []models.Transfer{
		testutils.MakeTransfer(1, 2, 0, 100),
		testutils.MakeTransfer(1, 2, 1, 100),
		testutils.MakeTransfer(1, 2, 2, 100),
	}
	invalidTransfer := testutils.MakeTransfer(2, 1, 1234, 100)

	// Calculate state root after applying 2 valid transfers
	s.cfg.MinTxsPerCommitment = 2
	s.cfg.MinCommitmentsPerBatch = 1

	s.hashSignAndAddTransfer(&s.wallets[0], &validTransfers[0])
	s.hashSignAndAddTransfer(&s.wallets[0], &validTransfers[1])

	tempTxsCtx := NewTxsContext(
		s.txsCtx.storage,
		s.txsCtx.client,
		s.cfg,
		metrics.NewCommanderMetrics(),
		context.Background(),
		batchtype.Transfer,
	)
	batchData, err := tempTxsCtx.CreateCommitments()
	s.NoError(err)
	s.Equal(batchData.Len(), 1)

	expectedPostStateRoot, err := tempTxsCtx.storage.StateTree.Root()
	s.NoError(err)

	tempTxsCtx.Rollback(nil)

	// Do the test
	s.cfg.MinTxsPerCommitment = 2
	s.cfg.MaxTxsPerCommitment = 2
	s.cfg.MinCommitmentsPerBatch = 1

	s.hashSignAndAddTransfer(&s.wallets[0], &validTransfers[2])
	s.hashSignAndAddTransfer(&s.wallets[1], &invalidTransfer)

	batchData, err = s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Equal(batchData.Len(), 1)

	stateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	s.Equal(expectedPostStateRoot, stateRoot)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_SupportsTransactionReplacement() {
	// Mine the transaction with higher fee in case there are two txs from the same sender with the same nonce
	s.cfg.MinTxsPerCommitment = 1

	transfer := testutils.MakeTransfer(1, 2, 0, 100)
	transfer.Hash = common.BytesToHash([]byte{1})
	err := s.storage.AddTransfer(&transfer)
	s.NoError(err)

	higherFeeTransfer := transfer
	higherFeeTransfer.Hash = common.BytesToHash([]byte{2})
	higherFeeTransfer.Fee = *transfer.Fee.MulN(2)
	err = s.storage.AddTransfer(&higherFeeTransfer)
	s.NoError(err)

	s.Less(transfer.Hash.String(), higherFeeTransfer.Hash.String())

	_, err = s.txsCtx.CreateCommitments()
	s.NoError(err)

	s.Len(s.txsCtx.txErrorsToStore, 1)
	txErr := s.txsCtx.txErrorsToStore[0]
	s.Equal(transfer.Hash, txErr.TxHash)
	s.Equal(applier.ErrNonceTooLow.Error(), txErr.ErrorMessage)

	minedHigherFeeTransfer, err := s.storage.GetTransfer(higherFeeTransfer.Hash)
	s.NoError(err)

	expectedCommitmentID := models.CommitmentID{
		BatchID:      models.MakeUint256(1),
		IndexInBatch: 0,
	}
	s.Equal(expectedCommitmentID, *minedHigherFeeTransfer.CommitmentID)
}

func (s *CreateCommitmentsTestSuite) hashSignAndAddTransfer(wallet *bls.Wallet, transfer *models.Transfer) {
	hash, err := encoder.HashTransfer(transfer)
	s.NoError(err)
	transfer.Hash = *hash

	encodedTransfer, err := encoder.EncodeTransferForSigning(transfer)
	s.NoError(err)
	signature, err := wallet.Sign(encodedTransfer)
	s.NoError(err)
	transfer.Signature = *signature.ModelsSignature()

	err = s.storage.AddTransfer(transfer)
	s.NoError(err)
}

func (s *CreateCommitmentsTestSuite) addUserStates() {
	feeReceiver := &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	}
	_, err := s.storage.StateTree.Set(0, feeReceiver)
	s.NoError(err)

	_, err = s.storage.StateTree.Set(1, &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)

	_, err = s.storage.StateTree.Set(2, &models.UserState{
		PubKeyID: 2,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	})
	s.NoError(err)
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

func TestCreateCommitmentsTestSuite(t *testing.T) {
	suite.Run(t, new(CreateCommitmentsTestSuite))
}
