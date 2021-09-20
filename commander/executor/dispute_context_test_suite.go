package executor

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/stretchr/testify/require"
)

type TestSuiteWithDisputeContext struct {
	TestSuiteWithRollupContext
	disputeCtx *DisputeContext
}

func (s *TestSuiteWithDisputeContext) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TestSuiteWithDisputeContext) SetupTest(batchType batchtype.BatchType) {
	s.TestSuiteWithRollupContext.SetupTestWithConfig(batchType, config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
		MinTxsPerCommitment:    1,
		MaxTxsPerCommitment:    1,
		DisableSignatures:      false,
	})

	s.disputeCtx = NewDisputeContext(s.executionCtx.storage, s.executionCtx.client)
}
