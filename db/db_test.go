package db

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type DbTestSuite struct {
	*require.Assertions
	suite.Suite
	db     *sqlx.DB
	config *config.Config
}

func (s *DbTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *DbTestSuite) SetupTest() {
	cfg := config.GetTestConfig()
	err := recreateDatabase(&cfg)
	s.NoError(err)

	db, err := GetDB(&cfg)
	s.NoError(err)
	s.db = db
	s.config = &cfg
}

func (s *DbTestSuite) TearDownTest() {
	err := s.db.Close()
	s.NoError(err)
}

func (s *DbTestSuite) TestGetDB() {
	s.NoError(s.db.Ping())
}

func (s *DbTestSuite) TestMigrations() {
	migrator, err := GetMigrator(s.config)
	s.NoError(err)

	s.NoError(migrator.Up())
	_, err = sq.Select("*").From("transaction").
		RunWith(s.db).Query()
	s.NoError(err)

	s.NoError(migrator.Down())
	_, err = sq.Select("*").From("transaction").
		RunWith(s.db).Query()
	s.Error(err)
}

func TestDbTestSuite(t *testing.T) {
	suite.Run(t, new(DbTestSuite))
}
