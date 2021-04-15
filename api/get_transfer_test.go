package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetTransferTestSuite struct {
	*require.Assertions
	suite.Suite
	api      *API
	db       *db.TestDB
	transfer *models.Transfer
}

func (s *GetTransferTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetTransferTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)

	storage := st.NewTestStorage(testDB.DB)
	s.api = &API{nil, storage, nil}
	s.db = testDB

	userState := models.UserState{
		AccountIndex: 1,
		TokenIndex:   models.MakeUint256(1),
		Balance:      models.MakeUint256(420),
		Nonce:        models.MakeUint256(0),
	}

	tree := st.NewStateTree(storage)
	err = tree.Set(1, &userState)
	s.NoError(err)

	transfer := &models.Transfer{
		FromStateID: 1,
		ToStateID:   2,
		Amount:      *models.NewUint256(50),
		Fee:         *models.NewUint256(10),
		Nonce:       *models.NewUint256(0),
		Signature:   []byte{1, 2, 3, 4},
	}

	s.transfer = transfer
}

func (s *GetTransferTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *GetTransferTestSuite) TestApi_GetTransfer() {
	transfer := dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   []byte{1, 2, 3, 4},
	}

	hash, err := s.api.SendTransaction(dto.MakeTransaction(transfer))
	s.NoError(err)

	res, err := s.api.GetTransfer(*hash)
	s.NoError(err)

	s.Equal(models.Pending, res.Status)
}

func TestGetTransferTestSuite(t *testing.T) {
	suite.Run(t, new(GetTransferTestSuite))
}
