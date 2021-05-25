// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rollup

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

// TypesBatch is an auto generated low-level Go binding around an user-defined struct.
type TypesBatch struct {
	CommitmentRoot [32]byte
	Meta           [32]byte
}

// TypesCommitment is an auto generated low-level Go binding around an user-defined struct.
type TypesCommitment struct {
	StateRoot [32]byte
	BodyRoot  [32]byte
}

// TypesCommitmentInclusionProof is an auto generated low-level Go binding around an user-defined struct.
type TypesCommitmentInclusionProof struct {
	Commitment TypesCommitment
	Path       *big.Int
	Witness    [][32]byte
}

// TypesMMCommitmentInclusionProof is an auto generated low-level Go binding around an user-defined struct.
type TypesMMCommitmentInclusionProof struct {
	Commitment TypesMassMigrationCommitment
	Path       *big.Int
	Witness    [][32]byte
}

// TypesMassMigrationBody is an auto generated low-level Go binding around an user-defined struct.
type TypesMassMigrationBody struct {
	AccountRoot  [32]byte
	Signature    [2]*big.Int
	SpokeID      *big.Int
	WithdrawRoot [32]byte
	TokenID      *big.Int
	Amount       *big.Int
	FeeReceiver  *big.Int
	Txs          []byte
}

// TypesMassMigrationCommitment is an auto generated low-level Go binding around an user-defined struct.
type TypesMassMigrationCommitment struct {
	StateRoot [32]byte
	Body      TypesMassMigrationBody
}

// TypesSignatureProof is an auto generated low-level Go binding around an user-defined struct.
type TypesSignatureProof struct {
	States          []TypesUserState
	StateWitnesses  [][][32]byte
	Pubkeys         [][4]*big.Int
	PubkeyWitnesses [][][32]byte
}

// TypesSignatureProofWithReceiver is an auto generated low-level Go binding around an user-defined struct.
type TypesSignatureProofWithReceiver struct {
	States                  []TypesUserState
	StateWitnesses          [][][32]byte
	PubkeysSender           [][4]*big.Int
	PubkeyWitnessesSender   [][][32]byte
	PubkeysReceiver         [][4]*big.Int
	PubkeyWitnessesReceiver [][][32]byte
}

// TypesStateMerkleProof is an auto generated low-level Go binding around an user-defined struct.
type TypesStateMerkleProof struct {
	State   TypesUserState
	Witness [][32]byte
}

// TypesSubtreeVacancyProof is an auto generated low-level Go binding around an user-defined struct.
type TypesSubtreeVacancyProof struct {
	Depth       *big.Int
	PathAtDepth *big.Int
	Witness     [][32]byte
}

// TypesTransferBody is an auto generated low-level Go binding around an user-defined struct.
type TypesTransferBody struct {
	AccountRoot [32]byte
	Signature   [2]*big.Int
	FeeReceiver *big.Int
	Txs         []byte
}

// TypesTransferCommitment is an auto generated low-level Go binding around an user-defined struct.
type TypesTransferCommitment struct {
	StateRoot [32]byte
	Body      TypesTransferBody
}

// TypesTransferCommitmentInclusionProof is an auto generated low-level Go binding around an user-defined struct.
type TypesTransferCommitmentInclusionProof struct {
	Commitment TypesTransferCommitment
	Path       *big.Int
	Witness    [][32]byte
}

// TypesUserState is an auto generated low-level Go binding around an user-defined struct.
type TypesUserState struct {
	PubkeyID *big.Int
	TokenID  *big.Int
	Balance  *big.Int
	Nonce    *big.Int
}

// RollupABI is the input ABI used to generate the binding from.
const RollupABI = "[{\"inputs\":[{\"internalType\":\"contractChooser\",\"name\":\"_chooser\",\"type\":\"address\"},{\"internalType\":\"contractIDepositManager\",\"name\":\"_depositManager\",\"type\":\"address\"},{\"internalType\":\"contractBLSAccountRegistry\",\"name\":\"_accountRegistry\",\"type\":\"address\"},{\"internalType\":\"contractTransfer\",\"name\":\"_transfer\",\"type\":\"address\"},{\"internalType\":\"contractMassMigration\",\"name\":\"_massMigration\",\"type\":\"address\"},{\"internalType\":\"contractCreate2Transfer\",\"name\":\"_create2Transfer\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"genesisStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"stakeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blocksToFinalise\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minGasLeft\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxTxsPerCommit\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subtreeID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"depositSubTreeRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"pathToSubTree\",\"type\":\"uint256\"}],\"name\":\"DepositsFinalised\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumTypes.Usage\",\"name\":\"batchType\",\"type\":\"uint8\"}],\"name\":\"NewBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nDeleted\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"completed\",\"type\":\"bool\"}],\"name\":\"RollbackStatus\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"committed\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"}],\"name\":\"StakeWithdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ZERO_BYTES32\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"accountRegistry\",\"outputs\":[{\"internalType\":\"contractBLSAccountRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"appID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"batches\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitmentRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"meta\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"chooser\",\"outputs\":[{\"internalType\":\"contractChooser\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"create2Transfer\",\"outputs\":[{\"internalType\":\"contractCreate2Transfer\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositManager\",\"outputs\":[{\"internalType\":\"contractIDepositManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"deposits\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.TransferBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"internalType\":\"structTypes.TransferCommitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.TransferCommitmentInclusionProof\",\"name\":\"target\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState[]\",\"name\":\"states\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"stateWitnesses\",\"type\":\"bytes32[][]\"},{\"internalType\":\"uint256[4][]\",\"name\":\"pubkeysSender\",\"type\":\"uint256[4][]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"pubkeyWitnessesSender\",\"type\":\"bytes32[][]\"},{\"internalType\":\"uint256[4][]\",\"name\":\"pubkeysReceiver\",\"type\":\"uint256[4][]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"pubkeyWitnessesReceiver\",\"type\":\"bytes32[][]\"}],\"internalType\":\"structTypes.SignatureProofWithReceiver\",\"name\":\"signatureProof\",\"type\":\"tuple\"}],\"name\":\"disputeSignatureCreate2Transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"spokeID\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"withdrawRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.MassMigrationBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"internalType\":\"structTypes.MassMigrationCommitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.MMCommitmentInclusionProof\",\"name\":\"target\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState[]\",\"name\":\"states\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"stateWitnesses\",\"type\":\"bytes32[][]\"},{\"internalType\":\"uint256[4][]\",\"name\":\"pubkeys\",\"type\":\"uint256[4][]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"pubkeyWitnesses\",\"type\":\"bytes32[][]\"}],\"internalType\":\"structTypes.SignatureProof\",\"name\":\"signatureProof\",\"type\":\"tuple\"}],\"name\":\"disputeSignatureMassMigration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.TransferBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"internalType\":\"structTypes.TransferCommitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.TransferCommitmentInclusionProof\",\"name\":\"target\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState[]\",\"name\":\"states\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"stateWitnesses\",\"type\":\"bytes32[][]\"},{\"internalType\":\"uint256[4][]\",\"name\":\"pubkeys\",\"type\":\"uint256[4][]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"pubkeyWitnesses\",\"type\":\"bytes32[][]\"}],\"internalType\":\"structTypes.SignatureProof\",\"name\":\"signatureProof\",\"type\":\"tuple\"}],\"name\":\"disputeSignatureTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"bodyRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structTypes.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.CommitmentInclusionProof\",\"name\":\"previous\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.TransferBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"internalType\":\"structTypes.TransferCommitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.TransferCommitmentInclusionProof\",\"name\":\"target\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState\",\"name\":\"state\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.StateMerkleProof[]\",\"name\":\"proofs\",\"type\":\"tuple[]\"}],\"name\":\"disputeTransitionCreate2Transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"bodyRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structTypes.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.CommitmentInclusionProof\",\"name\":\"previous\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"spokeID\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"withdrawRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.MassMigrationBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"internalType\":\"structTypes.MassMigrationCommitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.MMCommitmentInclusionProof\",\"name\":\"target\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState\",\"name\":\"state\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.StateMerkleProof[]\",\"name\":\"proofs\",\"type\":\"tuple[]\"}],\"name\":\"disputeTransitionMassMigration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"bodyRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structTypes.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.CommitmentInclusionProof\",\"name\":\"previous\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.TransferBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"internalType\":\"structTypes.TransferCommitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.TransferCommitmentInclusionProof\",\"name\":\"target\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState\",\"name\":\"state\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.StateMerkleProof[]\",\"name\":\"proofs\",\"type\":\"tuple[]\"}],\"name\":\"disputeTransitionTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"}],\"name\":\"getBatch\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"commitmentRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"meta\",\"type\":\"bytes32\"}],\"internalType\":\"structTypes.Batch\",\"name\":\"batch\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"invalidBatchMarker\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"keepRollingBack\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"massMigration\",\"outputs\":[{\"internalType\":\"contractMassMigration\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextBatchID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paramBlocksToFinalise\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paramMaxTxsPerCommit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paramMinGasLeft\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paramStakeAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"stateRoots\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[2][]\",\"name\":\"signatures\",\"type\":\"uint256[2][]\"},{\"internalType\":\"uint256[]\",\"name\":\"feeReceivers\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"txss\",\"type\":\"bytes[]\"}],\"name\":\"submitCreate2Transfer\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"bodyRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structTypes.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.CommitmentInclusionProof\",\"name\":\"previous\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"depth\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"pathAtDepth\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.SubtreeVacancyProof\",\"name\":\"vacant\",\"type\":\"tuple\"}],\"name\":\"submitDeposits\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"stateRoots\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[2][]\",\"name\":\"signatures\",\"type\":\"uint256[2][]\"},{\"internalType\":\"uint256[4][]\",\"name\":\"meta\",\"type\":\"uint256[4][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"withdrawRoots\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes[]\",\"name\":\"txss\",\"type\":\"bytes[]\"}],\"name\":\"submitMassMigration\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"stateRoots\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[2][]\",\"name\":\"signatures\",\"type\":\"uint256[2][]\"},{\"internalType\":\"uint256[]\",\"name\":\"feeReceivers\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"txss\",\"type\":\"bytes[]\"}],\"name\":\"submitTransfer\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"contractTransfer\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"}],\"name\":\"withdrawStake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"withdrawalBitmap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// RollupBin is the compiled bytecode used for deploying new contracts.
var RollupBin = "0x608060405267016345785d8a00006000556000600155612710600255602060035560006007553480156200003257600080fd5b50604051620049ff380380620049ff83398101604081905262000055916200027d565b600580546001600160a01b03199081166001600160a01b038e8116919091179092556004805482168d8416179055600b805482168c8416179055600c805482168b8416179055600d805482168a8416179055600e80549091169188169190911790556000848155600184905560028390556003829055604051620000ef908790600080516020620049df833981519152906020016200035d565b60408051601f1981840301815290829052805160209182012092506000916200012c918491600080516020620049df83398151915291016200035d565b60408051601f19818403018152828252805160209182012083830190925281835290925081016200016d6000600133436200021e60201b62001d821760201c565b9052600780546000908152600660209081526040808320855181559490910151600190940193909355905491517f3c367559be37298dc40fed468321ed4d44d99ab593f8e0bb7c82dbc84ab0bad692620001cb92909181906200036b565b60405180910390a1600780546001019055604051620001ef90309060200162000345565b60405160208183030381529060405280519060200120600f8190555050505050505050505050505050620003ab565b69ffffffff000000000000603082901b16600160501b600160f01b03605084901b1660ff60f01b60f086901b167fff0000000000000000000000000000000000000000000000000000000000000060f888901b16171717949350505050565b60008060008060008060008060008060006101608c8e0312156200029f578687fd5b8b51620002ac8162000392565b60208d0151909b50620002bf8162000392565b60408d0151909a50620002d28162000392565b60608d0151909950620002e58162000392565b60808d0151909850620002f88162000392565b60a08d01519097506200030b8162000392565b8096505060c08c0151945060e08c015193506101008c015192506101208c015191506101408c015190509295989b509295989b9093969950565b60609190911b6001600160601b031916815260140190565b918252602082015260400190565b8381526020810183905260608101600583106200038457fe5b826040830152949350505050565b6001600160a01b0381168114620003a857600080fd5b50565b61462480620003bb6000396000f3fe6080604052600436106101cd5760003560e01c80638a4068dd116100f7578063c08caaef11610095578063e42f755411610064578063e42f75541461048f578063eb84e076146104a4578063f9a6cbae146104c4578063fe5f9604146104d9576101cd565b8063c08caaef14610430578063ca858d9d14610450578063d089e11a14610465578063d1b243e21461047a576101cd565b8063a72598ce116100d1578063a72598ce146103b8578063ac96f0cd146103cd578063b02c43d0146103e2578063b32c4d8d14610402576101cd565b80638a4068dd1461037b5780639251597f1461039057806398d17621146103a3576101cd565b806350980f401161016f5780635f5b95b81161013e5780635f5b95b81461031e5780635f6e91d5146103315780636c7ac9d8146103515780637ae8c56814610366576101cd565b806350980f40146102a9578063562a2bca146102c95780635ac44282146102dc5780635b097d3714610309576101cd565b80632538507d116101ab5780632538507d1461023457806325d5971f1461025457806331c2b7db146102745780634e23e8c314610289576101cd565b8063035695e3146101d2578063069321b0146101e75780630ed75b9c14610212575b600080fd5b6101e56101e0366004613447565b6104f9565b005b3480156101f357600080fd5b506101fc6108a0565b6040516102099190613da0565b60405180910390f35b34801561021e57600080fd5b506102276108b2565b6040516102099190613e8c565b34801561024057600080fd5b506101e561024f36600461381d565b6108c1565b34801561026057600080fd5b506101e561026f3660046136f6565b610a6d565b34801561028057600080fd5b506101fc610bb8565b34801561029557600080fd5b506101e56102a43660046137be565b610bbe565b3480156102b557600080fd5b506101e56102c4366004613886565b610e12565b6101e56102d7366004613652565b610f42565b3480156102e857600080fd5b506102fc6102f73660046136f6565b61126e565b60405161020991906144ca565b34801561031557600080fd5b506101fc6112ce565b6101e561032c366004613531565b6112d4565b34801561033d57600080fd5b506101e561034c3660046137be565b611572565b34801561035d57600080fd5b50610227611747565b34801561037257600080fd5b50610227611756565b34801561038757600080fd5b50610227611765565b6101e561039e366004613531565b611774565b3480156103af57600080fd5b50610227611a05565b3480156103c457600080fd5b506101fc611a14565b3480156103d957600080fd5b506101e5611a1a565b3480156103ee57600080fd5b506101fc6103fd3660046136f6565b611a46565b34801561040e57600080fd5b5061042261041d3660046136f6565b611a58565b604051610209929190613d79565b34801561043c57600080fd5b506101e561044b3660046139ce565b611a71565b34801561045c57600080fd5b506101fc611ba1565b34801561047157600080fd5b50610227611ba7565b34801561048657600080fd5b506101fc611bb6565b34801561049b57600080fd5b506101fc611bbc565b3480156104b057600080fd5b506101e56104bf366004613731565b611bc2565b3480156104d057600080fd5b506101fc611d6a565b3480156104e557600080fd5b506101fc6104f43660046136f6565b611d70565b600560009054906101000a90046001600160a01b03166001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b15801561054757600080fd5b505afa15801561055b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061057f9190613420565b6001600160a01b0316336001600160a01b0316146105b85760405162461bcd60e51b81526004016105af90614234565b60405180910390fd5b600a54156105d85760405162461bcd60e51b81526004016105af9061410b565b6060896001600160401b03811180156105f057600080fd5b5060405190808252806020026020018201604052801561061a578160200160208202803683370190505b5090506000600b60009054906101000a90046001600160a01b03166001600160a01b031663ebf0c7176040518163ffffffff1660e01b815260040160206040518083038186803b15801561066d57600080fd5b505afa158015610681573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906106a591906135ef565b905060005b8b81101561087c576106ba612b51565b6040518061010001604052808481526020018d8d858181106106d857fe5b9050604002016002806020026040519081016040528092919082600260200280828437600092019190915250505081526020018b8b8581811061071757fe5b90506080020160006004811061072957fe5b6020020135815260200189898581811061073f57fe5b9050602002013581526020018b8b8581811061075757fe5b90506080020160016004811061076957fe5b602002013581526020018b8b8581811061077f57fe5b90506080020160026004811061079157fe5b602002013581526020018b8b858181106107a757fe5b9050608002016003600481106107b957fe5b602002013581526020018787858181106107cf57fe5b90506020028101906107e19190614535565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250505091525090508d8d8381811061082557fe5b9050602002013561083582611dc9565b604051602001610846929190613d79565b6040516020818303038152906040528051906020012084838151811061086857fe5b6020908102919091010152506001016106aa565b5061089261088983611e20565b8c8360026120b6565b505050505050505050505050565b6000805160206145f883398151915281565b600d546001600160a01b031681565b600083815260066020908152604091829020825180840190935280548352600101549082015283906108f290612175565b43106109105760405162461bcd60e51b81526004016105af90614142565b600a548110806109205750600a54155b61093c5760405162461bcd60e51b81526004016105af906142ff565b6000848152600660205260409020546109559084612185565b6109715760405162461bcd60e51b81526004016105af90614179565b610979612ba0565b506040805160a081018252845160209081015181015182528551518183015285518101515182840152600f546060830152855181015160e001516080830152600d548651909101518301519251632f90b1f160e21b815291926000926001600160a01b039092169163be42c7c4916109f79186918991600401614494565b60206040518083038186803b158015610a0f57600080fd5b505afa158015610a23573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a479190613636565b9050600081600a811115610a5757fe5b14610a6557610a65866121a1565b505050505050565b6000818152600660209081526040918290208251808401909352805483526001015490820152610a9c906121b1565b6001600160a01b0316336001600160a01b031614610acc5760405162461bcd60e51b81526004016105af90613f99565b6000818152600660209081526040918290208251808401909352805483526001015490820152610afb90612175565b4311610b195760405162461bcd60e51b81526004016105af906140ba565b610b248160096121c7565b15610b415760405162461bcd60e51b81526004016105af90614352565b610b4c8160096121ee565b60008054604051339282156108fc02929190818181858888f19350505050158015610b7b573d6000803e3d6000fd5b507f1248d48e2de900a1010c7fce73506969ecec243600bfc08b641b158f26d857cd3382604051610bad929190613d87565b60405180910390a150565b60015481565b60008481526006602090815260409182902082518084019093528054835260010154908201528490610bef90612175565b4310610c0d5760405162461bcd60e51b81526004016105af90614142565b600a54811080610c1d5750600a54155b610c395760405162461bcd60e51b81526004016105af906142ff565b848484602001516000808260001415610c8f575060001984016000818152600660209081526040918290208251808401909352805483526001908101549183019190915290610c8790612216565b039150610c98565b50506000198101835b81846020015114610cbb5760405162461bcd60e51b81526004016105af90614078565b600081815260066020526040902054610cd49085612223565b610cf05760405162461bcd60e51b81526004016105af90613ef4565b60008a815260066020526040902054610d099089612236565b610d255760405162461bcd60e51b81526004016105af90613fe9565b600080600c60009054906101000a90046001600160a01b03166001600160a01b031663929314928c60000151600001516003548d6000015160200151604001518e6000015160200151606001518e6040518663ffffffff1660e01b8152600401610d93959493929190613e49565b604080518083038186803b158015610daa57600080fd5b505afa158015610dbe573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610de29190613607565b9092509050600081600a811115610df557fe5b141580610e0457508951518214155b15610892576108928c6121a1565b60008381526006602090815260409182902082518084019093528054835260010154908201528390610e4390612175565b4310610e615760405162461bcd60e51b81526004016105af90614142565b600a54811080610e715750600a54155b610e8d5760405162461bcd60e51b81526004016105af906142ff565b600084815260066020526040902054610ea69084612236565b610ec25760405162461bcd60e51b81526004016105af90614031565b610eca612ba0565b506040805160a081018252845160209081015181015182528551518183015285518101515182840152600f5460608084019190915286519091015101516080820152600e549151639c57ceb560e01b815290916000916001600160a01b0390911690639c57ceb5906109f790859088906004016143c0565b600560009054906101000a90046001600160a01b03166001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b158015610f9057600080fd5b505afa158015610fa4573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610fc89190613420565b6001600160a01b0316336001600160a01b031614610ff85760405162461bcd60e51b81526004016105af90614234565b600a54156110185760405162461bcd60e51b81526004016105af9061410b565b60075460001901600081815260066020908152604091829020825180840190935280548352600190810154918301919091529061105490612216565b038360200151146110775760405162461bcd60e51b81526004016105af90614078565b6000818152600660205260409020546110909084612223565b6110ac5760405162461bcd60e51b81526004016105af90613f46565b82515182516110ce91906110bf90612249565b84602001518560400151612732565b6110ea5760405162461bcd60e51b81526004016105af906141f1565b600480546040805163d86ee48d60e01b8152815160009485946001600160a01b03169363d86ee48d938083019391929082900301818787803b15801561112f57600080fd5b505af1158015611143573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611167919061370e565b600185016000818152600860209081526040918290208490558801519051939550919350917f1e6cd0ae9aa6f6e6ea5cf7f07396a96fc1ea5ae4ec7d5054f53e204c47f79bd8916111bb9186918691614507565b60405180910390a160006111d8838760200151886040015161274a565b90506000816000805160206145f883398151915260001b604051602001611200929190613d79565b60408051601f19818403018152908290528051602091820120925060009161123a9184916000805160206145f88339815191529101613d79565b6040516020818303038152906040528051906020012090506112638160016000801b60046120b6565b505050505050505050565b611276612bd7565b6000828152600660205260409020600101546112a45760405162461bcd60e51b81526004016105af9061426b565b50600090815260066020908152604091829020825180840190935280548352600101549082015290565b600a5481565b600560009054906101000a90046001600160a01b03166001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b15801561132257600080fd5b505afa158015611336573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061135a9190613420565b6001600160a01b0316336001600160a01b03161461138a5760405162461bcd60e51b81526004016105af90614234565b600a54156113aa5760405162461bcd60e51b81526004016105af9061410b565b6060876001600160401b03811180156113c257600080fd5b506040519080825280602002602001820160405280156113ec578160200160208202803683370190505b5090506000600b60009054906101000a90046001600160a01b03166001600160a01b031663ebf0c7176040518163ffffffff1660e01b815260040160206040518083038186803b15801561143f57600080fd5b505afa158015611453573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061147791906135ef565b90506000805b8a81101561154f57828a8a8381811061149257fe5b9050604002018989848181106114a457fe5b905060200201358888858181106114b757fe5b90506020028101906114c99190614535565b6040516020016114dd959493929190613cba565b6040516020818303038152906040528051906020012091508b8b8281811061150157fe5b905060200201358260405160200161151a929190613d79565b6040516020818303038152906040528051906020012084828151811061153c57fe5b602090810291909101015260010161147d565b5061156561155c84611e20565b8b8460016120b6565b5050505050505050505050565b600084815260066020908152604091829020825180840190935280548352600101549082015284906115a390612175565b43106115c15760405162461bcd60e51b81526004016105af90614142565b600a548110806115d15750600a54155b6115ed5760405162461bcd60e51b81526004016105af906142ff565b84848460200151600080826000141561164357506000198401600081815260066020908152604091829020825180840190935280548352600190810154918301919091529061163b90612216565b03915061164c565b50506000198101835b8184602001511461166f5760405162461bcd60e51b81526004016105af90614078565b6000818152600660205260409020546116889085612223565b6116a45760405162461bcd60e51b81526004016105af90613ef4565b60008a8152600660205260409020546116bd9089612236565b6116d95760405162461bcd60e51b81526004016105af90613fe9565b600080600e60009054906101000a90046001600160a01b03166001600160a01b031663336920368c60000151600001516003548d6000015160200151604001518e6000015160200151606001518e6040518663ffffffff1660e01b8152600401610d93959493929190613e49565b6004546001600160a01b031681565b600e546001600160a01b031681565b600c546001600160a01b031681565b600560009054906101000a90046001600160a01b03166001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b1580156117c257600080fd5b505afa1580156117d6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117fa9190613420565b6001600160a01b0316336001600160a01b03161461182a5760405162461bcd60e51b81526004016105af90614234565b600a541561184a5760405162461bcd60e51b81526004016105af9061410b565b6060876001600160401b038111801561186257600080fd5b5060405190808252806020026020018201604052801561188c578160200160208202803683370190505b5090506000600b60009054906101000a90046001600160a01b03166001600160a01b031663ebf0c7176040518163ffffffff1660e01b815260040160206040518083038186803b1580156118df57600080fd5b505afa1580156118f3573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061191791906135ef565b90506000805b8a8110156119ef57828a8a8381811061193257fe5b90506040020189898481811061194457fe5b9050602002013588888581811061195757fe5b90506020028101906119699190614535565b60405160200161197d959493929190613cba565b6040516020818303038152906040528051906020012091508b8b828181106119a157fe5b90506020020135826040516020016119ba929190613d79565b604051602081830303815290604052805190602001208482815181106119dc57fe5b602090810291909101015260010161191d565b506115656119fc84611e20565b8b8460036120b6565b6005546001600160a01b031681565b60075481565b6000600a5411611a3c5760405162461bcd60e51b81526004016105af90613eb3565b611a446127f5565b565b60086020526000908152604090205481565b6006602052600090815260409020805460019091015482565b60008381526006602090815260409182902082518084019093528054835260010154908201528390611aa290612175565b4310611ac05760405162461bcd60e51b81526004016105af90614142565b600a54811080611ad05750600a54155b611aec5760405162461bcd60e51b81526004016105af906142ff565b600084815260066020526040902054611b059084612236565b611b215760405162461bcd60e51b81526004016105af90614031565b611b29612ba0565b506040805160a081018252845160209081015181015182528551518183015285518101515182840152600f5460608084019190915286519091015101516080820152600c5491516344ec5a7760e01b815290916000916001600160a01b03909116906344ec5a77906109f7908590889060040161446f565b60025481565b600b546001600160a01b031681565b60005481565b60035481565b60008481526006602090815260409182902082518084019093528054835260010154908201528490611bf390612175565b4310611c115760405162461bcd60e51b81526004016105af90614142565b600a54811080611c215750600a54155b611c3d5760405162461bcd60e51b81526004016105af906142ff565b848484602001516000808260001415611c93575060001984016000818152600660209081526040918290208251808401909352805483526001908101549183019190915290611c8b90612216565b039150611c9c565b50506000198101835b81846020015114611cbf5760405162461bcd60e51b81526004016105af90614078565b600081815260066020526040902054611cd89085612223565b611cf45760405162461bcd60e51b81526004016105af90613ef4565b60008a815260066020526040902054611d0d9089612185565b611d295760405162461bcd60e51b81526004016105af90613fe9565b600d548951516003548a516020015160405163ab5a164f60e01b815260009485946001600160a01b039091169363ab5a164f93610d93938f90600401613da9565b600f5481565b60096020526000908152604090205481565b69ffffffff000000000000603082901b16600160501b600160f01b03605084901b1660ff60f01b60f086901b166001600160f81b031960f888901b16171717949350505050565b80516020808301516040808501516060860151608087015160a088015160c089015160e08a01519551600099611e03999098979101613ce5565b604051602081830303815290604052805190602001209050919050565b6000602082511115611e445760405162461bcd60e51b81526004016105af906142c8565b611e4c612bee565b6000805160206145f883398151915281527f633dc4d7da7256660a892f8f1604a44b5432649cc8ec5cb3ced4c4e6ac94dd1d60208201527f890740a8eb06ce9be422cb8da5cdafc2b58c0a5e24036c578de2a433c828ff7d60408201527f3b8ec09e026fdc305365dfc94e189a81b38c7597b3d941c279f042e8206e0bd86060808301919091527fecd50eee38e386bd62be9bedb990706951b65fe053bd9d8a521af753d139e2da60808301527fdefff6d330bb5403f63b14f33b578274160de3a50df4efecf0e0db73bcdd3da560a08301528351600181811692918101901c90816001600160401b0381118015611f4357600080fd5b50604051908082528060200260200182016040528015611f6d578160200160208202803683370190505b5090508560005b60005b858503811015611ffc576000600182901b9050838181518110611f9657fe5b6020026020010151848260010181518110611fad57fe5b6020026020010151604051602001611fc6929190613d79565b60405160208183030381529060405280519060200120858381518110611fe857fe5b602090810291909101015250600101611f77565b856001141561206e5782600182901b8151811061201557fe5b602002602001015187836006811061202957fe5b602002015160405160200161203f929190613d79565b6040516020818303038152906040528051906020012084828151811061206157fe5b6020026020010181815250505b846001141561207d5750612094565b5060018085169550938401841c9383925001611f74565b826000815181106120a157fe5b60200260200101519650505050505050919050565b6000543410156120d85760405162461bcd60e51b81526004016105af90614389565b60405180604001604052808581526020016121048360048111156120f857fe5b86336001544301611d82565b905260078054600090815260066020908152604091829020845181559301516001909301929092555490517f3c367559be37298dc40fed468321ed4d44d99ab593f8e0bb7c82dbc84ab0bad69161215e91859085906144e1565b60405180910390a150506007805460010190555050565b6020015160301c63ffffffff1690565b6000612198836110bf84600001516129b1565b90505b92915050565b600a8190556121ae6127f5565b50565b602081015160501c6001600160a01b0316919050565b610100820460009081526020919091526040902054600160ff9092169190911b9081161490565b61010082046000908152602091909152604090208054600160ff9093169290921b9091179055565b6020015160f01c60ff1690565b6000612198836110bf84600001516129d6565b6000612198836110bf84600001516129f3565b6000612253612c0c565b6000805160206145f883398151915281527f633dc4d7da7256660a892f8f1604a44b5432649cc8ec5cb3ced4c4e6ac94dd1d6020808301919091527f890740a8eb06ce9be422cb8da5cdafc2b58c0a5e24036c578de2a433c828ff7d60408301527f3b8ec09e026fdc305365dfc94e189a81b38c7597b3d941c279f042e8206e0bd860608301527fecd50eee38e386bd62be9bedb990706951b65fe053bd9d8a521af753d139e2da60808301527fdefff6d330bb5403f63b14f33b578274160de3a50df4efecf0e0db73bcdd3da560a08301527f617bdd11f7c0a11f49db22f629387a12da7596f9d1704d7465177c63d88ec7d760c08301527f292c23a9aa1d8bea7e2435e555a4a60e379a5a35f3f452bae60121073fb6eead60e08301527fe1cea92ed99acdcb045a6726b2f87107e8a61620a232cf4d7d5b5766b3952e106101008301527f7ad66c0a68c72cb89e4fb4303841966e4062a76ab97451e3b9fb526a5ceb7f826101208301527fe026cc5a4aed3c22a58cbd3d2ac754c9352c5436f638042dca99034e836365166101408301527f3d04cffd8b46a874edf5cfae63077de85f849a660426697b06a829c70dd1409c6101608301527fad676aa337a485e4728a0b240d92b3ef7b3c372d06d189322bfd5f61f1e7203e6101808301527fa2fca4a49658f9fab7aa63289c91b7c7b6c832a6d0e69334ff5b0a3483d09dab6101a08301527f4ebfd9cd7bca2505f7bef59cc1c12ecc708fff26ae4af19abe852afe9e20c8626101c08301527f2def10d13dd169f550f578bda343d9717a138562e0093b380a1120789d53cf106101e08301527f776a31db34a1a0a7caaf862cffdfff1789297ffadc380bd3d39281d340abd3ad6102008301527fe2e7610b87a5fdf3a72ebe271287d923ab990eefac64b6e59d79f8b7e08c46e36102208301527f504364a5c6858bf98fff714ab5be9de19ed31a976860efbd0e772a2efe23e2e06102408301527f4f05f4acb83f5b65168d9fef89d56d4d77b8944015e6b1eed81b0238e2d0dba36102608301527f44a6d974c75b07423e1d6d33f481916fdd45830aea11b6347e700cd8b9f0767c6102808301527fedf260291f734ddac396a956127dde4c34c0cfb8d8052f88ac139658ccf2d5076102a08301527f6075c657a105351e7f0fce53bc320113324a522e8fd52dc878c762551e01a46e6102c08301527f6ca6a3f763a9395f7da16014725ca7ee17e4815c0ff8119bf33f273dee11833b6102e08301527f1c25ef10ffeb3c7d08aa707d17286e0b0d3cbcb50f1bd3b6523b63ba3b52dd0f6103008301527ffffc43bd08273ccf135fd3cacbeef055418e09eb728d727c4d5d5c556cdea7e36103208301527fc5ab8111456b1f28f3c7a0a604b4553ce905cb019c463ee159137af83c350b226103408301527f0ff273fcbf4ae0f2bd88d6cf319ff4004f8d7dca70d4ced4e74d2c74139739e66103608301527f7fa06ba11241ddd5efdc65d4e39c9f6991b74fd4b81b62230808216c876f827c6103808301527f7e275adf313a996c7e2950cac67caba02a5ff925ebf9906b58949f3e77aec5b96103a08301527f8f6162fa308d2b3a15dc33cffac85f13ab349173121645aedf00f471663108be6103c08301527f78ccaaab73373552f207a63599de54d7d8d0c1805f86ce7da15818d09f4cff626103e083015281908490811061272657fe5b60200201519392505050565b60008461274085858561274a565b1495945050505050565b600083815b83518110156127ec57600185821c166127a5578184828151811061276f57fe5b6020026020010151604051602001612788929190613d79565b6040516020818303038152906040528051906020012091506127e4565b8381815181106127b157fe5b6020026020010151826040516020016127cb929190613d79565b6040516020818303038152906040528051906020012091505b60010161274f565b50949350505050565b600a546007546000918291908103908290600019015b828210156128bd576002545a11612821576128bd565b818103600081815260066020908152604080832083815560010183905560089091529020549550935084156128b257600480546040516356f0001360e11b81526001600160a01b039091169163ade000269161287f91899101613da0565b600060405180830381600087803b15801561289957600080fd5b505af11580156128ad573d6000803e3d6000fd5b505050505b60019091019061280b565b60078054839003905582821480156128d5576000600a555b7f595cb200f38fcdbf73429ffce272e53b2e923c231b2059ce20d1f4909cff1bec8284836040516129089392919061451d565b60405180910390a16000805461291f908590612a07565b905060006129396003612933846002612a07565b90612a41565b905060006129478383612a83565b604051909150339083156108fc029084906000818181858888f19350505050158015612977573d6000803e3d6000fd5b5060405160009082156108fc0290839083818181858288f193505050501580156129a5573d6000803e3d6000fd5b50505050505050505050565b600081600001516129c58360200151611dc9565b604051602001611e03929190613d79565b600081600001518260200151604051602001611e03929190613d79565b600081600001516129c58360200151612ac5565b600082612a165750600061219b565b82820282848281612a2357fe5b04146121985760405162461bcd60e51b81526004016105af906141b0565b600061219883836040518060400160405280601a81526020017f536166654d6174683a206469766973696f6e206279207a65726f000000000000815250612aee565b600061219883836040518060400160405280601e81526020017f536166654d6174683a207375627472616374696f6e206f766572666c6f770000815250612b25565b60008160000151826020015183604001518460600151604051602001611e039493929190613d3e565b60008183612b0f5760405162461bcd60e51b81526004016105af9190613ea0565b506000838581612b1b57fe5b0495945050505050565b60008184841115612b495760405162461bcd60e51b81526004016105af9190613ea0565b505050900390565b6040805161010081019091526000815260208101612b6d612c2b565b81526020016000815260200160008019168152602001600081526020016000815260200160008152602001606081525090565b6040518060a00160405280612bb3612c2b565b81526000602082018190526040820181905260608083019190915260809091015290565b604080518082019091526000808252602082015290565b6040518060c001604052806006906020820280368337509192915050565b6040518061040001604052806020906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b600082601f830112612c59578081fd5b8135612c6c612c678261459f565b614579565b818152915060208083019084810160005b84811015612ca657612c94888484358a0101612e34565b84529282019290820190600101612c7d565b505050505092915050565b60008083601f840112612cc2578081fd5b5081356001600160401b03811115612cd8578182fd5b602083019150836020604083028501011115612cf357600080fd5b9250929050565b60008083601f840112612d0b578182fd5b5081356001600160401b03811115612d21578182fd5b602083019150836020608083028501011115612cf357600080fd5b6000601f8381840112612d4d578182fd5b8235612d5b612c678261459f565b8181529250602080840190858101608080850288018301891015612d7e57600080fd5b60005b85811015612de6578987840112612d9757600080fd5b612da082614579565b80848486018d811115612db257600080fd5b60005b6004811015612dd257823585529388019391880191600101612db5565b509188525095850195935050600101612d81565b5050505050505092915050565b60008083601f840112612e04578182fd5b5081356001600160401b03811115612e1a578182fd5b6020830191508360208083028501011115612cf357600080fd5b600082601f830112612e44578081fd5b8135612e52612c678261459f565b818152915060208083019084810181840286018201871015612e7357600080fd5b60005b84811015612ca657813584529282019290820190600101612e76565b600082601f830112612ea2578081fd5b8135612eb0612c678261459f565b818152915060208083019084810160005b84811015612ca6578135870160a080601f19838c03011215612ee257600080fd5b612eec6040614579565b612ef88b8785016133db565b815290820135906001600160401b03821115612f1357600080fd5b612f218b8784860101612e34565b81870152865250509282019290820190600101612ec1565b600082601f830112612f49578081fd5b8135612f57612c678261459f565b8181529150602080830190848101608080850287018301881015612f7a57600080fd5b60005b85811015612fa157612f8f89846133db565b85529383019391810191600101612f7d565b50505050505092915050565b600082601f830112612fbd578081fd5b612fc76040614579565b9050808284604085011115612fdb57600080fd5b60005b6002811015612ffd578135835260209283019290910190600101612fde565b50505092915050565b600082601f830112613016578081fd5b81356001600160401b0381111561302b578182fd5b61303e601f8201601f1916602001614579565b915080825283602082850101111561305557600080fd5b8060208401602084013760009082016020015292915050565b60008183036080811215613080578182fd5b61308a6060614579565b9150604081121561309a57600080fd5b506130a56040614579565b8235815260208301356020820152808252506040820135602082015260608201356001600160401b038111156130da57600080fd5b6130e684828501612e34565b60408301525092915050565b600060608284031215613103578081fd5b61310d6060614579565b905081356001600160401b038082111561312657600080fd5b908301906040828603121561313a57600080fd5b6131446040614579565b8235815260208301358281111561315a57600080fd5b6131668782860161319b565b6020830152508084525060208401356020840152604084013591508082111561318e57600080fd5b506130e684828501612e34565b600061012082840312156131ad578081fd5b6101006131b981614579565b9150823582526131cc8460208501612fad565b6020830152606083013560408301526080830135606083015260a0830135608083015260c083013560a083015260e083013560c08301528083013590506001600160401b0381111561321d57600080fd5b61322984828501613006565b60e08301525092915050565b600060808284031215613246578081fd5b6132506080614579565b905081356001600160401b038082111561326957600080fd5b61327585838601612f39565b8352602084013591508082111561328b57600080fd5b61329785838601612c49565b602084015260408401359150808211156132b057600080fd5b6132bc85838601612d3c565b604084015260608401359150808211156132d557600080fd5b506132e284828501612c49565b60608301525092915050565b6000606082840312156132ff578081fd5b6133096060614579565b905081356001600160401b038082111561332257600080fd5b908301906040828603121561333657600080fd5b6133406040614579565b823581526020808401358381111561335757600080fd5b939093019260a0848803121561336c57600080fd5b6133766080614579565b8435815261338688838701612fad565b82820152606085013560408201526080850135848111156133a657600080fd5b6133b289828801613006565b606083015250828201529084528481013590840152604084013591508082111561318e57600080fd5b6000608082840312156133ec578081fd5b6133f66080614579565b90508135815260208201356020820152604082013560408201526060820135606082015292915050565b600060208284031215613431578081fd5b81516001600160a01b0381168114612198578182fd5b60008060008060008060008060008060a08b8d031215613465578586fd5b8a356001600160401b038082111561347b578788fd5b6134878e838f01612df3565b909c509a5060208d013591508082111561349f578788fd5b6134ab8e838f01612cb1565b909a50985060408d01359150808211156134c3578788fd5b6134cf8e838f01612cfa565b909850965060608d01359150808211156134e7578586fd5b6134f38e838f01612df3565b909650945060808d013591508082111561350b578384fd5b506135188d828e01612df3565b915080935050809150509295989b9194979a5092959850565b6000806000806000806000806080898b03121561354c578182fd5b88356001600160401b0380821115613562578384fd5b61356e8c838d01612df3565b909a50985060208b0135915080821115613586578384fd5b6135928c838d01612cb1565b909850965060408b01359150808211156135aa578384fd5b6135b68c838d01612df3565b909650945060608b01359150808211156135ce578384fd5b506135db8b828c01612df3565b999c989b5096995094979396929594505050565b600060208284031215613600578081fd5b5051919050565b60008060408385031215613619578182fd5b82519150602083015161362b816145ea565b809150509250929050565b600060208284031215613647578081fd5b8151612198816145ea565b60008060408385031215613664578182fd5b82356001600160401b038082111561367a578384fd5b6136868683870161306e565b9350602085013591508082111561369b578283fd5b90840190606082870312156136ae578283fd5b6136b86060614579565b82358152602083013560208201526040830135828111156136d7578485fd5b6136e388828601612e34565b6040830152508093505050509250929050565b600060208284031215613707578081fd5b5035919050565b60008060408385031215613720578182fd5b505080516020909101519092909150565b60008060008060808587031215613746578182fd5b8435935060208501356001600160401b0380821115613763578384fd5b61376f8883890161306e565b94506040870135915080821115613784578384fd5b613790888389016130f2565b935060608701359150808211156137a5578283fd5b506137b287828801612e92565b91505092959194509250565b600080600080608085870312156137d3578182fd5b8435935060208501356001600160401b03808211156137f0578384fd5b6137fc8883890161306e565b94506040870135915080821115613811578384fd5b613790888389016132ee565b600080600060608486031215613831578081fd5b8335925060208401356001600160401b038082111561384e578283fd5b61385a878388016130f2565b9350604086013591508082111561386f578283fd5b5061387c86828701613235565b9150509250925092565b60008060006060848603121561389a578081fd5b8335925060208401356001600160401b03808211156138b7578283fd5b6138c3878388016132ee565b935060408601359150808211156138d8578283fd5b9085019060c082880312156138eb578283fd5b6138f560c0614579565b823582811115613903578485fd5b61390f89828601612f39565b825250602083013582811115613923578485fd5b61392f89828601612c49565b602083015250604083013582811115613946578485fd5b61395289828601612d3c565b604083015250606083013582811115613969578485fd5b61397589828601612c49565b60608301525060808301358281111561398c578485fd5b61399889828601612d3c565b60808301525060a0830135828111156139af578485fd5b6139bb89828601612c49565b60a0830152508093505050509250925092565b6000806000606084860312156139e2578081fd5b8335925060208401356001600160401b03808211156139ff578283fd5b61385a878388016132ee565b6000815180845260208085018081965082840281019150828601855b85811015613a51578284038952613a3f848351613abc565b98850198935090840190600101613a27565b5091979650505050505050565b6000815180845260208085019450808401835b83811015613ab157815187865b6004811015613a9b57825182529185019190850190600101613a7e565b5050506080969096019590820190600101613a71565b509495945050505050565b6000815180845260208085019450808401835b83811015613ab157815187529582019590820190600101613acf565b6000815180845260208085018081965082840281019150828601855b85811015613a51578284038952815160a0613b23868351613c94565b868201519150806080870152613b3b81870183613abc565b9a87019a9550505090840190600101613b07565b6000815180845260208085019450808401835b83811015613ab157613b75878351613c94565b6080969096019590820190600101613b62565b8060005b6002811015613bab578151845260209384019390910190600101613b8c565b50505050565b60008151808452613bc98160208601602086016145be565b601f01601f19169290920160200192915050565b6000613bea838351613b88565b602082015160408401526040820151606084015260608201516080840152608082015160c060a0850152613c2160c0850182613bb1565b949350505050565b6000815160808452613c3e6080850182613b4f565b905060208301518482036020860152613c578282613a0b565b91505060408301518482036040860152613c718282613a5e565b91505060608301518482036060860152613c8b8282613a0b565b95945050505050565b805182526020810151602083015260408101516040830152606081015160608301525050565b6000868252604086602084013760608201859052828460808401379101608001908152949350505050565b6000898252613cf7602083018a613b88565b8760608301528660808301528560a08301528460c08301528360e08301526101008351613d2a81838601602088016145be565b929092019091019998505050505050505050565b6000858252613d506020830186613b88565b8360608301528251613d698160808501602087016145be565b9190910160800195945050505050565b918252602082015260400190565b6001600160a01b03929092168252602082015260400190565b90815260200190565b600085825284602083015260806040830152835160808301526020840151613dd460a0840182613b88565b50604084015160e083015260608401516101008301526080840151610120818185015260a086015161014085015260c086015161016085015260e086015191508061018085015250613e2a6101a0840182613bb1565b90508281036060840152613e3e8185613aeb565b979650505050505050565b600086825285602083015284604083015260a06060830152613e6e60a0830185613bb1565b8281036080840152613e808185613aeb565b98975050505050505050565b6001600160a01b0391909116815260200190565b6000602082526121986020830184613bb1565b60208082526021908201527f42617463684d616e616765723a204973206e6f7420726f6c6c696e67206261636040820152606b60f81b606082015260800190565b60208082526032908201527f70726576696f757320636f6d6d69746d656e7420697320616273656e7420696e604082015271040e8d0ca40c6eae4e4cadce840c4c2e8c6d60731b606082015260800190565b60208082526033908201527f70726576696f757320636f6d6d69746d656e7420697320616273656e7420696e604082015272040e8d0ca40e0e4caecd2deeae640c4c2e8c6d606b1b606082015260800190565b60208082526030908201527f596f7520617265206e6f742074686520636f727265637420636f6d6d6974746560408201526f0e440ccdee440e8d0d2e640c4c2e8c6d60831b606082015260800190565b60208082526028908201527f54617267657420636f6d6d69746d656e7420697320616273656e7420696e20746040820152670d0ca40c4c2e8c6d60c31b606082015260800190565b60208082526027908201527f526f6c6c75703a20436f6d6d69746d656e74206e6f742070726573656e7420696040820152660dc40c4c2e8c6d60cb1b606082015260800190565b60208082526022908201527f70726576696f757320636f6d6d69746d656e74206861732077726f6e672070616040820152610e8d60f31b606082015260800190565b60208082526031908201527f54686973206261746368206973206e6f74207965742066696e616c697365642c60408201527020636865636b206261636b20736f6f6e2160781b606082015260800190565b6020808252601d908201527f42617463684d616e616765723a20497320726f6c6c696e67206261636b000000604082015260600190565b60208082526017908201527f426174636820616c72656164792066696e616c69736564000000000000000000604082015260600190565b6020808252601f908201527f436f6d6d69746d656e74206e6f742070726573656e7420696e20626174636800604082015260600190565b60208082526021908201527f536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f6040820152607760f81b606082015260800190565b60208082526023908201527f526f6c6c75703a2053746174652073756274726565206973206e6f7420766163604082015262185b9d60ea1b606082015260800190565b60208082526018908201527f526f6c6c75703a20496e76616c69642070726f706f7365720000000000000000604082015260600190565b6020808252603f908201527f42617463682069642067726561746572207468616e20746f74616c206e756d6260408201527f6572206f6620626174636865732c20696e76616c696420626174636820696400606082015260800190565b6020808252601b908201527f4d65726b6c65547265653a20546f6f206d616e79206c65617665730000000000604082015260600190565b60208082526033908201527f416c7265616479207375636365737366756c6c792064697370757465642e20526040820152726f6c6c206261636b20696e2070726f6365737360681b606082015260800190565b60208082526019908201527f526f6c6c75703a20416c72656164792077697468647261776e00000000000000604082015260600190565b6020808252601a908201527f526f6c6c75703a2077726f6e67207374616b6520616d6f756e74000000000000604082015260600190565b6000604082526143d36040830185613bdd565b8281036020840152835160c082526143ee60c0830182613b4f565b9050602085015182820360208401526144078282613a0b565b915050604085015182820360408401526144218282613a5e565b9150506060850151828203606084015261443b8282613a0b565b915050608085015182820360808401526144558282613a5e565b91505060a085015182820360a0840152613e3e8282613a0b565b6000604082526144826040830185613bdd565b8281036020840152613c8b8185613c29565b6000606082526144a76060830186613bdd565b82810360208401526144b98186613c29565b915050826040830152949350505050565b815181526020918201519181019190915260400190565b8381526020810183905260608101600583106144f957fe5b826040830152949350505050565b9283526020830191909152604082015260600190565b92835260208301919091521515604082015260600190565b6000808335601e1984360301811261454b578283fd5b8301803591506001600160401b03821115614564578283fd5b602001915036819003821315612cf357600080fd5b6040518181016001600160401b038111828210171561459757600080fd5b604052919050565b60006001600160401b038211156145b4578081fd5b5060209081020190565b60005b838110156145d95781810151838201526020016145c1565b83811115613bab5750506000910152565b600b81106121ae57600080fdfe290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563a164736f6c634300060c000a290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563"

// DeployRollup deploys a new Ethereum contract, binding an instance of Rollup to it.
func DeployRollup(auth *bind.TransactOpts, backend bind.ContractBackend, _chooser common.Address, _depositManager common.Address, _accountRegistry common.Address, _transfer common.Address, _massMigration common.Address, _create2Transfer common.Address, genesisStateRoot [32]byte, stakeAmount *big.Int, blocksToFinalise *big.Int, minGasLeft *big.Int, maxTxsPerCommit *big.Int) (common.Address, *types.Transaction, *Rollup, error) {
	parsed, err := abi.JSON(strings.NewReader(RollupABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(RollupBin), backend, _chooser, _depositManager, _accountRegistry, _transfer, _massMigration, _create2Transfer, genesisStateRoot, stakeAmount, blocksToFinalise, minGasLeft, maxTxsPerCommit)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Rollup{RollupCaller: RollupCaller{contract: contract}, RollupTransactor: RollupTransactor{contract: contract}, RollupFilterer: RollupFilterer{contract: contract}}, nil
}

// Rollup is an auto generated Go binding around an Ethereum contract.
type Rollup struct {
	RollupCaller     // Read-only binding to the contract
	RollupTransactor // Write-only binding to the contract
	RollupFilterer   // Log filterer for contract events
}

// RollupCaller is an auto generated read-only Go binding around an Ethereum contract.
type RollupCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RollupTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RollupFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RollupSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RollupSession struct {
	Contract     *Rollup           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RollupCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RollupCallerSession struct {
	Contract *RollupCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// RollupTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RollupTransactorSession struct {
	Contract     *RollupTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RollupRaw is an auto generated low-level Go binding around an Ethereum contract.
type RollupRaw struct {
	Contract *Rollup // Generic contract binding to access the raw methods on
}

// RollupCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RollupCallerRaw struct {
	Contract *RollupCaller // Generic read-only contract binding to access the raw methods on
}

// RollupTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RollupTransactorRaw struct {
	Contract *RollupTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRollup creates a new instance of Rollup, bound to a specific deployed contract.
func NewRollup(address common.Address, backend bind.ContractBackend) (*Rollup, error) {
	contract, err := bindRollup(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Rollup{RollupCaller: RollupCaller{contract: contract}, RollupTransactor: RollupTransactor{contract: contract}, RollupFilterer: RollupFilterer{contract: contract}}, nil
}

// NewRollupCaller creates a new read-only instance of Rollup, bound to a specific deployed contract.
func NewRollupCaller(address common.Address, caller bind.ContractCaller) (*RollupCaller, error) {
	contract, err := bindRollup(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RollupCaller{contract: contract}, nil
}

// NewRollupTransactor creates a new write-only instance of Rollup, bound to a specific deployed contract.
func NewRollupTransactor(address common.Address, transactor bind.ContractTransactor) (*RollupTransactor, error) {
	contract, err := bindRollup(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RollupTransactor{contract: contract}, nil
}

// NewRollupFilterer creates a new log filterer instance of Rollup, bound to a specific deployed contract.
func NewRollupFilterer(address common.Address, filterer bind.ContractFilterer) (*RollupFilterer, error) {
	contract, err := bindRollup(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RollupFilterer{contract: contract}, nil
}

// bindRollup binds a generic wrapper to an already deployed contract.
func bindRollup(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RollupABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rollup *RollupRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Rollup.Contract.RollupCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rollup *RollupRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rollup.Contract.RollupTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rollup *RollupRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rollup.Contract.RollupTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rollup *RollupCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Rollup.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rollup *RollupTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rollup.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rollup *RollupTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rollup.Contract.contract.Transact(opts, method, params...)
}

// ZEROBYTES32 is a free data retrieval call binding the contract method 0x069321b0.
//
// Solidity: function ZERO_BYTES32() view returns(bytes32)
func (_Rollup *RollupCaller) ZEROBYTES32(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "ZERO_BYTES32")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ZEROBYTES32 is a free data retrieval call binding the contract method 0x069321b0.
//
// Solidity: function ZERO_BYTES32() view returns(bytes32)
func (_Rollup *RollupSession) ZEROBYTES32() ([32]byte, error) {
	return _Rollup.Contract.ZEROBYTES32(&_Rollup.CallOpts)
}

// ZEROBYTES32 is a free data retrieval call binding the contract method 0x069321b0.
//
// Solidity: function ZERO_BYTES32() view returns(bytes32)
func (_Rollup *RollupCallerSession) ZEROBYTES32() ([32]byte, error) {
	return _Rollup.Contract.ZEROBYTES32(&_Rollup.CallOpts)
}

// AccountRegistry is a free data retrieval call binding the contract method 0xd089e11a.
//
// Solidity: function accountRegistry() view returns(address)
func (_Rollup *RollupCaller) AccountRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "accountRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AccountRegistry is a free data retrieval call binding the contract method 0xd089e11a.
//
// Solidity: function accountRegistry() view returns(address)
func (_Rollup *RollupSession) AccountRegistry() (common.Address, error) {
	return _Rollup.Contract.AccountRegistry(&_Rollup.CallOpts)
}

// AccountRegistry is a free data retrieval call binding the contract method 0xd089e11a.
//
// Solidity: function accountRegistry() view returns(address)
func (_Rollup *RollupCallerSession) AccountRegistry() (common.Address, error) {
	return _Rollup.Contract.AccountRegistry(&_Rollup.CallOpts)
}

// AppID is a free data retrieval call binding the contract method 0xf9a6cbae.
//
// Solidity: function appID() view returns(bytes32)
func (_Rollup *RollupCaller) AppID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "appID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// AppID is a free data retrieval call binding the contract method 0xf9a6cbae.
//
// Solidity: function appID() view returns(bytes32)
func (_Rollup *RollupSession) AppID() ([32]byte, error) {
	return _Rollup.Contract.AppID(&_Rollup.CallOpts)
}

// AppID is a free data retrieval call binding the contract method 0xf9a6cbae.
//
// Solidity: function appID() view returns(bytes32)
func (_Rollup *RollupCallerSession) AppID() ([32]byte, error) {
	return _Rollup.Contract.AppID(&_Rollup.CallOpts)
}

// Batches is a free data retrieval call binding the contract method 0xb32c4d8d.
//
// Solidity: function batches(uint256 ) view returns(bytes32 commitmentRoot, bytes32 meta)
func (_Rollup *RollupCaller) Batches(opts *bind.CallOpts, arg0 *big.Int) (struct {
	CommitmentRoot [32]byte
	Meta           [32]byte
}, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "batches", arg0)

	outstruct := new(struct {
		CommitmentRoot [32]byte
		Meta           [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.CommitmentRoot = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.Meta = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// Batches is a free data retrieval call binding the contract method 0xb32c4d8d.
//
// Solidity: function batches(uint256 ) view returns(bytes32 commitmentRoot, bytes32 meta)
func (_Rollup *RollupSession) Batches(arg0 *big.Int) (struct {
	CommitmentRoot [32]byte
	Meta           [32]byte
}, error) {
	return _Rollup.Contract.Batches(&_Rollup.CallOpts, arg0)
}

// Batches is a free data retrieval call binding the contract method 0xb32c4d8d.
//
// Solidity: function batches(uint256 ) view returns(bytes32 commitmentRoot, bytes32 meta)
func (_Rollup *RollupCallerSession) Batches(arg0 *big.Int) (struct {
	CommitmentRoot [32]byte
	Meta           [32]byte
}, error) {
	return _Rollup.Contract.Batches(&_Rollup.CallOpts, arg0)
}

// Chooser is a free data retrieval call binding the contract method 0x98d17621.
//
// Solidity: function chooser() view returns(address)
func (_Rollup *RollupCaller) Chooser(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "chooser")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Chooser is a free data retrieval call binding the contract method 0x98d17621.
//
// Solidity: function chooser() view returns(address)
func (_Rollup *RollupSession) Chooser() (common.Address, error) {
	return _Rollup.Contract.Chooser(&_Rollup.CallOpts)
}

// Chooser is a free data retrieval call binding the contract method 0x98d17621.
//
// Solidity: function chooser() view returns(address)
func (_Rollup *RollupCallerSession) Chooser() (common.Address, error) {
	return _Rollup.Contract.Chooser(&_Rollup.CallOpts)
}

// Create2Transfer is a free data retrieval call binding the contract method 0x7ae8c568.
//
// Solidity: function create2Transfer() view returns(address)
func (_Rollup *RollupCaller) Create2Transfer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "create2Transfer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Create2Transfer is a free data retrieval call binding the contract method 0x7ae8c568.
//
// Solidity: function create2Transfer() view returns(address)
func (_Rollup *RollupSession) Create2Transfer() (common.Address, error) {
	return _Rollup.Contract.Create2Transfer(&_Rollup.CallOpts)
}

// Create2Transfer is a free data retrieval call binding the contract method 0x7ae8c568.
//
// Solidity: function create2Transfer() view returns(address)
func (_Rollup *RollupCallerSession) Create2Transfer() (common.Address, error) {
	return _Rollup.Contract.Create2Transfer(&_Rollup.CallOpts)
}

// DepositManager is a free data retrieval call binding the contract method 0x6c7ac9d8.
//
// Solidity: function depositManager() view returns(address)
func (_Rollup *RollupCaller) DepositManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "depositManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DepositManager is a free data retrieval call binding the contract method 0x6c7ac9d8.
//
// Solidity: function depositManager() view returns(address)
func (_Rollup *RollupSession) DepositManager() (common.Address, error) {
	return _Rollup.Contract.DepositManager(&_Rollup.CallOpts)
}

// DepositManager is a free data retrieval call binding the contract method 0x6c7ac9d8.
//
// Solidity: function depositManager() view returns(address)
func (_Rollup *RollupCallerSession) DepositManager() (common.Address, error) {
	return _Rollup.Contract.DepositManager(&_Rollup.CallOpts)
}

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits(uint256 ) view returns(bytes32)
func (_Rollup *RollupCaller) Deposits(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "deposits", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits(uint256 ) view returns(bytes32)
func (_Rollup *RollupSession) Deposits(arg0 *big.Int) ([32]byte, error) {
	return _Rollup.Contract.Deposits(&_Rollup.CallOpts, arg0)
}

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits(uint256 ) view returns(bytes32)
func (_Rollup *RollupCallerSession) Deposits(arg0 *big.Int) ([32]byte, error) {
	return _Rollup.Contract.Deposits(&_Rollup.CallOpts, arg0)
}

// GetBatch is a free data retrieval call binding the contract method 0x5ac44282.
//
// Solidity: function getBatch(uint256 batchID) view returns((bytes32,bytes32) batch)
func (_Rollup *RollupCaller) GetBatch(opts *bind.CallOpts, batchID *big.Int) (TypesBatch, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "getBatch", batchID)

	if err != nil {
		return *new(TypesBatch), err
	}

	out0 := *abi.ConvertType(out[0], new(TypesBatch)).(*TypesBatch)

	return out0, err

}

// GetBatch is a free data retrieval call binding the contract method 0x5ac44282.
//
// Solidity: function getBatch(uint256 batchID) view returns((bytes32,bytes32) batch)
func (_Rollup *RollupSession) GetBatch(batchID *big.Int) (TypesBatch, error) {
	return _Rollup.Contract.GetBatch(&_Rollup.CallOpts, batchID)
}

// GetBatch is a free data retrieval call binding the contract method 0x5ac44282.
//
// Solidity: function getBatch(uint256 batchID) view returns((bytes32,bytes32) batch)
func (_Rollup *RollupCallerSession) GetBatch(batchID *big.Int) (TypesBatch, error) {
	return _Rollup.Contract.GetBatch(&_Rollup.CallOpts, batchID)
}

// InvalidBatchMarker is a free data retrieval call binding the contract method 0x5b097d37.
//
// Solidity: function invalidBatchMarker() view returns(uint256)
func (_Rollup *RollupCaller) InvalidBatchMarker(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "invalidBatchMarker")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InvalidBatchMarker is a free data retrieval call binding the contract method 0x5b097d37.
//
// Solidity: function invalidBatchMarker() view returns(uint256)
func (_Rollup *RollupSession) InvalidBatchMarker() (*big.Int, error) {
	return _Rollup.Contract.InvalidBatchMarker(&_Rollup.CallOpts)
}

// InvalidBatchMarker is a free data retrieval call binding the contract method 0x5b097d37.
//
// Solidity: function invalidBatchMarker() view returns(uint256)
func (_Rollup *RollupCallerSession) InvalidBatchMarker() (*big.Int, error) {
	return _Rollup.Contract.InvalidBatchMarker(&_Rollup.CallOpts)
}

// MassMigration is a free data retrieval call binding the contract method 0x0ed75b9c.
//
// Solidity: function massMigration() view returns(address)
func (_Rollup *RollupCaller) MassMigration(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "massMigration")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MassMigration is a free data retrieval call binding the contract method 0x0ed75b9c.
//
// Solidity: function massMigration() view returns(address)
func (_Rollup *RollupSession) MassMigration() (common.Address, error) {
	return _Rollup.Contract.MassMigration(&_Rollup.CallOpts)
}

// MassMigration is a free data retrieval call binding the contract method 0x0ed75b9c.
//
// Solidity: function massMigration() view returns(address)
func (_Rollup *RollupCallerSession) MassMigration() (common.Address, error) {
	return _Rollup.Contract.MassMigration(&_Rollup.CallOpts)
}

// NextBatchID is a free data retrieval call binding the contract method 0xa72598ce.
//
// Solidity: function nextBatchID() view returns(uint256)
func (_Rollup *RollupCaller) NextBatchID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "nextBatchID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextBatchID is a free data retrieval call binding the contract method 0xa72598ce.
//
// Solidity: function nextBatchID() view returns(uint256)
func (_Rollup *RollupSession) NextBatchID() (*big.Int, error) {
	return _Rollup.Contract.NextBatchID(&_Rollup.CallOpts)
}

// NextBatchID is a free data retrieval call binding the contract method 0xa72598ce.
//
// Solidity: function nextBatchID() view returns(uint256)
func (_Rollup *RollupCallerSession) NextBatchID() (*big.Int, error) {
	return _Rollup.Contract.NextBatchID(&_Rollup.CallOpts)
}

// ParamBlocksToFinalise is a free data retrieval call binding the contract method 0x31c2b7db.
//
// Solidity: function paramBlocksToFinalise() view returns(uint256)
func (_Rollup *RollupCaller) ParamBlocksToFinalise(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "paramBlocksToFinalise")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ParamBlocksToFinalise is a free data retrieval call binding the contract method 0x31c2b7db.
//
// Solidity: function paramBlocksToFinalise() view returns(uint256)
func (_Rollup *RollupSession) ParamBlocksToFinalise() (*big.Int, error) {
	return _Rollup.Contract.ParamBlocksToFinalise(&_Rollup.CallOpts)
}

// ParamBlocksToFinalise is a free data retrieval call binding the contract method 0x31c2b7db.
//
// Solidity: function paramBlocksToFinalise() view returns(uint256)
func (_Rollup *RollupCallerSession) ParamBlocksToFinalise() (*big.Int, error) {
	return _Rollup.Contract.ParamBlocksToFinalise(&_Rollup.CallOpts)
}

// ParamMaxTxsPerCommit is a free data retrieval call binding the contract method 0xe42f7554.
//
// Solidity: function paramMaxTxsPerCommit() view returns(uint256)
func (_Rollup *RollupCaller) ParamMaxTxsPerCommit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "paramMaxTxsPerCommit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ParamMaxTxsPerCommit is a free data retrieval call binding the contract method 0xe42f7554.
//
// Solidity: function paramMaxTxsPerCommit() view returns(uint256)
func (_Rollup *RollupSession) ParamMaxTxsPerCommit() (*big.Int, error) {
	return _Rollup.Contract.ParamMaxTxsPerCommit(&_Rollup.CallOpts)
}

// ParamMaxTxsPerCommit is a free data retrieval call binding the contract method 0xe42f7554.
//
// Solidity: function paramMaxTxsPerCommit() view returns(uint256)
func (_Rollup *RollupCallerSession) ParamMaxTxsPerCommit() (*big.Int, error) {
	return _Rollup.Contract.ParamMaxTxsPerCommit(&_Rollup.CallOpts)
}

// ParamMinGasLeft is a free data retrieval call binding the contract method 0xca858d9d.
//
// Solidity: function paramMinGasLeft() view returns(uint256)
func (_Rollup *RollupCaller) ParamMinGasLeft(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "paramMinGasLeft")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ParamMinGasLeft is a free data retrieval call binding the contract method 0xca858d9d.
//
// Solidity: function paramMinGasLeft() view returns(uint256)
func (_Rollup *RollupSession) ParamMinGasLeft() (*big.Int, error) {
	return _Rollup.Contract.ParamMinGasLeft(&_Rollup.CallOpts)
}

// ParamMinGasLeft is a free data retrieval call binding the contract method 0xca858d9d.
//
// Solidity: function paramMinGasLeft() view returns(uint256)
func (_Rollup *RollupCallerSession) ParamMinGasLeft() (*big.Int, error) {
	return _Rollup.Contract.ParamMinGasLeft(&_Rollup.CallOpts)
}

// ParamStakeAmount is a free data retrieval call binding the contract method 0xd1b243e2.
//
// Solidity: function paramStakeAmount() view returns(uint256)
func (_Rollup *RollupCaller) ParamStakeAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "paramStakeAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ParamStakeAmount is a free data retrieval call binding the contract method 0xd1b243e2.
//
// Solidity: function paramStakeAmount() view returns(uint256)
func (_Rollup *RollupSession) ParamStakeAmount() (*big.Int, error) {
	return _Rollup.Contract.ParamStakeAmount(&_Rollup.CallOpts)
}

// ParamStakeAmount is a free data retrieval call binding the contract method 0xd1b243e2.
//
// Solidity: function paramStakeAmount() view returns(uint256)
func (_Rollup *RollupCallerSession) ParamStakeAmount() (*big.Int, error) {
	return _Rollup.Contract.ParamStakeAmount(&_Rollup.CallOpts)
}

// Transfer is a free data retrieval call binding the contract method 0x8a4068dd.
//
// Solidity: function transfer() view returns(address)
func (_Rollup *RollupCaller) Transfer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "transfer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Transfer is a free data retrieval call binding the contract method 0x8a4068dd.
//
// Solidity: function transfer() view returns(address)
func (_Rollup *RollupSession) Transfer() (common.Address, error) {
	return _Rollup.Contract.Transfer(&_Rollup.CallOpts)
}

// Transfer is a free data retrieval call binding the contract method 0x8a4068dd.
//
// Solidity: function transfer() view returns(address)
func (_Rollup *RollupCallerSession) Transfer() (common.Address, error) {
	return _Rollup.Contract.Transfer(&_Rollup.CallOpts)
}

// WithdrawalBitmap is a free data retrieval call binding the contract method 0xfe5f9604.
//
// Solidity: function withdrawalBitmap(uint256 ) view returns(uint256)
func (_Rollup *RollupCaller) WithdrawalBitmap(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "withdrawalBitmap", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// WithdrawalBitmap is a free data retrieval call binding the contract method 0xfe5f9604.
//
// Solidity: function withdrawalBitmap(uint256 ) view returns(uint256)
func (_Rollup *RollupSession) WithdrawalBitmap(arg0 *big.Int) (*big.Int, error) {
	return _Rollup.Contract.WithdrawalBitmap(&_Rollup.CallOpts, arg0)
}

// WithdrawalBitmap is a free data retrieval call binding the contract method 0xfe5f9604.
//
// Solidity: function withdrawalBitmap(uint256 ) view returns(uint256)
func (_Rollup *RollupCallerSession) WithdrawalBitmap(arg0 *big.Int) (*big.Int, error) {
	return _Rollup.Contract.WithdrawalBitmap(&_Rollup.CallOpts, arg0)
}

// DisputeSignatureCreate2Transfer is a paid mutator transaction binding the contract method 0x50980f40.
//
// Solidity: function disputeSignatureCreate2Transfer(uint256 batchID, ((bytes32,(bytes32,uint256[2],uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256)[],bytes32[][],uint256[4][],bytes32[][],uint256[4][],bytes32[][]) signatureProof) returns()
func (_Rollup *RollupTransactor) DisputeSignatureCreate2Transfer(opts *bind.TransactOpts, batchID *big.Int, target TypesTransferCommitmentInclusionProof, signatureProof TypesSignatureProofWithReceiver) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "disputeSignatureCreate2Transfer", batchID, target, signatureProof)
}

// DisputeSignatureCreate2Transfer is a paid mutator transaction binding the contract method 0x50980f40.
//
// Solidity: function disputeSignatureCreate2Transfer(uint256 batchID, ((bytes32,(bytes32,uint256[2],uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256)[],bytes32[][],uint256[4][],bytes32[][],uint256[4][],bytes32[][]) signatureProof) returns()
func (_Rollup *RollupSession) DisputeSignatureCreate2Transfer(batchID *big.Int, target TypesTransferCommitmentInclusionProof, signatureProof TypesSignatureProofWithReceiver) (*types.Transaction, error) {
	return _Rollup.Contract.DisputeSignatureCreate2Transfer(&_Rollup.TransactOpts, batchID, target, signatureProof)
}

// DisputeSignatureCreate2Transfer is a paid mutator transaction binding the contract method 0x50980f40.
//
// Solidity: function disputeSignatureCreate2Transfer(uint256 batchID, ((bytes32,(bytes32,uint256[2],uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256)[],bytes32[][],uint256[4][],bytes32[][],uint256[4][],bytes32[][]) signatureProof) returns()
func (_Rollup *RollupTransactorSession) DisputeSignatureCreate2Transfer(batchID *big.Int, target TypesTransferCommitmentInclusionProof, signatureProof TypesSignatureProofWithReceiver) (*types.Transaction, error) {
	return _Rollup.Contract.DisputeSignatureCreate2Transfer(&_Rollup.TransactOpts, batchID, target, signatureProof)
}

// DisputeSignatureMassMigration is a paid mutator transaction binding the contract method 0x2538507d.
//
// Solidity: function disputeSignatureMassMigration(uint256 batchID, ((bytes32,(bytes32,uint256[2],uint256,bytes32,uint256,uint256,uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256)[],bytes32[][],uint256[4][],bytes32[][]) signatureProof) returns()
func (_Rollup *RollupTransactor) DisputeSignatureMassMigration(opts *bind.TransactOpts, batchID *big.Int, target TypesMMCommitmentInclusionProof, signatureProof TypesSignatureProof) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "disputeSignatureMassMigration", batchID, target, signatureProof)
}

// DisputeSignatureMassMigration is a paid mutator transaction binding the contract method 0x2538507d.
//
// Solidity: function disputeSignatureMassMigration(uint256 batchID, ((bytes32,(bytes32,uint256[2],uint256,bytes32,uint256,uint256,uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256)[],bytes32[][],uint256[4][],bytes32[][]) signatureProof) returns()
func (_Rollup *RollupSession) DisputeSignatureMassMigration(batchID *big.Int, target TypesMMCommitmentInclusionProof, signatureProof TypesSignatureProof) (*types.Transaction, error) {
	return _Rollup.Contract.DisputeSignatureMassMigration(&_Rollup.TransactOpts, batchID, target, signatureProof)
}

// DisputeSignatureMassMigration is a paid mutator transaction binding the contract method 0x2538507d.
//
// Solidity: function disputeSignatureMassMigration(uint256 batchID, ((bytes32,(bytes32,uint256[2],uint256,bytes32,uint256,uint256,uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256)[],bytes32[][],uint256[4][],bytes32[][]) signatureProof) returns()
func (_Rollup *RollupTransactorSession) DisputeSignatureMassMigration(batchID *big.Int, target TypesMMCommitmentInclusionProof, signatureProof TypesSignatureProof) (*types.Transaction, error) {
	return _Rollup.Contract.DisputeSignatureMassMigration(&_Rollup.TransactOpts, batchID, target, signatureProof)
}

// DisputeSignatureTransfer is a paid mutator transaction binding the contract method 0xc08caaef.
//
// Solidity: function disputeSignatureTransfer(uint256 batchID, ((bytes32,(bytes32,uint256[2],uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256)[],bytes32[][],uint256[4][],bytes32[][]) signatureProof) returns()
func (_Rollup *RollupTransactor) DisputeSignatureTransfer(opts *bind.TransactOpts, batchID *big.Int, target TypesTransferCommitmentInclusionProof, signatureProof TypesSignatureProof) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "disputeSignatureTransfer", batchID, target, signatureProof)
}

// DisputeSignatureTransfer is a paid mutator transaction binding the contract method 0xc08caaef.
//
// Solidity: function disputeSignatureTransfer(uint256 batchID, ((bytes32,(bytes32,uint256[2],uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256)[],bytes32[][],uint256[4][],bytes32[][]) signatureProof) returns()
func (_Rollup *RollupSession) DisputeSignatureTransfer(batchID *big.Int, target TypesTransferCommitmentInclusionProof, signatureProof TypesSignatureProof) (*types.Transaction, error) {
	return _Rollup.Contract.DisputeSignatureTransfer(&_Rollup.TransactOpts, batchID, target, signatureProof)
}

// DisputeSignatureTransfer is a paid mutator transaction binding the contract method 0xc08caaef.
//
// Solidity: function disputeSignatureTransfer(uint256 batchID, ((bytes32,(bytes32,uint256[2],uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256)[],bytes32[][],uint256[4][],bytes32[][]) signatureProof) returns()
func (_Rollup *RollupTransactorSession) DisputeSignatureTransfer(batchID *big.Int, target TypesTransferCommitmentInclusionProof, signatureProof TypesSignatureProof) (*types.Transaction, error) {
	return _Rollup.Contract.DisputeSignatureTransfer(&_Rollup.TransactOpts, batchID, target, signatureProof)
}

// DisputeTransitionCreate2Transfer is a paid mutator transaction binding the contract method 0x5f6e91d5.
//
// Solidity: function disputeTransitionCreate2Transfer(uint256 batchID, ((bytes32,bytes32),uint256,bytes32[]) previous, ((bytes32,(bytes32,uint256[2],uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256),bytes32[])[] proofs) returns()
func (_Rollup *RollupTransactor) DisputeTransitionCreate2Transfer(opts *bind.TransactOpts, batchID *big.Int, previous TypesCommitmentInclusionProof, target TypesTransferCommitmentInclusionProof, proofs []TypesStateMerkleProof) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "disputeTransitionCreate2Transfer", batchID, previous, target, proofs)
}

// DisputeTransitionCreate2Transfer is a paid mutator transaction binding the contract method 0x5f6e91d5.
//
// Solidity: function disputeTransitionCreate2Transfer(uint256 batchID, ((bytes32,bytes32),uint256,bytes32[]) previous, ((bytes32,(bytes32,uint256[2],uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256),bytes32[])[] proofs) returns()
func (_Rollup *RollupSession) DisputeTransitionCreate2Transfer(batchID *big.Int, previous TypesCommitmentInclusionProof, target TypesTransferCommitmentInclusionProof, proofs []TypesStateMerkleProof) (*types.Transaction, error) {
	return _Rollup.Contract.DisputeTransitionCreate2Transfer(&_Rollup.TransactOpts, batchID, previous, target, proofs)
}

// DisputeTransitionCreate2Transfer is a paid mutator transaction binding the contract method 0x5f6e91d5.
//
// Solidity: function disputeTransitionCreate2Transfer(uint256 batchID, ((bytes32,bytes32),uint256,bytes32[]) previous, ((bytes32,(bytes32,uint256[2],uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256),bytes32[])[] proofs) returns()
func (_Rollup *RollupTransactorSession) DisputeTransitionCreate2Transfer(batchID *big.Int, previous TypesCommitmentInclusionProof, target TypesTransferCommitmentInclusionProof, proofs []TypesStateMerkleProof) (*types.Transaction, error) {
	return _Rollup.Contract.DisputeTransitionCreate2Transfer(&_Rollup.TransactOpts, batchID, previous, target, proofs)
}

// DisputeTransitionMassMigration is a paid mutator transaction binding the contract method 0xeb84e076.
//
// Solidity: function disputeTransitionMassMigration(uint256 batchID, ((bytes32,bytes32),uint256,bytes32[]) previous, ((bytes32,(bytes32,uint256[2],uint256,bytes32,uint256,uint256,uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256),bytes32[])[] proofs) returns()
func (_Rollup *RollupTransactor) DisputeTransitionMassMigration(opts *bind.TransactOpts, batchID *big.Int, previous TypesCommitmentInclusionProof, target TypesMMCommitmentInclusionProof, proofs []TypesStateMerkleProof) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "disputeTransitionMassMigration", batchID, previous, target, proofs)
}

// DisputeTransitionMassMigration is a paid mutator transaction binding the contract method 0xeb84e076.
//
// Solidity: function disputeTransitionMassMigration(uint256 batchID, ((bytes32,bytes32),uint256,bytes32[]) previous, ((bytes32,(bytes32,uint256[2],uint256,bytes32,uint256,uint256,uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256),bytes32[])[] proofs) returns()
func (_Rollup *RollupSession) DisputeTransitionMassMigration(batchID *big.Int, previous TypesCommitmentInclusionProof, target TypesMMCommitmentInclusionProof, proofs []TypesStateMerkleProof) (*types.Transaction, error) {
	return _Rollup.Contract.DisputeTransitionMassMigration(&_Rollup.TransactOpts, batchID, previous, target, proofs)
}

// DisputeTransitionMassMigration is a paid mutator transaction binding the contract method 0xeb84e076.
//
// Solidity: function disputeTransitionMassMigration(uint256 batchID, ((bytes32,bytes32),uint256,bytes32[]) previous, ((bytes32,(bytes32,uint256[2],uint256,bytes32,uint256,uint256,uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256),bytes32[])[] proofs) returns()
func (_Rollup *RollupTransactorSession) DisputeTransitionMassMigration(batchID *big.Int, previous TypesCommitmentInclusionProof, target TypesMMCommitmentInclusionProof, proofs []TypesStateMerkleProof) (*types.Transaction, error) {
	return _Rollup.Contract.DisputeTransitionMassMigration(&_Rollup.TransactOpts, batchID, previous, target, proofs)
}

// DisputeTransitionTransfer is a paid mutator transaction binding the contract method 0x4e23e8c3.
//
// Solidity: function disputeTransitionTransfer(uint256 batchID, ((bytes32,bytes32),uint256,bytes32[]) previous, ((bytes32,(bytes32,uint256[2],uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256),bytes32[])[] proofs) returns()
func (_Rollup *RollupTransactor) DisputeTransitionTransfer(opts *bind.TransactOpts, batchID *big.Int, previous TypesCommitmentInclusionProof, target TypesTransferCommitmentInclusionProof, proofs []TypesStateMerkleProof) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "disputeTransitionTransfer", batchID, previous, target, proofs)
}

// DisputeTransitionTransfer is a paid mutator transaction binding the contract method 0x4e23e8c3.
//
// Solidity: function disputeTransitionTransfer(uint256 batchID, ((bytes32,bytes32),uint256,bytes32[]) previous, ((bytes32,(bytes32,uint256[2],uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256),bytes32[])[] proofs) returns()
func (_Rollup *RollupSession) DisputeTransitionTransfer(batchID *big.Int, previous TypesCommitmentInclusionProof, target TypesTransferCommitmentInclusionProof, proofs []TypesStateMerkleProof) (*types.Transaction, error) {
	return _Rollup.Contract.DisputeTransitionTransfer(&_Rollup.TransactOpts, batchID, previous, target, proofs)
}

// DisputeTransitionTransfer is a paid mutator transaction binding the contract method 0x4e23e8c3.
//
// Solidity: function disputeTransitionTransfer(uint256 batchID, ((bytes32,bytes32),uint256,bytes32[]) previous, ((bytes32,(bytes32,uint256[2],uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256),bytes32[])[] proofs) returns()
func (_Rollup *RollupTransactorSession) DisputeTransitionTransfer(batchID *big.Int, previous TypesCommitmentInclusionProof, target TypesTransferCommitmentInclusionProof, proofs []TypesStateMerkleProof) (*types.Transaction, error) {
	return _Rollup.Contract.DisputeTransitionTransfer(&_Rollup.TransactOpts, batchID, previous, target, proofs)
}

// KeepRollingBack is a paid mutator transaction binding the contract method 0xac96f0cd.
//
// Solidity: function keepRollingBack() returns()
func (_Rollup *RollupTransactor) KeepRollingBack(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "keepRollingBack")
}

// KeepRollingBack is a paid mutator transaction binding the contract method 0xac96f0cd.
//
// Solidity: function keepRollingBack() returns()
func (_Rollup *RollupSession) KeepRollingBack() (*types.Transaction, error) {
	return _Rollup.Contract.KeepRollingBack(&_Rollup.TransactOpts)
}

// KeepRollingBack is a paid mutator transaction binding the contract method 0xac96f0cd.
//
// Solidity: function keepRollingBack() returns()
func (_Rollup *RollupTransactorSession) KeepRollingBack() (*types.Transaction, error) {
	return _Rollup.Contract.KeepRollingBack(&_Rollup.TransactOpts)
}

// SubmitCreate2Transfer is a paid mutator transaction binding the contract method 0x9251597f.
//
// Solidity: function submitCreate2Transfer(bytes32[] stateRoots, uint256[2][] signatures, uint256[] feeReceivers, bytes[] txss) payable returns()
func (_Rollup *RollupTransactor) SubmitCreate2Transfer(opts *bind.TransactOpts, stateRoots [][32]byte, signatures [][2]*big.Int, feeReceivers []*big.Int, txss [][]byte) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "submitCreate2Transfer", stateRoots, signatures, feeReceivers, txss)
}

// SubmitCreate2Transfer is a paid mutator transaction binding the contract method 0x9251597f.
//
// Solidity: function submitCreate2Transfer(bytes32[] stateRoots, uint256[2][] signatures, uint256[] feeReceivers, bytes[] txss) payable returns()
func (_Rollup *RollupSession) SubmitCreate2Transfer(stateRoots [][32]byte, signatures [][2]*big.Int, feeReceivers []*big.Int, txss [][]byte) (*types.Transaction, error) {
	return _Rollup.Contract.SubmitCreate2Transfer(&_Rollup.TransactOpts, stateRoots, signatures, feeReceivers, txss)
}

// SubmitCreate2Transfer is a paid mutator transaction binding the contract method 0x9251597f.
//
// Solidity: function submitCreate2Transfer(bytes32[] stateRoots, uint256[2][] signatures, uint256[] feeReceivers, bytes[] txss) payable returns()
func (_Rollup *RollupTransactorSession) SubmitCreate2Transfer(stateRoots [][32]byte, signatures [][2]*big.Int, feeReceivers []*big.Int, txss [][]byte) (*types.Transaction, error) {
	return _Rollup.Contract.SubmitCreate2Transfer(&_Rollup.TransactOpts, stateRoots, signatures, feeReceivers, txss)
}

// SubmitDeposits is a paid mutator transaction binding the contract method 0x562a2bca.
//
// Solidity: function submitDeposits(((bytes32,bytes32),uint256,bytes32[]) previous, (uint256,uint256,bytes32[]) vacant) payable returns()
func (_Rollup *RollupTransactor) SubmitDeposits(opts *bind.TransactOpts, previous TypesCommitmentInclusionProof, vacant TypesSubtreeVacancyProof) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "submitDeposits", previous, vacant)
}

// SubmitDeposits is a paid mutator transaction binding the contract method 0x562a2bca.
//
// Solidity: function submitDeposits(((bytes32,bytes32),uint256,bytes32[]) previous, (uint256,uint256,bytes32[]) vacant) payable returns()
func (_Rollup *RollupSession) SubmitDeposits(previous TypesCommitmentInclusionProof, vacant TypesSubtreeVacancyProof) (*types.Transaction, error) {
	return _Rollup.Contract.SubmitDeposits(&_Rollup.TransactOpts, previous, vacant)
}

// SubmitDeposits is a paid mutator transaction binding the contract method 0x562a2bca.
//
// Solidity: function submitDeposits(((bytes32,bytes32),uint256,bytes32[]) previous, (uint256,uint256,bytes32[]) vacant) payable returns()
func (_Rollup *RollupTransactorSession) SubmitDeposits(previous TypesCommitmentInclusionProof, vacant TypesSubtreeVacancyProof) (*types.Transaction, error) {
	return _Rollup.Contract.SubmitDeposits(&_Rollup.TransactOpts, previous, vacant)
}

// SubmitMassMigration is a paid mutator transaction binding the contract method 0x035695e3.
//
// Solidity: function submitMassMigration(bytes32[] stateRoots, uint256[2][] signatures, uint256[4][] meta, bytes32[] withdrawRoots, bytes[] txss) payable returns()
func (_Rollup *RollupTransactor) SubmitMassMigration(opts *bind.TransactOpts, stateRoots [][32]byte, signatures [][2]*big.Int, meta [][4]*big.Int, withdrawRoots [][32]byte, txss [][]byte) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "submitMassMigration", stateRoots, signatures, meta, withdrawRoots, txss)
}

// SubmitMassMigration is a paid mutator transaction binding the contract method 0x035695e3.
//
// Solidity: function submitMassMigration(bytes32[] stateRoots, uint256[2][] signatures, uint256[4][] meta, bytes32[] withdrawRoots, bytes[] txss) payable returns()
func (_Rollup *RollupSession) SubmitMassMigration(stateRoots [][32]byte, signatures [][2]*big.Int, meta [][4]*big.Int, withdrawRoots [][32]byte, txss [][]byte) (*types.Transaction, error) {
	return _Rollup.Contract.SubmitMassMigration(&_Rollup.TransactOpts, stateRoots, signatures, meta, withdrawRoots, txss)
}

// SubmitMassMigration is a paid mutator transaction binding the contract method 0x035695e3.
//
// Solidity: function submitMassMigration(bytes32[] stateRoots, uint256[2][] signatures, uint256[4][] meta, bytes32[] withdrawRoots, bytes[] txss) payable returns()
func (_Rollup *RollupTransactorSession) SubmitMassMigration(stateRoots [][32]byte, signatures [][2]*big.Int, meta [][4]*big.Int, withdrawRoots [][32]byte, txss [][]byte) (*types.Transaction, error) {
	return _Rollup.Contract.SubmitMassMigration(&_Rollup.TransactOpts, stateRoots, signatures, meta, withdrawRoots, txss)
}

// SubmitTransfer is a paid mutator transaction binding the contract method 0x5f5b95b8.
//
// Solidity: function submitTransfer(bytes32[] stateRoots, uint256[2][] signatures, uint256[] feeReceivers, bytes[] txss) payable returns()
func (_Rollup *RollupTransactor) SubmitTransfer(opts *bind.TransactOpts, stateRoots [][32]byte, signatures [][2]*big.Int, feeReceivers []*big.Int, txss [][]byte) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "submitTransfer", stateRoots, signatures, feeReceivers, txss)
}

// SubmitTransfer is a paid mutator transaction binding the contract method 0x5f5b95b8.
//
// Solidity: function submitTransfer(bytes32[] stateRoots, uint256[2][] signatures, uint256[] feeReceivers, bytes[] txss) payable returns()
func (_Rollup *RollupSession) SubmitTransfer(stateRoots [][32]byte, signatures [][2]*big.Int, feeReceivers []*big.Int, txss [][]byte) (*types.Transaction, error) {
	return _Rollup.Contract.SubmitTransfer(&_Rollup.TransactOpts, stateRoots, signatures, feeReceivers, txss)
}

// SubmitTransfer is a paid mutator transaction binding the contract method 0x5f5b95b8.
//
// Solidity: function submitTransfer(bytes32[] stateRoots, uint256[2][] signatures, uint256[] feeReceivers, bytes[] txss) payable returns()
func (_Rollup *RollupTransactorSession) SubmitTransfer(stateRoots [][32]byte, signatures [][2]*big.Int, feeReceivers []*big.Int, txss [][]byte) (*types.Transaction, error) {
	return _Rollup.Contract.SubmitTransfer(&_Rollup.TransactOpts, stateRoots, signatures, feeReceivers, txss)
}

// WithdrawStake is a paid mutator transaction binding the contract method 0x25d5971f.
//
// Solidity: function withdrawStake(uint256 batchID) returns()
func (_Rollup *RollupTransactor) WithdrawStake(opts *bind.TransactOpts, batchID *big.Int) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "withdrawStake", batchID)
}

// WithdrawStake is a paid mutator transaction binding the contract method 0x25d5971f.
//
// Solidity: function withdrawStake(uint256 batchID) returns()
func (_Rollup *RollupSession) WithdrawStake(batchID *big.Int) (*types.Transaction, error) {
	return _Rollup.Contract.WithdrawStake(&_Rollup.TransactOpts, batchID)
}

// WithdrawStake is a paid mutator transaction binding the contract method 0x25d5971f.
//
// Solidity: function withdrawStake(uint256 batchID) returns()
func (_Rollup *RollupTransactorSession) WithdrawStake(batchID *big.Int) (*types.Transaction, error) {
	return _Rollup.Contract.WithdrawStake(&_Rollup.TransactOpts, batchID)
}

// RollupDepositsFinalisedIterator is returned from FilterDepositsFinalised and is used to iterate over the raw logs and unpacked data for DepositsFinalised events raised by the Rollup contract.
type RollupDepositsFinalisedIterator struct {
	Event *RollupDepositsFinalised // Event containing the contract specifics and raw log

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
func (it *RollupDepositsFinalisedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupDepositsFinalised)
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
		it.Event = new(RollupDepositsFinalised)
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
func (it *RollupDepositsFinalisedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupDepositsFinalisedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupDepositsFinalised represents a DepositsFinalised event raised by the Rollup contract.
type RollupDepositsFinalised struct {
	SubtreeID          *big.Int
	DepositSubTreeRoot [32]byte
	PathToSubTree      *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterDepositsFinalised is a free log retrieval operation binding the contract event 0x1e6cd0ae9aa6f6e6ea5cf7f07396a96fc1ea5ae4ec7d5054f53e204c47f79bd8.
//
// Solidity: event DepositsFinalised(uint256 subtreeID, bytes32 depositSubTreeRoot, uint256 pathToSubTree)
func (_Rollup *RollupFilterer) FilterDepositsFinalised(opts *bind.FilterOpts) (*RollupDepositsFinalisedIterator, error) {

	logs, sub, err := _Rollup.contract.FilterLogs(opts, "DepositsFinalised")
	if err != nil {
		return nil, err
	}
	return &RollupDepositsFinalisedIterator{contract: _Rollup.contract, event: "DepositsFinalised", logs: logs, sub: sub}, nil
}

// WatchDepositsFinalised is a free log subscription operation binding the contract event 0x1e6cd0ae9aa6f6e6ea5cf7f07396a96fc1ea5ae4ec7d5054f53e204c47f79bd8.
//
// Solidity: event DepositsFinalised(uint256 subtreeID, bytes32 depositSubTreeRoot, uint256 pathToSubTree)
func (_Rollup *RollupFilterer) WatchDepositsFinalised(opts *bind.WatchOpts, sink chan<- *RollupDepositsFinalised) (event.Subscription, error) {

	logs, sub, err := _Rollup.contract.WatchLogs(opts, "DepositsFinalised")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupDepositsFinalised)
				if err := _Rollup.contract.UnpackLog(event, "DepositsFinalised", log); err != nil {
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

// ParseDepositsFinalised is a log parse operation binding the contract event 0x1e6cd0ae9aa6f6e6ea5cf7f07396a96fc1ea5ae4ec7d5054f53e204c47f79bd8.
//
// Solidity: event DepositsFinalised(uint256 subtreeID, bytes32 depositSubTreeRoot, uint256 pathToSubTree)
func (_Rollup *RollupFilterer) ParseDepositsFinalised(log types.Log) (*RollupDepositsFinalised, error) {
	event := new(RollupDepositsFinalised)
	if err := _Rollup.contract.UnpackLog(event, "DepositsFinalised", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollupNewBatchIterator is returned from FilterNewBatch and is used to iterate over the raw logs and unpacked data for NewBatch events raised by the Rollup contract.
type RollupNewBatchIterator struct {
	Event *RollupNewBatch // Event containing the contract specifics and raw log

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
func (it *RollupNewBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupNewBatch)
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
		it.Event = new(RollupNewBatch)
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
func (it *RollupNewBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupNewBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupNewBatch represents a NewBatch event raised by the Rollup contract.
type RollupNewBatch struct {
	BatchID     *big.Int
	AccountRoot [32]byte
	BatchType   uint8
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterNewBatch is a free log retrieval operation binding the contract event 0x3c367559be37298dc40fed468321ed4d44d99ab593f8e0bb7c82dbc84ab0bad6.
//
// Solidity: event NewBatch(uint256 batchID, bytes32 accountRoot, uint8 batchType)
func (_Rollup *RollupFilterer) FilterNewBatch(opts *bind.FilterOpts) (*RollupNewBatchIterator, error) {

	logs, sub, err := _Rollup.contract.FilterLogs(opts, "NewBatch")
	if err != nil {
		return nil, err
	}
	return &RollupNewBatchIterator{contract: _Rollup.contract, event: "NewBatch", logs: logs, sub: sub}, nil
}

// WatchNewBatch is a free log subscription operation binding the contract event 0x3c367559be37298dc40fed468321ed4d44d99ab593f8e0bb7c82dbc84ab0bad6.
//
// Solidity: event NewBatch(uint256 batchID, bytes32 accountRoot, uint8 batchType)
func (_Rollup *RollupFilterer) WatchNewBatch(opts *bind.WatchOpts, sink chan<- *RollupNewBatch) (event.Subscription, error) {

	logs, sub, err := _Rollup.contract.WatchLogs(opts, "NewBatch")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupNewBatch)
				if err := _Rollup.contract.UnpackLog(event, "NewBatch", log); err != nil {
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

// ParseNewBatch is a log parse operation binding the contract event 0x3c367559be37298dc40fed468321ed4d44d99ab593f8e0bb7c82dbc84ab0bad6.
//
// Solidity: event NewBatch(uint256 batchID, bytes32 accountRoot, uint8 batchType)
func (_Rollup *RollupFilterer) ParseNewBatch(log types.Log) (*RollupNewBatch, error) {
	event := new(RollupNewBatch)
	if err := _Rollup.contract.UnpackLog(event, "NewBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollupRollbackStatusIterator is returned from FilterRollbackStatus and is used to iterate over the raw logs and unpacked data for RollbackStatus events raised by the Rollup contract.
type RollupRollbackStatusIterator struct {
	Event *RollupRollbackStatus // Event containing the contract specifics and raw log

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
func (it *RollupRollbackStatusIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupRollbackStatus)
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
		it.Event = new(RollupRollbackStatus)
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
func (it *RollupRollbackStatusIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupRollbackStatusIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupRollbackStatus represents a RollbackStatus event raised by the Rollup contract.
type RollupRollbackStatus struct {
	StartID   *big.Int
	NDeleted  *big.Int
	Completed bool
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRollbackStatus is a free log retrieval operation binding the contract event 0x595cb200f38fcdbf73429ffce272e53b2e923c231b2059ce20d1f4909cff1bec.
//
// Solidity: event RollbackStatus(uint256 startID, uint256 nDeleted, bool completed)
func (_Rollup *RollupFilterer) FilterRollbackStatus(opts *bind.FilterOpts) (*RollupRollbackStatusIterator, error) {

	logs, sub, err := _Rollup.contract.FilterLogs(opts, "RollbackStatus")
	if err != nil {
		return nil, err
	}
	return &RollupRollbackStatusIterator{contract: _Rollup.contract, event: "RollbackStatus", logs: logs, sub: sub}, nil
}

// WatchRollbackStatus is a free log subscription operation binding the contract event 0x595cb200f38fcdbf73429ffce272e53b2e923c231b2059ce20d1f4909cff1bec.
//
// Solidity: event RollbackStatus(uint256 startID, uint256 nDeleted, bool completed)
func (_Rollup *RollupFilterer) WatchRollbackStatus(opts *bind.WatchOpts, sink chan<- *RollupRollbackStatus) (event.Subscription, error) {

	logs, sub, err := _Rollup.contract.WatchLogs(opts, "RollbackStatus")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupRollbackStatus)
				if err := _Rollup.contract.UnpackLog(event, "RollbackStatus", log); err != nil {
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

// ParseRollbackStatus is a log parse operation binding the contract event 0x595cb200f38fcdbf73429ffce272e53b2e923c231b2059ce20d1f4909cff1bec.
//
// Solidity: event RollbackStatus(uint256 startID, uint256 nDeleted, bool completed)
func (_Rollup *RollupFilterer) ParseRollbackStatus(log types.Log) (*RollupRollbackStatus, error) {
	event := new(RollupRollbackStatus)
	if err := _Rollup.contract.UnpackLog(event, "RollbackStatus", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// RollupStakeWithdrawIterator is returned from FilterStakeWithdraw and is used to iterate over the raw logs and unpacked data for StakeWithdraw events raised by the Rollup contract.
type RollupStakeWithdrawIterator struct {
	Event *RollupStakeWithdraw // Event containing the contract specifics and raw log

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
func (it *RollupStakeWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupStakeWithdraw)
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
		it.Event = new(RollupStakeWithdraw)
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
func (it *RollupStakeWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupStakeWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupStakeWithdraw represents a StakeWithdraw event raised by the Rollup contract.
type RollupStakeWithdraw struct {
	Committed common.Address
	BatchID   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterStakeWithdraw is a free log retrieval operation binding the contract event 0x1248d48e2de900a1010c7fce73506969ecec243600bfc08b641b158f26d857cd.
//
// Solidity: event StakeWithdraw(address committed, uint256 batchID)
func (_Rollup *RollupFilterer) FilterStakeWithdraw(opts *bind.FilterOpts) (*RollupStakeWithdrawIterator, error) {

	logs, sub, err := _Rollup.contract.FilterLogs(opts, "StakeWithdraw")
	if err != nil {
		return nil, err
	}
	return &RollupStakeWithdrawIterator{contract: _Rollup.contract, event: "StakeWithdraw", logs: logs, sub: sub}, nil
}

// WatchStakeWithdraw is a free log subscription operation binding the contract event 0x1248d48e2de900a1010c7fce73506969ecec243600bfc08b641b158f26d857cd.
//
// Solidity: event StakeWithdraw(address committed, uint256 batchID)
func (_Rollup *RollupFilterer) WatchStakeWithdraw(opts *bind.WatchOpts, sink chan<- *RollupStakeWithdraw) (event.Subscription, error) {

	logs, sub, err := _Rollup.contract.WatchLogs(opts, "StakeWithdraw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupStakeWithdraw)
				if err := _Rollup.contract.UnpackLog(event, "StakeWithdraw", log); err != nil {
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

// ParseStakeWithdraw is a log parse operation binding the contract event 0x1248d48e2de900a1010c7fce73506969ecec243600bfc08b641b158f26d857cd.
//
// Solidity: event StakeWithdraw(address committed, uint256 batchID)
func (_Rollup *RollupFilterer) ParseStakeWithdraw(log types.Log) (*RollupStakeWithdraw, error) {
	event := new(RollupStakeWithdraw)
	if err := _Rollup.contract.UnpackLog(event, "StakeWithdraw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
