package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type RegisteredToken struct {
	ID       Uint256
	Contract common.Address
}
