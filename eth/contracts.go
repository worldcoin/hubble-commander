package eth

import (
	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/contracts/depositmanager"
	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/contracts/spokeregistry"
	"github.com/Worldcoin/hubble-commander/contracts/tokenregistry"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type Contract struct {
	ABI           *abi.ABI
	BoundContract *bind.BoundContract
}

func MakeContract(ABI *abi.ABI, boundContract *bind.BoundContract) Contract {
	return Contract{ABI: ABI, BoundContract: boundContract}
}

type AccountRegistry struct {
	*accountregistry.AccountRegistry
	Contract
}

type Rollup struct {
	*rollup.Rollup
	Contract
}

type TokenRegistry struct {
	*tokenregistry.TokenRegistry
	Contract
}

type DepositManager struct {
	*depositmanager.DepositManager
	Contract
}

type SpokeRegistry struct {
	*spokeregistry.SpokeRegistry
	Contract
}
