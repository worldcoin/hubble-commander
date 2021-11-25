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
	Bin: "0x61028060405260006001553480156200001757600080fd5b5060405162005816380380620058168339810160408190526200003a91620008f0565b60405180604001604052806006815260200165487562626c6560d01b815250604051806040016040528060018152602001603160f81b8152508585858f8f84608081815250508360a081815250508260c08181525050816001600160a01b0316610100816001600160a01b031660601b81525050806001600160a01b031660e0816001600160a01b031660601b81525050505050505060008280519060200120905060008280519060200120905060007f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f90508261016081815250508161018081815250506200012f6200033c60201b60201c565b610140526200014081848462000340565b610120526101a0525050506001600160601b031960608b811b82166101c0528a811b82166101e05289811b82166102005288901b166102205250610240819052604080516338080fa560e21b8152905162000215916001600160a01b038d169163e0203e94916004808201926020929091908290030181600087803b158015620001c957600080fd5b505af1158015620001de573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002049190620009b8565b6200038560201b6200216d1760201c565b6102605260405160009062000240908790600080516020620057f683398151915290602001620009d1565b60408051601f1981840301815290829052805160209182012092506000916200027d918491600080516020620057f68339815191529101620009d1565b60408051601f1981840301815282825280516020918201208383019092528183529092508101620002be6000600133436200087260201b620026561760201c565b90526001805460009081526020818152604080832085518155949091015193830193909355905491517f3c367559be37298dc40fed468321ed4d44d99ab593f8e0bb7c82dbc84ab0bad69262000318929091819062000a0b565b60405180910390a1505060018054810190555062000a4b9950505050505050505050565b4690565b60008383836200034f6200033c565b3060405160200162000366959493929190620009df565b6040516020818303038152906040528051906020012090509392505050565b600062000391620008d1565b600080516020620057f683398151915281527f633dc4d7da7256660a892f8f1604a44b5432649cc8ec5cb3ced4c4e6ac94dd1d6020808301919091527f890740a8eb06ce9be422cb8da5cdafc2b58c0a5e24036c578de2a433c828ff7d60408301527f3b8ec09e026fdc305365dfc94e189a81b38c7597b3d941c279f042e8206e0bd860608301527fecd50eee38e386bd62be9bedb990706951b65fe053bd9d8a521af753d139e2da60808301527fdefff6d330bb5403f63b14f33b578274160de3a50df4efecf0e0db73bcdd3da560a08301527f617bdd11f7c0a11f49db22f629387a12da7596f9d1704d7465177c63d88ec7d760c08301527f292c23a9aa1d8bea7e2435e555a4a60e379a5a35f3f452bae60121073fb6eead60e08301527fe1cea92ed99acdcb045a6726b2f87107e8a61620a232cf4d7d5b5766b3952e106101008301527f7ad66c0a68c72cb89e4fb4303841966e4062a76ab97451e3b9fb526a5ceb7f826101208301527fe026cc5a4aed3c22a58cbd3d2ac754c9352c5436f638042dca99034e836365166101408301527f3d04cffd8b46a874edf5cfae63077de85f849a660426697b06a829c70dd1409c6101608301527fad676aa337a485e4728a0b240d92b3ef7b3c372d06d189322bfd5f61f1e7203e6101808301527fa2fca4a49658f9fab7aa63289c91b7c7b6c832a6d0e69334ff5b0a3483d09dab6101a08301527f4ebfd9cd7bca2505f7bef59cc1c12ecc708fff26ae4af19abe852afe9e20c8626101c08301527f2def10d13dd169f550f578bda343d9717a138562e0093b380a1120789d53cf106101e08301527f776a31db34a1a0a7caaf862cffdfff1789297ffadc380bd3d39281d340abd3ad6102008301527fe2e7610b87a5fdf3a72ebe271287d923ab990eefac64b6e59d79f8b7e08c46e36102208301527f504364a5c6858bf98fff714ab5be9de19ed31a976860efbd0e772a2efe23e2e06102408301527f4f05f4acb83f5b65168d9fef89d56d4d77b8944015e6b1eed81b0238e2d0dba36102608301527f44a6d974c75b07423e1d6d33f481916fdd45830aea11b6347e700cd8b9f0767c6102808301527fedf260291f734ddac396a956127dde4c34c0cfb8d8052f88ac139658ccf2d5076102a08301527f6075c657a105351e7f0fce53bc320113324a522e8fd52dc878c762551e01a46e6102c08301527f6ca6a3f763a9395f7da16014725ca7ee17e4815c0ff8119bf33f273dee11833b6102e08301527f1c25ef10ffeb3c7d08aa707d17286e0b0d3cbcb50f1bd3b6523b63ba3b52dd0f6103008301527ffffc43bd08273ccf135fd3cacbeef055418e09eb728d727c4d5d5c556cdea7e36103208301527fc5ab8111456b1f28f3c7a0a604b4553ce905cb019c463ee159137af83c350b226103408301527f0ff273fcbf4ae0f2bd88d6cf319ff4004f8d7dca70d4ced4e74d2c74139739e66103608301527f7fa06ba11241ddd5efdc65d4e39c9f6991b74fd4b81b62230808216c876f827c6103808301527f7e275adf313a996c7e2950cac67caba02a5ff925ebf9906b58949f3e77aec5b96103a08301527f8f6162fa308d2b3a15dc33cffac85f13ab349173121645aedf00f471663108be6103c08301527f78ccaaab73373552f207a63599de54d7d8d0c1805f86ce7da15818d09f4cff626103e08301528190849081106200086657fe5b60200201519392505050565b69ffffffff000000000000603082901b16600160501b600160f01b03605084901b1660ff60f01b60f086901b167fff0000000000000000000000000000000000000000000000000000000000000060f888901b16171717949350505050565b6040518061040001604052806020906020820280368337509192915050565b60008060008060008060008060008060006101608c8e03121562000912578687fd5b8b516200091f8162000a32565b60208d0151909b50620009328162000a32565b60408d0151909a50620009458162000a32565b60608d0151909950620009588162000a32565b60808d01519098506200096b8162000a32565b60a08d01519097506200097e8162000a32565b8096505060c08c0151945060e08c015193506101008c015192506101208c015191506101408c015190509295989b509295989b9093969950565b600060208284031215620009ca578081fd5b5051919050565b918252602082015260400190565b9485526020850193909352604084019190915260608301526001600160a01b0316608082015260a00190565b83815260208101839052606081016005831062000a2457fe5b826040830152949350505050565b6001600160a01b038116811462000a4857600080fd5b50565b60805160a05160c05160e05160601c6101005160601c610120516101405161016051610180516101a0516101c05160601c6101e05160601c6102005160601c6102205160601c6102405161026051614c7362000b83600039806119bb5280611bac52508061141e528061185e5280611f4752806121205250806116235280611826528061190052508061095b5280610a8452806120f15250806113e652806119245280611ea25250806106bb5280610e0552806111015280611eff525080612ed0525080612f12525080612ef1525080612e77525080612ea75250806105685280610cbb5280610fb752806119485280611a0a5250806118ba5280611bf25280612c2c525080611edb5280612bbe525080610f9352806129ee525080610c375280611f23528061298c5280612cfa5250614c736000f3fe6080604052600436106101ee5760003560e01c80638a4068dd1161010d578063b6df3e6e116100a0578063d1b243e21161006f578063d1b243e2146104e7578063e42f7554146104fc578063eb84e07614610511578063f698da2514610531578063fe5f960414610546576101ee565b8063b6df3e6e1461048a578063c08caaef1461049d578063ca858d9d146104bd578063d089e11a146104d2576101ee565b8063acb8cc49116100dc578063acb8cc4914610412578063acf5b54d14610427578063b02c43d01461043c578063b32c4d8d1461045c576101ee565b80638a4068dd146103be57806398d17621146103d3578063a72598ce146103e8578063ac96f0cd146103fd576101ee565b80634e23e8c3116101855780635f6e91d5116101545780635f6e91d5146103525780636c7ac9d814610372578063796f077b146103875780637ae8c568146103a9576101ee565b80634e23e8c3146102d05780634f6de740146102f05780635ac44282146103105780635b097d371461033d576101ee565b806325d5971f116101c157806325d5971f146102755780632e4518d81461029557806331c2b7db146102a857806339c13983146102bd576101ee565b80630564e8b3146101f3578063069321b0146102085780630ed75b9c146102335780632538507d14610255575b600080fd5b6102066102013660046139fa565b610566565b005b34801561021457600080fd5b5061021d610947565b60405161022a91906142f2565b60405180910390f35b34801561023f57600080fd5b50610248610959565b60405161022a919061440a565b34801561026157600080fd5b50610206610270366004613d6f565b61097d565b34801561028157600080fd5b506102066102903660046139e2565b610b55565b6102066102a3366004613af3565b610cb9565b3480156102b457600080fd5b5061021d610f91565b6102066102cb366004613af3565b610fb5565b3480156102dc57600080fd5b506102066102eb366004613d10565b61127e565b3480156102fc57600080fd5b5061020661030b366004613dd8565b611506565b34801561031c57600080fd5b5061033061032b3660046139e2565b61165a565b60405161022a9190614af8565b34801561034957600080fd5b5061021d6116b8565b34801561035e57600080fd5b5061020661036d366004613d10565b6116be565b34801561037e57600080fd5b506102486118b8565b34801561039357600080fd5b5061039c6118dc565b60405161022a919061441e565b3480156103b557600080fd5b506102486118fe565b3480156103ca57600080fd5b50610248611922565b3480156103df57600080fd5b50610248611946565b3480156103f457600080fd5b5061021d61196a565b34801561040957600080fd5b50610206611970565b34801561041e57600080fd5b5061039c61199c565b34801561043357600080fd5b5061021d6119b9565b34801561044857600080fd5b5061021d6104573660046139e2565b6119dd565b34801561046857600080fd5b5061047c6104773660046139e2565b6119ef565b60405161022a9291906142cb565b610206610498366004613c6d565b611a08565b3480156104a957600080fd5b506102066104b8366004613f20565b611d88565b3480156104c957600080fd5b5061021d611ed9565b3480156104de57600080fd5b50610248611efd565b3480156104f357600080fd5b5061021d611f21565b34801561050857600080fd5b5061021d611f45565b34801561051d57600080fd5b5061020661052c366004613be0565b611f69565b34801561053d57600080fd5b5061021d61214b565b34801561055257600080fd5b5061021d6105613660046139e2565b61215b565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b1580156105bf57600080fd5b505afa1580156105d3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105f79190613958565b6001600160a01b0316336001600160a01b0316146106305760405162461bcd60e51b815260040161062790614820565b60405180910390fd5b600454156106505760405162461bcd60e51b8152600401610627906146f7565b8a60015481146106725760405162461bcd60e51b81526004016106279061493e565b60608a6001600160401b038111801561068a57600080fd5b506040519080825280602002602001820160405280156106b4578160200160208202803683370190505b50905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663ebf0c7176040518163ffffffff1660e01b815260040160206040518083038186803b15801561071257600080fd5b505afa158015610726573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061074a919061397f565b905060005b8c8110156109215761075f613095565b6040518061010001604052808481526020018e8e8581811061077d57fe5b9050604002016002806020026040519081016040528092919082600260200280828437600092019190915250505081526020018c8c858181106107bc57fe5b9050608002016000600481106107ce57fe5b602002013581526020018a8a858181106107e457fe5b9050602002013581526020018c8c858181106107fc57fe5b90506080020160016004811061080e57fe5b602002013581526020018c8c8581811061082457fe5b90506080020160026004811061083657fe5b602002013581526020018c8c8581811061084c57fe5b90506080020160036004811061085e57fe5b6020020135815260200188888581811061087457fe5b90506020028101906108869190614b81565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250505091525090508e8e838181106108ca57fe5b905060200201356108da8261269d565b6040516020016108eb9291906142cb565b6040516020818303038152906040528051906020012084838151811061090d57fe5b60209081029190910101525060010161074f565b5061093761092e836126f4565b8d83600261298a565b5050505050505050505050505050565b600080516020614c4783398151915281565b7f000000000000000000000000000000000000000000000000000000000000000081565b60008381526020818152604091829020825180840190935280548352600101549082015283906109ac90612a80565b43106109ca5760405162461bcd60e51b81526004016106279061472e565b6004548110806109da5750600454155b6109f65760405162461bcd60e51b8152600401610627906148eb565b600084815260208190526040902054610a0f9084612a90565b610a2b5760405162461bcd60e51b815260040161062790614765565b610a336130e4565b6040805160a081018252855160209081015181015182528651518183015286510151519181019190915260608101610a6961214b565b815260200185600001516020015160e00151815250905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663be42c7c48386886000015160200151604001516040518463ffffffff1660e01b8152600401610ade93929190614ac2565b60206040518083038186803b158015610af657600080fd5b505afa158015610b0a573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b2e91906139c6565b9050600081600b811115610b3e57fe5b14610b4d57610b4d8682612abb565b505050505050565b600081815260208181526040918290208251808401909352805483526001015490820152610b8290612b06565b6001600160a01b0316336001600160a01b031614610bb25760405162461bcd60e51b81526004016106279061454e565b600081815260208181526040918290208251808401909352805483526001015490820152610bdf90612a80565b4311610bfd5760405162461bcd60e51b8152600401610627906146a6565b610c08816003612b1c565b15610c255760405162461bcd60e51b815260040161062790614980565b610c30816003612b43565b60405133907f000000000000000000000000000000000000000000000000000000000000000080156108fc02916000818181858888f19350505050158015610c7c573d6000803e3d6000fd5b507f1248d48e2de900a1010c7fce73506969ecec243600bfc08b641b158f26d857cd3382604051610cae9291906142d9565b60405180910390a150565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b158015610d1257600080fd5b505afa158015610d26573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d4a9190613958565b6001600160a01b0316336001600160a01b031614610d7a5760405162461bcd60e51b815260040161062790614820565b60045415610d9a5760405162461bcd60e51b8152600401610627906146f7565b886001548114610dbc5760405162461bcd60e51b81526004016106279061493e565b6060886001600160401b0381118015610dd457600080fd5b50604051908082528060200260200182016040528015610dfe578160200160208202803683370190505b50905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663ebf0c7176040518163ffffffff1660e01b815260040160206040518083038186803b158015610e5c57600080fd5b505afa158015610e70573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e94919061397f565b90506000805b8b811015610f6c57828b8b83818110610eaf57fe5b9050604002018a8a84818110610ec157fe5b90506020020135898985818110610ed457fe5b9050602002810190610ee69190614b81565b604051602001610efa95949392919061420c565b6040516020818303038152906040528051906020012091508c8c82818110610f1e57fe5b9050602002013582604051602001610f379291906142cb565b60405160208183030381529060405280519060200120848281518110610f5957fe5b6020908102919091010152600101610e9a565b50610f82610f79846126f4565b8c84600361298a565b50505050505050505050505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b15801561100e57600080fd5b505afa158015611022573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110469190613958565b6001600160a01b0316336001600160a01b0316146110765760405162461bcd60e51b815260040161062790614820565b600454156110965760405162461bcd60e51b8152600401610627906146f7565b8860015481146110b85760405162461bcd60e51b81526004016106279061493e565b6060886001600160401b03811180156110d057600080fd5b506040519080825280602002602001820160405280156110fa578160200160208202803683370190505b50905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663ebf0c7176040518163ffffffff1660e01b815260040160206040518083038186803b15801561115857600080fd5b505afa15801561116c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611190919061397f565b90506000805b8b81101561126857828b8b838181106111ab57fe5b9050604002018a8a848181106111bd57fe5b905060200201358989858181106111d057fe5b90506020028101906111e29190614b81565b6040516020016111f695949392919061420c565b6040516020818303038152906040528051906020012091508c8c8281811061121a57fe5b90506020020135826040516020016112339291906142cb565b6040516020818303038152906040528051906020012084828151811061125557fe5b6020908102919091010152600101611196565b50610f82611275846126f4565b8c84600161298a565b60008481526020818152604091829020825180840190935280548352600101549082015284906112ad90612a80565b43106112cb5760405162461bcd60e51b81526004016106279061472e565b6004548110806112db5750600454155b6112f75760405162461bcd60e51b8152600401610627906148eb565b84848460200151600080826000141561134b5750600019840160008181526020818152604091829020825180840190935280548352600190810154918301919091529061134390612b6b565b039150611354565b50506000198101835b818460200151146113775760405162461bcd60e51b815260040161062790614664565b6000818152602081905260409020546113909085612b78565b6113ac5760405162461bcd60e51b815260040161062790614472565b60008a8152602081905260409020546113c59089612b8b565b6113e15760405162461bcd60e51b8152600401610627906145d5565b6000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663929314928c60000151600001517f00000000000000000000000000000000000000000000000000000000000000008d6000015160200151604001518e6000015160200151606001518e6040518663ffffffff1660e01b81526004016114789594939291906143c7565b604080518083038186803b15801561148f57600080fd5b505afa1580156114a3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114c79190613997565b9092509050600081600b8111156114da57fe5b1415806114e957508951518214155b156114f8576114f88c82612abb565b505050505050505050505050565b600083815260208181526040918290208251808401909352805483526001015490820152839061153590612a80565b43106115535760405162461bcd60e51b81526004016106279061472e565b6004548110806115635750600454155b61157f5760405162461bcd60e51b8152600401610627906148eb565b6000848152602081905260409020546115989084612b8b565b6115b45760405162461bcd60e51b81526004016106279061461d565b6115bc6130e4565b6040805160a0810182528551602090810151810151825286515181830152865101515191810191909152606081016115f261214b565b81528551602090810151606001519101526040516001627e8b6b60e01b031981529091506000906001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063ff81749590610ade90859088906004016149ee565b61166261311b565b6000828152602081905260409020600101546116905760405162461bcd60e51b815260040161062790614857565b5060009081526020818152604091829020825180840190935280548352600101549082015290565b60045481565b60008481526020818152604091829020825180840190935280548352600101549082015284906116ed90612a80565b431061170b5760405162461bcd60e51b81526004016106279061472e565b60045481108061171b5750600454155b6117375760405162461bcd60e51b8152600401610627906148eb565b84848460200151600080826000141561178b5750600019840160008181526020818152604091829020825180840190935280548352600190810154918301919091529061178390612b6b565b039150611794565b50506000198101835b818460200151146117b75760405162461bcd60e51b815260040161062790614664565b6000818152602081905260409020546117d09085612b78565b6117ec5760405162461bcd60e51b815260040161062790614472565b60008a8152602081905260409020546118059089612b8b565b6118215760405162461bcd60e51b8152600401610627906145d5565b6000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663336920368c60000151600001517f00000000000000000000000000000000000000000000000000000000000000008d6000015160200151604001518e6000015160200151606001518e6040518663ffffffff1660e01b81526004016114789594939291906143c7565b7f000000000000000000000000000000000000000000000000000000000000000081565b60405180604001604052806006815260200165487562626c6560d01b81525081565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b60015481565b6000600454116119925760405162461bcd60e51b815260040161062790614431565b61199a612b9e565b565b604051806040016040528060018152602001603160f81b81525081565b7f000000000000000000000000000000000000000000000000000000000000000081565b60026020526000908152604090205481565b6000602081905290815260409020805460019091015482565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b158015611a6157600080fd5b505afa158015611a75573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611a999190613958565b6001600160a01b0316336001600160a01b031614611ac95760405162461bcd60e51b815260040161062790614820565b60045415611ae95760405162461bcd60e51b8152600401610627906146f7565b826001548114611b0b5760405162461bcd60e51b81526004016106279061493e565b6000198401600081815260208181526040918290208251808401909352805483526001908101549183019190915290611b4390612b6b565b03846020015114611b665760405162461bcd60e51b815260040161062790614664565b600081815260208190526040902054611b7f9085612b78565b611b9b5760405162461bcd60e51b8152600401610627906144c4565b83515183516020850151611bd192917f000000000000000000000000000000000000000000000000000000000000000091612db0565b611bed5760405162461bcd60e51b8152600401610627906147dd565b6000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663d86ee48d6040518163ffffffff1660e01b81526004016040805180830381600087803b158015611c4a57600080fd5b505af1158015611c5e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611c829190613bbd565b600185016000818152600260205260409081902083905588519051939550919350917f1e6cd0ae9aa6f6e6ea5cf7f07396a96fc1ea5ae4ec7d5054f53e204c47f79bd891611cd39186918691614b35565b60405180910390a16000611cf08388600001518960200151612dc8565b9050600081600080516020614c4783398151915260001b604051602001611d189291906142cb565b60408051601f198184030181529082905280516020918201209250600091611d52918491600080516020614c4783398151915291016142cb565b604051602081830303815290604052805190602001209050611d7b8160016000801b600461298a565b5050505050505050505050565b6000838152602081815260409182902082518084019093528054835260010154908201528390611db790612a80565b4310611dd55760405162461bcd60e51b81526004016106279061472e565b600454811080611de55750600454155b611e015760405162461bcd60e51b8152600401610627906148eb565b600084815260208190526040902054611e1a9084612b8b565b611e365760405162461bcd60e51b81526004016106279061461d565b611e3e6130e4565b6040805160a081018252855160209081015181015182528651518183015286510151519181019190915260608101611e7461214b565b81528551602090810151606001519101526040516344ec5a7760e01b81529091506000906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906344ec5a7790610ade9085908890600401614a9d565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000848152602081815260409182902082518084019093528054835260010154908201528490611f9890612a80565b4310611fb65760405162461bcd60e51b81526004016106279061472e565b600454811080611fc65750600454155b611fe25760405162461bcd60e51b8152600401610627906148eb565b8484846020015160008082600014156120365750600019840160008181526020818152604091829020825180840190935280548352600190810154918301919091529061202e90612b6b565b03915061203f565b50506000198101835b818460200151146120625760405162461bcd60e51b815260040161062790614664565b60008181526020819052604090205461207b9085612b78565b6120975760405162461bcd60e51b815260040161062790614472565b60008a8152602081905260409020546120b09089612a90565b6120cc5760405162461bcd60e51b8152600401610627906145d5565b88515188516020015160405163ab5a164f60e01b815260009283926001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169263ab5a164f926114789290917f000000000000000000000000000000000000000000000000000000000000000091908f90600401614327565b6000612155612e73565b90505b90565b60036020526000908152604090205481565b6000612177613132565b600080516020614c4783398151915281527f633dc4d7da7256660a892f8f1604a44b5432649cc8ec5cb3ced4c4e6ac94dd1d6020808301919091527f890740a8eb06ce9be422cb8da5cdafc2b58c0a5e24036c578de2a433c828ff7d60408301527f3b8ec09e026fdc305365dfc94e189a81b38c7597b3d941c279f042e8206e0bd860608301527fecd50eee38e386bd62be9bedb990706951b65fe053bd9d8a521af753d139e2da60808301527fdefff6d330bb5403f63b14f33b578274160de3a50df4efecf0e0db73bcdd3da560a08301527f617bdd11f7c0a11f49db22f629387a12da7596f9d1704d7465177c63d88ec7d760c08301527f292c23a9aa1d8bea7e2435e555a4a60e379a5a35f3f452bae60121073fb6eead60e08301527fe1cea92ed99acdcb045a6726b2f87107e8a61620a232cf4d7d5b5766b3952e106101008301527f7ad66c0a68c72cb89e4fb4303841966e4062a76ab97451e3b9fb526a5ceb7f826101208301527fe026cc5a4aed3c22a58cbd3d2ac754c9352c5436f638042dca99034e836365166101408301527f3d04cffd8b46a874edf5cfae63077de85f849a660426697b06a829c70dd1409c6101608301527fad676aa337a485e4728a0b240d92b3ef7b3c372d06d189322bfd5f61f1e7203e6101808301527fa2fca4a49658f9fab7aa63289c91b7c7b6c832a6d0e69334ff5b0a3483d09dab6101a08301527f4ebfd9cd7bca2505f7bef59cc1c12ecc708fff26ae4af19abe852afe9e20c8626101c08301527f2def10d13dd169f550f578bda343d9717a138562e0093b380a1120789d53cf106101e08301527f776a31db34a1a0a7caaf862cffdfff1789297ffadc380bd3d39281d340abd3ad6102008301527fe2e7610b87a5fdf3a72ebe271287d923ab990eefac64b6e59d79f8b7e08c46e36102208301527f504364a5c6858bf98fff714ab5be9de19ed31a976860efbd0e772a2efe23e2e06102408301527f4f05f4acb83f5b65168d9fef89d56d4d77b8944015e6b1eed81b0238e2d0dba36102608301527f44a6d974c75b07423e1d6d33f481916fdd45830aea11b6347e700cd8b9f0767c6102808301527fedf260291f734ddac396a956127dde4c34c0cfb8d8052f88ac139658ccf2d5076102a08301527f6075c657a105351e7f0fce53bc320113324a522e8fd52dc878c762551e01a46e6102c08301527f6ca6a3f763a9395f7da16014725ca7ee17e4815c0ff8119bf33f273dee11833b6102e08301527f1c25ef10ffeb3c7d08aa707d17286e0b0d3cbcb50f1bd3b6523b63ba3b52dd0f6103008301527ffffc43bd08273ccf135fd3cacbeef055418e09eb728d727c4d5d5c556cdea7e36103208301527fc5ab8111456b1f28f3c7a0a604b4553ce905cb019c463ee159137af83c350b226103408301527f0ff273fcbf4ae0f2bd88d6cf319ff4004f8d7dca70d4ced4e74d2c74139739e66103608301527f7fa06ba11241ddd5efdc65d4e39c9f6991b74fd4b81b62230808216c876f827c6103808301527f7e275adf313a996c7e2950cac67caba02a5ff925ebf9906b58949f3e77aec5b96103a08301527f8f6162fa308d2b3a15dc33cffac85f13ab349173121645aedf00f471663108be6103c08301527f78ccaaab73373552f207a63599de54d7d8d0c1805f86ce7da15818d09f4cff626103e083015281908490811061264a57fe5b60200201519392505050565b69ffffffff000000000000603082901b16600160501b600160f01b03605084901b1660ff60f01b60f086901b166001600160f81b031960f888901b16171717949350505050565b80516020808301516040808501516060860151608087015160a088015160c089015160e08a015195516000996126d7999098979101614237565b604051602081830303815290604052805190602001209050919050565b60006020825111156127185760405162461bcd60e51b8152600401610627906148b4565b612720613151565b600080516020614c4783398151915281527f633dc4d7da7256660a892f8f1604a44b5432649cc8ec5cb3ced4c4e6ac94dd1d60208201527f890740a8eb06ce9be422cb8da5cdafc2b58c0a5e24036c578de2a433c828ff7d60408201527f3b8ec09e026fdc305365dfc94e189a81b38c7597b3d941c279f042e8206e0bd86060808301919091527fecd50eee38e386bd62be9bedb990706951b65fe053bd9d8a521af753d139e2da60808301527fdefff6d330bb5403f63b14f33b578274160de3a50df4efecf0e0db73bcdd3da560a08301528351600181811692918101901c90816001600160401b038111801561281757600080fd5b50604051908082528060200260200182016040528015612841578160200160208202803683370190505b5090508560005b60005b8585038110156128d0576000600182901b905083818151811061286a57fe5b602002602001015184826001018151811061288157fe5b602002602001015160405160200161289a9291906142cb565b604051602081830303815290604052805190602001208583815181106128bc57fe5b60209081029190910101525060010161284b565b85600114156129425782600182901b815181106128e957fe5b60200260200101518783600681106128fd57fe5b60200201516040516020016129139291906142cb565b6040516020818303038152906040528051906020012084828151811061293557fe5b6020026020010181815250505b84600114156129515750612968565b5060018085169550938401841c9383925001612848565b8260008151811061297557fe5b60200260200101519650505050505050919050565b7f00000000000000000000000000000000000000000000000000000000000000003410156129ca5760405162461bcd60e51b8152600401610627906149b7565b6040518060400160405280858152602001612a148360048111156129ea57fe5b86337f00000000000000000000000000000000000000000000000000000000000000004301612656565b9052600180546000908152602081815260409182902084518155930151928201929092555490517f3c367559be37298dc40fed468321ed4d44d99ab593f8e0bb7c82dbc84ab0bad691612a6a9185908590614b0f565b60405180910390a1505060018054810190555050565b6020015160301c63ffffffff1690565b6000612ab283612aa38460000151612f3d565b84602001518560400151612db0565b90505b92915050565b60048290556040517f79e7023a0e95197fd8e57153afa9dca5bb6ccf11c7b8a99b391eb11ba75d070490612af29084908490614b4b565b60405180910390a1612b02612b9e565b5050565b602081015160501c6001600160a01b0316919050565b610100820460009081526020919091526040902054600160ff9092169190911b9081161490565b61010082046000908152602091909152604090208054600160ff9093169290921b9091179055565b6020015160f01c60ff1690565b6000612ab283612aa38460000151612f62565b6000612ab283612aa38460000151612f7f565b6004546001546000918291908103908290600019015b82821015612c9f577f00000000000000000000000000000000000000000000000000000000000000005a11612be857612c9f565b8181036000818152602081815260408083208381556001018390556002909152902054955093508415612c94576040516356f0001360e11b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063ade0002690612c619088906004016142f2565b600060405180830381600087803b158015612c7b57600080fd5b505af1158015612c8f573d6000803e3d6000fd5b505050505b600190910190612bb4565b6001805483900390558282148015612cb75760006004555b7f595cb200f38fcdbf73429ffce272e53b2e923c231b2059ce20d1f4909cff1bec828483604051612cea93929190614b69565b60405180910390a16000612d1e847f0000000000000000000000000000000000000000000000000000000000000000612f93565b90506000612d386003612d32846002612f93565b90612fcd565b90506000612d468383612fff565b604051909150339083156108fc029084906000818181858888f19350505050158015612d76573d6000803e3d6000fd5b5060405160009082156108fc0290839083818181858288f19350505050158015612da4573d6000803e3d6000fd5b50505050505050505050565b600084612dbe858585612dc8565b1495945050505050565b600083815b8351811015612e6a57600185821c16612e235781848281518110612ded57fe5b6020026020010151604051602001612e069291906142cb565b604051602081830303815290604052805190602001209150612e62565b838181518110612e2f57fe5b602002602001015182604051602001612e499291906142cb565b6040516020818303038152906040528051906020012091505b600101612dcd565b50949350505050565b60007f0000000000000000000000000000000000000000000000000000000000000000612e9e613027565b1415612ecb57507f0000000000000000000000000000000000000000000000000000000000000000612158565b612f367f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f000000000000000000000000000000000000000000000000000000000000000061302b565b9050612158565b60008160000151612f51836020015161269d565b6040516020016126d79291906142cb565b6000816000015182602001516040516020016126d79291906142cb565b60008160000151612f51836020015161306c565b600082612fa257506000612ab5565b82820282848281612faf57fe5b0414612ab25760405162461bcd60e51b81526004016106279061479c565b6000808211612fee5760405162461bcd60e51b81526004016106279061459e565b818381612ff757fe5b049392505050565b6000828211156130215760405162461bcd60e51b815260040161062790614517565b50900390565b4690565b6000838383613038613027565b3060405160200161304d9594939291906142fb565b6040516020818303038152906040528051906020012090509392505050565b600081600001518260200151836040015184606001516040516020016126d79493929190614290565b60408051610100810190915260008152602081016130b161316f565b81526020016000815260200160008019168152602001600081526020016000815260200160008152602001606081525090565b6040518060a001604052806130f761316f565b81526000602082018190526040820181905260608083019190915260809091015290565b604080518082019091526000808252602082015290565b6040518061040001604052806020906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b600082601f83011261319d578081fd5b81356131b06131ab82614beb565b614bc5565b818152915060208083019084810160005b848110156131ea576131d8888484358a0101613378565b845292820192908201906001016131c1565b505050505092915050565b60008083601f840112613206578081fd5b5081356001600160401b0381111561321c578182fd5b60208301915083602060408302850101111561323757600080fd5b9250929050565b60008083601f84011261324f578182fd5b5081356001600160401b03811115613265578182fd5b60208301915083602060808302850101111561323757600080fd5b6000601f8381840112613291578182fd5b823561329f6131ab82614beb565b81815292506020808401908581016080808502880183018910156132c257600080fd5b60005b8581101561332a5789878401126132db57600080fd5b6132e482614bc5565b80848486018d8111156132f657600080fd5b60005b6004811015613316578235855293880193918801916001016132f9565b5091885250958501959350506001016132c5565b5050505050505092915050565b60008083601f840112613348578182fd5b5081356001600160401b0381111561335e578182fd5b602083019150836020808302850101111561323757600080fd5b600082601f830112613388578081fd5b81356133966131ab82614beb565b8181529150602080830190848101818402860182018710156133b757600080fd5b60005b848110156131ea578135845292820192908201906001016133ba565b600082601f8301126133e6578081fd5b81356133f46131ab82614beb565b818152915060208083019084810160005b848110156131ea578135870160a080601f19838c0301121561342657600080fd5b6134306040614bc5565b61343c8b878501613913565b815290820135906001600160401b0382111561345757600080fd5b6134658b8784860101613378565b81870152865250509282019290820190600101613405565b600082601f83011261348d578081fd5b813561349b6131ab82614beb565b81815291506020808301908481016080808502870183018810156134be57600080fd5b60005b858110156134e5576134d38984613913565b855293830193918101916001016134c1565b50505050505092915050565b600082601f830112613501578081fd5b61350b6040614bc5565b905080828460408501111561351f57600080fd5b60005b6002811015613541578135835260209283019290910190600101613522565b50505092915050565b600082601f83011261355a578081fd5b81356001600160401b0381111561356f578182fd5b613582601f8201601f1916602001614bc5565b915080825283602082850101111561359957600080fd5b8060208401602084013760009082016020015292915050565b600081830360808112156135c4578182fd5b6135ce6060614bc5565b915060408112156135de57600080fd5b506135e96040614bc5565b8235815260208301356020820152808252506040820135602082015260608201356001600160401b0381111561361e57600080fd5b61362a84828501613378565b60408301525092915050565b600060608284031215613647578081fd5b6136516060614bc5565b905081356001600160401b038082111561366a57600080fd5b908301906040828603121561367e57600080fd5b6136886040614bc5565b823581526020808401358381111561369f57600080fd5b939093019261012084880312156136b557600080fd5b6101006136c181614bc5565b853581526136d1898488016134f1565b83820152606086013560408201526080860135606082015260a0860135608082015260c086013560a082015260e086013560c08201528186013591508482111561371a57600080fd5b6137268983880161354a565b60e082015283830152509084528481013590840152604084013591508082111561374f57600080fd5b5061362a84828501613378565b60006080828403121561376d578081fd5b6137776080614bc5565b905081356001600160401b038082111561379057600080fd5b61379c8583860161347d565b835260208401359150808211156137b257600080fd5b6137be8583860161318d565b602084015260408401359150808211156137d757600080fd5b6137e385838601613280565b604084015260608401359150808211156137fc57600080fd5b506138098482850161318d565b60608301525092915050565b600060a08284031215613826578081fd5b6138306080614bc5565b90508135815261384383602084016134f1565b60208201526060820135604082015260808201356001600160401b0381111561386b57600080fd5b6138098482850161354a565b600060608284031215613888578081fd5b6138926060614bc5565b905081356001600160401b03808211156138ab57600080fd5b90830190604082860312156138bf57600080fd5b6138c96040614bc5565b823581526020830135828111156138df57600080fd5b6138eb87828601613815565b6020830152508084525060208401356020840152604084013591508082111561374f57600080fd5b600060808284031215613924578081fd5b61392e6080614bc5565b90508135815260208201356020820152604082013560408201526060820135606082015292915050565b600060208284031215613969578081fd5b81516001600160a01b0381168114612ab2578182fd5b600060208284031215613990578081fd5b5051919050565b600080604083850312156139a9578081fd5b8251915060208301516139bb81614c36565b809150509250929050565b6000602082840312156139d7578081fd5b8151612ab281614c36565b6000602082840312156139f3578081fd5b5035919050565b600080600080600080600080600080600060c08c8e031215613a1a578687fd5b8b359a506001600160401b038060208e01351115613a36578788fd5b613a468e60208f01358f01613337565b909b50995060408d0135811015613a5b578788fd5b613a6b8e60408f01358f016131f5565b909950975060608d0135811015613a80578687fd5b613a908e60608f01358f0161323e565b909750955060808d0135811015613aa5578485fd5b613ab58e60808f01358f01613337565b909550935060a08d0135811015613aca578283fd5b50613adb8d60a08e01358e01613337565b81935080925050509295989b509295989b9093969950565b600080600080600080600080600060a08a8c031215613b10578283fd5b8935985060208a01356001600160401b0380821115613b2d578485fd5b613b398d838e01613337565b909a50985060408c0135915080821115613b51578485fd5b613b5d8d838e016131f5565b909850965060608c0135915080821115613b75578485fd5b613b818d838e01613337565b909650945060808c0135915080821115613b99578384fd5b50613ba68c828d01613337565b915080935050809150509295985092959850929598565b60008060408385031215613bcf578182fd5b505080516020909101519092909150565b60008060008060808587031215613bf5578182fd5b8435935060208501356001600160401b0380821115613c12578384fd5b613c1e888389016135b2565b94506040870135915080821115613c33578384fd5b613c3f88838901613636565b93506060870135915080821115613c54578283fd5b50613c61878288016133d6565b91505092959194509250565b600080600060608486031215613c81578081fd5b8335925060208401356001600160401b0380821115613c9e578283fd5b613caa878388016135b2565b93506040860135915080821115613cbf578283fd5b9085019060408288031215613cd2578283fd5b613cdc6040614bc5565b82358152602083013582811115613cf1578485fd5b613cfd89828601613378565b6020830152508093505050509250925092565b60008060008060808587031215613d25578182fd5b8435935060208501356001600160401b0380821115613d42578384fd5b613d4e888389016135b2565b94506040870135915080821115613d63578384fd5b613c3f88838901613877565b600080600060608486031215613d83578081fd5b8335925060208401356001600160401b0380821115613da0578283fd5b613dac87838801613636565b93506040860135915080821115613dc1578283fd5b50613dce8682870161375c565b9150509250925092565b600080600060608486031215613dec578081fd5b8335925060208401356001600160401b0380821115613e09578283fd5b613e1587838801613877565b93506040860135915080821115613e2a578283fd5b9085019060c08288031215613e3d578283fd5b613e4760c0614bc5565b823582811115613e55578485fd5b613e618982860161347d565b825250602083013582811115613e75578485fd5b613e818982860161318d565b602083015250604083013582811115613e98578485fd5b613ea489828601613280565b604083015250606083013582811115613ebb578485fd5b613ec78982860161318d565b606083015250608083013582811115613ede578485fd5b613eea89828601613378565b60808301525060a083013582811115613f01578485fd5b613f0d8982860161318d565b60a0830152508093505050509250925092565b600080600060608486031215613f34578081fd5b8335925060208401356001600160401b0380821115613f51578283fd5b613dac87838801613877565b6000815180845260208085018081965082840281019150828601855b85811015613fa3578284038952613f9184835161400e565b98850198935090840190600101613f79565b5091979650505050505050565b6000815180845260208085019450808401835b8381101561400357815187865b6004811015613fed57825182529185019190850190600101613fd0565b5050506080969096019590820190600101613fc3565b509495945050505050565b6000815180845260208085019450808401835b8381101561400357815187529582019590820190600101614021565b6000815180845260208085018081965082840281019150828601855b85811015613fa3578284038952815160a06140758683516141e6565b86820151915080608087015261408d8187018361400e565b9a87019a9550505090840190600101614059565b6000815180845260208085019450808401835b83811015614003576140c78783516141e6565b60809690960195908201906001016140b4565b8060005b60028110156140fd5781518452602093840193909101906001016140de565b50505050565b6000815180845261411b816020860160208601614c0a565b601f01601f19169290920160200192915050565b600061413c8383516140da565b602082015160408401526040820151606084015260608201516080840152608082015160c060a085015261417360c0850182614103565b949350505050565b600081516080845261419060808501826140a1565b9050602083015184820360208601526141a98282613f5d565b915050604083015184820360408601526141c38282613fb0565b915050606083015184820360608601526141dd8282613f5d565b95945050505050565b805182526020810151602083015260408101516040830152606081015160608301525050565b6000868252604086602084013760608201859052828460808401379101608001908152949350505050565b6000898252614249602083018a6140da565b8760608301528660808301528560a08301528460c08301528360e0830152610100835161427c8183860160208801614c0a565b929092019091019998505050505050505050565b60008582526142a260208301866140da565b83606083015282516142bb816080850160208701614c0a565b9190910160800195945050505050565b918252602082015260400190565b6001600160a01b03929092168252602082015260400190565b90815260200190565b9485526020850193909352604084019190915260608301526001600160a01b0316608082015260a00190565b60008582528460208301526080604083015283516080830152602084015161435260a08401826140da565b50604084015160e083015260608401516101008301526080840151610120818185015260a086015161014085015260c086015161016085015260e0860151915080610180850152506143a86101a0840182614103565b905082810360608401526143bc818561403d565b979650505050505050565b600086825285602083015284604083015260a060608301526143ec60a0830185614103565b82810360808401526143fe818561403d565b98975050505050505050565b6001600160a01b0391909116815260200190565b600060208252612ab26020830184614103565b60208082526021908201527f42617463684d616e616765723a204973206e6f7420726f6c6c696e67206261636040820152606b60f81b606082015260800190565b60208082526032908201527f70726576696f757320636f6d6d69746d656e7420697320616273656e7420696e604082015271040e8d0ca40c6eae4e4cadce840c4c2e8c6d60731b606082015260800190565b60208082526033908201527f70726576696f757320636f6d6d69746d656e7420697320616273656e7420696e604082015272040e8d0ca40e0e4caecd2deeae640c4c2e8c6d606b1b606082015260800190565b6020808252601e908201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604082015260600190565b60208082526030908201527f596f7520617265206e6f742074686520636f727265637420636f6d6d6974746560408201526f0e440ccdee440e8d0d2e640c4c2e8c6d60831b606082015260800190565b6020808252601a908201527f536166654d6174683a206469766973696f6e206279207a65726f000000000000604082015260600190565b60208082526028908201527f54617267657420636f6d6d69746d656e7420697320616273656e7420696e20746040820152670d0ca40c4c2e8c6d60c31b606082015260800190565b60208082526027908201527f526f6c6c75703a20436f6d6d69746d656e74206e6f742070726573656e7420696040820152660dc40c4c2e8c6d60cb1b606082015260800190565b60208082526022908201527f70726576696f757320636f6d6d69746d656e74206861732077726f6e672070616040820152610e8d60f31b606082015260800190565b60208082526031908201527f54686973206261746368206973206e6f74207965742066696e616c697365642c60408201527020636865636b206261636b20736f6f6e2160781b606082015260800190565b6020808252601d908201527f42617463684d616e616765723a20497320726f6c6c696e67206261636b000000604082015260600190565b60208082526017908201527f426174636820616c72656164792066696e616c69736564000000000000000000604082015260600190565b6020808252601f908201527f436f6d6d69746d656e74206e6f742070726573656e7420696e20626174636800604082015260600190565b60208082526021908201527f536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f6040820152607760f81b606082015260800190565b60208082526023908201527f526f6c6c75703a2053746174652073756274726565206973206e6f7420766163604082015262185b9d60ea1b606082015260800190565b60208082526018908201527f526f6c6c75703a20496e76616c69642070726f706f7365720000000000000000604082015260600190565b6020808252603f908201527f42617463682069642067726561746572207468616e20746f74616c206e756d6260408201527f6572206f6620626174636865732c20696e76616c696420626174636820696400606082015260800190565b6020808252601b908201527f4d65726b6c65547265653a20546f6f206d616e79206c65617665730000000000604082015260600190565b60208082526033908201527f416c7265616479207375636365737366756c6c792064697370757465642e20526040820152726f6c6c206261636b20696e2070726f6365737360681b606082015260800190565b60208082526022908201527f6261746368494420646f6573206e6f74206d61746368206e6578744261746368604082015261125160f21b606082015260800190565b60208082526019908201527f526f6c6c75703a20416c72656164792077697468647261776e00000000000000604082015260600190565b6020808252601a908201527f526f6c6c75703a2077726f6e67207374616b6520616d6f756e74000000000000604082015260600190565b600060408252614a01604083018561412f565b8281036020840152835160c08252614a1c60c08301826140a1565b905060208501518282036020840152614a358282613f5d565b91505060408501518282036040840152614a4f8282613fb0565b91505060608501518282036060840152614a698282613f5d565b91505060808501518282036080840152614a83828261400e565b91505060a085015182820360a08401526143bc8282613f5d565b600060408252614ab0604083018561412f565b82810360208401526141dd818561417b565b600060608252614ad5606083018661412f565b8281036020840152614ae7818661417b565b915050826040830152949350505050565b815181526020918201519181019190915260400190565b838152602081018390526060810160058310614b2757fe5b826040830152949350505050565b9283526020830191909152604082015260600190565b82815260408101600c8310614b5c57fe5b8260208301529392505050565b92835260208301919091521515604082015260600190565b6000808335601e19843603018112614b97578283fd5b8301803591506001600160401b03821115614bb0578283fd5b60200191503681900382131561323757600080fd5b6040518181016001600160401b0381118282101715614be357600080fd5b604052919050565b60006001600160401b03821115614c00578081fd5b5060209081020190565b60005b83811015614c25578181015183820152602001614c0d565b838111156140fd5750506000910152565b600c8110614c4357600080fd5b5056fe290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563a164736f6c634300060c000a290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563",
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
