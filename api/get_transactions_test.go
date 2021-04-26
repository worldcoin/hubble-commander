package api

import (
	"math/big"
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/txstatus"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetTransactionsTestSuite struct {
	*require.Assertions
	suite.Suite
	api     *API
	storage *st.Storage
	db      *db.TestDB
	tree    *st.StateTree
}

func (s *GetTransactionsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetTransactionsTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)

	s.storage = st.NewTestStorage(testDB.DB)
	s.api = &API{nil, s.storage, nil}
	s.db = testDB
	s.tree = st.NewStateTree(s.storage)
}

func (s *GetTransactionsTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

// nolint:funlen
func (s *GetTransactionsTestSuite) TestGetTransactions() {
	accounts := []models.Account{
		{
			PubKeyID:  1,
			PublicKey: models.PublicKey{1, 2, 3},
		},
		{
			PubKeyID:  2,
			PublicKey: models.PublicKey{2, 3, 4},
		},
	}
	for i := range accounts {
		err := s.storage.AddAccountIfNotExists(&accounts[i])
		s.NoError(err)
	}

	userStates := []models.UserState{
		{
			PubKeyID:   accounts[0].PubKeyID,
			TokenIndex: models.MakeUint256(1),
			Balance:    models.MakeUint256(420),
			Nonce:      models.MakeUint256(0),
		},
		{
			PubKeyID:   accounts[1].PubKeyID,
			TokenIndex: models.MakeUint256(2),
			Balance:    models.MakeUint256(500),
			Nonce:      models.MakeUint256(0),
		},
		{
			PubKeyID:   accounts[0].PubKeyID,
			TokenIndex: models.MakeUint256(25),
			Balance:    models.MakeUint256(1),
			Nonce:      models.MakeUint256(73),
		},
	}

	err := s.tree.Set(0, &userStates[0])
	s.NoError(err)
	err = s.tree.Set(1, &userStates[1])
	s.NoError(err)
	err = s.tree.Set(2, &userStates[2])
	s.NoError(err)

	transfers := []models.Transfer{
		{
			TransactionBase: models.TransactionBase{
				Hash:                 common.BigToHash(big.NewInt(1234)),
				TxType:               txtype.Transfer,
				FromStateID:          0,
				Amount:               models.MakeUint256(1),
				Fee:                  models.MakeUint256(5),
				Nonce:                models.MakeUint256(0),
				Signature:            []byte{1, 2, 3, 4, 5},
				IncludedInCommitment: nil,
			},
			ToStateID: 1,
		},
		{
			TransactionBase: models.TransactionBase{
				Hash:                 common.BigToHash(big.NewInt(2345)),
				TxType:               txtype.Transfer,
				FromStateID:          0,
				Amount:               models.MakeUint256(2),
				Fee:                  models.MakeUint256(5),
				Nonce:                models.MakeUint256(1),
				Signature:            []byte{2, 3, 4, 5, 6},
				IncludedInCommitment: nil,
			},
			ToStateID: 1,
		},
		{
			TransactionBase: models.TransactionBase{
				Hash:                 common.BigToHash(big.NewInt(3456)),
				TxType:               txtype.Transfer,
				FromStateID:          1,
				Amount:               models.MakeUint256(3),
				Fee:                  models.MakeUint256(5),
				Nonce:                models.MakeUint256(0),
				Signature:            []byte{3, 4, 5, 6, 7},
				IncludedInCommitment: nil,
			},
			ToStateID: 0,
		},
		{
			TransactionBase: models.TransactionBase{
				Hash:                 common.BigToHash(big.NewInt(4567)),
				TxType:               txtype.Transfer,
				FromStateID:          0,
				Amount:               models.MakeUint256(2),
				Fee:                  models.MakeUint256(5),
				Nonce:                models.MakeUint256(2),
				Signature:            []byte{2, 3, 4, 5, 6},
				IncludedInCommitment: nil,
			},
			ToStateID: 1,
		},
	}

	err = s.storage.AddTransfer(&transfers[0])
	s.NoError(err)
	err = s.storage.AddTransfer(&transfers[1])
	s.NoError(err)
	err = s.storage.AddTransfer(&transfers[2])
	s.NoError(err)
	err = s.storage.AddTransfer(&transfers[3])
	s.NoError(err)

	create2Transfer := models.Create2Transfer{
		TransactionBase: models.TransactionBase{
			Hash:        common.BigToHash(big.NewInt(1111)),
			TxType:      txtype.Create2Transfer,
			FromStateID: 0,
			Amount:      models.MakeUint256(100),
			Fee:         models.MakeUint256(5),
			Nonce:       models.MakeUint256(1),
		},
		ToStateID:  1,
		ToPubKeyID: accounts[1].PubKeyID,
	}
	err = s.storage.AddCreate2Transfer(&create2Transfer)
	s.NoError(err)

	userTransfers, err := s.api.GetTransactions(&accounts[0].PublicKey)
	s.NoError(err)

	s.Len(userTransfers, 4)
	s.Contains(userTransfers, &dto.TransferReceipt{
		Transfer: transfers[0],
		Status:   txstatus.Pending,
	})
	s.NotContains(userTransfers, &dto.TransferReceipt{
		Transfer: transfers[2],
		Status:   txstatus.Pending,
	})
}

func (s *GetTransactionsTestSuite) TestGetTransactions_NoTransactions() {
	account := models.Account{
		PubKeyID:  1,
		PublicKey: models.PublicKey{1, 2, 3},
	}

	err := s.storage.AddAccountIfNotExists(&account)
	s.NoError(err)

	userState := &models.UserState{
		PubKeyID:   account.PubKeyID,
		TokenIndex: models.MakeUint256(1),
		Balance:    models.MakeUint256(420),
		Nonce:      models.MakeUint256(0),
	}
	err = s.tree.Set(0, userState)
	s.NoError(err)

	userTransfers, err := s.api.GetTransactions(&account.PublicKey)
	s.NoError(err)

	s.Len(userTransfers, 0)
	s.NotNil(userTransfers)
}

func TestGetTransactionsTestSuite(t *testing.T) {
	suite.Run(t, new(GetTransactionsTestSuite))
}
