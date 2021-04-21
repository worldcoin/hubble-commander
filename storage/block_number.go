package storage

var latestBlockNumber uint32

func (s *Storage) SetLatestBlockNumber(blockNumber uint32) {
	latestBlockNumber = blockNumber
}

func (s *Storage) GetLatestBlockNumber() uint32 {
	return latestBlockNumber
}
