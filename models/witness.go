package models

import "github.com/ethereum/go-ethereum/common"

type Witness []common.Hash

func (w Witness) Bytes() [][32]byte {
	result := make([][32]byte, 0, len(w))
	for i := range w {
		result = append(result, w[i])
	}
	return result
}
