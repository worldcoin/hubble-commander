// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package estimator

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

// CostEstimatorABI is the input ABI used to generate the binding from.
const CostEstimatorABI = "[{\"inputs\":[],\"name\":\"baseCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"pairCount\",\"type\":\"uint256\"}],\"name\":\"getGasCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"perPairCost\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"run\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// CostEstimatorBin is the compiled bytecode used for deploying new contracts.
var CostEstimatorBin = "0x608060405234801561001057600080fd5b506106b7806100206000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80634e79f8ca146100515780639382255714610080578063c040622614610088578063ebfd94b214610092575b600080fd5b61006e6004803603602081101561006757600080fd5b503561009a565b60408051918252519081900360200190f35b61006e6100a8565b6100906100ae565b005b61006e6100b8565b600054600154820201919050565b60005481565b6100b66100be565b565b60015481565b60006100c86100e5565b905060006100d461028e565b829003600181905590910360005550565b60006100ef6104d7565b6040518060c0016040528060018152602001600281526020017f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c281526020017f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed81526020017f090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b81526020017f12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa81525090506101a86104f5565b6000806107d05a0390506107d05a116101f25760405162461bcd60e51b815260040180806020018281038252603d8152602001806105b3603d913960400191505060405180910390fd5b60005a905060208460c087600886fa925060005a82039050836102465760405162461bcd60e51b815260040180806020018281038252603c815260200180610633603c913960400191505060405180910390fd5b8451156102845760405162461bcd60e51b81526004018080602001828103825260438152602001806105706043913960600191505060405180910390fd5b9550505050505090565b6000610298610513565b60405180610180016040528060018152602001600281526020017f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c281526020017f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed81526020017f090689d0585ff075ec9e99ad690c3395bc4b313370b38ef355acdadcd122975b81526020017f12c85ea5db8c6deb4aab71808dcb408fe3d1e7690c43d37b4ce6cc0166fa7daa815260200160018152602001600281526020017f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c281526020017f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed81526020017f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec81526020017f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d81525090506103f86104f5565b6000806107d05a0390506107d05a116104425760405162461bcd60e51b815260040180806020018281038252603d815260200180610533603d913960400191505060405180910390fd5b60005a905060208461018087600886fa925060005a82039050836104975760405162461bcd60e51b815260040180806020018281038252603c81526020018061066f603c913960400191505060405180910390fd5b84516001146102845760405162461bcd60e51b81526004018080602001828103825260438152602001806105f06043913960600191505060405180910390fd5b6040518060c001604052806006906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b604051806101800160405280600c90602082028036833750919291505056fe424e50616972696e67507265636f6d70696c65436f7374457374696d61746f723a206e6f7420656e6f756768206761732c20636f75706c652070616972424e50616972696e67507265636f6d70696c65436f7374457374696d61746f723a2073696e676c6520706169722063616c6c20726573756c74206d7573742062652030424e50616972696e67507265636f6d70696c65436f7374457374696d61746f723a206e6f7420656e6f756768206761732c2073696e676c652070616972424e50616972696e67507265636f6d70696c65436f7374457374696d61746f723a20636f75706c6520706169722063616c6c20726573756c74206d7573742062652031424e50616972696e67507265636f6d70696c65436f7374457374696d61746f723a2073696e676c6520706169722063616c6c206973206661696c6564424e50616972696e67507265636f6d70696c65436f7374457374696d61746f723a20636f75706c6520706169722063616c6c206973206661696c6564a164736f6c634300060c000a"

// DeployCostEstimator deploys a new Ethereum contract, binding an instance of CostEstimator to it.
func DeployCostEstimator(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CostEstimator, error) {
	parsed, err := abi.JSON(strings.NewReader(CostEstimatorABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(CostEstimatorBin), backend)
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
