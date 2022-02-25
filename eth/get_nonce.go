package eth

import (
	"context"

	"github.com/pkg/errors"
)

func (c *Client) GetNonce() (uint64, error) {
	nonce, err := c.Blockchain.GetBackend().PendingNonceAt(context.Background(), c.Blockchain.GetAccount().From)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return nonce, nil
}
