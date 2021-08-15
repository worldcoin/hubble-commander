// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package chooser

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

// ChooserMetaData contains all meta data concerning the Chooser contract.
var ChooserMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getProposer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"proposer\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// ChooserABI is the input ABI used to generate the binding from.
// Deprecated: Use ChooserMetaData.ABI instead.
var ChooserABI = ChooserMetaData.ABI

// Chooser is an auto generated Go binding around an Ethereum contract.
type Chooser struct {
	ChooserCaller     // Read-only binding to the contract
	ChooserTransactor // Write-only binding to the contract
	ChooserFilterer   // Log filterer for contract events
}

// ChooserCaller is an auto generated read-only Go binding around an Ethereum contract.
type ChooserCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChooserTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ChooserTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChooserFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ChooserFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ChooserSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ChooserSession struct {
	Contract     *Chooser          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ChooserCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ChooserCallerSession struct {
	Contract *ChooserCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ChooserTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ChooserTransactorSession struct {
	Contract     *ChooserTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ChooserRaw is an auto generated low-level Go binding around an Ethereum contract.
type ChooserRaw struct {
	Contract *Chooser // Generic contract binding to access the raw methods on
}

// ChooserCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ChooserCallerRaw struct {
	Contract *ChooserCaller // Generic read-only contract binding to access the raw methods on
}

// ChooserTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ChooserTransactorRaw struct {
	Contract *ChooserTransactor // Generic write-only contract binding to access the raw methods on
}

// NewChooser creates a new instance of Chooser, bound to a specific deployed contract.
func NewChooser(address common.Address, backend bind.ContractBackend) (*Chooser, error) {
	contract, err := bindChooser(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Chooser{ChooserCaller: ChooserCaller{contract: contract}, ChooserTransactor: ChooserTransactor{contract: contract}, ChooserFilterer: ChooserFilterer{contract: contract}}, nil
}

// NewChooserCaller creates a new read-only instance of Chooser, bound to a specific deployed contract.
func NewChooserCaller(address common.Address, caller bind.ContractCaller) (*ChooserCaller, error) {
	contract, err := bindChooser(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ChooserCaller{contract: contract}, nil
}

// NewChooserTransactor creates a new write-only instance of Chooser, bound to a specific deployed contract.
func NewChooserTransactor(address common.Address, transactor bind.ContractTransactor) (*ChooserTransactor, error) {
	contract, err := bindChooser(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ChooserTransactor{contract: contract}, nil
}

// NewChooserFilterer creates a new log filterer instance of Chooser, bound to a specific deployed contract.
func NewChooserFilterer(address common.Address, filterer bind.ContractFilterer) (*ChooserFilterer, error) {
	contract, err := bindChooser(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ChooserFilterer{contract: contract}, nil
}

// bindChooser binds a generic wrapper to an already deployed contract.
func bindChooser(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ChooserABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Chooser *ChooserRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Chooser.Contract.ChooserCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Chooser *ChooserRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chooser.Contract.ChooserTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Chooser *ChooserRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Chooser.Contract.ChooserTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Chooser *ChooserCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Chooser.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Chooser *ChooserTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Chooser.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Chooser *ChooserTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Chooser.Contract.contract.Transact(opts, method, params...)
}

// GetProposer is a free data retrieval call binding the contract method 0xe9790d02.
//
// Solidity: function getProposer() view returns(address proposer)
func (_Chooser *ChooserCaller) GetProposer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Chooser.contract.Call(opts, &out, "getProposer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetProposer is a free data retrieval call binding the contract method 0xe9790d02.
//
// Solidity: function getProposer() view returns(address proposer)
func (_Chooser *ChooserSession) GetProposer() (common.Address, error) {
	return _Chooser.Contract.GetProposer(&_Chooser.CallOpts)
}

// GetProposer is a free data retrieval call binding the contract method 0xe9790d02.
//
// Solidity: function getProposer() view returns(address proposer)
func (_Chooser *ChooserCallerSession) GetProposer() (common.Address, error) {
	return _Chooser.Contract.GetProposer(&_Chooser.CallOpts)
}
