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

	addTransfer(t, testDB.DB)

	clonedDB, err := testDB.Clone(config.GetTestConfig().Postgres)
	require.NoError(t, err)

	checkTransfer(t, clonedDB.DB, 1)
	checkTransfer(t, testDB.DB, 1)

	err = testDB.Teardown()
	require.NoError(t, err)

	err = clonedDB.Teardown()
	require.NoError(t, err)
}

func checkTransfer(t *testing.T, db *Database, expectedLength int) {
	res := make([]models.TransactionBase, 0, 1)
	err := db.Query(
		sq.Select("*").From("transaction_base"),
	).Into(&res)
	require.NoError(t, err)
	require.Len(t, res, expectedLength)
}

func addTransfer(t *testing.T, db *Database) {
	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	transfer := models.TransactionBase{
		Hash:        utils.RandomHash(),
		TxType:      txtype.Transfer,
		FromStateID: 5,
	}
	query, args, err := qb.Insert("transaction_base").
		Values(
			transfer.Hash,
			transfer.TxType,
			transfer.FromStateID,
			transfer.Amount,
			transfer.Fee,
			transfer.Nonce,
			transfer.Signature,
		).ToSql()
	require.NoError(t, err)

	_, err = db.Exec(query, args...)
	require.NoError(t, err)
}
