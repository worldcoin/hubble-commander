package testutils

import (
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/proofofburn"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type Rollup struct {
	Chooser         *proofofburn.ProofOfBurn
	AccountRegistry *accountregistry.AccountRegistry
}

func DeployRollup(deployer *bind.TransactOpts, backend bind.ContractBackend) (*Rollup, error) {
	_, _, proofOfBurn, err := proofofburn.DeployProofOfBurn(deployer, backend)
	if err != nil {
		return nil, err
	}

	_, _, accountRegistry, err := accountregistry.DeployAccountRegistry(deployer, backend)
	if err != nil {
		return nil, err
	}

	return &Rollup{
		Chooser:         proofOfBurn,
		AccountRegistry: accountRegistry,
	}, nil
}
