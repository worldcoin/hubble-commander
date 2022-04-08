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
	Bin: "0x6080604052600060025534801561001557600080fd5b5061046a806100256000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c806309824a80146100465780630a7973b71461005b578063f101e48114610092575b600080fd5b6100596100543660046103cf565b6100a9565b005b61006e6100693660046103fd565b61031d565b604080516001600160a01b0390931683526020830191909152015b60405180910390f35b61009b60025481565b604051908152602001610089565b6001600160a01b03811660009081526001602052604090205460ff16156101175760405162461bcd60e51b815260206004820152601960248201527f546f6b656e20616c726561647920726567697374657265642e0000000000000060448201526064015b60405180910390fd5b6012816001600160a01b031663313ce5676040518163ffffffff1660e01b815260040160206040518083038186803b15801561015257600080fd5b505afa158015610166573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061018a9190610415565b60ff1611156101db5760405162461bcd60e51b815260206004820152601960248201527f446f6e277420736572766520646563696d616c73203e20313800000000000000604482015260640161010e565b60006009826001600160a01b031663313ce5676040518163ffffffff1660e01b815260040160206040518083038186803b15801561021857600080fd5b505afa15801561022c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102509190610415565b6040805180820182526001600160a01b0386811680835260ff94909416949094101560208083018281526002805460009081528084528681209551865493511515600160a01b026001600160a81b031990941699169890981791909117909355848652600180825295849020805460ff191690961790955590548251908152938401929092529092507f5495d5c40ac716093397abf16b93fda5ccd72fb25488d6c41865b87898abea7c910160405180910390a16002805490600061031483610436565b91905055505050565b6000818152602081815260408083208151808301909252546001600160a01b038116808352600160a01b90910460ff1615159282019290925282916103b05760405162461bcd60e51b815260206004820152602360248201527f546f6b656e52656769737472793a20556e7265676973746572656420746f6b656044820152621b925160ea1b606482015260840161010e565b600191508060200151156103c657633b9aca0091505b51939092509050565b6000602082840312156103e0578081fd5b81356001600160a01b03811681146103f6578182fd5b9392505050565b60006020828403121561040e578081fd5b5035919050565b600060208284031215610426578081fd5b815160ff811681146103f6578182fd5b600060001982141561045657634e487b7160e01b81526011600452602481fd5b506001019056fea164736f6c6343000804000a",
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
