// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rollup

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
	PubkeyHashesReceiver    [][32]byte
	PubkeyWitnessesReceiver [][][32]byte
}

// TypesStateMerkleProof is an auto generated low-level Go binding around an user-defined struct.
type TypesStateMerkleProof struct {
	State   TypesUserState
	Witness [][32]byte
}

// TypesSubtreeVacancyProof is an auto generated low-level Go binding around an user-defined struct.
type TypesSubtreeVacancyProof struct {
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

// RollupMetaData contains all meta data concerning the Rollup contract.
var RollupMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractChooser\",\"name\":\"_chooser\",\"type\":\"address\"},{\"internalType\":\"contractIDepositManager\",\"name\":\"_depositManager\",\"type\":\"address\"},{\"internalType\":\"contractBLSAccountRegistry\",\"name\":\"_accountRegistry\",\"type\":\"address\"},{\"internalType\":\"contractTransfer\",\"name\":\"_transfer\",\"type\":\"address\"},{\"internalType\":\"contractMassMigration\",\"name\":\"_massMigration\",\"type\":\"address\"},{\"internalType\":\"contractCreate2Transfer\",\"name\":\"_create2Transfer\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"genesisStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"stakeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blocksToFinalise\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minGasLeft\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxTxsPerCommit\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subtreeID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"depositSubTreeRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"pathToSubTree\",\"type\":\"uint256\"}],\"name\":\"DepositsFinalised\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumTypes.Usage\",\"name\":\"batchType\",\"type\":\"uint8\"}],\"name\":\"NewBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nDeleted\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"completed\",\"type\":\"bool\"}],\"name\":\"RollbackStatus\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumTypes.Result\",\"name\":\"result\",\"type\":\"uint8\"}],\"name\":\"RollbackTriggered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"committed\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"}],\"name\":\"StakeWithdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DOMAIN_NAME\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DOMAIN_VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ZERO_BYTES32\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"accountRegistry\",\"outputs\":[{\"internalType\":\"contractBLSAccountRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"batches\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitmentRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"meta\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"chooser\",\"outputs\":[{\"internalType\":\"contractChooser\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"create2Transfer\",\"outputs\":[{\"internalType\":\"contractCreate2Transfer\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositManager\",\"outputs\":[{\"internalType\":\"contractIDepositManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"deposits\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.TransferBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"internalType\":\"structTypes.TransferCommitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.TransferCommitmentInclusionProof\",\"name\":\"target\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState[]\",\"name\":\"states\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"stateWitnesses\",\"type\":\"bytes32[][]\"},{\"internalType\":\"uint256[4][]\",\"name\":\"pubkeysSender\",\"type\":\"uint256[4][]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"pubkeyWitnessesSender\",\"type\":\"bytes32[][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"pubkeyHashesReceiver\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"pubkeyWitnessesReceiver\",\"type\":\"bytes32[][]\"}],\"internalType\":\"structTypes.SignatureProofWithReceiver\",\"name\":\"signatureProof\",\"type\":\"tuple\"}],\"name\":\"disputeSignatureCreate2Transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"spokeID\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"withdrawRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.MassMigrationBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"internalType\":\"structTypes.MassMigrationCommitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.MMCommitmentInclusionProof\",\"name\":\"target\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState[]\",\"name\":\"states\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"stateWitnesses\",\"type\":\"bytes32[][]\"},{\"internalType\":\"uint256[4][]\",\"name\":\"pubkeys\",\"type\":\"uint256[4][]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"pubkeyWitnesses\",\"type\":\"bytes32[][]\"}],\"internalType\":\"structTypes.SignatureProof\",\"name\":\"signatureProof\",\"type\":\"tuple\"}],\"name\":\"disputeSignatureMassMigration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.TransferBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"internalType\":\"structTypes.TransferCommitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.TransferCommitmentInclusionProof\",\"name\":\"target\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState[]\",\"name\":\"states\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"stateWitnesses\",\"type\":\"bytes32[][]\"},{\"internalType\":\"uint256[4][]\",\"name\":\"pubkeys\",\"type\":\"uint256[4][]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"pubkeyWitnesses\",\"type\":\"bytes32[][]\"}],\"internalType\":\"structTypes.SignatureProof\",\"name\":\"signatureProof\",\"type\":\"tuple\"}],\"name\":\"disputeSignatureTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"bodyRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structTypes.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.CommitmentInclusionProof\",\"name\":\"previous\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.TransferBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"internalType\":\"structTypes.TransferCommitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.TransferCommitmentInclusionProof\",\"name\":\"target\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState\",\"name\":\"state\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.StateMerkleProof[]\",\"name\":\"proofs\",\"type\":\"tuple[]\"}],\"name\":\"disputeTransitionCreate2Transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"bodyRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structTypes.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.CommitmentInclusionProof\",\"name\":\"previous\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"spokeID\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"withdrawRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.MassMigrationBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"internalType\":\"structTypes.MassMigrationCommitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.MMCommitmentInclusionProof\",\"name\":\"target\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState\",\"name\":\"state\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.StateMerkleProof[]\",\"name\":\"proofs\",\"type\":\"tuple[]\"}],\"name\":\"disputeTransitionMassMigration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"bodyRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structTypes.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.CommitmentInclusionProof\",\"name\":\"previous\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.TransferBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"internalType\":\"structTypes.TransferCommitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.TransferCommitmentInclusionProof\",\"name\":\"target\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState\",\"name\":\"state\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.StateMerkleProof[]\",\"name\":\"proofs\",\"type\":\"tuple[]\"}],\"name\":\"disputeTransitionTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"domainSeparator\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"}],\"name\":\"getBatch\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"commitmentRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"meta\",\"type\":\"bytes32\"}],\"internalType\":\"structTypes.Batch\",\"name\":\"batch\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"invalidBatchMarker\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"keepRollingBack\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"massMigration\",\"outputs\":[{\"internalType\":\"contractMassMigration\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextBatchID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paramBlocksToFinalise\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paramMaxTxsPerCommit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paramMinGasLeft\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paramStakeAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"stateRoots\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[2][]\",\"name\":\"signatures\",\"type\":\"uint256[2][]\"},{\"internalType\":\"uint256[]\",\"name\":\"feeReceivers\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"txss\",\"type\":\"bytes[]\"}],\"name\":\"submitCreate2Transfer\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"bodyRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structTypes.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.CommitmentInclusionProof\",\"name\":\"previous\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pathAtDepth\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.SubtreeVacancyProof\",\"name\":\"vacant\",\"type\":\"tuple\"}],\"name\":\"submitDeposits\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"stateRoots\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[2][]\",\"name\":\"signatures\",\"type\":\"uint256[2][]\"},{\"internalType\":\"uint256[4][]\",\"name\":\"meta\",\"type\":\"uint256[4][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"withdrawRoots\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes[]\",\"name\":\"txss\",\"type\":\"bytes[]\"}],\"name\":\"submitMassMigration\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"stateRoots\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[2][]\",\"name\":\"signatures\",\"type\":\"uint256[2][]\"},{\"internalType\":\"uint256[]\",\"name\":\"feeReceivers\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"txss\",\"type\":\"bytes[]\"}],\"name\":\"submitTransfer\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"contractTransfer\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"}],\"name\":\"withdrawStake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"withdrawalBitmap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"zeroHashAtSubtreeDepth\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6102a060405260006001553480156200001757600080fd5b5060405162005bf738038062005bf78339810160408190526200003a9162000921565b60405180604001604052806006815260200165487562626c6560d01b815250604051806040016040528060018152602001603160f81b8152508585858f8f84608081815250508360a081815250508260c08181525050816001600160a01b0316610100816001600160a01b031660601b81525050806001600160a01b031660e0816001600160a01b031660601b81525050505050505060008280519060200120905060008280519060200120905060007f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f9050826101808181525050816101a081815250504661014081815250506200013b8184846200035960201b60201c565b6101205230606090811b610160526101c0919091526001600160601b03198e821b81166101e0528d821b8116610200528c821b811661022052908b901b166102405250505061026082905250604080516338080fa560e21b815290516200021c916001600160a01b038d169163e0203e94916004808201926020929091908290030181600087803b158015620001d057600080fd5b505af1158015620001e5573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200020b9190620009e9565b620003a260201b620025e61760201c565b610280526040516000906200024f90879060008051602062005bd783398151915290602001918252602082015260400190565b60408051808303601f19018152828252805160209182012090830181905260008051602062005bd783398151915291830191909152915060009060600160408051601f1981840301815282825280516020918201208383019092528183529092508101620002cd600060013343620008a360201b62002ae31760201c565b90526001805460009081526020818152604080832085518155949091015193830193909355905491517f3c367559be37298dc40fed468321ed4d44d99ab593f8e0bb7c82dbc84ab0bad69262000327929091819062000a02565b60405180910390a160018054906000620003418362000a3d565b91905055505050505050505050505050505062000a7e565b6040805160208101859052908101839052606081018290524660808201523060a082015260009060c0016040516020818303038152906040528051906020012090509392505050565b6000620003ae62000902565b60008051602062005bd783398151915281527f633dc4d7da7256660a892f8f1604a44b5432649cc8ec5cb3ced4c4e6ac94dd1d6020808301919091527f890740a8eb06ce9be422cb8da5cdafc2b58c0a5e24036c578de2a433c828ff7d60408301527f3b8ec09e026fdc305365dfc94e189a81b38c7597b3d941c279f042e8206e0bd860608301527fecd50eee38e386bd62be9bedb990706951b65fe053bd9d8a521af753d139e2da60808301527fdefff6d330bb5403f63b14f33b578274160de3a50df4efecf0e0db73bcdd3da560a08301527f617bdd11f7c0a11f49db22f629387a12da7596f9d1704d7465177c63d88ec7d760c08301527f292c23a9aa1d8bea7e2435e555a4a60e379a5a35f3f452bae60121073fb6eead60e08301527fe1cea92ed99acdcb045a6726b2f87107e8a61620a232cf4d7d5b5766b3952e106101008301527f7ad66c0a68c72cb89e4fb4303841966e4062a76ab97451e3b9fb526a5ceb7f826101208301527fe026cc5a4aed3c22a58cbd3d2ac754c9352c5436f638042dca99034e836365166101408301527f3d04cffd8b46a874edf5cfae63077de85f849a660426697b06a829c70dd1409c6101608301527fad676aa337a485e4728a0b240d92b3ef7b3c372d06d189322bfd5f61f1e7203e6101808301527fa2fca4a49658f9fab7aa63289c91b7c7b6c832a6d0e69334ff5b0a3483d09dab6101a08301527f4ebfd9cd7bca2505f7bef59cc1c12ecc708fff26ae4af19abe852afe9e20c8626101c08301527f2def10d13dd169f550f578bda343d9717a138562e0093b380a1120789d53cf106101e08301527f776a31db34a1a0a7caaf862cffdfff1789297ffadc380bd3d39281d340abd3ad6102008301527fe2e7610b87a5fdf3a72ebe271287d923ab990eefac64b6e59d79f8b7e08c46e36102208301527f504364a5c6858bf98fff714ab5be9de19ed31a976860efbd0e772a2efe23e2e06102408301527f4f05f4acb83f5b65168d9fef89d56d4d77b8944015e6b1eed81b0238e2d0dba36102608301527f44a6d974c75b07423e1d6d33f481916fdd45830aea11b6347e700cd8b9f0767c6102808301527fedf260291f734ddac396a956127dde4c34c0cfb8d8052f88ac139658ccf2d5076102a08301527f6075c657a105351e7f0fce53bc320113324a522e8fd52dc878c762551e01a46e6102c08301527f6ca6a3f763a9395f7da16014725ca7ee17e4815c0ff8119bf33f273dee11833b6102e08301527f1c25ef10ffeb3c7d08aa707d17286e0b0d3cbcb50f1bd3b6523b63ba3b52dd0f6103008301527ffffc43bd08273ccf135fd3cacbeef055418e09eb728d727c4d5d5c556cdea7e36103208301527fc5ab8111456b1f28f3c7a0a604b4553ce905cb019c463ee159137af83c350b226103408301527f0ff273fcbf4ae0f2bd88d6cf319ff4004f8d7dca70d4ced4e74d2c74139739e66103608301527f7fa06ba11241ddd5efdc65d4e39c9f6991b74fd4b81b62230808216c876f827c6103808301527f7e275adf313a996c7e2950cac67caba02a5ff925ebf9906b58949f3e77aec5b96103a08301527f8f6162fa308d2b3a15dc33cffac85f13ab349173121645aedf00f471663108be6103c08301527f78ccaaab73373552f207a63599de54d7d8d0c1805f86ce7da15818d09f4cff626103e08301528190849081106200089757634e487b7160e01b600052603260045260246000fd5b60200201519392505050565b69ffffffff000000000000603082901b16600160501b600160f01b03605084901b1660ff60f01b60f086901b167fff0000000000000000000000000000000000000000000000000000000000000060f888901b16171717949350505050565b6040518061040001604052806020906020820280368337509192915050565b60008060008060008060008060008060006101608c8e03121562000943578687fd5b8b51620009508162000a65565b60208d0151909b50620009638162000a65565b60408d0151909a50620009768162000a65565b60608d0151909950620009898162000a65565b60808d01519098506200099c8162000a65565b60a08d0151909750620009af8162000a65565b8096505060c08c0151945060e08c015193506101008c015192506101208c015191506101408c015190509295989b509295989b9093969950565b600060208284031215620009fb578081fd5b5051919050565b83815260208101839052606081016005831062000a2f57634e487b7160e01b600052602160045260246000fd5b826040830152949350505050565b600060001982141562000a5e57634e487b7160e01b81526011600452602481fd5b5060010190565b6001600160a01b038116811462000a7b57600080fd5b50565b60805160a05160c05160e05160601c6101005160601c61012051610140516101605160601c610180516101a0516101c0516101e05160601c6102005160601c6102205160601c6102405160601c6102605161028051614fd362000c046000396000818161054d015261206a0152600081816106c6015281816118e901528181611daa01526125ac01526000818161045901528181611afe0152611d7201526000818161024f01528181610d76015261257d01526000818161048d015281816118b101526123a901526000818161065e015281816108ae015281816111e2015261154001526000613568015260006135b701526000613592015260006134eb015260006135150152600061353f0152600081816104c10152818161074c01528181611089015281816113e70152611e6c0152600081816103e6015281816120ec015261323b01526000818161062a01526131c70152600081816102ee0152612fd70152600081816106920152818161100501528181612f2e01526133210152614fd36000f3fe6080604052600436106101ee5760003560e01c80638a4068dd1161010d578063b6df3e6e116100a0578063d1b243e21161006f578063d1b243e214610680578063e42f7554146106b4578063eb84e076146106e8578063f698da2514610708578063fe5f96041461071d57600080fd5b8063b6df3e6e146105e5578063c08caaef146105f8578063ca858d9d14610618578063d089e11a1461064c57600080fd5b8063acb8cc49116100dc578063acb8cc491461050e578063acf5b54d1461053b578063b02c43d01461056f578063b32c4d8d1461059c57600080fd5b80638a4068dd1461047b57806398d17621146104af578063a72598ce146104e3578063ac96f0cd146104f957600080fd5b80634e23e8c3116101855780635f6e91d5116101545780635f6e91d5146103b45780636c7ac9d8146103d4578063796f077b146104085780637ae8c5681461044757600080fd5b80634e23e8c3146103235780634f6de740146103435780635ac44282146103635780635b097d371461039e57600080fd5b806325d5971f116101c157806325d5971f146102a95780632e4518d8146102c957806331c2b7db146102dc57806339c139831461031057600080fd5b80630564e8b3146101f3578063069321b0146102085780630ed75b9c1461023d5780632538507d14610289575b600080fd5b610206610201366004613f9b565b61074a565b005b34801561021457600080fd5b5061022a600080516020614fa783398151915281565b6040519081526020015b60405180910390f35b34801561024957600080fd5b506102717f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b039091168152602001610234565b34801561029557600080fd5b506102066102a4366004614324565b610c43565b3480156102b557600080fd5b506102066102c4366004613f83565b610e5b565b6102066102d7366004614094565b611087565b3480156102e857600080fd5b5061022a7f000000000000000000000000000000000000000000000000000000000000000081565b61020661031e366004614094565b6113e5565b34801561032f57600080fd5b5061020661033e3660046142c5565b611734565b34801561034f57600080fd5b5061020661035e36600461438d565b6119e5565b34801561036f57600080fd5b5061038361037e366004613f83565b611b35565b60408051825181526020928301519281019290925201610234565b3480156103aa57600080fd5b5061022a60045481565b3480156103c057600080fd5b506102066103cf3660046142c5565b611bf5565b3480156103e057600080fd5b506102717f000000000000000000000000000000000000000000000000000000000000000081565b34801561041457600080fd5b5061043a60405180604001604052806006815260200165487562626c6560d01b81525081565b604051610234919061497f565b34801561045357600080fd5b506102717f000000000000000000000000000000000000000000000000000000000000000081565b34801561048757600080fd5b506102717f000000000000000000000000000000000000000000000000000000000000000081565b3480156104bb57600080fd5b506102717f000000000000000000000000000000000000000000000000000000000000000081565b3480156104ef57600080fd5b5061022a60015481565b34801561050557600080fd5b50610206611e04565b34801561051a57600080fd5b5061043a604051806040016040528060018152602001603160f81b81525081565b34801561054757600080fd5b5061022a7f000000000000000000000000000000000000000000000000000000000000000081565b34801561057b57600080fd5b5061022a61058a366004613f83565b60026020526000908152604090205481565b3480156105a857600080fd5b506105d06105b7366004613f83565b6000602081905290815260409020805460019091015482565b60408051928352602083019190915201610234565b6102066105f336600461420e565b611e6a565b34801561060457600080fd5b506102066106133660046144d3565b612293565b34801561062457600080fd5b5061022a7f000000000000000000000000000000000000000000000000000000000000000081565b34801561065857600080fd5b506102717f000000000000000000000000000000000000000000000000000000000000000081565b34801561068c57600080fd5b5061022a7f000000000000000000000000000000000000000000000000000000000000000081565b3480156106c057600080fd5b5061022a7f000000000000000000000000000000000000000000000000000000000000000081565b3480156106f457600080fd5b50610206610703366004614181565b6123e0565b34801561071457600080fd5b5061022a6125d7565b34801561072957600080fd5b5061022a610738366004613f83565b60036020526000908152604090205481565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b1580156107a357600080fd5b505afa1580156107b7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107db9190613eff565b6001600160a01b0316336001600160a01b0316146108145760405162461bcd60e51b815260040161080b90614b23565b60405180910390fd5b600454156108345760405162461bcd60e51b815260040161080b90614ab5565b8a60015481146108565760405162461bcd60e51b815260040161080b90614bad565b60008a6001600160401b0381111561087e57634e487b7160e01b600052604160045260246000fd5b6040519080825280602002602001820160405280156108a7578160200160208202803683370190505b50905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663ebf0c7176040518163ffffffff1660e01b815260040160206040518083038186803b15801561090557600080fd5b505afa158015610919573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061093d9190613f26565b905060005b8c811015610c1d5760006040518061010001604052808481526020018e8e8581811061097e57634e487b7160e01b600052603260045260246000fd5b9050604002016002806020026040519081016040528092919082600260200280828437600092019190915250505081526020018c8c858181106109d157634e487b7160e01b600052603260045260246000fd5b9050608002016000600481106109f757634e487b7160e01b600052603260045260246000fd5b602002013581526020018a8a85818110610a2157634e487b7160e01b600052603260045260246000fd5b9050602002013581526020018c8c85818110610a4d57634e487b7160e01b600052603260045260246000fd5b905060800201600160048110610a7357634e487b7160e01b600052603260045260246000fd5b602002013581526020018c8c85818110610a9d57634e487b7160e01b600052603260045260246000fd5b905060800201600260048110610ac357634e487b7160e01b600052603260045260246000fd5b602002013581526020018c8c85818110610aed57634e487b7160e01b600052603260045260246000fd5b905060800201600360048110610b1357634e487b7160e01b600052603260045260246000fd5b60200201358152602001888885818110610b3d57634e487b7160e01b600052603260045260246000fd5b9050602002810190610b4f9190614d49565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250505091525090508e8e83818110610ba757634e487b7160e01b600052603260045260246000fd5b90506020020135610bb782612b2a565b60408051602081019390935282015260600160405160208183030381529060405280519060200120848381518110610bff57634e487b7160e01b600052603260045260246000fd5b60209081029190910101525080610c1581614f1f565b915050610942565b50610c33610c2a83612b81565b8d836002612f2c565b5050505050505050505050505050565b60008381526020818152604091829020825180840190935280548352600101549101819052839060301c63ffffffff164310610c915760405162461bcd60e51b815260040161080b90614aec565b600454811080610ca15750600454155b610cbd5760405162461bcd60e51b815260040161080b90614b5a565b600084815260208190526040902054610cd6908461307a565b610d225760405162461bcd60e51b815260206004820152601f60248201527f436f6d6d69746d656e74206e6f742070726573656e7420696e20626174636800604482015260640161080b565b6040805160a081018252845160209081015181015182528551518183015285510151519181019190915260009060608101610d5b6125d7565b815260200185600001516020015160e00151815250905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663be42c7c48386886000015160200151604001516040518463ffffffff1660e01b8152600401610dd093929190614cc3565b60206040518083038186803b158015610de857600080fd5b505afa158015610dfc573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e209190613f69565b9050600081600b811115610e4457634e487b7160e01b600052602160045260246000fd5b14610e5357610e5386826130a3565b505050505050565b6000818152602081815260409182902082518084019093528054835260010154910181905260501c6001600160a01b03166001600160a01b0316336001600160a01b031614610f055760405162461bcd60e51b815260206004820152603060248201527f596f7520617265206e6f742074686520636f727265637420636f6d6d6974746560448201526f0e440ccdee440e8d0d2e640c4c2e8c6d60831b606482015260840161080b565b6000818152602081815260409182902082518084019093528054835260010154910181905260301c63ffffffff164311610f9b5760405162461bcd60e51b815260206004820152603160248201527f54686973206261746368206973206e6f74207965742066696e616c697365642c60448201527020636865636b206261636b20736f6f6e2160781b606482015260840161080b565b610fa68160036130ee565b15610ff35760405162461bcd60e51b815260206004820152601960248201527f526f6c6c75703a20416c72656164792077697468647261776e00000000000000604482015260640161080b565b610ffe81600361312f565b60405133907f000000000000000000000000000000000000000000000000000000000000000080156108fc02916000818181858888f1935050505015801561104a573d6000803e3d6000fd5b5060408051338152602081018390527f1248d48e2de900a1010c7fce73506969ecec243600bfc08b641b158f26d857cd910160405180910390a150565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b1580156110e057600080fd5b505afa1580156110f4573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906111189190613eff565b6001600160a01b0316336001600160a01b0316146111485760405162461bcd60e51b815260040161080b90614b23565b600454156111685760405162461bcd60e51b815260040161080b90614ab5565b88600154811461118a5760405162461bcd60e51b815260040161080b90614bad565b6000886001600160401b038111156111b257634e487b7160e01b600052604160045260246000fd5b6040519080825280602002602001820160405280156111db578160200160208202803683370190505b50905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663ebf0c7176040518163ffffffff1660e01b815260040160206040518083038186803b15801561123957600080fd5b505afa15801561124d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906112719190613f26565b90506000805b8b8110156113c057828b8b838181106112a057634e487b7160e01b600052603260045260246000fd5b9050604002018a8a848181106112c657634e487b7160e01b600052603260045260246000fd5b905060200201358989858181106112ed57634e487b7160e01b600052603260045260246000fd5b90506020028101906112ff9190614d49565b6040516020016113139594939291906147dd565b6040516020818303038152906040528051906020012091508c8c8281811061134b57634e487b7160e01b600052603260045260246000fd5b905060200201358260405160200161136d929190918252602082015260400190565b604051602081830303815290604052805190602001208482815181106113a357634e487b7160e01b600052603260045260246000fd5b6020908102919091010152806113b881614f1f565b915050611277565b506113d66113cd84612b81565b8c846003612f2c565b50505050505050505050505050565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b15801561143e57600080fd5b505afa158015611452573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114769190613eff565b6001600160a01b0316336001600160a01b0316146114a65760405162461bcd60e51b815260040161080b90614b23565b600454156114c65760405162461bcd60e51b815260040161080b90614ab5565b8860015481146114e85760405162461bcd60e51b815260040161080b90614bad565b6000886001600160401b0381111561151057634e487b7160e01b600052604160045260246000fd5b604051908082528060200260200182016040528015611539578160200160208202803683370190505b50905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663ebf0c7176040518163ffffffff1660e01b815260040160206040518083038186803b15801561159757600080fd5b505afa1580156115ab573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115cf9190613f26565b90506000805b8b81101561171e57828b8b838181106115fe57634e487b7160e01b600052603260045260246000fd5b9050604002018a8a8481811061162457634e487b7160e01b600052603260045260246000fd5b9050602002013589898581811061164b57634e487b7160e01b600052603260045260246000fd5b905060200281019061165d9190614d49565b6040516020016116719594939291906147dd565b6040516020818303038152906040528051906020012091508c8c828181106116a957634e487b7160e01b600052603260045260246000fd5b90506020020135826040516020016116cb929190918252602082015260400190565b6040516020818303038152906040528051906020012084828151811061170157634e487b7160e01b600052603260045260246000fd5b60209081029190910101528061171681614f1f565b9150506115d5565b506113d661172b84612b81565b8c846001612f2c565b60008481526020818152604091829020825180840190935280548352600101549101819052849060301c63ffffffff1643106117825760405162461bcd60e51b815260040161080b90614aec565b6004548110806117925750600454155b6117ae5760405162461bcd60e51b815260040161080b90614b5a565b84848460200151600080826000141561180f576117cc600186614edc565b600081815260208181526040918290208251808401909352805483526001908101549290910182905291925060f01c60ff166118089190614edc565b915061181f565b508361181c600184614edc565b91505b818460200151146118425760405162461bcd60e51b815260040161080b90614a73565b60008181526020819052604090205461185b908561316c565b6118775760405162461bcd60e51b815260040161080b90614992565b60008a815260208190526040902054611890908961317f565b6118ac5760405162461bcd60e51b815260040161080b906149e4565b6000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663929314928c60000151600001517f00000000000000000000000000000000000000000000000000000000000000008d6000015160200151604001518e6000015160200151606001518e6040518663ffffffff1660e01b815260040161194395949392919061493c565b604080518083038186803b15801561195a57600080fd5b505afa15801561196e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119929190613f3e565b9092509050600081600b8111156119b957634e487b7160e01b600052602160045260246000fd5b1415806119c857508951518214155b156119d7576119d78c826130a3565b505050505050505050505050565b60008381526020818152604091829020825180840190935280548352600101549101819052839060301c63ffffffff164310611a335760405162461bcd60e51b815260040161080b90614aec565b600454811080611a435750600454155b611a5f5760405162461bcd60e51b815260040161080b90614b5a565b600084815260208190526040902054611a78908461317f565b611a945760405162461bcd60e51b815260040161080b90614a2c565b6040805160a081018252845160209081015181015182528551518183015285510151519181019190915260009060608101611acd6125d7565b81528551602090810151606001519101526040516001627e8b6b60e01b031981529091506000906001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063ff81749590610dd09085908890600401614bef565b6040805180820190915260008082526020820152600082815260208190526040902060010154611bcd5760405162461bcd60e51b815260206004820152603f60248201527f42617463682069642067726561746572207468616e20746f74616c206e756d6260448201527f6572206f6620626174636865732c20696e76616c696420626174636820696400606482015260840161080b565b5060009081526020818152604091829020825180840190935280548352600101549082015290565b60008481526020818152604091829020825180840190935280548352600101549101819052849060301c63ffffffff164310611c435760405162461bcd60e51b815260040161080b90614aec565b600454811080611c535750600454155b611c6f5760405162461bcd60e51b815260040161080b90614b5a565b848484602001516000808260001415611cd057611c8d600186614edc565b600081815260208181526040918290208251808401909352805483526001908101549290910182905291925060f01c60ff16611cc99190614edc565b9150611ce0565b5083611cdd600184614edc565b91505b81846020015114611d035760405162461bcd60e51b815260040161080b90614a73565b600081815260208190526040902054611d1c908561316c565b611d385760405162461bcd60e51b815260040161080b90614992565b60008a815260208190526040902054611d51908961317f565b611d6d5760405162461bcd60e51b815260040161080b906149e4565b6000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663336920368c60000151600001517f00000000000000000000000000000000000000000000000000000000000000008d6000015160200151604001518e6000015160200151606001518e6040518663ffffffff1660e01b815260040161194395949392919061493c565b600060045411611e605760405162461bcd60e51b815260206004820152602160248201527f42617463684d616e616765723a204973206e6f7420726f6c6c696e67206261636044820152606b60f81b606482015260840161080b565b611e68613192565b565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b158015611ec357600080fd5b505afa158015611ed7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611efb9190613eff565b6001600160a01b0316336001600160a01b031614611f2b5760405162461bcd60e51b815260040161080b90614b23565b60045415611f4b5760405162461bcd60e51b815260040161080b90614ab5565b826001548114611f6d5760405162461bcd60e51b815260040161080b90614bad565b6000611f7a600186614edc565b600081815260208181526040918290208251808401909352805483526001908101549290910182905291925060f01c60ff16611fb69190614edc565b846020015114611fd85760405162461bcd60e51b815260040161080b90614a73565b600081815260208190526040902054611ff1908561316c565b6120595760405162461bcd60e51b815260206004820152603360248201527f70726576696f757320636f6d6d69746d656e7420697320616273656e7420696e604482015272040e8d0ca40e0e4caecd2deeae640c4c2e8c6d606b1b606482015260840161080b565b8351518351602085015161208f92917f0000000000000000000000000000000000000000000000000000000000000000916133d7565b6120e75760405162461bcd60e51b815260206004820152602360248201527f526f6c6c75703a2053746174652073756274726565206973206e6f7420766163604482015262185b9d60ea1b606482015260840161080b565b6000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663d86ee48d6040518163ffffffff1660e01b81526004016040805180830381600087803b15801561214457600080fd5b505af1158015612158573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061217c919061415e565b9092509050600061218e846001614e91565b600081815260026020908152604091829020859055885182518781529182018690528183015290519192507f1e6cd0ae9aa6f6e6ea5cf7f07396a96fc1ea5ae4ec7d5054f53e204c47f79bd8919081900360600190a160006121f983886000015189602001516133ef565b9050600081600080516020614fa783398151915260001b60405160200161222a929190918252602082015260400190565b60408051601f198184030181528282528051602091820120818401819052600080516020614fa7833981519152848401528251808503840181526060909401909252825192019190912090915061228681600160006004612f2c565b5050505050505050505050565b60008381526020818152604091829020825180840190935280548352600101549101819052839060301c63ffffffff1643106122e15760405162461bcd60e51b815260040161080b90614aec565b6004548110806122f15750600454155b61230d5760405162461bcd60e51b815260040161080b90614b5a565b600084815260208190526040902054612326908461317f565b6123425760405162461bcd60e51b815260040161080b90614a2c565b6040805160a08101825284516020908101518101518252855151818301528551015151918101919091526000906060810161237b6125d7565b81528551602090810151606001519101526040516344ec5a7760e01b81529091506000906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906344ec5a7790610dd09085908890600401614c9e565b60008481526020818152604091829020825180840190935280548352600101549101819052849060301c63ffffffff16431061242e5760405162461bcd60e51b815260040161080b90614aec565b60045481108061243e5750600454155b61245a5760405162461bcd60e51b815260040161080b90614b5a565b8484846020015160008082600014156124bb57612478600186614edc565b600081815260208181526040918290208251808401909352805483526001908101549290910182905291925060f01c60ff166124b49190614edc565b91506124cb565b50836124c8600184614edc565b91505b818460200151146124ee5760405162461bcd60e51b815260040161080b90614a73565b600081815260208190526040902054612507908561316c565b6125235760405162461bcd60e51b815260040161080b90614992565b60008a81526020819052604090205461253c908961307a565b6125585760405162461bcd60e51b815260040161080b906149e4565b88515188516020015160405163ab5a164f60e01b815260009283926001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169263ab5a164f926119439290917f000000000000000000000000000000000000000000000000000000000000000091908f9060040161489c565b60006125e16134de565b905090565b60006125f06136b7565b600080516020614fa783398151915281527f633dc4d7da7256660a892f8f1604a44b5432649cc8ec5cb3ced4c4e6ac94dd1d6020808301919091527f890740a8eb06ce9be422cb8da5cdafc2b58c0a5e24036c578de2a433c828ff7d60408301527f3b8ec09e026fdc305365dfc94e189a81b38c7597b3d941c279f042e8206e0bd860608301527fecd50eee38e386bd62be9bedb990706951b65fe053bd9d8a521af753d139e2da60808301527fdefff6d330bb5403f63b14f33b578274160de3a50df4efecf0e0db73bcdd3da560a08301527f617bdd11f7c0a11f49db22f629387a12da7596f9d1704d7465177c63d88ec7d760c08301527f292c23a9aa1d8bea7e2435e555a4a60e379a5a35f3f452bae60121073fb6eead60e08301527fe1cea92ed99acdcb045a6726b2f87107e8a61620a232cf4d7d5b5766b3952e106101008301527f7ad66c0a68c72cb89e4fb4303841966e4062a76ab97451e3b9fb526a5ceb7f826101208301527fe026cc5a4aed3c22a58cbd3d2ac754c9352c5436f638042dca99034e836365166101408301527f3d04cffd8b46a874edf5cfae63077de85f849a660426697b06a829c70dd1409c6101608301527fad676aa337a485e4728a0b240d92b3ef7b3c372d06d189322bfd5f61f1e7203e6101808301527fa2fca4a49658f9fab7aa63289c91b7c7b6c832a6d0e69334ff5b0a3483d09dab6101a08301527f4ebfd9cd7bca2505f7bef59cc1c12ecc708fff26ae4af19abe852afe9e20c8626101c08301527f2def10d13dd169f550f578bda343d9717a138562e0093b380a1120789d53cf106101e08301527f776a31db34a1a0a7caaf862cffdfff1789297ffadc380bd3d39281d340abd3ad6102008301527fe2e7610b87a5fdf3a72ebe271287d923ab990eefac64b6e59d79f8b7e08c46e36102208301527f504364a5c6858bf98fff714ab5be9de19ed31a976860efbd0e772a2efe23e2e06102408301527f4f05f4acb83f5b65168d9fef89d56d4d77b8944015e6b1eed81b0238e2d0dba36102608301527f44a6d974c75b07423e1d6d33f481916fdd45830aea11b6347e700cd8b9f0767c6102808301527fedf260291f734ddac396a956127dde4c34c0cfb8d8052f88ac139658ccf2d5076102a08301527f6075c657a105351e7f0fce53bc320113324a522e8fd52dc878c762551e01a46e6102c08301527f6ca6a3f763a9395f7da16014725ca7ee17e4815c0ff8119bf33f273dee11833b6102e08301527f1c25ef10ffeb3c7d08aa707d17286e0b0d3cbcb50f1bd3b6523b63ba3b52dd0f6103008301527ffffc43bd08273ccf135fd3cacbeef055418e09eb728d727c4d5d5c556cdea7e36103208301527fc5ab8111456b1f28f3c7a0a604b4553ce905cb019c463ee159137af83c350b226103408301527f0ff273fcbf4ae0f2bd88d6cf319ff4004f8d7dca70d4ced4e74d2c74139739e66103608301527f7fa06ba11241ddd5efdc65d4e39c9f6991b74fd4b81b62230808216c876f827c6103808301527f7e275adf313a996c7e2950cac67caba02a5ff925ebf9906b58949f3e77aec5b96103a08301527f8f6162fa308d2b3a15dc33cffac85f13ab349173121645aedf00f471663108be6103c08301527f78ccaaab73373552f207a63599de54d7d8d0c1805f86ce7da15818d09f4cff626103e0830152819084908110612ad757634e487b7160e01b600052603260045260246000fd5b60200201519392505050565b69ffffffff000000000000603082901b16600160501b600160f01b03605084901b1660ff60f01b60f086901b166001600160f81b031960f888901b16171717949350505050565b80516020808301516040808501516060860151608087015160a088015160c089015160e08a01519551600099612b64999098979101614808565b604051602081830303815290604052805190602001209050919050565b6000602082511115612bd55760405162461bcd60e51b815260206004820152601b60248201527f4d65726b6c65547265653a20546f6f206d616e79206c65617665730000000000604482015260640161080b565b612bdd6136d6565b600080516020614fa783398151915281527f633dc4d7da7256660a892f8f1604a44b5432649cc8ec5cb3ced4c4e6ac94dd1d60208201527f890740a8eb06ce9be422cb8da5cdafc2b58c0a5e24036c578de2a433c828ff7d60408201527f3b8ec09e026fdc305365dfc94e189a81b38c7597b3d941c279f042e8206e0bd860608201527fecd50eee38e386bd62be9bedb990706951b65fe053bd9d8a521af753d139e2da60808201527fdefff6d330bb5403f63b14f33b578274160de3a50df4efecf0e0db73bcdd3da560a0820152825160018181169160009190612cc29082614e91565b901c90506000816001600160401b03811115612cee57634e487b7160e01b600052604160045260246000fd5b604051908082528060200260200182016040528015612d17578160200160208202803683370190505b5090508560005b60005b612d2b8686614edc565b811015612e05576000600182901b9050838181518110612d5b57634e487b7160e01b600052603260045260246000fd5b602002602001015184826001612d719190614e91565b81518110612d8f57634e487b7160e01b600052603260045260246000fd5b6020026020010151604051602001612db1929190918252602082015260400190565b60405160208183030381529060405280519060200120858381518110612de757634e487b7160e01b600052603260045260246000fd5b60209081029190910101525080612dfd81614f1f565b915050612d21565b8560011415612ebc5782600182901b81518110612e3257634e487b7160e01b600052603260045260246000fd5b6020026020010151878360068110612e5a57634e487b7160e01b600052603260045260246000fd5b6020020151604051602001612e79929190918252602082015260400190565b60405160208183030381529060405280519060200120848281518110612eaf57634e487b7160e01b600052603260045260246000fd5b6020026020010181815250505b8460011415612ecb5750612ef6565b60018086169650612edc8682614e91565b901c9450612eeb600183614e91565b915083925050612d1e565b82600081518110612f1757634e487b7160e01b600052603260045260246000fd5b60200260200101519650505050505050919050565b7f0000000000000000000000000000000000000000000000000000000000000000341015612f9c5760405162461bcd60e51b815260206004820152601a60248201527f526f6c6c75703a2077726f6e67207374616b6520616d6f756e74000000000000604482015260640161080b565b6040518060400160405280858152602001613001836004811115612fd057634e487b7160e01b600052602160045260246000fd5b8633612ffc7f000000000000000000000000000000000000000000000000000000000000000043614e91565b612ae3565b9052600180546000908152602081815260409182902084518155930151928201929092555490517f3c367559be37298dc40fed468321ed4d44d99ab593f8e0bb7c82dbc84ab0bad6916130579185908590614cf9565b60405180910390a16001805490600061306f83614f1f565b919050555050505050565b600061309c8361308d8460000151613605565b846020015185604001516133d7565b9392505050565b60048290556040517f79e7023a0e95197fd8e57153afa9dca5bb6ccf11c7b8a99b391eb11ba75d0704906130da9084908490614d25565b60405180910390a16130ea613192565b5050565b6000806130fd61010085614ea9565b9050600061310d61010086614f3a565b6000928352602094909452506040902054600190921b91821690911492915050565b600061313d61010084614ea9565b9050600061314d61010085614f3a565b600092835260209390935250604090208054600190921b909117905550565b600061309c8361308d8460000151613630565b600061309c8361308d8460000151613656565b60008060006004546001546131a79190614edc565b9050600080600180546131ba9190614edc565b90505b828210156132b2577f00000000000000000000000000000000000000000000000000000000000000005a116131f1576132b2565b6131fb8282614edc565b60008181526020818152604080832083815560010183905560029091529020549550935084156132a0576040516356f0001360e11b8152600481018690527f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169063ade0002690602401600060405180830381600087803b15801561328757600080fd5b505af115801561329b573d6000803e3d6000fd5b505050505b816132aa81614f1f565b9250506131bd565b81600160008282546132c49190614edc565b909155505081831480156132d85760006004555b60408051838152602081018590528215158183015290517f595cb200f38fcdbf73429ffce272e53b2e923c231b2059ce20d1f4909cff1bec9181900360600190a16000613345847f000000000000000000000000000000000000000000000000000000000000000061366a565b9050600061335f600361335984600261366a565b90613676565b9050600061336d8383613682565b604051909150339083156108fc029084906000818181858888f1935050505015801561339d573d6000803e3d6000fd5b5060405160009082156108fc0290839083818181858288f193505050501580156133cb573d6000803e3d6000fd5b50505050505050505050565b6000846133e58585856133ef565b1495945050505050565b600083815b83518110156134d557600185821c16613467578184828151811061342857634e487b7160e01b600052603260045260246000fd5b602002602001015160405160200161344a929190918252602082015260400190565b6040516020818303038152906040528051906020012091506134c3565b83818151811061348757634e487b7160e01b600052603260045260246000fd5b6020026020010151826040516020016134aa929190918252602082015260400190565b6040516020818303038152906040528051906020012091505b806134cd81614f1f565b9150506133f4565b50949350505050565b6000306001600160a01b037f00000000000000000000000000000000000000000000000000000000000000001614801561353757507f000000000000000000000000000000000000000000000000000000000000000046145b1561356157507f000000000000000000000000000000000000000000000000000000000000000090565b50604080517f00000000000000000000000000000000000000000000000000000000000000006020808301919091527f0000000000000000000000000000000000000000000000000000000000000000828401527f000000000000000000000000000000000000000000000000000000000000000060608301524660808301523060a0808401919091528351808403909101815260c0909201909252805191012090565b600081600001516136198360200151612b2a565b604080516020810193909352820152606001612b64565b600081600001518260200151604051602001612b64929190918252602082015260400190565b60008160000151613619836020015161368e565b600061309c8284614ebd565b600061309c8284614ea9565b600061309c8284614edc565b60008160000151826020015183604001518460600151604051602001612b649493929190614861565b6040518061040001604052806020906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b600082601f830112613704578081fd5b8135602061371961371483614e6e565b614e3e565b80838252828201915082860187848660051b8901011115613738578586fd5b855b858110156137785781356001600160401b03811115613757578788fd5b6137658a87838c0101613904565b855250928401929084019060010161373a565b5090979650505050505050565b60008083601f840112613796578182fd5b5081356001600160401b038111156137ac578182fd5b6020830191508360208260061b85010111156137c757600080fd5b9250929050565b60008083601f8401126137df578182fd5b5081356001600160401b038111156137f5578182fd5b6020830191508360208260071b85010111156137c757600080fd5b6000601f8381840112613821578182fd5b8235602061383161371483614e6e565b80838252828201915082870188848660071b8a01011115613850578687fd5b865b858110156138b4578987830112613867578788fd5b61386f614dfa565b8083608085018d811115613881578b8cfd5b8b5b60048110156138a057823585529389019391890191600101613883565b509187525094860194925050600101613852565b509098975050505050505050565b60008083601f8401126138d3578182fd5b5081356001600160401b038111156138e9578182fd5b6020830191508360208260051b85010111156137c757600080fd5b600082601f830112613914578081fd5b8135602061392461371483614e6e565b80838252828201915082860187848660051b8901011115613943578586fd5b855b8581101561377857813584529284019290840190600101613945565b600082601f830112613971578081fd5b8135602061398161371483614e6e565b80838252828201915082860187848660051b89010111156139a0578586fd5b855b858110156137785781356001600160401b03808211156139c0578889fd5b9089019060a0828c03601f19018113156139d857898afd5b6139e0614d8d565b6139ec8d8a8601613ebc565b81529083013590828211156139ff578a8bfd5b613a0d8d8a84870101613904565b818a0152875250505092840192908401906001016139a2565b600082601f830112613a36578081fd5b81356020613a4661371483614e6e565b80838252828201915082860187848660071b8901011115613a65578586fd5b855b8581101561377857613a798983613ebc565b84529284019260809190910190600101613a67565b600082601f830112613a9e578081fd5b604051604081018181106001600160401b0382111715613ac057613ac0614f90565b8060405250808385604086011115613ad6578384fd5b835b6002811015613af7578135835260209283019290910190600101613ad8565b509195945050505050565b600082601f830112613b12578081fd5b81356001600160401b03811115613b2b57613b2b614f90565b613b3e601f8201601f1916602001614e3e565b818152846020838601011115613b52578283fd5b816020850160208301379081016020019190915292915050565b8051600c8110613b7b57600080fd5b919050565b60008183036080811215613b92578182fd5b613b9a614db5565b91506040811215613baa57600080fd5b50613bb3614d8d565b8235815260208301356020820152808252506040820135602082015260608201356001600160401b03811115613be857600080fd5b613bf484828501613904565b60408301525092915050565b600060608284031215613c11578081fd5b613c19614db5565b905081356001600160401b0380821115613c3257600080fd5b9083019060408286031215613c4657600080fd5b613c4e614d8d565b8235815260208084013583811115613c6557600080fd5b93909301926101208488031215613c7b57600080fd5b613c83614dd7565b84358152613c9388838701613a8e565b82820152606085013560408201526080850135606082015260a0850135608082015260c085013560a082015260e085013560c082015261010085013584811115613cdc57600080fd5b613ce889828801613b02565b60e0830152508282015290845284810135908401526040840135915080821115613d1157600080fd5b50613bf484828501613904565b600060808284031215613d2f578081fd5b613d37614dfa565b905081356001600160401b0380821115613d5057600080fd5b613d5c85838601613a26565b83526020840135915080821115613d7257600080fd5b613d7e858386016136f4565b60208401526040840135915080821115613d9757600080fd5b613da385838601613810565b60408401526060840135915080821115613dbc57600080fd5b50613dc9848285016136f4565b60608301525092915050565b600060608284031215613de6578081fd5b613dee614db5565b905081356001600160401b0380821115613e0757600080fd5b9083019060408286031215613e1b57600080fd5b613e23614d8d565b8235815260208084013583811115613e3a57600080fd5b939093019260a08488031215613e4f57600080fd5b613e57614dfa565b84358152613e6788838701613a8e565b8282015260608501356040820152608085013584811115613e8757600080fd5b613e9389828801613b02565b6060830152508282015290845284810135908401526040840135915080821115613d1157600080fd5b600060808284031215613ecd578081fd5b613ed5614dfa565b90508135815260208201356020820152604082013560408201526060820135606082015292915050565b600060208284031215613f10578081fd5b81516001600160a01b038116811461309c578182fd5b600060208284031215613f37578081fd5b5051919050565b60008060408385031215613f50578081fd5b82519150613f6060208401613b6c565b90509250929050565b600060208284031215613f7a578081fd5b61309c82613b6c565b600060208284031215613f94578081fd5b5035919050565b600080600080600080600080600080600060c08c8e031215613fbb578687fd5b8b359a506001600160401b038060208e01351115613fd7578788fd5b613fe78e60208f01358f016138c2565b909b50995060408d0135811015613ffc578788fd5b61400c8e60408f01358f01613785565b909950975060608d0135811015614021578687fd5b6140318e60608f01358f016137ce565b909750955060808d0135811015614046578485fd5b6140568e60808f01358f016138c2565b909550935060a08d013581101561406b578283fd5b5061407c8d60a08e01358e016138c2565b81935080925050509295989b509295989b9093969950565b600080600080600080600080600060a08a8c0312156140b1578283fd5b8935985060208a01356001600160401b03808211156140ce578485fd5b6140da8d838e016138c2565b909a50985060408c01359150808211156140f2578485fd5b6140fe8d838e01613785565b909850965060608c0135915080821115614116578485fd5b6141228d838e016138c2565b909650945060808c013591508082111561413a578384fd5b506141478c828d016138c2565b915080935050809150509295985092959850929598565b60008060408385031215614170578182fd5b505080516020909101519092909150565b60008060008060808587031215614196578182fd5b8435935060208501356001600160401b03808211156141b3578384fd5b6141bf88838901613b80565b945060408701359150808211156141d4578384fd5b6141e088838901613c00565b935060608701359150808211156141f5578283fd5b5061420287828801613961565b91505092959194509250565b600080600060608486031215614222578081fd5b8335925060208401356001600160401b038082111561423f578283fd5b61424b87838801613b80565b93506040860135915080821115614260578283fd5b9085019060408288031215614273578283fd5b60405160408101818110838211171561428e5761428e614f90565b604052823581526020830135828111156142a6578485fd5b6142b289828601613904565b6020830152508093505050509250925092565b600080600080608085870312156142da578182fd5b8435935060208501356001600160401b03808211156142f7578384fd5b61430388838901613b80565b94506040870135915080821115614318578384fd5b6141e088838901613dd5565b600080600060608486031215614338578081fd5b8335925060208401356001600160401b0380821115614355578283fd5b61436187838801613c00565b93506040860135915080821115614376578283fd5b5061438386828701613d1e565b9150509250925092565b6000806000606084860312156143a1578081fd5b8335925060208401356001600160401b03808211156143be578283fd5b6143ca87838801613dd5565b935060408601359150808211156143df578283fd5b9085019060c082880312156143f2578283fd5b6143fa614e1c565b823582811115614408578485fd5b61441489828601613a26565b825250602083013582811115614428578485fd5b614434898286016136f4565b60208301525060408301358281111561444b578485fd5b61445789828601613810565b60408301525060608301358281111561446e578485fd5b61447a898286016136f4565b606083015250608083013582811115614491578485fd5b61449d89828601613904565b60808301525060a0830135828111156144b4578485fd5b6144c0898286016136f4565b60a0830152508093505050509250925092565b6000806000606084860312156144e7578081fd5b8335925060208401356001600160401b0380821115614504578283fd5b61436187838801613dd5565b600081518084526020808501808196508360051b81019150828601855b858110156145575782840389526145458483516145c2565b9885019893509084019060010161452d565b5091979650505050505050565b6000815180845260208085019450808401835b838110156145b757815187865b60048110156145a157825182529185019190850190600101614584565b5050506080969096019590820190600101614577565b509495945050505050565b6000815180845260208085019450808401835b838110156145b7578151875295820195908201906001016145d5565b600081518084526020808501808196508360051b81019150828601855b85811015614557578284038952815160a061464b868351805182526020810151602083015260408101516040830152606081015160608301525050565b868201519150806080870152614663818701836145c2565b9a87019a955050509084019060010161460e565b6000815180845260208085019450808401835b838110156145b7576146be878351805182526020810151602083015260408101516040830152606081015160608301525050565b608096909601959082019060010161468a565b8060005b60028110156146f45781518452602093840193909101906001016146d5565b50505050565b60008151808452614712816020860160208601614ef3565b601f01601f19169290920160200192915050565b6147318282516146d1565b6020810151604083015260408101516060830152606081015160808301526000608082015160c060a085015261476a60c08501826146fa565b949350505050565b60008151608084526147876080850182614677565b9050602083015184820360208601526147a08282614510565b915050604083015184820360408601526147ba8282614564565b915050606083015184820360608601526147d48282614510565b95945050505050565b8581526040856020830137606081018490526000828460808401379101608001908152949350505050565b88815261481860208201896146d1565b8660608201528560808201528460a08201528360c08201528260e08201526000610100835161484d8183860160208801614ef3565b929092019091019998505050505050505050565b84815261487160208201856146d1565b8260608201526000825161488c816080850160208701614ef3565b9190910160800195945050505050565b8481528360208201526080604082015282516080820152600060208401516148c760a08401826146d1565b50604084015160e083015260608401516101008301526080840151610120818185015260a086015161014085015260c086015161016085015260e08601519150806101808501525061491d6101a08401826146fa565b9050828103606084015261493181856145f1565b979650505050505050565b85815284602082015283604082015260a06060820152600061496160a08301856146fa565b828103608084015261497381856145f1565b98975050505050505050565b60208152600061309c60208301846146fa565b60208082526032908201527f70726576696f757320636f6d6d69746d656e7420697320616273656e7420696e604082015271040e8d0ca40c6eae4e4cadce840c4c2e8c6d60731b606082015260800190565b60208082526028908201527f54617267657420636f6d6d69746d656e7420697320616273656e7420696e20746040820152670d0ca40c4c2e8c6d60c31b606082015260800190565b60208082526027908201527f526f6c6c75703a20436f6d6d69746d656e74206e6f742070726573656e7420696040820152660dc40c4c2e8c6d60cb1b606082015260800190565b60208082526022908201527f70726576696f757320636f6d6d69746d656e74206861732077726f6e672070616040820152610e8d60f31b606082015260800190565b6020808252601d908201527f42617463684d616e616765723a20497320726f6c6c696e67206261636b000000604082015260600190565b60208082526017908201527f426174636820616c72656164792066696e616c69736564000000000000000000604082015260600190565b60208082526018908201527f526f6c6c75703a20496e76616c69642070726f706f7365720000000000000000604082015260600190565b60208082526033908201527f416c7265616479207375636365737366756c6c792064697370757465642e20526040820152726f6c6c206261636b20696e2070726f6365737360681b606082015260800190565b60208082526022908201527f6261746368494420646f6573206e6f74206d61746368206e6578744261746368604082015261125160f21b606082015260800190565b604081526000614c026040830185614726565b8281036020840152835160c08252614c1d60c0830182614677565b905060208501518282036020840152614c368282614510565b91505060408501518282036040840152614c508282614564565b91505060608501518282036060840152614c6a8282614510565b91505060808501518282036080840152614c8482826145c2565b91505060a085015182820360a08401526149318282614510565b604081526000614cb16040830185614726565b82810360208401526147d48185614772565b606081526000614cd66060830186614726565b8281036020840152614ce88186614772565b915050826040830152949350505050565b838152602081018390526060810160058310614d1757614d17614f7a565b826040830152949350505050565b82815260408101600c8310614d3c57614d3c614f7a565b8260208301529392505050565b6000808335601e19843603018112614d5f578283fd5b8301803591506001600160401b03821115614d78578283fd5b6020019150368190038213156137c757600080fd5b604080519081016001600160401b0381118282101715614daf57614daf614f90565b60405290565b604051606081016001600160401b0381118282101715614daf57614daf614f90565b60405161010081016001600160401b0381118282101715614daf57614daf614f90565b604051608081016001600160401b0381118282101715614daf57614daf614f90565b60405160c081016001600160401b0381118282101715614daf57614daf614f90565b604051601f8201601f191681016001600160401b0381118282101715614e6657614e66614f90565b604052919050565b60006001600160401b03821115614e8757614e87614f90565b5060051b60200190565b60008219821115614ea457614ea4614f4e565b500190565b600082614eb857614eb8614f64565b500490565b6000816000190483118215151615614ed757614ed7614f4e565b500290565b600082821015614eee57614eee614f4e565b500390565b60005b83811015614f0e578181015183820152602001614ef6565b838111156146f45750506000910152565b6000600019821415614f3357614f33614f4e565b5060010190565b600082614f4957614f49614f64565b500690565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052601260045260246000fd5b634e487b7160e01b600052602160045260246000fd5b634e487b7160e01b600052604160045260246000fdfe290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563a164736f6c6343000804000a290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563",
}

// RollupABI is the input ABI used to generate the binding from.
// Deprecated: Use RollupMetaData.ABI instead.
var RollupABI = RollupMetaData.ABI

// RollupBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use RollupMetaData.Bin instead.
var RollupBin = RollupMetaData.Bin

// DeployRollup deploys a new Ethereum contract, binding an instance of Rollup to it.
func DeployRollup(auth *bind.TransactOpts, backend bind.ContractBackend, _chooser common.Address, _depositManager common.Address, _accountRegistry common.Address, _transfer common.Address, _massMigration common.Address, _create2Transfer common.Address, genesisStateRoot [32]byte, stakeAmount *big.Int, blocksToFinalise *big.Int, minGasLeft *big.Int, maxTxsPerCommit *big.Int) (common.Address, *types.Transaction, *Rollup, error) {
	parsed, err := RollupMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(RollupBin), backend, _chooser, _depositManager, _accountRegistry, _transfer, _massMigration, _create2Transfer, genesisStateRoot, stakeAmount, blocksToFinalise, minGasLeft, maxTxsPerCommit)
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

// DOMAINNAME is a free data retrieval call binding the contract method 0x796f077b.
//
// Solidity: function DOMAIN_NAME() view returns(string)
func (_Rollup *RollupCaller) DOMAINNAME(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "DOMAIN_NAME")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// DOMAINNAME is a free data retrieval call binding the contract method 0x796f077b.
//
// Solidity: function DOMAIN_NAME() view returns(string)
func (_Rollup *RollupSession) DOMAINNAME() (string, error) {
	return _Rollup.Contract.DOMAINNAME(&_Rollup.CallOpts)
}

// DOMAINNAME is a free data retrieval call binding the contract method 0x796f077b.
//
// Solidity: function DOMAIN_NAME() view returns(string)
func (_Rollup *RollupCallerSession) DOMAINNAME() (string, error) {
	return _Rollup.Contract.DOMAINNAME(&_Rollup.CallOpts)
}

// DOMAINVERSION is a free data retrieval call binding the contract method 0xacb8cc49.
//
// Solidity: function DOMAIN_VERSION() view returns(string)
func (_Rollup *RollupCaller) DOMAINVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "DOMAIN_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// DOMAINVERSION is a free data retrieval call binding the contract method 0xacb8cc49.
//
// Solidity: function DOMAIN_VERSION() view returns(string)
func (_Rollup *RollupSession) DOMAINVERSION() (string, error) {
	return _Rollup.Contract.DOMAINVERSION(&_Rollup.CallOpts)
}

// DOMAINVERSION is a free data retrieval call binding the contract method 0xacb8cc49.
//
// Solidity: function DOMAIN_VERSION() view returns(string)
func (_Rollup *RollupCallerSession) DOMAINVERSION() (string, error) {
	return _Rollup.Contract.DOMAINVERSION(&_Rollup.CallOpts)
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

// DomainSeparator is a free data retrieval call binding the contract method 0xf698da25.
//
// Solidity: function domainSeparator() view returns(bytes32)
func (_Rollup *RollupCaller) DomainSeparator(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "domainSeparator")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DomainSeparator is a free data retrieval call binding the contract method 0xf698da25.
//
// Solidity: function domainSeparator() view returns(bytes32)
func (_Rollup *RollupSession) DomainSeparator() ([32]byte, error) {
	return _Rollup.Contract.DomainSeparator(&_Rollup.CallOpts)
}

// DomainSeparator is a free data retrieval call binding the contract method 0xf698da25.
//
// Solidity: function domainSeparator() view returns(bytes32)
func (_Rollup *RollupCallerSession) DomainSeparator() ([32]byte, error) {
	return _Rollup.Contract.DomainSeparator(&_Rollup.CallOpts)
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

// ZeroHashAtSubtreeDepth is a free data retrieval call binding the contract method 0xacf5b54d.
//
// Solidity: function zeroHashAtSubtreeDepth() view returns(bytes32)
func (_Rollup *RollupCaller) ZeroHashAtSubtreeDepth(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Rollup.contract.Call(opts, &out, "zeroHashAtSubtreeDepth")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ZeroHashAtSubtreeDepth is a free data retrieval call binding the contract method 0xacf5b54d.
//
// Solidity: function zeroHashAtSubtreeDepth() view returns(bytes32)
func (_Rollup *RollupSession) ZeroHashAtSubtreeDepth() ([32]byte, error) {
	return _Rollup.Contract.ZeroHashAtSubtreeDepth(&_Rollup.CallOpts)
}

// ZeroHashAtSubtreeDepth is a free data retrieval call binding the contract method 0xacf5b54d.
//
// Solidity: function zeroHashAtSubtreeDepth() view returns(bytes32)
func (_Rollup *RollupCallerSession) ZeroHashAtSubtreeDepth() ([32]byte, error) {
	return _Rollup.Contract.ZeroHashAtSubtreeDepth(&_Rollup.CallOpts)
}

// DisputeSignatureCreate2Transfer is a paid mutator transaction binding the contract method 0x4f6de740.
//
// Solidity: function disputeSignatureCreate2Transfer(uint256 batchID, ((bytes32,(bytes32,uint256[2],uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256)[],bytes32[][],uint256[4][],bytes32[][],bytes32[],bytes32[][]) signatureProof) returns()
func (_Rollup *RollupTransactor) DisputeSignatureCreate2Transfer(opts *bind.TransactOpts, batchID *big.Int, target TypesTransferCommitmentInclusionProof, signatureProof TypesSignatureProofWithReceiver) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "disputeSignatureCreate2Transfer", batchID, target, signatureProof)
}

// DisputeSignatureCreate2Transfer is a paid mutator transaction binding the contract method 0x4f6de740.
//
// Solidity: function disputeSignatureCreate2Transfer(uint256 batchID, ((bytes32,(bytes32,uint256[2],uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256)[],bytes32[][],uint256[4][],bytes32[][],bytes32[],bytes32[][]) signatureProof) returns()
func (_Rollup *RollupSession) DisputeSignatureCreate2Transfer(batchID *big.Int, target TypesTransferCommitmentInclusionProof, signatureProof TypesSignatureProofWithReceiver) (*types.Transaction, error) {
	return _Rollup.Contract.DisputeSignatureCreate2Transfer(&_Rollup.TransactOpts, batchID, target, signatureProof)
}

// DisputeSignatureCreate2Transfer is a paid mutator transaction binding the contract method 0x4f6de740.
//
// Solidity: function disputeSignatureCreate2Transfer(uint256 batchID, ((bytes32,(bytes32,uint256[2],uint256,bytes)),uint256,bytes32[]) target, ((uint256,uint256,uint256,uint256)[],bytes32[][],uint256[4][],bytes32[][],bytes32[],bytes32[][]) signatureProof) returns()
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

// SubmitCreate2Transfer is a paid mutator transaction binding the contract method 0x2e4518d8.
//
// Solidity: function submitCreate2Transfer(uint256 batchID, bytes32[] stateRoots, uint256[2][] signatures, uint256[] feeReceivers, bytes[] txss) payable returns()
func (_Rollup *RollupTransactor) SubmitCreate2Transfer(opts *bind.TransactOpts, batchID *big.Int, stateRoots [][32]byte, signatures [][2]*big.Int, feeReceivers []*big.Int, txss [][]byte) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "submitCreate2Transfer", batchID, stateRoots, signatures, feeReceivers, txss)
}

// SubmitCreate2Transfer is a paid mutator transaction binding the contract method 0x2e4518d8.
//
// Solidity: function submitCreate2Transfer(uint256 batchID, bytes32[] stateRoots, uint256[2][] signatures, uint256[] feeReceivers, bytes[] txss) payable returns()
func (_Rollup *RollupSession) SubmitCreate2Transfer(batchID *big.Int, stateRoots [][32]byte, signatures [][2]*big.Int, feeReceivers []*big.Int, txss [][]byte) (*types.Transaction, error) {
	return _Rollup.Contract.SubmitCreate2Transfer(&_Rollup.TransactOpts, batchID, stateRoots, signatures, feeReceivers, txss)
}

// SubmitCreate2Transfer is a paid mutator transaction binding the contract method 0x2e4518d8.
//
// Solidity: function submitCreate2Transfer(uint256 batchID, bytes32[] stateRoots, uint256[2][] signatures, uint256[] feeReceivers, bytes[] txss) payable returns()
func (_Rollup *RollupTransactorSession) SubmitCreate2Transfer(batchID *big.Int, stateRoots [][32]byte, signatures [][2]*big.Int, feeReceivers []*big.Int, txss [][]byte) (*types.Transaction, error) {
	return _Rollup.Contract.SubmitCreate2Transfer(&_Rollup.TransactOpts, batchID, stateRoots, signatures, feeReceivers, txss)
}

// SubmitDeposits is a paid mutator transaction binding the contract method 0xb6df3e6e.
//
// Solidity: function submitDeposits(uint256 batchID, ((bytes32,bytes32),uint256,bytes32[]) previous, (uint256,bytes32[]) vacant) payable returns()
func (_Rollup *RollupTransactor) SubmitDeposits(opts *bind.TransactOpts, batchID *big.Int, previous TypesCommitmentInclusionProof, vacant TypesSubtreeVacancyProof) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "submitDeposits", batchID, previous, vacant)
}

// SubmitDeposits is a paid mutator transaction binding the contract method 0xb6df3e6e.
//
// Solidity: function submitDeposits(uint256 batchID, ((bytes32,bytes32),uint256,bytes32[]) previous, (uint256,bytes32[]) vacant) payable returns()
func (_Rollup *RollupSession) SubmitDeposits(batchID *big.Int, previous TypesCommitmentInclusionProof, vacant TypesSubtreeVacancyProof) (*types.Transaction, error) {
	return _Rollup.Contract.SubmitDeposits(&_Rollup.TransactOpts, batchID, previous, vacant)
}

// SubmitDeposits is a paid mutator transaction binding the contract method 0xb6df3e6e.
//
// Solidity: function submitDeposits(uint256 batchID, ((bytes32,bytes32),uint256,bytes32[]) previous, (uint256,bytes32[]) vacant) payable returns()
func (_Rollup *RollupTransactorSession) SubmitDeposits(batchID *big.Int, previous TypesCommitmentInclusionProof, vacant TypesSubtreeVacancyProof) (*types.Transaction, error) {
	return _Rollup.Contract.SubmitDeposits(&_Rollup.TransactOpts, batchID, previous, vacant)
}

// SubmitMassMigration is a paid mutator transaction binding the contract method 0x0564e8b3.
//
// Solidity: function submitMassMigration(uint256 batchID, bytes32[] stateRoots, uint256[2][] signatures, uint256[4][] meta, bytes32[] withdrawRoots, bytes[] txss) payable returns()
func (_Rollup *RollupTransactor) SubmitMassMigration(opts *bind.TransactOpts, batchID *big.Int, stateRoots [][32]byte, signatures [][2]*big.Int, meta [][4]*big.Int, withdrawRoots [][32]byte, txss [][]byte) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "submitMassMigration", batchID, stateRoots, signatures, meta, withdrawRoots, txss)
}

// SubmitMassMigration is a paid mutator transaction binding the contract method 0x0564e8b3.
//
// Solidity: function submitMassMigration(uint256 batchID, bytes32[] stateRoots, uint256[2][] signatures, uint256[4][] meta, bytes32[] withdrawRoots, bytes[] txss) payable returns()
func (_Rollup *RollupSession) SubmitMassMigration(batchID *big.Int, stateRoots [][32]byte, signatures [][2]*big.Int, meta [][4]*big.Int, withdrawRoots [][32]byte, txss [][]byte) (*types.Transaction, error) {
	return _Rollup.Contract.SubmitMassMigration(&_Rollup.TransactOpts, batchID, stateRoots, signatures, meta, withdrawRoots, txss)
}

// SubmitMassMigration is a paid mutator transaction binding the contract method 0x0564e8b3.
//
// Solidity: function submitMassMigration(uint256 batchID, bytes32[] stateRoots, uint256[2][] signatures, uint256[4][] meta, bytes32[] withdrawRoots, bytes[] txss) payable returns()
func (_Rollup *RollupTransactorSession) SubmitMassMigration(batchID *big.Int, stateRoots [][32]byte, signatures [][2]*big.Int, meta [][4]*big.Int, withdrawRoots [][32]byte, txss [][]byte) (*types.Transaction, error) {
	return _Rollup.Contract.SubmitMassMigration(&_Rollup.TransactOpts, batchID, stateRoots, signatures, meta, withdrawRoots, txss)
}

// SubmitTransfer is a paid mutator transaction binding the contract method 0x39c13983.
//
// Solidity: function submitTransfer(uint256 batchID, bytes32[] stateRoots, uint256[2][] signatures, uint256[] feeReceivers, bytes[] txss) payable returns()
func (_Rollup *RollupTransactor) SubmitTransfer(opts *bind.TransactOpts, batchID *big.Int, stateRoots [][32]byte, signatures [][2]*big.Int, feeReceivers []*big.Int, txss [][]byte) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "submitTransfer", batchID, stateRoots, signatures, feeReceivers, txss)
}

// SubmitTransfer is a paid mutator transaction binding the contract method 0x39c13983.
//
// Solidity: function submitTransfer(uint256 batchID, bytes32[] stateRoots, uint256[2][] signatures, uint256[] feeReceivers, bytes[] txss) payable returns()
func (_Rollup *RollupSession) SubmitTransfer(batchID *big.Int, stateRoots [][32]byte, signatures [][2]*big.Int, feeReceivers []*big.Int, txss [][]byte) (*types.Transaction, error) {
	return _Rollup.Contract.SubmitTransfer(&_Rollup.TransactOpts, batchID, stateRoots, signatures, feeReceivers, txss)
}

// SubmitTransfer is a paid mutator transaction binding the contract method 0x39c13983.
//
// Solidity: function submitTransfer(uint256 batchID, bytes32[] stateRoots, uint256[2][] signatures, uint256[] feeReceivers, bytes[] txss) payable returns()
func (_Rollup *RollupTransactorSession) SubmitTransfer(batchID *big.Int, stateRoots [][32]byte, signatures [][2]*big.Int, feeReceivers []*big.Int, txss [][]byte) (*types.Transaction, error) {
	return _Rollup.Contract.SubmitTransfer(&_Rollup.TransactOpts, batchID, stateRoots, signatures, feeReceivers, txss)
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

// RollupRollbackTriggeredIterator is returned from FilterRollbackTriggered and is used to iterate over the raw logs and unpacked data for RollbackTriggered events raised by the Rollup contract.
type RollupRollbackTriggeredIterator struct {
	Event *RollupRollbackTriggered // Event containing the contract specifics and raw log

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
func (it *RollupRollbackTriggeredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RollupRollbackTriggered)
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
		it.Event = new(RollupRollbackTriggered)
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
func (it *RollupRollbackTriggeredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RollupRollbackTriggeredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RollupRollbackTriggered represents a RollbackTriggered event raised by the Rollup contract.
type RollupRollbackTriggered struct {
	BatchID *big.Int
	Result  uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRollbackTriggered is a free log retrieval operation binding the contract event 0x79e7023a0e95197fd8e57153afa9dca5bb6ccf11c7b8a99b391eb11ba75d0704.
//
// Solidity: event RollbackTriggered(uint256 batchID, uint8 result)
func (_Rollup *RollupFilterer) FilterRollbackTriggered(opts *bind.FilterOpts) (*RollupRollbackTriggeredIterator, error) {

	logs, sub, err := _Rollup.contract.FilterLogs(opts, "RollbackTriggered")
	if err != nil {
		return nil, err
	}
	return &RollupRollbackTriggeredIterator{contract: _Rollup.contract, event: "RollbackTriggered", logs: logs, sub: sub}, nil
}

// WatchRollbackTriggered is a free log subscription operation binding the contract event 0x79e7023a0e95197fd8e57153afa9dca5bb6ccf11c7b8a99b391eb11ba75d0704.
//
// Solidity: event RollbackTriggered(uint256 batchID, uint8 result)
func (_Rollup *RollupFilterer) WatchRollbackTriggered(opts *bind.WatchOpts, sink chan<- *RollupRollbackTriggered) (event.Subscription, error) {

	logs, sub, err := _Rollup.contract.WatchLogs(opts, "RollbackTriggered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RollupRollbackTriggered)
				if err := _Rollup.contract.UnpackLog(event, "RollbackTriggered", log); err != nil {
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

// ParseRollbackTriggered is a log parse operation binding the contract event 0x79e7023a0e95197fd8e57153afa9dca5bb6ccf11c7b8a99b391eb11ba75d0704.
//
// Solidity: event RollbackTriggered(uint256 batchID, uint8 result)
func (_Rollup *RollupFilterer) ParseRollbackTriggered(log types.Log) (*RollupRollbackTriggered, error) {
	event := new(RollupRollbackTriggered)
	if err := _Rollup.contract.UnpackLog(event, "RollbackTriggered", log); err != nil {
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
