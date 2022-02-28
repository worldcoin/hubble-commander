package eth

import (
	"strings"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/accountregistry"
	"github.com/Worldcoin/hubble-commander/eth/chain"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

type AccountManager struct {
	Blockchain                       chain.Connection
	AccountRegistry                  *AccountRegistry
	batchAccountRegistrationGasLimit uint64
	mineTimeout                      time.Duration
	requestsChan                     chan<- *TxSendingRequest
}

//goland:noinspection GoDeprecation
func NewAccountManager(blockchain chain.Connection, params *AccountManagerParams) (*AccountManager, error) {
	accountRegistryAbi, err := abi.JSON(strings.NewReader(accountregistry.AccountRegistryABI))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	backend := blockchain.GetBackend()
	accountRegistryContract := bind.NewBoundContract(params.AccountRegistryAddress, accountRegistryAbi, backend, backend, backend)
	return &AccountManager{
		Blockchain: blockchain,
		AccountRegistry: &AccountRegistry{
			AccountRegistry: params.AccountRegistry,
			Contract:        MakeContract(&accountRegistryAbi, accountRegistryContract),
		},
		batchAccountRegistrationGasLimit: params.BatchAccountRegistrationGasLimit,
		mineTimeout:                      params.MineTimeout,
		requestsChan:                     params.RequestsChan,
	}, nil
}

func (a *AccountManager) packAndRequest(
	contract *Contract,
	opts *bind.TransactOpts,
	shouldTrackTx bool,
	method string,
	data ...interface{},
) (*types.Transaction, error) {
	return packAndRequest(a.requestsChan, contract, opts, shouldTrackTx, method, data...)
}

type AccountManagerParams struct {
	AccountRegistry                  *accountregistry.AccountRegistry
	AccountRegistryAddress           common.Address
	BatchAccountRegistrationGasLimit uint64
	MineTimeout                      time.Duration
	RequestsChan                     chan<- *TxSendingRequest
}
