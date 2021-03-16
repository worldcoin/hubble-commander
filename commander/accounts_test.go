package commander

import (
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/db"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils"
	"github.com/Worldcoin/hubble-commander/testutils/deployer"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AccountsTestSuite struct {
	*require.Assertions
	suite.Suite
	storage         *st.Storage
	sim             *simulator.Simulator
	accountRegistry *accountregistry.AccountRegistry
	client          *eth.Client
}

func (s *AccountsTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *AccountsTestSuite) SetupTest() {
	testDB, err := db.NewTestDB()
	s.NoError(err)
	storage := st.NewTestStorage(testDB.DB)
	s.storage = storage
	sim, err := simulator.NewAutominingSimulator()
	s.NoError(err)
	s.sim = sim

	contracts, err := deployer.DeployRollup(sim)
	s.NoError(err)
	s.NoError(err)

	s.accountRegistry = contracts.AccountRegistry
	s.client = eth.NewTestClient(sim.Account, contracts.Rollup, contracts.AccountRegistry)
}

func (s *AccountsTestSuite) TearDownTest() {
	s.sim.Close()
}

func (s *AccountsTestSuite) Test_WatchAccounts_PreviousAccounts() {
	publicKey := models.PublicKey{2, 3, 4}
	_, err := s.accountRegistry.Register(s.sim.Account, publicKey.IntArray())
	s.NoError(err)

	go func() {
		err := WatchAccounts(s.storage, s.client)
		s.NoError(err)
	}()

	var accounts []models.Account
	testutils.WaitToPass(func() bool {
		accounts, err = s.storage.GetAccounts(&publicKey)
		s.NoError(err)
		return len(accounts) > 0
	})

	s.Len(accounts, 1)
}

func (s *AccountsTestSuite) Test_WatchAccounts_NewAccounts() {
	go func() {
		err := WatchAccounts(s.storage, s.client)
		s.NoError(err)
	}()

	time.Sleep(10 * time.Millisecond)

	publicKey := models.PublicKey{2, 3, 4}
	_, err := s.accountRegistry.Register(s.sim.Account, publicKey.IntArray())
	s.NoError(err)

	var accounts []models.Account
	testutils.WaitToPass(func() bool {
		accounts, err = s.storage.GetAccounts(&publicKey)
		s.NoError(err)
		return len(accounts) > 0
	})

	s.Len(accounts, 1)
}

func TestAccountsTestSuite(t *testing.T) {
	suite.Run(t, new(AccountsTestSuite))
}
