package storage

import (
	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/ethereum/go-ethereum/common"
)

var (
	zeroHashes [33]common.Hash
)

func init() {
	// Same as keccak256(abi.encode(0))
	zeroHashes[0] = common.HexToHash("0x290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563")
	for i := 1; i <= 32; i++ {
		zeroHashes[i] = utils.HashTwo(zeroHashes[i-1], zeroHashes[i-1])
	}
}

func GetZeroHash(level uint) common.Hash {
	if level > 32 {
		panic("level > 32")
	}

	return zeroHashes[level]
}
