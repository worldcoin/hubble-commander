package eth

import (
	"github.com/Worldcoin/hubble-commander/metrics"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

type Iterator interface {
	SetData(contract *bind.BoundContract, event string, logs chan types.Log, sub ethereum.Subscription)
	Error() error
}

func (c *Client) FilterLogs(contract *bind.BoundContract, eventName string, opts *bind.FilterOpts, it Iterator) (err error) {
	var (
		logs chan types.Log
		sub  event.Subscription
	)

	err = it.Error()
	if err != nil {
		return err
	}

	duration, err := metrics.MeasureDuration(func() error {
		logs, sub, err = contract.FilterLogs(opts, eventName)
		return err
	})
	if err != nil {
		return err
	}

	c.Metrics.SaveBlockchainCallDurationMeasurement(*duration, metrics.EventNameToMetricsEventFilterCallLabel(eventName))

	it.SetData(contract, eventName, logs, sub)

	return nil
}
