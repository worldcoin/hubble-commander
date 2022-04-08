// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package defaultaccountregistry

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

// DefaultAccountRegistryMetaData contains all meta data concerning the DefaultAccountRegistry contract.
var DefaultAccountRegistryMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractChooser\",\"name\":\"_chooser\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"endID\",\"type\":\"uint256\"}],\"name\":\"BatchPubkeyRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"}],\"name\":\"SinglePubkeyRegistered\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"BATCH_DEPTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"BATCH_SIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DEPTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIAL_LEFT_ROOT\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SET_SIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WITNESS_LENGTH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"chooser\",\"outputs\":[{\"internalType\":\"contractChooser\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256[4]\",\"name\":\"pubkey\",\"type\":\"uint256[4]\"},{\"internalType\":\"bytes32[31]\",\"name\":\"witness\",\"type\":\"bytes32[31]\"}],\"name\":\"exists\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"filledSubtreesLeft\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"filledSubtreesRight\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"initialSubtrees\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"leafIndexLeft\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"leafIndexRight\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[4]\",\"name\":\"pubkey\",\"type\":\"uint256[4]\"}],\"name\":\"register\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[4][16]\",\"name\":\"pubkeys\",\"type\":\"uint256[4][16]\"}],\"name\":\"registerBatch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"root\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rootLeft\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rootRight\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"zeros\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x610480604052600080516020620016ca83398151915260a09081527f633dc4d7da7256660a892f8f1604a44b5432649cc8ec5cb3ced4c4e6ac94dd1d60c0527f890740a8eb06ce9be422cb8da5cdafc2b58c0a5e24036c578de2a433c828ff7d60e0527f3b8ec09e026fdc305365dfc94e189a81b38c7597b3d941c279f042e8206e0bd8610100527fecd50eee38e386bd62be9bedb990706951b65fe053bd9d8a521af753d139e2da610120527fdefff6d330bb5403f63b14f33b578274160de3a50df4efecf0e0db73bcdd3da5610140527f617bdd11f7c0a11f49db22f629387a12da7596f9d1704d7465177c63d88ec7d7610160527f292c23a9aa1d8bea7e2435e555a4a60e379a5a35f3f452bae60121073fb6eead610180527fe1cea92ed99acdcb045a6726b2f87107e8a61620a232cf4d7d5b5766b3952e106101a0527f7ad66c0a68c72cb89e4fb4303841966e4062a76ab97451e3b9fb526a5ceb7f826101c0527fe026cc5a4aed3c22a58cbd3d2ac754c9352c5436f638042dca99034e836365166101e0527f3d04cffd8b46a874edf5cfae63077de85f849a660426697b06a829c70dd1409c610200527fad676aa337a485e4728a0b240d92b3ef7b3c372d06d189322bfd5f61f1e7203e610220527fa2fca4a49658f9fab7aa63289c91b7c7b6c832a6d0e69334ff5b0a3483d09dab610240527f4ebfd9cd7bca2505f7bef59cc1c12ecc708fff26ae4af19abe852afe9e20c862610260527f2def10d13dd169f550f578bda343d9717a138562e0093b380a1120789d53cf10610280527f776a31db34a1a0a7caaf862cffdfff1789297ffadc380bd3d39281d340abd3ad6102a0527fe2e7610b87a5fdf3a72ebe271287d923ab990eefac64b6e59d79f8b7e08c46e36102c0527f504364a5c6858bf98fff714ab5be9de19ed31a976860efbd0e772a2efe23e2e06102e0527f4f05f4acb83f5b65168d9fef89d56d4d77b8944015e6b1eed81b0238e2d0dba3610300527f44a6d974c75b07423e1d6d33f481916fdd45830aea11b6347e700cd8b9f0767c610320527fedf260291f734ddac396a956127dde4c34c0cfb8d8052f88ac139658ccf2d507610340527f6075c657a105351e7f0fce53bc320113324a522e8fd52dc878c762551e01a46e610360527f6ca6a3f763a9395f7da16014725ca7ee17e4815c0ff8119bf33f273dee11833b610380527f1c25ef10ffeb3c7d08aa707d17286e0b0d3cbcb50f1bd3b6523b63ba3b52dd0f6103a0527ffffc43bd08273ccf135fd3cacbeef055418e09eb728d727c4d5d5c556cdea7e36103c0527fc5ab8111456b1f28f3c7a0a604b4553ce905cb019c463ee159137af83c350b226103e0527f0ff273fcbf4ae0f2bd88d6cf319ff4004f8d7dca70d4ced4e74d2c74139739e6610400527f7fa06ba11241ddd5efdc65d4e39c9f6991b74fd4b81b62230808216c876f827c610420527f7e275adf313a996c7e2950cac67caba02a5ff925ebf9906b58949f3e77aec5b9610440527f8f6162fa308d2b3a15dc33cffac85f13ab349173121645aedf00f471663108be610460526200047e90600090601f6200075d565b50600060225560006023553480156200049657600080fd5b50604051620016ea380380620016ea833981016040819052620004b991620007b7565b604080516103e081019182905282917f78ccaaab73373552f207a63599de54d7d8d0c1805f86ce7da15818d09f4cff62916000918290601f9082845b815481526020019060010190808311620004f5575050600080516020620016ca833981519152602481905586935085925084915060015b601f8110156200066457602462000545600183620007e7565b601f81106200056457634e487b7160e01b600052603260045260246000fd5b0154602462000575600184620007e7565b601f81106200059457634e487b7160e01b600052603260045260246000fd5b015460408051602081019390935282015260600160405160208183030381529060405280519060200120602482601f8110620005e057634e487b7160e01b600052603260045260246000fd5b015560048110801590620005f4575080601f115b156200064f57602481601f81106200061c57634e487b7160e01b600052603260045260246000fd5b015460626200062d600484620007e7565b601b81106200064c57634e487b7160e01b600052603260045260246000fd5b01555b806200065b8162000801565b9150506200052c565b5062000674604383601f6200075d565b506022839055601f8481556024906200069090600190620007e7565b601f8110620006af57634e487b7160e01b600052603260045260246000fd5b01546024620006c16001601f620007e7565b601f8110620006e057634e487b7160e01b600052603260045260246000fd5b015460408051602081019390935282015260600160408051808303601f190181528282528051602091820120808255601f54918401919091529082015260600160408051808303601f190181529190528051602090910120602155505050505060609290921b6001600160601b0319166080525062000835915050565b82601f81019282156200078e579160200282015b828111156200078e57825182559160200191906001019062000771565b506200079c929150620007a0565b5090565b5b808211156200079c5760008155600101620007a1565b600060208284031215620007c9578081fd5b81516001600160a01b0381168114620007e0578182fd5b9392505050565b600082821015620007fc57620007fc6200081f565b500390565b60006000198214156200081857620008186200081f565b5060010190565b634e487b7160e01b600052601160045260246000fd5b60805160601c610e6f6200085b600039600081816101e1015261029f0152610e6f6000f3fe608060405234801561001057600080fd5b50600436106101215760003560e01c806395e4bf03116100ad578063d7c53ea711610071578063d7c53ea71461022f578063d828946314610238578063e829558814610241578063ebf0c71714610254578063f2aebf051461025d57600080fd5b806395e4bf03146101c957806398366e351461016f57806398d17621146101dc578063cab2da9b1461021b578063d0383d681461022457600080fd5b80635e71468b116100f45780635e71468b1461016f578063693c1db714610177578063709a8b2a14610180578063724a20f7146101a35780638d037962146101b657600080fd5b8063034a29ae146101265780631c4a7a941461014c5780631c76e77e1461015f57806349faa4d414610167575b600080fd5b610139610134366004610d52565b610284565b6040519081526020015b60405180910390f35b61013961015a366004610d0e565b61029b565b610139600481565b610139601081565b610139601f81565b61013960225481565b61019361018e366004610d6a565b6104a1565b6040519015158152602001610143565b6101396101b1366004610d52565b61050c565b6101396101c4366004610d52565b61051c565b6101396101d7366004610d37565b61052c565b6102037f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b039091168152602001610143565b61013960205481565b610139638000000081565b61013960235481565b610139601f5481565b61013961024f366004610d52565b610596565b61013960215481565b6101397f78ccaaab73373552f207a63599de54d7d8d0c1805f86ce7da15818d09f4cff6281565b604381601f811061029457600080fd5b0154905081565b60007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b1580156102f657600080fd5b505afa15801561030a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061032e9190610ce7565b6001600160a01b0316336001600160a01b03161461039f5760405162461bcd60e51b8152602060048201526024808201527f424c534163636f756e7452656769737472793a20496e76616c69642070726f7060448201526337b9b2b960e11b60648201526084015b60405180910390fd5b6103a7610c92565b60005b601081101561043a5760008482601081106103d557634e487b7160e01b600052603260045260246000fd5b608002016040516020016103e99190610daf565b6040516020818303038152906040528051906020012090508083836010811061042257634e487b7160e01b600052603260045260246000fd5b6020020152508061043281610e07565b9150506103aa565b506000610446826105a6565b90507f3154b80a7d9f6a143c37dde575f47deb78dacc7f280d8efc7e3ae102758a841b816001610477601083610dc4565b6104819190610df0565b604080519283526020830191909152015b60405180910390a19392505050565b600080836040516020016104b59190610daf565b604051602081830303815290604052805190602001209050610501818685601f806020026040519081016040528092919082601f602002808284376000920191909152506109a1915050565b9150505b9392505050565b600081601f811061029457600080fd5b606281601b811061029457600080fd5b600080826040516020016105409190610daf565b604051602081830303815290604052805190602001209050600061056382610ac3565b90507f59056afed767866d7c194cf26e24ebe16974ef943cccb729452e9adc265a9ac78160405161049291815260200190565b602481601f811061029457600080fd5b600060106105b960016380000000610df0565b6105c39190610df0565b602354106106135760405162461bcd60e51b815260206004820152601f60248201527f4163636f756e74547265653a207269676874207365742069732066756c6c20006044820152606401610396565b61061b610cb1565b60005b61062a60016004610df0565b6001901b8110156106f657600181901b84816010811061065a57634e487b7160e01b600052603260045260246000fd5b60200201518561066b836001610dc4565b6010811061068957634e487b7160e01b600052603260045260246000fd5b60200201516040516020016106a8929190918252602082015260400190565b604051602081830303815290604052805190602001208383600881106106de57634e487b7160e01b600052603260045260246000fd5b602002015250806106ee81610e07565b91505061061e565b5060015b60048110156108035760006001610712836004610df0565b61071c9190610df0565b6001901b905060005b818110156107ee57600181901b84816008811061075257634e487b7160e01b600052603260045260246000fd5b602002015185610763836001610dc4565b6008811061078157634e487b7160e01b600052603260045260246000fd5b60200201516040516020016107a0929190918252602082015260400190565b604051602081830303815290604052805190602001208583600881106107d657634e487b7160e01b600052603260045260246000fd5b602002015250806107e681610e07565b915050610725565b505080806107fb90610e07565b9150506106fa565b50805160235460009061081890601090610ddc565b90506000805b61082a6004601f610df0565b81101561093157826001166001141561089157606281601b811061085e57634e487b7160e01b600052603260045260246000fd5b01546040805160208101929092528101859052606001604051602081830303815290604052805190602001209350610918565b816108bf5783606282601b81106108b857634e487b7160e01b600052603260045260246000fd5b0155600191505b8360246108cd600484610dc4565b601f81106108eb57634e487b7160e01b600052603260045260246000fd5b01546040805160208101939093528201526060016040516020818303038152906040528051906020012093505b60019290921c918061092981610e07565b91505061081e565b506020838155601f5460405161095292869101918252602082015260400190565b60408051601f1981840301815291905280516020909101206021556023805460109190600090610983908490610dc4565b909155505060235461099790601090610df0565b9695505050505050565b6000806109b2638000000085610e22565b90508460005b601f811015610a9c578260011660011415610a2a578481601f81106109ed57634e487b7160e01b600052603260045260246000fd5b602002015182604051602001610a0d929190918252602082015260400190565b604051602081830303815290604052805190602001209150610a83565b818582601f8110610a4b57634e487b7160e01b600052603260045260246000fd5b6020020151604051602001610a6a929190918252602082015260400190565b6040516020818303038152906040528051906020012091505b60019290921c9180610a9481610e07565b9150506109b8565b506380000000851015610ab657601f541491506105059050565b6020541491506105059050565b6000610ad460016380000000610df0565b60225410610b245760405162461bcd60e51b815260206004820152601e60248201527f4163636f756e74547265653a206c656674207365742069732066756c6c2000006044820152606401610396565b60225482906000805b601f811015610c2c578260011660011415610b9657604381601f8110610b6357634e487b7160e01b600052603260045260246000fd5b01546040805160208101929092528101859052606001604051602081830303815290604052805190602001209350610c13565b81610bc45783604382601f8110610bbd57634e487b7160e01b600052603260045260246000fd5b0155600191505b83602482601f8110610be657634e487b7160e01b600052603260045260246000fd5b01546040805160208101939093528201526060016040516020818303038152906040528051906020012093505b60019290921c9180610c2481610e07565b915050610b2d565b50601f839055602080546040805192830186905282015260600160405160208183030381529060405280519060200120602181905550600160226000828254610c759190610dc4565b9091555050602254610c8990600190610df0565b95945050505050565b6040518061020001604052806010906020820280368337509192915050565b6040518061010001604052806008906020820280368337509192915050565b8060808101831015610ce157600080fd5b92915050565b600060208284031215610cf8578081fd5b81516001600160a01b0381168114610505578182fd5b6000610800808385031215610d21578182fd5b838184011115610d2f578182fd5b509092915050565b600060808284031215610d48578081fd5b6105058383610cd0565b600060208284031215610d63578081fd5b5035919050565b6000806000610480808587031215610d80578283fd5b84359350610d918660208701610cd0565b9250858186011115610da1578182fd5b5060a0840190509250925092565b60808282376000608091909101908152919050565b60008219821115610dd757610dd7610e36565b500190565b600082610deb57610deb610e4c565b500490565b600082821015610e0257610e02610e36565b500390565b6000600019821415610e1b57610e1b610e36565b5060010190565b600082610e3157610e31610e4c565b500690565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052601260045260246000fdfea164736f6c6343000804000a290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563",
}

// DefaultAccountRegistryABI is the input ABI used to generate the binding from.
// Deprecated: Use DefaultAccountRegistryMetaData.ABI instead.
var DefaultAccountRegistryABI = DefaultAccountRegistryMetaData.ABI

// DefaultAccountRegistryBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DefaultAccountRegistryMetaData.Bin instead.
var DefaultAccountRegistryBin = DefaultAccountRegistryMetaData.Bin

// DeployDefaultAccountRegistry deploys a new Ethereum contract, binding an instance of DefaultAccountRegistry to it.
func DeployDefaultAccountRegistry(auth *bind.TransactOpts, backend bind.ContractBackend, _chooser common.Address) (common.Address, *types.Transaction, *DefaultAccountRegistry, error) {
	parsed, err := DefaultAccountRegistryMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DefaultAccountRegistryBin), backend, _chooser)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DefaultAccountRegistry{DefaultAccountRegistryCaller: DefaultAccountRegistryCaller{contract: contract}, DefaultAccountRegistryTransactor: DefaultAccountRegistryTransactor{contract: contract}, DefaultAccountRegistryFilterer: DefaultAccountRegistryFilterer{contract: contract}}, nil
}

// DefaultAccountRegistry is an auto generated Go binding around an Ethereum contract.
type DefaultAccountRegistry struct {
	DefaultAccountRegistryCaller     // Read-only binding to the contract
	DefaultAccountRegistryTransactor // Write-only binding to the contract
	DefaultAccountRegistryFilterer   // Log filterer for contract events
}

// DefaultAccountRegistryCaller is an auto generated read-only Go binding around an Ethereum contract.
type DefaultAccountRegistryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DefaultAccountRegistryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DefaultAccountRegistryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DefaultAccountRegistryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DefaultAccountRegistryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DefaultAccountRegistrySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DefaultAccountRegistrySession struct {
	Contract     *DefaultAccountRegistry // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// DefaultAccountRegistryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DefaultAccountRegistryCallerSession struct {
	Contract *DefaultAccountRegistryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// DefaultAccountRegistryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DefaultAccountRegistryTransactorSession struct {
	Contract     *DefaultAccountRegistryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// DefaultAccountRegistryRaw is an auto generated low-level Go binding around an Ethereum contract.
type DefaultAccountRegistryRaw struct {
	Contract *DefaultAccountRegistry // Generic contract binding to access the raw methods on
}

// DefaultAccountRegistryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DefaultAccountRegistryCallerRaw struct {
	Contract *DefaultAccountRegistryCaller // Generic read-only contract binding to access the raw methods on
}

// DefaultAccountRegistryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DefaultAccountRegistryTransactorRaw struct {
	Contract *DefaultAccountRegistryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDefaultAccountRegistry creates a new instance of DefaultAccountRegistry, bound to a specific deployed contract.
func NewDefaultAccountRegistry(address common.Address, backend bind.ContractBackend) (*DefaultAccountRegistry, error) {
	contract, err := bindDefaultAccountRegistry(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DefaultAccountRegistry{DefaultAccountRegistryCaller: DefaultAccountRegistryCaller{contract: contract}, DefaultAccountRegistryTransactor: DefaultAccountRegistryTransactor{contract: contract}, DefaultAccountRegistryFilterer: DefaultAccountRegistryFilterer{contract: contract}}, nil
}

// NewDefaultAccountRegistryCaller creates a new read-only instance of DefaultAccountRegistry, bound to a specific deployed contract.
func NewDefaultAccountRegistryCaller(address common.Address, caller bind.ContractCaller) (*DefaultAccountRegistryCaller, error) {
	contract, err := bindDefaultAccountRegistry(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DefaultAccountRegistryCaller{contract: contract}, nil
}

// NewDefaultAccountRegistryTransactor creates a new write-only instance of DefaultAccountRegistry, bound to a specific deployed contract.
func NewDefaultAccountRegistryTransactor(address common.Address, transactor bind.ContractTransactor) (*DefaultAccountRegistryTransactor, error) {
	contract, err := bindDefaultAccountRegistry(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DefaultAccountRegistryTransactor{contract: contract}, nil
}

// NewDefaultAccountRegistryFilterer creates a new log filterer instance of DefaultAccountRegistry, bound to a specific deployed contract.
func NewDefaultAccountRegistryFilterer(address common.Address, filterer bind.ContractFilterer) (*DefaultAccountRegistryFilterer, error) {
	contract, err := bindDefaultAccountRegistry(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DefaultAccountRegistryFilterer{contract: contract}, nil
}

// bindDefaultAccountRegistry binds a generic wrapper to an already deployed contract.
func bindDefaultAccountRegistry(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DefaultAccountRegistryABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DefaultAccountRegistry *DefaultAccountRegistryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DefaultAccountRegistry.Contract.DefaultAccountRegistryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DefaultAccountRegistry *DefaultAccountRegistryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DefaultAccountRegistry.Contract.DefaultAccountRegistryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DefaultAccountRegistry *DefaultAccountRegistryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DefaultAccountRegistry.Contract.DefaultAccountRegistryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DefaultAccountRegistry *DefaultAccountRegistryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DefaultAccountRegistry.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DefaultAccountRegistry *DefaultAccountRegistryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DefaultAccountRegistry.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DefaultAccountRegistry *DefaultAccountRegistryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DefaultAccountRegistry.Contract.contract.Transact(opts, method, params...)
}

// BATCHDEPTH is a free data retrieval call binding the contract method 0x1c76e77e.
//
// Solidity: function BATCH_DEPTH() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistryCaller) BATCHDEPTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DefaultAccountRegistry.contract.Call(opts, &out, "BATCH_DEPTH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BATCHDEPTH is a free data retrieval call binding the contract method 0x1c76e77e.
//
// Solidity: function BATCH_DEPTH() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistrySession) BATCHDEPTH() (*big.Int, error) {
	return _DefaultAccountRegistry.Contract.BATCHDEPTH(&_DefaultAccountRegistry.CallOpts)
}

// BATCHDEPTH is a free data retrieval call binding the contract method 0x1c76e77e.
//
// Solidity: function BATCH_DEPTH() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistryCallerSession) BATCHDEPTH() (*big.Int, error) {
	return _DefaultAccountRegistry.Contract.BATCHDEPTH(&_DefaultAccountRegistry.CallOpts)
}

// BATCHSIZE is a free data retrieval call binding the contract method 0x49faa4d4.
//
// Solidity: function BATCH_SIZE() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistryCaller) BATCHSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DefaultAccountRegistry.contract.Call(opts, &out, "BATCH_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BATCHSIZE is a free data retrieval call binding the contract method 0x49faa4d4.
//
// Solidity: function BATCH_SIZE() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistrySession) BATCHSIZE() (*big.Int, error) {
	return _DefaultAccountRegistry.Contract.BATCHSIZE(&_DefaultAccountRegistry.CallOpts)
}

// BATCHSIZE is a free data retrieval call binding the contract method 0x49faa4d4.
//
// Solidity: function BATCH_SIZE() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistryCallerSession) BATCHSIZE() (*big.Int, error) {
	return _DefaultAccountRegistry.Contract.BATCHSIZE(&_DefaultAccountRegistry.CallOpts)
}

// DEPTH is a free data retrieval call binding the contract method 0x98366e35.
//
// Solidity: function DEPTH() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistryCaller) DEPTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DefaultAccountRegistry.contract.Call(opts, &out, "DEPTH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DEPTH is a free data retrieval call binding the contract method 0x98366e35.
//
// Solidity: function DEPTH() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistrySession) DEPTH() (*big.Int, error) {
	return _DefaultAccountRegistry.Contract.DEPTH(&_DefaultAccountRegistry.CallOpts)
}

// DEPTH is a free data retrieval call binding the contract method 0x98366e35.
//
// Solidity: function DEPTH() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistryCallerSession) DEPTH() (*big.Int, error) {
	return _DefaultAccountRegistry.Contract.DEPTH(&_DefaultAccountRegistry.CallOpts)
}

// INITIALLEFTROOT is a free data retrieval call binding the contract method 0xf2aebf05.
//
// Solidity: function INITIAL_LEFT_ROOT() view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistryCaller) INITIALLEFTROOT(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DefaultAccountRegistry.contract.Call(opts, &out, "INITIAL_LEFT_ROOT")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// INITIALLEFTROOT is a free data retrieval call binding the contract method 0xf2aebf05.
//
// Solidity: function INITIAL_LEFT_ROOT() view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistrySession) INITIALLEFTROOT() ([32]byte, error) {
	return _DefaultAccountRegistry.Contract.INITIALLEFTROOT(&_DefaultAccountRegistry.CallOpts)
}

// INITIALLEFTROOT is a free data retrieval call binding the contract method 0xf2aebf05.
//
// Solidity: function INITIAL_LEFT_ROOT() view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistryCallerSession) INITIALLEFTROOT() ([32]byte, error) {
	return _DefaultAccountRegistry.Contract.INITIALLEFTROOT(&_DefaultAccountRegistry.CallOpts)
}

// SETSIZE is a free data retrieval call binding the contract method 0xd0383d68.
//
// Solidity: function SET_SIZE() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistryCaller) SETSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DefaultAccountRegistry.contract.Call(opts, &out, "SET_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SETSIZE is a free data retrieval call binding the contract method 0xd0383d68.
//
// Solidity: function SET_SIZE() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistrySession) SETSIZE() (*big.Int, error) {
	return _DefaultAccountRegistry.Contract.SETSIZE(&_DefaultAccountRegistry.CallOpts)
}

// SETSIZE is a free data retrieval call binding the contract method 0xd0383d68.
//
// Solidity: function SET_SIZE() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistryCallerSession) SETSIZE() (*big.Int, error) {
	return _DefaultAccountRegistry.Contract.SETSIZE(&_DefaultAccountRegistry.CallOpts)
}

// WITNESSLENGTH is a free data retrieval call binding the contract method 0x5e71468b.
//
// Solidity: function WITNESS_LENGTH() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistryCaller) WITNESSLENGTH(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DefaultAccountRegistry.contract.Call(opts, &out, "WITNESS_LENGTH")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WITNESSLENGTH is a free data retrieval call binding the contract method 0x5e71468b.
//
// Solidity: function WITNESS_LENGTH() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistrySession) WITNESSLENGTH() (*big.Int, error) {
	return _DefaultAccountRegistry.Contract.WITNESSLENGTH(&_DefaultAccountRegistry.CallOpts)
}

// WITNESSLENGTH is a free data retrieval call binding the contract method 0x5e71468b.
//
// Solidity: function WITNESS_LENGTH() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistryCallerSession) WITNESSLENGTH() (*big.Int, error) {
	return _DefaultAccountRegistry.Contract.WITNESSLENGTH(&_DefaultAccountRegistry.CallOpts)
}

// Chooser is a free data retrieval call binding the contract method 0x98d17621.
//
// Solidity: function chooser() view returns(address)
func (_DefaultAccountRegistry *DefaultAccountRegistryCaller) Chooser(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DefaultAccountRegistry.contract.Call(opts, &out, "chooser")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Chooser is a free data retrieval call binding the contract method 0x98d17621.
//
// Solidity: function chooser() view returns(address)
func (_DefaultAccountRegistry *DefaultAccountRegistrySession) Chooser() (common.Address, error) {
	return _DefaultAccountRegistry.Contract.Chooser(&_DefaultAccountRegistry.CallOpts)
}

// Chooser is a free data retrieval call binding the contract method 0x98d17621.
//
// Solidity: function chooser() view returns(address)
func (_DefaultAccountRegistry *DefaultAccountRegistryCallerSession) Chooser() (common.Address, error) {
	return _DefaultAccountRegistry.Contract.Chooser(&_DefaultAccountRegistry.CallOpts)
}

// Exists is a free data retrieval call binding the contract method 0x709a8b2a.
//
// Solidity: function exists(uint256 pubkeyID, uint256[4] pubkey, bytes32[31] witness) view returns(bool)
func (_DefaultAccountRegistry *DefaultAccountRegistryCaller) Exists(opts *bind.CallOpts, pubkeyID *big.Int, pubkey [4]*big.Int, witness [31][32]byte) (bool, error) {
	var out []interface{}
	err := _DefaultAccountRegistry.contract.Call(opts, &out, "exists", pubkeyID, pubkey, witness)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Exists is a free data retrieval call binding the contract method 0x709a8b2a.
//
// Solidity: function exists(uint256 pubkeyID, uint256[4] pubkey, bytes32[31] witness) view returns(bool)
func (_DefaultAccountRegistry *DefaultAccountRegistrySession) Exists(pubkeyID *big.Int, pubkey [4]*big.Int, witness [31][32]byte) (bool, error) {
	return _DefaultAccountRegistry.Contract.Exists(&_DefaultAccountRegistry.CallOpts, pubkeyID, pubkey, witness)
}

// Exists is a free data retrieval call binding the contract method 0x709a8b2a.
//
// Solidity: function exists(uint256 pubkeyID, uint256[4] pubkey, bytes32[31] witness) view returns(bool)
func (_DefaultAccountRegistry *DefaultAccountRegistryCallerSession) Exists(pubkeyID *big.Int, pubkey [4]*big.Int, witness [31][32]byte) (bool, error) {
	return _DefaultAccountRegistry.Contract.Exists(&_DefaultAccountRegistry.CallOpts, pubkeyID, pubkey, witness)
}

// FilledSubtreesLeft is a free data retrieval call binding the contract method 0x034a29ae.
//
// Solidity: function filledSubtreesLeft(uint256 ) view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistryCaller) FilledSubtreesLeft(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _DefaultAccountRegistry.contract.Call(opts, &out, "filledSubtreesLeft", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FilledSubtreesLeft is a free data retrieval call binding the contract method 0x034a29ae.
//
// Solidity: function filledSubtreesLeft(uint256 ) view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistrySession) FilledSubtreesLeft(arg0 *big.Int) ([32]byte, error) {
	return _DefaultAccountRegistry.Contract.FilledSubtreesLeft(&_DefaultAccountRegistry.CallOpts, arg0)
}

// FilledSubtreesLeft is a free data retrieval call binding the contract method 0x034a29ae.
//
// Solidity: function filledSubtreesLeft(uint256 ) view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistryCallerSession) FilledSubtreesLeft(arg0 *big.Int) ([32]byte, error) {
	return _DefaultAccountRegistry.Contract.FilledSubtreesLeft(&_DefaultAccountRegistry.CallOpts, arg0)
}

// FilledSubtreesRight is a free data retrieval call binding the contract method 0x8d037962.
//
// Solidity: function filledSubtreesRight(uint256 ) view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistryCaller) FilledSubtreesRight(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _DefaultAccountRegistry.contract.Call(opts, &out, "filledSubtreesRight", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// FilledSubtreesRight is a free data retrieval call binding the contract method 0x8d037962.
//
// Solidity: function filledSubtreesRight(uint256 ) view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistrySession) FilledSubtreesRight(arg0 *big.Int) ([32]byte, error) {
	return _DefaultAccountRegistry.Contract.FilledSubtreesRight(&_DefaultAccountRegistry.CallOpts, arg0)
}

// FilledSubtreesRight is a free data retrieval call binding the contract method 0x8d037962.
//
// Solidity: function filledSubtreesRight(uint256 ) view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistryCallerSession) FilledSubtreesRight(arg0 *big.Int) ([32]byte, error) {
	return _DefaultAccountRegistry.Contract.FilledSubtreesRight(&_DefaultAccountRegistry.CallOpts, arg0)
}

// InitialSubtrees is a free data retrieval call binding the contract method 0x724a20f7.
//
// Solidity: function initialSubtrees(uint256 ) view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistryCaller) InitialSubtrees(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _DefaultAccountRegistry.contract.Call(opts, &out, "initialSubtrees", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// InitialSubtrees is a free data retrieval call binding the contract method 0x724a20f7.
//
// Solidity: function initialSubtrees(uint256 ) view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistrySession) InitialSubtrees(arg0 *big.Int) ([32]byte, error) {
	return _DefaultAccountRegistry.Contract.InitialSubtrees(&_DefaultAccountRegistry.CallOpts, arg0)
}

// InitialSubtrees is a free data retrieval call binding the contract method 0x724a20f7.
//
// Solidity: function initialSubtrees(uint256 ) view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistryCallerSession) InitialSubtrees(arg0 *big.Int) ([32]byte, error) {
	return _DefaultAccountRegistry.Contract.InitialSubtrees(&_DefaultAccountRegistry.CallOpts, arg0)
}

// LeafIndexLeft is a free data retrieval call binding the contract method 0x693c1db7.
//
// Solidity: function leafIndexLeft() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistryCaller) LeafIndexLeft(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DefaultAccountRegistry.contract.Call(opts, &out, "leafIndexLeft")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LeafIndexLeft is a free data retrieval call binding the contract method 0x693c1db7.
//
// Solidity: function leafIndexLeft() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistrySession) LeafIndexLeft() (*big.Int, error) {
	return _DefaultAccountRegistry.Contract.LeafIndexLeft(&_DefaultAccountRegistry.CallOpts)
}

// LeafIndexLeft is a free data retrieval call binding the contract method 0x693c1db7.
//
// Solidity: function leafIndexLeft() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistryCallerSession) LeafIndexLeft() (*big.Int, error) {
	return _DefaultAccountRegistry.Contract.LeafIndexLeft(&_DefaultAccountRegistry.CallOpts)
}

// LeafIndexRight is a free data retrieval call binding the contract method 0xd7c53ea7.
//
// Solidity: function leafIndexRight() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistryCaller) LeafIndexRight(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DefaultAccountRegistry.contract.Call(opts, &out, "leafIndexRight")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LeafIndexRight is a free data retrieval call binding the contract method 0xd7c53ea7.
//
// Solidity: function leafIndexRight() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistrySession) LeafIndexRight() (*big.Int, error) {
	return _DefaultAccountRegistry.Contract.LeafIndexRight(&_DefaultAccountRegistry.CallOpts)
}

// LeafIndexRight is a free data retrieval call binding the contract method 0xd7c53ea7.
//
// Solidity: function leafIndexRight() view returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistryCallerSession) LeafIndexRight() (*big.Int, error) {
	return _DefaultAccountRegistry.Contract.LeafIndexRight(&_DefaultAccountRegistry.CallOpts)
}

// Root is a free data retrieval call binding the contract method 0xebf0c717.
//
// Solidity: function root() view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistryCaller) Root(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DefaultAccountRegistry.contract.Call(opts, &out, "root")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Root is a free data retrieval call binding the contract method 0xebf0c717.
//
// Solidity: function root() view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistrySession) Root() ([32]byte, error) {
	return _DefaultAccountRegistry.Contract.Root(&_DefaultAccountRegistry.CallOpts)
}

// Root is a free data retrieval call binding the contract method 0xebf0c717.
//
// Solidity: function root() view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistryCallerSession) Root() ([32]byte, error) {
	return _DefaultAccountRegistry.Contract.Root(&_DefaultAccountRegistry.CallOpts)
}

// RootLeft is a free data retrieval call binding the contract method 0xd8289463.
//
// Solidity: function rootLeft() view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistryCaller) RootLeft(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DefaultAccountRegistry.contract.Call(opts, &out, "rootLeft")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RootLeft is a free data retrieval call binding the contract method 0xd8289463.
//
// Solidity: function rootLeft() view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistrySession) RootLeft() ([32]byte, error) {
	return _DefaultAccountRegistry.Contract.RootLeft(&_DefaultAccountRegistry.CallOpts)
}

// RootLeft is a free data retrieval call binding the contract method 0xd8289463.
//
// Solidity: function rootLeft() view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistryCallerSession) RootLeft() ([32]byte, error) {
	return _DefaultAccountRegistry.Contract.RootLeft(&_DefaultAccountRegistry.CallOpts)
}

// RootRight is a free data retrieval call binding the contract method 0xcab2da9b.
//
// Solidity: function rootRight() view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistryCaller) RootRight(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DefaultAccountRegistry.contract.Call(opts, &out, "rootRight")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// RootRight is a free data retrieval call binding the contract method 0xcab2da9b.
//
// Solidity: function rootRight() view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistrySession) RootRight() ([32]byte, error) {
	return _DefaultAccountRegistry.Contract.RootRight(&_DefaultAccountRegistry.CallOpts)
}

// RootRight is a free data retrieval call binding the contract method 0xcab2da9b.
//
// Solidity: function rootRight() view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistryCallerSession) RootRight() ([32]byte, error) {
	return _DefaultAccountRegistry.Contract.RootRight(&_DefaultAccountRegistry.CallOpts)
}

// Zeros is a free data retrieval call binding the contract method 0xe8295588.
//
// Solidity: function zeros(uint256 ) view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistryCaller) Zeros(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _DefaultAccountRegistry.contract.Call(opts, &out, "zeros", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Zeros is a free data retrieval call binding the contract method 0xe8295588.
//
// Solidity: function zeros(uint256 ) view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistrySession) Zeros(arg0 *big.Int) ([32]byte, error) {
	return _DefaultAccountRegistry.Contract.Zeros(&_DefaultAccountRegistry.CallOpts, arg0)
}

// Zeros is a free data retrieval call binding the contract method 0xe8295588.
//
// Solidity: function zeros(uint256 ) view returns(bytes32)
func (_DefaultAccountRegistry *DefaultAccountRegistryCallerSession) Zeros(arg0 *big.Int) ([32]byte, error) {
	return _DefaultAccountRegistry.Contract.Zeros(&_DefaultAccountRegistry.CallOpts, arg0)
}

// Register is a paid mutator transaction binding the contract method 0x95e4bf03.
//
// Solidity: function register(uint256[4] pubkey) returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistryTransactor) Register(opts *bind.TransactOpts, pubkey [4]*big.Int) (*types.Transaction, error) {
	return _DefaultAccountRegistry.contract.Transact(opts, "register", pubkey)
}

// Register is a paid mutator transaction binding the contract method 0x95e4bf03.
//
// Solidity: function register(uint256[4] pubkey) returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistrySession) Register(pubkey [4]*big.Int) (*types.Transaction, error) {
	return _DefaultAccountRegistry.Contract.Register(&_DefaultAccountRegistry.TransactOpts, pubkey)
}

// Register is a paid mutator transaction binding the contract method 0x95e4bf03.
//
// Solidity: function register(uint256[4] pubkey) returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistryTransactorSession) Register(pubkey [4]*big.Int) (*types.Transaction, error) {
	return _DefaultAccountRegistry.Contract.Register(&_DefaultAccountRegistry.TransactOpts, pubkey)
}

// RegisterBatch is a paid mutator transaction binding the contract method 0x1c4a7a94.
//
// Solidity: function registerBatch(uint256[4][16] pubkeys) returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistryTransactor) RegisterBatch(opts *bind.TransactOpts, pubkeys [16][4]*big.Int) (*types.Transaction, error) {
	return _DefaultAccountRegistry.contract.Transact(opts, "registerBatch", pubkeys)
}

// RegisterBatch is a paid mutator transaction binding the contract method 0x1c4a7a94.
//
// Solidity: function registerBatch(uint256[4][16] pubkeys) returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistrySession) RegisterBatch(pubkeys [16][4]*big.Int) (*types.Transaction, error) {
	return _DefaultAccountRegistry.Contract.RegisterBatch(&_DefaultAccountRegistry.TransactOpts, pubkeys)
}

// RegisterBatch is a paid mutator transaction binding the contract method 0x1c4a7a94.
//
// Solidity: function registerBatch(uint256[4][16] pubkeys) returns(uint256)
func (_DefaultAccountRegistry *DefaultAccountRegistryTransactorSession) RegisterBatch(pubkeys [16][4]*big.Int) (*types.Transaction, error) {
	return _DefaultAccountRegistry.Contract.RegisterBatch(&_DefaultAccountRegistry.TransactOpts, pubkeys)
}

// DefaultAccountRegistryBatchPubkeyRegisteredIterator is returned from FilterBatchPubkeyRegistered and is used to iterate over the raw logs and unpacked data for BatchPubkeyRegistered events raised by the DefaultAccountRegistry contract.
type DefaultAccountRegistryBatchPubkeyRegisteredIterator struct {
	Event *DefaultAccountRegistryBatchPubkeyRegistered // Event containing the contract specifics and raw log

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
func (it *DefaultAccountRegistryBatchPubkeyRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DefaultAccountRegistryBatchPubkeyRegistered)
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
		it.Event = new(DefaultAccountRegistryBatchPubkeyRegistered)
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
func (it *DefaultAccountRegistryBatchPubkeyRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DefaultAccountRegistryBatchPubkeyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DefaultAccountRegistryBatchPubkeyRegistered represents a BatchPubkeyRegistered event raised by the DefaultAccountRegistry contract.
type DefaultAccountRegistryBatchPubkeyRegistered struct {
	StartID *big.Int
	EndID   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBatchPubkeyRegistered is a free log retrieval operation binding the contract event 0x3154b80a7d9f6a143c37dde575f47deb78dacc7f280d8efc7e3ae102758a841b.
//
// Solidity: event BatchPubkeyRegistered(uint256 startID, uint256 endID)
func (_DefaultAccountRegistry *DefaultAccountRegistryFilterer) FilterBatchPubkeyRegistered(opts *bind.FilterOpts) (*DefaultAccountRegistryBatchPubkeyRegisteredIterator, error) {

	logs, sub, err := _DefaultAccountRegistry.contract.FilterLogs(opts, "BatchPubkeyRegistered")
	if err != nil {
		return nil, err
	}
	return &DefaultAccountRegistryBatchPubkeyRegisteredIterator{contract: _DefaultAccountRegistry.contract, event: "BatchPubkeyRegistered", logs: logs, sub: sub}, nil
}

// WatchBatchPubkeyRegistered is a free log subscription operation binding the contract event 0x3154b80a7d9f6a143c37dde575f47deb78dacc7f280d8efc7e3ae102758a841b.
//
// Solidity: event BatchPubkeyRegistered(uint256 startID, uint256 endID)
func (_DefaultAccountRegistry *DefaultAccountRegistryFilterer) WatchBatchPubkeyRegistered(opts *bind.WatchOpts, sink chan<- *DefaultAccountRegistryBatchPubkeyRegistered) (event.Subscription, error) {

	logs, sub, err := _DefaultAccountRegistry.contract.WatchLogs(opts, "BatchPubkeyRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DefaultAccountRegistryBatchPubkeyRegistered)
				if err := _DefaultAccountRegistry.contract.UnpackLog(event, "BatchPubkeyRegistered", log); err != nil {
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
func (_DefaultAccountRegistry *DefaultAccountRegistryFilterer) ParseBatchPubkeyRegistered(log types.Log) (*DefaultAccountRegistryBatchPubkeyRegistered, error) {
	event := new(DefaultAccountRegistryBatchPubkeyRegistered)
	if err := _DefaultAccountRegistry.contract.UnpackLog(event, "BatchPubkeyRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DefaultAccountRegistrySinglePubkeyRegisteredIterator is returned from FilterSinglePubkeyRegistered and is used to iterate over the raw logs and unpacked data for SinglePubkeyRegistered events raised by the DefaultAccountRegistry contract.
type DefaultAccountRegistrySinglePubkeyRegisteredIterator struct {
	Event *DefaultAccountRegistrySinglePubkeyRegistered // Event containing the contract specifics and raw log

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
func (it *DefaultAccountRegistrySinglePubkeyRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DefaultAccountRegistrySinglePubkeyRegistered)
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
		it.Event = new(DefaultAccountRegistrySinglePubkeyRegistered)
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
func (it *DefaultAccountRegistrySinglePubkeyRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DefaultAccountRegistrySinglePubkeyRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DefaultAccountRegistrySinglePubkeyRegistered represents a SinglePubkeyRegistered event raised by the DefaultAccountRegistry contract.
type DefaultAccountRegistrySinglePubkeyRegistered struct {
	PubkeyID *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSinglePubkeyRegistered is a free log retrieval operation binding the contract event 0x59056afed767866d7c194cf26e24ebe16974ef943cccb729452e9adc265a9ac7.
//
// Solidity: event SinglePubkeyRegistered(uint256 pubkeyID)
func (_DefaultAccountRegistry *DefaultAccountRegistryFilterer) FilterSinglePubkeyRegistered(opts *bind.FilterOpts) (*DefaultAccountRegistrySinglePubkeyRegisteredIterator, error) {

	logs, sub, err := _DefaultAccountRegistry.contract.FilterLogs(opts, "SinglePubkeyRegistered")
	if err != nil {
		return nil, err
	}
	return &DefaultAccountRegistrySinglePubkeyRegisteredIterator{contract: _DefaultAccountRegistry.contract, event: "SinglePubkeyRegistered", logs: logs, sub: sub}, nil
}

// WatchSinglePubkeyRegistered is a free log subscription operation binding the contract event 0x59056afed767866d7c194cf26e24ebe16974ef943cccb729452e9adc265a9ac7.
//
// Solidity: event SinglePubkeyRegistered(uint256 pubkeyID)
func (_DefaultAccountRegistry *DefaultAccountRegistryFilterer) WatchSinglePubkeyRegistered(opts *bind.WatchOpts, sink chan<- *DefaultAccountRegistrySinglePubkeyRegistered) (event.Subscription, error) {

	logs, sub, err := _DefaultAccountRegistry.contract.WatchLogs(opts, "SinglePubkeyRegistered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DefaultAccountRegistrySinglePubkeyRegistered)
				if err := _DefaultAccountRegistry.contract.UnpackLog(event, "SinglePubkeyRegistered", log); err != nil {
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
func (_DefaultAccountRegistry *DefaultAccountRegistryFilterer) ParseSinglePubkeyRegistered(log types.Log) (*DefaultAccountRegistrySinglePubkeyRegistered, error) {
	event := new(DefaultAccountRegistrySinglePubkeyRegistered)
	if err := _DefaultAccountRegistry.contract.UnpackLog(event, "SinglePubkeyRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
