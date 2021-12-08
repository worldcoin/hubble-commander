package executor

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
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
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

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

func (s *CreateCommitmentsTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
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
	transfersCount := uint32(4)
	s.preparePendingTransfers(transfersCount)

	pendingTransfers, err := s.storage.GetPendingTransfers()
	s.NoError(err)
	s.Len(pendingTransfers, int(transfersCount))

	batchData, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(batchData.Commitments(), 1)

	for i := range pendingTransfers {
		tx, err := s.storage.GetTransfer(pendingTransfers[i].Hash)
		s.NoError(err)
		s.Equal(batchData.Commitments()[0].ID, *tx.CommitmentID)
	}
}

func (s *CreateCommitmentsTestSuite) TestCreateCommitments_SkipsNonceTooHighTx() {
	validTransfersCount := 4
	s.preparePendingTransfers(uint32(validTransfersCount))

	nonceTooHighTx := testutils.GenerateValidTransfers(1)[0]
	nonceTooHighTx.Nonce = models.MakeUint256(21)
	err := s.storage.AddTransfer(&nonceTooHighTx)
	s.NoError(err)

	pendingTransfers, err := s.storage.GetPendingTransfers()
	s.NoError(err)
	s.Len(pendingTransfers, validTransfersCount+1)

	batchData, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(batchData.Commitments(), 1)

	for i := 0; i < validTransfersCount; i++ {
		var tx *models.Transfer
		tx, err = s.storage.GetTransfer(pendingTransfers[i].Hash)
		s.NoError(err)
		s.Equal(batchData.Commitments()[0].ID, *tx.CommitmentID)
	}

	tx, err := s.storage.GetTransfer(nonceTooHighTx.Hash)
	s.NoError(err)
	s.Nil(tx.CommitmentID)

	pendingTransfers, err = s.storage.GetPendingTransfers()
	s.NoError(err)
	s.Len(pendingTransfers, 1)
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

	preStateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	batchData, err := s.txsCtx.CreateCommitments()
	s.Nil(batchData)
	s.ErrorIs(err, ErrNotEnoughCommitments)

	postStateRoot, err := s.storage.StateTree.Root()
	s.NoError(err)

	s.Equal(preStateRoot, postStateRoot)
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
