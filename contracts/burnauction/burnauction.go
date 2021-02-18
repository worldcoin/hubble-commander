// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package burnauction

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

// BurnAuctionABI is the input ABI used to generate the binding from.
const BurnAuctionABI = "[{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_donationAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_donationNumerator\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"slot\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"coordinator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"NewBestBid\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BLOCKS_PER_SLOT\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DELTA_BLOCKS_INITIAL_SLOT\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DONATION_DENOMINATOR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"name\":\"auction\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"coordinator\",\"type\":\"address\"},{\"internalType\":\"uint128\",\"name\":\"amount\",\"type\":\"uint128\"},{\"internalType\":\"bool\",\"name\":\"initialized\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"bidAmount\",\"type\":\"uint256\"}],\"name\":\"bid\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numBlock\",\"type\":\"uint256\"}],\"name\":\"block2slot\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"currentSlot\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"deposits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"donationAddress\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"donationNumerator\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"genesisBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBlockNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getProposer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"witdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdrawDonation\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// BurnAuctionBin is the compiled bytecode used for deploying new contracts.
var BurnAuctionBin = "0x608060405234801561001057600080fd5b50604051610c35380380610c358339818101604052604081101561003357600080fd5b508051602090910151600054612710101561007f5760405162461bcd60e51b8152600401808060200182810382526029815260200180610c0c6029913960400191505060405180910390fd5b6103e861008a6100b7565b01600255600180546001600160a01b0319166001600160a01b0393909316929092179091556000556100bb565b4390565b610b42806100ca6000396000f3fe6080604052600436106100e85760003560e01c8063b404930e1161008a578063ec034bed11610059578063ec034bed1461029c578063f2182cf9146102b1578063f4cc8eea146102c6578063fc7e286d146102db576100e8565b8063b404930e146101ca578063d3dd08e2146101f4578063e81e1ccc14610256578063e9790d021461026b576100e8565b8063454a2ab3116100c6578063454a2ab3146101575780634cdc9c63146101765780639a478d241461018b578063a87a2ead146101a0576100e8565b80632243de47146100ed5780633359632e1461011b57806342cbb15c14610130575b600080fd5b3480156100f957600080fd5b5061010261030e565b6040805163ffffffff9092168252519081900360200190f35b34801561012757600080fd5b50610102610313565b34801561013c57600080fd5b5061014561032a565b60408051918252519081900360200190f35b6101746004803603602081101561016d57600080fd5b503561032e565b005b34801561018257600080fd5b5061014561051d565b34801561019757600080fd5b50610145610523565b3480156101ac57600080fd5b50610102600480360360208110156101c357600080fd5b5035610529565b3480156101d657600080fd5b50610174600480360360208110156101ed57600080fd5b503561054d565b34801561020057600080fd5b506102246004803603602081101561021757600080fd5b503563ffffffff166105e2565b604080516001600160a01b0390941684526001600160801b039092166020840152151582820152519081900360600190f35b34801561026257600080fd5b5061017461061c565b34801561027757600080fd5b506102806106c8565b604080516001600160a01b039092168252519081900360200190f35b3480156102a857600080fd5b5061028061076d565b3480156102bd57600080fd5b5061014561077c565b3480156102d257600080fd5b50610102610782565b3480156102e757600080fd5b50610145600480360360208110156102fe57600080fd5b50356001600160a01b0316610788565b606481565b600061032561032061032a565b610529565b905090565b4390565b6000610338610313565b60020163ffffffff81166000908152600360205260409020600101549091506001600160801b031680831161039e5760405162461bcd60e51b8152600401808060200182810382526023815260200180610a9d6023913960400191505060405180910390fd5b3360008181526004602052604090205434018411156103ee5760405162461bcd60e51b8152600401808060200182810382526030815260200180610a4c6030913960400191505060405180910390fd5b63ffffffff831660009081526003602052604081205461041b916001600160a01b0390911690849061079a565b600154600054610468916001600160a01b031690610448906127109061044290899061080f565b90610871565b6104636127106104426000548861080f90919063ffffffff16565b61079a565b61047381348661079a565b63ffffffff831660008181526003602090815260409182902080546001600160a01b03191633178155600101805460ff60801b196001600160801b038a166fffffffffffffffffffffffffffffffff199092169190911716600160801b17905581519283526001600160a01b03841690830152818101869052517f304f693446955254ce103ccf22f2ee397d8c2517f63076ee63e0dfa22fc5ba559181900360600190a150505050565b60025481565b61271081565b600060025482101561053d57506000610548565b506002546064908203045b919050565b3360008181526004602052604090205482111561059b5760405162461bcd60e51b815260040180806020018281038252603f815260200180610af7603f913960400191505060405180910390fd5b6105a78160008461079a565b6040516001600160a01b0382169083156108fc029084906000818181858888f193505050501580156105dd573d6000803e3d6000fd5b505050565b600360205260009081526040902080546001909101546001600160a01b03909116906001600160801b03811690600160801b900460ff1683565b6001546001600160a01b03166000908152600460205260409020546106725760405162461bcd60e51b8152600401808060200182810382526037815260200180610ac06037913960400191505060405180910390fd5b600180546001600160a01b039081166000908152600460205260408082208054908390559354905192169183156108fc0291849190818181858888f193505050501580156106c4573d6000803e3d6000fd5b5050565b6000806106d3610313565b63ffffffff8116600090815260036020526040902060010154909150600160801b900460ff1661074a576040805162461bcd60e51b815260206004820181905260248201527f41756374696f6e20686173206e6f74206265656e20696e697469616c697a6564604482015290519081900360640190fd5b63ffffffff166000908152600360205260409020546001600160a01b0316905090565b6001546001600160a01b031681565b60005481565b6103e881565b60046020526000908152604090205481565b6001600160a01b038316156105dd576001600160a01b0383166000908152600460205260409020546107cc90836108b3565b6001600160a01b03841660009081526004602052604090208190556107f1908261090d565b6001600160a01b038416600090815260046020526040902055505050565b60008261081e5750600061086b565b8282028284828161082b57fe5b04146108685760405162461bcd60e51b8152600401808060200182810382526021815260200180610a7c6021913960400191505060405180910390fd5b90505b92915050565b600061086883836040518060400160405280601a81526020017f536166654d6174683a206469766973696f6e206279207a65726f00000000000081525061094f565b600082820183811015610868576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b600061086883836040518060400160405280601e81526020017f536166654d6174683a207375627472616374696f6e206f766572666c6f7700008152506109f1565b600081836109db5760405162461bcd60e51b81526004018080602001828103825283818151815260200191508051906020019080838360005b838110156109a0578181015183820152602001610988565b50505050905090810190601f1680156109cd5780820380516001836020036101000a031916815260200191505b509250505060405180910390fd5b5060008385816109e757fe5b0495945050505050565b60008184841115610a435760405162461bcd60e51b81526020600482018181528351602484015283519092839260449091019190850190808383600083156109a0578181015183820152602001610988565b50505090039056fe4275726e41756374696f6e2c206269643a20696e73756666696369656e742066756e647320666f722062696464696e67536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f774275726e41756374696f6e2c206269643a206c657373207468656e2063757272656e744275726e41756374696f6e2c207769746864726177446f6e6174696f6e3a20646f6e6174696f6e206465706f736974206973207a65726f4275726e41756374696f6e2c2077697468647261773a20696e73756666696369656e74206465706f73697420616d6f756e7420666f72207769746864726177a164736f6c634300060c000a4275726e41756374696f6e2c20636f6e7374727563746f723a206261642064656e6f6d696e61746f72"

// DeployBurnAuction deploys a new Ethereum contract, binding an instance of BurnAuction to it.
func DeployBurnAuction(auth *bind.TransactOpts, backend bind.ContractBackend, _donationAddress common.Address, _donationNumerator *big.Int) (common.Address, *types.Transaction, *BurnAuction, error) {
	parsed, err := abi.JSON(strings.NewReader(BurnAuctionABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BurnAuctionBin), backend, _donationAddress, _donationNumerator)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BurnAuction{BurnAuctionCaller: BurnAuctionCaller{contract: contract}, BurnAuctionTransactor: BurnAuctionTransactor{contract: contract}, BurnAuctionFilterer: BurnAuctionFilterer{contract: contract}}, nil
}

// BurnAuction is an auto generated Go binding around an Ethereum contract.
type BurnAuction struct {
	BurnAuctionCaller     // Read-only binding to the contract
	BurnAuctionTransactor // Write-only binding to the contract
	BurnAuctionFilterer   // Log filterer for contract events
}

// BurnAuctionCaller is an auto generated read-only Go binding around an Ethereum contract.
type BurnAuctionCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BurnAuctionTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BurnAuctionTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BurnAuctionFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BurnAuctionFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BurnAuctionSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BurnAuctionSession struct {
	Contract     *BurnAuction      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BurnAuctionCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BurnAuctionCallerSession struct {
	Contract *BurnAuctionCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// BurnAuctionTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BurnAuctionTransactorSession struct {
	Contract     *BurnAuctionTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// BurnAuctionRaw is an auto generated low-level Go binding around an Ethereum contract.
type BurnAuctionRaw struct {
	Contract *BurnAuction // Generic contract binding to access the raw methods on
}

// BurnAuctionCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BurnAuctionCallerRaw struct {
	Contract *BurnAuctionCaller // Generic read-only contract binding to access the raw methods on
}

// BurnAuctionTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BurnAuctionTransactorRaw struct {
	Contract *BurnAuctionTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBurnAuction creates a new instance of BurnAuction, bound to a specific deployed contract.
func NewBurnAuction(address common.Address, backend bind.ContractBackend) (*BurnAuction, error) {
	contract, err := bindBurnAuction(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BurnAuction{BurnAuctionCaller: BurnAuctionCaller{contract: contract}, BurnAuctionTransactor: BurnAuctionTransactor{contract: contract}, BurnAuctionFilterer: BurnAuctionFilterer{contract: contract}}, nil
}

// NewBurnAuctionCaller creates a new read-only instance of BurnAuction, bound to a specific deployed contract.
func NewBurnAuctionCaller(address common.Address, caller bind.ContractCaller) (*BurnAuctionCaller, error) {
	contract, err := bindBurnAuction(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BurnAuctionCaller{contract: contract}, nil
}

// NewBurnAuctionTransactor creates a new write-only instance of BurnAuction, bound to a specific deployed contract.
func NewBurnAuctionTransactor(address common.Address, transactor bind.ContractTransactor) (*BurnAuctionTransactor, error) {
	contract, err := bindBurnAuction(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BurnAuctionTransactor{contract: contract}, nil
}

// NewBurnAuctionFilterer creates a new log filterer instance of BurnAuction, bound to a specific deployed contract.
func NewBurnAuctionFilterer(address common.Address, filterer bind.ContractFilterer) (*BurnAuctionFilterer, error) {
	contract, err := bindBurnAuction(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BurnAuctionFilterer{contract: contract}, nil
}

// bindBurnAuction binds a generic wrapper to an already deployed contract.
func bindBurnAuction(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BurnAuctionABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BurnAuction *BurnAuctionRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnAuction.Contract.BurnAuctionCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BurnAuction *BurnAuctionRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnAuction.Contract.BurnAuctionTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BurnAuction *BurnAuctionRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnAuction.Contract.BurnAuctionTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BurnAuction *BurnAuctionCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BurnAuction.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BurnAuction *BurnAuctionTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnAuction.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BurnAuction *BurnAuctionTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BurnAuction.Contract.contract.Transact(opts, method, params...)
}

// BLOCKSPERSLOT is a free data retrieval call binding the contract method 0x2243de47.
//
// Solidity: function BLOCKS_PER_SLOT() view returns(uint32)
func (_BurnAuction *BurnAuctionCaller) BLOCKSPERSLOT(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _BurnAuction.contract.Call(opts, &out, "BLOCKS_PER_SLOT")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// BLOCKSPERSLOT is a free data retrieval call binding the contract method 0x2243de47.
//
// Solidity: function BLOCKS_PER_SLOT() view returns(uint32)
func (_BurnAuction *BurnAuctionSession) BLOCKSPERSLOT() (uint32, error) {
	return _BurnAuction.Contract.BLOCKSPERSLOT(&_BurnAuction.CallOpts)
}

// BLOCKSPERSLOT is a free data retrieval call binding the contract method 0x2243de47.
//
// Solidity: function BLOCKS_PER_SLOT() view returns(uint32)
func (_BurnAuction *BurnAuctionCallerSession) BLOCKSPERSLOT() (uint32, error) {
	return _BurnAuction.Contract.BLOCKSPERSLOT(&_BurnAuction.CallOpts)
}

// DELTABLOCKSINITIALSLOT is a free data retrieval call binding the contract method 0xf4cc8eea.
//
// Solidity: function DELTA_BLOCKS_INITIAL_SLOT() view returns(uint32)
func (_BurnAuction *BurnAuctionCaller) DELTABLOCKSINITIALSLOT(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _BurnAuction.contract.Call(opts, &out, "DELTA_BLOCKS_INITIAL_SLOT")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// DELTABLOCKSINITIALSLOT is a free data retrieval call binding the contract method 0xf4cc8eea.
//
// Solidity: function DELTA_BLOCKS_INITIAL_SLOT() view returns(uint32)
func (_BurnAuction *BurnAuctionSession) DELTABLOCKSINITIALSLOT() (uint32, error) {
	return _BurnAuction.Contract.DELTABLOCKSINITIALSLOT(&_BurnAuction.CallOpts)
}

// DELTABLOCKSINITIALSLOT is a free data retrieval call binding the contract method 0xf4cc8eea.
//
// Solidity: function DELTA_BLOCKS_INITIAL_SLOT() view returns(uint32)
func (_BurnAuction *BurnAuctionCallerSession) DELTABLOCKSINITIALSLOT() (uint32, error) {
	return _BurnAuction.Contract.DELTABLOCKSINITIALSLOT(&_BurnAuction.CallOpts)
}

// DONATIONDENOMINATOR is a free data retrieval call binding the contract method 0x9a478d24.
//
// Solidity: function DONATION_DENOMINATOR() view returns(uint256)
func (_BurnAuction *BurnAuctionCaller) DONATIONDENOMINATOR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BurnAuction.contract.Call(opts, &out, "DONATION_DENOMINATOR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DONATIONDENOMINATOR is a free data retrieval call binding the contract method 0x9a478d24.
//
// Solidity: function DONATION_DENOMINATOR() view returns(uint256)
func (_BurnAuction *BurnAuctionSession) DONATIONDENOMINATOR() (*big.Int, error) {
	return _BurnAuction.Contract.DONATIONDENOMINATOR(&_BurnAuction.CallOpts)
}

// DONATIONDENOMINATOR is a free data retrieval call binding the contract method 0x9a478d24.
//
// Solidity: function DONATION_DENOMINATOR() view returns(uint256)
func (_BurnAuction *BurnAuctionCallerSession) DONATIONDENOMINATOR() (*big.Int, error) {
	return _BurnAuction.Contract.DONATIONDENOMINATOR(&_BurnAuction.CallOpts)
}

// Auction is a free data retrieval call binding the contract method 0xd3dd08e2.
//
// Solidity: function auction(uint32 ) view returns(address coordinator, uint128 amount, bool initialized)
func (_BurnAuction *BurnAuctionCaller) Auction(opts *bind.CallOpts, arg0 uint32) (struct {
	Coordinator common.Address
	Amount      *big.Int
	Initialized bool
}, error) {
	var out []interface{}
	err := _BurnAuction.contract.Call(opts, &out, "auction", arg0)

	outstruct := new(struct {
		Coordinator common.Address
		Amount      *big.Int
		Initialized bool
	})

	outstruct.Coordinator = out[0].(common.Address)
	outstruct.Amount = out[1].(*big.Int)
	outstruct.Initialized = out[2].(bool)

	return *outstruct, err

}

// Auction is a free data retrieval call binding the contract method 0xd3dd08e2.
//
// Solidity: function auction(uint32 ) view returns(address coordinator, uint128 amount, bool initialized)
func (_BurnAuction *BurnAuctionSession) Auction(arg0 uint32) (struct {
	Coordinator common.Address
	Amount      *big.Int
	Initialized bool
}, error) {
	return _BurnAuction.Contract.Auction(&_BurnAuction.CallOpts, arg0)
}

// Auction is a free data retrieval call binding the contract method 0xd3dd08e2.
//
// Solidity: function auction(uint32 ) view returns(address coordinator, uint128 amount, bool initialized)
func (_BurnAuction *BurnAuctionCallerSession) Auction(arg0 uint32) (struct {
	Coordinator common.Address
	Amount      *big.Int
	Initialized bool
}, error) {
	return _BurnAuction.Contract.Auction(&_BurnAuction.CallOpts, arg0)
}

// Block2slot is a free data retrieval call binding the contract method 0xa87a2ead.
//
// Solidity: function block2slot(uint256 numBlock) view returns(uint32)
func (_BurnAuction *BurnAuctionCaller) Block2slot(opts *bind.CallOpts, numBlock *big.Int) (uint32, error) {
	var out []interface{}
	err := _BurnAuction.contract.Call(opts, &out, "block2slot", numBlock)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// Block2slot is a free data retrieval call binding the contract method 0xa87a2ead.
//
// Solidity: function block2slot(uint256 numBlock) view returns(uint32)
func (_BurnAuction *BurnAuctionSession) Block2slot(numBlock *big.Int) (uint32, error) {
	return _BurnAuction.Contract.Block2slot(&_BurnAuction.CallOpts, numBlock)
}

// Block2slot is a free data retrieval call binding the contract method 0xa87a2ead.
//
// Solidity: function block2slot(uint256 numBlock) view returns(uint32)
func (_BurnAuction *BurnAuctionCallerSession) Block2slot(numBlock *big.Int) (uint32, error) {
	return _BurnAuction.Contract.Block2slot(&_BurnAuction.CallOpts, numBlock)
}

// CurrentSlot is a free data retrieval call binding the contract method 0x3359632e.
//
// Solidity: function currentSlot() view returns(uint32)
func (_BurnAuction *BurnAuctionCaller) CurrentSlot(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _BurnAuction.contract.Call(opts, &out, "currentSlot")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// CurrentSlot is a free data retrieval call binding the contract method 0x3359632e.
//
// Solidity: function currentSlot() view returns(uint32)
func (_BurnAuction *BurnAuctionSession) CurrentSlot() (uint32, error) {
	return _BurnAuction.Contract.CurrentSlot(&_BurnAuction.CallOpts)
}

// CurrentSlot is a free data retrieval call binding the contract method 0x3359632e.
//
// Solidity: function currentSlot() view returns(uint32)
func (_BurnAuction *BurnAuctionCallerSession) CurrentSlot() (uint32, error) {
	return _BurnAuction.Contract.CurrentSlot(&_BurnAuction.CallOpts)
}

// Deposits is a free data retrieval call binding the contract method 0xfc7e286d.
//
// Solidity: function deposits(address ) view returns(uint256)
func (_BurnAuction *BurnAuctionCaller) Deposits(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _BurnAuction.contract.Call(opts, &out, "deposits", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Deposits is a free data retrieval call binding the contract method 0xfc7e286d.
//
// Solidity: function deposits(address ) view returns(uint256)
func (_BurnAuction *BurnAuctionSession) Deposits(arg0 common.Address) (*big.Int, error) {
	return _BurnAuction.Contract.Deposits(&_BurnAuction.CallOpts, arg0)
}

// Deposits is a free data retrieval call binding the contract method 0xfc7e286d.
//
// Solidity: function deposits(address ) view returns(uint256)
func (_BurnAuction *BurnAuctionCallerSession) Deposits(arg0 common.Address) (*big.Int, error) {
	return _BurnAuction.Contract.Deposits(&_BurnAuction.CallOpts, arg0)
}

// DonationAddress is a free data retrieval call binding the contract method 0xec034bed.
//
// Solidity: function donationAddress() view returns(address)
func (_BurnAuction *BurnAuctionCaller) DonationAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnAuction.contract.Call(opts, &out, "donationAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DonationAddress is a free data retrieval call binding the contract method 0xec034bed.
//
// Solidity: function donationAddress() view returns(address)
func (_BurnAuction *BurnAuctionSession) DonationAddress() (common.Address, error) {
	return _BurnAuction.Contract.DonationAddress(&_BurnAuction.CallOpts)
}

// DonationAddress is a free data retrieval call binding the contract method 0xec034bed.
//
// Solidity: function donationAddress() view returns(address)
func (_BurnAuction *BurnAuctionCallerSession) DonationAddress() (common.Address, error) {
	return _BurnAuction.Contract.DonationAddress(&_BurnAuction.CallOpts)
}

// DonationNumerator is a free data retrieval call binding the contract method 0xf2182cf9.
//
// Solidity: function donationNumerator() view returns(uint256)
func (_BurnAuction *BurnAuctionCaller) DonationNumerator(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BurnAuction.contract.Call(opts, &out, "donationNumerator")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DonationNumerator is a free data retrieval call binding the contract method 0xf2182cf9.
//
// Solidity: function donationNumerator() view returns(uint256)
func (_BurnAuction *BurnAuctionSession) DonationNumerator() (*big.Int, error) {
	return _BurnAuction.Contract.DonationNumerator(&_BurnAuction.CallOpts)
}

// DonationNumerator is a free data retrieval call binding the contract method 0xf2182cf9.
//
// Solidity: function donationNumerator() view returns(uint256)
func (_BurnAuction *BurnAuctionCallerSession) DonationNumerator() (*big.Int, error) {
	return _BurnAuction.Contract.DonationNumerator(&_BurnAuction.CallOpts)
}

// GenesisBlock is a free data retrieval call binding the contract method 0x4cdc9c63.
//
// Solidity: function genesisBlock() view returns(uint256)
func (_BurnAuction *BurnAuctionCaller) GenesisBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BurnAuction.contract.Call(opts, &out, "genesisBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GenesisBlock is a free data retrieval call binding the contract method 0x4cdc9c63.
//
// Solidity: function genesisBlock() view returns(uint256)
func (_BurnAuction *BurnAuctionSession) GenesisBlock() (*big.Int, error) {
	return _BurnAuction.Contract.GenesisBlock(&_BurnAuction.CallOpts)
}

// GenesisBlock is a free data retrieval call binding the contract method 0x4cdc9c63.
//
// Solidity: function genesisBlock() view returns(uint256)
func (_BurnAuction *BurnAuctionCallerSession) GenesisBlock() (*big.Int, error) {
	return _BurnAuction.Contract.GenesisBlock(&_BurnAuction.CallOpts)
}

// GetBlockNumber is a free data retrieval call binding the contract method 0x42cbb15c.
//
// Solidity: function getBlockNumber() view returns(uint256)
func (_BurnAuction *BurnAuctionCaller) GetBlockNumber(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _BurnAuction.contract.Call(opts, &out, "getBlockNumber")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBlockNumber is a free data retrieval call binding the contract method 0x42cbb15c.
//
// Solidity: function getBlockNumber() view returns(uint256)
func (_BurnAuction *BurnAuctionSession) GetBlockNumber() (*big.Int, error) {
	return _BurnAuction.Contract.GetBlockNumber(&_BurnAuction.CallOpts)
}

// GetBlockNumber is a free data retrieval call binding the contract method 0x42cbb15c.
//
// Solidity: function getBlockNumber() view returns(uint256)
func (_BurnAuction *BurnAuctionCallerSession) GetBlockNumber() (*big.Int, error) {
	return _BurnAuction.Contract.GetBlockNumber(&_BurnAuction.CallOpts)
}

// GetProposer is a free data retrieval call binding the contract method 0xe9790d02.
//
// Solidity: function getProposer() view returns(address)
func (_BurnAuction *BurnAuctionCaller) GetProposer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BurnAuction.contract.Call(opts, &out, "getProposer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetProposer is a free data retrieval call binding the contract method 0xe9790d02.
//
// Solidity: function getProposer() view returns(address)
func (_BurnAuction *BurnAuctionSession) GetProposer() (common.Address, error) {
	return _BurnAuction.Contract.GetProposer(&_BurnAuction.CallOpts)
}

// GetProposer is a free data retrieval call binding the contract method 0xe9790d02.
//
// Solidity: function getProposer() view returns(address)
func (_BurnAuction *BurnAuctionCallerSession) GetProposer() (common.Address, error) {
	return _BurnAuction.Contract.GetProposer(&_BurnAuction.CallOpts)
}

// Bid is a paid mutator transaction binding the contract method 0x454a2ab3.
//
// Solidity: function bid(uint256 bidAmount) payable returns()
func (_BurnAuction *BurnAuctionTransactor) Bid(opts *bind.TransactOpts, bidAmount *big.Int) (*types.Transaction, error) {
	return _BurnAuction.contract.Transact(opts, "bid", bidAmount)
}

// Bid is a paid mutator transaction binding the contract method 0x454a2ab3.
//
// Solidity: function bid(uint256 bidAmount) payable returns()
func (_BurnAuction *BurnAuctionSession) Bid(bidAmount *big.Int) (*types.Transaction, error) {
	return _BurnAuction.Contract.Bid(&_BurnAuction.TransactOpts, bidAmount)
}

// Bid is a paid mutator transaction binding the contract method 0x454a2ab3.
//
// Solidity: function bid(uint256 bidAmount) payable returns()
func (_BurnAuction *BurnAuctionTransactorSession) Bid(bidAmount *big.Int) (*types.Transaction, error) {
	return _BurnAuction.Contract.Bid(&_BurnAuction.TransactOpts, bidAmount)
}

// Witdraw is a paid mutator transaction binding the contract method 0xb404930e.
//
// Solidity: function witdraw(uint256 amount) returns()
func (_BurnAuction *BurnAuctionTransactor) Witdraw(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _BurnAuction.contract.Transact(opts, "witdraw", amount)
}

// Witdraw is a paid mutator transaction binding the contract method 0xb404930e.
//
// Solidity: function witdraw(uint256 amount) returns()
func (_BurnAuction *BurnAuctionSession) Witdraw(amount *big.Int) (*types.Transaction, error) {
	return _BurnAuction.Contract.Witdraw(&_BurnAuction.TransactOpts, amount)
}

// Witdraw is a paid mutator transaction binding the contract method 0xb404930e.
//
// Solidity: function witdraw(uint256 amount) returns()
func (_BurnAuction *BurnAuctionTransactorSession) Witdraw(amount *big.Int) (*types.Transaction, error) {
	return _BurnAuction.Contract.Witdraw(&_BurnAuction.TransactOpts, amount)
}

// WithdrawDonation is a paid mutator transaction binding the contract method 0xe81e1ccc.
//
// Solidity: function withdrawDonation() returns()
func (_BurnAuction *BurnAuctionTransactor) WithdrawDonation(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BurnAuction.contract.Transact(opts, "withdrawDonation")
}

// WithdrawDonation is a paid mutator transaction binding the contract method 0xe81e1ccc.
//
// Solidity: function withdrawDonation() returns()
func (_BurnAuction *BurnAuctionSession) WithdrawDonation() (*types.Transaction, error) {
	return _BurnAuction.Contract.WithdrawDonation(&_BurnAuction.TransactOpts)
}

// WithdrawDonation is a paid mutator transaction binding the contract method 0xe81e1ccc.
//
// Solidity: function withdrawDonation() returns()
func (_BurnAuction *BurnAuctionTransactorSession) WithdrawDonation() (*types.Transaction, error) {
	return _BurnAuction.Contract.WithdrawDonation(&_BurnAuction.TransactOpts)
}

// BurnAuctionNewBestBidIterator is returned from FilterNewBestBid and is used to iterate over the raw logs and unpacked data for NewBestBid events raised by the BurnAuction contract.
type BurnAuctionNewBestBidIterator struct {
	Event *BurnAuctionNewBestBid // Event containing the contract specifics and raw log

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
func (it *BurnAuctionNewBestBidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BurnAuctionNewBestBid)
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
		it.Event = new(BurnAuctionNewBestBid)
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
func (it *BurnAuctionNewBestBidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BurnAuctionNewBestBidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BurnAuctionNewBestBid represents a NewBestBid event raised by the BurnAuction contract.
type BurnAuctionNewBestBid struct {
	Slot        uint32
	Coordinator common.Address
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterNewBestBid is a free log retrieval operation binding the contract event 0x304f693446955254ce103ccf22f2ee397d8c2517f63076ee63e0dfa22fc5ba55.
//
// Solidity: event NewBestBid(uint32 slot, address coordinator, uint256 amount)
func (_BurnAuction *BurnAuctionFilterer) FilterNewBestBid(opts *bind.FilterOpts) (*BurnAuctionNewBestBidIterator, error) {

	logs, sub, err := _BurnAuction.contract.FilterLogs(opts, "NewBestBid")
	if err != nil {
		return nil, err
	}
	return &BurnAuctionNewBestBidIterator{contract: _BurnAuction.contract, event: "NewBestBid", logs: logs, sub: sub}, nil
}

// WatchNewBestBid is a free log subscription operation binding the contract event 0x304f693446955254ce103ccf22f2ee397d8c2517f63076ee63e0dfa22fc5ba55.
//
// Solidity: event NewBestBid(uint32 slot, address coordinator, uint256 amount)
func (_BurnAuction *BurnAuctionFilterer) WatchNewBestBid(opts *bind.WatchOpts, sink chan<- *BurnAuctionNewBestBid) (event.Subscription, error) {

	logs, sub, err := _BurnAuction.contract.WatchLogs(opts, "NewBestBid")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BurnAuctionNewBestBid)
				if err := _BurnAuction.contract.UnpackLog(event, "NewBestBid", log); err != nil {
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

// ParseNewBestBid is a log parse operation binding the contract event 0x304f693446955254ce103ccf22f2ee397d8c2517f63076ee63e0dfa22fc5ba55.
//
// Solidity: event NewBestBid(uint32 slot, address coordinator, uint256 amount)
func (_BurnAuction *BurnAuctionFilterer) ParseNewBestBid(log types.Log) (*BurnAuctionNewBestBid, error) {
	event := new(BurnAuctionNewBestBid)
	if err := _BurnAuction.contract.UnpackLog(event, "NewBestBid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
