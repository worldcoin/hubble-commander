package storage

import (
	"context"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db"
	bdg "github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var databaseTracer = otel.Tracer("database")

type Database struct {
	Badger *db.Database
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	badgerDB, err := db.NewDatabase(cfg.Badger)
	if err != nil {
		return nil, err
	}

	database := &Database{
		Badger: badgerDB,
	}

	if cfg.Bootstrap.Prune {
		err = database.Badger.Prune()
		if err != nil {
			return nil, err
		}
		log.Debug("Badger database was pruned")
	}

	return database, nil
}

func (d *Database) BeginTransaction(opts TxOptions) (*db.TxController, *Database) {
	database := *d

	badgerTx, badgerDB := d.Badger.BeginTransaction(!opts.ReadOnly)
	database.Badger = badgerDB

	return badgerTx, &database
}

func (d *Database) ExecuteInTransactionWithSpan(
	ctx context.Context,
	opts TxOptions,
	fn func(txCtx context.Context, txDatabase *Database) error,
) error {
	retries := 0
	err := d.unsafeExecuteInTransactionWithSpan(ctx, retries, opts, fn)
	for errors.Is(err, bdg.ErrConflict) {
		// nb. if we were already inside a transaction when this function is
		//     called then we run inside the outer transaction, so the `Commit`
		//     call is a no-op and this retry logic will never get a chance to
		//     fire
		log.WithError(err).Warn("Retrying transaction due to conflict")
		err = d.unsafeExecuteInTransactionWithSpan(ctx, retries, opts, fn)
		retries += 1
	}

	return err
}

// all errors are already wrapped w stack traces, except errors fn returns
func (d *Database) ExecuteInTransaction(opts TxOptions, fn func(txDatabase *Database) error) error {
	err := d.unsafeExecuteInTransaction(opts, fn)
	if errors.Is(err, bdg.ErrConflict) {
		// nb. if we were already inside a transaction when this function is
		//     called then we run inside the outer transaction, so the `Commit`
		//     call is a no-op and this retry logic will never get a chance to
		//     fire
		log.Debug("ExecuteInTransaction transaction conflicted, trying again")
		return d.ExecuteInTransaction(opts, fn)
	}
	return err
}

func (d *Database) unsafeExecuteInTransaction(opts TxOptions, fn func(txDatabase *Database) error) (err error) {
	txController, txDatabase := d.BeginTransaction(opts)
	defer txController.Rollback(&err)

	err = fn(txDatabase)
	if err != nil {
		return err
	}

	return txController.Commit()
}

func (d *Database) unsafeExecuteInTransactionWithSpan(
	ctx context.Context,
	retries int,
	opts TxOptions,
	fn func(txCtx context.Context, txDatabase *Database) error,
) (err error) {
	spanCtx, span := databaseTracer.Start(ctx, "database.ExecuteInTransaction")
	defer span.End()

	span.SetAttributes(attribute.Int("hubble.database.tx_retries", retries))

	txController, txDatabase := d.BeginTransaction(opts)
	defer txController.Rollback(&err)

	err = fn(spanCtx, txDatabase)
	if err != nil {
		return err
	}

	return func() error {
		_, innerSpan := databaseTracer.Start(spanCtx, "database.Commit")
		defer innerSpan.End()

		return txController.Commit()
	}()
}

func (d *Database) Close() error {
	return d.Badger.Close()
}
