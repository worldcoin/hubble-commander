// we are testing a method from `api/admin` but we call methods in `api`. To break the import loop
// this file lives in api.
package api

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/api/admin"
	"github.com/Worldcoin/hubble-commander/api/rpc"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const authKeyValue = "secret key"

type MempoolTestSuite struct {
	*require.Assertions
	suite.Suite
	hubbleAPI *API
	adminAPI  *admin.API
	client    *eth.TestClient
	storage   *storage.TestStorage
}

func (s *MempoolTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *MempoolTestSuite) SetupTest() {
	var err error
	s.storage, err = storage.NewTestStorage()
	s.NoError(err)
	s.client, err = eth.NewTestClient()
	s.NoError(err)
	s.adminAPI = admin.NewTestAPI(
		&config.APIConfig{AuthenticationKey: authKeyValue},
		s.storage.Storage,
		s.client.Client,
	)
	s.hubbleAPI = NewTestAPI(
		s.storage.Storage,
		s.client.Client,
	)
}

func contextWithAuthKey(authKeyValue string) context.Context {
	return context.WithValue(context.Background(), rpc.AuthKey, authKeyValue)
}

func (s *MempoolTestSuite) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

// TODO: this is a good spot for a proptest!
//nolint:funlen
func (s *MempoolTestSuite) TestRecomputeState() {
	// I. Setup: create some accounts

	domain, err := s.client.GetDomain()
	s.NoError(err)

	firstWallet, err := bls.NewRandomWallet(*domain)
	s.NoError(err)

	firstPubkeyID := uint32(1)
	firstAccount := models.AccountLeaf{
		PubKeyID:  firstPubkeyID,
		PublicKey: *firstWallet.PublicKey(),
	}
	err = s.storage.AccountTree.SetSingle(&firstAccount)
	s.NoError(err)

	firstStateID := uint32(1)
	_, err = s.storage.StateTree.Set(
		firstStateID,
		&models.UserState{
			PubKeyID: firstAccount.PubKeyID,
			TokenID:  models.MakeUint256(0),
			Balance:  models.MakeUint256(100),
			Nonce:    models.MakeUint256(0),
		},
	)
	s.NoError(err)

	secondWallet, err := bls.NewRandomWallet(*domain)
	s.NoError(err)

	secondPubkeyID := uint32(2)
	secondAccount := models.AccountLeaf{
		PubKeyID:  secondPubkeyID,
		PublicKey: *secondWallet.PublicKey(),
	}
	err = s.storage.AccountTree.SetSingle(&secondAccount)
	s.NoError(err)

	secondStateID := uint32(2)
	_, err = s.storage.StateTree.Set(
		secondStateID,
		&models.UserState{
			PubKeyID: secondAccount.PubKeyID,
			TokenID:  models.MakeUint256(0),
			Balance:  models.MakeUint256(100),
			Nonce:    models.MakeUint256(0),
		},
	)
	s.NoError(err)

	// II. Insert some mempool transactions

	c2t := dto.Create2Transfer{
		FromStateID: ref.Uint32(firstStateID),
		ToPublicKey: secondWallet.PublicKey(),
		Amount:      models.NewUint256(10),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   &models.Signature{},
	}

	hash, err := s.hubbleAPI.SendTransaction(context.Background(), dto.MakeTransaction(c2t))
	s.NoError(err)
	s.NotNil(hash)

	transfer := dto.Transfer{
		FromStateID: ref.Uint32(firstStateID),
		ToStateID:   ref.Uint32(secondStateID),
		Amount:      models.NewUint256(10),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(1),
		Signature:   &models.Signature{},
	}

	hash, err = s.hubbleAPI.SendTransaction(context.Background(), dto.MakeTransaction(transfer))
	s.NoError(err)
	s.NotNil(hash)

	transfer = dto.Transfer{
		FromStateID: ref.Uint32(secondStateID),
		ToStateID:   ref.Uint32(firstStateID),
		Amount:      models.NewUint256(10),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		Signature:   &models.Signature{},
	}

	hash, err = s.hubbleAPI.SendTransaction(context.Background(), dto.MakeTransaction(transfer))
	s.NoError(err)
	s.NotNil(hash)

	// III. Now that we have some mempool transactions, manually ruin the pending state

	err = s.storage.UnsafeSetPendingState(firstStateID, models.MakeUint256(0), models.MakeUint256(0))
	s.NoError(err)

	// IV. With mutate=false the pending state should not be changed

	doNotMutate := false
	result, err := s.adminAPI.RecomputePendingState(contextWithAuthKey(authKeyValue), firstStateID, doNotMutate)
	s.NoError(err)
	s.Equal(&dto.RecomputePendingState{
		OldNonce:   models.MakeUint256(0),
		OldBalance: models.MakeUint256(0),
		NewNonce:   models.MakeUint256(2),
		NewBalance: models.MakeUint256(90),
	}, result)

	// V. With mutate=true the pending state should be fixed!

	pleaseMutate := true
	result, err = s.adminAPI.RecomputePendingState(contextWithAuthKey(authKeyValue), firstStateID, pleaseMutate)
	s.NoError(err)
	s.Equal(&dto.RecomputePendingState{
		OldNonce:   models.MakeUint256(0),
		OldBalance: models.MakeUint256(0),
		NewNonce:   models.MakeUint256(2),
		NewBalance: models.MakeUint256(90),
	}, result)

	result, err = s.adminAPI.RecomputePendingState(contextWithAuthKey(authKeyValue), firstStateID, doNotMutate)
	s.NoError(err)
	s.Equal(&dto.RecomputePendingState{
		OldNonce:   models.MakeUint256(2),
		OldBalance: models.MakeUint256(90),
		NewNonce:   models.MakeUint256(2),
		NewBalance: models.MakeUint256(90),
	}, result)
}

func TestMempoolTestSuite(t *testing.T) {
	suite.Run(t, new(MempoolTestSuite))
}
