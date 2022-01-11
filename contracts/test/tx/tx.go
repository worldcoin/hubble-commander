// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tx

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

// TxCreate2Transfer is an auto generated low-level Go binding around an user-defined struct.
type TxCreate2Transfer struct {
	FromIndex  *big.Int
	ToIndex    *big.Int
	ToPubkeyID *big.Int
	Amount     *big.Int
	Fee        *big.Int
}

// TxMassMigration is an auto generated low-level Go binding around an user-defined struct.
type TxMassMigration struct {
	FromIndex *big.Int
	Amount    *big.Int
	Fee       *big.Int
}

// TxTransfer is an auto generated low-level Go binding around an user-defined struct.
type TxTransfer struct {
	FromIndex *big.Int
	ToIndex   *big.Int
	Amount    *big.Int
	Fee       *big.Int
}

// TestTxMetaData contains all meta data concerning the TestTx contract.
var TestTxMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"create2TransferDecode\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"fromIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"toIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"toPubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"internalType\":\"structTx.Create2Transfer\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"to\",\"type\":\"bytes32\"}],\"name\":\"create2TransferMessageOf\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"name\":\"create2transferHasExcessData\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"fromIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"toIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"toPubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"internalType\":\"structTx.Create2Transfer[]\",\"name\":\"txs\",\"type\":\"tuple[]\"}],\"name\":\"create2transferSerialize\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"name\":\"create2transferSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"massMigrationDecode\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"fromIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"internalType\":\"structTx.MassMigration\",\"name\":\"_tx\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"name\":\"massMigrationSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"testEncodeDecimal\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"fromIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"internalType\":\"structTx.MassMigration\",\"name\":\"_tx\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"spokeID\",\"type\":\"uint256\"}],\"name\":\"testMassMigrationMessageOf\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"transferDecode\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"fromIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"toIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"internalType\":\"structTx.Transfer\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"name\":\"transferHasExcessData\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"transferMessageOf\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"fromIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"toIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"internalType\":\"structTx.Transfer[]\",\"name\":\"txs\",\"type\":\"tuple[]\"}],\"name\":\"transferSerialize\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"name\":\"transferSize\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610f8a806100206000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c80639a06678f1161008c578063c0bea13c11610066578063c0bea13c146101ea578063c2b378361461020a578063e177b4fe1461022a578063e9b344cb1461023d576100ea565b80639a06678f146101b1578063b116085a146101c4578063bf20a9df146101d7576100ea565b806361dd4671116100c857806361dd46711461014b57806364ab8e801461016b5780636e0af50f1461017e5780637cef563f1461019e576100ea565b80631c5f31f7146100ef5780633283ae43146101185780633bad0c7f14610138575b600080fd5b6101026100fd366004610a2d565b610250565b60405161010f9190610e31565b60405180910390f35b61012b610126366004610b97565b610261565b60405161010f9190610f2d565b61012b610146366004610b97565b61026c565b61015e610159366004610b97565b610277565b60405161010f9190610e26565b61012b610179366004610d05565b610282565b61019161018c366004610bca565b61028d565b60405161010f9190610f02565b61015e6101ac366004610b97565b6102a6565b6101026101bf366004610b07565b6102b1565b6101026101d2366004610c59565b6102bc565b6101026101e5366004610c0d565b6102db565b6101fd6101f8366004610bca565b6102f8565b60405161010f9190610ee1565b61021d610218366004610bca565b61030a565b60405161010f9190610ea7565b61012b610238366004610b97565b61031c565b61010261024b366004610cab565b610327565b606061025b82610334565b92915050565b600061025b826104c6565b600061025b826104da565b600061025b826104e7565b600061025b826104fd565b610295610907565b61029f8383610575565b9392505050565b600061025b826105e8565b606061025b826105f5565b60606102d26102cb8686610760565b84846107e3565b95945050505050565b60606102f06102ea8585610575565b83610825565b949350505050565b61030061092f565b61029f838361086a565b610312610950565b61029f8383610760565b600061025b826108d0565b60606102f08484846108dd565b8051606090816010820267ffffffffffffffff8111801561035457600080fd5b506040519080825280601f01601f19166020018201604052801561037f576020820181803683370190505b50905060005b828110156104be57600085828151811061039b57fe5b602002602001015160000151905060008683815181106103b757fe5b602002602001015160200151905060008784815181106103d357fe5b602002602001015160400151905060008885815181106103ef57fe5b6020026020010151606001519050600089868151811061040b57fe5b60200260200101516080015190506060858585610427866104fd565b610430866104fd565b604051602001610444959493929190610d87565b60408051601f1981840301815291905290506010870260005b60108110156104aa5782818151811061047257fe5b602001015160f81c60f81b8a8383018151811061048b57fe5b60200101906001600160f81b031916908160001a90535060010161045d565b505060019096019550610385945050505050565b509392505050565b600060108251816104d357fe5b0492915050565b6000600c8251816104d357fe5b600060108251816104f457fe5b06151592915050565b60008181805b600f81101561053f57821580159061051c5750600a8306155b1561053257600a83049250600182019150610537565b61053f565b600101610503565b50610fff82111561056b5760405162461bcd60e51b815260040161056290610e84565b60405180910390fd5b600c1b0192915050565b61057d610907565b506004600c8281028401918201516008830151600a80850151948401516040805160808101825263ffffffff9586168152939094166020840152600f86861c8116830a610fff97881602948401949094529384901c90921690910a9190921602606082015292915050565b6000600c8251816104f457fe5b805160609081600c820267ffffffffffffffff8111801561061557600080fd5b506040519080825280601f01601f191660200182016040528015610640576020820181803683370190505b50905060005b828110156104be57600085828151811061065c57fe5b6020026020010151600001519050600086838151811061067857fe5b602002602001015160200151905060006106a888858151811061069757fe5b6020026020010151604001516104fd565b905060006106cc8986815181106106bb57fe5b6020026020010151606001516104fd565b90506060848484846040516020016106e79493929190610d45565b60408051601f198184030181529190529050600c860260005b600c81101561074d5782818151811061071557fe5b602001015160f81c60f81b898383018151811061072e57fe5b60200101906001600160f81b031916908160001a905350600101610700565b5050600190950194506106469350505050565b610768610950565b50600460108281028401918201516008830151600c80850151600e86015195909401516040805160a08101825263ffffffff958616815293851660208501529490931693820193909352600f84841c8116600a90810a610fff9687160260608401529383901c1690920a921691909102608082015292915050565b60606003846000015183858760600151886080015160405160200161080d96959493929190610d1d565b60405160208183030381529060405290509392505050565b6060600183600001518460200151848660400151876060015160405160200161085396959493929190610d1d565b604051602081830303815290604052905092915050565b61087261092f565b506008908102919091016004810151600682015191909201516040805160608101825263ffffffff9094168452610fff808416600c94851c600f908116600a90810a9290920260208801529184169390941c1690920a029082015290565b600060088251816104d357fe5b60606002846000015185602001518660400151868660405160200161080d96959493929190610dd6565b6040518060800160405280600081526020016000815260200160008152602001600081525090565b60405180606001604052806000815260200160008152602001600081525090565b6040518060a0016040528060008152602001600081526020016000815260200160008152602001600081525090565b600082601f83011261098f578081fd5b813567ffffffffffffffff8111156109a5578182fd5b6109b8601f8201601f1916602001610f36565b91508082528360208285010111156109cf57600080fd5b8060208401602084013760009082016020015292915050565b6000608082840312156109f9578081fd5b610a036080610f36565b90508135815260208201356020820152604082013560408201526060820135606082015292915050565b60006020808385031215610a3f578182fd5b823567ffffffffffffffff811115610a55578283fd5b8301601f81018513610a65578283fd5b8035610a78610a7382610f5d565b610f36565b8181528381019083850160a0808502860187018a1015610a96578788fd5b8795505b84861015610af95780828b031215610ab0578788fd5b610ab981610f36565b8235815287830135888201526040808401359082015260608084013590820152608080840135908201528452600195909501949286019290810190610a9a565b509098975050505050505050565b60006020808385031215610b19578182fd5b823567ffffffffffffffff811115610b2f578283fd5b8301601f81018513610b3f578283fd5b8035610b4d610a7382610f5d565b818152838101908385016080808502860187018a1015610b6b578788fd5b8795505b84861015610af957610b818a836109e8565b8452600195909501949286019290810190610b6f565b600060208284031215610ba8578081fd5b813567ffffffffffffffff811115610bbe578182fd5b6102f08482850161097f565b60008060408385031215610bdc578081fd5b823567ffffffffffffffff811115610bf2578182fd5b610bfe8582860161097f565b95602094909401359450505050565b600080600060608486031215610c21578081fd5b833567ffffffffffffffff811115610c37578182fd5b610c438682870161097f565b9660208601359650604090950135949350505050565b60008060008060808587031215610c6e578081fd5b843567ffffffffffffffff811115610c84578182fd5b610c908782880161097f565b97602087013597506040870135966060013595509350505050565b600080600083850360a0811215610cc0578384fd5b6060811215610ccd578384fd5b50610cd86060610f36565b84358152602080860135908201526040808601359082015295606085013595506080909401359392505050565b600060208284031215610d16578081fd5b5035919050565b958652602086019490945260408501929092526060840152608083015260a082015260c00190565b6001600160e01b031960e095861b811682529390941b90921660048401526001600160f01b031960f091821b8116600885015291901b16600a820152600c0190565b6001600160e01b031960e096871b8116825294861b851660048201529290941b90921660088201526001600160f01b031960f092831b8116600c8301529290911b909116600e82015260100190565b60f89690961b6001600160f81b031916865260e094851b6001600160e01b0319908116600188015260058701949094526025860192909252831b8216604585015290911b166049820152604d0190565b901515815260200190565b6000602080835283518082850152825b81811015610e5d57858101830151858201604001528201610e41565b81811115610e6e5783604083870101525b50601f01601f1916929092016040019392505050565b602080825260099082015268109859081a5b9c1d5d60ba1b604082015260600190565b600060a082019050825182526020830151602083015260408301516040830152606083015160608301526080830151608083015292915050565b81518152602080830151908201526040918201519181019190915260600190565b8151815260208083015190820152604080830151908201526060918201519181019190915260800190565b90815260200190565b60405181810167ffffffffffffffff81118282101715610f5557600080fd5b604052919050565b600067ffffffffffffffff821115610f73578081fd5b506020908102019056fea164736f6c634300060c000a",
}

// TestTxABI is the input ABI used to generate the binding from.
// Deprecated: Use TestTxMetaData.ABI instead.
var TestTxABI = TestTxMetaData.ABI

// TestTxBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TestTxMetaData.Bin instead.
var TestTxBin = TestTxMetaData.Bin

// DeployTestTx deploys a new Ethereum contract, binding an instance of TestTx to it.
func DeployTestTx(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TestTx, error) {
	parsed, err := TestTxMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TestTxBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TestTx{TestTxCaller: TestTxCaller{contract: contract}, TestTxTransactor: TestTxTransactor{contract: contract}, TestTxFilterer: TestTxFilterer{contract: contract}}, nil
}

// TestTx is an auto generated Go binding around an Ethereum contract.
type TestTx struct {
	TestTxCaller     // Read-only binding to the contract
	TestTxTransactor // Write-only binding to the contract
	TestTxFilterer   // Log filterer for contract events
}

// TestTxCaller is an auto generated read-only Go binding around an Ethereum contract.
type TestTxCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestTxTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TestTxTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestTxFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TestTxFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TestTxSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TestTxSession struct {
	Contract     *TestTx           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestTxCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TestTxCallerSession struct {
	Contract *TestTxCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// TestTxTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TestTxTransactorSession struct {
	Contract     *TestTxTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TestTxRaw is an auto generated low-level Go binding around an Ethereum contract.
type TestTxRaw struct {
	Contract *TestTx // Generic contract binding to access the raw methods on
}

// TestTxCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TestTxCallerRaw struct {
	Contract *TestTxCaller // Generic read-only contract binding to access the raw methods on
}

// TestTxTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TestTxTransactorRaw struct {
	Contract *TestTxTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTestTx creates a new instance of TestTx, bound to a specific deployed contract.
func NewTestTx(address common.Address, backend bind.ContractBackend) (*TestTx, error) {
	contract, err := bindTestTx(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TestTx{TestTxCaller: TestTxCaller{contract: contract}, TestTxTransactor: TestTxTransactor{contract: contract}, TestTxFilterer: TestTxFilterer{contract: contract}}, nil
}

// NewTestTxCaller creates a new read-only instance of TestTx, bound to a specific deployed contract.
func NewTestTxCaller(address common.Address, caller bind.ContractCaller) (*TestTxCaller, error) {
	contract, err := bindTestTx(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TestTxCaller{contract: contract}, nil
}

// NewTestTxTransactor creates a new write-only instance of TestTx, bound to a specific deployed contract.
func NewTestTxTransactor(address common.Address, transactor bind.ContractTransactor) (*TestTxTransactor, error) {
	contract, err := bindTestTx(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TestTxTransactor{contract: contract}, nil
}

// NewTestTxFilterer creates a new log filterer instance of TestTx, bound to a specific deployed contract.
func NewTestTxFilterer(address common.Address, filterer bind.ContractFilterer) (*TestTxFilterer, error) {
	contract, err := bindTestTx(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TestTxFilterer{contract: contract}, nil
}

// bindTestTx binds a generic wrapper to an already deployed contract.
func bindTestTx(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TestTxABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestTx *TestTxRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestTx.Contract.TestTxCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestTx *TestTxRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestTx.Contract.TestTxTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestTx *TestTxRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestTx.Contract.TestTxTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TestTx *TestTxCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TestTx.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TestTx *TestTxTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TestTx.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TestTx *TestTxTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TestTx.Contract.contract.Transact(opts, method, params...)
}

// Create2TransferDecode is a free data retrieval call binding the contract method 0xc2b37836.
//
// Solidity: function create2TransferDecode(bytes txs, uint256 index) pure returns((uint256,uint256,uint256,uint256,uint256))
func (_TestTx *TestTxCaller) Create2TransferDecode(opts *bind.CallOpts, txs []byte, index *big.Int) (TxCreate2Transfer, error) {
	var out []interface{}
	err := _TestTx.contract.Call(opts, &out, "create2TransferDecode", txs, index)

	if err != nil {
		return *new(TxCreate2Transfer), err
	}

	out0 := *abi.ConvertType(out[0], new(TxCreate2Transfer)).(*TxCreate2Transfer)

	return out0, err

}

// Create2TransferDecode is a free data retrieval call binding the contract method 0xc2b37836.
//
// Solidity: function create2TransferDecode(bytes txs, uint256 index) pure returns((uint256,uint256,uint256,uint256,uint256))
func (_TestTx *TestTxSession) Create2TransferDecode(txs []byte, index *big.Int) (TxCreate2Transfer, error) {
	return _TestTx.Contract.Create2TransferDecode(&_TestTx.CallOpts, txs, index)
}

// Create2TransferDecode is a free data retrieval call binding the contract method 0xc2b37836.
//
// Solidity: function create2TransferDecode(bytes txs, uint256 index) pure returns((uint256,uint256,uint256,uint256,uint256))
func (_TestTx *TestTxCallerSession) Create2TransferDecode(txs []byte, index *big.Int) (TxCreate2Transfer, error) {
	return _TestTx.Contract.Create2TransferDecode(&_TestTx.CallOpts, txs, index)
}

// Create2TransferMessageOf is a free data retrieval call binding the contract method 0xb116085a.
//
// Solidity: function create2TransferMessageOf(bytes txs, uint256 index, uint256 nonce, bytes32 to) pure returns(bytes)
func (_TestTx *TestTxCaller) Create2TransferMessageOf(opts *bind.CallOpts, txs []byte, index *big.Int, nonce *big.Int, to [32]byte) ([]byte, error) {
	var out []interface{}
	err := _TestTx.contract.Call(opts, &out, "create2TransferMessageOf", txs, index, nonce, to)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Create2TransferMessageOf is a free data retrieval call binding the contract method 0xb116085a.
//
// Solidity: function create2TransferMessageOf(bytes txs, uint256 index, uint256 nonce, bytes32 to) pure returns(bytes)
func (_TestTx *TestTxSession) Create2TransferMessageOf(txs []byte, index *big.Int, nonce *big.Int, to [32]byte) ([]byte, error) {
	return _TestTx.Contract.Create2TransferMessageOf(&_TestTx.CallOpts, txs, index, nonce, to)
}

// Create2TransferMessageOf is a free data retrieval call binding the contract method 0xb116085a.
//
// Solidity: function create2TransferMessageOf(bytes txs, uint256 index, uint256 nonce, bytes32 to) pure returns(bytes)
func (_TestTx *TestTxCallerSession) Create2TransferMessageOf(txs []byte, index *big.Int, nonce *big.Int, to [32]byte) ([]byte, error) {
	return _TestTx.Contract.Create2TransferMessageOf(&_TestTx.CallOpts, txs, index, nonce, to)
}

// Create2transferHasExcessData is a free data retrieval call binding the contract method 0x61dd4671.
//
// Solidity: function create2transferHasExcessData(bytes txs) pure returns(bool)
func (_TestTx *TestTxCaller) Create2transferHasExcessData(opts *bind.CallOpts, txs []byte) (bool, error) {
	var out []interface{}
	err := _TestTx.contract.Call(opts, &out, "create2transferHasExcessData", txs)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Create2transferHasExcessData is a free data retrieval call binding the contract method 0x61dd4671.
//
// Solidity: function create2transferHasExcessData(bytes txs) pure returns(bool)
func (_TestTx *TestTxSession) Create2transferHasExcessData(txs []byte) (bool, error) {
	return _TestTx.Contract.Create2transferHasExcessData(&_TestTx.CallOpts, txs)
}

// Create2transferHasExcessData is a free data retrieval call binding the contract method 0x61dd4671.
//
// Solidity: function create2transferHasExcessData(bytes txs) pure returns(bool)
func (_TestTx *TestTxCallerSession) Create2transferHasExcessData(txs []byte) (bool, error) {
	return _TestTx.Contract.Create2transferHasExcessData(&_TestTx.CallOpts, txs)
}

// Create2transferSerialize is a free data retrieval call binding the contract method 0x1c5f31f7.
//
// Solidity: function create2transferSerialize((uint256,uint256,uint256,uint256,uint256)[] txs) pure returns(bytes)
func (_TestTx *TestTxCaller) Create2transferSerialize(opts *bind.CallOpts, txs []TxCreate2Transfer) ([]byte, error) {
	var out []interface{}
	err := _TestTx.contract.Call(opts, &out, "create2transferSerialize", txs)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// Create2transferSerialize is a free data retrieval call binding the contract method 0x1c5f31f7.
//
// Solidity: function create2transferSerialize((uint256,uint256,uint256,uint256,uint256)[] txs) pure returns(bytes)
func (_TestTx *TestTxSession) Create2transferSerialize(txs []TxCreate2Transfer) ([]byte, error) {
	return _TestTx.Contract.Create2transferSerialize(&_TestTx.CallOpts, txs)
}

// Create2transferSerialize is a free data retrieval call binding the contract method 0x1c5f31f7.
//
// Solidity: function create2transferSerialize((uint256,uint256,uint256,uint256,uint256)[] txs) pure returns(bytes)
func (_TestTx *TestTxCallerSession) Create2transferSerialize(txs []TxCreate2Transfer) ([]byte, error) {
	return _TestTx.Contract.Create2transferSerialize(&_TestTx.CallOpts, txs)
}

// Create2transferSize is a free data retrieval call binding the contract method 0x3283ae43.
//
// Solidity: function create2transferSize(bytes txs) pure returns(uint256)
func (_TestTx *TestTxCaller) Create2transferSize(opts *bind.CallOpts, txs []byte) (*big.Int, error) {
	var out []interface{}
	err := _TestTx.contract.Call(opts, &out, "create2transferSize", txs)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Create2transferSize is a free data retrieval call binding the contract method 0x3283ae43.
//
// Solidity: function create2transferSize(bytes txs) pure returns(uint256)
func (_TestTx *TestTxSession) Create2transferSize(txs []byte) (*big.Int, error) {
	return _TestTx.Contract.Create2transferSize(&_TestTx.CallOpts, txs)
}

// Create2transferSize is a free data retrieval call binding the contract method 0x3283ae43.
//
// Solidity: function create2transferSize(bytes txs) pure returns(uint256)
func (_TestTx *TestTxCallerSession) Create2transferSize(txs []byte) (*big.Int, error) {
	return _TestTx.Contract.Create2transferSize(&_TestTx.CallOpts, txs)
}

// MassMigrationDecode is a free data retrieval call binding the contract method 0xc0bea13c.
//
// Solidity: function massMigrationDecode(bytes txs, uint256 index) pure returns((uint256,uint256,uint256) _tx)
func (_TestTx *TestTxCaller) MassMigrationDecode(opts *bind.CallOpts, txs []byte, index *big.Int) (TxMassMigration, error) {
	var out []interface{}
	err := _TestTx.contract.Call(opts, &out, "massMigrationDecode", txs, index)

	if err != nil {
		return *new(TxMassMigration), err
	}

	out0 := *abi.ConvertType(out[0], new(TxMassMigration)).(*TxMassMigration)

	return out0, err

}

// MassMigrationDecode is a free data retrieval call binding the contract method 0xc0bea13c.
//
// Solidity: function massMigrationDecode(bytes txs, uint256 index) pure returns((uint256,uint256,uint256) _tx)
func (_TestTx *TestTxSession) MassMigrationDecode(txs []byte, index *big.Int) (TxMassMigration, error) {
	return _TestTx.Contract.MassMigrationDecode(&_TestTx.CallOpts, txs, index)
}

// MassMigrationDecode is a free data retrieval call binding the contract method 0xc0bea13c.
//
// Solidity: function massMigrationDecode(bytes txs, uint256 index) pure returns((uint256,uint256,uint256) _tx)
func (_TestTx *TestTxCallerSession) MassMigrationDecode(txs []byte, index *big.Int) (TxMassMigration, error) {
	return _TestTx.Contract.MassMigrationDecode(&_TestTx.CallOpts, txs, index)
}

// MassMigrationSize is a free data retrieval call binding the contract method 0xe177b4fe.
//
// Solidity: function massMigrationSize(bytes txs) pure returns(uint256)
func (_TestTx *TestTxCaller) MassMigrationSize(opts *bind.CallOpts, txs []byte) (*big.Int, error) {
	var out []interface{}
	err := _TestTx.contract.Call(opts, &out, "massMigrationSize", txs)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MassMigrationSize is a free data retrieval call binding the contract method 0xe177b4fe.
//
// Solidity: function massMigrationSize(bytes txs) pure returns(uint256)
func (_TestTx *TestTxSession) MassMigrationSize(txs []byte) (*big.Int, error) {
	return _TestTx.Contract.MassMigrationSize(&_TestTx.CallOpts, txs)
}

// MassMigrationSize is a free data retrieval call binding the contract method 0xe177b4fe.
//
// Solidity: function massMigrationSize(bytes txs) pure returns(uint256)
func (_TestTx *TestTxCallerSession) MassMigrationSize(txs []byte) (*big.Int, error) {
	return _TestTx.Contract.MassMigrationSize(&_TestTx.CallOpts, txs)
}

// TestEncodeDecimal is a free data retrieval call binding the contract method 0x64ab8e80.
//
// Solidity: function testEncodeDecimal(uint256 amount) pure returns(uint256)
func (_TestTx *TestTxCaller) TestEncodeDecimal(opts *bind.CallOpts, amount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _TestTx.contract.Call(opts, &out, "testEncodeDecimal", amount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TestEncodeDecimal is a free data retrieval call binding the contract method 0x64ab8e80.
//
// Solidity: function testEncodeDecimal(uint256 amount) pure returns(uint256)
func (_TestTx *TestTxSession) TestEncodeDecimal(amount *big.Int) (*big.Int, error) {
	return _TestTx.Contract.TestEncodeDecimal(&_TestTx.CallOpts, amount)
}

// TestEncodeDecimal is a free data retrieval call binding the contract method 0x64ab8e80.
//
// Solidity: function testEncodeDecimal(uint256 amount) pure returns(uint256)
func (_TestTx *TestTxCallerSession) TestEncodeDecimal(amount *big.Int) (*big.Int, error) {
	return _TestTx.Contract.TestEncodeDecimal(&_TestTx.CallOpts, amount)
}

// TestMassMigrationMessageOf is a free data retrieval call binding the contract method 0xe9b344cb.
//
// Solidity: function testMassMigrationMessageOf((uint256,uint256,uint256) _tx, uint256 nonce, uint256 spokeID) pure returns(bytes)
func (_TestTx *TestTxCaller) TestMassMigrationMessageOf(opts *bind.CallOpts, _tx TxMassMigration, nonce *big.Int, spokeID *big.Int) ([]byte, error) {
	var out []interface{}
	err := _TestTx.contract.Call(opts, &out, "testMassMigrationMessageOf", _tx, nonce, spokeID)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// TestMassMigrationMessageOf is a free data retrieval call binding the contract method 0xe9b344cb.
//
// Solidity: function testMassMigrationMessageOf((uint256,uint256,uint256) _tx, uint256 nonce, uint256 spokeID) pure returns(bytes)
func (_TestTx *TestTxSession) TestMassMigrationMessageOf(_tx TxMassMigration, nonce *big.Int, spokeID *big.Int) ([]byte, error) {
	return _TestTx.Contract.TestMassMigrationMessageOf(&_TestTx.CallOpts, _tx, nonce, spokeID)
}

// TestMassMigrationMessageOf is a free data retrieval call binding the contract method 0xe9b344cb.
//
// Solidity: function testMassMigrationMessageOf((uint256,uint256,uint256) _tx, uint256 nonce, uint256 spokeID) pure returns(bytes)
func (_TestTx *TestTxCallerSession) TestMassMigrationMessageOf(_tx TxMassMigration, nonce *big.Int, spokeID *big.Int) ([]byte, error) {
	return _TestTx.Contract.TestMassMigrationMessageOf(&_TestTx.CallOpts, _tx, nonce, spokeID)
}

// TransferDecode is a free data retrieval call binding the contract method 0x6e0af50f.
//
// Solidity: function transferDecode(bytes txs, uint256 index) pure returns((uint256,uint256,uint256,uint256))
func (_TestTx *TestTxCaller) TransferDecode(opts *bind.CallOpts, txs []byte, index *big.Int) (TxTransfer, error) {
	var out []interface{}
	err := _TestTx.contract.Call(opts, &out, "transferDecode", txs, index)

	if err != nil {
		return *new(TxTransfer), err
	}

	out0 := *abi.ConvertType(out[0], new(TxTransfer)).(*TxTransfer)

	return out0, err

}

// TransferDecode is a free data retrieval call binding the contract method 0x6e0af50f.
//
// Solidity: function transferDecode(bytes txs, uint256 index) pure returns((uint256,uint256,uint256,uint256))
func (_TestTx *TestTxSession) TransferDecode(txs []byte, index *big.Int) (TxTransfer, error) {
	return _TestTx.Contract.TransferDecode(&_TestTx.CallOpts, txs, index)
}

// TransferDecode is a free data retrieval call binding the contract method 0x6e0af50f.
//
// Solidity: function transferDecode(bytes txs, uint256 index) pure returns((uint256,uint256,uint256,uint256))
func (_TestTx *TestTxCallerSession) TransferDecode(txs []byte, index *big.Int) (TxTransfer, error) {
	return _TestTx.Contract.TransferDecode(&_TestTx.CallOpts, txs, index)
}

// TransferHasExcessData is a free data retrieval call binding the contract method 0x7cef563f.
//
// Solidity: function transferHasExcessData(bytes txs) pure returns(bool)
func (_TestTx *TestTxCaller) TransferHasExcessData(opts *bind.CallOpts, txs []byte) (bool, error) {
	var out []interface{}
	err := _TestTx.contract.Call(opts, &out, "transferHasExcessData", txs)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// TransferHasExcessData is a free data retrieval call binding the contract method 0x7cef563f.
//
// Solidity: function transferHasExcessData(bytes txs) pure returns(bool)
func (_TestTx *TestTxSession) TransferHasExcessData(txs []byte) (bool, error) {
	return _TestTx.Contract.TransferHasExcessData(&_TestTx.CallOpts, txs)
}

// TransferHasExcessData is a free data retrieval call binding the contract method 0x7cef563f.
//
// Solidity: function transferHasExcessData(bytes txs) pure returns(bool)
func (_TestTx *TestTxCallerSession) TransferHasExcessData(txs []byte) (bool, error) {
	return _TestTx.Contract.TransferHasExcessData(&_TestTx.CallOpts, txs)
}

// TransferMessageOf is a free data retrieval call binding the contract method 0xbf20a9df.
//
// Solidity: function transferMessageOf(bytes txs, uint256 index, uint256 nonce) pure returns(bytes)
func (_TestTx *TestTxCaller) TransferMessageOf(opts *bind.CallOpts, txs []byte, index *big.Int, nonce *big.Int) ([]byte, error) {
	var out []interface{}
	err := _TestTx.contract.Call(opts, &out, "transferMessageOf", txs, index, nonce)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// TransferMessageOf is a free data retrieval call binding the contract method 0xbf20a9df.
//
// Solidity: function transferMessageOf(bytes txs, uint256 index, uint256 nonce) pure returns(bytes)
func (_TestTx *TestTxSession) TransferMessageOf(txs []byte, index *big.Int, nonce *big.Int) ([]byte, error) {
	return _TestTx.Contract.TransferMessageOf(&_TestTx.CallOpts, txs, index, nonce)
}

// TransferMessageOf is a free data retrieval call binding the contract method 0xbf20a9df.
//
// Solidity: function transferMessageOf(bytes txs, uint256 index, uint256 nonce) pure returns(bytes)
func (_TestTx *TestTxCallerSession) TransferMessageOf(txs []byte, index *big.Int, nonce *big.Int) ([]byte, error) {
	return _TestTx.Contract.TransferMessageOf(&_TestTx.CallOpts, txs, index, nonce)
}

// TransferSerialize is a free data retrieval call binding the contract method 0x9a06678f.
//
// Solidity: function transferSerialize((uint256,uint256,uint256,uint256)[] txs) pure returns(bytes)
func (_TestTx *TestTxCaller) TransferSerialize(opts *bind.CallOpts, txs []TxTransfer) ([]byte, error) {
	var out []interface{}
	err := _TestTx.contract.Call(opts, &out, "transferSerialize", txs)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// TransferSerialize is a free data retrieval call binding the contract method 0x9a06678f.
//
// Solidity: function transferSerialize((uint256,uint256,uint256,uint256)[] txs) pure returns(bytes)
func (_TestTx *TestTxSession) TransferSerialize(txs []TxTransfer) ([]byte, error) {
	return _TestTx.Contract.TransferSerialize(&_TestTx.CallOpts, txs)
}

// TransferSerialize is a free data retrieval call binding the contract method 0x9a06678f.
//
// Solidity: function transferSerialize((uint256,uint256,uint256,uint256)[] txs) pure returns(bytes)
func (_TestTx *TestTxCallerSession) TransferSerialize(txs []TxTransfer) ([]byte, error) {
	return _TestTx.Contract.TransferSerialize(&_TestTx.CallOpts, txs)
}

// TransferSize is a free data retrieval call binding the contract method 0x3bad0c7f.
//
// Solidity: function transferSize(bytes txs) pure returns(uint256)
func (_TestTx *TestTxCaller) TransferSize(opts *bind.CallOpts, txs []byte) (*big.Int, error) {
	var out []interface{}
	err := _TestTx.contract.Call(opts, &out, "transferSize", txs)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TransferSize is a free data retrieval call binding the contract method 0x3bad0c7f.
//
// Solidity: function transferSize(bytes txs) pure returns(uint256)
func (_TestTx *TestTxSession) TransferSize(txs []byte) (*big.Int, error) {
	return _TestTx.Contract.TransferSize(&_TestTx.CallOpts, txs)
}

// TransferSize is a free data retrieval call binding the contract method 0x3bad0c7f.
//
// Solidity: function transferSize(bytes txs) pure returns(uint256)
func (_TestTx *TestTxCallerSession) TransferSize(txs []byte) (*big.Int, error) {
	return _TestTx.Contract.TransferSize(&_TestTx.CallOpts, txs)
}
