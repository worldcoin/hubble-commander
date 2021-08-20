package postgres

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

	checkTransfer(s.T(), s.db, 0)

	s.NoError(migrator.Down())

	res := make([]models.Batch, 0)
	err = s.db.Query(
		sq.Select("*").From("transfers"),
	).Into(&res)
	s.Error(err)
}

func (s *DBTestSuite) TestClone() {
	migrator, err := GetMigrator(s.config)
	s.NoError(err)

	s.NoError(migrator.Up())

	addTransfer(s.T(), s.db)

	clonedDB, err := s.db.Clone(s.config)
	s.NoError(err)

	checkTransfer(s.T(), clonedDB, 1)
	checkTransfer(s.T(), s.db, 1)
}

func (s *DBTestSuite) TestClone_DoesNotChangeReceiverDB() {
	initialName := s.getDBName(s.db)

	clonedDB, err := s.db.Clone(s.config)
	s.NoError(err)

	s.Equal(initialName, s.getDBName(s.db))
	s.NotEqual(initialName, s.getDBName(clonedDB))
}

func (s *DBTestSuite) getDBName(db *Database) string {
	dbName := make([]string, 0, 1)
	err := db.Select(&dbName, "SELECT current_database()")
	s.NoError(err)
	return dbName[0]
}

func TestDbTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
