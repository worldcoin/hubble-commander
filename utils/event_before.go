package utils

import "github.com/ethereum/go-ethereum/core/types"

func EventBefore(left, right *types.Log) bool {
	return left.BlockNumber < right.BlockNumber ||
		(left.BlockNumber == right.BlockNumber && left.Index < right.Index)
}
