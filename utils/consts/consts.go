package consts

import "github.com/ethereum/go-ethereum/common"

const (
	L2Unit        = 1e9
	AuthKeyHeader = "Auth-Key"
)

// ZeroHash is the same as keccak256(abi.encode(0))
var ZeroHash = common.HexToHash("0x290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563")
