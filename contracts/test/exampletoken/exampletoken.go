// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package exampletoken

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

// TestExampleTokenMetaData contains all meta data concerning the TestExampleToken contract.
var TestExampleTokenMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b50604051806040016040528060078152602001664578616d706c6560c81b815250604051806040016040528060038152602001620454d560ec1b81525081600390805190602001906200006692919062000286565b5080516200007c90600490602084019062000286565b50506005805460ff19166012179055506000620000986200010c565b60058054610100600160a81b0319166101006001600160a01b03841690810291909117909155604051919250906000907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a35062000106336b033b2e3c9fd0803ce800000062000110565b62000322565b3390565b6001600160a01b0382166200016c576040805162461bcd60e51b815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015290519081900360640190fd5b6200017a600083836200021f565b62000196816002546200022460201b620008001790919060201c565b6002556001600160a01b03821660009081526020818152604090912054620001c99183906200080062000224821b17901c565b6001600160a01b0383166000818152602081815260408083209490945583518581529351929391927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9281900390910190a35050565b505050565b6000828201838110156200027f576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10620002c957805160ff1916838001178555620002f9565b82800160010185558215620002f9579182015b82811115620002f9578251825591602001919060010190620002dc565b50620003079291506200030b565b5090565b5b808211156200030757600081556001016200030c565b610c7c80620003326000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c8063715018a61161008c578063a457c2d711610066578063a457c2d7146102a2578063a9059cbb146102ce578063dd62ed3e146102fa578063f2fde38b14610328576100ea565b8063715018a61461026c5780638da5cb5b1461027657806395d89b411461029a576100ea565b806323b872dd116100c857806323b872dd146101c6578063313ce567146101fc578063395093511461021a57806370a0823114610246576100ea565b806306fdde03146100ef578063095ea7b31461016c57806318160ddd146101ac575b600080fd5b6100f761034e565b6040805160208082528351818301528351919283929083019185019080838360005b83811015610131578181015183820152602001610119565b50505050905090810190601f16801561015e5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6101986004803603604081101561018257600080fd5b506001600160a01b0381351690602001356103e4565b604080519115158252519081900360200190f35b6101b4610401565b60408051918252519081900360200190f35b610198600480360360608110156101dc57600080fd5b506001600160a01b03813581169160208101359091169060400135610407565b61020461048e565b6040805160ff9092168252519081900360200190f35b6101986004803603604081101561023057600080fd5b506001600160a01b038135169060200135610497565b6101b46004803603602081101561025c57600080fd5b50356001600160a01b03166104e5565b610274610500565b005b61027e6105c4565b604080516001600160a01b039092168252519081900360200190f35b6100f76105d8565b610198600480360360408110156102b857600080fd5b506001600160a01b038135169060200135610639565b610198600480360360408110156102e457600080fd5b506001600160a01b0381351690602001356106a1565b6101b46004803603604081101561031057600080fd5b506001600160a01b03813581169160200135166106b5565b6102746004803603602081101561033e57600080fd5b50356001600160a01b03166106e0565b60038054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156103da5780601f106103af576101008083540402835291602001916103da565b820191906000526020600020905b8154815290600101906020018083116103bd57829003601f168201915b5050505050905090565b60006103f86103f1610861565b8484610865565b50600192915050565b60025490565b6000610414848484610951565b61048484610420610861565b61047f85604051806060016040528060288152602001610bda602891396001600160a01b038a1660009081526001602052604081209061045e610861565b6001600160a01b031681526020810191909152604001600020549190610aac565b610865565b5060019392505050565b60055460ff1690565b60006103f86104a4610861565b8461047f85600160006104b5610861565b6001600160a01b03908116825260208083019390935260409182016000908120918c168152925290205490610800565b6001600160a01b031660009081526020819052604090205490565b610508610861565b6001600160a01b03166105196105c4565b6001600160a01b031614610574576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b60055460405160009161010090046001600160a01b0316907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a360058054610100600160a81b0319169055565b60055461010090046001600160a01b031690565b60048054604080516020601f60026000196101006001881615020190951694909404938401819004810282018101909252828152606093909290918301828280156103da5780601f106103af576101008083540402835291602001916103da565b60006103f8610646610861565b8461047f85604051806060016040528060258152602001610c4b6025913960016000610670610861565b6001600160a01b03908116825260208083019390935260409182016000908120918d16815292529020549190610aac565b60006103f86106ae610861565b8484610951565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b6106e8610861565b6001600160a01b03166106f96105c4565b6001600160a01b031614610754576040805162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015290519081900360640190fd5b6001600160a01b0381166107995760405162461bcd60e51b8152600401808060200182810382526026815260200180610b6c6026913960400191505060405180910390fd5b6005546040516001600160a01b0380841692610100900416907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3600580546001600160a01b0390921661010002610100600160a81b0319909216919091179055565b60008282018381101561085a576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b9392505050565b3390565b6001600160a01b0383166108aa5760405162461bcd60e51b8152600401808060200182810382526024815260200180610c276024913960400191505060405180910390fd5b6001600160a01b0382166108ef5760405162461bcd60e51b8152600401808060200182810382526022815260200180610b926022913960400191505060405180910390fd5b6001600160a01b03808416600081815260016020908152604080832094871680845294825291829020859055815185815291517f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259281900390910190a3505050565b6001600160a01b0383166109965760405162461bcd60e51b8152600401808060200182810382526025815260200180610c026025913960400191505060405180910390fd5b6001600160a01b0382166109db5760405162461bcd60e51b8152600401808060200182810382526023815260200180610b496023913960400191505060405180910390fd5b6109e6838383610b43565b610a2381604051806060016040528060268152602001610bb4602691396001600160a01b0386166000908152602081905260409020549190610aac565b6001600160a01b038085166000908152602081905260408082209390935590841681522054610a529082610800565b6001600160a01b038084166000818152602081815260409182902094909455805185815290519193928716927fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef92918290030190a3505050565b60008184841115610b3b5760405162461bcd60e51b81526004018080602001828103825283818151815260200191508051906020019080838360005b83811015610b00578181015183820152602001610ae8565b50505050905090810190601f168015610b2d5780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b505050900390565b50505056fe45524332303a207472616e7366657220746f20746865207a65726f20616464726573734f776e61626c653a206e6577206f776e657220697320746865207a65726f206164647265737345524332303a20617070726f766520746f20746865207a65726f206164647265737345524332303a207472616e7366657220616d6f756e7420657863656564732062616c616e636545524332303a207472616e7366657220616d6f756e74206578636565647320616c6c6f77616e636545524332303a207472616e736665722066726f6d20746865207a65726f206164647265737345524332303a20617070726f76652066726f6d20746865207a65726f206164647265737345524332303a2064656372656173656420616c6c6f77616e63652062656c6f77207a65726fa164736f6c634300060c000a",
}

// TestExampleTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use TestExampleTokenMetaData.ABI instead.
var TestExampleTokenABI = TestExampleTokenMetaData.ABI

// TestExampleTokenBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestExampleTokenMetaData.Bin instead.
var TestExampleTokenBin = TestExampleTokenMetaData.Bin

// DeployTestExampleToken deploys a new Ethereum contract, binding an instance of TestExampleToken to it.
func DeployTestExampleToken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TestExampleToken, error) {
	parsed, err := TestExampleTokenMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestExampleTokenBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestExampleToken{TestExampleTokenCaller: TestExampleTokenCaller{contract: contract}, TestExampleTokenTransactor: TestExampleTokenTransactor{contract: contract}, TestExampleTokenFilterer: TestExampleTokenFilterer{contract: contract}}, nil
}

// TestExampleToken is an auto generated Go binding around an Ethereum contract.
type TestExampleToken struct {
	TestExampleTokenCaller     // Read-only binding to the contract
	TestExampleTokenTransactor // Write-only binding to the contract
	TestExampleTokenFilterer   // Log filterer for contract events
}

// TestExampleTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestExampleTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestExampleTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestExampleTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestExampleTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestExampleTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestExampleTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestExampleTokenSession struct {
	Contract     *TestExampleToken // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestExampleTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestExampleTokenCallerSession struct {
	Contract *TestExampleTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// TestExampleTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestExampleTokenTransactorSession struct {
	Contract     *TestExampleTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// TestExampleTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestExampleTokenRaw struct {
	Contract *TestExampleToken // Generic contract binding to access the raw methods on
}

// TestExampleTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestExampleTokenCallerRaw struct {
	Contract *TestExampleTokenCaller // Generic read-only contract binding to access the raw methods on
}

// TestExampleTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestExampleTokenTransactorRaw struct {
	Contract *TestExampleTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestExampleToken creates a new instance of TestExampleToken, bound to a specific deployed contract.
func NewTestExampleToken(address common.Address, backend bind.ContractBackend) (*TestExampleToken, error) {
	contract, err := bindTestExampleToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestExampleToken{TestExampleTokenCaller: TestExampleTokenCaller{contract: contract}, TestExampleTokenTransactor: TestExampleTokenTransactor{contract: contract}, TestExampleTokenFilterer: TestExampleTokenFilterer{contract: contract}}, nil
}

// NewTestExampleTokenCaller creates a new read-only instance of TestExampleToken, bound to a specific deployed contract.
func NewTestExampleTokenCaller(address common.Address, caller bind.ContractCaller) (*TestExampleTokenCaller, error) {
	contract, err := bindTestExampleToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestExampleTokenCaller{contract: contract}, nil
}

// NewTestExampleTokenTransactor creates a new write-only instance of TestExampleToken, bound to a specific deployed contract.
func NewTestExampleTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*TestExampleTokenTransactor, error) {
	contract, err := bindTestExampleToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestExampleTokenTransactor{contract: contract}, nil
}

// NewTestExampleTokenFilterer creates a new log filterer instance of TestExampleToken, bound to a specific deployed contract.
func NewTestExampleTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*TestExampleTokenFilterer, error) {
	contract, err := bindTestExampleToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestExampleTokenFilterer{contract: contract}, nil
}

// bindTestExampleToken binds a generic wrapper to an already deployed contract.
func bindTestExampleToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TestExampleTokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestExampleToken *TestExampleTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestExampleToken.Contract.TestExampleTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestExampleToken *TestExampleTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestExampleToken.Contract.TestExampleTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestExampleToken *TestExampleTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestExampleToken.Contract.TestExampleTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestExampleToken *TestExampleTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestExampleToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestExampleToken *TestExampleTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestExampleToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestExampleToken *TestExampleTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestExampleToken.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TestExampleToken *TestExampleTokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TestExampleToken.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TestExampleToken *TestExampleTokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TestExampleToken.Contract.Allowance(&_TestExampleToken.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_TestExampleToken *TestExampleTokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _TestExampleToken.Contract.Allowance(&_TestExampleToken.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TestExampleToken *TestExampleTokenCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TestExampleToken.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TestExampleToken *TestExampleTokenSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _TestExampleToken.Contract.BalanceOf(&_TestExampleToken.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_TestExampleToken *TestExampleTokenCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _TestExampleToken.Contract.BalanceOf(&_TestExampleToken.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TestExampleToken *TestExampleTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _TestExampleToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TestExampleToken *TestExampleTokenSession) Decimals() (uint8, error) {
	return _TestExampleToken.Contract.Decimals(&_TestExampleToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_TestExampleToken *TestExampleTokenCallerSession) Decimals() (uint8, error) {
	return _TestExampleToken.Contract.Decimals(&_TestExampleToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TestExampleToken *TestExampleTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestExampleToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TestExampleToken *TestExampleTokenSession) Name() (string, error) {
	return _TestExampleToken.Contract.Name(&_TestExampleToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TestExampleToken *TestExampleTokenCallerSession) Name() (string, error) {
	return _TestExampleToken.Contract.Name(&_TestExampleToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TestExampleToken *TestExampleTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TestExampleToken.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TestExampleToken *TestExampleTokenSession) Owner() (common.Address, error) {
	return _TestExampleToken.Contract.Owner(&_TestExampleToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TestExampleToken *TestExampleTokenCallerSession) Owner() (common.Address, error) {
	return _TestExampleToken.Contract.Owner(&_TestExampleToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TestExampleToken *TestExampleTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TestExampleToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TestExampleToken *TestExampleTokenSession) Symbol() (string, error) {
	return _TestExampleToken.Contract.Symbol(&_TestExampleToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TestExampleToken *TestExampleTokenCallerSession) Symbol() (string, error) {
	return _TestExampleToken.Contract.Symbol(&_TestExampleToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TestExampleToken *TestExampleTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TestExampleToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TestExampleToken *TestExampleTokenSession) TotalSupply() (*big.Int, error) {
	return _TestExampleToken.Contract.TotalSupply(&_TestExampleToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_TestExampleToken *TestExampleTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _TestExampleToken.Contract.TotalSupply(&_TestExampleToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TestExampleToken *TestExampleTokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestExampleToken.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TestExampleToken *TestExampleTokenSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestExampleToken.Contract.Approve(&_TestExampleToken.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_TestExampleToken *TestExampleTokenTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestExampleToken.Contract.Approve(&_TestExampleToken.TransactOpts, spender, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_TestExampleToken *TestExampleTokenTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _TestExampleToken.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_TestExampleToken *TestExampleTokenSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _TestExampleToken.Contract.DecreaseAllowance(&_TestExampleToken.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_TestExampleToken *TestExampleTokenTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _TestExampleToken.Contract.DecreaseAllowance(&_TestExampleToken.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_TestExampleToken *TestExampleTokenTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _TestExampleToken.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_TestExampleToken *TestExampleTokenSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _TestExampleToken.Contract.IncreaseAllowance(&_TestExampleToken.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_TestExampleToken *TestExampleTokenTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _TestExampleToken.Contract.IncreaseAllowance(&_TestExampleToken.TransactOpts, spender, addedValue)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TestExampleToken *TestExampleTokenTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestExampleToken.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TestExampleToken *TestExampleTokenSession) RenounceOwnership() (*types.Transaction, error) {
	return _TestExampleToken.Contract.RenounceOwnership(&_TestExampleToken.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TestExampleToken *TestExampleTokenTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TestExampleToken.Contract.RenounceOwnership(&_TestExampleToken.TransactOpts)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_TestExampleToken *TestExampleTokenTransactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestExampleToken.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_TestExampleToken *TestExampleTokenSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestExampleToken.Contract.Transfer(&_TestExampleToken.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_TestExampleToken *TestExampleTokenTransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestExampleToken.Contract.Transfer(&_TestExampleToken.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_TestExampleToken *TestExampleTokenTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestExampleToken.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_TestExampleToken *TestExampleTokenSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestExampleToken.Contract.TransferFrom(&_TestExampleToken.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_TestExampleToken *TestExampleTokenTransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TestExampleToken.Contract.TransferFrom(&_TestExampleToken.TransactOpts, sender, recipient, amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TestExampleToken *TestExampleTokenTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TestExampleToken.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TestExampleToken *TestExampleTokenSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TestExampleToken.Contract.TransferOwnership(&_TestExampleToken.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TestExampleToken *TestExampleTokenTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TestExampleToken.Contract.TransferOwnership(&_TestExampleToken.TransactOpts, newOwner)
}

// TestExampleTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the TestExampleToken contract.
type TestExampleTokenApprovalIterator struct {
	Event *TestExampleTokenApproval // Event containing the contract specifics and raw log

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
func (it *TestExampleTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestExampleTokenApproval)
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
		it.Event = new(TestExampleTokenApproval)
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
func (it *TestExampleTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestExampleTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestExampleTokenApproval represents a Approval event raised by the TestExampleToken contract.
type TestExampleTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TestExampleToken *TestExampleTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*TestExampleTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TestExampleToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &TestExampleTokenApprovalIterator{contract: _TestExampleToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TestExampleToken *TestExampleTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TestExampleTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _TestExampleToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestExampleTokenApproval)
				if err := _TestExampleToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_TestExampleToken *TestExampleTokenFilterer) ParseApproval(log types.Log) (*TestExampleTokenApproval, error) {
	event := new(TestExampleTokenApproval)
	if err := _TestExampleToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestExampleTokenOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TestExampleToken contract.
type TestExampleTokenOwnershipTransferredIterator struct {
	Event *TestExampleTokenOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *TestExampleTokenOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestExampleTokenOwnershipTransferred)
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
		it.Event = new(TestExampleTokenOwnershipTransferred)
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
func (it *TestExampleTokenOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestExampleTokenOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestExampleTokenOwnershipTransferred represents a OwnershipTransferred event raised by the TestExampleToken contract.
type TestExampleTokenOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TestExampleToken *TestExampleTokenFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TestExampleTokenOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TestExampleToken.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TestExampleTokenOwnershipTransferredIterator{contract: _TestExampleToken.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TestExampleToken *TestExampleTokenFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TestExampleTokenOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TestExampleToken.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestExampleTokenOwnershipTransferred)
				if err := _TestExampleToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TestExampleToken *TestExampleTokenFilterer) ParseOwnershipTransferred(log types.Log) (*TestExampleTokenOwnershipTransferred, error) {
	event := new(TestExampleTokenOwnershipTransferred)
	if err := _TestExampleToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TestExampleTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the TestExampleToken contract.
type TestExampleTokenTransferIterator struct {
	Event *TestExampleTokenTransfer // Event containing the contract specifics and raw log

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
func (it *TestExampleTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TestExampleTokenTransfer)
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
		it.Event = new(TestExampleTokenTransfer)
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
func (it *TestExampleTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TestExampleTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TestExampleTokenTransfer represents a Transfer event raised by the TestExampleToken contract.
type TestExampleTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TestExampleToken *TestExampleTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*TestExampleTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestExampleToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &TestExampleTokenTransferIterator{contract: _TestExampleToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TestExampleToken *TestExampleTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TestExampleTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _TestExampleToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TestExampleTokenTransfer)
				if err := _TestExampleToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_TestExampleToken *TestExampleTokenFilterer) ParseTransfer(log types.Log) (*TestExampleTokenTransfer, error) {
	event := new(TestExampleTokenTransfer)
	if err := _TestExampleToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
