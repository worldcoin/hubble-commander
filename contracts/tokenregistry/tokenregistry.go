// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tokenregistry

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

// TokenRegistryMetaData contains all meta data concerning the TokenRegistry contract.
var TokenRegistryMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenContract\",\"type\":\"address\"}],\"name\":\"TokenRegistered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"nextTokenID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenContract\",\"type\":\"address\"}],\"name\":\"registerToken\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"}],\"name\":\"safeGetRecord\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"l2Unit\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052600060025534801561001557600080fd5b50610412806100256000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806309824a80146100465780630a7973b71461006e578063f101e481146100ae575b600080fd5b61006c6004803603602081101561005c57600080fd5b50356001600160a01b03166100c8565b005b61008b6004803603602081101561008457600080fd5b5035610326565b604080516001600160a01b03909316835260208301919091528051918290030190f35b6100b66103c5565b60408051918252519081900360200190f35b6001600160a01b03811660009081526001602052604090205460ff1615610136576040805162461bcd60e51b815260206004820152601960248201527f546f6b656e20616c726561647920726567697374657265642e00000000000000604482015290519081900360640190fd5b6012816001600160a01b031663313ce5676040518163ffffffff1660e01b815260040160206040518083038186803b15801561017157600080fd5b505afa158015610185573d6000803e3d6000fd5b505050506040513d602081101561019b57600080fd5b505160ff1611156101f3576040805162461bcd60e51b815260206004820152601960248201527f446f6e277420736572766520646563696d616c73203e20313800000000000000604482015290519081900360640190fd5b60006009826001600160a01b031663313ce5676040518163ffffffff1660e01b815260040160206040518083038186803b15801561023057600080fd5b505afa158015610244573d6000803e3d6000fd5b505050506040513d602081101561025a57600080fd5b50516040805180820182526001600160a01b0386811680835260ff909416949094101560208083018281526002805460009081528084528681209551865493516001600160a01b031990941699169890981760ff60a01b1916600160a01b9215159290920291909117909355848652600180825295849020805460ff1916909617909555905482519081529384019290925280519193507f5495d5c40ac716093397abf16b93fda5ccd72fb25488d6c41865b87898abea7c928290030190a15050600280546001019055565b6000806103316103cb565b50600083815260208181526040918290208251808401909352546001600160a01b038116808452600160a01b90910460ff161515918301919091526103a75760405162461bcd60e51b81526004018080602001828103825260238152602001806103e36023913960400191505060405180910390fd5b600191508060200151156103bd57633b9aca0091505b519150915091565b60025481565b60408051808201909152600080825260208201529056fe546f6b656e52656769737472793a20556e7265676973746572656420746f6b656e4944a164736f6c634300060c000a",
}

// TokenRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use TokenRegistryMetaData.ABI instead.
var TokenRegistryABI = TokenRegistryMetaData.ABI

// TokenRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TokenRegistryMetaData.Bin instead.
var TokenRegistryBin = TokenRegistryMetaData.Bin

// DeployTokenRegistry deploys a new Ethereum contract, binding an instance of TokenRegistry to it.
func DeployTokenRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TokenRegistry, error) {
	parsed, err := TokenRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TokenRegistryBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TokenRegistry{TokenRegistryCaller: TokenRegistryCaller{contract: contract}, TokenRegistryTransactor: TokenRegistryTransactor{contract: contract}, TokenRegistryFilterer: TokenRegistryFilterer{contract: contract}}, nil
}

// TokenRegistry is an auto generated Go binding around an Ethereum contract.
type TokenRegistry struct {
	TokenRegistryCaller     // Read-only binding to the contract
	TokenRegistryTransactor // Write-only binding to the contract
	TokenRegistryFilterer   // Log filterer for contract events
}

// TokenRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenRegistrySession struct {
	Contract     *TokenRegistry    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenRegistryCallerSession struct {
	Contract *TokenRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// TokenRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenRegistryTransactorSession struct {
	Contract     *TokenRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// TokenRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenRegistryRaw struct {
	Contract *TokenRegistry // Generic contract binding to access the raw methods on
}

// TokenRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenRegistryCallerRaw struct {
	Contract *TokenRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// TokenRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenRegistryTransactorRaw struct {
	Contract *TokenRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenRegistry creates a new instance of TokenRegistry, bound to a specific deployed contract.
func NewTokenRegistry(address common.Address, backend bind.ContractBackend) (*TokenRegistry, error) {
	contract, err := bindTokenRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenRegistry{TokenRegistryCaller: TokenRegistryCaller{contract: contract}, TokenRegistryTransactor: TokenRegistryTransactor{contract: contract}, TokenRegistryFilterer: TokenRegistryFilterer{contract: contract}}, nil
}

// NewTokenRegistryCaller creates a new read-only instance of TokenRegistry, bound to a specific deployed contract.
func NewTokenRegistryCaller(address common.Address, caller bind.ContractCaller) (*TokenRegistryCaller, error) {
	contract, err := bindTokenRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenRegistryCaller{contract: contract}, nil
}

// NewTokenRegistryTransactor creates a new write-only instance of TokenRegistry, bound to a specific deployed contract.
func NewTokenRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenRegistryTransactor, error) {
	contract, err := bindTokenRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenRegistryTransactor{contract: contract}, nil
}

// NewTokenRegistryFilterer creates a new log filterer instance of TokenRegistry, bound to a specific deployed contract.
func NewTokenRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenRegistryFilterer, error) {
	contract, err := bindTokenRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenRegistryFilterer{contract: contract}, nil
}

// bindTokenRegistry binds a generic wrapper to an already deployed contract.
func bindTokenRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenRegistry *TokenRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenRegistry.Contract.TokenRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenRegistry *TokenRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenRegistry.Contract.TokenRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenRegistry *TokenRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenRegistry.Contract.TokenRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenRegistry *TokenRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenRegistry *TokenRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenRegistry *TokenRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenRegistry.Contract.contract.Transact(opts, method, params...)
}

// NextTokenID is a free data retrieval call binding the contract method 0xf101e481.
//
// Solidity: function nextTokenID() view returns(uint256)
func (_TokenRegistry *TokenRegistryCaller) NextTokenID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TokenRegistry.contract.Call(opts, &out, "nextTokenID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextTokenID is a free data retrieval call binding the contract method 0xf101e481.
//
// Solidity: function nextTokenID() view returns(uint256)
func (_TokenRegistry *TokenRegistrySession) NextTokenID() (*big.Int, error) {
	return _TokenRegistry.Contract.NextTokenID(&_TokenRegistry.CallOpts)
}

// NextTokenID is a free data retrieval call binding the contract method 0xf101e481.
//
// Solidity: function nextTokenID() view returns(uint256)
func (_TokenRegistry *TokenRegistryCallerSession) NextTokenID() (*big.Int, error) {
	return _TokenRegistry.Contract.NextTokenID(&_TokenRegistry.CallOpts)
}

// SafeGetRecord is a free data retrieval call binding the contract method 0x0a7973b7.
//
// Solidity: function safeGetRecord(uint256 tokenID) view returns(address, uint256 l2Unit)
func (_TokenRegistry *TokenRegistryCaller) SafeGetRecord(opts *bind.CallOpts, tokenID *big.Int) (common.Address, *big.Int, error) {
	var out []interface{}
	err := _TokenRegistry.contract.Call(opts, &out, "safeGetRecord", tokenID)

	if err != nil {
		return *new(common.Address), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// SafeGetRecord is a free data retrieval call binding the contract method 0x0a7973b7.
//
// Solidity: function safeGetRecord(uint256 tokenID) view returns(address, uint256 l2Unit)
func (_TokenRegistry *TokenRegistrySession) SafeGetRecord(tokenID *big.Int) (common.Address, *big.Int, error) {
	return _TokenRegistry.Contract.SafeGetRecord(&_TokenRegistry.CallOpts, tokenID)
}

// SafeGetRecord is a free data retrieval call binding the contract method 0x0a7973b7.
//
// Solidity: function safeGetRecord(uint256 tokenID) view returns(address, uint256 l2Unit)
func (_TokenRegistry *TokenRegistryCallerSession) SafeGetRecord(tokenID *big.Int) (common.Address, *big.Int, error) {
	return _TokenRegistry.Contract.SafeGetRecord(&_TokenRegistry.CallOpts, tokenID)
}

// RegisterToken is a paid mutator transaction binding the contract method 0x09824a80.
//
// Solidity: function registerToken(address tokenContract) returns()
func (_TokenRegistry *TokenRegistryTransactor) RegisterToken(opts *bind.TransactOpts, tokenContract common.Address) (*types.Transaction, error) {
	return _TokenRegistry.contract.Transact(opts, "registerToken", tokenContract)
}

// RegisterToken is a paid mutator transaction binding the contract method 0x09824a80.
//
// Solidity: function registerToken(address tokenContract) returns()
func (_TokenRegistry *TokenRegistrySession) RegisterToken(tokenContract common.Address) (*types.Transaction, error) {
	return _TokenRegistry.Contract.RegisterToken(&_TokenRegistry.TransactOpts, tokenContract)
}

// RegisterToken is a paid mutator transaction binding the contract method 0x09824a80.
//
// Solidity: function registerToken(address tokenContract) returns()
func (_TokenRegistry *TokenRegistryTransactorSession) RegisterToken(tokenContract common.Address) (*types.Transaction, error) {
	return _TokenRegistry.Contract.RegisterToken(&_TokenRegistry.TransactOpts, tokenContract)
}

// TokenRegistryTokenRegisteredIterator is returned from FilterTokenRegistered and is used to iterate over the raw logs and unpacked data for TokenRegistered events raised by the TokenRegistry contract.
type TokenRegistryTokenRegisteredIterator struct {
	Event *TokenRegistryTokenRegistered // Event containing the contract specifics and raw log

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
func (it *TokenRegistryTokenRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenRegistryTokenRegistered)
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
		it.Event = new(TokenRegistryTokenRegistered)
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
func (it *TokenRegistryTokenRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenRegistryTokenRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenRegistryTokenRegistered represents a TokenRegistered event raised by the TokenRegistry contract.
type TokenRegistryTokenRegistered struct {
	TokenID       *big.Int
	TokenContract common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterTokenRegistered is a free log retrieval operation binding the contract event 0x5495d5c40ac716093397abf16b93fda5ccd72fb25488d6c41865b87898abea7c.
//
// Solidity: event TokenRegistered(uint256 tokenID, address tokenContract)
func (_TokenRegistry *TokenRegistryFilterer) FilterTokenRegistered(opts *bind.FilterOpts) (*TokenRegistryTokenRegisteredIterator, error) {

	logs, sub, err := _TokenRegistry.contract.FilterLogs(opts, "TokenRegistered")
	if err != nil {
		return nil, err
	}
	return &TokenRegistryTokenRegisteredIterator{contract: _TokenRegistry.contract, event: "TokenRegistered", logs: logs, sub: sub}, nil
}

// WatchTokenRegistered is a free log subscription operation binding the contract event 0x5495d5c40ac716093397abf16b93fda5ccd72fb25488d6c41865b87898abea7c.
//
// Solidity: event TokenRegistered(uint256 tokenID, address tokenContract)
func (_TokenRegistry *TokenRegistryFilterer) WatchTokenRegistered(opts *bind.WatchOpts, sink chan<- *TokenRegistryTokenRegistered) (event.Subscription, error) {

	logs, sub, err := _TokenRegistry.contract.WatchLogs(opts, "TokenRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenRegistryTokenRegistered)
				if err := _TokenRegistry.contract.UnpackLog(event, "TokenRegistered", log); err != nil {
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

// ParseTokenRegistered is a log parse operation binding the contract event 0x5495d5c40ac716093397abf16b93fda5ccd72fb25488d6c41865b87898abea7c.
//
// Solidity: event TokenRegistered(uint256 tokenID, address tokenContract)
func (_TokenRegistry *TokenRegistryFilterer) ParseTokenRegistered(log types.Log) (*TokenRegistryTokenRegistered, error) {
	event := new(TokenRegistryTokenRegistered)
	if err := _TokenRegistry.contract.UnpackLog(event, "TokenRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
