package postgres

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TxControllerTestSuite struct {
	*require.Assertions
	suite.Suite
	db *TestDB
}

func (s *TxControllerTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *TxControllerTestSuite) SetupTest() {
	testDB, err := NewTestDB()
	s.NoError(err)
	s.db = testDB
}

func (s *TxControllerTestSuite) TearDownTest() {
	err := s.db.Teardown()
	s.NoError(err)
}

func (s *TxControllerTestSuite) TestRollback_DoesNotModifyCauseError() {
	tx, _, err := s.db.DB.BeginTransaction()
	s.NoError(err)

	cause := errors.New("cause of rollback")
	tx.Rollback(&cause)
	s.Equal("cause of rollback", cause.Error())
}

type MockController struct{}

func (m MockController) Rollback() error {
	return errors.New("rollback error #1")
}

func (m MockController) Commit() error {
	return nil
}

func (s *TxControllerTestSuite) TestRollback_WrapsCauseErrorWithRollbackError() {
	tx := TransactionController{
		tx:       MockController{},
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
