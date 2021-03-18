package db

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DBTestSuite struct {
	*require.Assertions
	suite.Suite
	db     *Database
	config *config.DBConfig
}

func (s *DBTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *DBTestSuite) SetupTest() {
	cfg := config.GetTestConfig()
	err := recreateDatabase(&cfg.DB)
	s.NoError(err)

	db, err := NewDatabase(&cfg.DB)
	s.NoError(err)
	s.db = db
	s.config = &cfg.DB
}

func (s *DBTestSuite) TearDownTest() {
	err := s.db.Close()
	s.NoError(err)
}

func (s *DBTestSuite) TestGetDB() {
	s.NoError(s.db.Ping())
}

func (s *DBTestSuite) TestMigrations() {
	migrator, err := GetMigrator(s.config)
	s.NoError(err)

	s.NoError(migrator.Up())

	res := make([]models.Transaction, 0, 1)
	err = s.db.Query(
		sq.Select("*").From("transaction"),
	).Into(&res)
	s.NoError(err)

	s.NoError(migrator.Down())

	err = s.db.Query(
		sq.Select("*").From("transaction"),
	).Into(&res)
	s.Error(err)
}

func TestDbTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
