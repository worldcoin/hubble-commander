package storage

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

var (
	zeroHashes [32]common.Hash
)

func init() {
	// keccak256(abi.encode(0))
	zeroHashes[0] = common.HexToHash("0x290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563")
	for i := 1; i < 32; i++ {
		buf := make([]byte, 64, 64)
		copy(buf[0:32], zeroHashes[i-1].Bytes())
		copy(buf[32:64], zeroHashes[i-1].Bytes())
		zeroHashes[i] = keccack256(buf)
	}
}

func GetZeroHash(level uint) common.Hash {
	if level > 31 {
		panic("level >= 32")
	}

	return zeroHashes[level]
}

func keccack256(x []byte) common.Hash {
	hashWriter := sha3.NewLegacyKeccak256()
	hashWriter.Write(x)
	hash := common.Hash{}
	hashWriter.Sum(hash[:0])
	return hash
}

func (storage *Storage) GetWitness(path models.MerklePath) ([]common.Hash, error) {
	storage.DB
}
