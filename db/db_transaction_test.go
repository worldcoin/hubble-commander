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

func (s *TxControllerTestSuite) Test_Rollback_ReturnsCauseError() {
	tx, _, err := s.db.DB.BeginTransaction()
	s.NoError(err)

	cause := errors.New("cause of rollback")

	err = tx.Rollback(cause)
	s.Equal(cause, err)
}

func TestTxControllerTestSuite(t *testing.T) {
	suite.Run(t, new(TxControllerTestSuite))
}
