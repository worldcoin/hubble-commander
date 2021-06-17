// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package transfer

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

// TypesAuthCommon is an auto generated low-level Go binding around an user-defined struct.
type TypesAuthCommon struct {
	Signature   [2]*big.Int
	StateRoot   [32]byte
	AccountRoot [32]byte
	Domain      [32]byte
	Txs         []byte
}

// TypesSignatureProof is an auto generated low-level Go binding around an user-defined struct.
type TypesSignatureProof struct {
	States          []TypesUserState
	StateWitnesses  [][][32]byte
	Pubkeys         [][4]*big.Int
	PubkeyWitnesses [][][32]byte
}

// TypesStateMerkleProof is an auto generated low-level Go binding around an user-defined struct.
type TypesStateMerkleProof struct {
	State   TypesUserState
	Witness [][32]byte
}

// TypesUserState is an auto generated low-level Go binding around an user-defined struct.
type TypesUserState struct {
	PubkeyID *big.Int
	TokenID  *big.Int
	Balance  *big.Int
	Nonce    *big.Int
}

// TransferABI is the input ABI used to generate the binding from.
const TransferABI = "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"domain\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.AuthCommon\",\"name\":\"common\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState[]\",\"name\":\"states\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"stateWitnesses\",\"type\":\"bytes32[][]\"},{\"internalType\":\"uint256[4][]\",\"name\":\"pubkeys\",\"type\":\"uint256[4][]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"pubkeyWitnesses\",\"type\":\"bytes32[][]\"}],\"internalType\":\"structTypes.SignatureProof\",\"name\":\"proof\",\"type\":\"tuple\"}],\"name\":\"checkSignature\",\"outputs\":[{\"internalType\":\"enumTypes.Result\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"maxTxSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState\",\"name\":\"state\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.StateMerkleProof[]\",\"name\":\"proofs\",\"type\":\"tuple[]\"}],\"name\":\"processTransferCommit\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"enumTypes.Result\",\"name\":\"result\",\"type\":\"uint8\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// TransferBin is the compiled bytecode used for deploying new contracts.
var TransferBin = "0x608060405234801561001057600080fd5b50612f23806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806344ec5a771461003b5780639293149214610064575b600080fd5b61004e6100493660046129f5565b610085565b60405161005b9190612b94565b60405180910390f35b610077610072366004612928565b61009a565b60405161005b929190612b80565b600061009183836101cf565b90505b92915050565b6000806100a684610500565b156100b6575085905060086101c5565b60006100c185610518565b9050868111156100d85787600992509250506101c5565b600080856000815181106100e857fe5b602002602001015160000151602001519050610102612481565b60005b8481101561019457610117898261052c565b91506101558c83858b856002028151811061012e57fe5b60200260200101518c866002026001018151811061014857fe5b602002602001015161059f565b909c509550600086600a81111561016857fe5b1461017a578b965050505050506101c5565b606082015161018a9085906105fc565b9350600101610105565b506101b88b8a84868b89600202815181106101ab57fe5b6020026020010151610621565b909b508b96509450505050505b9550959350505050565b6000806101df8460800151610518565b905060608167ffffffffffffffff811180156101fa57600080fd5b5060405190808252806020026020018201604052801561023457816020015b6102216124a9565b8152602001906001900390816102195790505b50905060608267ffffffffffffffff8111801561025057600080fd5b5060405190808252806020026020018201604052801561027a578160200160208202803683370190505b50905060005b838110156104b15760001981850301610297612481565b60808901516102a6908361052c565b90506102f989602001516102d08a6000015185815181106102c357fe5b60200260200101516106cf565b8051906020012083600001518b6020015186815181106102ec57fe5b602002602001015161070e565b61031e5760405162461bcd60e51b815260040161031590612bd9565b60405180910390fd5b60008860000151838151811061033057fe5b602002602001015160600151116103595760405162461bcd60e51b815260040161031590612ce0565b6103cb89604001518960400151848151811061037157fe5b60200260200101516040516020016103899190612ac2565b604051602081830303815290604052805190602001208a6000015185815181106103af57fe5b6020026020010151600001518b6060015186815181106102ec57fe5b6103e75760405162461bcd60e51b815260040161031590612ddc565b60006001896000015184815181106103fb57fe5b6020026020010151606001510390506000826000015166038d7ea4c6800001905060005b85811015610452578187828151811061043457fe5b6020026020010151141561044a57600019909201915b60010161041f565b508086868151811061046057fe5b60200260200101818152505060606104788484610726565b90506104888c606001518261076b565b88868151811061049457fe5b602002602001018190525050505050508080600101915050610280565b506000806104c8886000015188604001518661082a565b92509050816104df57600a95505050505050610094565b806104f257600595505050505050610094565b506000979650505050505050565b6000600c82518161050d57fe5b06151590505b919050565b6000600c82518161052557fe5b0492915050565b610534612481565b506004600c8281028401918201516008830151600a80850151948401516040805160808101825263ffffffff9586168152939094166020840152600f86861c8116830a610fff97881602948401949094529384901c90921690910a9190921602606082015292915050565b6000806105bc8787600001518789604001518a6060015189610c5d565b9092509050600081600a8111156105cf57fe5b146105d9576101c5565b6105ee82876020015187896040015187610621565b909890975095505050505050565b6000828201838110156100915760405162461bcd60e51b815260040161031590612c1d565b6000806106478761063585600001516106cf565b8051906020012088866020015161070e565b6106635760405162461bcd60e51b815260040161031590612e63565b61066b612481565b600061067c87878760000151610d10565b9092509050600081600a81111561068f57fe5b146106a1576000935091506101c59050565b6106bf6106ad836106cf565b80519060200120898760200151610d81565b9960009950975050505050505050565b606081600001518260200151836040015184606001516040516020016106f89493929190612b2f565b6040516020818303038152906040529050919050565b60008461071c858585610d81565b1495945050505050565b6060600183600001518460200151848660400151876060015160405160200161075496959493929190612b4a565b604051602081830303815290604052905092915050565b6107736124a9565b61077b6124a9565b6107858484610e2e565b905061078f6124a9565b6107a08260005b6020020151610ee9565b90506107aa6124a9565b6107b5836001610796565b90506107bf6124c7565b825181526020808401518282015282516040808401919091529083015160608301526000908460808460066107d05a03fa90508080156107fe57610800565bfe5b508061081e5760405162461bcd60e51b815260040161031590612ba2565b50919695505050505050565b815160009081908061084e5760405162461bcd60e51b815260040161031590612d58565b8351811461086e5760405162461bcd60e51b815260040161031590612c8b565b6006600182010260608167ffffffffffffffff8111801561088e57600080fd5b506040519080825280602002602001820160405280156108b8578160200160208202803683370190505b5090508760006020020151816000815181106108d057fe5b60209081029190910101528760016020020151816001815181106108f057fe5b6020026020010181815250507f198e9393920d483a7260bfb731fb5d25f1aa493335a9e71297e485b7aef312c28160028151811061092a57fe5b6020026020010181815250507f1800deef121f1e76426a00665e5c4479674322d4f75edadd46debd5cd992f6ed8160038151811061096457fe5b6020026020010181815250507f275dc4a288d1afb3cbb1ac09187524c7db36395df7be3b99e673b13a075a65ec8160048151811061099e57fe5b6020026020010181815250507f1d9befcd05a5323e6da4d435f3b617cdb3af83285c2df711ef39c01571827f9d816005815181106109d857fe5b60200260200101818152505060005b83811015610b88578681815181106109fb57fe5b6020026020010151600060028110610a0f57fe5b6020020151828260060260060181518110610a2657fe5b602002602001018181525050868181518110610a3e57fe5b6020026020010151600160028110610a5257fe5b6020020151828260060260070181518110610a6957fe5b602002602001018181525050878181518110610a8157fe5b6020026020010151600160048110610a9557fe5b6020020151828260060260080181518110610aac57fe5b602002602001018181525050878181518110610ac457fe5b6020026020010151600060048110610ad857fe5b6020020151828260060260090181518110610aef57fe5b602002602001018181525050878181518110610b0757fe5b6020026020010151600360048110610b1b57fe5b60200201518282600602600a0181518110610b3257fe5b602002602001018181525050878181518110610b4a57fe5b6020026020010151600260048110610b5e57fe5b60200201518282600602600b0181518110610b7557fe5b60209081029190910101526001016109e7565b50610b916124e5565b60405163273cfc6560e11b815260009073079d8077c465bd0bf0fc502ad2b846757e41566190634e79f8ca90610bce906001890190600401612ea6565b60206040518083038186803b158015610be657600080fd5b505afa158015610bfa573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c1e9190612910565b90506020826020860260208601600885fa955085610c4757600080965096505050505050610c55565b505115159450600193505050505b935093915050565b600080610c8388610c7185600001516106cf565b8051906020012089866020015161070e565b610c9f5760405162461bcd60e51b815260040161031590612d17565b610ca7612481565b6000610cb98888888860000151611265565b9092509050600081600a811115610ccc57fe5b14610cde57600093509150610d059050565b610cfc610cea836106cf565b805190602001208a8760200151610d81565b93506000925050505b965096945050505050565b610d18612481565b600084836020015114610d3057508190506004610c55565b60405180608001604052808460000151815260200184602001518152602001610d668686604001516105fc90919063ffffffff16565b81526060850151602090910152915060009050935093915050565b600083815b8351811015610e2357600185821c16610ddc5781848281518110610da657fe5b6020026020010151604051602001610dbf929190612b72565b604051602081830303815290604052805190602001209150610e1b565b838181518110610de857fe5b602002602001015182604051602001610e02929190612b72565b6040516020818303038152906040528051906020012091505b600101610d86565b5090505b9392505050565b610e366124a9565b6060610e42848461132d565b9050600080600080601885016001600160c01b0381511693506030860190506001600160c01b038151169450600080516020612ef783398151915285600080516020612ef7833981519152600160c01b870908604887015160608801516001600160c01b0390811697501694509250600080516020612ef783398151915290508481600160c01b860908604080518082019091529283526020830152509695505050505050565b610ef16124a9565b600080516020612ef78339815191528210610f1e5760405162461bcd60e51b815260040161031590612d99565b816000610f2a826115a9565b9150506000600080516020612ef783398151915280610f4557fe5b8384099050600080516020612ef78339815191526004820890506000600080516020612ef783398151915277b3c4d79d41a91759a9e4c7e359b6b89eaec68e62effffffd850990506000600080516020612ef78339815191528383099050610fac816115d2565b9050600080516020612ef78339815191528283099150600080516020612ef78339815191528183099150600080516020612ef78339815191528286099150600080516020612ef783398151915282600080516020612ef7833981519152037759e26bcea0d48bacd4f263f1acdb5c4f5763473177fffffe089450600080516020612ef78339815191528586099150600080516020612ef78339815191528583099150600080516020612ef7833981519152600383089150600061106e836115a9565b909350905080156110b157846110925782600080516020612ef78339815191520392505b5050604080518082019091529384526020840152509091506105139050565b600080516020612ef783398151915260018708600080516020612ef7833981519152039550600080516020612ef7833981519152806110ec57fe5b8687099250600080516020612ef78339815191528684099250600080516020612ef7833981519152600384089250611123836115a9565b90935090508015611162578461109257505060408051808201909152938452600080516020612ef7833981519152036020840152509091506105139050565b600080516020612ef78339815191528485099550600080516020612ef78339815191528687099550600080516020612ef78339815191528287099550600080516020612ef78339815191528287099550600080516020612ef7833981519152600187089550600080516020612ef78339815191528687099250600080516020612ef78339815191528684099250600080516020612ef783398151915260038408925061120d836115a9565b90935090508061122f5760405162461bcd60e51b815260040161031590612e21565b846112485782600080516020612ef78339815191520392505b505060408051808201909152938452602084015250909392505050565b61126d612481565b60008461127f57508190506001611324565b600061128b86866105fc565b905080846040015110156112a6578360029250925050611324565b868460200151146112be578360039250925050611324565b6112c6612481565b604051806080016040528086600001518152602001866020015181526020016112fc8488604001516115dd90919063ffffffff16565b8152602001611319600188606001516105fc90919063ffffffff16565b905293506000925050505b94509492505050565b8051606090816064820167ffffffffffffffff8111801561134d57600080fd5b506040519080825280601f01601f191660200182016040528015611378576020820181803683370190505b506040805160608082526080820190925291925090816020820181803683370190505090506060820160005b848110156113be57602081880181015183830152016113a4565b50830160008153600101606081536001016000815360018101879052602101602081535060006002836040516113f49190612af6565b602060405180830381855afa158015611411573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906114349190612910565b905060006042945084845281602085015260016040850153604184018890526020606185015360028460405161146a9190612af6565b602060405180830381855afa158015611487573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906114aa9190612910565b9050806020840152808218806020860152506002604085015360418401889052602060618501536002846040516114e19190612af6565b602060405180830381855afa1580156114fe573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906115219190612910565b9050806040840152808218806020860152506003604085015360418401889052602060618501536002846040516115589190612af6565b602060405180830381855afa158015611575573d6000803e3d6000fd5b5050506040513d601f19601f820116820180604052508101906115989190612910565b606084015250909695505050505050565b6000806115b583611605565b915082600080516020612ef7833981519152838409149050915091565b600061009482611d3a565b6000828211156115ff5760405162461bcd60e51b815260040161031590612c54565b50900390565b6000600080516020612ef78339815191528083840991508083830981838209828283098385830984848309858484098684850997508684840987858409945087898a09985087898a09985087898a09985087898a09985087898a09985087878a09985087898a09985087898a09985087898a099850878a8a09985087898a09985087898a09985087898a09985087898a099850878a8a09985087898a09985087898a09985087898a09985087898a09985087898a09985087848a09985087898a09985087898a09985087898a09985087898a09985087898a09985087848a09985087898a09985087898a09985087898a099850878a8a09985087898a09985087898a09985087898a09985087898a09985087848a09985087898a09985087898a09985087898a09985087898a09985087898a099850878a8a09985087898a09985087898a09985087898a09985087898a09985087878a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087818a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087868a09985087898a09985087898a09985087898a09985087898a09985087878a09985087898a09985087898a09985087898a09985087898a09985087848a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087868a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a099850878a8a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087828a09985087898a09985087898a09985087898a09985087898a09985087898a09985087818a09985087898a09985087898a09985087898a09985087868a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087878a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087868a09985087898a09985087898a09985087898a09985087878a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087828a09985087898a09985087898a09985087898a09985087898a09985087828a09985087898a09985087898a09985087898a09985087898a09985087898a09985087868a09985087898a09985087898a09985087898a09985087848a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087828a09985087898a09985087898a09985087898a09985087898a09985087868a09985087898a09985087898a09985087898a09985087898a09985087898a09985087838a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087828a09985087898a09985087898a099850878a8a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087848a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087848a09985087898a09985087898a09985087898a09985087898a09985087898a09985087868a09985087898a09985087898a099850878a8a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087818a09985050868889099750868889099750868889099750868889099750868889099750868889099750868489099750868889099750868889099750868889099750868889099750868889099750868989099750868889099750868889099750868889099750868889099750868889099750868889099750868989099750868889099750868889099750868889099750868889099750868889099750868689099750868889099750868889099750868889099750868889099750868889099750868889099750868889099750868889099750868889099750868189099750508587880996508587880996508587880996508585880996508587880996508587880996508587880996508585880996508587880996508587880996508587880996508587880996508587880996508587880996508587880996508587880996508583880996508587880996508587880996508587880996508587880996508581880996508587880996508587880996508587880996508587880996508583880996508587880996508587880996508587880996508584880996508587880996508587880996508587880996508587880996508587880996508581880996505050505050808283099392505050565b6000600080516020612ef78339815191528083840991508083830981838209828283098385830984848309858484098684850997508684840987858409945087898a09985087898a09985087898a09985087898a09985087898a09985087878a09985087898a09985087898a09985087898a099850878a8a09985087898a09985087898a09985087898a09985087898a099850878a8a09985087898a09985087898a09985087898a09985087898a09985087898a09985087848a09985087898a09985087898a09985087898a09985087898a09985087898a09985087848a09985087898a09985087898a09985087898a099850878a8a09985087898a09985087898a09985087898a09985087898a09985087848a09985087898a09985087898a09985087898a09985087898a09985087898a099850878a8a09985087898a09985087898a09985087898a09985087898a09985087878a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087818a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087868a09985087898a09985087898a09985087898a09985087898a09985087878a09985087898a09985087898a09985087898a09985087898a09985087848a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087868a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a099850878a8a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087828a09985087898a09985087898a09985087898a09985087898a09985087898a09985087818a09985087898a09985087898a09985087898a09985087868a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087878a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087868a09985087898a09985087898a09985087898a09985087878a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087828a09985087898a09985087898a09985087898a09985087898a09985087828a09985087898a09985087898a09985087898a09985087898a09985087898a09985087868a09985087898a09985087898a09985087898a09985087848a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087828a09985087898a09985087898a09985087898a09985087898a09985087868a09985087898a09985087898a09985087898a09985087898a09985087898a09985087838a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087828a09985087898a09985087898a099850878a8a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087848a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087848a09985087898a09985087898a09985087898a09985087898a09985087898a09985087868a09985087898a09985087898a099850878a8a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087898a09985087818a09985050868889099750868889099750868889099750868889099750868889099750868889099750868489099750868889099750868889099750868889099750868889099750868889099750868989099750868889099750868889099750868889099750868889099750868889099750868889099750868989099750868889099750868889099750868889099750868889099750868889099750868689099750868889099750868889099750868889099750868889099750868889099750868889099750868889099750868889099750868889099750868189099750508587880996508587880996508587880996508585880996508587880996508587880996508587880996508585880996508587880996508587880996508587880996508587880996508587880996508587880996508587880996508587880996508583880996508587880996508587880996508587880996508587880996508581880996505050838586099450838586099450838586099450838586099450838186099450508284850993508284850993508284850993508281850993508284850993508284850993508285850993508284850993508284850993508284850993508284850993508284850993508284850993508281850995945050505050565b6040518060800160405280600081526020016000815260200160008152602001600081525090565b60405180604001604052806002906020820280368337509192915050565b60405180608001604052806004906020820280368337509192915050565b60405180602001604052806001906020820280368337509192915050565b600082601f830112612513578081fd5b813561252661252182612ed6565b612eaf565b818152915060208083019084810160005b848110156125605761254e888484358a0101612622565b84529282019290820190600101612537565b505050505092915050565b6000601f838184011261257c578182fd5b823561258a61252182612ed6565b81815292506020808401908581016080808502880183018910156125ad57600080fd5b60005b858110156126155789878401126125c657600080fd5b6125cf82612eaf565b80848486018d8111156125e157600080fd5b60005b6004811015612601578235855293880193918801916001016125e4565b5091885250958501959350506001016125b0565b5050505050505092915050565b600082601f830112612632578081fd5b813561264061252182612ed6565b81815291506020808301908481018184028601820187101561266157600080fd5b60005b8481101561256057813584529282019290820190600101612664565b600082601f830112612690578081fd5b813561269e61252182612ed6565b81815291506020808301908481016080808502870183018810156126c157600080fd5b60005b858110156126e8576126d689846128cb565b855293830193918101916001016126c4565b50505050505092915050565b600082601f830112612704578081fd5b61270e6040612eaf565b905080828460408501111561272257600080fd5b60005b6002811015612744578135835260209283019290910190600101612725565b50505092915050565b600082601f83011261275d578081fd5b813567ffffffffffffffff811115612773578182fd5b612786601f8201601f1916602001612eaf565b915080825283602082850101111561279d57600080fd5b8060208401602084013760009082016020015292915050565b6000608082840312156127c7578081fd5b6127d16080612eaf565b9050813567ffffffffffffffff808211156127eb57600080fd5b6127f785838601612680565b8352602084013591508082111561280d57600080fd5b61281985838601612503565b6020840152604084013591508082111561283257600080fd5b61283e8583860161256b565b6040840152606084013591508082111561285757600080fd5b5061286484828501612503565b60608301525092915050565b600060a08284031215612881578081fd5b61288b6040612eaf565b905061289783836128cb565b8152608082013567ffffffffffffffff8111156128b357600080fd5b6128bf84828501612622565b60208301525092915050565b6000608082840312156128dc578081fd5b6128e66080612eaf565b90508135815260208201356020820152604082013560408201526060820135606082015292915050565b600060208284031215612921578081fd5b5051919050565b600080600080600060a0868803121561293f578081fd5b85359450602080870135945060408701359350606087013567ffffffffffffffff8082111561296c578384fd5b6129788a838b0161274d565b9450608089013591508082111561298d578384fd5b508701601f8101891361299e578283fd5b80356129ac61252182612ed6565b81815283810190838501865b848110156129e1576129cf8e888435890101612870565b845292860192908601906001016129b8565b505080955050505050509295509295909350565b60008060408385031215612a07578182fd5b823567ffffffffffffffff80821115612a1e578384fd5b9084019060c08287031215612a31578384fd5b612a3b60a0612eaf565b612a4587846126f4565b815260408301356020820152606083013560408201526080830135606082015260a083013582811115612a76578586fd5b612a828882860161274d565b60808301525093506020850135915080821115612a9d578283fd5b50612aaa858286016127b6565b9150509250929050565b600b8110612abe57fe5b9052565b60008183825b6004811015612ae7578151835260209283019290910190600101612ac8565b50505060808201905092915050565b60008251815b81811015612b165760208186018101518583015201612afc565b81811115612b245782828501525b509190910192915050565b93845260208401929092526040830152606082015260800190565b958652602086019490945260408501929092526060840152608083015260a082015260c00190565b918252602082015260400190565b82815260408101610e276020830184612ab4565b602081016100948284612ab4565b60208082526017908201527f424c533a20626e206164642063616c6c206661696c6564000000000000000000604082015260600190565b60208082526024908201527f41757468656e7469636974793a20737461746520696e636c7573696f6e20736960408201526333b732b960e11b606082015260800190565b6020808252601b908201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604082015260600190565b6020808252601e908201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604082015260600190565b60208082526035908201527f424c533a206e756d626572206f66207075626c6963206b65797320616e64206d604082015274195cdcd859d95cc81b5d5cdd08189948195c5d585b605a1b606082015260800190565b60208082526018908201527f41757468656e7469636974793a207a65726f206e6f6e63650000000000000000604082015260600190565b60208082526021908201527f5472616e736974696f6e3a2053656e64657220646f6573206e6f7420657869736040820152601d60fa1b606082015260800190565b60208082526021908201527f424c533a206e756d626572206f66207075626c6963206b6579206973207a65726040820152606f60f81b606082015260800190565b60208082526023908201527f6d6170546f506f696e7446543a20696e76616c6964206669656c6420656c656d604082015262195b9d60ea1b606082015260800190565b60208082526025908201527f41757468656e7469636974793a206163636f756e7420646f6573206e6f742065604082015264786973747360d81b606082015260800190565b60208082526022908201527f424c533a20626164206674206d617070696e6720696d706c656d656e7461746960408201526137b760f11b606082015260800190565b60208082526023908201527f5472616e736974696f6e3a20726563656976657220646f6573206e6f742065786040820152621a5cdd60ea1b606082015260800190565b90815260200190565b60405181810167ffffffffffffffff81118282101715612ece57600080fd5b604052919050565b600067ffffffffffffffff821115612eec578081fd5b506020908102019056fe30644e72e131a029b85045b68181585d97816a916871ca8d3c208c16d87cfd47a164736f6c634300060c000a"

// DeployTransfer deploys a new Ethereum contract, binding an instance of Transfer to it.
func DeployTransfer(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Transfer, error) {
	parsed, err := abi.JSON(strings.NewReader(TransferABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TransferBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Transfer{TransferCaller: TransferCaller{contract: contract}, TransferTransactor: TransferTransactor{contract: contract}, TransferFilterer: TransferFilterer{contract: contract}}, nil
}

// Transfer is an auto generated Go binding around an Ethereum contract.
type Transfer struct {
	TransferCaller     // Read-only binding to the contract
	TransferTransactor // Write-only binding to the contract
	TransferFilterer   // Log filterer for contract events
}

// TransferCaller is an auto generated read-only Go binding around an Ethereum contract.
type TransferCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TransferTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TransferFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TransferSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TransferSession struct {
	Contract     *Transfer         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TransferCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TransferCallerSession struct {
	Contract *TransferCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// TransferTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TransferTransactorSession struct {
	Contract     *TransferTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// TransferRaw is an auto generated low-level Go binding around an Ethereum contract.
type TransferRaw struct {
	Contract *Transfer // Generic contract binding to access the raw methods on
}

// TransferCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TransferCallerRaw struct {
	Contract *TransferCaller // Generic read-only contract binding to access the raw methods on
}

// TransferTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TransferTransactorRaw struct {
	Contract *TransferTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTransfer creates a new instance of Transfer, bound to a specific deployed contract.
func NewTransfer(address common.Address, backend bind.ContractBackend) (*Transfer, error) {
	contract, err := bindTransfer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Transfer{TransferCaller: TransferCaller{contract: contract}, TransferTransactor: TransferTransactor{contract: contract}, TransferFilterer: TransferFilterer{contract: contract}}, nil
}

// NewTransferCaller creates a new read-only instance of Transfer, bound to a specific deployed contract.
func NewTransferCaller(address common.Address, caller bind.ContractCaller) (*TransferCaller, error) {
	contract, err := bindTransfer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TransferCaller{contract: contract}, nil
}

// NewTransferTransactor creates a new write-only instance of Transfer, bound to a specific deployed contract.
func NewTransferTransactor(address common.Address, transactor bind.ContractTransactor) (*TransferTransactor, error) {
	contract, err := bindTransfer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TransferTransactor{contract: contract}, nil
}

// NewTransferFilterer creates a new log filterer instance of Transfer, bound to a specific deployed contract.
func NewTransferFilterer(address common.Address, filterer bind.ContractFilterer) (*TransferFilterer, error) {
	contract, err := bindTransfer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TransferFilterer{contract: contract}, nil
}

// bindTransfer binds a generic wrapper to an already deployed contract.
func bindTransfer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TransferABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Transfer *TransferRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Transfer.Contract.TransferCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Transfer *TransferRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Transfer.Contract.TransferTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Transfer *TransferRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Transfer.Contract.TransferTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Transfer *TransferCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Transfer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Transfer *TransferTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Transfer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Transfer *TransferTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Transfer.Contract.contract.Transact(opts, method, params...)
}

// CheckSignature is a free data retrieval call binding the contract method 0x44ec5a77.
//
// Solidity: function checkSignature((uint256[2],bytes32,bytes32,bytes32,bytes) common, ((uint256,uint256,uint256,uint256)[],bytes32[][],uint256[4][],bytes32[][]) proof) view returns(uint8)
func (_Transfer *TransferCaller) CheckSignature(opts *bind.CallOpts, common TypesAuthCommon, proof TypesSignatureProof) (uint8, error) {
	var out []interface{}
	err := _Transfer.contract.Call(opts, &out, "checkSignature", common, proof)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// CheckSignature is a free data retrieval call binding the contract method 0x44ec5a77.
//
// Solidity: function checkSignature((uint256[2],bytes32,bytes32,bytes32,bytes) common, ((uint256,uint256,uint256,uint256)[],bytes32[][],uint256[4][],bytes32[][]) proof) view returns(uint8)
func (_Transfer *TransferSession) CheckSignature(common TypesAuthCommon, proof TypesSignatureProof) (uint8, error) {
	return _Transfer.Contract.CheckSignature(&_Transfer.CallOpts, common, proof)
}

// CheckSignature is a free data retrieval call binding the contract method 0x44ec5a77.
//
// Solidity: function checkSignature((uint256[2],bytes32,bytes32,bytes32,bytes) common, ((uint256,uint256,uint256,uint256)[],bytes32[][],uint256[4][],bytes32[][]) proof) view returns(uint8)
func (_Transfer *TransferCallerSession) CheckSignature(common TypesAuthCommon, proof TypesSignatureProof) (uint8, error) {
	return _Transfer.Contract.CheckSignature(&_Transfer.CallOpts, common, proof)
}

// ProcessTransferCommit is a free data retrieval call binding the contract method 0x92931492.
//
// Solidity: function processTransferCommit(bytes32 stateRoot, uint256 maxTxSize, uint256 feeReceiver, bytes txs, ((uint256,uint256,uint256,uint256),bytes32[])[] proofs) pure returns(bytes32, uint8 result)
func (_Transfer *TransferCaller) ProcessTransferCommit(opts *bind.CallOpts, stateRoot [32]byte, maxTxSize *big.Int, feeReceiver *big.Int, txs []byte, proofs []TypesStateMerkleProof) ([32]byte, uint8, error) {
	var out []interface{}
	err := _Transfer.contract.Call(opts, &out, "processTransferCommit", stateRoot, maxTxSize, feeReceiver, txs, proofs)

	if err != nil {
		return *new([32]byte), *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	out1 := *abi.ConvertType(out[1], new(uint8)).(*uint8)

	return out0, out1, err

}

// ProcessTransferCommit is a free data retrieval call binding the contract method 0x92931492.
//
// Solidity: function processTransferCommit(bytes32 stateRoot, uint256 maxTxSize, uint256 feeReceiver, bytes txs, ((uint256,uint256,uint256,uint256),bytes32[])[] proofs) pure returns(bytes32, uint8 result)
func (_Transfer *TransferSession) ProcessTransferCommit(stateRoot [32]byte, maxTxSize *big.Int, feeReceiver *big.Int, txs []byte, proofs []TypesStateMerkleProof) ([32]byte, uint8, error) {
	return _Transfer.Contract.ProcessTransferCommit(&_Transfer.CallOpts, stateRoot, maxTxSize, feeReceiver, txs, proofs)
}

// ProcessTransferCommit is a free data retrieval call binding the contract method 0x92931492.
//
// Solidity: function processTransferCommit(bytes32 stateRoot, uint256 maxTxSize, uint256 feeReceiver, bytes txs, ((uint256,uint256,uint256,uint256),bytes32[])[] proofs) pure returns(bytes32, uint8 result)
func (_Transfer *TransferCallerSession) ProcessTransferCommit(stateRoot [32]byte, maxTxSize *big.Int, feeReceiver *big.Int, txs []byte, proofs []TypesStateMerkleProof) ([32]byte, uint8, error) {
	return _Transfer.Contract.ProcessTransferCommit(&_Transfer.CallOpts, stateRoot, maxTxSize, feeReceiver, txs, proofs)
}
