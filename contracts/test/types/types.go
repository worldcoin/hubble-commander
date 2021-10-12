// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package types

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

// TypesTransferBody is an auto generated low-level Go binding around an user-defined struct.
type TypesTransferBody struct {
	AccountRoot [32]byte
	Signature   [2]*big.Int
	FeeReceiver *big.Int
	Txs         []byte
}

// TestTypesMetaData contains all meta data concerning the TestTypes contract.
var TestTypesMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"meta\",\"type\":\"bytes32\"}],\"name\":\"decodeMeta\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"batchType\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"size\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"committer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"finaliseOn\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchType\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"commitmentLength\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"committer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"finaliseOn\",\"type\":\"uint256\"}],\"name\":\"encodeMeta\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.TransferBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"name\":\"hashTransferBody\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506104a0806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80634b2b0063146100465780637a7f59af1461006f578063ed7adbec14610082575b600080fd5b610059610054366004610287565b6100a5565b6040516100669190610422565b60405180910390f35b61005961007d366004610363565b6100b6565b61009561009036600461026f565b6100cd565b604051610066949392919061042b565b60006100b082610125565b92915050565b60006100c48585858561016b565b95945050505050565b6000806000806100db6101ef565b506040805180820190915260008152602081018690526100fa816101b2565b9450610105816101bc565b9350610110816101c9565b925061011b816101df565b9150509193509193565b6000816000015182602001518360400151846060015160405160200161014e94939291906103ac565b604051602081830303815290604052805190602001209050919050565b69ffffffff000000000000603082901b16600160501b600160f01b03605084901b1660ff60f01b60f086901b166001600160f81b031960f888901b16171717949350505050565b6020015160f81c90565b6020015160f01c60ff1690565b602081015160501c6001600160a01b0316919050565b6020015160301c63ffffffff1690565b604080518082019091526000808252602082015290565b600082601f830112610216578081fd5b813567ffffffffffffffff81111561022c578182fd5b61023f601f8201601f191660200161044f565b915080825283602082850101111561025657600080fd5b8060208401602084013760009082016020015292915050565b600060208284031215610280578081fd5b5035919050565b60006020808385031215610299578182fd5b823567ffffffffffffffff808211156102b0578384fd5b9084019060a082870312156102c3578384fd5b6102cd608061044f565b8235815286603f8401126102df578485fd5b60026102f26102ed82610476565b61044f565b80868601606087018b81111561030657898afd5b895b8581101561032457823585529389019391890191600101610308565b50978501919091525050509235604084015260808201359281841115610348578485fd5b61035487858501610206565b60608201529695505050505050565b60008060008060808587031215610378578283fd5b843593506020850135925060408501356001600160a01b038116811461039c578283fd5b9396929550929360600135925050565b6000858252602080830186835b60028110156103d6578151835291830191908301906001016103b9565b5050508460608401528351825b818110156103ff578581018301518582016080015282016103e3565b818111156104105783608083870101525b50929092016080019695505050505050565b90815260200190565b93845260208401929092526001600160a01b03166040830152606082015260800190565b60405181810167ffffffffffffffff8111828210171561046e57600080fd5b604052919050565b600067ffffffffffffffff82111561048c578081fd5b506020029056fea164736f6c634300060c000a",
}

// TestTypesABI is the input ABI used to generate the binding from.
// Deprecated: Use TestTypesMetaData.ABI instead.
var TestTypesABI = TestTypesMetaData.ABI

// TestTypesBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestTypesMetaData.Bin instead.
var TestTypesBin = TestTypesMetaData.Bin

// DeployTestTypes deploys a new Ethereum contract, binding an instance of TestTypes to it.
func DeployTestTypes(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TestTypes, error) {
	parsed, err := TestTypesMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestTypesBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestTypes{TestTypesCaller: TestTypesCaller{contract: contract}, TestTypesTransactor: TestTypesTransactor{contract: contract}, TestTypesFilterer: TestTypesFilterer{contract: contract}}, nil
}

// TestTypes is an auto generated Go binding around an Ethereum contract.
type TestTypes struct {
	TestTypesCaller     // Read-only binding to the contract
	TestTypesTransactor // Write-only binding to the contract
	TestTypesFilterer   // Log filterer for contract events
}

// TestTypesCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestTypesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestTypesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestTypesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestTypesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestTypesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestTypesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestTypesSession struct {
	Contract     *TestTypes        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestTypesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestTypesCallerSession struct {
	Contract *TestTypesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// TestTypesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestTypesTransactorSession struct {
	Contract     *TestTypesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// TestTypesRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestTypesRaw struct {
	Contract *TestTypes // Generic contract binding to access the raw methods on
}

// TestTypesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestTypesCallerRaw struct {
	Contract *TestTypesCaller // Generic read-only contract binding to access the raw methods on
}

// TestTypesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestTypesTransactorRaw struct {
	Contract *TestTypesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestTypes creates a new instance of TestTypes, bound to a specific deployed contract.
func NewTestTypes(address common.Address, backend bind.ContractBackend) (*TestTypes, error) {
	contract, err := bindTestTypes(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestTypes{TestTypesCaller: TestTypesCaller{contract: contract}, TestTypesTransactor: TestTypesTransactor{contract: contract}, TestTypesFilterer: TestTypesFilterer{contract: contract}}, nil
}

// NewTestTypesCaller creates a new read-only instance of TestTypes, bound to a specific deployed contract.
func NewTestTypesCaller(address common.Address, caller bind.ContractCaller) (*TestTypesCaller, error) {
	contract, err := bindTestTypes(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestTypesCaller{contract: contract}, nil
}

// NewTestTypesTransactor creates a new write-only instance of TestTypes, bound to a specific deployed contract.
func NewTestTypesTransactor(address common.Address, transactor bind.ContractTransactor) (*TestTypesTransactor, error) {
	contract, err := bindTestTypes(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestTypesTransactor{contract: contract}, nil
}

// NewTestTypesFilterer creates a new log filterer instance of TestTypes, bound to a specific deployed contract.
func NewTestTypesFilterer(address common.Address, filterer bind.ContractFilterer) (*TestTypesFilterer, error) {
	contract, err := bindTestTypes(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestTypesFilterer{contract: contract}, nil
}

// bindTestTypes binds a generic wrapper to an already deployed contract.
func bindTestTypes(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TestTypesABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestTypes *TestTypesRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestTypes.Contract.TestTypesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestTypes *TestTypesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestTypes.Contract.TestTypesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestTypes *TestTypesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestTypes.Contract.TestTypesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestTypes *TestTypesCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestTypes.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestTypes *TestTypesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestTypes.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestTypes *TestTypesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestTypes.Contract.contract.Transact(opts, method, params...)
}

// DecodeMeta is a free data retrieval call binding the contract method 0xed7adbec.
//
// Solidity: function decodeMeta(bytes32 meta) pure returns(uint256 batchType, uint256 size, address committer, uint256 finaliseOn)
func (_TestTypes *TestTypesCaller) DecodeMeta(opts *bind.CallOpts, meta [32]byte) (struct {
	BatchType  *big.Int
	Size       *big.Int
	Committer  common.Address
	FinaliseOn *big.Int
}, error) {
	var out []interface{}
	err := _TestTypes.contract.Call(opts, &out, "decodeMeta", meta)

	outstruct := new(struct {
		BatchType  *big.Int
		Size       *big.Int
		Committer  common.Address
		FinaliseOn *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.BatchType = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Size = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Committer = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.FinaliseOn = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// DecodeMeta is a free data retrieval call binding the contract method 0xed7adbec.
//
// Solidity: function decodeMeta(bytes32 meta) pure returns(uint256 batchType, uint256 size, address committer, uint256 finaliseOn)
func (_TestTypes *TestTypesSession) DecodeMeta(meta [32]byte) (struct {
	BatchType  *big.Int
	Size       *big.Int
	Committer  common.Address
	FinaliseOn *big.Int
}, error) {
	return _TestTypes.Contract.DecodeMeta(&_TestTypes.CallOpts, meta)
}

// DecodeMeta is a free data retrieval call binding the contract method 0xed7adbec.
//
// Solidity: function decodeMeta(bytes32 meta) pure returns(uint256 batchType, uint256 size, address committer, uint256 finaliseOn)
func (_TestTypes *TestTypesCallerSession) DecodeMeta(meta [32]byte) (struct {
	BatchType  *big.Int
	Size       *big.Int
	Committer  common.Address
	FinaliseOn *big.Int
}, error) {
	return _TestTypes.Contract.DecodeMeta(&_TestTypes.CallOpts, meta)
}

// EncodeMeta is a free data retrieval call binding the contract method 0x7a7f59af.
//
// Solidity: function encodeMeta(uint256 batchType, uint256 commitmentLength, address committer, uint256 finaliseOn) pure returns(bytes32)
func (_TestTypes *TestTypesCaller) EncodeMeta(opts *bind.CallOpts, batchType *big.Int, commitmentLength *big.Int, committer common.Address, finaliseOn *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _TestTypes.contract.Call(opts, &out, "encodeMeta", batchType, commitmentLength, committer, finaliseOn)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EncodeMeta is a free data retrieval call binding the contract method 0x7a7f59af.
//
// Solidity: function encodeMeta(uint256 batchType, uint256 commitmentLength, address committer, uint256 finaliseOn) pure returns(bytes32)
func (_TestTypes *TestTypesSession) EncodeMeta(batchType *big.Int, commitmentLength *big.Int, committer common.Address, finaliseOn *big.Int) ([32]byte, error) {
	return _TestTypes.Contract.EncodeMeta(&_TestTypes.CallOpts, batchType, commitmentLength, committer, finaliseOn)
}

// EncodeMeta is a free data retrieval call binding the contract method 0x7a7f59af.
//
// Solidity: function encodeMeta(uint256 batchType, uint256 commitmentLength, address committer, uint256 finaliseOn) pure returns(bytes32)
func (_TestTypes *TestTypesCallerSession) EncodeMeta(batchType *big.Int, commitmentLength *big.Int, committer common.Address, finaliseOn *big.Int) ([32]byte, error) {
	return _TestTypes.Contract.EncodeMeta(&_TestTypes.CallOpts, batchType, commitmentLength, committer, finaliseOn)
}

// HashTransferBody is a free data retrieval call binding the contract method 0x4b2b0063.
//
// Solidity: function hashTransferBody((bytes32,uint256[2],uint256,bytes) body) pure returns(bytes32)
func (_TestTypes *TestTypesCaller) HashTransferBody(opts *bind.CallOpts, body TypesTransferBody) ([32]byte, error) {
	var out []interface{}
	err := _TestTypes.contract.Call(opts, &out, "hashTransferBody", body)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// HashTransferBody is a free data retrieval call binding the contract method 0x4b2b0063.
//
// Solidity: function hashTransferBody((bytes32,uint256[2],uint256,bytes) body) pure returns(bytes32)
func (_TestTypes *TestTypesSession) HashTransferBody(body TypesTransferBody) ([32]byte, error) {
	return _TestTypes.Contract.HashTransferBody(&_TestTypes.CallOpts, body)
}

// HashTransferBody is a free data retrieval call binding the contract method 0x4b2b0063.
//
// Solidity: function hashTransferBody((bytes32,uint256[2],uint256,bytes) body) pure returns(bytes32)
func (_TestTypes *TestTypesCallerSession) HashTransferBody(body TypesTransferBody) ([32]byte, error) {
	return _TestTypes.Contract.HashTransferBody(&_TestTypes.CallOpts, body)
}
