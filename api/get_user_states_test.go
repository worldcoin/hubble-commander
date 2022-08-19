package api

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GetUserStatesTestSuite struct {
	*require.Assertions
	suite.Suite
	api      *API
	teardown func() error
}

func (s *GetUserStatesTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *GetUserStatesTestSuite) SetupTest() {
	testStorage, err := st.NewTestStorage()
	s.NoError(err)
	s.teardown = testStorage.Teardown
	s.api = NewTestAPI(
		testStorage.Storage,
		eth.DomainOnlyTestClient,
	)
}

func (s *GetUserStatesTestSuite) TearDownTest() {
	err := s.teardown()
	s.NoError(err)
}

func (s *GetUserStatesTestSuite) TestGetUserStates_NoSuchState() {
	pubkey := testutils.RandomPublicKey()
	_, err := s.api.GetUserStates(context.Background(), &pubkey)
	s.Error(err)
	s.Equal(err.Error(), "user states not found")

	// TODO: why is the following not true?
	// s.True(storage.IsNotFoundError(err))
}

func (s *GetUserStatesTestSuite) TestGetUserStates_ZipsBatchednAndPendingStates() {
	// user states should return both the batched and pending states
	senderAccount := models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: models.PublicKey{1, 2, 3},
	}
	err := s.api.storage.AccountTree.SetSingle(&senderAccount)
	s.NoError(err)

	senderStateID := 1
	_, err = s.api.storage.StateTree.Set(
		uint32(senderStateID),
		&models.UserState{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(0),
			Balance:  models.MakeUint256(100),
			Nonce:    models.MakeUint256(0),
		},
	)
	s.NoError(err)

	domain, err := s.api.client.GetDomain()
	s.NoError(err)

	receiverWallet, err := bls.NewRandomWallet(*domain)
	s.NoError(err)

	receiverAccount := models.AccountLeaf{
		PubKeyID:  2,
		PublicKey: *receiverWallet.PublicKey(),
	}
	err = s.api.storage.AccountTree.SetSingle(&receiverAccount)
	s.NoError(err)

	receiverStateID := 2
	_, err = s.api.storage.StateTree.Set(
		uint32(receiverStateID),
		&models.UserState{
			PubKeyID: 2,
			TokenID:  models.MakeUint256(0),
			Balance:  models.MakeUint256(10),
			Nonce:    models.MakeUint256(0),
		},
	)
	s.NoError(err)

	c2t := dto.Create2Transfer{
		FromStateID: ref.Uint32(1),
		ToPublicKey: receiverWallet.PublicKey(),
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   &models.Signature{},
	}

	hash, err := s.api.SendTransaction(context.Background(), dto.MakeTransaction(c2t))
	s.NoError(err)
	s.NotNil(hash)

	userStates, err := s.api.GetUserStates(context.Background(), &receiverAccount.PublicKey)
	s.NoError(err)
	s.Len(userStates, 2)

	s.Equal(dto.UserStateWithID{
		StateID: 2,
		UserState: dto.UserState{
			PubKeyID: 2,
			TokenID:  models.MakeUint256(0),
			Balance:  models.MakeUint256(10),
			Nonce:    models.MakeUint256(0),
		},
	}, userStates[0])
	s.Equal(dto.UserStateWithID{
		StateID: -1,
		UserState: dto.UserState{
			PubKeyID: 2,
			TokenID:  models.MakeUint256(0),
			Balance:  models.MakeUint256(50),
			Nonce:    models.MakeUint256(0),
		},
	}, userStates[1])
}

func (s *GetUserStatesTestSuite) TestGetUserStates_HasPendingC2T() {
	account := models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: models.PublicKey{1, 2, 3},
	}
	err := s.api.storage.AccountTree.SetSingle(&account)
	s.NoError(err)

	senderStateID := 1
	_, err = s.api.storage.StateTree.Set(
		uint32(senderStateID),
		&models.UserState{
			PubKeyID: 1,
			TokenID:  models.MakeUint256(1),
			Balance:  models.MakeUint256(100),
			Nonce:    models.MakeUint256(0),
		},
	)
	s.NoError(err)

	pubkey := testutils.RandomPublicKey()

	c2t := dto.Create2Transfer{
		FromStateID: ref.Uint32(1),
		ToPublicKey: &pubkey,
		Amount:      models.NewUint256(50),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   &models.Signature{},
	}

	hash, err := s.api.SendTransaction(context.Background(), dto.MakeTransaction(c2t))
	s.NoError(err)
	s.NotNil(hash)

	userStates, err := s.api.GetUserStates(context.Background(), &pubkey)
	s.NoError(err)
	s.Len(userStates, 1)

	s.Equal(userStates[0], dto.UserStateWithID{
		StateID: -1,
		UserState: dto.UserState{
			PubKeyID: -1,
			TokenID:  models.MakeUint256(0),
			Balance:  models.MakeUint256(50),
			Nonce:    models.MakeUint256(0),
		},
	})
}

func (s *GetUserStatesTestSuite) TestGetUserStates() {
	accounts := []models.AccountLeaf{
		{
			PubKeyID:  1,
			PublicKey: models.PublicKey{1, 2, 3},
		},
		{
			PubKeyID:  2,
			PublicKey: models.PublicKey{1, 2, 3},
		},
	}
	for i := range accounts {
		err := s.api.storage.AccountTree.SetSingle(&accounts[i])
		s.NoError(err)
	}

	leaves := []models.StateLeaf{
		{
			StateID: 0,
			UserState: models.UserState{
				PubKeyID: accounts[0].PubKeyID,
				TokenID:  models.MakeUint256(1),
				Balance:  models.MakeUint256(420),
				Nonce:    models.MakeUint256(0),
			},
		},
		{
			StateID: 1,
			UserState: models.UserState{
				PubKeyID: accounts[1].PubKeyID,
				TokenID:  models.MakeUint256(2),
				Balance:  models.MakeUint256(500),
				Nonce:    models.MakeUint256(0),
			},
		},
		{
			StateID: 2,
			UserState: models.UserState{
				PubKeyID: accounts[0].PubKeyID,
				TokenID:  models.MakeUint256(25),
				Balance:  models.MakeUint256(1),
				Nonce:    models.MakeUint256(73),
			},
		},
	}
	for i := range leaves {
		_, err := s.api.storage.StateTree.Set(leaves[i].StateID, &leaves[i].UserState)
		s.NoError(err)
	}

	userStates, err := s.api.GetUserStates(context.Background(), &accounts[0].PublicKey)
	s.NoError(err)

	s.Len(userStates, 3)
	s.Equal(userStates[0], dto.UserStateWithID{
		StateID:   0,
		UserState: dto.MakeUserState(&leaves[0].UserState),
	})
	s.Equal(userStates[1], dto.UserStateWithID{
		StateID:   1,
		UserState: dto.MakeUserState(&leaves[1].UserState),
	})
	s.Equal(userStates[2], dto.UserStateWithID{
		StateID:   2,
		UserState: dto.MakeUserState(&leaves[2].UserState),
	})
}

func TestGetUserStatesTestSuite(t *testing.T) {
	suite.Run(t, new(GetUserStatesTestSuite))
}
