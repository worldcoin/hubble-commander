// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package estimator

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

// CostEstimatorMetaData contains all meta data concerning the CostEstimator contract.
var CostEstimatorMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"baseCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"pairCount\",\"type\":\"uint256\"}],\"name\":\"getGasCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"perPairCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"run\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506106bd806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80634e79f8ca146100515780639382255714610076578063c04062261461007f578063ebfd94b214610089575b600080fd5b61006461005f366004610614565b610092565b60405190815260200160405180910390f35b61006460005481565b6100876100b3565b005b61006460015481565b600080546001546100a39084610644565b6100ad919061062c565b92915050565b6100bb6100bd565b565b60006100c76100f5565b905060006100d3610328565b90506100df8282610663565b60018190556100ee9083610663565b6000555050565b6000806040518060c0016040528060018152602001600281526020017f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c281526020017f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed81526020017f090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b81526020017f12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa81525090506101b16105f6565b6000806107d05a6101c29190610663565b90506107d05a1161022e5760405162461bcd60e51b815260206004820152603d602482015260008051602061069183398151915260448201527f3a206e6f7420656e6f756768206761732c2073696e676c65207061697200000060648201526084015b60405180910390fd5b60005a905060208460c087600886fa925060005a61024c9083610663565b9050836102af5760405162461bcd60e51b815260206004820152603c602482015260008051602061069183398151915260448201527f3a2073696e676c6520706169722063616c6c206973206661696c6564000000006064820152608401610225565b84511561031e5760405162461bcd60e51b8152602060048201526043602482015260008051602061069183398151915260448201527f3a2073696e676c6520706169722063616c6c20726573756c74206d757374206260648201526206520360ec1b608482015260a401610225565b9695505050505050565b60008060405180610180016040528060018152602001600281526020017f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c281526020017f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed81526020017f090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b81526020017f12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa815260200160018152602001600281526020017f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c281526020017f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed81526020017f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec81526020017f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d815250905061048b6105f6565b6000806107d05a61049c9190610663565b90506107d05a116105035760405162461bcd60e51b815260206004820152603d602482015260008051602061069183398151915260448201527f3a206e6f7420656e6f756768206761732c20636f75706c6520706169720000006064820152608401610225565b60005a905060208461018087600886fa925060005a6105229083610663565b9050836105855760405162461bcd60e51b815260206004820152603c602482015260008051602061069183398151915260448201527f3a20636f75706c6520706169722063616c6c206973206661696c6564000000006064820152608401610225565b845160011461031e5760405162461bcd60e51b8152602060048201526043602482015260008051602061069183398151915260448201527f3a20636f75706c6520706169722063616c6c20726573756c74206d757374206260648201526265203160e81b608482015260a401610225565b60405180602001604052806001906020820280368337509192915050565b600060208284031215610625578081fd5b5035919050565b6000821982111561063f5761063f61067a565b500190565b600081600019048311821515161561065e5761065e61067a565b500290565b6000828210156106755761067561067a565b500390565b634e487b7160e01b600052601160045260246000fdfe424e50616972696e67507265636f6d70696c65436f7374457374696d61746f72a164736f6c6343000804000a",
}

// CostEstimatorABI is the input ABI used to generate the binding from.
// Deprecated: Use CostEstimatorMetaData.ABI instead.
var CostEstimatorABI = CostEstimatorMetaData.ABI

// CostEstimatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CostEstimatorMetaData.Bin instead.
var CostEstimatorBin = CostEstimatorMetaData.Bin

// DeployCostEstimator deploys a new Ethereum contract, binding an instance of CostEstimator to it.
func DeployCostEstimator(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CostEstimator, error) {
	parsed, err := CostEstimatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CostEstimatorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CostEstimator{CostEstimatorCaller: CostEstimatorCaller{contract: contract}, CostEstimatorTransactor: CostEstimatorTransactor{contract: contract}, CostEstimatorFilterer: CostEstimatorFilterer{contract: contract}}, nil
}

// CostEstimator is an auto generated Go binding around an Ethereum contract.
type CostEstimator struct {
	CostEstimatorCaller     // Read-only binding to the contract
	CostEstimatorTransactor // Write-only binding to the contract
	CostEstimatorFilterer   // Log filterer for contract events
}

// CostEstimatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type CostEstimatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CostEstimatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CostEstimatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CostEstimatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CostEstimatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CostEstimatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CostEstimatorSession struct {
	Contract     *CostEstimator    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CostEstimatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CostEstimatorCallerSession struct {
	Contract *CostEstimatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// CostEstimatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CostEstimatorTransactorSession struct {
	Contract     *CostEstimatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// CostEstimatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type CostEstimatorRaw struct {
	Contract *CostEstimator // Generic contract binding to access the raw methods on
}

// CostEstimatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CostEstimatorCallerRaw struct {
	Contract *CostEstimatorCaller // Generic read-only contract binding to access the raw methods on
}

// CostEstimatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CostEstimatorTransactorRaw struct {
	Contract *CostEstimatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCostEstimator creates a new instance of CostEstimator, bound to a specific deployed contract.
func NewCostEstimator(address common.Address, backend bind.ContractBackend) (*CostEstimator, error) {
	contract, err := bindCostEstimator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CostEstimator{CostEstimatorCaller: CostEstimatorCaller{contract: contract}, CostEstimatorTransactor: CostEstimatorTransactor{contract: contract}, CostEstimatorFilterer: CostEstimatorFilterer{contract: contract}}, nil
}

// NewCostEstimatorCaller creates a new read-only instance of CostEstimator, bound to a specific deployed contract.
func NewCostEstimatorCaller(address common.Address, caller bind.ContractCaller) (*CostEstimatorCaller, error) {
	contract, err := bindCostEstimator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CostEstimatorCaller{contract: contract}, nil
}

// NewCostEstimatorTransactor creates a new write-only instance of CostEstimator, bound to a specific deployed contract.
func NewCostEstimatorTransactor(address common.Address, transactor bind.ContractTransactor) (*CostEstimatorTransactor, error) {
	contract, err := bindCostEstimator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CostEstimatorTransactor{contract: contract}, nil
}

// NewCostEstimatorFilterer creates a new log filterer instance of CostEstimator, bound to a specific deployed contract.
func NewCostEstimatorFilterer(address common.Address, filterer bind.ContractFilterer) (*CostEstimatorFilterer, error) {
	contract, err := bindCostEstimator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CostEstimatorFilterer{contract: contract}, nil
}

// bindCostEstimator binds a generic wrapper to an already deployed contract.
func bindCostEstimator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(CostEstimatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CostEstimator *CostEstimatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CostEstimator.Contract.CostEstimatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CostEstimator *CostEstimatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CostEstimator.Contract.CostEstimatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CostEstimator *CostEstimatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CostEstimator.Contract.CostEstimatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CostEstimator *CostEstimatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CostEstimator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CostEstimator *CostEstimatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CostEstimator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CostEstimator *CostEstimatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CostEstimator.Contract.contract.Transact(opts, method, params...)
}

// BaseCost is a free data retrieval call binding the contract method 0x93822557.
//
// Solidity: function baseCost() view returns(uint256)
func (_CostEstimator *CostEstimatorCaller) BaseCost(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CostEstimator.contract.Call(opts, &out, "baseCost")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BaseCost is a free data retrieval call binding the contract method 0x93822557.
//
// Solidity: function baseCost() view returns(uint256)
func (_CostEstimator *CostEstimatorSession) BaseCost() (*big.Int, error) {
	return _CostEstimator.Contract.BaseCost(&_CostEstimator.CallOpts)
}

// BaseCost is a free data retrieval call binding the contract method 0x93822557.
//
// Solidity: function baseCost() view returns(uint256)
func (_CostEstimator *CostEstimatorCallerSession) BaseCost() (*big.Int, error) {
	return _CostEstimator.Contract.BaseCost(&_CostEstimator.CallOpts)
}

// GetGasCost is a free data retrieval call binding the contract method 0x4e79f8ca.
//
// Solidity: function getGasCost(uint256 pairCount) view returns(uint256)
func (_CostEstimator *CostEstimatorCaller) GetGasCost(opts *bind.CallOpts, pairCount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _CostEstimator.contract.Call(opts, &out, "getGasCost", pairCount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetGasCost is a free data retrieval call binding the contract method 0x4e79f8ca.
//
// Solidity: function getGasCost(uint256 pairCount) view returns(uint256)
func (_CostEstimator *CostEstimatorSession) GetGasCost(pairCount *big.Int) (*big.Int, error) {
	return _CostEstimator.Contract.GetGasCost(&_CostEstimator.CallOpts, pairCount)
}

// GetGasCost is a free data retrieval call binding the contract method 0x4e79f8ca.
//
// Solidity: function getGasCost(uint256 pairCount) view returns(uint256)
func (_CostEstimator *CostEstimatorCallerSession) GetGasCost(pairCount *big.Int) (*big.Int, error) {
	return _CostEstimator.Contract.GetGasCost(&_CostEstimator.CallOpts, pairCount)
}

// PerPairCost is a free data retrieval call binding the contract method 0xebfd94b2.
//
// Solidity: function perPairCost() view returns(uint256)
func (_CostEstimator *CostEstimatorCaller) PerPairCost(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CostEstimator.contract.Call(opts, &out, "perPairCost")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PerPairCost is a free data retrieval call binding the contract method 0xebfd94b2.
//
// Solidity: function perPairCost() view returns(uint256)
func (_CostEstimator *CostEstimatorSession) PerPairCost() (*big.Int, error) {
	return _CostEstimator.Contract.PerPairCost(&_CostEstimator.CallOpts)
}

// PerPairCost is a free data retrieval call binding the contract method 0xebfd94b2.
//
// Solidity: function perPairCost() view returns(uint256)
func (_CostEstimator *CostEstimatorCallerSession) PerPairCost() (*big.Int, error) {
	return _CostEstimator.Contract.PerPairCost(&_CostEstimator.CallOpts)
}

// Run is a paid mutator transaction binding the contract method 0xc0406226.
//
// Solidity: function run() returns()
func (_CostEstimator *CostEstimatorTransactor) Run(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CostEstimator.contract.Transact(opts, "run")
}

// Run is a paid mutator transaction binding the contract method 0xc0406226.
//
// Solidity: function run() returns()
func (_CostEstimator *CostEstimatorSession) Run() (*types.Transaction, error) {
	return _CostEstimator.Contract.Run(&_CostEstimator.TransactOpts)
}

// Run is a paid mutator transaction binding the contract method 0xc0406226.
//
// Solidity: function run() returns()
func (_CostEstimator *CostEstimatorTransactorSession) Run() (*types.Transaction, error) {
	return _CostEstimator.Contract.Run(&_CostEstimator.TransactOpts)
}
