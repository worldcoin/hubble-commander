package executor

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

type TestSuiteWithRollupContext struct {
	TestSuiteWithExecutionContext
	rollupCtx *RollupContext
}

func (s *TestSuiteWithRollupContext) SetupTest(batchType txtype.TransactionType) {
	s.TestSuiteWithExecutionContext.SetupTest()
	s.rollupCtx = NewTestRollupContext(s.executionCtx, batchType)
}

func (s *TestSuiteWithRollupContext) SetupTestWithConfig(batchType txtype.TransactionType, cfg config.RollupConfig) {
	s.TestSuiteWithExecutionContext.SetupTestWithConfig(cfg)
	s.rollupCtx = NewTestRollupContext(s.executionCtx, batchType)
}
