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
const DepositManagerABI = "[{\"inputs\":[{\"internalType\":\"contractITokenRegistry\",\"name\":\"_tokenRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_vault\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"maxSubtreeDepth\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"l2Amount\",\"type\":\"uint256\"}],\"name\":\"DepositQueued\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subtreeID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"subtreeRoot\",\"type\":\"bytes32\"}],\"name\":\"DepositSubTreeReady\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"babyTrees\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"babyTreesLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"back\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"l1Amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"}],\"name\":\"depositFor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dequeueToSubmit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"subtreeID\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"subtreeRoot\",\"type\":\"bytes32\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"front\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paramMaxSubtreeDepth\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paramMaxSubtreeSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"queue\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"subtreeRoot\",\"type\":\"bytes32\"}],\"name\":\"reenqueue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollup\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_rollup\",\"type\":\"address\"}],\"name\":\"setRollupAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tokenRegistry\",\"outputs\":[{\"internalType\":\"contractITokenRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"vault\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// DepositManagerBin is the compiled bytecode used for deploying new contracts.
var DepositManagerBin = "0x6101206040526001805560006002556000600455600060055534801561002457600080fd5b506040516110d23803806110d283398101604081905261004391610075565b6001811b60805233606090811b60a05260c0919091526001600160601b031992811b8316610100521b1660e0526100cf565b600080600060608486031215610089578283fd5b8351610094816100b7565b60208501519093506100a5816100b7565b80925050604084015190509250925092565b6001600160a01b03811681146100cc57600080fd5b50565b60805160a05160601c60c05160e05160601c6101005160601c610fa761012b60003980610321528061056152508061049e528061068d52508061065752508061027152806102f45250806105c752806108545250610fa76000f3fe608060405234801561001057600080fd5b50600436106101005760003560e01c8063ba75bbd811610097578063ddf0b00911610066578063ddf0b009146101b9578063e0203e94146101cc578063ee9d68ce146101d4578063fbfa77cf146101e757610100565b8063ba75bbd81461018b578063c7accaa414610193578063cb23bcb51461019b578063d86ee48d146101a357610100565b80638dde0840116100d35780638dde084014610155578063966fda621461015d5780639d23c4c714610170578063ade000261461017857610100565b806307663706146101055780632dfdf0b51461011a578063425e97f2146101385780638da5cb5b14610140575b600080fd5b610118610113366004610b2e565b6101ef565b005b6101226102e6565b60405161012f9190610c7b565b60405180910390f35b6101226102ec565b6101486102f2565b60405161012f9190610c29565b610122610316565b61011861016b366004610bc7565b61031c565b61014861055f565b610118610186366004610b97565b610583565b6101226105bf565b6101226105c5565b6101486105e9565b6101ab6105fe565b60405161012f929190610c84565b6101226101c7366004610b97565b610643565b610122610655565b6101226101e2366004610b97565b610679565b61014861068b565b600654610100900460ff168061020857506102086106af565b80610216575060065460ff16155b61023b5760405162461bcd60e51b815260040161023290610d97565b60405180910390fd5b600654610100900460ff16158015610266576006805460ff1961ff0019909116610100171660011790555b336001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016146102ae5760405162461bcd60e51b815260040161023290610e79565b6006805462010000600160b01b031916620100006001600160a01b0385160217905580156102e2576006805461ff00191690555b5050565b60055481565b60045481565b7f000000000000000000000000000000000000000000000000000000000000000090565b60025481565b6000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316630a7973b7846040518263ffffffff1660e01b815260040161036b9190610c7b565b604080518083038186803b15801561038257600080fd5b505afa158015610396573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103ba9190610b4a565b9150915083600014806103d457508084816103d157fe5b06155b6103f05760405162461bcd60e51b815260040161023290610cd3565b604051636eb1769f60e11b815284906001600160a01b0384169063dd62ed3e906104209033903090600401610c3d565b60206040518083038186803b15801561043857600080fd5b505afa15801561044c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104709190610baf565b101561048e5760405162461bcd60e51b815260040161023290610d60565b6104c36001600160a01b038316337f0000000000000000000000000000000000000000000000000000000000000000876106c0565b60008185816104ce57fe5b0490506104d9610b06565b60405180608001604052808881526020018681526020018381526020016000815250905060606105088261071e565b90507fe61c9511bbd1fc8fe4fefa18e7e1b1a21da551206ff80e07f4dbcc8829878f3e88878560405161053d93929190610f43565b60405180910390a161055581805190602001206107b5565b5050505050505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b6006546201000090046001600160a01b031633146105b35760405162461bcd60e51b815260040161023290610de5565b6105bc816108bd565b50565b60015481565b7f000000000000000000000000000000000000000000000000000000000000000081565b6006546201000090046001600160a01b031681565b60065460009081906201000090046001600160a01b031633146106335760405162461bcd60e51b815260040161023290610de5565b61063b610917565b915091509091565b60006020819052908152604090205481565b7f000000000000000000000000000000000000000000000000000000000000000081565b60036020526000908152604090205481565b7f000000000000000000000000000000000000000000000000000000000000000081565b60006106ba3061095e565b15905090565b610718846323b872dd60e01b8585856040516024016106e193929190610c57565b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b031990931692909217909152610964565b50505050565b805160609015801561073257506020820151155b801561074057506040820151155b801561074e57506060820151155b1561077b5760006040516020016107659190610c92565b60405160208183030381529060405290506107b0565b81516020808401516040808601516060870151915161079e959491929101610c0e565b60405160208183030381529060405290505b919050565b600580546001908101918290556004546000818152600360205260409020849055015b6001821661084a5760011981016000908152600360209081526040808320546000198501845292819020549051610810939201610c84565b60408051601f19818403018152918152815160209283012060011984016000908152600390935291205560019190911c90600019016107d8565b60048190556005547f000000000000000000000000000000000000000000000000000000000000000014156108b8576000805260036020527f3617319a054d772f909f7c479a2cebe5066e836a939412e32403c99029b92eff546108ad906108bd565b600060048190556005555b505050565b600280546001019081905560008181526020819052604090819020839055517f4d3a5844ed7dad1eee8f5c6143c14063c6944b6000cc260955d11b0706ecb4929061090b9083908590610c84565b60405180910390a15050565b60015460025460009082111561093f5760405162461bcd60e51b815260040161023290610e32565b5060008181526020819052604081208054919055600180830190559091565b3b151590565b60606109b9826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b03166109f39092919063ffffffff16565b8051909150156108b857808060200190518101906109d79190610b77565b6108b85760405162461bcd60e51b815260040161023290610ef9565b6060610a028484600085610a0c565b90505b9392505050565b606082471015610a2e5760405162461bcd60e51b815260040161023290610d1a565b610a378561095e565b610a535760405162461bcd60e51b815260040161023290610ec2565b60006060866001600160a01b03168587604051610a709190610bf2565b60006040518083038185875af1925050503d8060008114610aad576040519150601f19603f3d011682016040523d82523d6000602084013e610ab2565b606091505b5091509150610ac2828286610acd565b979650505050505050565b60608315610adc575081610a05565b825115610aec5782518084602001fd5b8160405162461bcd60e51b81526004016102329190610ca0565b6040518060800160405280600081526020016000815260200160008152602001600081525090565b600060208284031215610b3f578081fd5b8135610a0581610f85565b60008060408385031215610b5c578081fd5b8251610b6781610f85565b6020939093015192949293505050565b600060208284031215610b88578081fd5b81518015158114610a05578182fd5b600060208284031215610ba8578081fd5b5035919050565b600060208284031215610bc0578081fd5b5051919050565b600080600060608486031215610bdb578081fd5b505081359360208301359350604090920135919050565b60008251610c04818460208701610f59565b9190910192915050565b93845260208401929092526040830152606082015260800190565b6001600160a01b0391909116815260200190565b6001600160a01b0392831681529116602082015260400190565b6001600160a01b039384168152919092166020820152604081019190915260600190565b90815260200190565b918252602082015260400190565b60ff91909116815260200190565b6000602082528251806020840152610cbf816040850160208701610f59565b601f01601f19169190910160400192915050565b60208082526027908201527f6c31416d6f756e742073686f756c642062652061206d756c7469706c65206f66604082015266081b0c955b9a5d60ca1b606082015260800190565b60208082526026908201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6040820152651c8818d85b1b60d21b606082015260800190565b6020808252601c908201527f746f6b656e20616c6c6f77616e6365206e6f7420617070726f76656400000000604082015260600190565b6020808252602e908201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160408201526d191e481a5b9a5d1a585b1a5e995960921b606082015260800190565b6020808252602d908201527f4465706f7369744d616e616765723a2073656e646572206973206e6f7420526f60408201526c1b1b1d5c0818dbdb9d1c9858dd609a1b606082015260800190565b60208082526027908201527f4465706f73697420436f72653a2051756575652073686f756c64206265206e6f6040820152666e2d656d70747960c81b606082015260800190565b60208082526029908201527f496d6d757461626c654f776e61626c653a2063616c6c6572206973206e6f74206040820152683a34329037bbb732b960b91b606082015260800190565b6020808252601d908201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604082015260600190565b6020808252602a908201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6040820152691bdd081cdd58d8d9595960b21b606082015260800190565b9283526020830191909152604082015260600190565b60005b83811015610f74578181015183820152602001610f5c565b838111156107185750506000910152565b6001600160a01b03811681146105bc57600080fdfea164736f6c634300060c000a"

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

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DepositManager *DepositManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DepositManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DepositManager *DepositManagerSession) Owner() (common.Address, error) {
	return _DepositManager.Contract.Owner(&_DepositManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DepositManager *DepositManagerCallerSession) Owner() (common.Address, error) {
	return _DepositManager.Contract.Owner(&_DepositManager.CallOpts)
}

// ParamMaxSubtreeDepth is a free data retrieval call binding the contract method 0xe0203e94.
//
// Solidity: function paramMaxSubtreeDepth() view returns(uint256)
func (_DepositManager *DepositManagerCaller) ParamMaxSubtreeDepth(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DepositManager.contract.Call(opts, &out, "paramMaxSubtreeDepth")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ParamMaxSubtreeDepth is a free data retrieval call binding the contract method 0xe0203e94.
//
// Solidity: function paramMaxSubtreeDepth() view returns(uint256)
func (_DepositManager *DepositManagerSession) ParamMaxSubtreeDepth() (*big.Int, error) {
	return _DepositManager.Contract.ParamMaxSubtreeDepth(&_DepositManager.CallOpts)
}

// ParamMaxSubtreeDepth is a free data retrieval call binding the contract method 0xe0203e94.
//
// Solidity: function paramMaxSubtreeDepth() view returns(uint256)
func (_DepositManager *DepositManagerCallerSession) ParamMaxSubtreeDepth() (*big.Int, error) {
	return _DepositManager.Contract.ParamMaxSubtreeDepth(&_DepositManager.CallOpts)
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
// Solidity: function depositFor(uint256 pubkeyID, uint256 l1Amount, uint256 tokenID) returns()
func (_DepositManager *DepositManagerTransactor) DepositFor(opts *bind.TransactOpts, pubkeyID *big.Int, l1Amount *big.Int, tokenID *big.Int) (*types.Transaction, error) {
	return _DepositManager.contract.Transact(opts, "depositFor", pubkeyID, l1Amount, tokenID)
}

// DepositFor is a paid mutator transaction binding the contract method 0x966fda62.
//
// Solidity: function depositFor(uint256 pubkeyID, uint256 l1Amount, uint256 tokenID) returns()
func (_DepositManager *DepositManagerSession) DepositFor(pubkeyID *big.Int, l1Amount *big.Int, tokenID *big.Int) (*types.Transaction, error) {
	return _DepositManager.Contract.DepositFor(&_DepositManager.TransactOpts, pubkeyID, l1Amount, tokenID)
}

// DepositFor is a paid mutator transaction binding the contract method 0x966fda62.
//
// Solidity: function depositFor(uint256 pubkeyID, uint256 l1Amount, uint256 tokenID) returns()
func (_DepositManager *DepositManagerTransactorSession) DepositFor(pubkeyID *big.Int, l1Amount *big.Int, tokenID *big.Int) (*types.Transaction, error) {
	return _DepositManager.Contract.DepositFor(&_DepositManager.TransactOpts, pubkeyID, l1Amount, tokenID)
}

// DequeueToSubmit is a paid mutator transaction binding the contract method 0xd86ee48d.
//
// Solidity: function dequeueToSubmit() returns(uint256 subtreeID, bytes32 subtreeRoot)
func (_DepositManager *DepositManagerTransactor) DequeueToSubmit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DepositManager.contract.Transact(opts, "dequeueToSubmit")
}

// DequeueToSubmit is a paid mutator transaction binding the contract method 0xd86ee48d.
//
// Solidity: function dequeueToSubmit() returns(uint256 subtreeID, bytes32 subtreeRoot)
func (_DepositManager *DepositManagerSession) DequeueToSubmit() (*types.Transaction, error) {
	return _DepositManager.Contract.DequeueToSubmit(&_DepositManager.TransactOpts)
}

// DequeueToSubmit is a paid mutator transaction binding the contract method 0xd86ee48d.
//
// Solidity: function dequeueToSubmit() returns(uint256 subtreeID, bytes32 subtreeRoot)
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
	TokenID  *big.Int
	L2Amount *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDepositQueued is a free log retrieval operation binding the contract event 0xe61c9511bbd1fc8fe4fefa18e7e1b1a21da551206ff80e07f4dbcc8829878f3e.
//
// Solidity: event DepositQueued(uint256 pubkeyID, uint256 tokenID, uint256 l2Amount)
func (_DepositManager *DepositManagerFilterer) FilterDepositQueued(opts *bind.FilterOpts) (*DepositManagerDepositQueuedIterator, error) {

	logs, sub, err := _DepositManager.contract.FilterLogs(opts, "DepositQueued")
	if err != nil {
		return nil, err
	}
	return &DepositManagerDepositQueuedIterator{contract: _DepositManager.contract, event: "DepositQueued", logs: logs, sub: sub}, nil
}

// WatchDepositQueued is a free log subscription operation binding the contract event 0xe61c9511bbd1fc8fe4fefa18e7e1b1a21da551206ff80e07f4dbcc8829878f3e.
//
// Solidity: event DepositQueued(uint256 pubkeyID, uint256 tokenID, uint256 l2Amount)
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

// ParseDepositQueued is a log parse operation binding the contract event 0xe61c9511bbd1fc8fe4fefa18e7e1b1a21da551206ff80e07f4dbcc8829878f3e.
//
// Solidity: event DepositQueued(uint256 pubkeyID, uint256 tokenID, uint256 l2Amount)
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
	SubtreeID   *big.Int
	SubtreeRoot [32]byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDepositSubTreeReady is a free log retrieval operation binding the contract event 0x4d3a5844ed7dad1eee8f5c6143c14063c6944b6000cc260955d11b0706ecb492.
//
// Solidity: event DepositSubTreeReady(uint256 subtreeID, bytes32 subtreeRoot)
func (_DepositManager *DepositManagerFilterer) FilterDepositSubTreeReady(opts *bind.FilterOpts) (*DepositManagerDepositSubTreeReadyIterator, error) {

	logs, sub, err := _DepositManager.contract.FilterLogs(opts, "DepositSubTreeReady")
	if err != nil {
		return nil, err
	}
	return &DepositManagerDepositSubTreeReadyIterator{contract: _DepositManager.contract, event: "DepositSubTreeReady", logs: logs, sub: sub}, nil
}

// WatchDepositSubTreeReady is a free log subscription operation binding the contract event 0x4d3a5844ed7dad1eee8f5c6143c14063c6944b6000cc260955d11b0706ecb492.
//
// Solidity: event DepositSubTreeReady(uint256 subtreeID, bytes32 subtreeRoot)
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

// ParseDepositSubTreeReady is a log parse operation binding the contract event 0x4d3a5844ed7dad1eee8f5c6143c14063c6944b6000cc260955d11b0706ecb492.
//
// Solidity: event DepositSubTreeReady(uint256 subtreeID, bytes32 subtreeRoot)
func (_DepositManager *DepositManagerFilterer) ParseDepositSubTreeReady(log types.Log) (*DepositManagerDepositSubTreeReady, error) {
	event := new(DepositManagerDepositSubTreeReady)
	if err := _DepositManager.contract.UnpackLog(event, "DepositSubTreeReady", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
