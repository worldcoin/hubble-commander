// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package spokeregistry

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// SpokeRegistryMetaData contains all meta data concerning the SpokeRegistry contract.
var SpokeRegistryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"spokeID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"spokeContract\",\"type\":\"address\"}],\"name\":\"RegisteredSpoke\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"spokeID\",\"type\":\"uint256\"}],\"name\":\"getSpokeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"numSpokes\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spokeContract\",\"type\":\"address\"}],\"name\":\"registerSpoke\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"registeredSpokes\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506101a4806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c806350177aef146100515780639307139714610079578063d9bc6c2414610093578063da2fd3df146100cc575b600080fd5b6100776004803603602081101561006757600080fd5b50356001600160a01b03166100e9565b005b61008161015b565b60408051918252519081900360200190f35b6100b0600480360360208110156100a957600080fd5b5035610161565b604080516001600160a01b039092168252519081900360200190f35b6100b0600480360360208110156100e257600080fd5b503561017c565b6001805481018082556000908152602081815260409182902080546001600160a01b0319166001600160a01b038616908117909155925482519081529081019290925280517f94ede2aea46f7f1f64cc4bf492cc08fb2119502829fb3f53acc0186e2a8e64d19281900390910190a150565b60015481565b6000602081905290815260409020546001600160a01b031681565b6000908152602081905260409020546001600160a01b03169056fea164736f6c634300060c000a",
}

// SpokeRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use SpokeRegistryMetaData.ABI instead.
var SpokeRegistryABI = SpokeRegistryMetaData.ABI

// SpokeRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use SpokeRegistryMetaData.Bin instead.
var SpokeRegistryBin = SpokeRegistryMetaData.Bin

// DeploySpokeRegistry deploys a new Ethereum contract, binding an instance of SpokeRegistry to it.
func DeploySpokeRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SpokeRegistry, error) {
	parsed, err := SpokeRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(SpokeRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SpokeRegistry{SpokeRegistryCaller: SpokeRegistryCaller{contract: contract}, SpokeRegistryTransactor: SpokeRegistryTransactor{contract: contract}, SpokeRegistryFilterer: SpokeRegistryFilterer{contract: contract}}, nil
}

// SpokeRegistry is an auto generated Go binding around an Ethereum contract.
type SpokeRegistry struct {
	SpokeRegistryCaller     // Read-only binding to the contract
	SpokeRegistryTransactor // Write-only binding to the contract
	SpokeRegistryFilterer   // Log filterer for contract events
}

// SpokeRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type SpokeRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SpokeRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SpokeRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SpokeRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SpokeRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SpokeRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SpokeRegistrySession struct {
	Contract     *SpokeRegistry    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SpokeRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SpokeRegistryCallerSession struct {
	Contract *SpokeRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// SpokeRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SpokeRegistryTransactorSession struct {
	Contract     *SpokeRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// SpokeRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type SpokeRegistryRaw struct {
	Contract *SpokeRegistry // Generic contract binding to access the raw methods on
}

// SpokeRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SpokeRegistryCallerRaw struct {
	Contract *SpokeRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// SpokeRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SpokeRegistryTransactorRaw struct {
	Contract *SpokeRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSpokeRegistry creates a new instance of SpokeRegistry, bound to a specific deployed contract.
func NewSpokeRegistry(address common.Address, backend bind.ContractBackend) (*SpokeRegistry, error) {
	contract, err := bindSpokeRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SpokeRegistry{SpokeRegistryCaller: SpokeRegistryCaller{contract: contract}, SpokeRegistryTransactor: SpokeRegistryTransactor{contract: contract}, SpokeRegistryFilterer: SpokeRegistryFilterer{contract: contract}}, nil
}

// NewSpokeRegistryCaller creates a new read-only instance of SpokeRegistry, bound to a specific deployed contract.
func NewSpokeRegistryCaller(address common.Address, caller bind.ContractCaller) (*SpokeRegistryCaller, error) {
	contract, err := bindSpokeRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SpokeRegistryCaller{contract: contract}, nil
}

// NewSpokeRegistryTransactor creates a new write-only instance of SpokeRegistry, bound to a specific deployed contract.
func NewSpokeRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*SpokeRegistryTransactor, error) {
	contract, err := bindSpokeRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SpokeRegistryTransactor{contract: contract}, nil
}

// NewSpokeRegistryFilterer creates a new log filterer instance of SpokeRegistry, bound to a specific deployed contract.
func NewSpokeRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*SpokeRegistryFilterer, error) {
	contract, err := bindSpokeRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SpokeRegistryFilterer{contract: contract}, nil
}

// bindSpokeRegistry binds a generic wrapper to an already deployed contract.
func bindSpokeRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SpokeRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SpokeRegistry *SpokeRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SpokeRegistry.Contract.SpokeRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SpokeRegistry *SpokeRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SpokeRegistry.Contract.SpokeRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SpokeRegistry *SpokeRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SpokeRegistry.Contract.SpokeRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SpokeRegistry *SpokeRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SpokeRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SpokeRegistry *SpokeRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SpokeRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SpokeRegistry *SpokeRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SpokeRegistry.Contract.contract.Transact(opts, method, params...)
}

// GetSpokeAddress is a free data retrieval call binding the contract method 0xda2fd3df.
//
// Solidity: function getSpokeAddress(uint256 spokeID) view returns(address)
func (_SpokeRegistry *SpokeRegistryCaller) GetSpokeAddress(opts *bind.CallOpts, spokeID *big.Int) (common.Address, error) {
	var out []interface{}
	err := _SpokeRegistry.contract.Call(opts, &out, "getSpokeAddress", spokeID)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetSpokeAddress is a free data retrieval call binding the contract method 0xda2fd3df.
//
// Solidity: function getSpokeAddress(uint256 spokeID) view returns(address)
func (_SpokeRegistry *SpokeRegistrySession) GetSpokeAddress(spokeID *big.Int) (common.Address, error) {
	return _SpokeRegistry.Contract.GetSpokeAddress(&_SpokeRegistry.CallOpts, spokeID)
}

// GetSpokeAddress is a free data retrieval call binding the contract method 0xda2fd3df.
//
// Solidity: function getSpokeAddress(uint256 spokeID) view returns(address)
func (_SpokeRegistry *SpokeRegistryCallerSession) GetSpokeAddress(spokeID *big.Int) (common.Address, error) {
	return _SpokeRegistry.Contract.GetSpokeAddress(&_SpokeRegistry.CallOpts, spokeID)
}

// NumSpokes is a free data retrieval call binding the contract method 0x93071397.
//
// Solidity: function numSpokes() view returns(uint256)
func (_SpokeRegistry *SpokeRegistryCaller) NumSpokes(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SpokeRegistry.contract.Call(opts, &out, "numSpokes")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumSpokes is a free data retrieval call binding the contract method 0x93071397.
//
// Solidity: function numSpokes() view returns(uint256)
func (_SpokeRegistry *SpokeRegistrySession) NumSpokes() (*big.Int, error) {
	return _SpokeRegistry.Contract.NumSpokes(&_SpokeRegistry.CallOpts)
}

// NumSpokes is a free data retrieval call binding the contract method 0x93071397.
//
// Solidity: function numSpokes() view returns(uint256)
func (_SpokeRegistry *SpokeRegistryCallerSession) NumSpokes() (*big.Int, error) {
	return _SpokeRegistry.Contract.NumSpokes(&_SpokeRegistry.CallOpts)
}

// RegisteredSpokes is a free data retrieval call binding the contract method 0xd9bc6c24.
//
// Solidity: function registeredSpokes(uint256 ) view returns(address)
func (_SpokeRegistry *SpokeRegistryCaller) RegisteredSpokes(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _SpokeRegistry.contract.Call(opts, &out, "registeredSpokes", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RegisteredSpokes is a free data retrieval call binding the contract method 0xd9bc6c24.
//
// Solidity: function registeredSpokes(uint256 ) view returns(address)
func (_SpokeRegistry *SpokeRegistrySession) RegisteredSpokes(arg0 *big.Int) (common.Address, error) {
	return _SpokeRegistry.Contract.RegisteredSpokes(&_SpokeRegistry.CallOpts, arg0)
}

// RegisteredSpokes is a free data retrieval call binding the contract method 0xd9bc6c24.
//
// Solidity: function registeredSpokes(uint256 ) view returns(address)
func (_SpokeRegistry *SpokeRegistryCallerSession) RegisteredSpokes(arg0 *big.Int) (common.Address, error) {
	return _SpokeRegistry.Contract.RegisteredSpokes(&_SpokeRegistry.CallOpts, arg0)
}

// RegisterSpoke is a paid mutator transaction binding the contract method 0x50177aef.
//
// Solidity: function registerSpoke(address spokeContract) returns()
func (_SpokeRegistry *SpokeRegistryTransactor) RegisterSpoke(opts *bind.TransactOpts, spokeContract common.Address) (*types.Transaction, error) {
	return _SpokeRegistry.contract.Transact(opts, "registerSpoke", spokeContract)
}

// RegisterSpoke is a paid mutator transaction binding the contract method 0x50177aef.
//
// Solidity: function registerSpoke(address spokeContract) returns()
func (_SpokeRegistry *SpokeRegistrySession) RegisterSpoke(spokeContract common.Address) (*types.Transaction, error) {
	return _SpokeRegistry.Contract.RegisterSpoke(&_SpokeRegistry.TransactOpts, spokeContract)
}

// RegisterSpoke is a paid mutator transaction binding the contract method 0x50177aef.
//
// Solidity: function registerSpoke(address spokeContract) returns()
func (_SpokeRegistry *SpokeRegistryTransactorSession) RegisterSpoke(spokeContract common.Address) (*types.Transaction, error) {
	return _SpokeRegistry.Contract.RegisterSpoke(&_SpokeRegistry.TransactOpts, spokeContract)
}

// SpokeRegistryRegisteredSpokeIterator is returned from FilterRegisteredSpoke and is used to iterate over the raw logs and unpacked data for RegisteredSpoke events raised by the SpokeRegistry contract.
type SpokeRegistryRegisteredSpokeIterator struct {
	Event *SpokeRegistryRegisteredSpoke // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SpokeRegistryRegisteredSpokeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SpokeRegistryRegisteredSpoke)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SpokeRegistryRegisteredSpoke)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SpokeRegistryRegisteredSpokeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SpokeRegistryRegisteredSpokeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SpokeRegistryRegisteredSpoke represents a RegisteredSpoke event raised by the SpokeRegistry contract.
type SpokeRegistryRegisteredSpoke struct {
	SpokeID       *big.Int
	SpokeContract common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterRegisteredSpoke is a free log retrieval operation binding the contract event 0x94ede2aea46f7f1f64cc4bf492cc08fb2119502829fb3f53acc0186e2a8e64d1.
//
// Solidity: event RegisteredSpoke(uint256 spokeID, address spokeContract)
func (_SpokeRegistry *SpokeRegistryFilterer) FilterRegisteredSpoke(opts *bind.FilterOpts) (*SpokeRegistryRegisteredSpokeIterator, error) {

	logs, sub, err := _SpokeRegistry.contract.FilterLogs(opts, "RegisteredSpoke")
	if err != nil {
		return nil, err
	}
	return &SpokeRegistryRegisteredSpokeIterator{contract: _SpokeRegistry.contract, event: "RegisteredSpoke", logs: logs, sub: sub}, nil
}

// WatchRegisteredSpoke is a free log subscription operation binding the contract event 0x94ede2aea46f7f1f64cc4bf492cc08fb2119502829fb3f53acc0186e2a8e64d1.
//
// Solidity: event RegisteredSpoke(uint256 spokeID, address spokeContract)
func (_SpokeRegistry *SpokeRegistryFilterer) WatchRegisteredSpoke(opts *bind.WatchOpts, sink chan<- *SpokeRegistryRegisteredSpoke) (event.Subscription, error) {

	logs, sub, err := _SpokeRegistry.contract.WatchLogs(opts, "RegisteredSpoke")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SpokeRegistryRegisteredSpoke)
				if err := _SpokeRegistry.contract.UnpackLog(event, "RegisteredSpoke", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRegisteredSpoke is a log parse operation binding the contract event 0x94ede2aea46f7f1f64cc4bf492cc08fb2119502829fb3f53acc0186e2a8e64d1.
//
// Solidity: event RegisteredSpoke(uint256 spokeID, address spokeContract)
func (_SpokeRegistry *SpokeRegistryFilterer) ParseRegisteredSpoke(log types.Log) (*SpokeRegistryRegisteredSpoke, error) {
	event := new(SpokeRegistryRegisteredSpoke)
	if err := _SpokeRegistry.contract.UnpackLog(event, "RegisteredSpoke", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
