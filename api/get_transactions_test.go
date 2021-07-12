package api

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
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
	tree    *st.StateTree
}

func (s *GetTransactionsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetTransactionsTestSuite) SetupTest() {
	var err error
	s.storage, err = st.NewTestStorageWithBadger()
	s.NoError(err)
	s.api = &API{storage: s.storage.Storage}
	s.tree = st.NewStateTree(s.storage.Storage)
}

func (s *GetTransactionsTestSuite) TearDownTest() {
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *GetTransactionsTestSuite) addAccounts() {
	accounts := []models.Account{
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
		err := s.storage.AddAccountIfNotExists(&accounts[i])
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
		_, err := s.tree.Set(uint32(i), &userStates[i])
		s.NoError(err)
	}
}

func (s *GetTransactionsTestSuite) makeTransfer(fromStateID, toStateID uint32) models.Transfer {
	return models.Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.Transfer,
			FromStateID: fromStateID,
			Amount:      models.MakeUint256(10),
			Fee:         models.MakeUint256(1),
			Nonce:       models.MakeUint256(0),
			Signature:   models.MakeRandomSignature(),
		},
		ToStateID: toStateID,
	}
}

func (s *GetTransactionsTestSuite) addTransfers() []models.Transfer {
	transfers := []models.Transfer{
		s.makeTransfer(0, 2),
		s.makeTransfer(1, 2),
		s.makeTransfer(2, 4),
		s.makeTransfer(2, 0),
		s.makeTransfer(3, 4),
		s.makeTransfer(4, 3),
	}

	for i := range transfers {
		receiveTime, err := s.storage.AddTransfer(&transfers[i])
		s.NoError(err)
		transfers[i].ReceiveTime = receiveTime
	}
	return transfers
}

func (s *GetTransactionsTestSuite) makeCreate2Transfer(
	fromStateID, toStateID uint32,
	toPublicKey *models.PublicKey,
) models.Create2Transfer {
	return models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        utils.RandomHash(),
			TxType:      txtype.Create2Transfer,
			FromStateID: fromStateID,
			Amount:      models.MakeUint256(10),
			Fee:         models.MakeUint256(1),
			Nonce:       models.MakeUint256(1),
			Signature:   models.MakeRandomSignature(),
		},
		ToStateID:   ref.Uint32(toStateID),
		ToPublicKey: *toPublicKey,
	}
}

func (s *GetTransactionsTestSuite) addCreate2Transfers() []models.Create2Transfer {
	transfers := []models.Create2Transfer{
		s.makeCreate2Transfer(0, 5, &models.PublicKey{3, 4, 5}),
		s.makeCreate2Transfer(4, 6, &models.PublicKey{3, 4, 5}),
		s.makeCreate2Transfer(4, 3, &models.PublicKey{3, 4, 5}),
	}

	for i := range transfers {
		receiveTime, err := s.storage.AddCreate2Transfer(&transfers[i])
		s.NoError(err)
		transfers[i].ReceiveTime = receiveTime
	}
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
			TransferWithBatchDetails: models.TransferWithBatchDetails{
				Transfer: transfer,
			},
			Status: txstatus.Pending,
		}
	}

	newCreate2Receipt := func(transfer models.Create2Transfer) *dto.Create2TransferReceipt {
		return &dto.Create2TransferReceipt{
			Create2TransferWithBatchDetails: models.Create2TransferWithBatchDetails{
				Create2Transfer: transfer,
			},
			Status: txstatus.Pending,
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
		Type:              txtype.Transfer,
		TransactionHash:   utils.RandomHash(),
		Hash:              ref.Hash(utils.RandomHash()),
		FinalisationBlock: ref.Uint32(1234),
		AccountTreeRoot:   ref.Hash(utils.RandomHash()),
		PrevStateRoot:     ref.Hash(utils.RandomHash()),
		SubmissionTime:    models.NewTimestamp(time.Now().UTC()),
	}
	err := s.storage.AddBatch(batch)
	s.NoError(err)

	commitment := &models.Commitment{
		ID:                1,
		Type:              txtype.Transfer,
		Transactions:      utils.RandomBytes(12),
		FeeReceiver:       1234,
		CombinedSignature: models.MakeRandomSignature(),
		PostStateRoot:     utils.RandomHash(),
		IncludedInBatch:   &batch.ID,
	}

	_, err = s.storage.AddCommitment(commitment)
	s.NoError(err)
	return batch
}

func (s *GetTransactionsTestSuite) addIncludedTransfer() models.Transfer {
	transfer := s.makeTransfer(0, 2)
	transfer.IncludedInCommitment = ref.Int32(1)
	receiveTime, err := s.storage.AddTransfer(&transfer)
	s.NoError(err)
	transfer.ReceiveTime = receiveTime
	return transfer
}

func (s *GetTransactionsTestSuite) addIncludedCreate2Transfer() models.Create2Transfer {
	create2Transfer := s.makeCreate2Transfer(0, 5, &models.PublicKey{3, 4, 5})
	create2Transfer.IncludedInCommitment = ref.Int32(1)
	receiveTime2, err := s.storage.AddCreate2Transfer(&create2Transfer)
	s.NoError(err)
	create2Transfer.ReceiveTime = receiveTime2
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
		return &dto.TransferReceipt{
			TransferWithBatchDetails: models.TransferWithBatchDetails{
				Transfer:  transfer,
				BatchHash: batch.Hash,
				BatchTime: batch.SubmissionTime,
			},
			Status: txstatus.InBatch,
		}
	}

	newIncludedCreate2Receipt := func(transfer models.Create2Transfer) *dto.Create2TransferReceipt {
		return &dto.Create2TransferReceipt{
			Create2TransferWithBatchDetails: models.Create2TransferWithBatchDetails{
				Create2Transfer: transfer,
				BatchHash:       batch.Hash,
				BatchTime:       batch.SubmissionTime,
			},
			Status: txstatus.InBatch,
		}
	}

	s.Len(txs, 2)
	s.Contains(txs, newIncludedTransferReceipt(transfer))
	s.Contains(txs, newIncludedCreate2Receipt(create2Transfer))
}

func (s *GetTransactionsTestSuite) TestGetTransactions_NoTransactions() {
	account := models.Account{
		PubKeyID:  1,
		PublicKey: models.PublicKey{1, 2, 3},
	}

	err := s.storage.AddAccountIfNotExists(&account)
	s.NoError(err)

	userState := &models.UserState{
		PubKeyID: account.PubKeyID,
		TokenID:  models.MakeUint256(1),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	}
	_, err = s.tree.Set(0, userState)
	s.NoError(err)

	userTransfers, err := s.api.GetTransactions(&account.PublicKey)
	s.NoError(err)

	s.NotNil(userTransfers)
	s.Len(userTransfers, 0)
}

func TestGetTransactionsTestSuite(t *testing.T) {
	suite.Run(t, new(GetTransactionsTestSuite))
}
