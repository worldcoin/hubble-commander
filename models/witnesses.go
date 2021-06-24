package models

import "github.com/ethereum/go-ethereum/common"

type Witnesses []common.Hash

func (w Witnesses) Bytes() [][32]byte {
	result := make([][32]byte, 0, len(w))
	for i := range w {
		result = append(result, w[i])
		w[i].Bytes()
	}
	return result
}
