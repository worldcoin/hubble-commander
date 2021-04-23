package storage

import "github.com/ethereum/go-ethereum/common"

func (s *Storage) GetTransaction(hash common.Hash) (interface{}, error) {
	transfer, err := s.GetTransfer(hash)
	if err != nil && !IsNotFoundError(err) {
		return nil, err
	}
	if transfer != nil {
		return transfer, nil
	}

	return s.GetCreate2Transfer(hash)
}
