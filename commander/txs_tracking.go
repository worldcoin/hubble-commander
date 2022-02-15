package commander

import (
	"github.com/Worldcoin/hubble-commander/commander/tracker"
)

func (c *Commander) startFailedTxsTracking() error {
	return tracker.StartFailedTxsTracking(c.workersContext, c.client, c.txsTrackingChannels.SentTxs)
}

func (c *Commander) startTxsRequestsSending() error {
	return tracker.StartTxsRequestsSending(c.workersContext, c.txsTrackingChannels.Requests)
}
