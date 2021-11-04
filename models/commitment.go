package models

import "github.com/ethereum/go-ethereum/common"

type Commitment interface {
	GetBodyHash() common.Hash
	GetPostStateRoot() common.Hash
	LeafHash() common.Hash
}
