package postgres

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
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

	s.checkBatch(s.db, 0)

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

	s.addBatch()

	cloneCfg := *s.config
	cloneCfg.Name += "_clone"
	clonedDB, err := s.db.Clone(&cloneCfg, s.config.Name)
	s.NoError(err)

	s.checkBatch(clonedDB, 1)
	s.checkBatch(s.db, 1)
}

func (s *DBTestSuite) checkBatch(db *Database, expectedLength int) {
	res := make([]models.Batch, 0, 1)
	err := db.Query(
		sq.Select("*").From("batch"),
	).Into(&res)
	s.NoError(err)
	s.Len(res, expectedLength)
}

func (s *DBTestSuite) addBatch() {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	batch := models.Batch{
		ID:              models.MakeUint256(1),
		Type:            txtype.Transfer,
		TransactionHash: utils.RandomHash(),
	}
	query, args, err := qb.Insert("batch").
		Values(
			batch.ID,
			batch.Type,
			batch.TransactionHash,
		).ToSql()
	s.NoError(err)

	_, err = s.db.Exec(query, args...)
	s.NoError(err)
}

func TestDbTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
