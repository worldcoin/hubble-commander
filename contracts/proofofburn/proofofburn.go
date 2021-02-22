// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package proofofburn

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ProofOfBurnABI is the input ABI used to generate the binding from.
const ProofOfBurnABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"coordinator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getProposer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// ProofOfBurnBin is the compiled bytecode used for deploying new contracts.
var ProofOfBurnBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b03191633179055608a806100316000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80630a009097146037578063e9790d02146059575b600080fd5b603d605f565b604080516001600160a01b039092168252519081900360200190f35b603d606e565b6000546001600160a01b031681565b6000546001600160a01b03169056fea164736f6c634300060c000a"

// DeployProofOfBurn deploys a new Ethereum contract, binding an instance of ProofOfBurn to it.
func DeployProofOfBurn(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ProofOfBurn, error) {
	parsed, err := abi.JSON(strings.NewReader(ProofOfBurnABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ProofOfBurnBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ProofOfBurn{ProofOfBurnCaller: ProofOfBurnCaller{contract: contract}, ProofOfBurnTransactor: ProofOfBurnTransactor{contract: contract}, ProofOfBurnFilterer: ProofOfBurnFilterer{contract: contract}}, nil
}

// ProofOfBurn is an auto generated Go binding around an Ethereum contract.
type ProofOfBurn struct {
	ProofOfBurnCaller     // Read-only binding to the contract
	ProofOfBurnTransactor // Write-only binding to the contract
	ProofOfBurnFilterer   // Log filterer for contract events
}

// ProofOfBurnCaller is an auto generated read-only Go binding around an Ethereum contract.
type ProofOfBurnCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofOfBurnTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ProofOfBurnTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofOfBurnFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ProofOfBurnFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofOfBurnSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ProofOfBurnSession struct {
	Contract     *ProofOfBurn      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProofOfBurnCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ProofOfBurnCallerSession struct {
	Contract *ProofOfBurnCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// ProofOfBurnTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ProofOfBurnTransactorSession struct {
	Contract     *ProofOfBurnTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// ProofOfBurnRaw is an auto generated low-level Go binding around an Ethereum contract.
type ProofOfBurnRaw struct {
	Contract *ProofOfBurn // Generic contract binding to access the raw methods on
}

// ProofOfBurnCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ProofOfBurnCallerRaw struct {
	Contract *ProofOfBurnCaller // Generic read-only contract binding to access the raw methods on
}

// ProofOfBurnTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ProofOfBurnTransactorRaw struct {
	Contract *ProofOfBurnTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProofOfBurn creates a new instance of ProofOfBurn, bound to a specific deployed contract.
func NewProofOfBurn(address common.Address, backend bind.ContractBackend) (*ProofOfBurn, error) {
	contract, err := bindProofOfBurn(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ProofOfBurn{ProofOfBurnCaller: ProofOfBurnCaller{contract: contract}, ProofOfBurnTransactor: ProofOfBurnTransactor{contract: contract}, ProofOfBurnFilterer: ProofOfBurnFilterer{contract: contract}}, nil
}

// NewProofOfBurnCaller creates a new read-only instance of ProofOfBurn, bound to a specific deployed contract.
func NewProofOfBurnCaller(address common.Address, caller bind.ContractCaller) (*ProofOfBurnCaller, error) {
	contract, err := bindProofOfBurn(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProofOfBurnCaller{contract: contract}, nil
}

// NewProofOfBurnTransactor creates a new write-only instance of ProofOfBurn, bound to a specific deployed contract.
func NewProofOfBurnTransactor(address common.Address, transactor bind.ContractTransactor) (*ProofOfBurnTransactor, error) {
	contract, err := bindProofOfBurn(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProofOfBurnTransactor{contract: contract}, nil
}

// NewProofOfBurnFilterer creates a new log filterer instance of ProofOfBurn, bound to a specific deployed contract.
func NewProofOfBurnFilterer(address common.Address, filterer bind.ContractFilterer) (*ProofOfBurnFilterer, error) {
	contract, err := bindProofOfBurn(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProofOfBurnFilterer{contract: contract}, nil
}

// bindProofOfBurn binds a generic wrapper to an already deployed contract.
func bindProofOfBurn(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ProofOfBurnABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProofOfBurn *ProofOfBurnRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ProofOfBurn.Contract.ProofOfBurnCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProofOfBurn *ProofOfBurnRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProofOfBurn.Contract.ProofOfBurnTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProofOfBurn *ProofOfBurnRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProofOfBurn.Contract.ProofOfBurnTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProofOfBurn *ProofOfBurnCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ProofOfBurn.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProofOfBurn *ProofOfBurnTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProofOfBurn.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProofOfBurn *ProofOfBurnTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProofOfBurn.Contract.contract.Transact(opts, method, params...)
}

// Coordinator is a free data retrieval call binding the contract method 0x0a009097.
//
// Solidity: function coordinator() view returns(address)
func (_ProofOfBurn *ProofOfBurnCaller) Coordinator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ProofOfBurn.contract.Call(opts, &out, "coordinator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Coordinator is a free data retrieval call binding the contract method 0x0a009097.
//
// Solidity: function coordinator() view returns(address)
func (_ProofOfBurn *ProofOfBurnSession) Coordinator() (common.Address, error) {
	return _ProofOfBurn.Contract.Coordinator(&_ProofOfBurn.CallOpts)
}

// Coordinator is a free data retrieval call binding the contract method 0x0a009097.
//
// Solidity: function coordinator() view returns(address)
func (_ProofOfBurn *ProofOfBurnCallerSession) Coordinator() (common.Address, error) {
	return _ProofOfBurn.Contract.Coordinator(&_ProofOfBurn.CallOpts)
}

// GetProposer is a free data retrieval call binding the contract method 0xe9790d02.
//
// Solidity: function getProposer() view returns(address)
func (_ProofOfBurn *ProofOfBurnCaller) GetProposer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ProofOfBurn.contract.Call(opts, &out, "getProposer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetProposer is a free data retrieval call binding the contract method 0xe9790d02.
//
// Solidity: function getProposer() view returns(address)
func (_ProofOfBurn *ProofOfBurnSession) GetProposer() (common.Address, error) {
	return _ProofOfBurn.Contract.GetProposer(&_ProofOfBurn.CallOpts)
}

// GetProposer is a free data retrieval call binding the contract method 0xe9790d02.
//
// Solidity: function getProposer() view returns(address)
func (_ProofOfBurn *ProofOfBurnCallerSession) GetProposer() (common.Address, error) {
	return _ProofOfBurn.Contract.GetProposer(&_ProofOfBurn.CallOpts)
}
