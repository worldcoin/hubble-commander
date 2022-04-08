// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package accountregistry

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

// AccountRegistryMetaData contains all meta data concerning the AccountRegistry contract.
var AccountRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractChooser\",\"name\":\"_chooser\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"rootLeft\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"leafIndexLeft\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[31]\",\"name\":\"filledSubtreesLeft\",\"type\":\"bytes32[31]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"endID\",\"type\":\"uint256\"}],\"name\":\"BatchPubkeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"}],\"name\":\"SinglePubkeyRegistered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BATCH_DEPTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"BATCH_SIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEPTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SET_SIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WITNESS_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"chooser\",\"outputs\":[{\"internalType\":\"contractChooser\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256[4]\",\"name\":\"pubkey\",\"type\":\"uint256[4]\"},{\"internalType\":\"bytes32[31]\",\"name\":\"witness\",\"type\":\"bytes32[31]\"}],\"name\":\"exists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"filledSubtreesLeft\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"filledSubtreesRight\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"leafIndexLeft\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"leafIndexRight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[4]\",\"name\":\"pubkey\",\"type\":\"uint256[4]\"}],\"name\":\"register\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[4][16]\",\"name\":\"pubkeys\",\"type\":\"uint256[4][16]\"}],\"name\":\"registerBatch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"root\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rootLeft\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rootRight\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"zeros\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60a0604052600060035560006004553480156200001b57600080fd5b5060405162001259380380620012598339810160408190526200003e91620002f7565b8282827f290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e5638060056000015560015b601f811015620001a457600562000085600183620003bb565b601f8110620000a457634e487b7160e01b600052603260045260246000fd5b01546005620000b5600184620003bb565b601f8110620000d457634e487b7160e01b600052603260045260246000fd5b015460408051602081019390935282015260600160405160208183030381529060405280519060200120600582601f81106200012057634e487b7160e01b600052603260045260246000fd5b01556004811080159062000134575080601f115b156200018f57600581601f81106200015c57634e487b7160e01b600052603260045260246000fd5b015460436200016d600484620003bb565b601b81106200018c57634e487b7160e01b600052603260045260246000fd5b01555b806200019b81620003d5565b9150506200006c565b50620001b4602483601f6200029d565b50600383905560008490556005620001cf6001601f620003bb565b601f8110620001ee57634e487b7160e01b600052603260045260246000fd5b01546005620002006001601f620003bb565b601f81106200021f57634e487b7160e01b600052603260045260246000fd5b015460408051602081019390935282015260600160408051808303601f1901815282825280516020918201206001819055600054918401919091529082015260600160408051808303601f190181529190528051602090910120600255505050505060609290921b6001600160601b031916608052506200041f9050565b82601f8101928215620002ce579160200282015b82811115620002ce578251825591602001919060010190620002b1565b50620002dc929150620002e0565b5090565b5b80821115620002dc5760008155600101620002e1565b6000806000806104408086880312156200030f578485fd5b85516001600160a01b038116811462000326578586fd5b8095505060208087015194506040870151935087607f88011262000348578283fd5b6040516103e081016001600160401b03811182821017156200036e576200036e62000409565b6040528060608901848a018b101562000385578586fd5b8594505b601f851015620003aa57805182526001949094019390830190830162000389565b505080935050505092959194509250565b600082821015620003d057620003d0620003f3565b500390565b6000600019821415620003ec57620003ec620003f3565b5060010190565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052604160045260246000fd5b60805160601c610e1462000445600039600081816101b8015261024f0152610e146000f3fe608060405234801561001057600080fd5b506004361061010b5760003560e01c806395e4bf03116100a2578063d0383d6811610071578063d0383d68146101fb578063d7c53ea714610206578063d82894631461020f578063e829558814610218578063ebf0c7171461022b57600080fd5b806395e4bf03146101a057806398366e351461015957806398d17621146101b3578063cab2da9b146101f257600080fd5b80635e71468b116100de5780635e71468b14610159578063693c1db714610161578063709a8b2a1461016a5780638d0379621461018d57600080fd5b8063034a29ae146101105780631c4a7a94146101365780631c76e77e1461014957806349faa4d414610151575b600080fd5b61012361011e366004610cf7565b610234565b6040519081526020015b60405180910390f35b610123610144366004610cb3565b61024b565b610123600481565b610123601081565b610123601f81565b61012360035481565b61017d610178366004610d0f565b610451565b604051901515815260200161012d565b61012361019b366004610cf7565b6104bc565b6101236101ae366004610cdc565b6104cc565b6101da7f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b03909116815260200161012d565b61012360015481565b610123638000000081565b61012360045481565b61012360005481565b610123610226366004610cf7565b610536565b61012360025481565b602481601f811061024457600080fd5b0154905081565b60007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b1580156102a657600080fd5b505afa1580156102ba573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102de9190610c8c565b6001600160a01b0316336001600160a01b03161461034f5760405162461bcd60e51b8152602060048201526024808201527f424c534163636f756e7452656769737472793a20496e76616c69642070726f7060448201526337b9b2b960e11b60648201526084015b60405180910390fd5b610357610c37565b60005b60108110156103ea57600084826010811061038557634e487b7160e01b600052603260045260246000fd5b608002016040516020016103999190610d54565b604051602081830303815290604052805190602001209050808383601081106103d257634e487b7160e01b600052603260045260246000fd5b602002015250806103e281610dac565b91505061035a565b5060006103f682610546565b90507f3154b80a7d9f6a143c37dde575f47deb78dacc7f280d8efc7e3ae102758a841b816001610427601083610d69565b6104319190610d95565b604080519283526020830191909152015b60405180910390a19392505050565b600080836040516020016104659190610d54565b6040516020818303038152906040528051906020012090506104b1818685601f806020026040519081016040528092919082601f6020028082843760009201919091525061093d915050565b9150505b9392505050565b604381601b811061024457600080fd5b600080826040516020016104e09190610d54565b604051602081830303815290604052805190602001209050600061050382610a5f565b90507f59056afed767866d7c194cf26e24ebe16974ef943cccb729452e9adc265a9ac78160405161044291815260200190565b600581601f811061024457600080fd5b6000601061055960016380000000610d95565b6105639190610d95565b600454106105b35760405162461bcd60e51b815260206004820152601f60248201527f4163636f756e74547265653a207269676874207365742069732066756c6c20006044820152606401610346565b6105bb610c56565b60005b6105ca60016004610d95565b6001901b81101561069657600181901b8481601081106105fa57634e487b7160e01b600052603260045260246000fd5b60200201518561060b836001610d69565b6010811061062957634e487b7160e01b600052603260045260246000fd5b6020020151604051602001610648929190918252602082015260400190565b6040516020818303038152906040528051906020012083836008811061067e57634e487b7160e01b600052603260045260246000fd5b6020020152508061068e81610dac565b9150506105be565b5060015b60048110156107a357600060016106b2836004610d95565b6106bc9190610d95565b6001901b905060005b8181101561078e57600181901b8481600881106106f257634e487b7160e01b600052603260045260246000fd5b602002015185610703836001610d69565b6008811061072157634e487b7160e01b600052603260045260246000fd5b6020020151604051602001610740929190918252602082015260400190565b6040516020818303038152906040528051906020012085836008811061077657634e487b7160e01b600052603260045260246000fd5b6020020152508061078681610dac565b9150506106c5565b5050808061079b90610dac565b91505061069a565b5080516004546000906107b890601090610d81565b90506000805b6107ca6004601f610d95565b8110156108d157826001166001141561083157604381601b81106107fe57634e487b7160e01b600052603260045260246000fd5b015460408051602081019290925281018590526060016040516020818303038152906040528051906020012093506108b8565b8161085f5783604382601b811061085857634e487b7160e01b600052603260045260246000fd5b0155600191505b83600561086d600484610d69565b601f811061088b57634e487b7160e01b600052603260045260246000fd5b01546040805160208101939093528201526060016040516020818303038152906040528051906020012093505b60019290921c91806108c981610dac565b9150506107be565b506001839055600054604080516020810192909252810184905260600160408051601f198184030181529190528051602090910120600255600480546010919060009061091f908490610d69565b909155505060045461093390601090610d95565b9695505050505050565b60008061094e638000000085610dc7565b90508460005b601f811015610a385782600116600114156109c6578481601f811061098957634e487b7160e01b600052603260045260246000fd5b6020020151826040516020016109a9929190918252602082015260400190565b604051602081830303815290604052805190602001209150610a1f565b818582601f81106109e757634e487b7160e01b600052603260045260246000fd5b6020020151604051602001610a06929190918252602082015260400190565b6040516020818303038152906040528051906020012091505b60019290921c9180610a3081610dac565b915050610954565b506380000000851015610a52576000541491506104b59050565b6001541491506104b59050565b6000610a7060016380000000610d95565b60035410610ac05760405162461bcd60e51b815260206004820152601e60248201527f4163636f756e74547265653a206c656674207365742069732066756c6c2000006044820152606401610346565b60035482906000805b601f811015610bc8578260011660011415610b3257602481601f8110610aff57634e487b7160e01b600052603260045260246000fd5b01546040805160208101929092528101859052606001604051602081830303815290604052805190602001209350610baf565b81610b605783602482601f8110610b5957634e487b7160e01b600052603260045260246000fd5b0155600191505b83600582601f8110610b8257634e487b7160e01b600052603260045260246000fd5b01546040805160208101939093528201526060016040516020818303038152906040528051906020012093505b60019290921c9180610bc081610dac565b915050610ac9565b506000839055600154604051610beb918591602001918252602082015260400190565b60405160208183030381529060405280519060200120600281905550600160036000828254610c1a9190610d69565b9091555050600354610c2e90600190610d95565b95945050505050565b6040518061020001604052806010906020820280368337509192915050565b6040518061010001604052806008906020820280368337509192915050565b8060808101831015610c8657600080fd5b92915050565b600060208284031215610c9d578081fd5b81516001600160a01b03811681146104b5578182fd5b6000610800808385031215610cc6578182fd5b838184011115610cd4578182fd5b509092915050565b600060808284031215610ced578081fd5b6104b58383610c75565b600060208284031215610d08578081fd5b5035919050565b6000806000610480808587031215610d25578283fd5b84359350610d368660208701610c75565b9250858186011115610d46578182fd5b5060a0840190509250925092565b60808282376000608091909101908152919050565b60008219821115610d7c57610d7c610ddb565b500190565b600082610d9057610d90610df1565b500490565b600082821015610da757610da7610ddb565b500390565b6000600019821415610dc057610dc0610ddb565b5060010190565b600082610dd657610dd6610df1565b500690565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052601260045260246000fdfea164736f6c6343000804000a",
}

// AccountRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use AccountRegistryMetaData.ABI instead.
var AccountRegistryABI = AccountRegistryMetaData.ABI

// AccountRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use AccountRegistryMetaData.Bin instead.
var AccountRegistryBin = AccountRegistryMetaData.Bin

// DeployAccountRegistry deploys a new Ethereum contract, binding an instance of AccountRegistry to it.
func DeployAccountRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, _chooser common.Address, rootLeft [32]byte, leafIndexLeft *big.Int, filledSubtreesLeft [31][32]byte) (common.Address, *types.Transaction, *AccountRegistry, error) {
	parsed, err := AccountRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(AccountRegistryBin), backend, _chooser, rootLeft, leafIndexLeft, filledSubtreesLeft)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AccountRegistry{AccountRegistryCaller: AccountRegistryCaller{contract: contract}, AccountRegistryTransactor: AccountRegistryTransactor{contract: contract}, AccountRegistryFilterer: AccountRegistryFilterer{contract: contract}}, nil
}

// AccountRegistry is an auto generated Go binding around an Ethereum contract.
type AccountRegistry struct {
	AccountRegistryCaller     // Read-only binding to the contract
	AccountRegistryTransactor // Write-only binding to the contract
	AccountRegistryFilterer   // Log filterer for contract events
}

// AccountRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type AccountRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AccountRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AccountRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccountRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AccountRegistrySession struct {
	Contract     *AccountRegistry  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AccountRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AccountRegistryCallerSession struct {
	Contract *AccountRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// AccountRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AccountRegistryTransactorSession struct {
	Contract     *AccountRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// AccountRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type AccountRegistryRaw struct {
	Contract *AccountRegistry // Generic contract binding to access the raw methods on
}

// AccountRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AccountRegistryCallerRaw struct {
	Contract *AccountRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// AccountRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AccountRegistryTransactorRaw struct {
	Contract *AccountRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAccountRegistry creates a new instance of AccountRegistry, bound to a specific deployed contract.
func NewAccountRegistry(address common.Address, backend bind.ContractBackend) (*AccountRegistry, error) {
	contract, err := bindAccountRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AccountRegistry{AccountRegistryCaller: AccountRegistryCaller{contract: contract}, AccountRegistryTransactor: AccountRegistryTransactor{contract: contract}, AccountRegistryFilterer: AccountRegistryFilterer{contract: contract}}, nil
}

// NewAccountRegistryCaller creates a new read-only instance of AccountRegistry, bound to a specific deployed contract.
func NewAccountRegistryCaller(address common.Address, caller bind.ContractCaller) (*AccountRegistryCaller, error) {
	contract, err := bindAccountRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccountRegistryCaller{contract: contract}, nil
}

// NewAccountRegistryTransactor creates a new write-only instance of AccountRegistry, bound to a specific deployed contract.
func NewAccountRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*AccountRegistryTransactor, error) {
	contract, err := bindAccountRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccountRegistryTransactor{contract: contract}, nil
}

// NewAccountRegistryFilterer creates a new log filterer instance of AccountRegistry, bound to a specific deployed contract.
func NewAccountRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*AccountRegistryFilterer, error) {
	contract, err := bindAccountRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccountRegistryFilterer{contract: contract}, nil
}

// bindAccountRegistry binds a generic wrapper to an already deployed contract.
func bindAccountRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AccountRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccountRegistry *AccountRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccountRegistry.Contract.AccountRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccountRegistry *AccountRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccountRegistry.Contract.AccountRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccountRegistry *AccountRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccountRegistry.Contract.AccountRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccountRegistry *AccountRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccountRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccountRegistry *AccountRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccountRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccountRegistry *AccountRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccountRegistry.Contract.contract.Transact(opts, method, params...)
}

// BATCHDEPTH is a free data retrieval call binding the contract method 0x1c76e77e.
//
// Solidity: function BATCH_DEPTH() view returns(uint256)
func (_AccountRegistry *AccountRegistryCaller) BATCHDEPTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "BATCH_DEPTH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BATCHDEPTH is a free data retrieval call binding the contract method 0x1c76e77e.
//
// Solidity: function BATCH_DEPTH() view returns(uint256)
func (_AccountRegistry *AccountRegistrySession) BATCHDEPTH() (*big.Int, error) {
	return _AccountRegistry.Contract.BATCHDEPTH(&_AccountRegistry.CallOpts)
}

// BATCHDEPTH is a free data retrieval call binding the contract method 0x1c76e77e.
//
// Solidity: function BATCH_DEPTH() view returns(uint256)
func (_AccountRegistry *AccountRegistryCallerSession) BATCHDEPTH() (*big.Int, error) {
	return _AccountRegistry.Contract.BATCHDEPTH(&_AccountRegistry.CallOpts)
}

// BATCHSIZE is a free data retrieval call binding the contract method 0x49faa4d4.
//
// Solidity: function BATCH_SIZE() view returns(uint256)
func (_AccountRegistry *AccountRegistryCaller) BATCHSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "BATCH_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BATCHSIZE is a free data retrieval call binding the contract method 0x49faa4d4.
//
// Solidity: function BATCH_SIZE() view returns(uint256)
func (_AccountRegistry *AccountRegistrySession) BATCHSIZE() (*big.Int, error) {
	return _AccountRegistry.Contract.BATCHSIZE(&_AccountRegistry.CallOpts)
}

// BATCHSIZE is a free data retrieval call binding the contract method 0x49faa4d4.
//
// Solidity: function BATCH_SIZE() view returns(uint256)
func (_AccountRegistry *AccountRegistryCallerSession) BATCHSIZE() (*big.Int, error) {
	return _AccountRegistry.Contract.BATCHSIZE(&_AccountRegistry.CallOpts)
}

// DEPTH is a free data retrieval call binding the contract method 0x98366e35.
//
// Solidity: function DEPTH() view returns(uint256)
func (_AccountRegistry *AccountRegistryCaller) DEPTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "DEPTH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DEPTH is a free data retrieval call binding the contract method 0x98366e35.
//
// Solidity: function DEPTH() view returns(uint256)
func (_AccountRegistry *AccountRegistrySession) DEPTH() (*big.Int, error) {
	return _AccountRegistry.Contract.DEPTH(&_AccountRegistry.CallOpts)
}

// DEPTH is a free data retrieval call binding the contract method 0x98366e35.
//
// Solidity: function DEPTH() view returns(uint256)
func (_AccountRegistry *AccountRegistryCallerSession) DEPTH() (*big.Int, error) {
	return _AccountRegistry.Contract.DEPTH(&_AccountRegistry.CallOpts)
}

// SETSIZE is a free data retrieval call binding the contract method 0xd0383d68.
//
// Solidity: function SET_SIZE() view returns(uint256)
func (_AccountRegistry *AccountRegistryCaller) SETSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "SET_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SETSIZE is a free data retrieval call binding the contract method 0xd0383d68.
//
// Solidity: function SET_SIZE() view returns(uint256)
func (_AccountRegistry *AccountRegistrySession) SETSIZE() (*big.Int, error) {
	return _AccountRegistry.Contract.SETSIZE(&_AccountRegistry.CallOpts)
}

// SETSIZE is a free data retrieval call binding the contract method 0xd0383d68.
//
// Solidity: function SET_SIZE() view returns(uint256)
func (_AccountRegistry *AccountRegistryCallerSession) SETSIZE() (*big.Int, error) {
	return _AccountRegistry.Contract.SETSIZE(&_AccountRegistry.CallOpts)
}

// WITNESSLENGTH is a free data retrieval call binding the contract method 0x5e71468b.
//
// Solidity: function WITNESS_LENGTH() view returns(uint256)
func (_AccountRegistry *AccountRegistryCaller) WITNESSLENGTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "WITNESS_LENGTH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WITNESSLENGTH is a free data retrieval call binding the contract method 0x5e71468b.
//
// Solidity: function WITNESS_LENGTH() view returns(uint256)
func (_AccountRegistry *AccountRegistrySession) WITNESSLENGTH() (*big.Int, error) {
	return _AccountRegistry.Contract.WITNESSLENGTH(&_AccountRegistry.CallOpts)
}

// WITNESSLENGTH is a free data retrieval call binding the contract method 0x5e71468b.
//
// Solidity: function WITNESS_LENGTH() view returns(uint256)
func (_AccountRegistry *AccountRegistryCallerSession) WITNESSLENGTH() (*big.Int, error) {
	return _AccountRegistry.Contract.WITNESSLENGTH(&_AccountRegistry.CallOpts)
}

// Chooser is a free data retrieval call binding the contract method 0x98d17621.
//
// Solidity: function chooser() view returns(address)
func (_AccountRegistry *AccountRegistryCaller) Chooser(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "chooser")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Chooser is a free data retrieval call binding the contract method 0x98d17621.
//
// Solidity: function chooser() view returns(address)
func (_AccountRegistry *AccountRegistrySession) Chooser() (common.Address, error) {
	return _AccountRegistry.Contract.Chooser(&_AccountRegistry.CallOpts)
}

// Chooser is a free data retrieval call binding the contract method 0x98d17621.
//
// Solidity: function chooser() view returns(address)
func (_AccountRegistry *AccountRegistryCallerSession) Chooser() (common.Address, error) {
	return _AccountRegistry.Contract.Chooser(&_AccountRegistry.CallOpts)
}

// Exists is a free data retrieval call binding the contract method 0x709a8b2a.
//
// Solidity: function exists(uint256 pubkeyID, uint256[4] pubkey, bytes32[31] witness) view returns(bool)
func (_AccountRegistry *AccountRegistryCaller) Exists(opts *bind.CallOpts, pubkeyID *big.Int, pubkey [4]*big.Int, witness [31][32]byte) (bool, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "exists", pubkeyID, pubkey, witness)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Exists is a free data retrieval call binding the contract method 0x709a8b2a.
//
// Solidity: function exists(uint256 pubkeyID, uint256[4] pubkey, bytes32[31] witness) view returns(bool)
func (_AccountRegistry *AccountRegistrySession) Exists(pubkeyID *big.Int, pubkey [4]*big.Int, witness [31][32]byte) (bool, error) {
	return _AccountRegistry.Contract.Exists(&_AccountRegistry.CallOpts, pubkeyID, pubkey, witness)
}

// Exists is a free data retrieval call binding the contract method 0x709a8b2a.
//
// Solidity: function exists(uint256 pubkeyID, uint256[4] pubkey, bytes32[31] witness) view returns(bool)
func (_AccountRegistry *AccountRegistryCallerSession) Exists(pubkeyID *big.Int, pubkey [4]*big.Int, witness [31][32]byte) (bool, error) {
	return _AccountRegistry.Contract.Exists(&_AccountRegistry.CallOpts, pubkeyID, pubkey, witness)
}

// FilledSubtreesLeft is a free data retrieval call binding the contract method 0x034a29ae.
//
// Solidity: function filledSubtreesLeft(uint256 ) view returns(bytes32)
func (_AccountRegistry *AccountRegistryCaller) FilledSubtreesLeft(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "filledSubtreesLeft", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FilledSubtreesLeft is a free data retrieval call binding the contract method 0x034a29ae.
//
// Solidity: function filledSubtreesLeft(uint256 ) view returns(bytes32)
func (_AccountRegistry *AccountRegistrySession) FilledSubtreesLeft(arg0 *big.Int) ([32]byte, error) {
	return _AccountRegistry.Contract.FilledSubtreesLeft(&_AccountRegistry.CallOpts, arg0)
}

// FilledSubtreesLeft is a free data retrieval call binding the contract method 0x034a29ae.
//
// Solidity: function filledSubtreesLeft(uint256 ) view returns(bytes32)
func (_AccountRegistry *AccountRegistryCallerSession) FilledSubtreesLeft(arg0 *big.Int) ([32]byte, error) {
	return _AccountRegistry.Contract.FilledSubtreesLeft(&_AccountRegistry.CallOpts, arg0)
}

// FilledSubtreesRight is a free data retrieval call binding the contract method 0x8d037962.
//
// Solidity: function filledSubtreesRight(uint256 ) view returns(bytes32)
func (_AccountRegistry *AccountRegistryCaller) FilledSubtreesRight(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "filledSubtreesRight", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FilledSubtreesRight is a free data retrieval call binding the contract method 0x8d037962.
//
// Solidity: function filledSubtreesRight(uint256 ) view returns(bytes32)
func (_AccountRegistry *AccountRegistrySession) FilledSubtreesRight(arg0 *big.Int) ([32]byte, error) {
	return _AccountRegistry.Contract.FilledSubtreesRight(&_AccountRegistry.CallOpts, arg0)
}

// FilledSubtreesRight is a free data retrieval call binding the contract method 0x8d037962.
//
// Solidity: function filledSubtreesRight(uint256 ) view returns(bytes32)
func (_AccountRegistry *AccountRegistryCallerSession) FilledSubtreesRight(arg0 *big.Int) ([32]byte, error) {
	return _AccountRegistry.Contract.FilledSubtreesRight(&_AccountRegistry.CallOpts, arg0)
}

// LeafIndexLeft is a free data retrieval call binding the contract method 0x693c1db7.
//
// Solidity: function leafIndexLeft() view returns(uint256)
func (_AccountRegistry *AccountRegistryCaller) LeafIndexLeft(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "leafIndexLeft")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LeafIndexLeft is a free data retrieval call binding the contract method 0x693c1db7.
//
// Solidity: function leafIndexLeft() view returns(uint256)
func (_AccountRegistry *AccountRegistrySession) LeafIndexLeft() (*big.Int, error) {
	return _AccountRegistry.Contract.LeafIndexLeft(&_AccountRegistry.CallOpts)
}

// LeafIndexLeft is a free data retrieval call binding the contract method 0x693c1db7.
//
// Solidity: function leafIndexLeft() view returns(uint256)
func (_AccountRegistry *AccountRegistryCallerSession) LeafIndexLeft() (*big.Int, error) {
	return _AccountRegistry.Contract.LeafIndexLeft(&_AccountRegistry.CallOpts)
}

// LeafIndexRight is a free data retrieval call binding the contract method 0xd7c53ea7.
//
// Solidity: function leafIndexRight() view returns(uint256)
func (_AccountRegistry *AccountRegistryCaller) LeafIndexRight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "leafIndexRight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LeafIndexRight is a free data retrieval call binding the contract method 0xd7c53ea7.
//
// Solidity: function leafIndexRight() view returns(uint256)
func (_AccountRegistry *AccountRegistrySession) LeafIndexRight() (*big.Int, error) {
	return _AccountRegistry.Contract.LeafIndexRight(&_AccountRegistry.CallOpts)
}

// LeafIndexRight is a free data retrieval call binding the contract method 0xd7c53ea7.
//
// Solidity: function leafIndexRight() view returns(uint256)
func (_AccountRegistry *AccountRegistryCallerSession) LeafIndexRight() (*big.Int, error) {
	return _AccountRegistry.Contract.LeafIndexRight(&_AccountRegistry.CallOpts)
}

// Root is a free data retrieval call binding the contract method 0xebf0c717.
//
// Solidity: function root() view returns(bytes32)
func (_AccountRegistry *AccountRegistryCaller) Root(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "root")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Root is a free data retrieval call binding the contract method 0xebf0c717.
//
// Solidity: function root() view returns(bytes32)
func (_AccountRegistry *AccountRegistrySession) Root() ([32]byte, error) {
	return _AccountRegistry.Contract.Root(&_AccountRegistry.CallOpts)
}

// Root is a free data retrieval call binding the contract method 0xebf0c717.
//
// Solidity: function root() view returns(bytes32)
func (_AccountRegistry *AccountRegistryCallerSession) Root() ([32]byte, error) {
	return _AccountRegistry.Contract.Root(&_AccountRegistry.CallOpts)
}

// RootLeft is a free data retrieval call binding the contract method 0xd8289463.
//
// Solidity: function rootLeft() view returns(bytes32)
func (_AccountRegistry *AccountRegistryCaller) RootLeft(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "rootLeft")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RootLeft is a free data retrieval call binding the contract method 0xd8289463.
//
// Solidity: function rootLeft() view returns(bytes32)
func (_AccountRegistry *AccountRegistrySession) RootLeft() ([32]byte, error) {
	return _AccountRegistry.Contract.RootLeft(&_AccountRegistry.CallOpts)
}

// RootLeft is a free data retrieval call binding the contract method 0xd8289463.
//
// Solidity: function rootLeft() view returns(bytes32)
func (_AccountRegistry *AccountRegistryCallerSession) RootLeft() ([32]byte, error) {
	return _AccountRegistry.Contract.RootLeft(&_AccountRegistry.CallOpts)
}

// RootRight is a free data retrieval call binding the contract method 0xcab2da9b.
//
// Solidity: function rootRight() view returns(bytes32)
func (_AccountRegistry *AccountRegistryCaller) RootRight(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "rootRight")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RootRight is a free data retrieval call binding the contract method 0xcab2da9b.
//
// Solidity: function rootRight() view returns(bytes32)
func (_AccountRegistry *AccountRegistrySession) RootRight() ([32]byte, error) {
	return _AccountRegistry.Contract.RootRight(&_AccountRegistry.CallOpts)
}

// RootRight is a free data retrieval call binding the contract method 0xcab2da9b.
//
// Solidity: function rootRight() view returns(bytes32)
func (_AccountRegistry *AccountRegistryCallerSession) RootRight() ([32]byte, error) {
	return _AccountRegistry.Contract.RootRight(&_AccountRegistry.CallOpts)
}

// Zeros is a free data retrieval call binding the contract method 0xe8295588.
//
// Solidity: function zeros(uint256 ) view returns(bytes32)
func (_AccountRegistry *AccountRegistryCaller) Zeros(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _AccountRegistry.contract.Call(opts, &out, "zeros", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Zeros is a free data retrieval call binding the contract method 0xe8295588.
//
// Solidity: function zeros(uint256 ) view returns(bytes32)
func (_AccountRegistry *AccountRegistrySession) Zeros(arg0 *big.Int) ([32]byte, error) {
	return _AccountRegistry.Contract.Zeros(&_AccountRegistry.CallOpts, arg0)
}

// Zeros is a free data retrieval call binding the contract method 0xe8295588.
//
// Solidity: function zeros(uint256 ) view returns(bytes32)
func (_AccountRegistry *AccountRegistryCallerSession) Zeros(arg0 *big.Int) ([32]byte, error) {
	return _AccountRegistry.Contract.Zeros(&_AccountRegistry.CallOpts, arg0)
}

// Register is a paid mutator transaction binding the contract method 0x95e4bf03.
//
// Solidity: function register(uint256[4] pubkey) returns(uint256)
func (_AccountRegistry *AccountRegistryTransactor) Register(opts *bind.TransactOpts, pubkey [4]*big.Int) (*types.Transaction, error) {
	return _AccountRegistry.contract.Transact(opts, "register", pubkey)
}

// Register is a paid mutator transaction binding the contract method 0x95e4bf03.
//
// Solidity: function register(uint256[4] pubkey) returns(uint256)
func (_AccountRegistry *AccountRegistrySession) Register(pubkey [4]*big.Int) (*types.Transaction, error) {
	return _AccountRegistry.Contract.Register(&_AccountRegistry.TransactOpts, pubkey)
}

// Register is a paid mutator transaction binding the contract method 0x95e4bf03.
//
// Solidity: function register(uint256[4] pubkey) returns(uint256)
func (_AccountRegistry *AccountRegistryTransactorSession) Register(pubkey [4]*big.Int) (*types.Transaction, error) {
	return _AccountRegistry.Contract.Register(&_AccountRegistry.TransactOpts, pubkey)
}

// RegisterBatch is a paid mutator transaction binding the contract method 0x1c4a7a94.
//
// Solidity: function registerBatch(uint256[4][16] pubkeys) returns(uint256)
func (_AccountRegistry *AccountRegistryTransactor) RegisterBatch(opts *bind.TransactOpts, pubkeys [16][4]*big.Int) (*types.Transaction, error) {
	return _AccountRegistry.contract.Transact(opts, "registerBatch", pubkeys)
}

// RegisterBatch is a paid mutator transaction binding the contract method 0x1c4a7a94.
//
// Solidity: function registerBatch(uint256[4][16] pubkeys) returns(uint256)
func (_AccountRegistry *AccountRegistrySession) RegisterBatch(pubkeys [16][4]*big.Int) (*types.Transaction, error) {
	return _AccountRegistry.Contract.RegisterBatch(&_AccountRegistry.TransactOpts, pubkeys)
}

// RegisterBatch is a paid mutator transaction binding the contract method 0x1c4a7a94.
//
// Solidity: function registerBatch(uint256[4][16] pubkeys) returns(uint256)
func (_AccountRegistry *AccountRegistryTransactorSession) RegisterBatch(pubkeys [16][4]*big.Int) (*types.Transaction, error) {
	return _AccountRegistry.Contract.RegisterBatch(&_AccountRegistry.TransactOpts, pubkeys)
}

// AccountRegistryBatchPubkeyRegisteredIterator is returned from FilterBatchPubkeyRegistered and is used to iterate over the raw logs and unpacked data for BatchPubkeyRegistered events raised by the AccountRegistry contract.
type AccountRegistryBatchPubkeyRegisteredIterator struct {
	Event *AccountRegistryBatchPubkeyRegistered // Event containing the contract specifics and raw log

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
func (it *AccountRegistryBatchPubkeyRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccountRegistryBatchPubkeyRegistered)
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
		it.Event = new(AccountRegistryBatchPubkeyRegistered)
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
func (it *AccountRegistryBatchPubkeyRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccountRegistryBatchPubkeyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccountRegistryBatchPubkeyRegistered represents a BatchPubkeyRegistered event raised by the AccountRegistry contract.
type AccountRegistryBatchPubkeyRegistered struct {
	StartID *big.Int
	EndID   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBatchPubkeyRegistered is a free log retrieval operation binding the contract event 0x3154b80a7d9f6a143c37dde575f47deb78dacc7f280d8efc7e3ae102758a841b.
//
// Solidity: event BatchPubkeyRegistered(uint256 startID, uint256 endID)
func (_AccountRegistry *AccountRegistryFilterer) FilterBatchPubkeyRegistered(opts *bind.FilterOpts) (*AccountRegistryBatchPubkeyRegisteredIterator, error) {

	logs, sub, err := _AccountRegistry.contract.FilterLogs(opts, "BatchPubkeyRegistered")
	if err != nil {
		return nil, err
	}
	return &AccountRegistryBatchPubkeyRegisteredIterator{contract: _AccountRegistry.contract, event: "BatchPubkeyRegistered", logs: logs, sub: sub}, nil
}

// WatchBatchPubkeyRegistered is a free log subscription operation binding the contract event 0x3154b80a7d9f6a143c37dde575f47deb78dacc7f280d8efc7e3ae102758a841b.
//
// Solidity: event BatchPubkeyRegistered(uint256 startID, uint256 endID)
func (_AccountRegistry *AccountRegistryFilterer) WatchBatchPubkeyRegistered(opts *bind.WatchOpts, sink chan<- *AccountRegistryBatchPubkeyRegistered) (event.Subscription, error) {

	logs, sub, err := _AccountRegistry.contract.WatchLogs(opts, "BatchPubkeyRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccountRegistryBatchPubkeyRegistered)
				if err := _AccountRegistry.contract.UnpackLog(event, "BatchPubkeyRegistered", log); err != nil {
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

// ParseBatchPubkeyRegistered is a log parse operation binding the contract event 0x3154b80a7d9f6a143c37dde575f47deb78dacc7f280d8efc7e3ae102758a841b.
//
// Solidity: event BatchPubkeyRegistered(uint256 startID, uint256 endID)
func (_AccountRegistry *AccountRegistryFilterer) ParseBatchPubkeyRegistered(log types.Log) (*AccountRegistryBatchPubkeyRegistered, error) {
	event := new(AccountRegistryBatchPubkeyRegistered)
	if err := _AccountRegistry.contract.UnpackLog(event, "BatchPubkeyRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AccountRegistrySinglePubkeyRegisteredIterator is returned from FilterSinglePubkeyRegistered and is used to iterate over the raw logs and unpacked data for SinglePubkeyRegistered events raised by the AccountRegistry contract.
type AccountRegistrySinglePubkeyRegisteredIterator struct {
	Event *AccountRegistrySinglePubkeyRegistered // Event containing the contract specifics and raw log

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
func (it *AccountRegistrySinglePubkeyRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccountRegistrySinglePubkeyRegistered)
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
		it.Event = new(AccountRegistrySinglePubkeyRegistered)
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
func (it *AccountRegistrySinglePubkeyRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccountRegistrySinglePubkeyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccountRegistrySinglePubkeyRegistered represents a SinglePubkeyRegistered event raised by the AccountRegistry contract.
type AccountRegistrySinglePubkeyRegistered struct {
	PubkeyID *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSinglePubkeyRegistered is a free log retrieval operation binding the contract event 0x59056afed767866d7c194cf26e24ebe16974ef943cccb729452e9adc265a9ac7.
//
// Solidity: event SinglePubkeyRegistered(uint256 pubkeyID)
func (_AccountRegistry *AccountRegistryFilterer) FilterSinglePubkeyRegistered(opts *bind.FilterOpts) (*AccountRegistrySinglePubkeyRegisteredIterator, error) {

	logs, sub, err := _AccountRegistry.contract.FilterLogs(opts, "SinglePubkeyRegistered")
	if err != nil {
		return nil, err
	}
	return &AccountRegistrySinglePubkeyRegisteredIterator{contract: _AccountRegistry.contract, event: "SinglePubkeyRegistered", logs: logs, sub: sub}, nil
}

// WatchSinglePubkeyRegistered is a free log subscription operation binding the contract event 0x59056afed767866d7c194cf26e24ebe16974ef943cccb729452e9adc265a9ac7.
//
// Solidity: event SinglePubkeyRegistered(uint256 pubkeyID)
func (_AccountRegistry *AccountRegistryFilterer) WatchSinglePubkeyRegistered(opts *bind.WatchOpts, sink chan<- *AccountRegistrySinglePubkeyRegistered) (event.Subscription, error) {

	logs, sub, err := _AccountRegistry.contract.WatchLogs(opts, "SinglePubkeyRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccountRegistrySinglePubkeyRegistered)
				if err := _AccountRegistry.contract.UnpackLog(event, "SinglePubkeyRegistered", log); err != nil {
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

// ParseSinglePubkeyRegistered is a log parse operation binding the contract event 0x59056afed767866d7c194cf26e24ebe16974ef943cccb729452e9adc265a9ac7.
//
// Solidity: event SinglePubkeyRegistered(uint256 pubkeyID)
func (_AccountRegistry *AccountRegistryFilterer) ParseSinglePubkeyRegistered(log types.Log) (*AccountRegistrySinglePubkeyRegistered, error) {
	event := new(AccountRegistrySinglePubkeyRegistered)
	if err := _AccountRegistry.contract.UnpackLog(event, "SinglePubkeyRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
