package executor

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/suite"
)

type Create2TransferCommitmentsTestSuite struct {
	testSuiteWithTxsContext
	maxTxBytesInCommitment int
}

func (s *Create2TransferCommitmentsTestSuite) SetupTest() {
	s.testSuiteWithTxsContext.SetupTestWithConfig(batchtype.Create2Transfer, &config.RollupConfig{
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    4,
		FeeReceiverPubKeyID:    2,
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 1,
	})
	s.maxTxBytesInCommitment = encoder.Create2TransferLength * int(s.cfg.MaxTxsPerCommitment)

	err := populateAccounts(s.storage.Storage, genesisBalances)
	s.NoError(err)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCommitments_WithMinTxsPerCommitment() {
	transfers := testutils.GenerateValidCreate2Transfers(1)
	s.addCreate2Transfers(transfers)

	preRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.Create2TransferLength * len(transfers)
	result, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(result.Commitments(), 1)
	s.Len(result.Commitments()[0].Transactions, expectedTxsLength)

	postRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(result.Commitments()[0].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCommitments_WithMoreThanMinTxsPerCommitment() {
	transfers := testutils.GenerateValidCreate2Transfers(3)
	s.addCreate2Transfers(transfers)

	preRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.Create2TransferLength * len(transfers)
	result, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(result.Commitments(), 1)
	s.Len(result.Commitments()[0].Transactions, expectedTxsLength)

	postRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(result.Commitments()[0].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCommitments_ForMultipleCommitmentsInBatch() {
	s.cfg = &config.RollupConfig{
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    4,
		FeeReceiverPubKeyID:    2,
		MaxCommitmentsPerBatch: 3,
	}

	executionCtx := NewTestExecutionContext(s.storage.Storage, s.client.Client, s.cfg)
	s.txsCtx = NewTestTxsContext(executionCtx, s.txsCtx.BatchType)

	addAccountWithHighNonce(s.Assertions, s.storage.Storage, 124)

	transfers := testutils.GenerateValidCreate2Transfers(9)
	s.invalidateCreate2Transfers(transfers[7:9])

	highNonceTransfers := []models.Create2Transfer{
		testutils.MakeCreate2Transfer(124, nil, 10, 1, &models.PublicKey{5, 4, 3, 2, 1}),
		testutils.MakeCreate2Transfer(124, nil, 11, 1, &models.PublicKey{5, 4, 3, 2, 1}),
	}

	transfers = append(transfers, highNonceTransfers...)
	s.addCreate2Transfers(transfers)

	preRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	result, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(result.Commitments(), 3)
	s.Len(result.Commitments()[0].Transactions, s.maxTxBytesInCommitment)
	s.Len(result.Commitments()[1].Transactions, s.maxTxBytesInCommitment)
	s.Len(result.Commitments()[2].Transactions, encoder.Create2TransferLength)

	postRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(result.Commitments()[2].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) invalidateCreate2Transfers(transfers []models.Create2Transfer) {
	for i := range transfers {
		tx := &transfers[i]
		tx.Amount = *genesisBalances[tx.FromStateID].MulN(10)
	}
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCommitments_ReturnsErrorWhenThereAreNotEnoughPendingTransfers() {
	preRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	result, err := s.txsCtx.CreateCommitments()
	s.Nil(result)
	s.ErrorIs(err, ErrNotEnoughTxs)

	postRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCommitments_ReturnsErrorWhenThereAreNotEnoughValidTransfers() {
	s.cfg = &config.RollupConfig{
		MinTxsPerCommitment:    32,
		MaxTxsPerCommitment:    32,
		FeeReceiverPubKeyID:    2,
		MaxCommitmentsPerBatch: 1,
		MinCommitmentsPerBatch: 1,
	}

	executionCtx := NewTestExecutionContext(s.storage.Storage, s.client.Client, s.cfg)
	s.txsCtx = NewTestTxsContext(executionCtx, s.txsCtx.BatchType)

	transfers := testutils.GenerateValidCreate2Transfers(2)
	s.addCreate2Transfers(transfers)

	pendingTransfers, err := s.storage.GetPendingCreate2Transfers()
	s.NoError(err)
	s.Len(pendingTransfers, 2)

	preRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	result, err := s.txsCtx.CreateCommitments()
	s.Nil(result)
	s.ErrorIs(err, ErrNotEnoughTxs)

	postRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	s.Equal(preRoot, postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCommitments_StoresCorrectCommitment() {
	transfersCount := uint32(4)
	s.preparePendingCreate2Transfers(transfersCount)

	preRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)

	expectedTxsLength := encoder.Create2TransferLength * int(transfersCount)
	result, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(result.Commitments(), 1)
	s.Len(result.Commitments()[0].Transactions, expectedTxsLength)
	s.Equal(result.Commitments()[0].FeeReceiver, uint32(2))

	postRoot, err := s.executionCtx.storage.StateTree.Root()
	s.NoError(err)
	s.NotEqual(preRoot, postRoot)
	s.Equal(result.Commitments()[0].PostStateRoot, *postRoot)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCommitments_CreatesMaximallyAsManyCommitmentsAsSpecifiedInConfig() {
	s.preparePendingCreate2Transfers(5)

	result, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(result.Commitments(), 1)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCommitments_UpdateTransfers() {
	s.preparePendingCreate2Transfers(2)

	pendingTransfers, err := s.storage.GetPendingCreate2Transfers()
	s.NoError(err)
	s.Len(pendingTransfers, 2)

	result, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(result.Commitments(), 1)

	for i := range pendingTransfers {
		tx, err := s.storage.GetCreate2Transfer(pendingTransfers[i].Hash)
		s.NoError(err)
		s.Equal(result.Commitments()[0].ID, *tx.CommitmentID)
		s.Equal(uint32(i+3), *tx.ToStateID)
	}
}

func (s *Create2TransferCommitmentsTestSuite) TestRegisterPendingAccounts_RegistersAccountsAndAddsMissingToAccountTree() {
	pendingAccounts := make([]models.AccountLeaf, st.AccountBatchSize-5)
	for i := 0; i < len(pendingAccounts); i++ {
		pendingAccounts[i] = models.AccountLeaf{
			PubKeyID:  uint32(st.AccountBatchOffset + i),
			PublicKey: models.PublicKey{byte(i), 8, 9},
		}
	}

	err := s.txsCtx.registerPendingAccounts(pendingAccounts)
	s.NoError(err)
	s.client.GetBackend().Commit()

	expectedAccounts := pendingAccounts
	for i := len(pendingAccounts); i < st.AccountBatchSize; i++ {
		expectedAccounts = append(expectedAccounts, models.AccountLeaf{
			PubKeyID:  uint32(st.AccountBatchOffset + i),
			PublicKey: mockPublicKey,
		})
	}

	registeredAccounts := s.getRegisteredAccounts(0)
	s.Equal(expectedAccounts, registeredAccounts)

	for i := len(pendingAccounts); i < st.AccountBatchSize; i++ {
		account, err := s.storage.AccountTree.Leaf(expectedAccounts[i].PubKeyID)
		s.NoError(err)
		s.Equal(expectedAccounts[i], *account)
	}
}

func (s *Create2TransferCommitmentsTestSuite) TestRegisterPendingAccounts_FillsMissingAccounts() {
	pendingAccounts := []models.AccountLeaf{
		{
			PubKeyID:  st.AccountBatchOffset,
			PublicKey: models.PublicKey{9, 8, 7},
		},
	}

	err := s.txsCtx.registerPendingAccounts(pendingAccounts)
	s.NoError(err)
	s.client.GetBackend().Commit()

	registeredAccounts := s.getRegisteredAccounts(0)
	s.Equal(pendingAccounts[0], registeredAccounts[0])
	for i := 1; i < len(pendingAccounts); i++ {
		s.Equal(mockPublicKey, registeredAccounts[i].PublicKey)
	}
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCommitments_RegistersAccounts() {
	transfers := testutils.GenerateValidCreate2Transfers(1)
	s.addCreate2Transfers(transfers)

	expectedTxsLength := encoder.Create2TransferLength * len(transfers)
	result, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(result.Commitments(), 1)
	s.Len(result.Commitments()[0].Transactions, expectedTxsLength)

	s.client.GetBackend().Commit()
	accounts := s.getRegisteredAccounts(0)
	s.Len(accounts, 16)
	s.Equal(transfers[0].ToPublicKey, accounts[0].PublicKey)
}

func (s *Create2TransferCommitmentsTestSuite) TestCreateCommitments_SkipsNonceTooHighTx() {
	transfers := testutils.GenerateValidCreate2Transfers(3)
	transfers[2].Nonce = models.MakeUint256(21)
	s.addCreate2Transfers(transfers)

	pendingTransfers, err := s.storage.GetPendingCreate2Transfers()
	s.NoError(err)
	s.Len(pendingTransfers, 3)

	result, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(result.Commitments(), 1)

	for i := 0; i < 2; i++ {
		var tx *models.Create2Transfer
		tx, err = s.storage.GetCreate2Transfer(pendingTransfers[i].Hash)
		s.NoError(err)
		s.Equal(result.Commitments()[0].ID, *tx.CommitmentID)
	}

	tx, err := s.storage.GetCreate2Transfer(transfers[2].Hash)
	s.NoError(err)
	s.Nil(tx.CommitmentID)

	pendingTransfers, err = s.storage.GetPendingCreate2Transfers()
	s.NoError(err)
	s.Len(pendingTransfers, 1)
}

func (s *Create2TransferCommitmentsTestSuite) getRegisteredAccounts(startBlockNumber uint64) []models.AccountLeaf {
	it, err := s.client.AccountRegistry.FilterBatchPubkeyRegistered(&bind.FilterOpts{Start: startBlockNumber})
	s.NoError(err)

	registeredAccounts := make([]models.AccountLeaf, 0)
	for it.Next() {
		tx, _, err := s.client.Blockchain.GetBackend().TransactionByHash(context.Background(), it.Event.Raw.TxHash)
		s.NoError(err)

		accounts, err := s.client.ExtractAccountsBatch(tx.Data(), it.Event)
		s.NoError(err)

		registeredAccounts = append(registeredAccounts, accounts...)
	}
	return registeredAccounts
}

func TestCreate2TransferCommitmentsTestSuite(t *testing.T) {
	suite.Run(t, new(Create2TransferCommitmentsTestSuite))
}

func (s *Create2TransferCommitmentsTestSuite) addCreate2Transfers(transfers []models.Create2Transfer) {
	err := s.storage.BatchAddCreate2Transfer(transfers)
	s.NoError(err)
}

func (s *Create2TransferCommitmentsTestSuite) preparePendingCreate2Transfers(transfersAmount uint32) {
	transfers := testutils.GenerateValidCreate2Transfers(transfersAmount)
	s.addCreate2Transfers(transfers)
}
