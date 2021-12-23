package api

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetTransactionsTestSuite struct {
	*require.Assertions
	suite.Suite
	api     *API
	storage *st.TestStorage
}

func (s *GetTransactionsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetTransactionsTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)
	s.api = &API{storage: s.storage.Storage}
}

func (s *GetTransactionsTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetTransactionsTestSuite) addAccounts() {
	accounts := []models.AccountLeaf{
		{
			PubKeyID:  0,
			PublicKey: models.PublicKey{1, 1, 1},
		},
		{
			PubKeyID:  1,
			PublicKey: models.PublicKey{1, 1, 1},
		},
		{
			PubKeyID:  2,
			PublicKey: models.PublicKey{2, 2, 2},
		},
		{
			PubKeyID:  3,
			PublicKey: models.PublicKey{3, 3, 3},
		},
	}
	for i := range accounts {
		err := s.storage.AccountTree.SetSingle(&accounts[i])
		s.NoError(err)
	}
}

func (s *GetTransactionsTestSuite) addUserStates() {
	makeUserState := func(tokenID uint64, pubKeyID uint32) models.UserState {
		return models.UserState{
			PubKeyID: pubKeyID,
			TokenID:  models.MakeUint256(tokenID),
			Balance:  models.MakeUint256(500),
			Nonce:    models.MakeUint256(0),
		}
	}

	userStates := []models.UserState{
		makeUserState(1, 0), // StateID = 0
		makeUserState(1, 1), // StateID = 1
		makeUserState(1, 2), // StateID = 2
		makeUserState(2, 0), // StateID = 3
		makeUserState(2, 2), // StateID = 4
	}

	for i := range userStates {
		_, err := s.storage.StateTree.Set(uint32(i), &userStates[i])
		s.NoError(err)
	}
}

func (s *GetTransactionsTestSuite) addTransfers() []models.Transfer {
	transfers := []models.Transfer{
		testutils.MakeTransfer(0, 2, 0, 10),
		testutils.MakeTransfer(1, 2, 0, 10),
		testutils.MakeTransfer(2, 4, 0, 10),
		testutils.MakeTransfer(2, 0, 0, 10),
		testutils.MakeTransfer(3, 4, 0, 10),
		testutils.MakeTransfer(4, 3, 0, 10),
	}

	err := s.storage.BatchAddTransfer(transfers)
	s.NoError(err)
	return transfers
}

func (s *GetTransactionsTestSuite) addCreate2Transfers() []models.Create2Transfer {
	transfers := []models.Create2Transfer{
		testutils.MakeCreate2Transfer(0, ref.Uint32(5), 0, 10, &models.PublicKey{3, 4, 5}),
		testutils.MakeCreate2Transfer(4, ref.Uint32(6), 0, 10, &models.PublicKey{3, 4, 5}),
		testutils.MakeCreate2Transfer(4, ref.Uint32(3), 0, 10, &models.PublicKey{3, 4, 5}),
	}

	err := s.storage.BatchAddCreate2Transfer(transfers)
	s.NoError(err)
	return transfers
}

func (s *GetTransactionsTestSuite) TestGetTransactions() {
	s.addAccounts()
	s.addUserStates()
	transfers := s.addTransfers()
	create2Transfers := s.addCreate2Transfers()

	txs, err := s.api.GetTransactions(&models.PublicKey{1, 1, 1})
	s.NoError(err)

	newTransferReceipt := func(transfer models.Transfer) *dto.TransferReceipt {
		return &dto.TransferReceipt{
			TransferWithBatchDetails: dto.MakeTransferWithBatchDetailsFromTransfer(&transfer),
			Status:                   txstatus.Pending,
		}
	}

	newCreate2Receipt := func(transfer models.Create2Transfer) *dto.Create2TransferReceipt {
		return &dto.Create2TransferReceipt{
			Create2TransferWithBatchDetails: dto.MakeCreate2TransferWithBatchDetailsFromCreate2Transfer(&transfer),
			Status:                          txstatus.Pending,
		}
	}

	s.Len(txs, 7)
	s.Contains(txs, newTransferReceipt(transfers[0]))
	s.Contains(txs, newTransferReceipt(transfers[1]))
	s.Contains(txs, newTransferReceipt(transfers[3]))
	s.Contains(txs, newTransferReceipt(transfers[4]))
	s.Contains(txs, newTransferReceipt(transfers[5]))
	s.Contains(txs, newCreate2Receipt(create2Transfers[0]))
	s.Contains(txs, newCreate2Receipt(create2Transfers[2]))
}

func (s *GetTransactionsTestSuite) addCommitmentAndBatch() *models.Batch {
	batch := &models.Batch{
		ID:                models.MakeUint256(1),
		Type:              batchtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              ref.Hash(utils.RandomHash()),
		FinalisationBlock: ref.Uint32(1234),
		AccountTreeRoot:   ref.Hash(utils.RandomHash()),
		PrevStateRoot:     ref.Hash(utils.RandomHash()),
		SubmissionTime:    models.NewTimestamp(time.Unix(140, 0).UTC()),
	}
	err := s.storage.AddBatch(batch)
	s.NoError(err)

	commitment := &models.TxCommitment{
		CommitmentBase: models.CommitmentBase{
			ID: models.CommitmentID{
				BatchID:      batch.ID,
				IndexInBatch: 0,
			},
			Type:          batchtype.Transfer,
			PostStateRoot: utils.RandomHash(),
		},
		FeeReceiver:       1234,
		CombinedSignature: models.MakeRandomSignature(),
	}

	err = s.storage.AddTxCommitment(commitment)
	s.NoError(err)
	return batch
}

func (s *GetTransactionsTestSuite) addIncludedTransfer() models.Transfer {
	transfer := testutils.MakeTransfer(0, 2, 0, 10)
	transfer.CommitmentID = &models.CommitmentID{
		BatchID:      models.MakeUint256(1),
		IndexInBatch: 0,
	}
	err := s.storage.AddTransfer(&transfer)
	s.NoError(err)
	return transfer
}

func (s *GetTransactionsTestSuite) addIncludedCreate2Transfer() models.Create2Transfer {
	create2Transfer := testutils.MakeCreate2Transfer(0, ref.Uint32(5), 0, 10, &models.PublicKey{3, 4, 5})
	create2Transfer.CommitmentID = &models.CommitmentID{
		BatchID:      models.MakeUint256(1),
		IndexInBatch: 0,
	}
	err := s.storage.AddCreate2Transfer(&create2Transfer)
	s.NoError(err)
	return create2Transfer
}

func (s *GetTransactionsTestSuite) TestGetTransactions_ReceiptsWithDetails() {
	s.addAccounts()
	s.addUserStates()
	batch := s.addCommitmentAndBatch()

	transfer := s.addIncludedTransfer()
	create2Transfer := s.addIncludedCreate2Transfer()

	txs, err := s.api.GetTransactions(&models.PublicKey{1, 1, 1})
	s.NoError(err)

	newIncludedTransferReceipt := func(transfer models.Transfer) *dto.TransferReceipt {
		transferWithBatchDetails := dto.MakeTransferWithBatchDetailsFromTransfer(&transfer)
		transferWithBatchDetails.BatchHash = batch.Hash
		transferWithBatchDetails.BatchTime = batch.SubmissionTime
		return &dto.TransferReceipt{
			TransferWithBatchDetails: transferWithBatchDetails,
			Status:                   txstatus.InBatch,
		}
	}

	newIncludedCreate2Receipt := func(transfer models.Create2Transfer) *dto.Create2TransferReceipt {
		create2TransferWithBatchDetails := dto.MakeCreate2TransferWithBatchDetailsFromCreate2Transfer(&transfer)
		create2TransferWithBatchDetails.BatchHash = batch.Hash
		create2TransferWithBatchDetails.BatchTime = batch.SubmissionTime
		return &dto.Create2TransferReceipt{
			Create2TransferWithBatchDetails: create2TransferWithBatchDetails,
			Status:                          txstatus.InBatch,
		}
	}

	s.Len(txs, 2)
	s.Contains(txs, newIncludedTransferReceipt(transfer))
	s.Contains(txs, newIncludedCreate2Receipt(create2Transfer))
}

func (s *GetTransactionsTestSuite) TestGetTransactions_NoTransactions() {
	account := models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: models.PublicKey{1, 2, 3},
	}

	err := s.storage.AccountTree.SetSingle(&account)
	s.NoError(err)

	userState := &models.UserState{
		PubKeyID: account.PubKeyID,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	}
	_, err = s.storage.StateTree.Set(0, userState)
	s.NoError(err)

	userTransfers, err := s.api.GetTransactions(&account.PublicKey)
	s.NoError(err)

	s.NotNil(userTransfers)
	s.Len(userTransfers, 0)
}

func (s *GetTransactionsTestSuite) TestGetTransactions_SingleCreate2Transfer() {
	s.addAccounts()
	s.addUserStates()

	c2t := testutils.MakeCreate2Transfer(0, nil, 0, 10, &models.PublicKey{1, 1, 1})
	err := s.storage.AddCreate2Transfer(&c2t)
	s.NoError(err)

	txs, err := s.api.GetTransactions(&models.PublicKey{1, 1, 1})
	s.NoError(err)

	s.Len(txs, 1)
}

func (s *GetTransactionsTestSuite) TestGetTransactions_SingleTransfer() {
	s.addAccounts()
	s.addUserStates()

	transfer := testutils.MakeTransfer(0, 2, 0, 10)
	err := s.storage.AddTransfer(&transfer)
	s.NoError(err)

	txs, err := s.api.GetTransactions(&models.PublicKey{1, 1, 1})
	s.NoError(err)

	s.Len(txs, 1)
}

func TestGetTransactionsTestSuite(t *testing.T) {
	suite.Run(t, new(GetTransactionsTestSuite))
}
