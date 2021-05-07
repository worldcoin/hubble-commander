package badger

import "github.com/dgraph-io/badger/v3"

type ControllerAdapter struct {
	*badger.Txn
}

func (a *ControllerAdapter) Rollback() error {
	a.Txn.Discard()
	return nil
}
