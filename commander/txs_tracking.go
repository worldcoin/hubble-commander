package commander

import (
	"github.com/Worldcoin/hubble-commander/commander/tracker"
)

func (c *Commander) trackSentTxs() error {
	return tracker.TrackSentTxs(c.workersContext, c.client, c.txsTrackingChannels.SentTxs)
}

func (c *Commander) sendRequestedTxs() error {
	return tracker.SendRequestedTxs(c.workersContext, c.txsTrackingChannels.Requests)
}
