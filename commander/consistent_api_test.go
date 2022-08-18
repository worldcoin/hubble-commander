package commander

// this is a test of the api more than anything, but putting it here because we touch
// the commander and this breaks any potential import loops

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type APIConsistencySuite struct {
	*require.Assertions
	suite.Suite
}

func (s *APIConsistencySuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

//nolint:funlen
func (s *APIConsistencySuite) TestGetUserStatesPendingToBatched() {
	testStorage, err := st.NewTestStorage()
	s.NoError(err)

	client, err := eth.NewTestClient()
	s.NoError(err)

	// we want to check that the state updates before the chain mines our batches
	client.StopAutomine()

	// we were updating this manually but that is no longer necessary
	acct := client.Simulator.GetAccount()
	acct.Nonce = nil

	theapi := api.NewTestAPI(
		testStorage.Storage,
		client.Client,
	)

	commander := &Commander{
		cfg: &config.Config{
			Rollup: &config.RollupConfig{
				MinTxsPerCommitment:    1,
				MaxTxsPerCommitment:    1,
				MinCommitmentsPerBatch: 1,
				MaxCommitmentsPerBatch: 32,
			},
		},
		storage: testStorage.Storage,
		client:  client.Client,
		metrics: metrics.NewCommanderMetrics(),
	}

	feeReceiver := &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	}
	_, err = testStorage.Storage.StateTree.Set(0, feeReceiver)
	s.NoError(err)

	domain, err := client.Client.GetDomain()
	s.NoError(err)

	wallet, err := bls.NewRandomWallet(*domain)
	s.NoError(err)

	receiverWallet, err := bls.NewRandomWallet(*domain)
	s.NoError(err)

	err = testStorage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  1,
		PublicKey: *wallet.PublicKey(),
	})
	s.NoError(err)

	userState := &models.UserState{
		PubKeyID: 1,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	}
	_, err = testStorage.StateTree.Set(1, userState)
	s.NoError(err)

	create2TransferWithoutSignature := dto.Create2Transfer{
		FromStateID: ref.Uint32(1),
		Amount:      models.NewUint256(20),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		ToPublicKey: receiverWallet.PublicKey(),
	}
	create2Transfer, err := api.SignCreate2Transfer(wallet, create2TransferWithoutSignature)
	s.NoError(err)

	hash, err := theapi.SendTransaction(context.Background(), dto.MakeTransaction(*create2Transfer))
	s.NoError(err)
	s.NotNil(hash)

	fetchedStates, err := theapi.GetUserStates(context.Background(), receiverWallet.PublicKey())
	s.NoError(err)
	s.Len(fetchedStates, 1)
	s.Equal(dto.UserStateWithID{
		StateID: ^uint32(0),
		UserState: dto.UserState{
			PubKeyID: ^uint32(0),
			TokenID:  models.MakeUint256(0),
			Balance:  models.MakeUint256(20),
			Nonce:    models.MakeUint256(0),
		},
	}, fetchedStates[0])

	nextPubkeyID, err := testStorage.AccountTree.NextBatchAccountPubKeyID()
	s.NoError(err)
	s.Equal(uint32(0x80000000), *nextPubkeyID) // see AccountBatchOffset

	currentBatchType := batchtype.Create2Transfer
	err = commander.rollupLoopIteration(context.Background(), &currentBatchType)
	s.NoError(err)

	fetchedStates, err = theapi.GetUserStates(context.Background(), receiverWallet.PublicKey())
	s.NoError(err)
	s.Len(fetchedStates, 1)
	s.Equal(dto.UserStateWithID{
		StateID: 2,
		UserState: dto.UserState{
			PubKeyID: *nextPubkeyID,
			TokenID:  models.MakeUint256(0),
			Balance:  models.MakeUint256(20),
			Nonce:    models.MakeUint256(0),
		},
	}, fetchedStates[0])
}

//nolint:funlen
func (s *APIConsistencySuite) TestConsistency() {
	// a bunch of setup

	testStorage, err := st.NewTestStorage()
	s.NoError(err)

	client, err := eth.NewTestClient()
	s.NoError(err)

	// we want to check that the state updates before the chain mines our batches
	client.StopAutomine()

	// we were updating this manually but that is no longer necessary
	acct := client.Simulator.GetAccount()
	acct.Nonce = nil

	theapi := api.NewTestAPI(
		testStorage.Storage,
		client.Client,
	)

	commander := &Commander{
		cfg: &config.Config{
			Rollup: &config.RollupConfig{
				MinTxsPerCommitment:    1,
				MaxTxsPerCommitment:    1,
				MinCommitmentsPerBatch: 1,
				MaxCommitmentsPerBatch: 32,
			},
		},
		storage: testStorage.Storage,
		client:  client.Client,
		metrics: metrics.NewCommanderMetrics(),
	}

	feeReceiver := &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	}
	_, err = testStorage.Storage.StateTree.Set(0, feeReceiver)
	s.NoError(err)

	domain, err := client.Client.GetDomain()
	s.NoError(err)

	wallet, err := bls.NewRandomWallet(*domain)
	s.NoError(err)

	receiverWallet, err := bls.NewRandomWallet(*domain)
	s.NoError(err)

	err = testStorage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  123,
		PublicKey: *wallet.PublicKey(),
	})
	s.NoError(err)

	userState := &models.UserState{
		PubKeyID: 123,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(420),
		Nonce:    models.MakeUint256(0),
	}

	_, err = testStorage.StateTree.Set(1, userState)
	s.NoError(err)

	create2TransferWithoutSignature := dto.Create2Transfer{
		FromStateID: ref.Uint32(1),
		Amount:      models.NewUint256(20),
		Fee:         models.NewUint256(10),
		Nonce:       models.NewUint256(0),
		ToPublicKey: receiverWallet.PublicKey(),
	}
	create2Transfer, err := api.SignCreate2Transfer(wallet, create2TransferWithoutSignature)
	s.NoError(err)

	// submit two transaction to the api which depend on each other

	// check the nonce of the account

	fetchedState, err := theapi.GetUserState(context.Background(), 1)
	s.NoError(err)
	s.Equal(fetchedState.Nonce, models.MakeUint256(0))
	s.Equal(fetchedState.Balance, models.MakeUint256(420))

	hash, err := theapi.SendTransaction(context.Background(), dto.MakeTransaction(*create2Transfer))
	s.NoError(err)
	s.NotNil(hash)

	// run an iteration of the rollup loop here in such a way that it generates
	// a commitment

	currentBatchType := batchtype.Create2Transfer
	err = commander.rollupLoopIteration(context.Background(), &currentBatchType)
	s.NoError(err)

	fetchedState, err = theapi.GetUserState(context.Background(), 1)
	s.NoError(err)
	s.Equal(models.MakeUint256(1), fetchedState.Nonce)
	s.Equal(models.MakeUint256(390), fetchedState.Balance)
}

func TestAPIConsistencySuite(t *testing.T) {
	suite.Run(t, new(APIConsistencySuite))
}
