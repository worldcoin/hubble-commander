package commander

import "time"

func (c *Commander) trackingTxs() error {
	for {
		select {
		case <-c.workersContext.Done():
			return nil
		case err := <-c.txsTracker.Fail():
			return err
		case err := <-c.client.AccountManager.RegistrationFail():
			return err
		default:
			time.Sleep(time.Millisecond * 300)
		}
	}
}
