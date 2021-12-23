package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type RegisteredSpoke struct {
	ID       Uint256
	Contract common.Address
}
