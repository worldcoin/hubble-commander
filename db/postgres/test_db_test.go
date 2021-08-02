package postgres

import (
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
)

func TestTestDB_Clone(t *testing.T) {
	testDB, err := NewTestDB()
	require.NoError(t, err)

	addBatch(t, testDB.DB)

	clonedDB, err := testDB.Clone(config.GetTestConfig().Postgres)
	require.NoError(t, err)

	checkBatch(t, clonedDB.DB, 1)
	checkBatch(t, testDB.DB, 1)

	err = testDB.Teardown()
	require.NoError(t, err)

	err = clonedDB.Teardown()
	require.NoError(t, err)
}

func checkBatch(t *testing.T, db *Database, expectedLength int) {
	res := make([]models.Batch, 0, 1)
	err := db.Query(
		sq.Select("*").From("batch"),
	).Into(&res)
	require.NoError(t, err)
	require.Len(t, res, expectedLength)
}

func addBatch(t *testing.T, db *Database) {
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
	require.NoError(t, err)

	_, err = db.Exec(query, args...)
	require.NoError(t, err)
}
