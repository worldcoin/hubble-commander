package executor

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/stretchr/testify/require"
)

type TestSuiteWithDisputeContext struct {
	TestSuiteWithExecutionContext
	disputeCtx *DisputeContext
}

func (s *TestSuiteWithDisputeContext) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TestSuiteWithDisputeContext) SetupTest() {
	s.TestSuiteWithExecutionContext.SetupTestWithConfig(config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    1,
		DisableSignatures:      false,
	})

	s.disputeCtx = NewDisputeContext(s.storage.Storage, s.client.Client)
}
