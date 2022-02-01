package eth

import (
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

const (
	msgIsNotRollingBack     = "Is not rolling back"
	gasEstimatorErrorPrefix = "failed to estimate gas needed"
	errorSeparator          = ":"
)

var errKeepRollingBackFailed = errors.New("keep rolling back failed")

func (c *Client) KeepRollingBack() error {
	transaction, err := c.rollup().KeepRollingBack()
	if err != nil {
		return handleKeepRollingBackError(err)
	}
	return c.waitForKeepRollingBack(transaction)
}

func (c *Client) waitForKeepRollingBack(tx *types.Transaction) error {
	receipt, err := c.WaitToBeMined(tx)
	if err != nil {
		return err
	}
	if receipt.Status == types.ReceiptStatusSuccessful {
		return nil
	}

	invalidBatchID, err := c.GetInvalidBatchID()
	if err != nil {
		return err
	}
	if invalidBatchID == nil {
		return nil
	}

	return errKeepRollingBackFailed
}

func handleKeepRollingBackError(err error) error {
	errMsg := getGasEstimateErrorMessage(err)
	if errMsg == msgIsNotRollingBack {
		return nil
	}
	return err
}

func getGasEstimateErrorMessage(err error) string {
	msg := err.Error()
	if !strings.HasPrefix(msg, gasEstimatorErrorPrefix) {
		return msg
	}
	parts := strings.Split(msg, errorSeparator)
	return strings.TrimSpace(parts[len(parts)-1])
}
