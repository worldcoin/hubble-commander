package api

import (
	"testing"

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
		err := s.tree.Set(uint32(i), &userStates[i])
		s.NoError(err)
	}
}

func (s *GetTransactionsTestSuite) addTransfers() []models.Transfer {
	makeTransfer := func(fromStateID, toStateID uint32) models.Transfer {
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

	transfers := []models.Transfer{
		makeTransfer(0, 2),
		makeTransfer(1, 2),
		makeTransfer(2, 4),
		makeTransfer(2, 0),
		makeTransfer(3, 4),
		makeTransfer(4, 3),
	}

	for i := range transfers {
		err := s.storage.AddTransfer(&transfers[i])
		s.NoError(err)
	}
	return transfers
}

func (s *GetTransactionsTestSuite) addCreate2Transfers() []models.Create2Transfer {
	makeCreate2Transfer := func(fromStateID, toStateID uint32, toPublicKey models.PublicKey) models.Create2Transfer {
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
			ToPublicKey: toPublicKey,
		}
	}

	transfers := []models.Create2Transfer{
		makeCreate2Transfer(0, 5, models.PublicKey{3, 4, 5}),
		makeCreate2Transfer(4, 6, models.PublicKey{3, 4, 5}),
		makeCreate2Transfer(4, 3, models.PublicKey{3, 4, 5}),
	}

	for i := range transfers {
		err := s.storage.AddCreate2Transfer(&transfers[i])
		s.NoError(err)
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
			TransferWithBatchHash: models.TransferWithBatchHash{
				Transfer: transfer,
			},
			Status: txstatus.Pending,
		}
	}

	newCreate2Receipt := func(transfer models.Create2Transfer) *dto.Create2TransferReceipt {
		return &dto.Create2TransferReceipt{
			Create2TransferWithBatchHash: models.Create2TransferWithBatchHash{
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
	err = s.tree.Set(0, userState)
	s.NoError(err)

	userTransfers, err := s.api.GetTransactions(&account.PublicKey)
	s.NoError(err)

	s.NotNil(userTransfers)
	s.Len(userTransfers, 0)
}

func TestGetTransactionsTestSuite(t *testing.T) {
	suite.Run(t, new(GetTransactionsTestSuite))
}
