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

type C2TCommitmentsTestSuite struct {
	testSuiteWithTxsContext
	maxTxBytesInCommitment int
}

func (s *C2TCommitmentsTestSuite) SetupTest() {
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

func (s *C2TCommitmentsTestSuite) TestCreateCommitments_UpdatesTransactions() {
	transfers := testutils.GenerateValidCreate2Transfers(2)
	s.addCreate2Transfers(transfers)

	batchData, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(batchData.Commitments(), 1)

	for i := range transfers {
		tx, err := s.storage.GetCreate2Transfer(transfers[i].Hash)
		s.NoError(err)
		s.Equal(batchData.Commitments()[0].ID, *tx.CommitmentID)
		s.Equal(uint32(i+3), *tx.ToStateID)
		s.Nil(tx.ErrorMessage)
	}
}

func (s *C2TCommitmentsTestSuite) TestCreateCommitments_RegistersAccounts() {
	transfers := testutils.GenerateValidCreate2Transfers(1)
	s.addCreate2Transfers(transfers)

	expectedTxsLength := encoder.Create2TransferLength * len(transfers)
	batchData, err := s.txsCtx.CreateCommitments()
	s.NoError(err)
	s.Len(batchData.Commitments(), 1)
	s.Len(batchData.Commitments()[0].Transactions, expectedTxsLength)

	s.client.GetBackend().Commit()
	accounts := s.getRegisteredAccounts(0)
	s.Len(accounts, 16)
	s.Equal(transfers[0].ToPublicKey, accounts[0].PublicKey)
}

func (s *C2TCommitmentsTestSuite) TestRegisterPendingAccounts_RegistersAccountsAndAddsMissingToAccountTree() {
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

func (s *C2TCommitmentsTestSuite) TestRegisterPendingAccounts_FillsMissingAccounts() {
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

func (s *C2TCommitmentsTestSuite) getRegisteredAccounts(startBlockNumber uint64) []models.AccountLeaf {
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

func TestC2TCommitmentsTestSuite(t *testing.T) {
	suite.Run(t, new(C2TCommitmentsTestSuite))
}

func (s *C2TCommitmentsTestSuite) addCreate2Transfers(transfers []models.Create2Transfer) {
	err := s.storage.BatchAddCreate2Transfer(transfers)
	s.NoError(err)
}
