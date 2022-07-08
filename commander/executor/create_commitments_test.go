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
	s.AcceptNewConfig()

	transfers := s.preparePendingTransfers(1)

	preRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.TransferLength * len(transfers)
	commitments, err := s.txsCtx.CreateCommitments(context.Background())
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].ToTxCommitmentWithTxs().Transactions, expectedTxsLength)

	postRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].ToTxCommitmentWithTxs().PostStateRoot, *postRoot)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_WithMoreThanMinTxsPerCommitment() {
	transfers := s.preparePendingTransfers(3)

	preRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.TransferLength * len(transfers)
	commitments, err := s.txsCtx.CreateCommitments(context.Background())
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].ToTxCommitmentWithTxs().Transactions, expectedTxsLength)

	postRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].ToTxCommitmentWithTxs().PostStateRoot, *postRoot)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_ForMultipleCommitmentsInBatch() {
	s.cfg = &config.RollupConfig{
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    4,
		FeeReceiverPubKeyID:    2,
		MaxCommitmentsPerBatch: 3,
	}
	s.AcceptNewConfig()

	addAccountWithHighNonce(s.Assertions, s.storage.Storage, 123)

	transfers := testutils.GenerateValidTransfers(9)
	s.invalidateTransfers(transfers[7:9])

	highNonceTransfers := []models.Transfer{
		testutils.MakeTransfer(123, 1, 10, 1),
		testutils.MakeTransfer(123, 1, 11, 1),
	}

	transfers = append(transfers, highNonceTransfers...)
	s.initTxs(transfers)

	preRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)

	commitments, err := s.txsCtx.CreateCommitments(context.Background())
	s.NoError(err)
	s.Len(commitments, 3)
	s.Len(commitments[0].ToTxCommitmentWithTxs().Transactions, s.maxTxBytesInCommitment)
	s.Len(commitments[1].ToTxCommitmentWithTxs().Transactions, s.maxTxBytesInCommitment)
	s.Len(commitments[2].ToTxCommitmentWithTxs().Transactions, encoder.TransferLength)

	postRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[2].ToTxCommitmentWithTxs().PostStateRoot, *postRoot)
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

	s.preparePendingTransfers(0)

	commitments, err := s.txsCtx.CreateCommitments(context.Background())
	s.Nil(commitments)
	s.ErrorIs(err, ErrNotEnoughTxs)

	postRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_ReturnsErrorWhenThereAreNotEnoughValidTransfers() {
	s.cfg = &config.RollupConfig{
		MinTxsPerCommitment:    32,
		MaxTxsPerCommitment:    32,
		FeeReceiverPubKeyID:    2,
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 1,
	}
	s.AcceptNewConfig()

	s.preparePendingTransfers(2)

	commitments, err := s.txsCtx.CreateCommitments(context.Background())
	s.Nil(commitments)
	s.ErrorIs(err, ErrNotEnoughTxs)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_StoresCorrectCommitment() {
	transfersCount := uint32(4)
	s.preparePendingTransfers(transfersCount)

	preRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.TransferLength * int(transfersCount)
	commitments, err := s.txsCtx.CreateCommitments(context.Background())
	s.NoError(err)
	s.Len(commitments, 1)
	s.Len(commitments[0].ToTxCommitmentWithTxs().Transactions, expectedTxsLength)
	s.Equal(commitments[0].ToTxCommitmentWithTxs().FeeReceiver, uint32(2))

	postRoot, err := s.txsCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(commitments[0].ToTxCommitmentWithTxs().PostStateRoot, *postRoot)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_CreatesMaximallyAsManyCommitmentsAsSpecifiedInConfig() {
	s.preparePendingTransfers(5)

	commitments, err := s.txsCtx.CreateCommitments(context.Background())
	s.NoError(err)
	s.Len(commitments, 1)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_MarksTransfersAsIncludedInCommitment() {
	transfers := s.preparePendingTransfers(4)

	commitments, err := s.txsCtx.CreateCommitments(context.Background())
	s.NoError(err)
	s.Len(commitments, 1)

	for i := range transfers {
		tx, err := s.storage.GetTransfer(transfers[i].Hash)
		s.NoError(err)
		s.Equal(commitments[0].ToTxCommitmentWithTxs().ID, *tx.CommitmentSlot.CommitmentID())
		s.Nil(tx.ErrorMessage)
	}
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_SkipsNonceTooHighTx() {
	txs := testutils.GenerateValidTransfers(5)
	validTxs := txs[:4]
	nonceTooHighTx := &txs[4]
	nonceTooHighTx.Nonce = models.MakeUint256(21)

	s.initTxs(validTxs.AppendOne(nonceTooHighTx))

	commitments, err := s.txsCtx.CreateCommitments(context.Background())
	s.NoError(err)
	s.Len(commitments, 1)

	for i := range validTxs {
		var tx *models.Transfer
		tx, err = s.storage.GetTransfer(validTxs[i].Hash)
		s.NoError(err)
		s.Equal(commitments[0].ToTxCommitmentWithTxs().ID, *tx.CommitmentSlot.CommitmentID())
	}

	tx, err := s.storage.GetTransfer(nonceTooHighTx.Hash)
	s.NoError(err)
	s.Nil(tx.CommitmentSlot)
	s.Nil(tx.ErrorMessage)
}

func (s *CreateCommitmentsTestSuite) preparePendingTransfers(transfersAmount uint32) models.TransferArray {
	txs := testutils.GenerateValidTransfers(transfersAmount)
	s.initTxs(txs)
	return txs
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_DoesNotCreateCommitmentsWithLessTxsThanRequired() {
	validTransfer := testutils.MakeTransfer(1, 2, 0, 100)
	invalidTransfer := testutils.MakeTransfer(2, 1, 1234, 100)
	s.initTxs(models.TransferArray{validTransfer, invalidTransfer})

	commitments, err := s.txsCtx.CreateCommitments(context.Background())
	s.Nil(commitments)
	s.ErrorIs(err, ErrNotEnoughTxs)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_ReadyTransactionSkipsMinCommitmentsCheck() {
	s.cfg.MinTxsPerCommitment = 1
	s.cfg.MaxTxsPerCommitment = 1
	s.cfg.MinCommitmentsPerBatch = 2
	s.cfg.MaxTxnDelay = 1 * time.Second
	s.AcceptNewConfig()

	validTransfer := testutils.MakeTransfer(1, 2, 0, 100)
	{
		twoSecondsAgo := time.Now().UTC().Add(time.Duration(-2) * time.Second)
		validTransfer.ReceiveTime = models.NewTimestamp(twoSecondsAgo)
	}
	s.initTxs(models.TransferArray{validTransfer})

	batchData, err := s.txsCtx.CreateCommitments(context.Background())
	s.NoError(err)
	s.NotNil(batchData)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_ReturnsErrorIfCouldNotCreateEnoughCommitments() {
	s.cfg.MinTxsPerCommitment = 1
	s.cfg.MaxTxsPerCommitment = 1
	s.cfg.MinCommitmentsPerBatch = 2
	s.AcceptNewConfig()

	validTransfer := testutils.MakeTransfer(1, 2, 0, 100)
	invalidTransfer := testutils.MakeTransfer(2, 1, 1234, 100)
	s.initTxs(models.TransferArray{validTransfer, invalidTransfer})

	commitments, err := s.txsCtx.CreateCommitments(context.Background())
	s.Nil(commitments)
	s.ErrorIs(err, ErrNotEnoughTxs)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_StoresErrorMessagesOfInvalidTransactions() {
	s.cfg.MinTxsPerCommitment = 1
	s.AcceptNewConfig()

	invalidTransfer := testutils.MakeTransfer(1, 1234, 0, 100)
	s.initTxs(models.TransferArray{invalidTransfer})

	commitments, err := s.txsCtx.CreateCommitments(context.Background())
	s.Nil(commitments)
	s.ErrorIs(err, ErrNotEnoughTxs)
	s.Len(s.txsCtx.txErrorsToStore, 1)

	expectedTxError := models.TxError{
		TxHash:        invalidTransfer.Hash,
		SenderStateID: invalidTransfer.FromStateID,
		ErrorMessage:  applier.ErrNonexistentReceiver.Error(),
	}
	s.Equal(expectedTxError, s.txsCtx.txErrorsToStore[0])
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_DoesNotCallRevertToWhenNotNecessary() {
	validTransfer := testutils.MakeTransfer(1, 2, 0, 100)
	invalidTransfer := testutils.MakeTransfer(2, 1, 1234, 100)
	s.initTxs(models.TransferArray{validTransfer, invalidTransfer})

	preStateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	commitments, err := s.txsCtx.CreateCommitments(context.Background())
	s.Nil(commitments)
	s.ErrorIs(err, ErrNotEnoughTxs)

	postStateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	s.NotEqual(preStateRoot, postStateRoot)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_CallsRevertToWhenNecessary() {
	validTransfers := models.TransferArray{
		testutils.MakeTransfer(1, 2, 0, 100),
		testutils.MakeTransfer(1, 2, 1, 100),
		testutils.MakeTransfer(1, 2, 2, 100),
	}
	invalidTransfer := testutils.MakeTransfer(2, 1, 1234, 100)

	// Calculate state root after applying 2 valid transfers
	s.cfg.MinTxsPerCommitment = 2
	s.cfg.MinCommitmentsPerBatch = 1
	s.AcceptNewConfig()

	tempTxsCtx := NewTxsContext(
		s.txsCtx.storage,
		s.txsCtx.client,
		s.cfg,
		metrics.NewCommanderMetrics(),
		s.txsCtx.Mempool,
		context.Background(),
		batchtype.Transfer,
	)
	initTxs(s.Assertions, tempTxsCtx, validTransfers[:2])

	commitments, err := tempTxsCtx.CreateCommitments(context.Background())
	s.NoError(err)
	s.Len(commitments, 1)

	expectedPostStateRoot, err := tempTxsCtx.storage.StateTree.Root()
	s.NoError(err)

	tempTxsCtx.Rollback(nil)

	// Do the test
	s.cfg.MinTxsPerCommitment = 2
	s.cfg.MaxTxsPerCommitment = 2
	s.cfg.MinCommitmentsPerBatch = 1
	s.AcceptNewConfig()

	s.initTxs(validTransfers.AppendOne(&invalidTransfer))

	commitments, err = s.txsCtx.CreateCommitments(context.Background())
	s.NoError(err)
	s.Len(commitments, 1)

	stateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	s.Equal(expectedPostStateRoot, stateRoot)
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_SupportsTransactionReplacement() {
	// Mine the transaction with higher fee in case there are two txs from the same sender with the same nonce
	s.cfg.MinTxsPerCommitment = 1
	s.AcceptNewConfig()

	transfer := testutils.MakeTransfer(1, 2, 0, 100)
	higherFeeTransfer := testutils.MakeTransfer(1, 2, 0, 100)
	higherFeeTransfer.Fee = *transfer.Fee.MulN(2)

	s.initTxs(models.TransferArray{transfer, higherFeeTransfer})

	_, err := s.txsCtx.CreateCommitments(context.Background())
	s.NoError(err)

	minedHigherFeeTransfer, err := s.storage.GetTransfer(higherFeeTransfer.Hash)
	s.NoError(err)

	expectedCommitmentID := models.CommitmentID{
		BatchID:      models.MakeUint256(1),
		IndexInBatch: 0,
	}
	s.Equal(expectedCommitmentID, *minedHigherFeeTransfer.CommitmentSlot.CommitmentID())
}

func (s *CreateCommitmentsTestSuite) initTxs(txs models.GenericTransactionArray) {
	initTxs(s.Assertions, s.txsCtx, txs)
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
