package eth

import (
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type NewClientParams struct {
	Rollup          *rollup.Rollup
	AccountRegistry *accountregistry.AccountRegistry
	ClientConfig
}

type ClientConfig struct {
	txTimeout   *time.Duration  // default 60s
	stakeAmount *models.Uint256 // default 0.1 ether
}

type Client struct {
	account         bind.TransactOpts
	config          ClientConfig
	Rollup          *rollup.Rollup
	AccountRegistry *accountregistry.AccountRegistry
}

func NewClient(account *bind.TransactOpts, params NewClientParams) (*Client, error) {
	fillWithDefaults(&params.ClientConfig)

	return &Client{
		account:         *account,
		config:          params.ClientConfig,
		Rollup:          params.Rollup,
		AccountRegistry: params.AccountRegistry,
	}, nil
}

func fillWithDefaults(c *ClientConfig) {
	if c.txTimeout == nil {
		c.txTimeout = ref.Duration(60 * time.Second)
	}
	if c.stakeAmount == nil {
		c.stakeAmount = models.NewUint256(1e17)
	}
}

func (c *Client) withValue(value *models.Uint256) *bind.TransactOpts {
	opts := c.account
	opts.Value = &value.Int
	return &opts
}

func (c *Client) SubmitTransfersBatch(commitments []*models.Commitment) (batchID *models.Uint256, err error) {
	count := len(commitments)

	stateRoots := make([][32]byte, 0, count)
	signatures := make([][2]*big.Int, 0, count)
	feeReceivers := make([]*big.Int, 0, count)
	transactions := make([][]byte, 0, count)

	for _, commitment := range commitments {
		stateRoots = append(stateRoots, commitment.PostStateRoot)
		signatures = append(signatures, commitment.CombinedSignature.ToBigIntPointers())
		feeReceivers = append(feeReceivers, new(big.Int).SetUint64(uint64(commitment.FeeReceiver)))
		transactions = append(transactions, commitment.Transactions)
	}

	sink := make(chan *rollup.RollupNewBatch)
	subscription, err := c.Rollup.WatchNewBatch(&bind.WatchOpts{}, sink)
	if err != nil {
		return
	}
	defer subscription.Unsubscribe()

	tx, err := c.Rollup.SubmitTransfer(
		c.withValue(c.config.stakeAmount),
		stateRoots,
		signatures,
		feeReceivers,
		transactions,
	)
	if err != nil {
		return
	}

	for {
		select {
		case newBatch := <-sink:
			if newBatch.Raw.TxHash == tx.Hash() {
				return models.NewUint256FromBig(*newBatch.Index), nil
			}
		case <-time.After(*c.config.txTimeout):
			return nil, fmt.Errorf("timeout")
		}
	}
}
