package executor

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
)

type testSuiteWithRollupContext struct {
	testSuiteWithExecutionContext
	rollupCtx *RollupContext
}

func (s *testSuiteWithRollupContext) SetupTest(batchType batchtype.BatchType) {
	s.testSuiteWithExecutionContext.SetupTest()
	s.rollupCtx = NewTestRollupContext(s.executionCtx, batchType)
}

func (s *testSuiteWithRollupContext) SetupTestWithConfig(batchType batchtype.BatchType, cfg config.RollupConfig) {
	s.testSuiteWithExecutionContext.SetupTestWithConfig(cfg)
	s.rollupCtx = NewTestRollupContext(s.executionCtx, batchType)
}
