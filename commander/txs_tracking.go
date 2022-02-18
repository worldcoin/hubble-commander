package commander

import (
	"github.com/Worldcoin/hubble-commander/commander/tracker"
)

func (c *Commander) trackSentTxs() error {
	return c.txsTracker.TrackSentTxs(c.workersContext)
}

func (c *Commander) sendRequestedTxs() error {
	return tracker.SendRequestedTxs(c.workersContext, c.txsTrackingChannels.Requests)
}
