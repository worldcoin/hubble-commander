package postgres

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

//const clonedDBSuffix = "_clone"

type DBTestSuite struct {
	*require.Assertions
	suite.Suite
	db     *Database
	config *config.PostgresConfig
}

func (s *DBTestSuite) SetupSuite() {
	s.Assertions = require.New(s.T())
}

func (s *DBTestSuite) SetupTest() {
	cfg := config.GetTestConfig().Postgres
	err := RecreateDatabase(cfg)
	s.NoError(err)

	db, err := NewDatabase(cfg)
	s.NoError(err)
	s.db = db
	s.config = cfg
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

	checkBatch(s.T(), s.db, 0)

	s.NoError(migrator.Down())

	res := make([]models.Batch, 0)
	err = s.db.Query(
		sq.Select("*").From("batch"),
	).Into(&res)
	s.Error(err)
}

func (s *DBTestSuite) TestClone() {
	migrator, err := GetMigrator(s.config)
	s.NoError(err)

	s.NoError(migrator.Up())

	addBatch(s.T(), s.db)

	clonedDB, err := s.db.Clone(s.config)
	s.NoError(err)

	checkBatch(s.T(), clonedDB, 1)
	checkBatch(s.T(), s.db, 1)
}

func TestDbTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
