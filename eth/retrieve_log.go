package eth

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

var eventTopics = map[string]common.Hash{
	"NewBatch":               crypto.Keccak256Hash([]byte("NewBatch(uint256,bytes32,uint8)")),
	"SinglePubkeyRegistered": crypto.Keccak256Hash([]byte("SinglePubkeyRegistered(uint256)")),
	"BatchPubkeyRegistered":  crypto.Keccak256Hash([]byte("BatchPubkeyRegistered(uint256,uint256)")),
	"RegisteredToken":        crypto.Keccak256Hash([]byte("RegisteredToken(uint256,address)")),
	"DepositQueued":          crypto.Keccak256Hash([]byte("DepositQueued(uint256,uint256,uint256)")),
}

var (
	ErrReceiptWithoutLogs = fmt.Errorf("the receipt contains no logs")
)

func retrieveLog(receipt *types.Receipt, logName string) (*types.Log, error) {
	if len(receipt.Logs) < 1 {
		return nil, errors.WithStack(ErrReceiptWithoutLogs)
	}

	for i := range receipt.Logs {
		log := receipt.Logs[i]
		if log.Topics[0] == eventTopics[logName] {
			return log, nil
		}
	}

	return nil, errors.WithStack(NewLogNotFoundError(logName))
}
