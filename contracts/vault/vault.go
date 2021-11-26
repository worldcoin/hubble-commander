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
	Bin: "0x60e060405234801561001057600080fd5b50604051610eb5380380610eb583398101604081905261002f91610053565b33606090811b6080526001600160601b031992811b831660c0521b1660a0526100a4565b60008060408385031215610065578182fd5b82516100708161008c565b60208401519092506100818161008c565b809150509250929050565b6001600160a01b03811681146100a157600080fd5b50565b60805160601c60a05160601c60c05160601c610dcd6100e8600039806103f9528061056352508061022c5280610587525080610182528061053f5250610dcd6000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c80638da5cb5b1161005b5780638da5cb5b146100d35780639d23c4c7146100e8578063b640e504146100f0578063cb23bcb5146100f85761007d565b806307663706146100825780630d82dcf4146100975780634abb54d9146100c0575b600080fd5b61009561009036600461093c565b610100565b005b6100aa6100a536600461098b565b6101f7565b6040516100b79190610bc9565b60405180910390f35b6100956100ce3660046109a3565b61020a565b6100db61053d565b6040516100b79190610b9c565b6100db610561565b6100db610585565b6100db6105a9565b600054610100900460ff168061011957506101196105be565b80610127575060005460ff16155b61014c5760405162461bcd60e51b815260040161014390610c54565b60405180910390fd5b600054610100900460ff16158015610177576000805460ff1961ff0019909116610100171660011790555b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146101bf5760405162461bcd60e51b815260040161014390610ce3565b6000805462010000600160b01b031916620100006001600160a01b0385160217905580156101f3576000805461ff00191690555b5050565b60006102048260016105cf565b92915050565b805160200151604090810151905163da2fd3df60e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169163da2fd3df916102609190600401610d78565b60206040518083038186803b15801561027857600080fd5b505afa15801561028c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102b091906108cc565b6001600160a01b0316336001600160a01b0316146102e05760405162461bcd60e51b815260040161014390610d2c565b6102e8610773565b600054604051632d62214160e11b8152620100009091046001600160a01b031690635ac442829061031d908690600401610d78565b604080518083038186803b15801561033457600080fd5b505afa158015610348573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061036c9190610958565b9050610377816105f6565b4310156103965760405162461bcd60e51b815260040161014390610ca2565b805182516103b791906103a890610606565b84602001518560400151610648565b6103d35760405162461bcd60e51b815260040161014390610bd4565b81516020015160800151604051630a7973b760e01b815260009182916001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001691630a7973b79161042d9190600401610d78565b604080518083038186803b15801561044457600080fd5b505afa158015610458573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061047c91906108ef565b9150915061048b856001610660565b83516020015160a0015160405163095ea7b360e01b8152908202906001600160a01b0384169063095ea7b3906104c79033908590600401610bb0565b602060405180830381600087803b1580156104e157600080fd5b505af11580156104f5573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610519919061091c565b6105355760405162461bcd60e51b815260040161014390610c1d565b505050505050565b7f000000000000000000000000000000000000000000000000000000000000000090565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000546201000090046001600160a01b031681565b60006105c930610688565b15905090565b610100820460009081526020919091526040902054600160ff9092169190911b9081161490565b6020015160301c63ffffffff1690565b6000816000015161061a836020015161068e565b60405160200161062b929190610b8e565b604051602081830303815290604052805190602001209050919050565b6000846106568585856106c8565b1495945050505050565b61010082046000908152602091909152604090208054600160ff9093169290921b9091179055565b3b151590565b80516020808301516040808501516060860151608087015160a088015160c089015160e08a0151955160009961062b999098979101610af9565b600083815b835181101561076a57600185821c1661072357818482815181106106ed57fe5b6020026020010151604051602001610706929190610b8e565b604051602081830303815290604052805190602001209150610762565b83818151811061072f57fe5b602002602001015182604051602001610749929190610b8e565b6040516020818303038152906040528051906020012091505b6001016106cd565b50949350505050565b604080518082019091526000808252602082015290565b600082601f83011261079a578081fd5b813567ffffffffffffffff8111156107b0578182fd5b60208082026107c0828201610d81565b838152935081840185830182870184018810156107dc57600080fd5b600092505b848310156107ff5780358252600192909201919083019083016107e1565b505050505092915050565b600082601f83011261081a578081fd5b6108246040610d81565b905080828460408501111561083857600080fd5b60005b600281101561085a57813583526020928301929091019060010161083b565b50505092915050565b600082601f830112610873578081fd5b813567ffffffffffffffff811115610889578182fd5b61089c601f8201601f1916602001610d81565b91508082528360208285010111156108b357600080fd5b8060208401602084013760009082016020015292915050565b6000602082840312156108dd578081fd5b81516108e881610da8565b9392505050565b60008060408385031215610901578081fd5b825161090c81610da8565b6020939093015192949293505050565b60006020828403121561092d578081fd5b815180151581146108e8578182fd5b60006020828403121561094d578081fd5b81356108e881610da8565b600060408284031215610969578081fd5b6109736040610d81565b82518152602083015160208201528091505092915050565b60006020828403121561099c578081fd5b5035919050565b60008060408084860312156109b6578283fd5b8335925060208085013567ffffffffffffffff808211156109d5578485fd5b90860190606082890312156109e8578485fd5b6109f26060610d81565b823582811115610a00578687fd5b8301808a03861315610a10578687fd5b610a1986610d81565b813581528582013584811115610a2d578889fd5b9190910190610120828c031215610a42578788fd5b610100610a4e81610d81565b83358152610a5e8d89860161080a565b888201526060840135898201526080840135606082015260a0840135608082015260c084013560a082015260e084013560c082015281840135915085821115610aa557898afd5b610ab18d838601610863565b60e082015282880152508252508284013593810193909352838201359281841115610ada578586fd5b610ae68985850161078a565b8582015280955050505050509250929050565b600089825260208083018a835b6002811015610b2357815183529183019190830190600101610b06565b5050508860608401528760808401528660a08401528560c08401528460e08401528351825b81811015610b655785810183015185820161010001528201610b48565b81811115610b77578361010083870101525b5092909201610100019a9950505050505050505050565b918252602082015260400190565b6001600160a01b0391909116815260200190565b6001600160a01b03929092168252602082015260400190565b901515815260200190565b60208082526029908201527f5661756c743a20436f6d6d69746d656e74206973206e6f742070726573656e74604082015268040d2dc40c4c2e8c6d60bb1b606082015260800190565b6020808252601c908201527f5661756c743a20546f6b656e20617070726f76616c206661696c656400000000604082015260600190565b6020808252602e908201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160408201526d191e481a5b9a5d1a585b1a5e995960921b606082015260800190565b60208082526021908201527f5661756c743a2042617463682073686f6f756c642062652066696e616c6973656040820152601960fa1b606082015260800190565b60208082526029908201527f496d6d757461626c654f776e61626c653a2063616c6c6572206973206e6f74206040820152683a34329037bbb732b960b91b606082015260800190565b6020808252602c908201527f5661756c743a206d73672e73656e6465722073686f756c64206265207468652060408201526b7461726765742073706f6b6560a01b606082015260800190565b90815260200190565b60405181810167ffffffffffffffff81118282101715610da057600080fd5b604052919050565b6001600160a01b0381168114610dbd57600080fd5b5056fea164736f6c634300060c000a",
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
