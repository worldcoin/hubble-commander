// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package proofofauthority

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

// ProofOfAuthorityMetaData contains all meta data concerning the ProofOfAuthority contract.
var ProofOfAuthorityMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_proposers\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"getProposer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506040516101a63803806101a68339818101604052602081101561003357600080fd5b810190808051604051939291908464010000000082111561005357600080fd5b90830190602082018581111561006857600080fd5b825186602082028301116401000000008211171561008557600080fd5b82525081516020918201928201910280838360005b838110156100b257818101518382015260200161009a565b5050505090500160405250505060005b81518110156101165760016000808484815181106100dc57fe5b6020908102919091018101516001600160a01b03168252810191909152604001600020805460ff19169115159190911790556001016100c2565b50506080806101266000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063e9790d0214602d575b600080fd5b6033604f565b604080516001600160a01b039092168252519081900360200190f35b3360009081526020819052604081205460ff1615606c5750336070565b5060005b9056fea164736f6c634300060c000a",
}

// ProofOfAuthorityABI is the input ABI used to generate the binding from.
// Deprecated: Use ProofOfAuthorityMetaData.ABI instead.
var ProofOfAuthorityABI = ProofOfAuthorityMetaData.ABI

// ProofOfAuthorityBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ProofOfAuthorityMetaData.Bin instead.
var ProofOfAuthorityBin = ProofOfAuthorityMetaData.Bin

// DeployProofOfAuthority deploys a new Ethereum contract, binding an instance of ProofOfAuthority to it.
func DeployProofOfAuthority(auth *bind.TransactOpts, backend bind.ContractBackend, _proposers []common.Address) (common.Address, *types.Transaction, *ProofOfAuthority, error) {
	parsed, err := ProofOfAuthorityMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ProofOfAuthorityBin), backend, _proposers)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ProofOfAuthority{ProofOfAuthorityCaller: ProofOfAuthorityCaller{contract: contract}, ProofOfAuthorityTransactor: ProofOfAuthorityTransactor{contract: contract}, ProofOfAuthorityFilterer: ProofOfAuthorityFilterer{contract: contract}}, nil
}

// ProofOfAuthority is an auto generated Go binding around an Ethereum contract.
type ProofOfAuthority struct {
	ProofOfAuthorityCaller     // Read-only binding to the contract
	ProofOfAuthorityTransactor // Write-only binding to the contract
	ProofOfAuthorityFilterer   // Log filterer for contract events
}

// ProofOfAuthorityCaller is an auto generated read-only Go binding around an Ethereum contract.
type ProofOfAuthorityCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofOfAuthorityTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ProofOfAuthorityTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofOfAuthorityFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ProofOfAuthorityFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ProofOfAuthoritySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ProofOfAuthoritySession struct {
	Contract     *ProofOfAuthority // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ProofOfAuthorityCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ProofOfAuthorityCallerSession struct {
	Contract *ProofOfAuthorityCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// ProofOfAuthorityTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ProofOfAuthorityTransactorSession struct {
	Contract     *ProofOfAuthorityTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ProofOfAuthorityRaw is an auto generated low-level Go binding around an Ethereum contract.
type ProofOfAuthorityRaw struct {
	Contract *ProofOfAuthority // Generic contract binding to access the raw methods on
}

// ProofOfAuthorityCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ProofOfAuthorityCallerRaw struct {
	Contract *ProofOfAuthorityCaller // Generic read-only contract binding to access the raw methods on
}

// ProofOfAuthorityTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ProofOfAuthorityTransactorRaw struct {
	Contract *ProofOfAuthorityTransactor // Generic write-only contract binding to access the raw methods on
}

// NewProofOfAuthority creates a new instance of ProofOfAuthority, bound to a specific deployed contract.
func NewProofOfAuthority(address common.Address, backend bind.ContractBackend) (*ProofOfAuthority, error) {
	contract, err := bindProofOfAuthority(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ProofOfAuthority{ProofOfAuthorityCaller: ProofOfAuthorityCaller{contract: contract}, ProofOfAuthorityTransactor: ProofOfAuthorityTransactor{contract: contract}, ProofOfAuthorityFilterer: ProofOfAuthorityFilterer{contract: contract}}, nil
}

// NewProofOfAuthorityCaller creates a new read-only instance of ProofOfAuthority, bound to a specific deployed contract.
func NewProofOfAuthorityCaller(address common.Address, caller bind.ContractCaller) (*ProofOfAuthorityCaller, error) {
	contract, err := bindProofOfAuthority(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ProofOfAuthorityCaller{contract: contract}, nil
}

// NewProofOfAuthorityTransactor creates a new write-only instance of ProofOfAuthority, bound to a specific deployed contract.
func NewProofOfAuthorityTransactor(address common.Address, transactor bind.ContractTransactor) (*ProofOfAuthorityTransactor, error) {
	contract, err := bindProofOfAuthority(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ProofOfAuthorityTransactor{contract: contract}, nil
}

// NewProofOfAuthorityFilterer creates a new log filterer instance of ProofOfAuthority, bound to a specific deployed contract.
func NewProofOfAuthorityFilterer(address common.Address, filterer bind.ContractFilterer) (*ProofOfAuthorityFilterer, error) {
	contract, err := bindProofOfAuthority(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ProofOfAuthorityFilterer{contract: contract}, nil
}

// bindProofOfAuthority binds a generic wrapper to an already deployed contract.
func bindProofOfAuthority(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ProofOfAuthorityABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProofOfAuthority *ProofOfAuthorityRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ProofOfAuthority.Contract.ProofOfAuthorityCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProofOfAuthority *ProofOfAuthorityRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProofOfAuthority.Contract.ProofOfAuthorityTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProofOfAuthority *ProofOfAuthorityRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProofOfAuthority.Contract.ProofOfAuthorityTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ProofOfAuthority *ProofOfAuthorityCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ProofOfAuthority.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ProofOfAuthority *ProofOfAuthorityTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ProofOfAuthority.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ProofOfAuthority *ProofOfAuthorityTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ProofOfAuthority.Contract.contract.Transact(opts, method, params...)
}

// GetProposer is a free data retrieval call binding the contract method 0xe9790d02.
//
// Solidity: function getProposer() view returns(address)
func (_ProofOfAuthority *ProofOfAuthorityCaller) GetProposer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ProofOfAuthority.contract.Call(opts, &out, "getProposer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetProposer is a free data retrieval call binding the contract method 0xe9790d02.
//
// Solidity: function getProposer() view returns(address)
func (_ProofOfAuthority *ProofOfAuthoritySession) GetProposer() (common.Address, error) {
	return _ProofOfAuthority.Contract.GetProposer(&_ProofOfAuthority.CallOpts)
}

// GetProposer is a free data retrieval call binding the contract method 0xe9790d02.
//
// Solidity: function getProposer() view returns(address)
func (_ProofOfAuthority *ProofOfAuthorityCallerSession) GetProposer() (common.Address, error) {
	return _ProofOfAuthority.Contract.GetProposer(&_ProofOfAuthority.CallOpts)
}
