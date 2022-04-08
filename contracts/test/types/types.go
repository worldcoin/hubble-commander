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
	Bin: "0x608060405234801561001057600080fd5b50610431806100206000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80634b2b0063146100465780637a7f59af1461006c578063ed7adbec146100c4575b600080fd5b610059610054366004610238565b610142565b6040519081526020015b60405180910390f35b61005961007a366004610307565b600060f885901b6001600160f81b03191660f085901b60ff60f01b1617605084901b600160501b600160f01b031617603083901b69ffffffff000000000000161795945050505050565b6101136100d2366004610220565b604080518082019091526000815260200181905260f881901c9160ff60f083901c16916001600160a01b03605082901c169160309190911c63ffffffff1690565b604051610063949392919093845260208401929092526001600160a01b03166040830152606082015260800190565b600061014d82610153565b92915050565b6000816000015182602001518360400151846060015160405160200161017c9493929190610350565b604051602081830303815290604052805190602001209050919050565b600082601f8301126101a9578081fd5b813567ffffffffffffffff808211156101c4576101c461040e565b604051601f8301601f19908116603f011681019082821181831017156101ec576101ec61040e565b81604052838152866020858801011115610204578485fd5b8360208701602083013792830160200193909352509392505050565b600060208284031215610231578081fd5b5035919050565b6000602080838503121561024a578182fd5b823567ffffffffffffffff80821115610261578384fd5b9084019060a08287031215610274578384fd5b61027c6103c6565b8235815286603f84011261028e578485fd5b6102966103ef565b80858501606086018a8111156102aa578889fd5b885b60028110156102c9578235855293880193918801916001016102ac565b5096840191909152505092356040840152608082013592818411156102ec578485fd5b6102f887858501610199565b60608201529695505050505050565b6000806000806080858703121561031c578283fd5b843593506020850135925060408501356001600160a01b0381168114610340578283fd5b9396929550929360600135925050565b8481526000602080830186835b600281101561037a5781518352918301919083019060010161035d565b5050508460608401528351825b818110156103a357858101830151858201608001528201610387565b818111156103b45783608083870101525b50929092016080019695505050505050565b6040516080810167ffffffffffffffff811182821017156103e9576103e961040e565b60405290565b6040805190810167ffffffffffffffff811182821017156103e9576103e95b634e487b7160e01b600052604160045260246000fdfea164736f6c6343000804000a",
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
