package db

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TxControllerTestSuite struct {
	*require.Assertions
	suite.Suite
}

func (s *TxControllerTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

type MockController struct {
	rollbackFails bool
}

func (m MockController) Rollback() error {
	if m.rollbackFails {
		return errors.New("rollback error #1")
	}
	return nil
}

func (m MockController) Commit() error {
	return nil
}

func (s *TxControllerTestSuite) TestRollback_DoesNotModifyCauseError() {
	tx := TxController{
		tx:       MockController{},
		isLocked: false,
	}

	cause := errors.New("cause of rollback")
	tx.Rollback(&cause)
	s.Equal("cause of rollback", cause.Error())
}

func (s *TxControllerTestSuite) TestRollback_WrapsCauseErrorWithRollbackError() {
	tx := TxController{
		tx:       MockController{rollbackFails: true},
		isLocked: false,
	}

	err := errors.New("cause of rollback")
	tx.Rollback(&err)
	s.Equal("rollback caused by: cause of rollback, failed with: rollback error #1", err.Error())
	s.Equal("cause of rollback", errors.Unwrap(err).Error())
}

func TestTxControllerTestSuite(t *testing.T) {
	suite.Run(t, new(TxControllerTestSuite))
}
