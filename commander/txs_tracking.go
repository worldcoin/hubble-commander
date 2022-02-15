package commander

import (
	"github.com/Worldcoin/hubble-commander/commander/tracker"
)

func (c *Commander) startTrackingSentTxs() error {
	return tracker.StartTrackingSentTxs(c.workersContext, c.client, c.txsTrackingChannels.SentTxs)
}

func (c *Commander) startSendingRequestedTxs() error {
	return tracker.StartTxsRequestsSending(c.workersContext, c.txsTrackingChannels.Requests)
}
