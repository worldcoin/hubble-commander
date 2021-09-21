package executor

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
)

type TestSuiteWithRollupContext struct {
	TestSuiteWithExecutionContext
	rollupCtx *RollupContext
}

func (s *TestSuiteWithRollupContext) SetupTest(batchType batchtype.BatchType) {
	s.TestSuiteWithExecutionContext.SetupTest()
	s.rollupCtx = NewTestRollupContext(s.executionCtx, batchType)
}

func (s *TestSuiteWithRollupContext) SetupTestWithConfig(batchType batchtype.BatchType, cfg config.RollupConfig) {
	s.TestSuiteWithExecutionContext.SetupTestWithConfig(cfg)
	s.rollupCtx = NewTestRollupContext(s.executionCtx, batchType)
}
