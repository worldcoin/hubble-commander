package eth

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

const (
	NewBatchEvent               = "NewBatch"
	SinglePubkeyRegisteredEvent = "SinglePubkeyRegistered"
	BatchPubkeyRegisteredEvent  = "BatchPubkeyRegistered"
	TokenRegisteredEvent        = "TokenRegistered"
	SpokeRegisteredEvent        = "SpokeRegistered"
	DepositQueuedEvent          = "DepositQueued"
	DepositSubTreeReadyEvent    = "DepositSubTreeReady"
	DepositsFinalisedEvent      = "DepositsFinalised"
)

var eventTopics = map[string]common.Hash{
	NewBatchEvent:               crypto.Keccak256Hash([]byte("NewBatch(uint256,bytes32,uint8)")),
	SinglePubkeyRegisteredEvent: crypto.Keccak256Hash([]byte("SinglePubkeyRegistered(uint256)")),
	BatchPubkeyRegisteredEvent:  crypto.Keccak256Hash([]byte("BatchPubkeyRegistered(uint256,uint256)")),
	TokenRegisteredEvent:        crypto.Keccak256Hash([]byte("TokenRegistered(uint256,address)")),
	SpokeRegisteredEvent:        crypto.Keccak256Hash([]byte("SpokeRegistered(uint256,address)")),
	DepositQueuedEvent:          crypto.Keccak256Hash([]byte("DepositQueued(uint256,uint256,uint256,uint256,uint256)")),
	DepositSubTreeReadyEvent:    crypto.Keccak256Hash([]byte("DepositSubTreeReady(uint256,bytes32)")),
	DepositsFinalisedEvent:      crypto.Keccak256Hash([]byte("DepositsFinalised(uint256,bytes32,uint256)")),
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
