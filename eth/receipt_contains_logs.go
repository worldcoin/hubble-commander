package eth

import "github.com/ethereum/go-ethereum/core/types"

func receiptContainsLogs(receipt *types.Receipt) bool {
	return len(receipt.Logs) < 1 || receipt.Logs[0] == nil
}
