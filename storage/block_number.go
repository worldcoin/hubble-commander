package storage

import (
	"sync/atomic"

	"github.com/Worldcoin/hubble-commander/models/stored"
	"github.com/dgraph-io/badger/v3"
	"github.com/pkg/errors"
)

var SyncedBlockKey = []byte("SyncedBlock")

func (s *ChainStateStorage) SetLatestBlockNumber(blockNumber uint32) {
	atomic.StoreUint32(&s.latestBlockNumber, blockNumber)
}

func (s *ChainStateStorage) GetLatestBlockNumber() uint32 {
	return atomic.LoadUint32(&s.latestBlockNumber)
}

func (s *ChainStateStorage) SetSyncedBlock(blockNumber uint64) error {
	atomic.StoreUint64(&s.syncedBlock, blockNumber)

	return s.database.Badger.RawUpdate(func(txn *badger.Txn) error {
		value := stored.EncodeUint64(blockNumber)
		return txn.Set(SyncedBlockKey, value)
	})
}

func (s *ChainStateStorage) GetSyncedBlock() (*uint64, error) {
	syncedBlock := atomic.LoadUint64(&s.syncedBlock)
	if syncedBlock != 0 {
		return &syncedBlock, nil
	}

	err := s.database.Badger.View(func(txn *badger.Txn) error {
		item, err := txn.Get(SyncedBlockKey)
		if err != nil {
			return errors.WithStack(err)
		}

		err = item.Value(func(val []byte) error {
			innerErr := stored.DecodeUint64(val, &syncedBlock)
			if innerErr != nil {
				panic("stored.DecodeUint64 never returns error")
			}
			return nil
		})
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
	if errors.Is(err, badger.ErrKeyNotFound) {
		zero := uint64(0)
		return &zero, nil
	}
	if err != nil {
		return nil, err
	}

	atomic.StoreUint64(&s.syncedBlock, syncedBlock)
	return &syncedBlock, nil
}
