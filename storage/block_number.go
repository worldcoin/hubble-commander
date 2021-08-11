package storage

import "github.com/Worldcoin/hubble-commander/utils/ref"

func (s *ChainStateStorage) SetLatestBlockNumber(blockNumber uint32) {
	s.latestBlockNumber = blockNumber
}

func (s *ChainStateStorage) GetLatestBlockNumber() uint32 {
	return s.latestBlockNumber
}

func (s *ChainStateStorage) SetSyncedBlock(blockNumber uint64) error {
	s.syncedBlock = &blockNumber
	chainState, err := s.GetChainState()
	if err != nil {
		return err
	}
	chainState.SyncedBlock = blockNumber
	return s.SetChainState(chainState)
}

func (s *ChainStateStorage) GetSyncedBlock() (*uint64, error) {
	if s.syncedBlock != nil {
		return s.syncedBlock, nil
	}
	chainState, err := s.GetChainState()
	if IsNotFoundError(err) {
		return ref.Uint64(0), nil
	}
	if err != nil {
		return nil, err
	}
	return &chainState.SyncedBlock, nil
}
