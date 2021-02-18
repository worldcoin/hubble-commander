// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package accountregistry

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

// AccountRegistryABI is the input ABI used to generate the binding from.
const AccountRegistryABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256[4]\",\"name\":\"pubkey\",\"type\":\"uint256[4]\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"}],\"name\":\"PubkeyRegistered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BATCH_DEPTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"BATCH_SIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEPTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SET_SIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WITNESS_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256[4]\",\"name\":\"pubkey\",\"type\":\"uint256[4]\"},{\"internalType\":\"bytes32[31]\",\"name\":\"witness\",\"type\":\"bytes32[31]\"}],\"name\":\"exists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"filledSubtreesLeft\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"filledSubtreesRight\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"leafIndexLeft\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"leafIndexRight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[4]\",\"name\":\"pubkey\",\"type\":\"uint256[4]\"}],\"name\":\"register\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[4][16]\",\"name\":\"pubkeys\",\"type\":\"uint256[4][16]\"}],\"name\":\"registerBatch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"root\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rootLeft\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rootRight\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"zeros\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// AccountRegistryBin is the compiled bytecode used for deploying new contracts.
var AccountRegistryBin = "0x60806040526000600355600060045534801561001a57600080fd5b507f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e5636005819055602481905560015b601f81101561011a57600560018203601f811061006257fe5b0154600560018303601f811061007457fe5b0154604051602001808381526020018281526020019250505060405160208183030381529060405280519060200120600582601f81106100b057fe5b0155601f8110156100da57600581601f81106100c857fe5b0154602482601f81106100d757fe5b01555b806004111580156100eb575080601f115b1561011257600581601f81106100fd57fe5b0154604360048303601b811061010f57fe5b01555b600101610049565b505060235460408051602080820184905281830184905282518083038401815260608301845280519082012060008181556080840186905260a0808501969096528451808503909601865260c084018552855195830195909520600181905560e0840191909152610100808401919091528351808403909101815261012090920190925280519101206002556109c99081906101b690396000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c806395e4bf0311610097578063d7c53ea711610066578063d7c53ea7146101f6578063d8289463146101fe578063e829558814610206578063ebf0c7171461022357610100565b806395e4bf03146101ca57806398366e3514610161578063cab2da9b146101e6578063d0383d68146101ee57610100565b80635e71468b116100d35780635e71468b14610161578063693c1db714610169578063709a8b2a146101715780638d037962146101ad57610100565b8063034a29ae146101055780631c4a7a94146101345780631c76e77e1461015157806349faa4d414610159575b600080fd5b6101226004803603602081101561011b57600080fd5b503561022b565b60408051918252519081900360200190f35b610122600480360361080081101561014b57600080fd5b5061023f565b610122610336565b61012261033b565b610122610340565b610122610345565b610199600480360361048081101561018857600080fd5b50803590602081019060a00161034b565b604080519115158252519081900360200190f35b610122600480360360208110156101c357600080fd5b50356103c1565b610122600480360360808110156101e057600080fd5b506103ce565b61012261046c565b610122610472565b61012261047a565b610122610480565b6101226004803603602081101561021c57600080fd5b5035610486565b610122610493565b602481601f811061023857fe5b0154905081565b600061024961097e565b60005b6010811015610322577ff0777e5cea47492e18df87dcc844efabdfad315d1a2b4883d87cc2b964eddff084826010811061028257fe5b6080020182601f6001901b60045401016040518083600460200280828437600083820152601f01601f19169091019283525050604051908190036020019150a160008482601081106102d057fe5b608002016040516020018082600460200280828437808301925050509150506040516020818303038152906040528051906020012090508083836010811061031457fe5b60200201525060010161024c565b50600061032e82610499565b949350505050565b600481565b601081565b601f81565b60035481565b600080836040516020018082600460200280828437808301925050509150506040516020818303038152906040528051906020012090506103b6818685601f806020026040519081016040528092919082601f60200280828437600092019190915250610733915050565b9150505b9392505050565b604381601b811061023857fe5b6000808260405160200180826004602002808284378083019250505091505060405160208183030381529060405280519060200120905060006104108261080f565b90507ff0777e5cea47492e18df87dcc844efabdfad315d1a2b4883d87cc2b964eddff084826040518083600460200280828437600083820152601f01601f19169091019283525050604051908190036020019150a19392505050565b60015481565b638000000081565b60045481565b60005481565b600581601f811061023857fe5b60025481565b600454600090637fffffef116104f6576040805162461bcd60e51b815260206004820152601f60248201527f4163636f756e74547265653a207269676874207365742069732066756c6c2000604482015290519081900360640190fd5b6104fe61099d565b60005b600881101561057b57600181901b84816010811061051b57fe5b602002015185826001016010811061052f57fe5b602002015160405160200180838152602001828152602001925050506040516020818303038152906040528051906020012083836008811061056d57fe5b602002015250600101610501565b5060015b600481101561061a5760016000196004839003011b60005b8181101561061057600181901b8481600881106105b057fe5b60200201518582600101600881106105c457fe5b602002015160405160200180838152602001828152602001925050506040516020818303038152906040528051906020012085836008811061060257fe5b602002015250600101610597565b505060010161057f565b508051600454601090046000805b601b8110156106ea57826001166001141561068157604381601b811061064a57fe5b01548460405160200180838152602001828152602001925050506040516020818303038152906040528051906020012093506106de565b8161069b5783604382601b811061069457fe5b0155600191505b83600560048301601f81106106ac57fe5b015460405160200180838152602001828152602001925050506040516020818303038152906040528051906020012093505b600192831c9201610628565b5050506001819055600054604080516020808201939093528082019390935280518084038201815260609093019052815191012060025550506004805460108101909155919050565b6000637fffffff831684825b601f8110156107e857826001166001141561079a578481601f811061076057fe5b60200201518260405160200180838152602001828152602001925050506040516020818303038152906040528051906020012091506107dc565b818582601f81106107a757fe5b602002015160405160200180838152602001828152602001925050506040516020818303038152906040528051906020012091505b600192831c920161073f565b506380000000851015610802576000541491506103ba9050565b6001541491506103ba9050565b600354600090637fffffff1161086c576040805162461bcd60e51b815260206004820152601e60248201527f4163636f756e74547265653a206c656674207365742069732066756c6c200000604482015290519081900360640190fd5b60035482906000805b601f8110156109345782600116600114156108ce57602481601f811061089757fe5b0154846040516020018083815260200182815260200192505050604051602081830303815290604052805190602001209350610928565b816108e85783602482601f81106108e157fe5b0155600191505b83600582601f81106108f657fe5b015460405160200180838152602001828152602001925050506040516020818303038152906040528051906020012093505b600192831c9201610875565b505050600081905560018054604080516020808201959095528082019290925280518083038201815260609092019052805192019190912060025560038054918201905592915050565b6040518061020001604052806010906020820280368337509192915050565b604051806101000160405280600890602082028036833750919291505056fea164736f6c634300060c000a"

// DeployAccountRegistry deploys a new Ethereum contract, binding an instance of AccountRegistry to it.
func DeployAccountRegistry(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *AccountRegistry, error) {
	parsed, err := abi.JSON(strings.NewReader(AccountRegistryABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AccountRegistryBin), backend)
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

// AccountRegistryPubkeyRegisteredIterator is returned from FilterPubkeyRegistered and is used to iterate over the raw logs and unpacked data for PubkeyRegistered events raised by the AccountRegistry contract.
type AccountRegistryPubkeyRegisteredIterator struct {
	Event *AccountRegistryPubkeyRegistered // Event containing the contract specifics and raw log

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
func (it *AccountRegistryPubkeyRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccountRegistryPubkeyRegistered)
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
		it.Event = new(AccountRegistryPubkeyRegistered)
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
func (it *AccountRegistryPubkeyRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccountRegistryPubkeyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccountRegistryPubkeyRegistered represents a PubkeyRegistered event raised by the AccountRegistry contract.
type AccountRegistryPubkeyRegistered struct {
	Pubkey   [4]*big.Int
	PubkeyID *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterPubkeyRegistered is a free log retrieval operation binding the contract event 0xf0777e5cea47492e18df87dcc844efabdfad315d1a2b4883d87cc2b964eddff0.
//
// Solidity: event PubkeyRegistered(uint256[4] pubkey, uint256 pubkeyID)
func (_AccountRegistry *AccountRegistryFilterer) FilterPubkeyRegistered(opts *bind.FilterOpts) (*AccountRegistryPubkeyRegisteredIterator, error) {

	logs, sub, err := _AccountRegistry.contract.FilterLogs(opts, "PubkeyRegistered")
	if err != nil {
		return nil, err
	}
	return &AccountRegistryPubkeyRegisteredIterator{contract: _AccountRegistry.contract, event: "PubkeyRegistered", logs: logs, sub: sub}, nil
}

// WatchPubkeyRegistered is a free log subscription operation binding the contract event 0xf0777e5cea47492e18df87dcc844efabdfad315d1a2b4883d87cc2b964eddff0.
//
// Solidity: event PubkeyRegistered(uint256[4] pubkey, uint256 pubkeyID)
func (_AccountRegistry *AccountRegistryFilterer) WatchPubkeyRegistered(opts *bind.WatchOpts, sink chan<- *AccountRegistryPubkeyRegistered) (event.Subscription, error) {

	logs, sub, err := _AccountRegistry.contract.WatchLogs(opts, "PubkeyRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccountRegistryPubkeyRegistered)
				if err := _AccountRegistry.contract.UnpackLog(event, "PubkeyRegistered", log); err != nil {
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

// ParsePubkeyRegistered is a log parse operation binding the contract event 0xf0777e5cea47492e18df87dcc844efabdfad315d1a2b4883d87cc2b964eddff0.
//
// Solidity: event PubkeyRegistered(uint256[4] pubkey, uint256 pubkeyID)
func (_AccountRegistry *AccountRegistryFilterer) ParsePubkeyRegistered(log types.Log) (*AccountRegistryPubkeyRegistered, error) {
	event := new(AccountRegistryPubkeyRegistered)
	if err := _AccountRegistry.contract.UnpackLog(event, "PubkeyRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
