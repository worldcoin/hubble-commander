// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vault

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

// VaultABI is the input ABI used to generate the binding from.
const VaultABI = "[{\"inputs\":[{\"internalType\":\"contractITokenRegistry\",\"name\":\"_tokenRegistry\",\"type\":\"address\"},{\"internalType\":\"contractSpokeRegistry\",\"name\":\"_spokes\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"}],\"name\":\"isBatchApproved\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"spokeID\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"withdrawRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.MassMigrationBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"internalType\":\"structTypes.MassMigrationCommitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.MMCommitmentInclusionProof\",\"name\":\"commitmentMP\",\"type\":\"tuple\"}],\"name\":\"requestApproval\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollup\",\"outputs\":[{\"internalType\":\"contractRollup\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractRollup\",\"name\":\"_rollup\",\"type\":\"address\"}],\"name\":\"setRollupAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"spokes\",\"outputs\":[{\"internalType\":\"contractSpokeRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenRegistry\",\"outputs\":[{\"internalType\":\"contractITokenRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// VaultBin is the compiled bytecode used for deploying new contracts.
var VaultBin = "0x608060405234801561001057600080fd5b50604051610c45380380610c4583398101604081905261002f91610060565b600280546001600160a01b039384166001600160a01b031991821617909155600180549290931691161790556100b1565b60008060408385031215610072578182fd5b825161007d81610099565b602084015190925061008e81610099565b809150509250929050565b6001600160a01b03811681146100ae57600080fd5b50565b610b85806100c06000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c806307663706146100675780630d82dcf41461007c5780634abb54d9146100a55780639d23c4c7146100b8578063b640e504146100cd578063cb23bcb5146100d5575b600080fd5b61007a610075366004610805565b6100dd565b005b61008f61008a366004610854565b6100ff565b60405161009c9190610a04565b60405180910390f35b61007a6100b336600461086c565b610112565b6100c0610403565b60405161009c9190610a0f565b6100c0610412565b6100c0610421565b600080546001600160a01b0319166001600160a01b0392909216919091179055565b600061010c826003610430565b92915050565b600154815160200151604090810151905163da2fd3df60e01b81526001600160a01b039092169163da2fd3df9161014b91600401610b30565b60206040518083038186803b15801561016357600080fd5b505afa158015610177573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061019b91906107c2565b6001600160a01b0316336001600160a01b0316146101d45760405162461bcd60e51b81526004016101cb90610ae4565b60405180910390fd5b6101dc6105ce565b600054604051632d62214160e11b81526001600160a01b0390911690635ac442829061020c908690600401610b30565b604080518083038186803b15801561022357600080fd5b505afa158015610237573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061025b9190610821565b905061026681610457565b4310156102855760405162461bcd60e51b81526004016101cb90610aa3565b805182516102a6919061029790610467565b846020015185604001516104a9565b6102c25760405162461bcd60e51b81526004016101cb90610a23565b600254825160200151608001516040516320f5ab4f60e11b81526000926001600160a01b0316916341eb569e916102fc9190600401610b30565b60206040518083038186803b15801561031457600080fd5b505afa158015610328573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061034c91906107c2565b90506103598460036104c1565b82516020015160a0015160405163095ea7b360e01b81526001600160a01b0383169163095ea7b39161038f9133916004016109eb565b602060405180830381600087803b1580156103a957600080fd5b505af11580156103bd573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103e191906107e5565b6103fd5760405162461bcd60e51b81526004016101cb90610a6c565b50505050565b6002546001600160a01b031681565b6001546001600160a01b031681565b6000546001600160a01b031681565b610100820460009081526020919091526040902054600160ff9092169190911b9081161490565b6020015160301c63ffffffff1690565b6000816000015161047b83602001516104e9565b60405160200161048c9291906109dd565b604051602081830303815290604052805190602001209050919050565b6000846104b7858585610523565b1495945050505050565b61010082046000908152602091909152604090208054600160ff9093169290921b9091179055565b80516020808301516040808501516060860151608087015160a088015160c089015160e08a0151955160009961048c999098979101610948565b600083815b83518110156105c557600185821c1661057e578184828151811061054857fe5b60200260200101516040516020016105619291906109dd565b6040516020818303038152906040528051906020012091506105bd565b83818151811061058a57fe5b6020026020010151826040516020016105a49291906109dd565b6040516020818303038152906040528051906020012091505b600101610528565b50949350505050565b604080518082019091526000808252602082015290565b600082601f8301126105f5578081fd5b813567ffffffffffffffff81111561060b578182fd5b602080820261061b828201610b39565b8381529350818401858301828701840188101561063757600080fd5b600092505b8483101561065a57803582526001929092019190830190830161063c565b505050505092915050565b600082601f830112610675578081fd5b61067f6040610b39565b905080828460408501111561069357600080fd5b60005b60028110156106b5578135835260209283019290910190600101610696565b50505092915050565b600082601f8301126106ce578081fd5b813567ffffffffffffffff8111156106e4578182fd5b6106f7601f8201601f1916602001610b39565b915080825283602082850101111561070e57600080fd5b8060208401602084013760009082016020015292915050565b60006101208284031215610739578081fd5b61010061074581610b39565b9150823582526107588460208501610665565b6020830152606083013560408301526080830135606083015260a0830135608083015260c083013560a083015260e083013560c083015280830135905067ffffffffffffffff8111156107aa57600080fd5b6107b6848285016106be565b60e08301525092915050565b6000602082840312156107d3578081fd5b81516107de81610b60565b9392505050565b6000602082840312156107f6578081fd5b815180151581146107de578182fd5b600060208284031215610816578081fd5b81356107de81610b60565b600060408284031215610832578081fd5b61083c6040610b39565b82518152602083015160208201528091505092915050565b600060208284031215610865578081fd5b5035919050565b6000806040838503121561087e578081fd5b82359150602083013567ffffffffffffffff8082111561089c578283fd5b90840190606082870312156108af578283fd5b6108b96060610b39565b8235828111156108c7578485fd5b8301604081890312156108d8578485fd5b6108e26040610b39565b813581526020820135848111156108f7578687fd5b6109038a828501610727565b602083015250808352505060208301356020820152604083013582811115610929578485fd5b610935888286016105e5565b6040830152508093505050509250929050565b600089825260208083018a835b600281101561097257815183529183019190830190600101610955565b5050508860608401528760808401528660a08401528560c08401528460e08401528351825b818110156109b45785810183015185820161010001528201610997565b818111156109c6578361010083870101525b5092909201610100019a9950505050505050505050565b918252602082015260400190565b6001600160a01b03929092168252602082015260400190565b901515815260200190565b6001600160a01b0391909116815260200190565b60208082526029908201527f5661756c743a20436f6d6d69746d656e74206973206e6f742070726573656e74604082015268040d2dc40c4c2e8c6d60bb1b606082015260800190565b6020808252601c908201527f5661756c743a20546f6b656e20617070726f76616c206661696c656400000000604082015260600190565b60208082526021908201527f5661756c743a2042617463682073686f6f756c642062652066696e616c6973656040820152601960fa1b606082015260800190565b6020808252602c908201527f5661756c743a206d73672e73656e6465722073686f756c64206265207468652060408201526b7461726765742073706f6b6560a01b606082015260800190565b90815260200190565b60405181810167ffffffffffffffff81118282101715610b5857600080fd5b604052919050565b6001600160a01b0381168114610b7557600080fd5b5056fea164736f6c634300060c000a"

// DeployVault deploys a new Ethereum contract, binding an instance of Vault to it.
func DeployVault(auth *bind.TransactOpts, backend bind.ContractBackend, _tokenRegistry common.Address, _spokes common.Address) (common.Address, *types.Transaction, *Vault, error) {
	parsed, err := abi.JSON(strings.NewReader(VaultABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(VaultBin), backend, _tokenRegistry, _spokes)
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
