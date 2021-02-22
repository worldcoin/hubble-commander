// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package depositmanager

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

// DepositManagerABI is the input ABI used to generate the binding from.
const DepositManagerABI = "[{\"inputs\":[{\"internalType\":\"contractITokenRegistry\",\"name\":\"_tokenRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_vault\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"maxSubtreeDepth\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"DepositQueued\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"DepositSubTreeReady\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"babyTrees\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"babyTreesLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"back\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"}],\"name\":\"depositFor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dequeueToSubmit\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"subtreeRoot\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"front\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paramMaxSubtreeSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"queue\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"subtreeRoot\",\"type\":\"bytes32\"}],\"name\":\"reenqueue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollup\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_rollup\",\"type\":\"address\"}],\"name\":\"setRollupAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenRegistry\",\"outputs\":[{\"internalType\":\"contractITokenRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"vault\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// DepositManagerBin is the compiled bytecode used for deploying new contracts.
var DepositManagerBin = "0x608060405260018055600060025560006004556000600555600260065534801561002857600080fd5b50604051610b25380380610b2583398101604081905261004791610083565b600980546001600160a01b039485166001600160a01b03199182161790915560078054939094169216919091179091556001901b6006556100dd565b600080600060608486031215610097578283fd5b83516100a2816100c5565b60208501519093506100b3816100c5565b80925050604084015190509250925092565b6001600160a01b03811681146100da57600080fd5b50565b610a39806100ec6000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c8063ba75bbd81161008c578063d86ee48d11610066578063d86ee48d14610185578063ddf0b0091461018d578063ee9d68ce146101a0578063fbfa77cf146101b3576100ea565b8063ba75bbd81461016d578063c7accaa414610175578063cb23bcb51461017d576100ea565b80638dde0840116100c85780638dde08401461012a578063966fda62146101325780639d23c4c714610145578063ade000261461015a576100ea565b806307663706146100ef5780632dfdf0b514610104578063425e97f214610122575b600080fd5b6101026100fd366004610738565b6101bb565b005b61010c6101dd565b604051610119919061085f565b60405180910390f35b61010c6101e3565b61010c6101e9565b6101026101403660046107c7565b6101ef565b61014d6104a7565b604051610119919061080d565b610102610168366004610797565b6104b6565b61010c6104ec565b61010c6104f2565b61014d6104f8565b61010c610507565b61010c61019b366004610797565b610541565b61010c6101ae366004610797565b610553565b61014d610565565b600880546001600160a01b0319166001600160a01b0392909216919091179055565b60055481565b60045481565b60025481565b600082116102185760405162461bcd60e51b815260040161020f90610941565b60405180910390fd5b6009546040516320f5ab4f60e11b81526000916001600160a01b0316906341eb569e9061024990859060040161085f565b60206040518083038186803b15801561026157600080fd5b505afa158015610275573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610299919061075b565b905082816001600160a01b031663dd62ed3e33306040518363ffffffff1660e01b81526004016102ca929190610821565b60206040518083038186803b1580156102e257600080fd5b505afa1580156102f6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061031a91906107af565b10156103385760405162461bcd60e51b815260040161020f90610876565b6007546040516323b872dd60e01b81526001600160a01b03808416926323b872dd9261036c9233921690889060040161083b565b602060405180830381600087803b15801561038657600080fd5b505af115801561039a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103be9190610777565b6103da5760405162461bcd60e51b815260040161020f90610985565b6103e2610710565b604051806080016040528086815260200184815260200185815260200160008152509050606061041182610574565b90507f5a1922090c12e28c132a961f6bb4d74350598a62e8981b5eb9bb5ccbbce9df9486826040516104449291906109bc565b60405180910390a1600061045e82805190602001206105b3565b9050801561049e577f4744f3a44c5716c9fa423a71cdaa806771a8bf469f4c007ca338b8e8e202a8b581604051610495919061085f565b60405180910390a15b50505050505050565b6009546001600160a01b031681565b6008546001600160a01b031633146104e05760405162461bcd60e51b815260040161020f906108ad565b6104e9816106ac565b50565b60015481565b60065481565b6008546001600160a01b031681565b6008546000906001600160a01b031633146105345760405162461bcd60e51b815260040161020f906108ad565b61053c6106c8565b905090565b60006020819052908152604090205481565b60036020526000908152604090205481565b6007546001600160a01b031681565b6060816000015182602001518360400151846060015160405160200161059d94939291906107f2565b6040516020818303038152906040529050919050565b6005805460019081019182905560045460008181526003602052604081208590559291015b6001821661064a5760011981016000908152600360209081526040808320546000198501845292819020549051610610939201610868565b60408051601f19818403018152918152815160209283012060011984016000908152600390935291205560019190911c90600019016105d8565b600481905560065460055414156106a0576000805260036020527f3617319a054d772f909f7c479a2cebe5066e836a939412e32403c99029b92eff549250610691836106ac565b600060048190556005556106a5565b600092505b5050919050565b6002805460010190819055600090815260208190526040902055565b600060015460025410156106ee5760405162461bcd60e51b815260040161020f906108fa565b5060018054600090815260208190526040812080549190558154820190915590565b6040518060800160405280600081526020016000815260200160008152602001600081525090565b600060208284031215610749578081fd5b813561075481610a17565b9392505050565b60006020828403121561076c578081fd5b815161075481610a17565b600060208284031215610788578081fd5b81518015158114610754578182fd5b6000602082840312156107a8578081fd5b5035919050565b6000602082840312156107c0578081fd5b5051919050565b6000806000606084860312156107db578182fd5b505081359360208301359350604090920135919050565b93845260208401929092526040830152606082015260800190565b6001600160a01b0391909116815260200190565b6001600160a01b0392831681529116602082015260400190565b6001600160a01b039384168152919092166020820152604081019190915260600190565b90815260200190565b918252602082015260400190565b6020808252601c908201527f746f6b656e20616c6c6f77616e6365206e6f7420617070726f76656400000000604082015260600190565b6020808252602d908201527f4465706f7369744d616e616765723a2073656e646572206973206e6f7420526f60408201526c1b1b1d5c0818dbdb9d1c9858dd609a1b606082015260800190565b60208082526027908201527f4465706f73697420436f72653a2051756575652073686f756c64206265206e6f6040820152666e2d656d70747960c81b606082015260800190565b60208082526024908201527f746f6b656e206465706f736974206d75737420626520677265617465722074686040820152630616e20360e41b606082015260800190565b6020808252601b908201527f746f6b656e207472616e73666572206e6f7420617070726f7665640000000000604082015260600190565b600083825260206040818401528351806040850152825b818110156109ef578581018301518582016060015282016109d3565b81811115610a005783606083870101525b50601f01601f191692909201606001949350505050565b6001600160a01b03811681146104e957600080fdfea164736f6c634300060c000a"

// DeployDepositManager deploys a new Ethereum contract, binding an instance of DepositManager to it.
func DeployDepositManager(auth *bind.TransactOpts, backend bind.ContractBackend, _tokenRegistry common.Address, _vault common.Address, maxSubtreeDepth *big.Int) (common.Address, *types.Transaction, *DepositManager, error) {
	parsed, err := abi.JSON(strings.NewReader(DepositManagerABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(DepositManagerBin), backend, _tokenRegistry, _vault, maxSubtreeDepth)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DepositManager{DepositManagerCaller: DepositManagerCaller{contract: contract}, DepositManagerTransactor: DepositManagerTransactor{contract: contract}, DepositManagerFilterer: DepositManagerFilterer{contract: contract}}, nil
}

// DepositManager is an auto generated Go binding around an Ethereum contract.
type DepositManager struct {
	DepositManagerCaller     // Read-only binding to the contract
	DepositManagerTransactor // Write-only binding to the contract
	DepositManagerFilterer   // Log filterer for contract events
}

// DepositManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type DepositManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DepositManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DepositManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DepositManagerSession struct {
	Contract     *DepositManager   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DepositManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DepositManagerCallerSession struct {
	Contract *DepositManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// DepositManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DepositManagerTransactorSession struct {
	Contract     *DepositManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// DepositManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type DepositManagerRaw struct {
	Contract *DepositManager // Generic contract binding to access the raw methods on
}

// DepositManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DepositManagerCallerRaw struct {
	Contract *DepositManagerCaller // Generic read-only contract binding to access the raw methods on
}

// DepositManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DepositManagerTransactorRaw struct {
	Contract *DepositManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDepositManager creates a new instance of DepositManager, bound to a specific deployed contract.
func NewDepositManager(address common.Address, backend bind.ContractBackend) (*DepositManager, error) {
	contract, err := bindDepositManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DepositManager{DepositManagerCaller: DepositManagerCaller{contract: contract}, DepositManagerTransactor: DepositManagerTransactor{contract: contract}, DepositManagerFilterer: DepositManagerFilterer{contract: contract}}, nil
}

// NewDepositManagerCaller creates a new read-only instance of DepositManager, bound to a specific deployed contract.
func NewDepositManagerCaller(address common.Address, caller bind.ContractCaller) (*DepositManagerCaller, error) {
	contract, err := bindDepositManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DepositManagerCaller{contract: contract}, nil
}

// NewDepositManagerTransactor creates a new write-only instance of DepositManager, bound to a specific deployed contract.
func NewDepositManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*DepositManagerTransactor, error) {
	contract, err := bindDepositManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DepositManagerTransactor{contract: contract}, nil
}

// NewDepositManagerFilterer creates a new log filterer instance of DepositManager, bound to a specific deployed contract.
func NewDepositManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*DepositManagerFilterer, error) {
	contract, err := bindDepositManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DepositManagerFilterer{contract: contract}, nil
}

// bindDepositManager binds a generic wrapper to an already deployed contract.
func bindDepositManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DepositManagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DepositManager *DepositManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DepositManager.Contract.DepositManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DepositManager *DepositManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositManager.Contract.DepositManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DepositManager *DepositManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DepositManager.Contract.DepositManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DepositManager *DepositManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DepositManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DepositManager *DepositManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DepositManager *DepositManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DepositManager.Contract.contract.Transact(opts, method, params...)
}

// BabyTrees is a free data retrieval call binding the contract method 0xee9d68ce.
//
// Solidity: function babyTrees(uint256 ) view returns(bytes32)
func (_DepositManager *DepositManagerCaller) BabyTrees(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _DepositManager.contract.Call(opts, &out, "babyTrees", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BabyTrees is a free data retrieval call binding the contract method 0xee9d68ce.
//
// Solidity: function babyTrees(uint256 ) view returns(bytes32)
func (_DepositManager *DepositManagerSession) BabyTrees(arg0 *big.Int) ([32]byte, error) {
	return _DepositManager.Contract.BabyTrees(&_DepositManager.CallOpts, arg0)
}

// BabyTrees is a free data retrieval call binding the contract method 0xee9d68ce.
//
// Solidity: function babyTrees(uint256 ) view returns(bytes32)
func (_DepositManager *DepositManagerCallerSession) BabyTrees(arg0 *big.Int) ([32]byte, error) {
	return _DepositManager.Contract.BabyTrees(&_DepositManager.CallOpts, arg0)
}

// BabyTreesLength is a free data retrieval call binding the contract method 0x425e97f2.
//
// Solidity: function babyTreesLength() view returns(uint256)
func (_DepositManager *DepositManagerCaller) BabyTreesLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DepositManager.contract.Call(opts, &out, "babyTreesLength")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BabyTreesLength is a free data retrieval call binding the contract method 0x425e97f2.
//
// Solidity: function babyTreesLength() view returns(uint256)
func (_DepositManager *DepositManagerSession) BabyTreesLength() (*big.Int, error) {
	return _DepositManager.Contract.BabyTreesLength(&_DepositManager.CallOpts)
}

// BabyTreesLength is a free data retrieval call binding the contract method 0x425e97f2.
//
// Solidity: function babyTreesLength() view returns(uint256)
func (_DepositManager *DepositManagerCallerSession) BabyTreesLength() (*big.Int, error) {
	return _DepositManager.Contract.BabyTreesLength(&_DepositManager.CallOpts)
}

// Back is a free data retrieval call binding the contract method 0x8dde0840.
//
// Solidity: function back() view returns(uint256)
func (_DepositManager *DepositManagerCaller) Back(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DepositManager.contract.Call(opts, &out, "back")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Back is a free data retrieval call binding the contract method 0x8dde0840.
//
// Solidity: function back() view returns(uint256)
func (_DepositManager *DepositManagerSession) Back() (*big.Int, error) {
	return _DepositManager.Contract.Back(&_DepositManager.CallOpts)
}

// Back is a free data retrieval call binding the contract method 0x8dde0840.
//
// Solidity: function back() view returns(uint256)
func (_DepositManager *DepositManagerCallerSession) Back() (*big.Int, error) {
	return _DepositManager.Contract.Back(&_DepositManager.CallOpts)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_DepositManager *DepositManagerCaller) DepositCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DepositManager.contract.Call(opts, &out, "depositCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_DepositManager *DepositManagerSession) DepositCount() (*big.Int, error) {
	return _DepositManager.Contract.DepositCount(&_DepositManager.CallOpts)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_DepositManager *DepositManagerCallerSession) DepositCount() (*big.Int, error) {
	return _DepositManager.Contract.DepositCount(&_DepositManager.CallOpts)
}

// Front is a free data retrieval call binding the contract method 0xba75bbd8.
//
// Solidity: function front() view returns(uint256)
func (_DepositManager *DepositManagerCaller) Front(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DepositManager.contract.Call(opts, &out, "front")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Front is a free data retrieval call binding the contract method 0xba75bbd8.
//
// Solidity: function front() view returns(uint256)
func (_DepositManager *DepositManagerSession) Front() (*big.Int, error) {
	return _DepositManager.Contract.Front(&_DepositManager.CallOpts)
}

// Front is a free data retrieval call binding the contract method 0xba75bbd8.
//
// Solidity: function front() view returns(uint256)
func (_DepositManager *DepositManagerCallerSession) Front() (*big.Int, error) {
	return _DepositManager.Contract.Front(&_DepositManager.CallOpts)
}

// ParamMaxSubtreeSize is a free data retrieval call binding the contract method 0xc7accaa4.
//
// Solidity: function paramMaxSubtreeSize() view returns(uint256)
func (_DepositManager *DepositManagerCaller) ParamMaxSubtreeSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DepositManager.contract.Call(opts, &out, "paramMaxSubtreeSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ParamMaxSubtreeSize is a free data retrieval call binding the contract method 0xc7accaa4.
//
// Solidity: function paramMaxSubtreeSize() view returns(uint256)
func (_DepositManager *DepositManagerSession) ParamMaxSubtreeSize() (*big.Int, error) {
	return _DepositManager.Contract.ParamMaxSubtreeSize(&_DepositManager.CallOpts)
}

// ParamMaxSubtreeSize is a free data retrieval call binding the contract method 0xc7accaa4.
//
// Solidity: function paramMaxSubtreeSize() view returns(uint256)
func (_DepositManager *DepositManagerCallerSession) ParamMaxSubtreeSize() (*big.Int, error) {
	return _DepositManager.Contract.ParamMaxSubtreeSize(&_DepositManager.CallOpts)
}

// Queue is a free data retrieval call binding the contract method 0xddf0b009.
//
// Solidity: function queue(uint256 ) view returns(bytes32)
func (_DepositManager *DepositManagerCaller) Queue(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _DepositManager.contract.Call(opts, &out, "queue", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Queue is a free data retrieval call binding the contract method 0xddf0b009.
//
// Solidity: function queue(uint256 ) view returns(bytes32)
func (_DepositManager *DepositManagerSession) Queue(arg0 *big.Int) ([32]byte, error) {
	return _DepositManager.Contract.Queue(&_DepositManager.CallOpts, arg0)
}

// Queue is a free data retrieval call binding the contract method 0xddf0b009.
//
// Solidity: function queue(uint256 ) view returns(bytes32)
func (_DepositManager *DepositManagerCallerSession) Queue(arg0 *big.Int) ([32]byte, error) {
	return _DepositManager.Contract.Queue(&_DepositManager.CallOpts, arg0)
}

// Rollup is a free data retrieval call binding the contract method 0xcb23bcb5.
//
// Solidity: function rollup() view returns(address)
func (_DepositManager *DepositManagerCaller) Rollup(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DepositManager.contract.Call(opts, &out, "rollup")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Rollup is a free data retrieval call binding the contract method 0xcb23bcb5.
//
// Solidity: function rollup() view returns(address)
func (_DepositManager *DepositManagerSession) Rollup() (common.Address, error) {
	return _DepositManager.Contract.Rollup(&_DepositManager.CallOpts)
}

// Rollup is a free data retrieval call binding the contract method 0xcb23bcb5.
//
// Solidity: function rollup() view returns(address)
func (_DepositManager *DepositManagerCallerSession) Rollup() (common.Address, error) {
	return _DepositManager.Contract.Rollup(&_DepositManager.CallOpts)
}

// TokenRegistry is a free data retrieval call binding the contract method 0x9d23c4c7.
//
// Solidity: function tokenRegistry() view returns(address)
func (_DepositManager *DepositManagerCaller) TokenRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DepositManager.contract.Call(opts, &out, "tokenRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenRegistry is a free data retrieval call binding the contract method 0x9d23c4c7.
//
// Solidity: function tokenRegistry() view returns(address)
func (_DepositManager *DepositManagerSession) TokenRegistry() (common.Address, error) {
	return _DepositManager.Contract.TokenRegistry(&_DepositManager.CallOpts)
}

// TokenRegistry is a free data retrieval call binding the contract method 0x9d23c4c7.
//
// Solidity: function tokenRegistry() view returns(address)
func (_DepositManager *DepositManagerCallerSession) TokenRegistry() (common.Address, error) {
	return _DepositManager.Contract.TokenRegistry(&_DepositManager.CallOpts)
}

// Vault is a free data retrieval call binding the contract method 0xfbfa77cf.
//
// Solidity: function vault() view returns(address)
func (_DepositManager *DepositManagerCaller) Vault(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DepositManager.contract.Call(opts, &out, "vault")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Vault is a free data retrieval call binding the contract method 0xfbfa77cf.
//
// Solidity: function vault() view returns(address)
func (_DepositManager *DepositManagerSession) Vault() (common.Address, error) {
	return _DepositManager.Contract.Vault(&_DepositManager.CallOpts)
}

// Vault is a free data retrieval call binding the contract method 0xfbfa77cf.
//
// Solidity: function vault() view returns(address)
func (_DepositManager *DepositManagerCallerSession) Vault() (common.Address, error) {
	return _DepositManager.Contract.Vault(&_DepositManager.CallOpts)
}

// DepositFor is a paid mutator transaction binding the contract method 0x966fda62.
//
// Solidity: function depositFor(uint256 pubkeyID, uint256 amount, uint256 tokenID) returns()
func (_DepositManager *DepositManagerTransactor) DepositFor(opts *bind.TransactOpts, pubkeyID *big.Int, amount *big.Int, tokenID *big.Int) (*types.Transaction, error) {
	return _DepositManager.contract.Transact(opts, "depositFor", pubkeyID, amount, tokenID)
}

// DepositFor is a paid mutator transaction binding the contract method 0x966fda62.
//
// Solidity: function depositFor(uint256 pubkeyID, uint256 amount, uint256 tokenID) returns()
func (_DepositManager *DepositManagerSession) DepositFor(pubkeyID *big.Int, amount *big.Int, tokenID *big.Int) (*types.Transaction, error) {
	return _DepositManager.Contract.DepositFor(&_DepositManager.TransactOpts, pubkeyID, amount, tokenID)
}

// DepositFor is a paid mutator transaction binding the contract method 0x966fda62.
//
// Solidity: function depositFor(uint256 pubkeyID, uint256 amount, uint256 tokenID) returns()
func (_DepositManager *DepositManagerTransactorSession) DepositFor(pubkeyID *big.Int, amount *big.Int, tokenID *big.Int) (*types.Transaction, error) {
	return _DepositManager.Contract.DepositFor(&_DepositManager.TransactOpts, pubkeyID, amount, tokenID)
}

// DequeueToSubmit is a paid mutator transaction binding the contract method 0xd86ee48d.
//
// Solidity: function dequeueToSubmit() returns(bytes32 subtreeRoot)
func (_DepositManager *DepositManagerTransactor) DequeueToSubmit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositManager.contract.Transact(opts, "dequeueToSubmit")
}

// DequeueToSubmit is a paid mutator transaction binding the contract method 0xd86ee48d.
//
// Solidity: function dequeueToSubmit() returns(bytes32 subtreeRoot)
func (_DepositManager *DepositManagerSession) DequeueToSubmit() (*types.Transaction, error) {
	return _DepositManager.Contract.DequeueToSubmit(&_DepositManager.TransactOpts)
}

// DequeueToSubmit is a paid mutator transaction binding the contract method 0xd86ee48d.
//
// Solidity: function dequeueToSubmit() returns(bytes32 subtreeRoot)
func (_DepositManager *DepositManagerTransactorSession) DequeueToSubmit() (*types.Transaction, error) {
	return _DepositManager.Contract.DequeueToSubmit(&_DepositManager.TransactOpts)
}

// Reenqueue is a paid mutator transaction binding the contract method 0xade00026.
//
// Solidity: function reenqueue(bytes32 subtreeRoot) returns()
func (_DepositManager *DepositManagerTransactor) Reenqueue(opts *bind.TransactOpts, subtreeRoot [32]byte) (*types.Transaction, error) {
	return _DepositManager.contract.Transact(opts, "reenqueue", subtreeRoot)
}

// Reenqueue is a paid mutator transaction binding the contract method 0xade00026.
//
// Solidity: function reenqueue(bytes32 subtreeRoot) returns()
func (_DepositManager *DepositManagerSession) Reenqueue(subtreeRoot [32]byte) (*types.Transaction, error) {
	return _DepositManager.Contract.Reenqueue(&_DepositManager.TransactOpts, subtreeRoot)
}

// Reenqueue is a paid mutator transaction binding the contract method 0xade00026.
//
// Solidity: function reenqueue(bytes32 subtreeRoot) returns()
func (_DepositManager *DepositManagerTransactorSession) Reenqueue(subtreeRoot [32]byte) (*types.Transaction, error) {
	return _DepositManager.Contract.Reenqueue(&_DepositManager.TransactOpts, subtreeRoot)
}

// SetRollupAddress is a paid mutator transaction binding the contract method 0x07663706.
//
// Solidity: function setRollupAddress(address _rollup) returns()
func (_DepositManager *DepositManagerTransactor) SetRollupAddress(opts *bind.TransactOpts, _rollup common.Address) (*types.Transaction, error) {
	return _DepositManager.contract.Transact(opts, "setRollupAddress", _rollup)
}

// SetRollupAddress is a paid mutator transaction binding the contract method 0x07663706.
//
// Solidity: function setRollupAddress(address _rollup) returns()
func (_DepositManager *DepositManagerSession) SetRollupAddress(_rollup common.Address) (*types.Transaction, error) {
	return _DepositManager.Contract.SetRollupAddress(&_DepositManager.TransactOpts, _rollup)
}

// SetRollupAddress is a paid mutator transaction binding the contract method 0x07663706.
//
// Solidity: function setRollupAddress(address _rollup) returns()
func (_DepositManager *DepositManagerTransactorSession) SetRollupAddress(_rollup common.Address) (*types.Transaction, error) {
	return _DepositManager.Contract.SetRollupAddress(&_DepositManager.TransactOpts, _rollup)
}

// DepositManagerDepositQueuedIterator is returned from FilterDepositQueued and is used to iterate over the raw logs and unpacked data for DepositQueued events raised by the DepositManager contract.
type DepositManagerDepositQueuedIterator struct {
	Event *DepositManagerDepositQueued // Event containing the contract specifics and raw log

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
func (it *DepositManagerDepositQueuedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositManagerDepositQueued)
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
		it.Event = new(DepositManagerDepositQueued)
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
func (it *DepositManagerDepositQueuedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositManagerDepositQueuedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositManagerDepositQueued represents a DepositQueued event raised by the DepositManager contract.
type DepositManagerDepositQueued struct {
	PubkeyID *big.Int
	Data     []byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDepositQueued is a free log retrieval operation binding the contract event 0x5a1922090c12e28c132a961f6bb4d74350598a62e8981b5eb9bb5ccbbce9df94.
//
// Solidity: event DepositQueued(uint256 pubkeyID, bytes data)
func (_DepositManager *DepositManagerFilterer) FilterDepositQueued(opts *bind.FilterOpts) (*DepositManagerDepositQueuedIterator, error) {

	logs, sub, err := _DepositManager.contract.FilterLogs(opts, "DepositQueued")
	if err != nil {
		return nil, err
	}
	return &DepositManagerDepositQueuedIterator{contract: _DepositManager.contract, event: "DepositQueued", logs: logs, sub: sub}, nil
}

// WatchDepositQueued is a free log subscription operation binding the contract event 0x5a1922090c12e28c132a961f6bb4d74350598a62e8981b5eb9bb5ccbbce9df94.
//
// Solidity: event DepositQueued(uint256 pubkeyID, bytes data)
func (_DepositManager *DepositManagerFilterer) WatchDepositQueued(opts *bind.WatchOpts, sink chan<- *DepositManagerDepositQueued) (event.Subscription, error) {

	logs, sub, err := _DepositManager.contract.WatchLogs(opts, "DepositQueued")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositManagerDepositQueued)
				if err := _DepositManager.contract.UnpackLog(event, "DepositQueued", log); err != nil {
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

// ParseDepositQueued is a log parse operation binding the contract event 0x5a1922090c12e28c132a961f6bb4d74350598a62e8981b5eb9bb5ccbbce9df94.
//
// Solidity: event DepositQueued(uint256 pubkeyID, bytes data)
func (_DepositManager *DepositManagerFilterer) ParseDepositQueued(log types.Log) (*DepositManagerDepositQueued, error) {
	event := new(DepositManagerDepositQueued)
	if err := _DepositManager.contract.UnpackLog(event, "DepositQueued", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositManagerDepositSubTreeReadyIterator is returned from FilterDepositSubTreeReady and is used to iterate over the raw logs and unpacked data for DepositSubTreeReady events raised by the DepositManager contract.
type DepositManagerDepositSubTreeReadyIterator struct {
	Event *DepositManagerDepositSubTreeReady // Event containing the contract specifics and raw log

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
func (it *DepositManagerDepositSubTreeReadyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositManagerDepositSubTreeReady)
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
		it.Event = new(DepositManagerDepositSubTreeReady)
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
func (it *DepositManagerDepositSubTreeReadyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositManagerDepositSubTreeReadyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositManagerDepositSubTreeReady represents a DepositSubTreeReady event raised by the DepositManager contract.
type DepositManagerDepositSubTreeReady struct {
	Root [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterDepositSubTreeReady is a free log retrieval operation binding the contract event 0x4744f3a44c5716c9fa423a71cdaa806771a8bf469f4c007ca338b8e8e202a8b5.
//
// Solidity: event DepositSubTreeReady(bytes32 root)
func (_DepositManager *DepositManagerFilterer) FilterDepositSubTreeReady(opts *bind.FilterOpts) (*DepositManagerDepositSubTreeReadyIterator, error) {

	logs, sub, err := _DepositManager.contract.FilterLogs(opts, "DepositSubTreeReady")
	if err != nil {
		return nil, err
	}
	return &DepositManagerDepositSubTreeReadyIterator{contract: _DepositManager.contract, event: "DepositSubTreeReady", logs: logs, sub: sub}, nil
}

// WatchDepositSubTreeReady is a free log subscription operation binding the contract event 0x4744f3a44c5716c9fa423a71cdaa806771a8bf469f4c007ca338b8e8e202a8b5.
//
// Solidity: event DepositSubTreeReady(bytes32 root)
func (_DepositManager *DepositManagerFilterer) WatchDepositSubTreeReady(opts *bind.WatchOpts, sink chan<- *DepositManagerDepositSubTreeReady) (event.Subscription, error) {

	logs, sub, err := _DepositManager.contract.WatchLogs(opts, "DepositSubTreeReady")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositManagerDepositSubTreeReady)
				if err := _DepositManager.contract.UnpackLog(event, "DepositSubTreeReady", log); err != nil {
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

// ParseDepositSubTreeReady is a log parse operation binding the contract event 0x4744f3a44c5716c9fa423a71cdaa806771a8bf469f4c007ca338b8e8e202a8b5.
//
// Solidity: event DepositSubTreeReady(bytes32 root)
func (_DepositManager *DepositManagerFilterer) ParseDepositSubTreeReady(log types.Log) (*DepositManagerDepositSubTreeReady, error) {
	event := new(DepositManagerDepositSubTreeReady)
	if err := _DepositManager.contract.UnpackLog(event, "DepositSubTreeReady", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
