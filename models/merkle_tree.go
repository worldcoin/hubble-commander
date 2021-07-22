package models

import "github.com/ethereum/go-ethereum/common"

type MerkleTreeNode struct {
	MerklePath MerklePath
	DataHash   common.Hash
}
