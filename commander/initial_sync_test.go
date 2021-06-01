package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type InitialSyncTestSuite struct {
	*require.Assertions
	suite.Suite
	client   *eth.TestClient
	cmd      *Commander
	tree     *st.StateTree
	teardown func() error
}

func (s *InitialSyncTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *InitialSyncTestSuite) SetupTest() {
	var err error
	s.client, err = eth.NewTestClient()
	s.NoError(err)

	cfg := config.GetTestConfig()
	cfg.Rollup.SyncSize = 1
	cfg.Rollup.MinCommitmentsPerBatch = 1
	cfg.Rollup.MaxCommitmentsPerBatch = 1
	cfg.Rollup.TxsPerCommitment = 1
	cfg.Rollup.FeeReceiverPubKeyID = 2
	s.setupDB(cfg)
}

func (s *InitialSyncTestSuite) setupDB(cfg *config.Config) {
	storage, err := st.NewTestStorageWithBadger()
	s.NoError(err)
	s.teardown = storage.Teardown
	s.tree = st.NewStateTree(storage.Storage)

	s.cmd = &Commander{
		cfg:     cfg,
		storage: storage.Storage,
		client:  s.client.Client,
	}

	s.seedDB()
}

func (s *InitialSyncTestSuite) TearDownTest() {
	s.client.Close()
	err := s.teardown()
	s.NoError(err)
}

func (s *InitialSyncTestSuite) TestInitialSync() {
	number, err := s.client.GetLatestBlockNumber()
	s.NoError(err)
	s.cmd.storage.SetLatestBlockNumber(*number + 3)

	accounts := []models.Account{
		{PublicKey: models.PublicKey{1, 1, 1}},
		{PublicKey: models.PublicKey{2, 2, 2}},
	}

	registrations, unsubscribe, err := s.client.WatchRegistrations(&bind.WatchOpts{})
	s.NoError(err)
	defer unsubscribe()
	for i := range accounts {
		var senderPubKeyID *uint32
		senderPubKeyID, err = s.client.RegisterAccount(&accounts[i].PublicKey, registrations)
		s.NoError(err)
		s.Equal(uint32(i), *senderPubKeyID)
		accounts[i].PubKeyID = *senderPubKeyID
	}

	s.addBatch()
	// Recreate database
	err = s.teardown()
	s.NoError(err)
	s.setupDB(s.cmd.cfg)

	//err = s.cmd.InitialSync()
	//s.NoError(err)
	for i := range accounts {
		var userAccounts []models.Account
		userAccounts, err = s.cmd.storage.GetAccounts(&accounts[i].PublicKey)
		s.NoError(err)
		s.Len(userAccounts, 1)
		s.Contains(accounts, userAccounts[0])
	}

	batches, err := s.cmd.storage.GetBatchesInRange(nil, nil)
	s.NoError(err)
	s.Len(batches, 1)
}

func (s *InitialSyncTestSuite) addBatch() {
	tx := models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 2,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(0),
			Nonce:       models.MakeUint256(0),
			Signature:   *mockSignature(s.T()),
		},
		ToStateID: 3,
	}
	err := s.cmd.storage.AddTransfer(&tx)
	s.NoError(err)

	commitments, err := createTransferCommitments([]models.Transfer{tx}, s.cmd.storage, s.cmd.cfg.Rollup, testDomain)
	s.NoError(err)
	s.Len(commitments, 1)

	err = submitBatch(txtype.Transfer, commitments, s.cmd.storage, s.client.Client, s.cmd.cfg.Rollup)
	s.NoError(err)
}

func (s *InitialSyncTestSuite) seedDB() {
	err := s.cmd.storage.AddAccountIfNotExists(&models.Account{
		PubKeyID:  2,
		PublicKey: models.PublicKey{1, 2, 3},
	})
	s.NoError(err)

	err = s.cmd.storage.AddAccountIfNotExists(&models.Account{
		PubKeyID:  3,
		PublicKey: models.PublicKey{2, 3, 4},
	})
	s.NoError(err)

	err = s.tree.Set(2, &models.UserState{
		PubKeyID:   2,
		TokenIndex: models.MakeUint256(0),
		Balance:    models.MakeUint256(1000),
		Nonce:      models.MakeUint256(0),
	})
	s.NoError(err)

	err = s.tree.Set(3, &models.UserState{
		PubKeyID:   3,
		TokenIndex: models.MakeUint256(0),
		Balance:    models.MakeUint256(0),
		Nonce:      models.MakeUint256(0),
	})
	s.NoError(err)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(InitialSyncTestSuite))
}
