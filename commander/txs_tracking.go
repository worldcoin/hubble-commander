package commander

import (
	"github.com/Worldcoin/hubble-commander/commander/tracker"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/ethereum/go-ethereum/core/types"
)

func (c *Commander) startFailedTxsTracking(txsHashChan <-chan *types.Transaction) error {
	return tracker.StartFailedTxsTracking(c.workersContext, c.client, txsHashChan)
}

func (c *Commander) startTxsRequestsSending(requestsChan <-chan *eth.TxSendingRequest) error {
	return tracker.StartTxsRequestsSending(c.workersContext, requestsChan)
}
