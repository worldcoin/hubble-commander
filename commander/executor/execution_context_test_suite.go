package executor

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testSuiteWithExecutionContext struct {
	*require.Assertions
	suite.Suite
	storage      *st.TestStorage
	cfg          *config.RollupConfig
	client       *eth.TestClient
	executionCtx *ExecutionContext
}

func (s *testSuiteWithExecutionContext) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *testSuiteWithExecutionContext) SetupTest() {
	s.SetupTestWithConfig(&config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
	})
}

func (s *testSuiteWithExecutionContext) SetupTestWithConfig(cfg *config.RollupConfig) {
	var err error
	s.storage, err = st.NewTestStorage()
	s.NoError(err)

	s.cfg = cfg

	s.client, err = eth.NewTestClient()
	s.NoError(err)

	s.executionCtx = NewTestExecutionContext(s.storage.Storage, s.client.Client, s.cfg)
}

func (s *testSuiteWithExecutionContext) TearDownTest() {
	s.client.Close()
	err := s.storage.Teardown()
	s.NoError(err)
}
