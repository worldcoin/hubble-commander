package applier

import (
	st "github.com/Worldcoin/hubble-commander/storage"
)

type Applier struct {
	storage *st.Storage
}

func NewApplier(storage *st.Storage) *Applier {
	return &Applier{
		storage: storage,
	}
}
