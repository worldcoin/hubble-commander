package commander

// this is a test of the api, but we touch the commander so it must be put here to break
// the import loop

import (
	"context"
	"testing"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/api/admin"
	"github.com/Worldcoin/hubble-commander/api/rpc"
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const API_KEY = "secret"

func contextWithAuthKey() context.Context {
	return context.WithValue(context.Background(), rpc.AuthKey, API_KEY)
}

type RevertBatchTestSuite struct {
	*require.Assertions
	suite.Suite

	storage *st.TestStorage
	client  *eth.TestClient

	commander *Commander

	api      *api.API
	adminAPI *admin.API
}

func (s *RevertBatchTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *RevertBatchTestSuite) SetupTest() {
	var err error

	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.client, err = eth.NewTestClient()
	s.NoError(err)

	s.client.StopAutomine()

	// we were updating this manually but that is no longer necessary
	acct := s.client.Simulator.GetAccount()
	acct.Nonce = nil

	s.commander = &Commander{
		cfg: &config.Config{
			Rollup: &config.RollupConfig{
				MinTxsPerCommitment:    1,
				MaxTxsPerCommitment:    1,
				MinCommitmentsPerBatch: 1,
				MaxCommitmentsPerBatch: 32,
			},
		},
		storage: s.storage.Storage,
		client:  s.client.Client,
		metrics: metrics.NewCommanderMetrics(),
	}

	s.api = api.NewTestAPI(
		s.storage.Storage,
		s.client.Client,
	)

	s.adminAPI = admin.NewTestAPI(
		&config.APIConfig{AuthenticationKey: API_KEY},
		s.storage.Storage,
		s.client.Client,
	)
}

func (s *RevertBatchTestSuite) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}

func (s *RevertBatchTestSuite) createFeeReceiver() {
	feeReceiver := &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	}
	_, err := s.storage.Storage.StateTree.Set(0, feeReceiver)
	s.NoError(err)
}

func (s *RevertBatchTestSuite) createAccount(stateID, balance uint32) *bls.Wallet {
	domain, err := s.client.Client.GetDomain()
	s.NoError(err)

	wallet, err := bls.NewRandomWallet(*domain)
	s.NoError(err)

	err = s.storage.AccountTree.SetSingle(&models.AccountLeaf{
		PubKeyID:  stateID,
		PublicKey: *wallet.PublicKey(),
	})
	s.NoError(err)

	userState := &models.UserState{
		PubKeyID: stateID,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(uint64(balance)),
		Nonce:    models.MakeUint256(0),
	}
	_, err = s.storage.StateTree.Set(stateID, userState)
	s.NoError(err)

	return wallet
}

func (s *RevertBatchTestSuite) sendTransfer(from, to, nonce, amount uint32) {
	transfer := dto.Transfer{
		FromStateID: ref.Uint32(from),
		ToStateID:   ref.Uint32(to),
		Amount:      models.NewUint256(uint64(amount)),
		Fee:         models.NewUint256(1),
		Nonce:       models.NewUint256(uint64(nonce)),
		Signature:   &models.Signature{},
	}

	hash, err := s.api.SendTransaction(context.Background(), dto.MakeTransaction(transfer))
	s.NoError(err)
	s.NotNil(hash)
}

func (s *RevertBatchTestSuite) assertAPIState(stateID, nonce, balance uint32) {
	userState, err := s.api.GetUserState(context.Background(), stateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(uint64(balance)), userState.Balance)
	s.Equal(models.MakeUint256(uint64(nonce)), userState.Nonce)
}

func (s *RevertBatchTestSuite) assertStorageState(stateID, nonce, balance uint32) {
	leaf, err := s.storage.StateTree.Leaf(stateID)
	s.NoError(err)
	s.Equal(models.MakeUint256(uint64(balance)), leaf.UserState.Balance)
	s.Equal(models.MakeUint256(uint64(nonce)), leaf.UserState.Nonce)
}

func (s *RevertBatchTestSuite) assertAPIAndStorageMatch(stateID, nonce, balance uint32) {
	s.assertStorageState(stateID, nonce, balance)
	s.assertAPIState(stateID, nonce, balance)
}

func (s *RevertBatchTestSuite) singleRollupIteration(batchtype batchtype.BatchType) {
	err := s.commander.rollupLoopIteration(context.Background(), &batchtype)
	s.NoError(err)
}

func (s *RevertBatchTestSuite) assertBatchCount(count uint32) {
	batches, err := s.api.GetBatches(models.NewUint256(0), models.NewUint256(1000))
	s.NoError(err)
	s.Len(batches, int(count))
}

func (s *RevertBatchTestSuite) assertExistsOneBatch() *dto.Batch {
	batches, err := s.api.GetBatches(models.NewUint256(0), models.NewUint256(1000))
	s.NoError(err)
	s.Len(batches, 1)

	return &batches[0]
}

func (s *RevertBatchTestSuite) assertMempoolSize(txType txtype.TransactionType, count uint32) {
	size, err := s.storage.CountPendingTxsOfType(txType)
	s.NoError(err)
	s.Equal(count, size)
}

func (s *RevertBatchTestSuite) Test_RevertSingleBatchNoMempoolTxs() {
	// an easy case: a batch with a transfer in it and no more mempool txs
	//               put the single tx back into the mempool, that's all

	s.createFeeReceiver()
	s.createAccount(1, 400)
	s.createAccount(2, 400)

	s.sendTransfer(1, 2, 0, 9) // includes a fee of 1
	apiStateHasNotChanged := func() {
		s.assertAPIState(1, 1, 390)
		s.assertAPIState(2, 0, 409)
	}
	apiStateHasNotChanged()
	s.assertMempoolSize(txtype.Transfer, 1)
	unbatchedStorageState := func() {
		s.assertStorageState(1, 0, 400)
		s.assertStorageState(2, 0, 400)
	}
	unbatchedStorageState()

	s.singleRollupIteration(batchtype.Transfer)
	firstBatch := s.assertExistsOneBatch()
	s.assertMempoolSize(txtype.Transfer, 0)
	batchedStorageState := func() {
		s.assertStorageState(1, 1, 390)
		s.assertStorageState(2, 0, 409)
	}
	batchedStorageState()

	err := s.adminAPI.RevertBatches(contextWithAuthKey(), firstBatch.ID)
	s.NoError(err)
	s.assertMempoolSize(txtype.Transfer, 1)
	apiStateHasNotChanged()
	unbatchedStorageState()

	s.singleRollupIteration(batchtype.Transfer)
	s.assertExistsOneBatch()
	s.assertMempoolSize(txtype.Transfer, 0)
	apiStateHasNotChanged()
	batchedStorageState()
}

func (s *RevertBatchTestSuite) Test_RevertMultipleBatches() {
	// a hard case: revert multiple batches at the same time,
	//              additional transactions have arrived since the batches were created

	s.createFeeReceiver()
	s.createAccount(1, 400)
	s.createAccount(2, 0)

	// first batch

	s.sendTransfer(1, 2, 0, 10) // sendTransfer adds a fee of 1, we spend 11
	s.sendTransfer(2, 1, 0, 4)

	s.singleRollupIteration(batchtype.Transfer)
	firstBatch := s.assertExistsOneBatch()

	// second batch

	s.sendTransfer(2, 1, 1, 4)

	s.singleRollupIteration(batchtype.Transfer)
	s.assertBatchCount(2)
	s.assertAPIAndStorageMatch(1, 1, 397)
	s.assertAPIAndStorageMatch(2, 2, 0)

	// some additional, unbatched transfers

	s.sendTransfer(1, 2, 1, 99)
	s.sendTransfer(2, 1, 2, 98)

	// revert it all

	err := s.adminAPI.RevertBatches(contextWithAuthKey(), firstBatch.ID)
	s.NoError(err)
	s.assertAPIState(1, 2, 297+98)
	s.assertAPIState(2, 3, 0)
	s.assertStorageState(1, 0, 400)
	s.assertStorageState(2, 0, 0)
	s.assertMempoolSize(txtype.Transfer, 5)

	// now batch it all up again

	s.singleRollupIteration(batchtype.Transfer)
	s.assertBatchCount(1)
	s.assertMempoolSize(txtype.Transfer, 0)
	s.assertAPIAndStorageMatch(1, 2, 297+98)
	s.assertAPIAndStorageMatch(2, 3, 0)
}

func TestRevertBatchTestSuite(t *testing.T) {
	suite.Run(t, new(RevertBatchTestSuite))
}
