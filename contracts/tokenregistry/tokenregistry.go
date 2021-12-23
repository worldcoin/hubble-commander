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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenContract\",\"type\":\"address\"}],\"name\":\"RegisteredToken\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenContract\",\"type\":\"address\"}],\"name\":\"RegistrationRequest\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenContract\",\"type\":\"address\"}],\"name\":\"finaliseRegistration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextTokenID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"pendingRegistrations\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenContract\",\"type\":\"address\"}],\"name\":\"requestRegistration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"}],\"name\":\"safeGetRecord\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"l2Unit\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052600060025534801561001557600080fd5b5061054c806100256000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c80630a7973b71461005c5780630b72ccbc1461009c578063cd4852e5146100c4578063e9b6bfe3146100fe578063f101e48114610124575b600080fd5b6100796004803603602081101561007257600080fd5b503561013e565b604080516001600160a01b03909316835260208301919091528051918290030190f35b6100c2600480360360208110156100b257600080fd5b50356001600160a01b03166101df565b005b6100ea600480360360208110156100da57600080fd5b50356001600160a01b0316610366565b604080519115158252519081900360200190f35b6100c26004803603602081101561011457600080fd5b50356001600160a01b031661037b565b61012c6104ff565b60408051918252519081900360200190f35b600080610149610505565b506000838152600160209081526040918290208251808401909352546001600160a01b038116808452600160a01b90910460ff161515918301919091526101c15760405162461bcd60e51b815260040180806020018281038252602381526020018061051d6023913960400191505060405180910390fd5b600191508060200151156101d757633b9aca0091505b519150915091565b6001600160a01b03811660009081526020819052604090205460ff1661024c576040805162461bcd60e51b815260206004820152601860248201527f546f6b656e20776173206e6f7420726567697374657265640000000000000000604482015290519081900360640190fd5b60006009826001600160a01b031663313ce5676040518163ffffffff1660e01b815260040160206040518083038186803b15801561028957600080fd5b505afa15801561029d573d6000803e3d6000fd5b505050506040513d60208110156102b357600080fd5b50516040805180820182526001600160a01b0386811680835260ff90941694909410156020808301828152600280546000908152600184528690209451855492516001600160a01b031990931698169790971760ff60a01b1916600160a01b9115159190910217909255935482519081529081019290925280519293507f5dbaa701a7acef513f72a61799f7e50f4653f462b9f780d88d1b9bec89de216892918290030190a15050600280546001019055565b60006020819052908152604090205460ff1681565b6001600160a01b03811660009081526020819052604090205460ff16156103e9576040805162461bcd60e51b815260206004820152601960248201527f546f6b656e20616c726561647920726567697374657265642e00000000000000604482015290519081900360640190fd5b6012816001600160a01b031663313ce5676040518163ffffffff1660e01b815260040160206040518083038186803b15801561042457600080fd5b505afa158015610438573d6000803e3d6000fd5b505050506040513d602081101561044e57600080fd5b505160ff1611156104a6576040805162461bcd60e51b815260206004820152601960248201527f446f6e277420736572766520646563696d616c73203e20313800000000000000604482015290519081900360640190fd5b6001600160a01b03811660008181526020818152604091829020805460ff19166001179055815192835290517fdc79fc57451962cfe3916e686997a49229af75ce2055deb4c0f0fdf3d5d2e7c19281900390910190a150565b60025481565b60408051808201909152600080825260208201529056fe546f6b656e52656769737472793a20556e7265676973746572656420746f6b656e4944a164736f6c634300060c000a",
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

// PendingRegistrations is a free data retrieval call binding the contract method 0xcd4852e5.
//
// Solidity: function pendingRegistrations(address ) view returns(bool)
func (_TokenRegistry *TokenRegistryCaller) PendingRegistrations(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _TokenRegistry.contract.Call(opts, &out, "pendingRegistrations", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// PendingRegistrations is a free data retrieval call binding the contract method 0xcd4852e5.
//
// Solidity: function pendingRegistrations(address ) view returns(bool)
func (_TokenRegistry *TokenRegistrySession) PendingRegistrations(arg0 common.Address) (bool, error) {
	return _TokenRegistry.Contract.PendingRegistrations(&_TokenRegistry.CallOpts, arg0)
}

// PendingRegistrations is a free data retrieval call binding the contract method 0xcd4852e5.
//
// Solidity: function pendingRegistrations(address ) view returns(bool)
func (_TokenRegistry *TokenRegistryCallerSession) PendingRegistrations(arg0 common.Address) (bool, error) {
	return _TokenRegistry.Contract.PendingRegistrations(&_TokenRegistry.CallOpts, arg0)
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

// FinaliseRegistration is a paid mutator transaction binding the contract method 0x0b72ccbc.
//
// Solidity: function finaliseRegistration(address tokenContract) returns()
func (_TokenRegistry *TokenRegistryTransactor) FinaliseRegistration(opts *bind.TransactOpts, tokenContract common.Address) (*types.Transaction, error) {
	return _TokenRegistry.contract.Transact(opts, "finaliseRegistration", tokenContract)
}

// FinaliseRegistration is a paid mutator transaction binding the contract method 0x0b72ccbc.
//
// Solidity: function finaliseRegistration(address tokenContract) returns()
func (_TokenRegistry *TokenRegistrySession) FinaliseRegistration(tokenContract common.Address) (*types.Transaction, error) {
	return _TokenRegistry.Contract.FinaliseRegistration(&_TokenRegistry.TransactOpts, tokenContract)
}

// FinaliseRegistration is a paid mutator transaction binding the contract method 0x0b72ccbc.
//
// Solidity: function finaliseRegistration(address tokenContract) returns()
func (_TokenRegistry *TokenRegistryTransactorSession) FinaliseRegistration(tokenContract common.Address) (*types.Transaction, error) {
	return _TokenRegistry.Contract.FinaliseRegistration(&_TokenRegistry.TransactOpts, tokenContract)
}

// RequestRegistration is a paid mutator transaction binding the contract method 0xe9b6bfe3.
//
// Solidity: function requestRegistration(address tokenContract) returns()
func (_TokenRegistry *TokenRegistryTransactor) RequestRegistration(opts *bind.TransactOpts, tokenContract common.Address) (*types.Transaction, error) {
	return _TokenRegistry.contract.Transact(opts, "requestRegistration", tokenContract)
}

// RequestRegistration is a paid mutator transaction binding the contract method 0xe9b6bfe3.
//
// Solidity: function requestRegistration(address tokenContract) returns()
func (_TokenRegistry *TokenRegistrySession) RequestRegistration(tokenContract common.Address) (*types.Transaction, error) {
	return _TokenRegistry.Contract.RequestRegistration(&_TokenRegistry.TransactOpts, tokenContract)
}

// RequestRegistration is a paid mutator transaction binding the contract method 0xe9b6bfe3.
//
// Solidity: function requestRegistration(address tokenContract) returns()
func (_TokenRegistry *TokenRegistryTransactorSession) RequestRegistration(tokenContract common.Address) (*types.Transaction, error) {
	return _TokenRegistry.Contract.RequestRegistration(&_TokenRegistry.TransactOpts, tokenContract)
}

// TokenRegistryRegisteredTokenIterator is returned from FilterRegisteredToken and is used to iterate over the raw logs and unpacked data for RegisteredToken events raised by the TokenRegistry contract.
type TokenRegistryRegisteredTokenIterator struct {
	Event *TokenRegistryRegisteredToken // Event containing the contract specifics and raw log

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
func (it *TokenRegistryRegisteredTokenIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenRegistryRegisteredToken)
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
		it.Event = new(TokenRegistryRegisteredToken)
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
func (it *TokenRegistryRegisteredTokenIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenRegistryRegisteredTokenIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenRegistryRegisteredToken represents a RegisteredToken event raised by the TokenRegistry contract.
type TokenRegistryRegisteredToken struct {
	TokenID       *big.Int
	TokenContract common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterRegisteredToken is a free log retrieval operation binding the contract event 0x5dbaa701a7acef513f72a61799f7e50f4653f462b9f780d88d1b9bec89de2168.
//
// Solidity: event RegisteredToken(uint256 tokenID, address tokenContract)
func (_TokenRegistry *TokenRegistryFilterer) FilterRegisteredToken(opts *bind.FilterOpts) (*TokenRegistryRegisteredTokenIterator, error) {

	logs, sub, err := _TokenRegistry.contract.FilterLogs(opts, "RegisteredToken")
	if err != nil {
		return nil, err
	}
	return &TokenRegistryRegisteredTokenIterator{contract: _TokenRegistry.contract, event: "RegisteredToken", logs: logs, sub: sub}, nil
}

// WatchRegisteredToken is a free log subscription operation binding the contract event 0x5dbaa701a7acef513f72a61799f7e50f4653f462b9f780d88d1b9bec89de2168.
//
// Solidity: event RegisteredToken(uint256 tokenID, address tokenContract)
func (_TokenRegistry *TokenRegistryFilterer) WatchRegisteredToken(opts *bind.WatchOpts, sink chan<- *TokenRegistryRegisteredToken) (event.Subscription, error) {

	logs, sub, err := _TokenRegistry.contract.WatchLogs(opts, "RegisteredToken")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenRegistryRegisteredToken)
				if err := _TokenRegistry.contract.UnpackLog(event, "RegisteredToken", log); err != nil {
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

// ParseRegisteredToken is a log parse operation binding the contract event 0x5dbaa701a7acef513f72a61799f7e50f4653f462b9f780d88d1b9bec89de2168.
//
// Solidity: event RegisteredToken(uint256 tokenID, address tokenContract)
func (_TokenRegistry *TokenRegistryFilterer) ParseRegisteredToken(log types.Log) (*TokenRegistryRegisteredToken, error) {
	event := new(TokenRegistryRegisteredToken)
	if err := _TokenRegistry.contract.UnpackLog(event, "RegisteredToken", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenRegistryRegistrationRequestIterator is returned from FilterRegistrationRequest and is used to iterate over the raw logs and unpacked data for RegistrationRequest events raised by the TokenRegistry contract.
type TokenRegistryRegistrationRequestIterator struct {
	Event *TokenRegistryRegistrationRequest // Event containing the contract specifics and raw log

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
func (it *TokenRegistryRegistrationRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenRegistryRegistrationRequest)
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
		it.Event = new(TokenRegistryRegistrationRequest)
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
func (it *TokenRegistryRegistrationRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenRegistryRegistrationRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenRegistryRegistrationRequest represents a RegistrationRequest event raised by the TokenRegistry contract.
type TokenRegistryRegistrationRequest struct {
	TokenContract common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterRegistrationRequest is a free log retrieval operation binding the contract event 0xdc79fc57451962cfe3916e686997a49229af75ce2055deb4c0f0fdf3d5d2e7c1.
//
// Solidity: event RegistrationRequest(address tokenContract)
func (_TokenRegistry *TokenRegistryFilterer) FilterRegistrationRequest(opts *bind.FilterOpts) (*TokenRegistryRegistrationRequestIterator, error) {

	logs, sub, err := _TokenRegistry.contract.FilterLogs(opts, "RegistrationRequest")
	if err != nil {
		return nil, err
	}
	return &TokenRegistryRegistrationRequestIterator{contract: _TokenRegistry.contract, event: "RegistrationRequest", logs: logs, sub: sub}, nil
}

// WatchRegistrationRequest is a free log subscription operation binding the contract event 0xdc79fc57451962cfe3916e686997a49229af75ce2055deb4c0f0fdf3d5d2e7c1.
//
// Solidity: event RegistrationRequest(address tokenContract)
func (_TokenRegistry *TokenRegistryFilterer) WatchRegistrationRequest(opts *bind.WatchOpts, sink chan<- *TokenRegistryRegistrationRequest) (event.Subscription, error) {

	logs, sub, err := _TokenRegistry.contract.WatchLogs(opts, "RegistrationRequest")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenRegistryRegistrationRequest)
				if err := _TokenRegistry.contract.UnpackLog(event, "RegistrationRequest", log); err != nil {
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

// ParseRegistrationRequest is a log parse operation binding the contract event 0xdc79fc57451962cfe3916e686997a49229af75ce2055deb4c0f0fdf3d5d2e7c1.
//
// Solidity: event RegistrationRequest(address tokenContract)
func (_TokenRegistry *TokenRegistryFilterer) ParseRegistrationRequest(log types.Log) (*TokenRegistryRegistrationRequest, error) {
	event := new(TokenRegistryRegistrationRequest)
	if err := _TokenRegistry.contract.UnpackLog(event, "RegistrationRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
