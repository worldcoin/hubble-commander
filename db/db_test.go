package db

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DBTestSuite struct {
	*require.Assertions
	suite.Suite
	db     *Database
	config *config.Config
}

func (s *DBTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *DBTestSuite) SetupTest() {
	cfg := config.GetTestConfig()
	err := recreateDatabase(&cfg)
	s.NoError(err)

	db, err := GetDB(&cfg)
	s.NoError(err)
	s.db = db
	s.config = &cfg
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

	// nolint
	_, err = sq.Select("*").From("transaction").
		RunWith(s.db).Query()
	s.NoError(err)

	s.NoError(migrator.Down())

	// nolint
	_, err = sq.Select("*").From("transaction").
		RunWith(s.db).Query()
	s.Error(err)
}

func TestDbTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
