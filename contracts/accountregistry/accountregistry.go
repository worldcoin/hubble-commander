// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package accountregistry

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

// AccountRegistryMetaData contains all meta data concerning the AccountRegistry contract.
var AccountRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractChooser\",\"name\":\"_chooser\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"rootLeft\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"leafIndexLeft\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[31]\",\"name\":\"filledSubtreesLeft\",\"type\":\"bytes32[31]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"endID\",\"type\":\"uint256\"}],\"name\":\"BatchPubkeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"}],\"name\":\"SinglePubkeyRegistered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BATCH_DEPTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"BATCH_SIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEPTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SET_SIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WITNESS_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"chooser\",\"outputs\":[{\"internalType\":\"contractChooser\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256[4]\",\"name\":\"pubkey\",\"type\":\"uint256[4]\"},{\"internalType\":\"bytes32[31]\",\"name\":\"witness\",\"type\":\"bytes32[31]\"}],\"name\":\"exists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"filledSubtreesLeft\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"filledSubtreesRight\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"leafIndexLeft\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"leafIndexRight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[4]\",\"name\":\"pubkey\",\"type\":\"uint256[4]\"}],\"name\":\"register\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[4][16]\",\"name\":\"pubkeys\",\"type\":\"uint256[4][16]\"}],\"name\":\"registerBatch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"root\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rootLeft\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rootRight\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"zeros\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a06040526000600355600060045534801561001a57600080fd5b50604051610cf4380380610cf4833981810160405261044081101561003e57600080fd5b5080516020820151604083015191929091906060018282827f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e5638060056000015560015b601f81101561012a57600560018203601f811061009a57fe5b0154600560018303601f81106100ac57fe5b0154604051602001808381526020018281526020019250505060405160208183030381529060405280519060200120600582601f81106100e857fe5b0155600481108015906100fb575080601f115b1561012257600581601f811061010d57fe5b0154604360048303601b811061011f57fe5b01555b600101610081565b50610138602483601f6101c0565b50505060035560008190556023546040805160208082018490528183019390935281518082038301815260608083018452815191850191909120600181905560808084019690965260a0808401919091528351808403909101815260c090920190925280519201919091206002559490941b6001600160601b03191690935250610213915050565b82601f81019282156101ee579160200282015b828111156101ee5782518255916020019190600101906101d3565b506101fa9291506101fe565b5090565b5b808211156101fa57600081556001016101ff565b60805160601c610abf61023560003980610272528061051c5250610abf6000f3fe608060405234801561001057600080fd5b506004361061010b5760003560e01c806395e4bf03116100a2578063d0383d6811610071578063d0383d681461021d578063d7c53ea714610225578063d82894631461022d578063e829558814610235578063ebf0c717146102525761010b565b806395e4bf03146101d557806398366e351461016c57806398d17621146101f1578063cab2da9b146102155761010b565b80635e71468b116100de5780635e71468b1461016c578063693c1db714610174578063709a8b2a1461017c5780638d037962146101b85761010b565b8063034a29ae146101105780631c4a7a941461013f5780631c76e77e1461015c57806349faa4d414610164575b600080fd5b61012d6004803603602081101561012657600080fd5b503561025a565b60408051918252519081900360200190f35b61012d600480360361080081101561015657600080fd5b5061026e565b61012d610402565b61012d610407565b61012d61040c565b61012d610411565b6101a4600480360361048081101561019357600080fd5b50803590602081019060a001610417565b604080519115158252519081900360200190f35b61012d600480360360208110156101ce57600080fd5b503561048d565b61012d600480360360808110156101eb57600080fd5b5061049a565b6101f961051a565b604080516001600160a01b039092168252519081900360200190f35b61012d61053e565b61012d610544565b61012d61054c565b61012d610552565b61012d6004803603602081101561024b57600080fd5b5035610558565b61012d610565565b602481601f811061026757fe5b0154905081565b60007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b1580156102c957600080fd5b505afa1580156102dd573d6000803e3d6000fd5b505050506040513d60208110156102f357600080fd5b50516001600160a01b0316331461033b5760405162461bcd60e51b8152600401808060200182810382526024815260200180610a8f6024913960400191505060405180910390fd5b610343610a50565b60005b60108110156103af57600084826010811061035d57fe5b60800201604051602001808260046020028082843780830192505050915050604051602081830303815290604052805190602001209050808383601081106103a157fe5b602002015250600101610346565b5060006103bb8261056b565b60408051828152600f8301602082015281519293507f3154b80a7d9f6a143c37dde575f47deb78dacc7f280d8efc7e3ae102758a841b929081900390910190a19392505050565b600481565b601081565b601f81565b60035481565b60008083604051602001808260046020028082843780830192505050915050604051602081830303815290604052805190602001209050610482818685601f806020026040519081016040528092919082601f60200280828437600092019190915250610805915050565b9150505b9392505050565b604381601b811061026757fe5b6000808260405160200180826004602002808284378083019250505091505060405160208183030381529060405280519060200120905060006104dc826108e1565b6040805182815290519192507f59056afed767866d7c194cf26e24ebe16974ef943cccb729452e9adc265a9ac7919081900360200190a19392505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b60015481565b638000000081565b60045481565b60005481565b600581601f811061026757fe5b60025481565b600454600090637fffffef116105c8576040805162461bcd60e51b815260206004820152601f60248201527f4163636f756e74547265653a207269676874207365742069732066756c6c2000604482015290519081900360640190fd5b6105d0610a6f565b60005b600881101561064d57600181901b8481601081106105ed57fe5b602002015185826001016010811061060157fe5b602002015160405160200180838152602001828152602001925050506040516020818303038152906040528051906020012083836008811061063f57fe5b6020020152506001016105d3565b5060015b60048110156106ec5760016000196004839003011b60005b818110156106e257600181901b84816008811061068257fe5b602002015185826001016008811061069657fe5b60200201516040516020018083815260200182815260200192505050604051602081830303815290604052805190602001208583600881106106d457fe5b602002015250600101610669565b5050600101610651565b508051600454601090046000805b601b8110156107bc57826001166001141561075357604381601b811061071c57fe5b01548460405160200180838152602001828152602001925050506040516020818303038152906040528051906020012093506107b0565b8161076d5783604382601b811061076657fe5b0155600191505b83600560048301601f811061077e57fe5b015460405160200180838152602001828152602001925050506040516020818303038152906040528051906020012093505b600192831c92016106fa565b5050506001819055600054604080516020808201939093528082019390935280518084038201815260609093019052815191012060025550506004805460108101909155919050565b6000637fffffff831684825b601f8110156108ba57826001166001141561086c578481601f811061083257fe5b60200201518260405160200180838152602001828152602001925050506040516020818303038152906040528051906020012091506108ae565b818582601f811061087957fe5b602002015160405160200180838152602001828152602001925050506040516020818303038152906040528051906020012091505b600192831c9201610811565b5063800000008510156108d4576000541491506104869050565b6001541491506104869050565b600354600090637fffffff1161093e576040805162461bcd60e51b815260206004820152601e60248201527f4163636f756e74547265653a206c656674207365742069732066756c6c200000604482015290519081900360640190fd5b60035482906000805b601f811015610a065782600116600114156109a057602481601f811061096957fe5b01548460405160200180838152602001828152602001925050506040516020818303038152906040528051906020012093506109fa565b816109ba5783602482601f81106109b357fe5b0155600191505b83600582601f81106109c857fe5b015460405160200180838152602001828152602001925050506040516020818303038152906040528051906020012093505b600192831c9201610947565b505050600081905560018054604080516020808201959095528082019290925280518083038201815260609092019052805192019190912060025560038054918201905592915050565b6040518061020001604052806010906020820280368337509192915050565b604051806101000160405280600890602082028036833750919291505056fe424c534163636f756e7452656769737472793a20496e76616c69642070726f706f736572a164736f6c634300060c000a",
}

// AccountRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use AccountRegistryMetaData.ABI instead.
var AccountRegistryABI = AccountRegistryMetaData.ABI

// AccountRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AccountRegistryMetaData.Bin instead.
var AccountRegistryBin = AccountRegistryMetaData.Bin

// DeployAccountRegistry deploys a new Ethereum contract, binding an instance of AccountRegistry to it.
func DeployAccountRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, _chooser common.Address, rootLeft [32]byte, leafIndexLeft *big.Int, filledSubtreesLeft [31][32]byte) (common.Address, *types.Transaction, *AccountRegistry, error) {
	parsed, err := AccountRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AccountRegistryBin), backend, _chooser, rootLeft, leafIndexLeft, filledSubtreesLeft)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AccountRegistry{AccountRegistryCaller: AccountRegistryCaller{contract: contract}, AccountRegistryTransactor: AccountRegistryTransactor{contract: contract}, AccountRegistryFilterer: AccountRegistryFilterer{contract: contract}}, nil
}

// AccountRegistry is an auto generated Go binding around an Ethereum contract.
type AccountRegistry struct {
	AccountRegistryCaller     // Read-only binding to the contract
	AccountRegistryTransactor // Write-only binding to the contract
	AccountRegistryFilterer   // Log filterer for contract events
}

// AccountRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type AccountRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AccountRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AccountRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AccountRegistrySession struct {
	Contract     *AccountRegistry  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AccountRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AccountRegistryCallerSession struct {
	Contract *AccountRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// AccountRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AccountRegistryTransactorSession struct {
	Contract     *AccountRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// AccountRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type AccountRegistryRaw struct {
	Contract *AccountRegistry // Generic contract binding to access the raw methods on
}

// AccountRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AccountRegistryCallerRaw struct {
	Contract *AccountRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// AccountRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AccountRegistryTransactorRaw struct {
	Contract *AccountRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAccountRegistry creates a new instance of AccountRegistry, bound to a specific deployed contract.
func NewAccountRegistry(address common.Address, backend bind.ContractBackend) (*AccountRegistry, error) {
	contract, err := bindAccountRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AccountRegistry{AccountRegistryCaller: AccountRegistryCaller{contract: contract}, AccountRegistryTransactor: AccountRegistryTransactor{contract: contract}, AccountRegistryFilterer: AccountRegistryFilterer{contract: contract}}, nil
}

// NewAccountRegistryCaller creates a new read-only instance of AccountRegistry, bound to a specific deployed contract.
func NewAccountRegistryCaller(address common.Address, caller bind.ContractCaller) (*AccountRegistryCaller, error) {
	contract, err := bindAccountRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccountRegistryCaller{contract: contract}, nil
}

// NewAccountRegistryTransactor creates a new write-only instance of AccountRegistry, bound to a specific deployed contract.
func NewAccountRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*AccountRegistryTransactor, error) {
	contract, err := bindAccountRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccountRegistryTransactor{contract: contract}, nil
}

// NewAccountRegistryFilterer creates a new log filterer instance of AccountRegistry, bound to a specific deployed contract.
func NewAccountRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*AccountRegistryFilterer, error) {
	contract, err := bindAccountRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccountRegistryFilterer{contract: contract}, nil
}

// bindAccountRegistry binds a generic wrapper to an already deployed contract.
func bindAccountRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AccountRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccountRegistry *AccountRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccountRegistry.Contract.AccountRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccountRegistry *AccountRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccountRegistry.Contract.AccountRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccountRegistry *AccountRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccountRegistry.Contract.AccountRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccountRegistry *AccountRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccountRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccountRegistry *AccountRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccountRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccountRegistry *AccountRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccountRegistry.Contract.contract.Transact(opts, method, params...)
}

// BATCHDEPTH is a free data retrieval call binding the contract method 0x1c76e77e.
//
// Solidity: function BATCH_DEPTH() view returns(uint256)
func (_AccountRegistry *AccountRegistryCaller) BATCHDEPTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "BATCH_DEPTH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BATCHDEPTH is a free data retrieval call binding the contract method 0x1c76e77e.
//
// Solidity: function BATCH_DEPTH() view returns(uint256)
func (_AccountRegistry *AccountRegistrySession) BATCHDEPTH() (*big.Int, error) {
	return _AccountRegistry.Contract.BATCHDEPTH(&_AccountRegistry.CallOpts)
}

// BATCHDEPTH is a free data retrieval call binding the contract method 0x1c76e77e.
//
// Solidity: function BATCH_DEPTH() view returns(uint256)
func (_AccountRegistry *AccountRegistryCallerSession) BATCHDEPTH() (*big.Int, error) {
	return _AccountRegistry.Contract.BATCHDEPTH(&_AccountRegistry.CallOpts)
}

// BATCHSIZE is a free data retrieval call binding the contract method 0x49faa4d4.
//
// Solidity: function BATCH_SIZE() view returns(uint256)
func (_AccountRegistry *AccountRegistryCaller) BATCHSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "BATCH_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BATCHSIZE is a free data retrieval call binding the contract method 0x49faa4d4.
//
// Solidity: function BATCH_SIZE() view returns(uint256)
func (_AccountRegistry *AccountRegistrySession) BATCHSIZE() (*big.Int, error) {
	return _AccountRegistry.Contract.BATCHSIZE(&_AccountRegistry.CallOpts)
}

// BATCHSIZE is a free data retrieval call binding the contract method 0x49faa4d4.
//
// Solidity: function BATCH_SIZE() view returns(uint256)
func (_AccountRegistry *AccountRegistryCallerSession) BATCHSIZE() (*big.Int, error) {
	return _AccountRegistry.Contract.BATCHSIZE(&_AccountRegistry.CallOpts)
}

// DEPTH is a free data retrieval call binding the contract method 0x98366e35.
//
// Solidity: function DEPTH() view returns(uint256)
func (_AccountRegistry *AccountRegistryCaller) DEPTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "DEPTH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DEPTH is a free data retrieval call binding the contract method 0x98366e35.
//
// Solidity: function DEPTH() view returns(uint256)
func (_AccountRegistry *AccountRegistrySession) DEPTH() (*big.Int, error) {
	return _AccountRegistry.Contract.DEPTH(&_AccountRegistry.CallOpts)
}

// DEPTH is a free data retrieval call binding the contract method 0x98366e35.
//
// Solidity: function DEPTH() view returns(uint256)
func (_AccountRegistry *AccountRegistryCallerSession) DEPTH() (*big.Int, error) {
	return _AccountRegistry.Contract.DEPTH(&_AccountRegistry.CallOpts)
}

// SETSIZE is a free data retrieval call binding the contract method 0xd0383d68.
//
// Solidity: function SET_SIZE() view returns(uint256)
func (_AccountRegistry *AccountRegistryCaller) SETSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "SET_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SETSIZE is a free data retrieval call binding the contract method 0xd0383d68.
//
// Solidity: function SET_SIZE() view returns(uint256)
func (_AccountRegistry *AccountRegistrySession) SETSIZE() (*big.Int, error) {
	return _AccountRegistry.Contract.SETSIZE(&_AccountRegistry.CallOpts)
}

// SETSIZE is a free data retrieval call binding the contract method 0xd0383d68.
//
// Solidity: function SET_SIZE() view returns(uint256)
func (_AccountRegistry *AccountRegistryCallerSession) SETSIZE() (*big.Int, error) {
	return _AccountRegistry.Contract.SETSIZE(&_AccountRegistry.CallOpts)
}

// WITNESSLENGTH is a free data retrieval call binding the contract method 0x5e71468b.
//
// Solidity: function WITNESS_LENGTH() view returns(uint256)
func (_AccountRegistry *AccountRegistryCaller) WITNESSLENGTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "WITNESS_LENGTH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WITNESSLENGTH is a free data retrieval call binding the contract method 0x5e71468b.
//
// Solidity: function WITNESS_LENGTH() view returns(uint256)
func (_AccountRegistry *AccountRegistrySession) WITNESSLENGTH() (*big.Int, error) {
	return _AccountRegistry.Contract.WITNESSLENGTH(&_AccountRegistry.CallOpts)
}

// WITNESSLENGTH is a free data retrieval call binding the contract method 0x5e71468b.
//
// Solidity: function WITNESS_LENGTH() view returns(uint256)
func (_AccountRegistry *AccountRegistryCallerSession) WITNESSLENGTH() (*big.Int, error) {
	return _AccountRegistry.Contract.WITNESSLENGTH(&_AccountRegistry.CallOpts)
}

// Chooser is a free data retrieval call binding the contract method 0x98d17621.
//
// Solidity: function chooser() view returns(address)
func (_AccountRegistry *AccountRegistryCaller) Chooser(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "chooser")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Chooser is a free data retrieval call binding the contract method 0x98d17621.
//
// Solidity: function chooser() view returns(address)
func (_AccountRegistry *AccountRegistrySession) Chooser() (common.Address, error) {
	return _AccountRegistry.Contract.Chooser(&_AccountRegistry.CallOpts)
}

// Chooser is a free data retrieval call binding the contract method 0x98d17621.
//
// Solidity: function chooser() view returns(address)
func (_AccountRegistry *AccountRegistryCallerSession) Chooser() (common.Address, error) {
	return _AccountRegistry.Contract.Chooser(&_AccountRegistry.CallOpts)
}

// Exists is a free data retrieval call binding the contract method 0x709a8b2a.
//
// Solidity: function exists(uint256 pubkeyID, uint256[4] pubkey, bytes32[31] witness) view returns(bool)
func (_AccountRegistry *AccountRegistryCaller) Exists(opts *bind.CallOpts, pubkeyID *big.Int, pubkey [4]*big.Int, witness [31][32]byte) (bool, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "exists", pubkeyID, pubkey, witness)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Exists is a free data retrieval call binding the contract method 0x709a8b2a.
//
// Solidity: function exists(uint256 pubkeyID, uint256[4] pubkey, bytes32[31] witness) view returns(bool)
func (_AccountRegistry *AccountRegistrySession) Exists(pubkeyID *big.Int, pubkey [4]*big.Int, witness [31][32]byte) (bool, error) {
	return _AccountRegistry.Contract.Exists(&_AccountRegistry.CallOpts, pubkeyID, pubkey, witness)
}

// Exists is a free data retrieval call binding the contract method 0x709a8b2a.
//
// Solidity: function exists(uint256 pubkeyID, uint256[4] pubkey, bytes32[31] witness) view returns(bool)
func (_AccountRegistry *AccountRegistryCallerSession) Exists(pubkeyID *big.Int, pubkey [4]*big.Int, witness [31][32]byte) (bool, error) {
	return _AccountRegistry.Contract.Exists(&_AccountRegistry.CallOpts, pubkeyID, pubkey, witness)
}

// FilledSubtreesLeft is a free data retrieval call binding the contract method 0x034a29ae.
//
// Solidity: function filledSubtreesLeft(uint256 ) view returns(bytes32)
func (_AccountRegistry *AccountRegistryCaller) FilledSubtreesLeft(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "filledSubtreesLeft", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FilledSubtreesLeft is a free data retrieval call binding the contract method 0x034a29ae.
//
// Solidity: function filledSubtreesLeft(uint256 ) view returns(bytes32)
func (_AccountRegistry *AccountRegistrySession) FilledSubtreesLeft(arg0 *big.Int) ([32]byte, error) {
	return _AccountRegistry.Contract.FilledSubtreesLeft(&_AccountRegistry.CallOpts, arg0)
}

// FilledSubtreesLeft is a free data retrieval call binding the contract method 0x034a29ae.
//
// Solidity: function filledSubtreesLeft(uint256 ) view returns(bytes32)
func (_AccountRegistry *AccountRegistryCallerSession) FilledSubtreesLeft(arg0 *big.Int) ([32]byte, error) {
	return _AccountRegistry.Contract.FilledSubtreesLeft(&_AccountRegistry.CallOpts, arg0)
}

// FilledSubtreesRight is a free data retrieval call binding the contract method 0x8d037962.
//
// Solidity: function filledSubtreesRight(uint256 ) view returns(bytes32)
func (_AccountRegistry *AccountRegistryCaller) FilledSubtreesRight(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "filledSubtreesRight", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FilledSubtreesRight is a free data retrieval call binding the contract method 0x8d037962.
//
// Solidity: function filledSubtreesRight(uint256 ) view returns(bytes32)
func (_AccountRegistry *AccountRegistrySession) FilledSubtreesRight(arg0 *big.Int) ([32]byte, error) {
	return _AccountRegistry.Contract.FilledSubtreesRight(&_AccountRegistry.CallOpts, arg0)
}

// FilledSubtreesRight is a free data retrieval call binding the contract method 0x8d037962.
//
// Solidity: function filledSubtreesRight(uint256 ) view returns(bytes32)
func (_AccountRegistry *AccountRegistryCallerSession) FilledSubtreesRight(arg0 *big.Int) ([32]byte, error) {
	return _AccountRegistry.Contract.FilledSubtreesRight(&_AccountRegistry.CallOpts, arg0)
}

// LeafIndexLeft is a free data retrieval call binding the contract method 0x693c1db7.
//
// Solidity: function leafIndexLeft() view returns(uint256)
func (_AccountRegistry *AccountRegistryCaller) LeafIndexLeft(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "leafIndexLeft")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LeafIndexLeft is a free data retrieval call binding the contract method 0x693c1db7.
//
// Solidity: function leafIndexLeft() view returns(uint256)
func (_AccountRegistry *AccountRegistrySession) LeafIndexLeft() (*big.Int, error) {
	return _AccountRegistry.Contract.LeafIndexLeft(&_AccountRegistry.CallOpts)
}

// LeafIndexLeft is a free data retrieval call binding the contract method 0x693c1db7.
//
// Solidity: function leafIndexLeft() view returns(uint256)
func (_AccountRegistry *AccountRegistryCallerSession) LeafIndexLeft() (*big.Int, error) {
	return _AccountRegistry.Contract.LeafIndexLeft(&_AccountRegistry.CallOpts)
}

// LeafIndexRight is a free data retrieval call binding the contract method 0xd7c53ea7.
//
// Solidity: function leafIndexRight() view returns(uint256)
func (_AccountRegistry *AccountRegistryCaller) LeafIndexRight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "leafIndexRight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LeafIndexRight is a free data retrieval call binding the contract method 0xd7c53ea7.
//
// Solidity: function leafIndexRight() view returns(uint256)
func (_AccountRegistry *AccountRegistrySession) LeafIndexRight() (*big.Int, error) {
	return _AccountRegistry.Contract.LeafIndexRight(&_AccountRegistry.CallOpts)
}

// LeafIndexRight is a free data retrieval call binding the contract method 0xd7c53ea7.
//
// Solidity: function leafIndexRight() view returns(uint256)
func (_AccountRegistry *AccountRegistryCallerSession) LeafIndexRight() (*big.Int, error) {
	return _AccountRegistry.Contract.LeafIndexRight(&_AccountRegistry.CallOpts)
}

// Root is a free data retrieval call binding the contract method 0xebf0c717.
//
// Solidity: function root() view returns(bytes32)
func (_AccountRegistry *AccountRegistryCaller) Root(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "root")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Root is a free data retrieval call binding the contract method 0xebf0c717.
//
// Solidity: function root() view returns(bytes32)
func (_AccountRegistry *AccountRegistrySession) Root() ([32]byte, error) {
	return _AccountRegistry.Contract.Root(&_AccountRegistry.CallOpts)
}

// Root is a free data retrieval call binding the contract method 0xebf0c717.
//
// Solidity: function root() view returns(bytes32)
func (_AccountRegistry *AccountRegistryCallerSession) Root() ([32]byte, error) {
	return _AccountRegistry.Contract.Root(&_AccountRegistry.CallOpts)
}

// RootLeft is a free data retrieval call binding the contract method 0xd8289463.
//
// Solidity: function rootLeft() view returns(bytes32)
func (_AccountRegistry *AccountRegistryCaller) RootLeft(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "rootLeft")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RootLeft is a free data retrieval call binding the contract method 0xd8289463.
//
// Solidity: function rootLeft() view returns(bytes32)
func (_AccountRegistry *AccountRegistrySession) RootLeft() ([32]byte, error) {
	return _AccountRegistry.Contract.RootLeft(&_AccountRegistry.CallOpts)
}

// RootLeft is a free data retrieval call binding the contract method 0xd8289463.
//
// Solidity: function rootLeft() view returns(bytes32)
func (_AccountRegistry *AccountRegistryCallerSession) RootLeft() ([32]byte, error) {
	return _AccountRegistry.Contract.RootLeft(&_AccountRegistry.CallOpts)
}

// RootRight is a free data retrieval call binding the contract method 0xcab2da9b.
//
// Solidity: function rootRight() view returns(bytes32)
func (_AccountRegistry *AccountRegistryCaller) RootRight(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "rootRight")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RootRight is a free data retrieval call binding the contract method 0xcab2da9b.
//
// Solidity: function rootRight() view returns(bytes32)
func (_AccountRegistry *AccountRegistrySession) RootRight() ([32]byte, error) {
	return _AccountRegistry.Contract.RootRight(&_AccountRegistry.CallOpts)
}

// RootRight is a free data retrieval call binding the contract method 0xcab2da9b.
//
// Solidity: function rootRight() view returns(bytes32)
func (_AccountRegistry *AccountRegistryCallerSession) RootRight() ([32]byte, error) {
	return _AccountRegistry.Contract.RootRight(&_AccountRegistry.CallOpts)
}

// Zeros is a free data retrieval call binding the contract method 0xe8295588.
//
// Solidity: function zeros(uint256 ) view returns(bytes32)
func (_AccountRegistry *AccountRegistryCaller) Zeros(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "zeros", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Zeros is a free data retrieval call binding the contract method 0xe8295588.
//
// Solidity: function zeros(uint256 ) view returns(bytes32)
func (_AccountRegistry *AccountRegistrySession) Zeros(arg0 *big.Int) ([32]byte, error) {
	return _AccountRegistry.Contract.Zeros(&_AccountRegistry.CallOpts, arg0)
}

// Zeros is a free data retrieval call binding the contract method 0xe8295588.
//
// Solidity: function zeros(uint256 ) view returns(bytes32)
func (_AccountRegistry *AccountRegistryCallerSession) Zeros(arg0 *big.Int) ([32]byte, error) {
	return _AccountRegistry.Contract.Zeros(&_AccountRegistry.CallOpts, arg0)
}

// Register is a paid mutator transaction binding the contract method 0x95e4bf03.
//
// Solidity: function register(uint256[4] pubkey) returns(uint256)
func (_AccountRegistry *AccountRegistryTransactor) Register(opts *bind.TransactOpts, pubkey [4]*big.Int) (*types.Transaction, error) {
	return _AccountRegistry.contract.Transact(opts, "register", pubkey)
}

// Register is a paid mutator transaction binding the contract method 0x95e4bf03.
//
// Solidity: function register(uint256[4] pubkey) returns(uint256)
func (_AccountRegistry *AccountRegistrySession) Register(pubkey [4]*big.Int) (*types.Transaction, error) {
	return _AccountRegistry.Contract.Register(&_AccountRegistry.TransactOpts, pubkey)
}

// Register is a paid mutator transaction binding the contract method 0x95e4bf03.
//
// Solidity: function register(uint256[4] pubkey) returns(uint256)
func (_AccountRegistry *AccountRegistryTransactorSession) Register(pubkey [4]*big.Int) (*types.Transaction, error) {
	return _AccountRegistry.Contract.Register(&_AccountRegistry.TransactOpts, pubkey)
}

// RegisterBatch is a paid mutator transaction binding the contract method 0x1c4a7a94.
//
// Solidity: function registerBatch(uint256[4][16] pubkeys) returns(uint256)
func (_AccountRegistry *AccountRegistryTransactor) RegisterBatch(opts *bind.TransactOpts, pubkeys [16][4]*big.Int) (*types.Transaction, error) {
	return _AccountRegistry.contract.Transact(opts, "registerBatch", pubkeys)
}

// RegisterBatch is a paid mutator transaction binding the contract method 0x1c4a7a94.
//
// Solidity: function registerBatch(uint256[4][16] pubkeys) returns(uint256)
func (_AccountRegistry *AccountRegistrySession) RegisterBatch(pubkeys [16][4]*big.Int) (*types.Transaction, error) {
	return _AccountRegistry.Contract.RegisterBatch(&_AccountRegistry.TransactOpts, pubkeys)
}

// RegisterBatch is a paid mutator transaction binding the contract method 0x1c4a7a94.
//
// Solidity: function registerBatch(uint256[4][16] pubkeys) returns(uint256)
func (_AccountRegistry *AccountRegistryTransactorSession) RegisterBatch(pubkeys [16][4]*big.Int) (*types.Transaction, error) {
	return _AccountRegistry.Contract.RegisterBatch(&_AccountRegistry.TransactOpts, pubkeys)
}

// AccountRegistryBatchPubkeyRegisteredIterator is returned from FilterBatchPubkeyRegistered and is used to iterate over the raw logs and unpacked data for BatchPubkeyRegistered events raised by the AccountRegistry contract.
type AccountRegistryBatchPubkeyRegisteredIterator struct {
	Event *AccountRegistryBatchPubkeyRegistered // Event containing the contract specifics and raw log

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
func (it *AccountRegistryBatchPubkeyRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccountRegistryBatchPubkeyRegistered)
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
		it.Event = new(AccountRegistryBatchPubkeyRegistered)
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
func (it *AccountRegistryBatchPubkeyRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccountRegistryBatchPubkeyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccountRegistryBatchPubkeyRegistered represents a BatchPubkeyRegistered event raised by the AccountRegistry contract.
type AccountRegistryBatchPubkeyRegistered struct {
	StartID *big.Int
	EndID   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBatchPubkeyRegistered is a free log retrieval operation binding the contract event 0x3154b80a7d9f6a143c37dde575f47deb78dacc7f280d8efc7e3ae102758a841b.
//
// Solidity: event BatchPubkeyRegistered(uint256 startID, uint256 endID)
func (_AccountRegistry *AccountRegistryFilterer) FilterBatchPubkeyRegistered(opts *bind.FilterOpts) (*AccountRegistryBatchPubkeyRegisteredIterator, error) {

	logs, sub, err := _AccountRegistry.contract.FilterLogs(opts, "BatchPubkeyRegistered")
	if err != nil {
		return nil, err
	}
	return &AccountRegistryBatchPubkeyRegisteredIterator{contract: _AccountRegistry.contract, event: "BatchPubkeyRegistered", logs: logs, sub: sub}, nil
}

// WatchBatchPubkeyRegistered is a free log subscription operation binding the contract event 0x3154b80a7d9f6a143c37dde575f47deb78dacc7f280d8efc7e3ae102758a841b.
//
// Solidity: event BatchPubkeyRegistered(uint256 startID, uint256 endID)
func (_AccountRegistry *AccountRegistryFilterer) WatchBatchPubkeyRegistered(opts *bind.WatchOpts, sink chan<- *AccountRegistryBatchPubkeyRegistered) (event.Subscription, error) {

	logs, sub, err := _AccountRegistry.contract.WatchLogs(opts, "BatchPubkeyRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccountRegistryBatchPubkeyRegistered)
				if err := _AccountRegistry.contract.UnpackLog(event, "BatchPubkeyRegistered", log); err != nil {
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

// ParseBatchPubkeyRegistered is a log parse operation binding the contract event 0x3154b80a7d9f6a143c37dde575f47deb78dacc7f280d8efc7e3ae102758a841b.
//
// Solidity: event BatchPubkeyRegistered(uint256 startID, uint256 endID)
func (_AccountRegistry *AccountRegistryFilterer) ParseBatchPubkeyRegistered(log types.Log) (*AccountRegistryBatchPubkeyRegistered, error) {
	event := new(AccountRegistryBatchPubkeyRegistered)
	if err := _AccountRegistry.contract.UnpackLog(event, "BatchPubkeyRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccountRegistrySinglePubkeyRegisteredIterator is returned from FilterSinglePubkeyRegistered and is used to iterate over the raw logs and unpacked data for SinglePubkeyRegistered events raised by the AccountRegistry contract.
type AccountRegistrySinglePubkeyRegisteredIterator struct {
	Event *AccountRegistrySinglePubkeyRegistered // Event containing the contract specifics and raw log

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
func (it *AccountRegistrySinglePubkeyRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccountRegistrySinglePubkeyRegistered)
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
		it.Event = new(AccountRegistrySinglePubkeyRegistered)
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
func (it *AccountRegistrySinglePubkeyRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccountRegistrySinglePubkeyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccountRegistrySinglePubkeyRegistered represents a SinglePubkeyRegistered event raised by the AccountRegistry contract.
type AccountRegistrySinglePubkeyRegistered struct {
	PubkeyID *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSinglePubkeyRegistered is a free log retrieval operation binding the contract event 0x59056afed767866d7c194cf26e24ebe16974ef943cccb729452e9adc265a9ac7.
//
// Solidity: event SinglePubkeyRegistered(uint256 pubkeyID)
func (_AccountRegistry *AccountRegistryFilterer) FilterSinglePubkeyRegistered(opts *bind.FilterOpts) (*AccountRegistrySinglePubkeyRegisteredIterator, error) {

	logs, sub, err := _AccountRegistry.contract.FilterLogs(opts, "SinglePubkeyRegistered")
	if err != nil {
		return nil, err
	}
	return &AccountRegistrySinglePubkeyRegisteredIterator{contract: _AccountRegistry.contract, event: "SinglePubkeyRegistered", logs: logs, sub: sub}, nil
}

// WatchSinglePubkeyRegistered is a free log subscription operation binding the contract event 0x59056afed767866d7c194cf26e24ebe16974ef943cccb729452e9adc265a9ac7.
//
// Solidity: event SinglePubkeyRegistered(uint256 pubkeyID)
func (_AccountRegistry *AccountRegistryFilterer) WatchSinglePubkeyRegistered(opts *bind.WatchOpts, sink chan<- *AccountRegistrySinglePubkeyRegistered) (event.Subscription, error) {

	logs, sub, err := _AccountRegistry.contract.WatchLogs(opts, "SinglePubkeyRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccountRegistrySinglePubkeyRegistered)
				if err := _AccountRegistry.contract.UnpackLog(event, "SinglePubkeyRegistered", log); err != nil {
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

// ParseSinglePubkeyRegistered is a log parse operation binding the contract event 0x59056afed767866d7c194cf26e24ebe16974ef943cccb729452e9adc265a9ac7.
//
// Solidity: event SinglePubkeyRegistered(uint256 pubkeyID)
func (_AccountRegistry *AccountRegistryFilterer) ParseSinglePubkeyRegistered(log types.Log) (*AccountRegistrySinglePubkeyRegistered, error) {
	event := new(AccountRegistrySinglePubkeyRegistered)
	if err := _AccountRegistry.contract.UnpackLog(event, "SinglePubkeyRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
