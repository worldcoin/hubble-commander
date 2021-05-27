package storage

func (s *Storage) SetLatestBlockNumber(blockNumber uint32) {
	s.latestBlockNumber = blockNumber
}

func (s *Storage) GetLatestBlockNumber() uint32 {
	return s.latestBlockNumber
}
