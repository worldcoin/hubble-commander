package api

// TODO: rename to mempool_test.go?

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/mempool"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetUserStateTestSuite struct {
	*require.Assertions
	suite.Suite
	api      *API
	teardown func() error
}

func (s *GetUserStateTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetUserStateTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorage()
	s.NoError(err)
	s.teardown = testStorage.Teardown

	s.api = NewTestAPI(
		testStorage.Storage,
		eth.DomainOnlyTestClient,
		mempool.NewTestTxPool(),
	)
}

func (s *GetUserStateTestSuite) TearDownTest() {
	err := s.teardown()
	s.NoError(err)
}

func (s *GetUserStateTestSuite) TestGetUserState() {

	err := s.api.storage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: models.PublicKey{1, 2, 3},
	})
	s.NoError(err)

	// fee receiver
	_, err = s.api.storage.StateTree.Set(
		0,
		&models.UserState{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(0),
			Nonce:    models.MakeUint256(0),
		},
	)
	s.NoError(err)

	// sender
	_, err = s.api.storage.StateTree.Set(
		1,
		&models.UserState{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(100),
			Nonce:    models.MakeUint256(0),
		},
	)
	s.NoError(err)

	// receiver
	_, err = s.api.storage.StateTree.Set(
		2,
		&models.UserState{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(0),
			Nonce:    models.MakeUint256(0),
		},
	)
	s.NoError(err)

	userState, err := s.api.GetUserState(context.Background(), 1)
	s.NoError(err)

	s.Equal(userState.StateID, uint32(1))
	s.Equal(userState.UserState.Nonce.Int.Uint64(), uint64(0))
	s.Equal(userState.UserState.Balance.Int.Uint64(), uint64(100))

	transfer := dto.Transfer{
		FromStateID: ref.Uint32(1),
		ToStateID:   ref.Uint32(2),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   &models.Signature{},
	}

	hash, err := s.api.SendTransaction(dto.MakeTransaction(transfer))
	s.NoError(err)
	s.NotNil(hash)

	// the rollup loop is not running (it doesn't even exist) but the api immediately
	// acts as if our transaction has been applied:

	userState, err = s.api.GetUserState(context.Background(), 1)
	s.NoError(err)

	s.Equal(uint32(1), userState.StateID)
	s.Equal(uint64(1), userState.UserState.Nonce.Int.Uint64())
	s.Equal(uint64(40), userState.UserState.Balance.Int.Uint64())
}

func TestGetUserStateTestSuite(t *testing.T) {
	suite.Run(t, new(GetUserStateTestSuite))
}
