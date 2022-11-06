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
	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/utils/consts"
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

func contextWithAuthKey() context.Context {
	return context.WithValue(context.Background(), rpc.AuthKey, authKeyValue)
}

func (s *MempoolTestSuite) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *MempoolTestSuite) TestGetPendingStates_EmptyMempool() {
	result, err := s.adminAPI.GetPendingStates(contextWithAuthKey(), 0, 1000)
	s.NoError(err)
	s.Len(result, 0)
}

// TODO: this is a good spot for a proptest!

func (s *MempoolTestSuite) TestRecomputeState() {
	// I. Setup: create some accounts

	firstStateID, _ := s.createState(1)
	secondStateID, secondWallet := s.createState(2)

	s.assertAPIBalance(firstStateID, 100)
	s.assertAPIBalance(secondStateID, 100)

	// II. Insert some mempool transactions

	s.sendC2T(firstStateID, 0, secondWallet.PublicKey())
	s.assertAPIBalance(firstStateID, 80)
	s.assertAPIBalance(secondStateID, 100)

	s.sendTransfer(firstStateID, 1, secondStateID)
	s.assertAPIBalance(firstStateID, 60)
	s.assertAPIBalance(secondStateID, 110)

	s.sendTransfer(secondStateID, 0, firstStateID)
	s.assertAPIBalance(firstStateID, 70)
	s.assertAPIBalance(secondStateID, 90)

	// III. Now that we have some mempool transactions manually ruin the pending state

	err := s.storage.UnsafeSetPendingState(firstStateID, models.MakeUint256(0), models.MakeUint256(0))
	s.NoError(err)

	s.assertAPIBalance(firstStateID, 0)

	// IV. With mutate=false the pending state should not be changed

	doNotMutate := false
	result, err := s.adminAPI.RecomputePendingState(contextWithAuthKey(), firstStateID, doNotMutate)
	s.NoError(err)
	s.Equal(&dto.RecomputePendingState{
		OldNonce:   models.MakeUint256(0),
		OldBalance: models.MakeUint256(0),
		NewNonce:   models.MakeUint256(2),
		NewBalance: models.MakeUint256(70),
	}, result)
	s.assertAPIBalance(firstStateID, 0)

	// V. With mutate=true the pending state should be fixed!

	pleaseMutate := true
	result, err = s.adminAPI.RecomputePendingState(contextWithAuthKey(), firstStateID, pleaseMutate)
	s.NoError(err)
	s.Equal(&dto.RecomputePendingState{
		OldNonce:   models.MakeUint256(0),
		OldBalance: models.MakeUint256(0),
		NewNonce:   models.MakeUint256(2),
		NewBalance: models.MakeUint256(70),
	}, result)
	s.assertAPIBalance(firstStateID, 70)

	result, err = s.adminAPI.RecomputePendingState(contextWithAuthKey(), firstStateID, doNotMutate)
	s.NoError(err)
	s.Equal(&dto.RecomputePendingState{
		OldNonce:   models.MakeUint256(2),
		OldBalance: models.MakeUint256(70),
		NewNonce:   models.MakeUint256(2),
		NewBalance: models.MakeUint256(70),
	}, result)
}

func (s *MempoolTestSuite) assertPendingStates(startID, pageSize uint32, expected []dto.UserStateWithID) {
	pendingStates, err := s.adminAPI.GetPendingStates(contextWithAuthKey(), startID, pageSize)
	s.NoError(err)
	s.Equal(expected, pendingStates)
}

func (s *MempoolTestSuite) assertUserStates(pubkey *models.PublicKey, expected []dto.UserStateWithID) {
	pendingStates, err := s.hubbleAPI.GetUserStates(contextWithAuthKey(), pubkey)
	s.NoError(err)
	s.Equal(expected, pendingStates)
}

func (s *MempoolTestSuite) assertNoUserStates(pubkey *models.PublicKey) {
	pendingStates, err := s.hubbleAPI.GetUserStates(contextWithAuthKey(), pubkey)
	s.ErrorContains(err, "user states not found")
	s.Len(pendingStates, 0)
}

func (s *MempoolTestSuite) expectedStateWithID(stateID, nonce, balance uint32) dto.UserStateWithID {
	return dto.UserStateWithID{
		StateID: stateID,
		UserState: dto.UserState{
			PubKeyID: stateID,
			TokenID:  models.MakeUint256(0),
			Nonce:    models.MakeUint256(uint64(nonce)),
			Balance:  models.MakeUint256(uint64(balance)),
		},
	}
}

func (s *MempoolTestSuite) TestGetPendingStates() {
	firstStateID, _ := s.createState(1)
	secondStateID, secondWallet := s.createState(2)

	s.assertAPIBalance(firstStateID, 100)
	s.assertAPIBalance(secondStateID, 100)

	s.assertPendingStates(0, 1000, []dto.UserStateWithID{})

	s.sendC2T(firstStateID, 0, secondWallet.PublicKey())

	s.assertPendingStates(0, 1000, []dto.UserStateWithID{
		s.expectedStateWithID(firstStateID, 1, 80),
	})
	s.assertPendingStates(2, 1000, []dto.UserStateWithID{})

	s.sendTransfer(firstStateID, 1, secondStateID)
	s.assertPendingStates(0, 1000, []dto.UserStateWithID{
		s.expectedStateWithID(firstStateID, 2, 60),
		s.expectedStateWithID(secondStateID, 0, 110),
	})
	s.assertPendingStates(1, 1, []dto.UserStateWithID{
		s.expectedStateWithID(firstStateID, 2, 60),
	})
	s.assertPendingStates(2, 1, []dto.UserStateWithID{
		s.expectedStateWithID(secondStateID, 0, 110),
	})
}

func (s *MempoolTestSuite) TestGetPendingC2TState() {
	firstStateID, _ := s.createState(1)

	randomWallet := s.randomWallet()
	s.sendC2T(firstStateID, 0, randomWallet.PublicKey())
	s.assertUserStates(randomWallet.PublicKey(), []dto.UserStateWithID{
		s.expectedStateWithID(consts.PendingID, 0, 10),
	})

	s.sendC2T(firstStateID, 1, randomWallet.PublicKey())
	s.assertUserStates(randomWallet.PublicKey(), []dto.UserStateWithID{
		s.expectedStateWithID(consts.PendingID, 0, 20),
	})
}

func (s *MempoolTestSuite) randomPublicKey() *models.PublicKey {
	domain := bls.Domain{1, 2, 3, 4}
	wallet, err := bls.NewRandomWallet(domain)
	s.NoError(err)

	return wallet.PublicKey()
}

func (s *MempoolTestSuite) rawInsert(tx models.GenericTransaction) {
	pendingTx := stored.NewPendingTx(tx)
	err := s.storage.UnsafeInsertPendingTxSkipValidation(pendingTx)
	s.NoError(err)
}

func (s *MempoolTestSuite) TestGetPendingC2TStateFromMigratedAccount() {
	firstStateID, _ := s.createState(1)
	destKey := s.randomPublicKey()
	s.assertNoUserStates(destKey)

	s.rawInsert(testutils.NewCreate2Transfer(firstStateID, nil, 0, 10, destKey))
	s.rawInsert(testutils.NewCreate2Transfer(firstStateID, nil, 1, 20, destKey))

	err := s.storage.MigratePubKeyPendingState()
	s.NoError(err)

	// the raw insert skipped the codepath which updates the pending state so we do it manually
	err = s.storage.UnsafeSetPendingState(firstStateID, models.MakeUint256(2), models.MakeUint256(50))
	s.NoError(err)

	// confirm the migration happened correctly
	s.assertUserStates(destKey, []dto.UserStateWithID{
		s.expectedStateWithID(consts.PendingID, 0, 30),
	})

	s.sendC2T(firstStateID, 2, destKey)
	s.assertUserStates(destKey, []dto.UserStateWithID{
		s.expectedStateWithID(consts.PendingID, 0, 40),
	})
}

func (s *MempoolTestSuite) bothPubkeyMethodsEqual(startPrefix []byte, pageSize uint32) []dto.PubkeyBalance {
	balancesGet, err := s.adminAPI.GetPendingPubkeyBalances(contextWithAuthKey(), startPrefix, pageSize)
	s.NoError(err)

	balancesRecompute, err := s.adminAPI.RecomputePubkeyBalances(contextWithAuthKey(), startPrefix, pageSize)
	s.NoError(err)

	s.Equal(balancesGet, balancesRecompute)
	return balancesGet
}

func (s *MempoolTestSuite) TestGetPendingPubkeyBalances() {
	firstStateID, _ := s.createState(1)

	firstDestKey := s.randomPublicKey()
	s.sendC2T(firstStateID, 0, firstDestKey)
	s.sendC2T(firstStateID, 1, firstDestKey)

	secondDestKey := s.randomPublicKey()
	s.sendC2T(firstStateID, 2, secondDestKey)

	thirdDestKey := s.randomPublicKey()
	s.sendC2T(firstStateID, 3, thirdDestKey)

	_, err := s.adminAPI.GetPendingPubkeyBalances(context.Background(), []byte{}, 0)
	s.ErrorContains(err, "missing authentication key")

	allbalances := s.bothPubkeyMethodsEqual([]byte{}, 0)
	s.Len(allbalances, 3)
	s.Contains(allbalances, dto.PubkeyBalance{
		PubKey:  *firstDestKey,
		Balance: models.MakeUint256(20),
	})
	s.Contains(allbalances, dto.PubkeyBalance{
		PubKey:  *secondDestKey,
		Balance: models.MakeUint256(10),
	})
	s.Contains(allbalances, dto.PubkeyBalance{
		PubKey:  *thirdDestKey,
		Balance: models.MakeUint256(10),
	})

	// test pagination
	balances := s.bothPubkeyMethodsEqual([]byte{}, 1)
	s.Len(balances, 1)
	s.Equal(allbalances[0], balances[0])

	balances = s.bothPubkeyMethodsEqual(allbalances[1].PubKey[:], 1)
	s.Len(balances, 1)
	s.Equal(allbalances[1], balances[0])
}

func (s *MempoolTestSuite) assertAPIBalance(stateID, balance uint32) {
	userState, err := s.hubbleAPI.GetUserState(context.Background(), stateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(uint64(balance)), userState.Balance)
}

func (s *MempoolTestSuite) randomWallet() *bls.Wallet {
	domain, err := s.client.GetDomain()
	s.NoError(err)

	wallet, err := bls.NewRandomWallet(*domain)
	s.NoError(err)

	return wallet
}

func (s *MempoolTestSuite) createState(stateID uint32) (uint32, *bls.Wallet) {
	wallet := s.randomWallet()

	account := models.AccountLeaf{
		PubKeyID:  stateID,
		PublicKey: *wallet.PublicKey(),
	}
	err := s.storage.AccountTree.SetSingle(&account)
	s.NoError(err)

	_, err = s.storage.StateTree.Set(
		stateID,
		&models.UserState{
			PubKeyID: account.PubKeyID,
			TokenID:  models.MakeUint256(0),
			Balance:  models.MakeUint256(100),
			Nonce:    models.MakeUint256(0),
		},
	)
	s.NoError(err)

	return stateID, wallet
}

func (s *MempoolTestSuite) sendC2T(from, nonce uint32, to *models.PublicKey) {
	c2t := dto.Create2Transfer{
		FromStateID: ref.Uint32(from),
		ToPublicKey: to,
		Amount:      models.NewUint256(10),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(uint64(nonce)),
		Signature:   &models.Signature{},
	}

	hash, err := s.hubbleAPI.SendTransaction(context.Background(), dto.MakeTransaction(c2t))
	s.NoError(err)
	s.NotNil(hash)
}

func (s *MempoolTestSuite) sendTransfer(from, nonce, to uint32) {
	transfer := dto.Transfer{
		FromStateID: ref.Uint32(from),
		ToStateID:   ref.Uint32(to),
		Amount:      models.NewUint256(10),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(uint64(nonce)),
		Signature:   &models.Signature{},
	}

	hash, err := s.hubbleAPI.SendTransaction(context.Background(), dto.MakeTransaction(transfer))
	s.NoError(err)
	s.NotNil(hash)
}

func TestMempoolTestSuite(t *testing.T) {
	suite.Run(t, new(MempoolTestSuite))
}
