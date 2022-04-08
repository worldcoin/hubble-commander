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
	Bin: "0x608060405234801561001057600080fd5b5061122c806100206000396000f3fe608060405234801561001057600080fd5b50600436106100ea5760003560e01c80639a06678f1161008c578063c0bea13c11610066578063c0bea13c14610214578063c2b3783614610249578063e177b4fe1461029e578063e9b344cb146102b157600080fd5b80639a06678f146101db578063b116085a146101ee578063bf20a9df1461020157600080fd5b806361dd4671116100c857806361dd46711461014c57806364ab8e801461016f5780636e0af50f146101825780637cef563f146101c857600080fd5b80631c5f31f7146100ef5780633283ae43146101185780633bad0c7f14610139575b600080fd5b6101026100fd366004610d37565b6102c4565b60405161010f9190611073565b60405180910390f35b61012b610126366004610ed6565b6102d5565b60405190815260200161010f565b61012b610147366004610ed6565b6102e0565b61015f61015a366004610ed6565b6102eb565b604051901515815260200161010f565b61012b61017d36600461105b565b6102f6565b610195610190366004610f08565b610301565b60405161010f91908151815260208083015190820152604080830151908201526060918201519181019190915260800190565b61015f6101d6366004610ed6565b61033d565b6101026101e9366004610e0f565b610348565b6101026101fc366004610f95565b610353565b61010261020f366004610f4a565b610372565b610227610222366004610f08565b6103e1565b604080518251815260208084015190820152918101519082015260600161010f565b61025c610257366004610f08565b61040f565b60405161010f9190600060a082019050825182526020830151602083015260408301516040830152606083015160608301526080830151608083015292915050565b61012b6102ac366004610ed6565b61044b565b6101026102bf366004610fe6565b610456565b60606102cf82610463565b92915050565b60006102cf826106fc565b60006102cf8261070c565b60006102cf8261071c565b60006102cf82610734565b61032c6040518060800160405280600081526020016000815260200160008152602001600081525090565b61033683836107e6565b9392505050565b60006102cf8261087c565b60606102cf8261088c565b60606103696103628686610ad9565b8484610b86565b95945050505050565b60606103d961038185856107e6565b8051602080830151604080850151606095860151825160019581019590955284830195909552948301919091526080820187905260a082019390935260c0808201929092528251808203909201825260e00190915290565b949350505050565b61040560405180606001604052806000815260200160008152602001600081525090565b6103368383610bdc565b6104416040518060a0016040528060008152602001600081526020016000815260200160008152602001600081525090565b6103368383610ad9565b60006102cf82610c5e565b60606103d9848484610c6e565b8051606090600061047582601061118f565b6001600160401b0381111561049a57634e487b7160e01b600052604160045260246000fd5b6040519080825280601f01601f1916602001820160405280156104c4576020820181803683370190505b50905060005b828110156106f45760008582815181106104f457634e487b7160e01b600052603260045260246000fd5b6020026020010151600001519050600086838151811061052457634e487b7160e01b600052603260045260246000fd5b6020026020010151602001519050600087848151811061055457634e487b7160e01b600052603260045260246000fd5b6020026020010151604001519050600088858151811061058457634e487b7160e01b600052603260045260246000fd5b602002602001015160600151905060008986815181106105b457634e487b7160e01b600052603260045260246000fd5b602002602001015160800151905060008585856105d086610734565b6105d986610734565b6040516001600160e01b031960e096871b8116602083015294861b851660248201529290941b90921660288201526001600160f01b031960f092831b8116602c8301529290911b909116602e82015260300160408051601f198184030181529190529050600061064a60108961118f565b905060005b60108110156106d95782818151811061067857634e487b7160e01b600052603260045260246000fd5b01602001516001600160f81b0319168a6106928484611163565b815181106106b057634e487b7160e01b600052603260045260246000fd5b60200101906001600160f81b031916908160001a905350806106d1816111ae565b91505061064f565b505050505050505080806106ec906111ae565b9150506104ca565b509392505050565b6000601082516102cf919061117b565b6000600c82516102cf919061117b565b60006010825161072c91906111c9565b151592915050565b60008181805b600f81101561079557821580159061075a5750610758600a846111c9565b155b1561077e5761076a600a8461117b565b9250610777600183611163565b9150610783565b610795565b8061078d816111ae565b91505061073a565b50610fff8211156107d85760405162461bcd60e51b8152602060048201526009602482015268109859081a5b9c1d5d60ba1b604482015260640160405180910390fd5b6103d982600c83901b611163565b6108116040518060800160405280600081526020016000815260200160008152602001600081525090565b506004600c8281028401918201516008830151600a80850151948401516040805160808101825263ffffffff9586168152939094166020840152600f86861c8116830a610fff97881602948401949094529384901c90921690910a9190921602606082015292915050565b6000600c825161072c91906111c9565b8051606090600061089e82600c61118f565b6001600160401b038111156108c357634e487b7160e01b600052604160045260246000fd5b6040519080825280601f01601f1916602001820160405280156108ed576020820181803683370190505b50905060005b828110156106f457600085828151811061091d57634e487b7160e01b600052603260045260246000fd5b6020026020010151600001519050600086838151811061094d57634e487b7160e01b600052603260045260246000fd5b6020026020010151602001519050600061099188858151811061098057634e487b7160e01b600052603260045260246000fd5b602002602001015160400151610734565b905060006109c98986815181106109b857634e487b7160e01b600052603260045260246000fd5b602002602001015160600151610734565b6040516001600160e01b031960e087811b8216602084015286901b1660248201526001600160f01b031960f085811b8216602884015283901b16602a820152909150600090602c0160408051601f1981840301815291905290506000610a30600c8861118f565b905060005b600c811015610abf57828181518110610a5e57634e487b7160e01b600052603260045260246000fd5b01602001516001600160f81b03191689610a788484611163565b81518110610a9657634e487b7160e01b600052603260045260246000fd5b60200101906001600160f81b031916908160001a90535080610ab7816111ae565b915050610a35565b505050505050508080610ad1906111ae565b9150506108f3565b610b0b6040518060a0016040528060008152602001600081526020016000815260200160008152602001600081525090565b50600460108281028401918201516008830151600c80850151600e86015195909401516040805160a08101825263ffffffff958616815293851660208501529490931693820193909352600f84841c8116600a90810a610fff9687160260608401529383901c1690920a921691909102608082015292915050565b825160608481015160808087015160408051600360208201529081019590955284840186905290840186905260a084019190915260c08301529060e0015b60405160208183030381529060405290509392505050565b610c0060405180606001604052806000815260200160008152602001600081525090565b506008908102919091016004810151600682015191909201516040805160608101825263ffffffff9094168452610fff808416600c94851c600f908116600a90810a9290920260208801529184169390941c1690920a029082015290565b6000600882516102cf919061117b565b82516020808501516040808701519051600160f91b938101939093526001600160e01b031960e094851b811660218501526025840192909252604583015284831b811660658301529183901b9091166069820152606090606d01610bc4565b600082601f830112610cdd578081fd5b81356001600160401b03811115610cf657610cf6611209565b610d09601f8201601f1916602001611110565b818152846020838601011115610d1d578283fd5b816020850160208301379081016020019190915292915050565b60006020808385031215610d49578182fd5b82356001600160401b03811115610d5e578283fd5b8301601f81018513610d6e578283fd5b8035610d81610d7c82611140565b611110565b8181528381019083850160a0808502860187018a1015610d9f578788fd5b8795505b84861015610e015780828b031215610db9578788fd5b610dc16110c6565b8235815287830135888201526040808401359082015260608084013590820152608080840135908201528452600195909501949286019290810190610da3565b509098975050505050505050565b60006020808385031215610e21578182fd5b82356001600160401b03811115610e36578283fd5b8301601f81018513610e46578283fd5b8035610e54610d7c82611140565b80828252848201915084840188868560071b8701011115610e73578687fd5b8694505b83851015610eca57608080828b031215610e8f578788fd5b610e976110ee565b82358152878301358882015260408084013590820152606080840135908201528452600195909501949286019201610e77565b50979650505050505050565b600060208284031215610ee7578081fd5b81356001600160401b03811115610efc578182fd5b6103d984828501610ccd565b60008060408385031215610f1a578081fd5b82356001600160401b03811115610f2f578182fd5b610f3b85828601610ccd565b95602094909401359450505050565b600080600060608486031215610f5e578081fd5b83356001600160401b03811115610f73578182fd5b610f7f86828701610ccd565b9660208601359650604090950135949350505050565b60008060008060808587031215610faa578081fd5b84356001600160401b03811115610fbf578182fd5b610fcb87828801610ccd565b97602087013597506040870135966060013595509350505050565b600080600083850360a0811215610ffb578384fd5b6060811215611008578384fd5b50604051606081018181106001600160401b038211171561102b5761102b611209565b60409081528535825260208087013590830152858101359082015295606085013595506080909401359392505050565b60006020828403121561106c578081fd5b5035919050565b6000602080835283518082850152825b8181101561109f57858101830151858201604001528201611083565b818111156110b05783604083870101525b50601f01601f1916929092016040019392505050565b60405160a081016001600160401b03811182821017156110e8576110e8611209565b60405290565b604051608081016001600160401b03811182821017156110e8576110e8611209565b604051601f8201601f191681016001600160401b038111828210171561113857611138611209565b604052919050565b60006001600160401b0382111561115957611159611209565b5060051b60200190565b60008219821115611176576111766111dd565b500190565b60008261118a5761118a6111f3565b500490565b60008160001904831182151516156111a9576111a96111dd565b500290565b60006000198214156111c2576111c26111dd565b5060010190565b6000826111d8576111d86111f3565b500690565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052601260045260246000fd5b634e487b7160e01b600052604160045260246000fdfea164736f6c6343000804000a",
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
