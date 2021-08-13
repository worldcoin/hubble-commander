// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package generic

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

// TypesUserState is an auto generated low-level Go binding around an user-defined struct.
type TypesUserState struct {
	PubkeyID *big.Int
	TokenID  *big.Int
	Balance  *big.Int
	Nonce    *big.Int
}

// FrontendGenericMetaData contains all meta data concerning the FrontendGeneric contract.
var FrontendGenericMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"stateBytes\",\"type\":\"bytes\"}],\"name\":\"decodeState\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState\",\"name\":\"state\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState\",\"name\":\"state\",\"type\":\"tuple\"}],\"name\":\"encode\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506103b1806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806317412b8a1461003b578063b3b8362114610064575b600080fd5b61004e610049366004610250565b610084565b60405161005b9190610318565b60405180910390f35b6100776100723660046101e3565b6100a5565b60405161005b9190610379565b606061009d61009836849003840184610267565b6100f3565b90505b919050565b6100ad6101bb565b6100ec83838080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061018992505050565b9392505050565b805160609015801561010757506020820151155b801561011557506040820151155b801561012357506060820151155b1561015057600060405160200161013a919061036b565b60405160208183030381529060405290506100a0565b8151602080840151604080860151606087015191516101739594919291016102fd565b6040516020818303038152906040529050919050565b6101916101bb565b818060200190518101906101a591906102c8565b6060850152604084015260208301528152919050565b6040518060800160405280600081526020016000815260200160008152602001600081525090565b600080602083850312156101f5578182fd5b823567ffffffffffffffff8082111561020c578384fd5b818501915085601f83011261021f578384fd5b81358181111561022d578485fd5b86602082850101111561023e578485fd5b60209290920196919550909350505050565b600060808284031215610261578081fd5b50919050565b600060808284031215610278578081fd5b6040516080810181811067ffffffffffffffff82111715610297578283fd5b8060405250823581526020830135602082015260408301356040820152606083013560608201528091505092915050565b600080600080608085870312156102dd578182fd5b505082516020840151604085015160609095015191969095509092509050565b93845260208401929092526040830152606082015260800190565b6000602080835283518082850152825b8181101561034457858101830151858201604001528201610328565b818111156103555783604083870101525b50601f01601f1916929092016040019392505050565b60ff91909116815260200190565b815181526020808301519082015260408083015190820152606091820151918101919091526080019056fea164736f6c634300060c000a",
}

// FrontendGenericABI is the input ABI used to generate the binding from.
// Deprecated: Use FrontendGenericMetaData.ABI instead.
var FrontendGenericABI = FrontendGenericMetaData.ABI

// FrontendGenericBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use FrontendGenericMetaData.Bin instead.
var FrontendGenericBin = FrontendGenericMetaData.Bin

// DeployFrontendGeneric deploys a new Ethereum contract, binding an instance of FrontendGeneric to it.
func DeployFrontendGeneric(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *FrontendGeneric, error) {
	parsed, err := FrontendGenericMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FrontendGenericBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FrontendGeneric{FrontendGenericCaller: FrontendGenericCaller{contract: contract}, FrontendGenericTransactor: FrontendGenericTransactor{contract: contract}, FrontendGenericFilterer: FrontendGenericFilterer{contract: contract}}, nil
}

// FrontendGeneric is an auto generated Go binding around an Ethereum contract.
type FrontendGeneric struct {
	FrontendGenericCaller     // Read-only binding to the contract
	FrontendGenericTransactor // Write-only binding to the contract
	FrontendGenericFilterer   // Log filterer for contract events
}

// FrontendGenericCaller is an auto generated read-only Go binding around an Ethereum contract.
type FrontendGenericCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FrontendGenericTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FrontendGenericTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FrontendGenericFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FrontendGenericFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FrontendGenericSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FrontendGenericSession struct {
	Contract     *FrontendGeneric  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FrontendGenericCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FrontendGenericCallerSession struct {
	Contract *FrontendGenericCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// FrontendGenericTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FrontendGenericTransactorSession struct {
	Contract     *FrontendGenericTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// FrontendGenericRaw is an auto generated low-level Go binding around an Ethereum contract.
type FrontendGenericRaw struct {
	Contract *FrontendGeneric // Generic contract binding to access the raw methods on
}

// FrontendGenericCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FrontendGenericCallerRaw struct {
	Contract *FrontendGenericCaller // Generic read-only contract binding to access the raw methods on
}

// FrontendGenericTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FrontendGenericTransactorRaw struct {
	Contract *FrontendGenericTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFrontendGeneric creates a new instance of FrontendGeneric, bound to a specific deployed contract.
func NewFrontendGeneric(address common.Address, backend bind.ContractBackend) (*FrontendGeneric, error) {
	contract, err := bindFrontendGeneric(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FrontendGeneric{FrontendGenericCaller: FrontendGenericCaller{contract: contract}, FrontendGenericTransactor: FrontendGenericTransactor{contract: contract}, FrontendGenericFilterer: FrontendGenericFilterer{contract: contract}}, nil
}

// NewFrontendGenericCaller creates a new read-only instance of FrontendGeneric, bound to a specific deployed contract.
func NewFrontendGenericCaller(address common.Address, caller bind.ContractCaller) (*FrontendGenericCaller, error) {
	contract, err := bindFrontendGeneric(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FrontendGenericCaller{contract: contract}, nil
}

// NewFrontendGenericTransactor creates a new write-only instance of FrontendGeneric, bound to a specific deployed contract.
func NewFrontendGenericTransactor(address common.Address, transactor bind.ContractTransactor) (*FrontendGenericTransactor, error) {
	contract, err := bindFrontendGeneric(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FrontendGenericTransactor{contract: contract}, nil
}

// NewFrontendGenericFilterer creates a new log filterer instance of FrontendGeneric, bound to a specific deployed contract.
func NewFrontendGenericFilterer(address common.Address, filterer bind.ContractFilterer) (*FrontendGenericFilterer, error) {
	contract, err := bindFrontendGeneric(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FrontendGenericFilterer{contract: contract}, nil
}

// bindFrontendGeneric binds a generic wrapper to an already deployed contract.
func bindFrontendGeneric(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FrontendGenericABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FrontendGeneric *FrontendGenericRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FrontendGeneric.Contract.FrontendGenericCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FrontendGeneric *FrontendGenericRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FrontendGeneric.Contract.FrontendGenericTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FrontendGeneric *FrontendGenericRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FrontendGeneric.Contract.FrontendGenericTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FrontendGeneric *FrontendGenericCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FrontendGeneric.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FrontendGeneric *FrontendGenericTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FrontendGeneric.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FrontendGeneric *FrontendGenericTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FrontendGeneric.Contract.contract.Transact(opts, method, params...)
}

// DecodeState is a free data retrieval call binding the contract method 0xb3b83621.
//
// Solidity: function decodeState(bytes stateBytes) pure returns((uint256,uint256,uint256,uint256) state)
func (_FrontendGeneric *FrontendGenericCaller) DecodeState(opts *bind.CallOpts, stateBytes []byte) (TypesUserState, error) {
	var out []interface{}
	err := _FrontendGeneric.contract.Call(opts, &out, "decodeState", stateBytes)

	if err != nil {
		return *new(TypesUserState), err
	}

	out0 := *abi.ConvertType(out[0], new(TypesUserState)).(*TypesUserState)

	return out0, err

}

// DecodeState is a free data retrieval call binding the contract method 0xb3b83621.
//
// Solidity: function decodeState(bytes stateBytes) pure returns((uint256,uint256,uint256,uint256) state)
func (_FrontendGeneric *FrontendGenericSession) DecodeState(stateBytes []byte) (TypesUserState, error) {
	return _FrontendGeneric.Contract.DecodeState(&_FrontendGeneric.CallOpts, stateBytes)
}

// DecodeState is a free data retrieval call binding the contract method 0xb3b83621.
//
// Solidity: function decodeState(bytes stateBytes) pure returns((uint256,uint256,uint256,uint256) state)
func (_FrontendGeneric *FrontendGenericCallerSession) DecodeState(stateBytes []byte) (TypesUserState, error) {
	return _FrontendGeneric.Contract.DecodeState(&_FrontendGeneric.CallOpts, stateBytes)
}

// Encode is a free data retrieval call binding the contract method 0x17412b8a.
//
// Solidity: function encode((uint256,uint256,uint256,uint256) state) pure returns(bytes)
func (_FrontendGeneric *FrontendGenericCaller) Encode(opts *bind.CallOpts, state TypesUserState) ([]byte, error) {
	var out []interface{}
	err := _FrontendGeneric.contract.Call(opts, &out, "encode", state)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Encode is a free data retrieval call binding the contract method 0x17412b8a.
//
// Solidity: function encode((uint256,uint256,uint256,uint256) state) pure returns(bytes)
func (_FrontendGeneric *FrontendGenericSession) Encode(state TypesUserState) ([]byte, error) {
	return _FrontendGeneric.Contract.Encode(&_FrontendGeneric.CallOpts, state)
}

// Encode is a free data retrieval call binding the contract method 0x17412b8a.
//
// Solidity: function encode((uint256,uint256,uint256,uint256) state) pure returns(bytes)
func (_FrontendGeneric *FrontendGenericCallerSession) Encode(state TypesUserState) ([]byte, error) {
	return _FrontendGeneric.Contract.Encode(&_FrontendGeneric.CallOpts, state)
}
