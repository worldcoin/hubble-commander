package commander

import (
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/Worldcoin/hubble-commander/utils/ref"
)

func GetChainConnection(cfg *config.EthereumConfig) (chain.Connection, error) {
	if cfg.RPCURL == "simulator" {
		return simulator.NewConfiguredSimulator(simulator.Config{
			FirstAccountPrivateKey: ref.String(cfg.PrivateKeys[0]),
			AutomineEnabled:        ref.Bool(true),
		})
	}
	return chain.NewRPCConnection(cfg)
}
