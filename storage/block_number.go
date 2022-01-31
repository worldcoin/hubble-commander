package storage

import (
	"sync/atomic"
)

func (s *ChainStateStorage) SetLatestBlockNumber(blockNumber uint32) {
	atomic.StoreUint32(&s.latestBlockNumber, blockNumber)
}

func (s *ChainStateStorage) GetLatestBlockNumber() uint32 {
	return atomic.LoadUint32(&s.latestBlockNumber)
}

func (s *ChainStateStorage) SetSyncedBlock(blockNumber uint64) error {
	atomic.StoreUint64(&s.syncedBlock, blockNumber)
	chainState, err := s.GetChainState()
	if err != nil {
		return err
	}
	chainState.SyncedBlock = blockNumber
	return s.SetChainState(chainState)
}

func (s *ChainStateStorage) GetSyncedBlock() (*uint64, error) {
	syncedBlock := atomic.LoadUint64(&s.syncedBlock)
	if syncedBlock != 0 {
		return &syncedBlock, nil
	}
	chainState, err := s.GetChainState()
	if IsNotFoundError(err) {
		return &syncedBlock, nil
	}
	if err != nil {
		return nil, err
	}

	atomic.StoreUint64(&s.syncedBlock, chainState.SyncedBlock)
	return &chainState.SyncedBlock, nil
}
