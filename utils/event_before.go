package utils

import "github.com/ethereum/go-ethereum/core/types"

func EventBefore(left *types.Log, right *types.Log) bool {
	return left.BlockNumber < right.BlockNumber ||
		(left.BlockNumber == right.BlockNumber && left.Index < right.Index)
}
