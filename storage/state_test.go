package storage

import "github.com/ethereum/go-ethereum/common"

func (s *StorageTestSuite) TestZeroHash() {
	s.Equal(common.HexToHash("0x78ccaaab73373552f207a63599de54d7d8d0c1805f86ce7da15818d09f4cff62"), GetZeroHash(31))
}
