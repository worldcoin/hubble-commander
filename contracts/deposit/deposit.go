// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package deposit

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

// DepositABI is the input ABI used to generate the binding from.
const DepositABI = "[{\"inputs\":[{\"internalType\":\"contractITokenRegistry\",\"name\":\"_tokenRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_vault\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"maxSubtreeDepth\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"DepositQueued\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"DepositSubTreeReady\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"babyTrees\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"babyTreesLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"back\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"}],\"name\":\"depositFor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dequeueToSubmit\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"subtreeRoot\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"front\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paramMaxSubtreeSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"queue\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"subtreeRoot\",\"type\":\"bytes32\"}],\"name\":\"reenqueue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollup\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_rollup\",\"type\":\"address\"}],\"name\":\"setRollupAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenRegistry\",\"outputs\":[{\"internalType\":\"contractITokenRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"vault\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// DepositBin is the compiled bytecode used for deploying new contracts.
var DepositBin = "0x608060405260018055600060025560006004556000600555600260065534801561002857600080fd5b50604051610b25380380610b2583398101604081905261004791610083565b600980546001600160a01b039485166001600160a01b03199182161790915560078054939094169216919091179091556001901b6006556100dd565b600080600060608486031215610097578283fd5b83516100a2816100c5565b60208501519093506100b3816100c5565b80925050604084015190509250925092565b6001600160a01b03811681146100da57600080fd5b50565b610a39806100ec6000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c8063ba75bbd81161008c578063d86ee48d11610066578063d86ee48d14610185578063ddf0b0091461018d578063ee9d68ce146101a0578063fbfa77cf146101b3576100ea565b8063ba75bbd81461016d578063c7accaa414610175578063cb23bcb51461017d576100ea565b80638dde0840116100c85780638dde08401461012a578063966fda62146101325780639d23c4c714610145578063ade000261461015a576100ea565b806307663706146100ef5780632dfdf0b514610104578063425e97f214610122575b600080fd5b6101026100fd366004610738565b6101bb565b005b61010c6101dd565b604051610119919061085f565b60405180910390f35b61010c6101e3565b61010c6101e9565b6101026101403660046107c7565b6101ef565b61014d6104a7565b604051610119919061080d565b610102610168366004610797565b6104b6565b61010c6104ec565b61010c6104f2565b61014d6104f8565b61010c610507565b61010c61019b366004610797565b610541565b61010c6101ae366004610797565b610553565b61014d610565565b600880546001600160a01b0319166001600160a01b0392909216919091179055565b60055481565b60045481565b60025481565b600082116102185760405162461bcd60e51b815260040161020f90610941565b60405180910390fd5b6009546040516320f5ab4f60e11b81526000916001600160a01b0316906341eb569e9061024990859060040161085f565b60206040518083038186803b15801561026157600080fd5b505afa158015610275573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610299919061075b565b905082816001600160a01b031663dd62ed3e33306040518363ffffffff1660e01b81526004016102ca929190610821565b60206040518083038186803b1580156102e257600080fd5b505afa1580156102f6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061031a91906107af565b10156103385760405162461bcd60e51b815260040161020f90610876565b6007546040516323b872dd60e01b81526001600160a01b03808416926323b872dd9261036c9233921690889060040161083b565b602060405180830381600087803b15801561038657600080fd5b505af115801561039a573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103be9190610777565b6103da5760405162461bcd60e51b815260040161020f90610985565b6103e2610710565b604051806080016040528086815260200184815260200185815260200160008152509050606061041182610574565b90507f5a1922090c12e28c132a961f6bb4d74350598a62e8981b5eb9bb5ccbbce9df9486826040516104449291906109bc565b60405180910390a1600061045e82805190602001206105b3565b9050801561049e577f4744f3a44c5716c9fa423a71cdaa806771a8bf469f4c007ca338b8e8e202a8b581604051610495919061085f565b60405180910390a15b50505050505050565b6009546001600160a01b031681565b6008546001600160a01b031633146104e05760405162461bcd60e51b815260040161020f906108ad565b6104e9816106ac565b50565b60015481565b60065481565b6008546001600160a01b031681565b6008546000906001600160a01b031633146105345760405162461bcd60e51b815260040161020f906108ad565b61053c6106c8565b905090565b60006020819052908152604090205481565b60036020526000908152604090205481565b6007546001600160a01b031681565b6060816000015182602001518360400151846060015160405160200161059d94939291906107f2565b6040516020818303038152906040529050919050565b6005805460019081019182905560045460008181526003602052604081208590559291015b6001821661064a5760011981016000908152600360209081526040808320546000198501845292819020549051610610939201610868565b60408051601f19818403018152918152815160209283012060011984016000908152600390935291205560019190911c90600019016105d8565b600481905560065460055414156106a0576000805260036020527f3617319a054d772f909f7c479a2cebe5066e836a939412e32403c99029b92eff549250610691836106ac565b600060048190556005556106a5565b600092505b5050919050565b6002805460010190819055600090815260208190526040902055565b600060015460025410156106ee5760405162461bcd60e51b815260040161020f906108fa565b5060018054600090815260208190526040812080549190558154820190915590565b6040518060800160405280600081526020016000815260200160008152602001600081525090565b600060208284031215610749578081fd5b813561075481610a17565b9392505050565b60006020828403121561076c578081fd5b815161075481610a17565b600060208284031215610788578081fd5b81518015158114610754578182fd5b6000602082840312156107a8578081fd5b5035919050565b6000602082840312156107c0578081fd5b5051919050565b6000806000606084860312156107db578182fd5b505081359360208301359350604090920135919050565b93845260208401929092526040830152606082015260800190565b6001600160a01b0391909116815260200190565b6001600160a01b0392831681529116602082015260400190565b6001600160a01b039384168152919092166020820152604081019190915260600190565b90815260200190565b918252602082015260400190565b6020808252601c908201527f746f6b656e20616c6c6f77616e6365206e6f7420617070726f76656400000000604082015260600190565b6020808252602d908201527f4465706f7369744d616e616765723a2073656e646572206973206e6f7420526f60408201526c1b1b1d5c0818dbdb9d1c9858dd609a1b606082015260800190565b60208082526027908201527f4465706f73697420436f72653a2051756575652073686f756c64206265206e6f6040820152666e2d656d70747960c81b606082015260800190565b60208082526024908201527f746f6b656e206465706f736974206d75737420626520677265617465722074686040820152630616e20360e41b606082015260800190565b6020808252601b908201527f746f6b656e207472616e73666572206e6f7420617070726f7665640000000000604082015260600190565b600083825260206040818401528351806040850152825b818110156109ef578581018301518582016060015282016109d3565b81811115610a005783606083870101525b50601f01601f191692909201606001949350505050565b6001600160a01b03811681146104e957600080fdfea164736f6c634300060c000a"

// DeployDeposit deploys a new Ethereum contract, binding an instance of Deposit to it.
func DeployDeposit(auth *bind.TransactOpts, backend bind.ContractBackend, _tokenRegistry common.Address, _vault common.Address, maxSubtreeDepth *big.Int) (common.Address, *types.Transaction, *Deposit, error) {
	parsed, err := abi.JSON(strings.NewReader(DepositABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(DepositBin), backend, _tokenRegistry, _vault, maxSubtreeDepth)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Deposit{DepositCaller: DepositCaller{contract: contract}, DepositTransactor: DepositTransactor{contract: contract}, DepositFilterer: DepositFilterer{contract: contract}}, nil
}

// Deposit is an auto generated Go binding around an Ethereum contract.
type Deposit struct {
	DepositCaller     // Read-only binding to the contract
	DepositTransactor // Write-only binding to the contract
	DepositFilterer   // Log filterer for contract events
}

// DepositCaller is an auto generated read-only Go binding around an Ethereum contract.
type DepositCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DepositTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DepositFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DepositSession struct {
	Contract     *Deposit          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DepositCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DepositCallerSession struct {
	Contract *DepositCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// DepositTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DepositTransactorSession struct {
	Contract     *DepositTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// DepositRaw is an auto generated low-level Go binding around an Ethereum contract.
type DepositRaw struct {
	Contract *Deposit // Generic contract binding to access the raw methods on
}

// DepositCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DepositCallerRaw struct {
	Contract *DepositCaller // Generic read-only contract binding to access the raw methods on
}

// DepositTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DepositTransactorRaw struct {
	Contract *DepositTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDeposit creates a new instance of Deposit, bound to a specific deployed contract.
func NewDeposit(address common.Address, backend bind.ContractBackend) (*Deposit, error) {
	contract, err := bindDeposit(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Deposit{DepositCaller: DepositCaller{contract: contract}, DepositTransactor: DepositTransactor{contract: contract}, DepositFilterer: DepositFilterer{contract: contract}}, nil
}

// NewDepositCaller creates a new read-only instance of Deposit, bound to a specific deployed contract.
func NewDepositCaller(address common.Address, caller bind.ContractCaller) (*DepositCaller, error) {
	contract, err := bindDeposit(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DepositCaller{contract: contract}, nil
}

// NewDepositTransactor creates a new write-only instance of Deposit, bound to a specific deployed contract.
func NewDepositTransactor(address common.Address, transactor bind.ContractTransactor) (*DepositTransactor, error) {
	contract, err := bindDeposit(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DepositTransactor{contract: contract}, nil
}

// NewDepositFilterer creates a new log filterer instance of Deposit, bound to a specific deployed contract.
func NewDepositFilterer(address common.Address, filterer bind.ContractFilterer) (*DepositFilterer, error) {
	contract, err := bindDeposit(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DepositFilterer{contract: contract}, nil
}

// bindDeposit binds a generic wrapper to an already deployed contract.
func bindDeposit(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DepositABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Deposit *DepositRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Deposit.Contract.DepositCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Deposit *DepositRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.Contract.DepositTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Deposit *DepositRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Deposit.Contract.DepositTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Deposit *DepositCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Deposit.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Deposit *DepositTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Deposit *DepositTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Deposit.Contract.contract.Transact(opts, method, params...)
}

// BabyTrees is a free data retrieval call binding the contract method 0xee9d68ce.
//
// Solidity: function babyTrees(uint256 ) view returns(bytes32)
func (_Deposit *DepositCaller) BabyTrees(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "babyTrees", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BabyTrees is a free data retrieval call binding the contract method 0xee9d68ce.
//
// Solidity: function babyTrees(uint256 ) view returns(bytes32)
func (_Deposit *DepositSession) BabyTrees(arg0 *big.Int) ([32]byte, error) {
	return _Deposit.Contract.BabyTrees(&_Deposit.CallOpts, arg0)
}

// BabyTrees is a free data retrieval call binding the contract method 0xee9d68ce.
//
// Solidity: function babyTrees(uint256 ) view returns(bytes32)
func (_Deposit *DepositCallerSession) BabyTrees(arg0 *big.Int) ([32]byte, error) {
	return _Deposit.Contract.BabyTrees(&_Deposit.CallOpts, arg0)
}

// BabyTreesLength is a free data retrieval call binding the contract method 0x425e97f2.
//
// Solidity: function babyTreesLength() view returns(uint256)
func (_Deposit *DepositCaller) BabyTreesLength(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "babyTreesLength")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BabyTreesLength is a free data retrieval call binding the contract method 0x425e97f2.
//
// Solidity: function babyTreesLength() view returns(uint256)
func (_Deposit *DepositSession) BabyTreesLength() (*big.Int, error) {
	return _Deposit.Contract.BabyTreesLength(&_Deposit.CallOpts)
}

// BabyTreesLength is a free data retrieval call binding the contract method 0x425e97f2.
//
// Solidity: function babyTreesLength() view returns(uint256)
func (_Deposit *DepositCallerSession) BabyTreesLength() (*big.Int, error) {
	return _Deposit.Contract.BabyTreesLength(&_Deposit.CallOpts)
}

// Back is a free data retrieval call binding the contract method 0x8dde0840.
//
// Solidity: function back() view returns(uint256)
func (_Deposit *DepositCaller) Back(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "back")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Back is a free data retrieval call binding the contract method 0x8dde0840.
//
// Solidity: function back() view returns(uint256)
func (_Deposit *DepositSession) Back() (*big.Int, error) {
	return _Deposit.Contract.Back(&_Deposit.CallOpts)
}

// Back is a free data retrieval call binding the contract method 0x8dde0840.
//
// Solidity: function back() view returns(uint256)
func (_Deposit *DepositCallerSession) Back() (*big.Int, error) {
	return _Deposit.Contract.Back(&_Deposit.CallOpts)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_Deposit *DepositCaller) DepositCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "depositCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_Deposit *DepositSession) DepositCount() (*big.Int, error) {
	return _Deposit.Contract.DepositCount(&_Deposit.CallOpts)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_Deposit *DepositCallerSession) DepositCount() (*big.Int, error) {
	return _Deposit.Contract.DepositCount(&_Deposit.CallOpts)
}

// Front is a free data retrieval call binding the contract method 0xba75bbd8.
//
// Solidity: function front() view returns(uint256)
func (_Deposit *DepositCaller) Front(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "front")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Front is a free data retrieval call binding the contract method 0xba75bbd8.
//
// Solidity: function front() view returns(uint256)
func (_Deposit *DepositSession) Front() (*big.Int, error) {
	return _Deposit.Contract.Front(&_Deposit.CallOpts)
}

// Front is a free data retrieval call binding the contract method 0xba75bbd8.
//
// Solidity: function front() view returns(uint256)
func (_Deposit *DepositCallerSession) Front() (*big.Int, error) {
	return _Deposit.Contract.Front(&_Deposit.CallOpts)
}

// ParamMaxSubtreeSize is a free data retrieval call binding the contract method 0xc7accaa4.
//
// Solidity: function paramMaxSubtreeSize() view returns(uint256)
func (_Deposit *DepositCaller) ParamMaxSubtreeSize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "paramMaxSubtreeSize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ParamMaxSubtreeSize is a free data retrieval call binding the contract method 0xc7accaa4.
//
// Solidity: function paramMaxSubtreeSize() view returns(uint256)
func (_Deposit *DepositSession) ParamMaxSubtreeSize() (*big.Int, error) {
	return _Deposit.Contract.ParamMaxSubtreeSize(&_Deposit.CallOpts)
}

// ParamMaxSubtreeSize is a free data retrieval call binding the contract method 0xc7accaa4.
//
// Solidity: function paramMaxSubtreeSize() view returns(uint256)
func (_Deposit *DepositCallerSession) ParamMaxSubtreeSize() (*big.Int, error) {
	return _Deposit.Contract.ParamMaxSubtreeSize(&_Deposit.CallOpts)
}

// Queue is a free data retrieval call binding the contract method 0xddf0b009.
//
// Solidity: function queue(uint256 ) view returns(bytes32)
func (_Deposit *DepositCaller) Queue(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "queue", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Queue is a free data retrieval call binding the contract method 0xddf0b009.
//
// Solidity: function queue(uint256 ) view returns(bytes32)
func (_Deposit *DepositSession) Queue(arg0 *big.Int) ([32]byte, error) {
	return _Deposit.Contract.Queue(&_Deposit.CallOpts, arg0)
}

// Queue is a free data retrieval call binding the contract method 0xddf0b009.
//
// Solidity: function queue(uint256 ) view returns(bytes32)
func (_Deposit *DepositCallerSession) Queue(arg0 *big.Int) ([32]byte, error) {
	return _Deposit.Contract.Queue(&_Deposit.CallOpts, arg0)
}

// Rollup is a free data retrieval call binding the contract method 0xcb23bcb5.
//
// Solidity: function rollup() view returns(address)
func (_Deposit *DepositCaller) Rollup(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "rollup")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Rollup is a free data retrieval call binding the contract method 0xcb23bcb5.
//
// Solidity: function rollup() view returns(address)
func (_Deposit *DepositSession) Rollup() (common.Address, error) {
	return _Deposit.Contract.Rollup(&_Deposit.CallOpts)
}

// Rollup is a free data retrieval call binding the contract method 0xcb23bcb5.
//
// Solidity: function rollup() view returns(address)
func (_Deposit *DepositCallerSession) Rollup() (common.Address, error) {
	return _Deposit.Contract.Rollup(&_Deposit.CallOpts)
}

// TokenRegistry is a free data retrieval call binding the contract method 0x9d23c4c7.
//
// Solidity: function tokenRegistry() view returns(address)
func (_Deposit *DepositCaller) TokenRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "tokenRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenRegistry is a free data retrieval call binding the contract method 0x9d23c4c7.
//
// Solidity: function tokenRegistry() view returns(address)
func (_Deposit *DepositSession) TokenRegistry() (common.Address, error) {
	return _Deposit.Contract.TokenRegistry(&_Deposit.CallOpts)
}

// TokenRegistry is a free data retrieval call binding the contract method 0x9d23c4c7.
//
// Solidity: function tokenRegistry() view returns(address)
func (_Deposit *DepositCallerSession) TokenRegistry() (common.Address, error) {
	return _Deposit.Contract.TokenRegistry(&_Deposit.CallOpts)
}

// Vault is a free data retrieval call binding the contract method 0xfbfa77cf.
//
// Solidity: function vault() view returns(address)
func (_Deposit *DepositCaller) Vault(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "vault")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Vault is a free data retrieval call binding the contract method 0xfbfa77cf.
//
// Solidity: function vault() view returns(address)
func (_Deposit *DepositSession) Vault() (common.Address, error) {
	return _Deposit.Contract.Vault(&_Deposit.CallOpts)
}

// Vault is a free data retrieval call binding the contract method 0xfbfa77cf.
//
// Solidity: function vault() view returns(address)
func (_Deposit *DepositCallerSession) Vault() (common.Address, error) {
	return _Deposit.Contract.Vault(&_Deposit.CallOpts)
}

// DepositFor is a paid mutator transaction binding the contract method 0x966fda62.
//
// Solidity: function depositFor(uint256 pubkeyID, uint256 amount, uint256 tokenID) returns()
func (_Deposit *DepositTransactor) DepositFor(opts *bind.TransactOpts, pubkeyID *big.Int, amount *big.Int, tokenID *big.Int) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "depositFor", pubkeyID, amount, tokenID)
}

// DepositFor is a paid mutator transaction binding the contract method 0x966fda62.
//
// Solidity: function depositFor(uint256 pubkeyID, uint256 amount, uint256 tokenID) returns()
func (_Deposit *DepositSession) DepositFor(pubkeyID *big.Int, amount *big.Int, tokenID *big.Int) (*types.Transaction, error) {
	return _Deposit.Contract.DepositFor(&_Deposit.TransactOpts, pubkeyID, amount, tokenID)
}

// DepositFor is a paid mutator transaction binding the contract method 0x966fda62.
//
// Solidity: function depositFor(uint256 pubkeyID, uint256 amount, uint256 tokenID) returns()
func (_Deposit *DepositTransactorSession) DepositFor(pubkeyID *big.Int, amount *big.Int, tokenID *big.Int) (*types.Transaction, error) {
	return _Deposit.Contract.DepositFor(&_Deposit.TransactOpts, pubkeyID, amount, tokenID)
}

// DequeueToSubmit is a paid mutator transaction binding the contract method 0xd86ee48d.
//
// Solidity: function dequeueToSubmit() returns(bytes32 subtreeRoot)
func (_Deposit *DepositTransactor) DequeueToSubmit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "dequeueToSubmit")
}

// DequeueToSubmit is a paid mutator transaction binding the contract method 0xd86ee48d.
//
// Solidity: function dequeueToSubmit() returns(bytes32 subtreeRoot)
func (_Deposit *DepositSession) DequeueToSubmit() (*types.Transaction, error) {
	return _Deposit.Contract.DequeueToSubmit(&_Deposit.TransactOpts)
}

// DequeueToSubmit is a paid mutator transaction binding the contract method 0xd86ee48d.
//
// Solidity: function dequeueToSubmit() returns(bytes32 subtreeRoot)
func (_Deposit *DepositTransactorSession) DequeueToSubmit() (*types.Transaction, error) {
	return _Deposit.Contract.DequeueToSubmit(&_Deposit.TransactOpts)
}

// Reenqueue is a paid mutator transaction binding the contract method 0xade00026.
//
// Solidity: function reenqueue(bytes32 subtreeRoot) returns()
func (_Deposit *DepositTransactor) Reenqueue(opts *bind.TransactOpts, subtreeRoot [32]byte) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "reenqueue", subtreeRoot)
}

// Reenqueue is a paid mutator transaction binding the contract method 0xade00026.
//
// Solidity: function reenqueue(bytes32 subtreeRoot) returns()
func (_Deposit *DepositSession) Reenqueue(subtreeRoot [32]byte) (*types.Transaction, error) {
	return _Deposit.Contract.Reenqueue(&_Deposit.TransactOpts, subtreeRoot)
}

// Reenqueue is a paid mutator transaction binding the contract method 0xade00026.
//
// Solidity: function reenqueue(bytes32 subtreeRoot) returns()
func (_Deposit *DepositTransactorSession) Reenqueue(subtreeRoot [32]byte) (*types.Transaction, error) {
	return _Deposit.Contract.Reenqueue(&_Deposit.TransactOpts, subtreeRoot)
}

// SetRollupAddress is a paid mutator transaction binding the contract method 0x07663706.
//
// Solidity: function setRollupAddress(address _rollup) returns()
func (_Deposit *DepositTransactor) SetRollupAddress(opts *bind.TransactOpts, _rollup common.Address) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "setRollupAddress", _rollup)
}

// SetRollupAddress is a paid mutator transaction binding the contract method 0x07663706.
//
// Solidity: function setRollupAddress(address _rollup) returns()
func (_Deposit *DepositSession) SetRollupAddress(_rollup common.Address) (*types.Transaction, error) {
	return _Deposit.Contract.SetRollupAddress(&_Deposit.TransactOpts, _rollup)
}

// SetRollupAddress is a paid mutator transaction binding the contract method 0x07663706.
//
// Solidity: function setRollupAddress(address _rollup) returns()
func (_Deposit *DepositTransactorSession) SetRollupAddress(_rollup common.Address) (*types.Transaction, error) {
	return _Deposit.Contract.SetRollupAddress(&_Deposit.TransactOpts, _rollup)
}

// DepositDepositQueuedIterator is returned from FilterDepositQueued and is used to iterate over the raw logs and unpacked data for DepositQueued events raised by the Deposit contract.
type DepositDepositQueuedIterator struct {
	Event *DepositDepositQueued // Event containing the contract specifics and raw log

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
func (it *DepositDepositQueuedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositDepositQueued)
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
		it.Event = new(DepositDepositQueued)
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
func (it *DepositDepositQueuedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositDepositQueuedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositDepositQueued represents a DepositQueued event raised by the Deposit contract.
type DepositDepositQueued struct {
	PubkeyID *big.Int
	Data     []byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDepositQueued is a free log retrieval operation binding the contract event 0x5a1922090c12e28c132a961f6bb4d74350598a62e8981b5eb9bb5ccbbce9df94.
//
// Solidity: event DepositQueued(uint256 pubkeyID, bytes data)
func (_Deposit *DepositFilterer) FilterDepositQueued(opts *bind.FilterOpts) (*DepositDepositQueuedIterator, error) {

	logs, sub, err := _Deposit.contract.FilterLogs(opts, "DepositQueued")
	if err != nil {
		return nil, err
	}
	return &DepositDepositQueuedIterator{contract: _Deposit.contract, event: "DepositQueued", logs: logs, sub: sub}, nil
}

// WatchDepositQueued is a free log subscription operation binding the contract event 0x5a1922090c12e28c132a961f6bb4d74350598a62e8981b5eb9bb5ccbbce9df94.
//
// Solidity: event DepositQueued(uint256 pubkeyID, bytes data)
func (_Deposit *DepositFilterer) WatchDepositQueued(opts *bind.WatchOpts, sink chan<- *DepositDepositQueued) (event.Subscription, error) {

	logs, sub, err := _Deposit.contract.WatchLogs(opts, "DepositQueued")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositDepositQueued)
				if err := _Deposit.contract.UnpackLog(event, "DepositQueued", log); err != nil {
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
func (_Deposit *DepositFilterer) ParseDepositQueued(log types.Log) (*DepositDepositQueued, error) {
	event := new(DepositDepositQueued)
	if err := _Deposit.contract.UnpackLog(event, "DepositQueued", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositDepositSubTreeReadyIterator is returned from FilterDepositSubTreeReady and is used to iterate over the raw logs and unpacked data for DepositSubTreeReady events raised by the Deposit contract.
type DepositDepositSubTreeReadyIterator struct {
	Event *DepositDepositSubTreeReady // Event containing the contract specifics and raw log

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
func (it *DepositDepositSubTreeReadyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositDepositSubTreeReady)
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
		it.Event = new(DepositDepositSubTreeReady)
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
func (it *DepositDepositSubTreeReadyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositDepositSubTreeReadyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositDepositSubTreeReady represents a DepositSubTreeReady event raised by the Deposit contract.
type DepositDepositSubTreeReady struct {
	Root [32]byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterDepositSubTreeReady is a free log retrieval operation binding the contract event 0x4744f3a44c5716c9fa423a71cdaa806771a8bf469f4c007ca338b8e8e202a8b5.
//
// Solidity: event DepositSubTreeReady(bytes32 root)
func (_Deposit *DepositFilterer) FilterDepositSubTreeReady(opts *bind.FilterOpts) (*DepositDepositSubTreeReadyIterator, error) {

	logs, sub, err := _Deposit.contract.FilterLogs(opts, "DepositSubTreeReady")
	if err != nil {
		return nil, err
	}
	return &DepositDepositSubTreeReadyIterator{contract: _Deposit.contract, event: "DepositSubTreeReady", logs: logs, sub: sub}, nil
}

// WatchDepositSubTreeReady is a free log subscription operation binding the contract event 0x4744f3a44c5716c9fa423a71cdaa806771a8bf469f4c007ca338b8e8e202a8b5.
//
// Solidity: event DepositSubTreeReady(bytes32 root)
func (_Deposit *DepositFilterer) WatchDepositSubTreeReady(opts *bind.WatchOpts, sink chan<- *DepositDepositSubTreeReady) (event.Subscription, error) {

	logs, sub, err := _Deposit.contract.WatchLogs(opts, "DepositSubTreeReady")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositDepositSubTreeReady)
				if err := _Deposit.contract.UnpackLog(event, "DepositSubTreeReady", log); err != nil {
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
func (_Deposit *DepositFilterer) ParseDepositSubTreeReady(log types.Log) (*DepositDepositSubTreeReady, error) {
	event := new(DepositDepositSubTreeReady)
	if err := _Deposit.contract.UnpackLog(event, "DepositSubTreeReady", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
