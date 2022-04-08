// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vault

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

// TypesMMCommitmentInclusionProof is an auto generated low-level Go binding around an user-defined struct.
type TypesMMCommitmentInclusionProof struct {
	Commitment TypesMassMigrationCommitment
	Path       *big.Int
	Witness    [][32]byte
}

// TypesMassMigrationBody is an auto generated low-level Go binding around an user-defined struct.
type TypesMassMigrationBody struct {
	AccountRoot  [32]byte
	Signature    [2]*big.Int
	SpokeID      *big.Int
	WithdrawRoot [32]byte
	TokenID      *big.Int
	Amount       *big.Int
	FeeReceiver  *big.Int
	Txs          []byte
}

// TypesMassMigrationCommitment is an auto generated low-level Go binding around an user-defined struct.
type TypesMassMigrationCommitment struct {
	StateRoot [32]byte
	Body      TypesMassMigrationBody
}

// VaultMetaData contains all meta data concerning the Vault contract.
var VaultMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractITokenRegistry\",\"name\":\"_tokenRegistry\",\"type\":\"address\"},{\"internalType\":\"contractSpokeRegistry\",\"name\":\"_spokes\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"}],\"name\":\"isBatchApproved\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"spokeID\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"withdrawRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.MassMigrationBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"internalType\":\"structTypes.MassMigrationCommitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.MMCommitmentInclusionProof\",\"name\":\"commitmentMP\",\"type\":\"tuple\"}],\"name\":\"requestApproval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollup\",\"outputs\":[{\"internalType\":\"contractRollup\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractRollup\",\"name\":\"_rollup\",\"type\":\"address\"}],\"name\":\"setRollupAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"spokes\",\"outputs\":[{\"internalType\":\"contractSpokeRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenRegistry\",\"outputs\":[{\"internalType\":\"contractITokenRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60e060405234801561001057600080fd5b50604051610f6a380380610f6a83398101604081905261002f91610053565b33606090811b6080526001600160601b031992811b831660c0521b1660a0526100a4565b60008060408385031215610065578182fd5b82516100708161008c565b60208401519092506100818161008c565b809150509250929050565b6001600160a01b03811681146100a157600080fd5b50565b60805160601c60a05160601c60c05160601c610e7a6100f0600039600081816101110152610590015260008181610138015261031301526000818160d401526102270152610e7a6000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c80638da5cb5b1161005b5780638da5cb5b146100d25780639d23c4c71461010c578063b640e50414610133578063cb23bcb51461015a57600080fd5b806307663706146100825780630d82dcf4146100975780634abb54d9146100bf575b600080fd5b610095610090366004610ac5565b610173565b005b6100aa6100a5366004610b12565b6102de565b60405190151581526020015b60405180910390f35b6100956100cd366004610b2a565b6102f1565b7f00000000000000000000000000000000000000000000000000000000000000005b6040516001600160a01b0390911681526020016100b6565b6100f47f000000000000000000000000000000000000000000000000000000000000000081565b6100f47f000000000000000000000000000000000000000000000000000000000000000081565b6000546100f4906201000090046001600160a01b031681565b600054610100900460ff1661018e5760005460ff1615610192565b303b155b6101fa5760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084015b60405180910390fd5b600054610100900460ff1615801561021c576000805461ffff19166101011790555b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146102a65760405162461bcd60e51b815260206004820152602960248201527f496d6d757461626c654f776e61626c653a2063616c6c6572206973206e6f74206044820152683a34329037bbb732b960b91b60648201526084016101f1565b6000805462010000600160b01b031916620100006001600160a01b0385160217905580156102da576000805461ff00191690555b5050565b60006102eb82600161070f565b92915050565b805160200151604090810151905163da2fd3df60e01b815260048101919091527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063da2fd3df9060240160206040518083038186803b15801561035d57600080fd5b505afa158015610371573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103959190610a55565b6001600160a01b0316336001600160a01b03161461040a5760405162461bcd60e51b815260206004820152602c60248201527f5661756c743a206d73672e73656e6465722073686f756c64206265207468652060448201526b7461726765742073706f6b6560a01b60648201526084016101f1565b60008054604051632d62214160e11b815260048101859052620100009091046001600160a01b031690635ac4428290602401604080518083038186803b15801561045357600080fd5b505afa158015610467573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061048b9190610ae1565b602081015190915060301c63ffffffff164310156104eb5760405162461bcd60e51b815260206004820181905260248201527f5661756c743a2042617463682073686f756c642062652066696e616c6973656460448201526064016101f1565b8051825161050c91906104fd90610750565b84602001518560400151610794565b61056a5760405162461bcd60e51b815260206004820152602960248201527f5661756c743a20436f6d6d69746d656e74206973206e6f742070726573656e74604482015268040d2dc40c4c2e8c6d60bb1b60648201526084016101f1565b81516020015160800151604051630a7973b760e01b8152600481019190915260009081907f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031690630a7973b790602401604080518083038186803b1580156105d957600080fd5b505afa1580156105ed573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106119190610a78565b915091506106208560016107ac565b83516020015160a00151600090610638908390610dc5565b60405163095ea7b360e01b8152336004820152602481018290529091506001600160a01b0384169063095ea7b390604401602060405180830381600087803b15801561068357600080fd5b505af1158015610697573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106bb9190610aa5565b6107075760405162461bcd60e51b815260206004820152601c60248201527f5661756c743a20546f6b656e20617070726f76616c206661696c65640000000060448201526064016101f1565b505050505050565b60008061071e61010085610db1565b9050600061072e61010086610dff565b6000928352602094909452506040902054600190921b91821690911492915050565b6000816000015161076483602001516107e9565b6040805160208101939093528201526060015b604051602081830303815290604052805190602001209050919050565b6000846107a2858585610823565b1495945050505050565b60006107ba61010084610db1565b905060006107ca61010085610dff565b600092835260209390935250604090208054600190921b909117905550565b80516020808301516040808501516060860151608087015160a088015160c089015160e08a01519551600099610777999098979101610c7b565b600083815b835181101561090957600185821c1661089b578184828151811061085c57634e487b7160e01b600052603260045260246000fd5b602002602001015160405160200161087e929190918252602082015260400190565b6040516020818303038152906040528051906020012091506108f7565b8381815181106108bb57634e487b7160e01b600052603260045260246000fd5b6020026020010151826040516020016108de929190918252602082015260400190565b6040516020818303038152906040528051906020012091505b8061090181610de4565b915050610828565b50949350505050565b600082601f830112610922578081fd5b8135602067ffffffffffffffff82111561093e5761093e610e3f565b8160051b61094d828201610d80565b838152828101908684018388018501891015610967578687fd5b8693505b8584101561098957803583526001939093019291840191840161096b565b50979650505050505050565b600082601f8301126109a5578081fd5b6109ad610d10565b8083856040860111156109be578384fd5b835b60028110156109df5781358452602093840193909101906001016109c0565b509095945050505050565b600082601f8301126109fa578081fd5b813567ffffffffffffffff811115610a1457610a14610e3f565b610a27601f8201601f1916602001610d80565b818152846020838601011115610a3b578283fd5b816020850160208301379081016020019190915292915050565b600060208284031215610a66578081fd5b8151610a7181610e55565b9392505050565b60008060408385031215610a8a578081fd5b8251610a9581610e55565b6020939093015192949293505050565b600060208284031215610ab6578081fd5b81518015158114610a71578182fd5b600060208284031215610ad6578081fd5b8135610a7181610e55565b600060408284031215610af2578081fd5b610afa610d10565b82518152602083015160208201528091505092915050565b600060208284031215610b23578081fd5b5035919050565b60008060408385031215610b3c578182fd5b8235915060208084013567ffffffffffffffff80821115610b5b578384fd5b9085019060608288031215610b6e578384fd5b610b76610d39565b823582811115610b84578586fd5b83016040818a031215610b95578586fd5b610b9d610d10565b813581528582013584811115610bb1578788fd5b9190910190610120828b031215610bc6578687fd5b610bce610d5c565b82358152610bde8b888501610995565b87820152606083013560408201526080830135606082015260a0830135608082015260c083013560a082015260e083013560c082015261010083013585811115610c26578889fd5b610c328c8286016109ea565b60e08301525081870152825250828401359381019390935260408201359281841115610c5c578485fd5b610c6888858501610912565b6040820152809450505050509250929050565b888152600060208083018a835b6002811015610ca557815183529183019190830190600101610c88565b5050508860608401528760808401528660a08401528560c08401528460e08401528351825b81811015610ce75785810183015185820161010001528201610cca565b81811115610cf9578361010083870101525b5092909201610100019a9950505050505050505050565b6040805190810167ffffffffffffffff81118282101715610d3357610d33610e3f565b60405290565b6040516060810167ffffffffffffffff81118282101715610d3357610d33610e3f565b604051610100810167ffffffffffffffff81118282101715610d3357610d33610e3f565b604051601f8201601f1916810167ffffffffffffffff81118282101715610da957610da9610e3f565b604052919050565b600082610dc057610dc0610e29565b500490565b6000816000190483118215151615610ddf57610ddf610e13565b500290565b6000600019821415610df857610df8610e13565b5060010190565b600082610e0e57610e0e610e29565b500690565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052601260045260246000fd5b634e487b7160e01b600052604160045260246000fd5b6001600160a01b0381168114610e6a57600080fd5b5056fea164736f6c6343000804000a",
}

// VaultABI is the input ABI used to generate the binding from.
// Deprecated: Use VaultMetaData.ABI instead.
var VaultABI = VaultMetaData.ABI

// VaultBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use VaultMetaData.Bin instead.
var VaultBin = VaultMetaData.Bin

// DeployVault deploys a new Ethereum contract, binding an instance of Vault to it.
func DeployVault(auth *bind.TransactOpts, backend bind.ContractBackend, _tokenRegistry common.Address, _spokes common.Address) (common.Address, *types.Transaction, *Vault, error) {
	parsed, err := VaultMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VaultBin), backend, _tokenRegistry, _spokes)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Vault{VaultCaller: VaultCaller{contract: contract}, VaultTransactor: VaultTransactor{contract: contract}, VaultFilterer: VaultFilterer{contract: contract}}, nil
}

// Vault is an auto generated Go binding around an Ethereum contract.
type Vault struct {
	VaultCaller     // Read-only binding to the contract
	VaultTransactor // Write-only binding to the contract
	VaultFilterer   // Log filterer for contract events
}

// VaultCaller is an auto generated read-only Go binding around an Ethereum contract.
type VaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VaultTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VaultFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VaultSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VaultSession struct {
	Contract     *Vault            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VaultCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VaultCallerSession struct {
	Contract *VaultCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// VaultTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VaultTransactorSession struct {
	Contract     *VaultTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VaultRaw is an auto generated low-level Go binding around an Ethereum contract.
type VaultRaw struct {
	Contract *Vault // Generic contract binding to access the raw methods on
}

// VaultCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VaultCallerRaw struct {
	Contract *VaultCaller // Generic read-only contract binding to access the raw methods on
}

// VaultTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VaultTransactorRaw struct {
	Contract *VaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVault creates a new instance of Vault, bound to a specific deployed contract.
func NewVault(address common.Address, backend bind.ContractBackend) (*Vault, error) {
	contract, err := bindVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Vault{VaultCaller: VaultCaller{contract: contract}, VaultTransactor: VaultTransactor{contract: contract}, VaultFilterer: VaultFilterer{contract: contract}}, nil
}

// NewVaultCaller creates a new read-only instance of Vault, bound to a specific deployed contract.
func NewVaultCaller(address common.Address, caller bind.ContractCaller) (*VaultCaller, error) {
	contract, err := bindVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VaultCaller{contract: contract}, nil
}

// NewVaultTransactor creates a new write-only instance of Vault, bound to a specific deployed contract.
func NewVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*VaultTransactor, error) {
	contract, err := bindVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VaultTransactor{contract: contract}, nil
}

// NewVaultFilterer creates a new log filterer instance of Vault, bound to a specific deployed contract.
func NewVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*VaultFilterer, error) {
	contract, err := bindVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VaultFilterer{contract: contract}, nil
}

// bindVault binds a generic wrapper to an already deployed contract.
func bindVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VaultABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Vault *VaultRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Vault.Contract.VaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Vault *VaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vault.Contract.VaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Vault *VaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Vault.Contract.VaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Vault *VaultCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Vault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Vault *VaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Vault *VaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Vault.Contract.contract.Transact(opts, method, params...)
}

// IsBatchApproved is a free data retrieval call binding the contract method 0x0d82dcf4.
//
// Solidity: function isBatchApproved(uint256 batchID) view returns(bool)
func (_Vault *VaultCaller) IsBatchApproved(opts *bind.CallOpts, batchID *big.Int) (bool, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "isBatchApproved", batchID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsBatchApproved is a free data retrieval call binding the contract method 0x0d82dcf4.
//
// Solidity: function isBatchApproved(uint256 batchID) view returns(bool)
func (_Vault *VaultSession) IsBatchApproved(batchID *big.Int) (bool, error) {
	return _Vault.Contract.IsBatchApproved(&_Vault.CallOpts, batchID)
}

// IsBatchApproved is a free data retrieval call binding the contract method 0x0d82dcf4.
//
// Solidity: function isBatchApproved(uint256 batchID) view returns(bool)
func (_Vault *VaultCallerSession) IsBatchApproved(batchID *big.Int) (bool, error) {
	return _Vault.Contract.IsBatchApproved(&_Vault.CallOpts, batchID)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Vault *VaultCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Vault *VaultSession) Owner() (common.Address, error) {
	return _Vault.Contract.Owner(&_Vault.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Vault *VaultCallerSession) Owner() (common.Address, error) {
	return _Vault.Contract.Owner(&_Vault.CallOpts)
}

// Rollup is a free data retrieval call binding the contract method 0xcb23bcb5.
//
// Solidity: function rollup() view returns(address)
func (_Vault *VaultCaller) Rollup(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "rollup")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Rollup is a free data retrieval call binding the contract method 0xcb23bcb5.
//
// Solidity: function rollup() view returns(address)
func (_Vault *VaultSession) Rollup() (common.Address, error) {
	return _Vault.Contract.Rollup(&_Vault.CallOpts)
}

// Rollup is a free data retrieval call binding the contract method 0xcb23bcb5.
//
// Solidity: function rollup() view returns(address)
func (_Vault *VaultCallerSession) Rollup() (common.Address, error) {
	return _Vault.Contract.Rollup(&_Vault.CallOpts)
}

// Spokes is a free data retrieval call binding the contract method 0xb640e504.
//
// Solidity: function spokes() view returns(address)
func (_Vault *VaultCaller) Spokes(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "spokes")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Spokes is a free data retrieval call binding the contract method 0xb640e504.
//
// Solidity: function spokes() view returns(address)
func (_Vault *VaultSession) Spokes() (common.Address, error) {
	return _Vault.Contract.Spokes(&_Vault.CallOpts)
}

// Spokes is a free data retrieval call binding the contract method 0xb640e504.
//
// Solidity: function spokes() view returns(address)
func (_Vault *VaultCallerSession) Spokes() (common.Address, error) {
	return _Vault.Contract.Spokes(&_Vault.CallOpts)
}

// TokenRegistry is a free data retrieval call binding the contract method 0x9d23c4c7.
//
// Solidity: function tokenRegistry() view returns(address)
func (_Vault *VaultCaller) TokenRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "tokenRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenRegistry is a free data retrieval call binding the contract method 0x9d23c4c7.
//
// Solidity: function tokenRegistry() view returns(address)
func (_Vault *VaultSession) TokenRegistry() (common.Address, error) {
	return _Vault.Contract.TokenRegistry(&_Vault.CallOpts)
}

// TokenRegistry is a free data retrieval call binding the contract method 0x9d23c4c7.
//
// Solidity: function tokenRegistry() view returns(address)
func (_Vault *VaultCallerSession) TokenRegistry() (common.Address, error) {
	return _Vault.Contract.TokenRegistry(&_Vault.CallOpts)
}

// RequestApproval is a paid mutator transaction binding the contract method 0x4abb54d9.
//
// Solidity: function requestApproval(uint256 batchID, ((bytes32,(bytes32,uint256[2],uint256,bytes32,uint256,uint256,uint256,bytes)),uint256,bytes32[]) commitmentMP) returns()
func (_Vault *VaultTransactor) RequestApproval(opts *bind.TransactOpts, batchID *big.Int, commitmentMP TypesMMCommitmentInclusionProof) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "requestApproval", batchID, commitmentMP)
}

// RequestApproval is a paid mutator transaction binding the contract method 0x4abb54d9.
//
// Solidity: function requestApproval(uint256 batchID, ((bytes32,(bytes32,uint256[2],uint256,bytes32,uint256,uint256,uint256,bytes)),uint256,bytes32[]) commitmentMP) returns()
func (_Vault *VaultSession) RequestApproval(batchID *big.Int, commitmentMP TypesMMCommitmentInclusionProof) (*types.Transaction, error) {
	return _Vault.Contract.RequestApproval(&_Vault.TransactOpts, batchID, commitmentMP)
}

// RequestApproval is a paid mutator transaction binding the contract method 0x4abb54d9.
//
// Solidity: function requestApproval(uint256 batchID, ((bytes32,(bytes32,uint256[2],uint256,bytes32,uint256,uint256,uint256,bytes)),uint256,bytes32[]) commitmentMP) returns()
func (_Vault *VaultTransactorSession) RequestApproval(batchID *big.Int, commitmentMP TypesMMCommitmentInclusionProof) (*types.Transaction, error) {
	return _Vault.Contract.RequestApproval(&_Vault.TransactOpts, batchID, commitmentMP)
}

// SetRollupAddress is a paid mutator transaction binding the contract method 0x07663706.
//
// Solidity: function setRollupAddress(address _rollup) returns()
func (_Vault *VaultTransactor) SetRollupAddress(opts *bind.TransactOpts, _rollup common.Address) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "setRollupAddress", _rollup)
}

// SetRollupAddress is a paid mutator transaction binding the contract method 0x07663706.
//
// Solidity: function setRollupAddress(address _rollup) returns()
func (_Vault *VaultSession) SetRollupAddress(_rollup common.Address) (*types.Transaction, error) {
	return _Vault.Contract.SetRollupAddress(&_Vault.TransactOpts, _rollup)
}

// SetRollupAddress is a paid mutator transaction binding the contract method 0x07663706.
//
// Solidity: function setRollupAddress(address _rollup) returns()
func (_Vault *VaultTransactorSession) SetRollupAddress(_rollup common.Address) (*types.Transaction, error) {
	return _Vault.Contract.SetRollupAddress(&_Vault.TransactOpts, _rollup)
}
