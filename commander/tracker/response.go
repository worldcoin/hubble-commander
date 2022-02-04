package tracker

import "github.com/ethereum/go-ethereum/core/types"

type txRequestResponse struct {
	Tx  *types.Transaction
	Err error
}
