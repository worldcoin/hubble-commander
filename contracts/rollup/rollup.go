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
	ABI: "[{\"inputs\":[{\"internalType\":\"contractChooser\",\"name\":\"_chooser\",\"type\":\"address\"},{\"internalType\":\"contractIDepositManager\",\"name\":\"_depositManager\",\"type\":\"address\"},{\"internalType\":\"contractBLSAccountRegistry\",\"name\":\"_accountRegistry\",\"type\":\"address\"},{\"internalType\":\"contractTransfer\",\"name\":\"_transfer\",\"type\":\"address\"},{\"internalType\":\"contractMassMigration\",\"name\":\"_massMigration\",\"type\":\"address\"},{\"internalType\":\"contractCreate2Transfer\",\"name\":\"_create2Transfer\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"genesisStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"stakeAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blocksToFinalise\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"minGasLeft\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"maxTxsPerCommit\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subtreeID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"depositSubTreeRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"pathToSubTree\",\"type\":\"uint256\"}],\"name\":\"DepositsFinalised\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"enumTypes.Usage\",\"name\":\"batchType\",\"type\":\"uint8\"}],\"name\":\"NewBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"startID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"nDeleted\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"completed\",\"type\":\"bool\"}],\"name\":\"RollbackStatus\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"committed\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"}],\"name\":\"StakeWithdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DOMAIN_NAME\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"DOMAIN_VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ZERO_BYTES32\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"accountRegistry\",\"outputs\":[{\"internalType\":\"contractBLSAccountRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"batches\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"commitmentRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"meta\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"chooser\",\"outputs\":[{\"internalType\":\"contractChooser\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"create2Transfer\",\"outputs\":[{\"internalType\":\"contractCreate2Transfer\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositManager\",\"outputs\":[{\"internalType\":\"contractIDepositManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"deposits\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.TransferBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"internalType\":\"structTypes.TransferCommitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.TransferCommitmentInclusionProof\",\"name\":\"target\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState[]\",\"name\":\"states\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"stateWitnesses\",\"type\":\"bytes32[][]\"},{\"internalType\":\"uint256[4][]\",\"name\":\"pubkeysSender\",\"type\":\"uint256[4][]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"pubkeyWitnessesSender\",\"type\":\"bytes32[][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"pubkeyHashesReceiver\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"pubkeyWitnessesReceiver\",\"type\":\"bytes32[][]\"}],\"internalType\":\"structTypes.SignatureProofWithReceiver\",\"name\":\"signatureProof\",\"type\":\"tuple\"}],\"name\":\"disputeSignatureCreate2Transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"spokeID\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"withdrawRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.MassMigrationBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"internalType\":\"structTypes.MassMigrationCommitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.MMCommitmentInclusionProof\",\"name\":\"target\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState[]\",\"name\":\"states\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"stateWitnesses\",\"type\":\"bytes32[][]\"},{\"internalType\":\"uint256[4][]\",\"name\":\"pubkeys\",\"type\":\"uint256[4][]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"pubkeyWitnesses\",\"type\":\"bytes32[][]\"}],\"internalType\":\"structTypes.SignatureProof\",\"name\":\"signatureProof\",\"type\":\"tuple\"}],\"name\":\"disputeSignatureMassMigration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.TransferBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"internalType\":\"structTypes.TransferCommitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.TransferCommitmentInclusionProof\",\"name\":\"target\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState[]\",\"name\":\"states\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"stateWitnesses\",\"type\":\"bytes32[][]\"},{\"internalType\":\"uint256[4][]\",\"name\":\"pubkeys\",\"type\":\"uint256[4][]\"},{\"internalType\":\"bytes32[][]\",\"name\":\"pubkeyWitnesses\",\"type\":\"bytes32[][]\"}],\"internalType\":\"structTypes.SignatureProof\",\"name\":\"signatureProof\",\"type\":\"tuple\"}],\"name\":\"disputeSignatureTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"bodyRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structTypes.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.CommitmentInclusionProof\",\"name\":\"previous\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.TransferBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"internalType\":\"structTypes.TransferCommitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.TransferCommitmentInclusionProof\",\"name\":\"target\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState\",\"name\":\"state\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.StateMerkleProof[]\",\"name\":\"proofs\",\"type\":\"tuple[]\"}],\"name\":\"disputeTransitionCreate2Transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"bodyRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structTypes.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.CommitmentInclusionProof\",\"name\":\"previous\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"spokeID\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"withdrawRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.MassMigrationBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"internalType\":\"structTypes.MassMigrationCommitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.MMCommitmentInclusionProof\",\"name\":\"target\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState\",\"name\":\"state\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.StateMerkleProof[]\",\"name\":\"proofs\",\"type\":\"tuple[]\"}],\"name\":\"disputeTransitionMassMigration\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"bodyRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structTypes.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.CommitmentInclusionProof\",\"name\":\"previous\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accountRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[2]\",\"name\":\"signature\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256\",\"name\":\"feeReceiver\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"txs\",\"type\":\"bytes\"}],\"internalType\":\"structTypes.TransferBody\",\"name\":\"body\",\"type\":\"tuple\"}],\"internalType\":\"structTypes.TransferCommitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.TransferCommitmentInclusionProof\",\"name\":\"target\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pubkeyID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tokenID\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"}],\"internalType\":\"structTypes.UserState\",\"name\":\"state\",\"type\":\"tuple\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.StateMerkleProof[]\",\"name\":\"proofs\",\"type\":\"tuple[]\"}],\"name\":\"disputeTransitionTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"domainSeparator\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"}],\"name\":\"getBatch\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"commitmentRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"meta\",\"type\":\"bytes32\"}],\"internalType\":\"structTypes.Batch\",\"name\":\"batch\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"invalidBatchMarker\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"keepRollingBack\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"massMigration\",\"outputs\":[{\"internalType\":\"contractMassMigration\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextBatchID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paramBlocksToFinalise\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paramMaxTxsPerCommit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paramMinGasLeft\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paramStakeAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"stateRoots\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[2][]\",\"name\":\"signatures\",\"type\":\"uint256[2][]\"},{\"internalType\":\"uint256[]\",\"name\":\"feeReceivers\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"txss\",\"type\":\"bytes[]\"}],\"name\":\"submitCreate2Transfer\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"bodyRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structTypes.Commitment\",\"name\":\"commitment\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"path\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.CommitmentInclusionProof\",\"name\":\"previous\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"pathAtDepth\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"witness\",\"type\":\"bytes32[]\"}],\"internalType\":\"structTypes.SubtreeVacancyProof\",\"name\":\"vacant\",\"type\":\"tuple\"}],\"name\":\"submitDeposits\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"stateRoots\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[2][]\",\"name\":\"signatures\",\"type\":\"uint256[2][]\"},{\"internalType\":\"uint256[4][]\",\"name\":\"meta\",\"type\":\"uint256[4][]\"},{\"internalType\":\"bytes32[]\",\"name\":\"withdrawRoots\",\"type\":\"bytes32[]\"},{\"internalType\":\"bytes[]\",\"name\":\"txss\",\"type\":\"bytes[]\"}],\"name\":\"submitMassMigration\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"stateRoots\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[2][]\",\"name\":\"signatures\",\"type\":\"uint256[2][]\"},{\"internalType\":\"uint256[]\",\"name\":\"feeReceivers\",\"type\":\"uint256[]\"},{\"internalType\":\"bytes[]\",\"name\":\"txss\",\"type\":\"bytes[]\"}],\"name\":\"submitTransfer\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"contractTransfer\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchID\",\"type\":\"uint256\"}],\"name\":\"withdrawStake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"withdrawalBitmap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"zeroHashAtSubtreeDepth\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x61028060405260006001553480156200001757600080fd5b50604051620056a8380380620056a88339810160408190526200003a91620008f0565b60405180604001604052806006815260200165487562626c6560d01b815250604051806040016040528060018152602001603160f81b8152508585858f8f84608081815250508360a081815250508260c08181525050816001600160a01b0316610100816001600160a01b031660601b81525050806001600160a01b031660e0816001600160a01b031660601b81525050505050505060008280519060200120905060008280519060200120905060007f8b73c3c69bb8fe3d512ecc4cf759cc79239f7b179b0ffacaa9a75d522b39400f90508261016081815250508161018081815250506200012f6200033c60201b60201c565b610140526200014081848462000340565b610120526101a0525050506001600160601b031960608b811b82166101c0528a811b82166101e05289811b82166102005288901b166102205250610240819052604080516338080fa560e21b8152905162000215916001600160a01b038d169163e0203e94916004808201926020929091908290030181600087803b158015620001c957600080fd5b505af1158015620001de573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190620002049190620009b8565b6200038560201b620020ce1760201c565b61026052604051600090620002409087906000805160206200568883398151915290602001620009d1565b60408051601f1981840301815290829052805160209182012092506000916200027d918491600080516020620056888339815191529101620009d1565b60408051601f1981840301815282825280516020918201208383019092528183529092508101620002be6000600133436200087260201b620025b71760201c565b90526001805460009081526020818152604080832085518155949091015193830193909355905491517f3c367559be37298dc40fed468321ed4d44d99ab593f8e0bb7c82dbc84ab0bad69262000318929091819062000a0b565b60405180910390a1505060018054810190555062000a4b9950505050505050505050565b4690565b60008383836200034f6200033c565b3060405160200162000366959493929190620009df565b6040516020818303038152906040528051906020012090509392505050565b600062000391620008d1565b6000805160206200568883398151915281527f633dc4d7da7256660a892f8f1604a44b5432649cc8ec5cb3ced4c4e6ac94dd1d6020808301919091527f890740a8eb06ce9be422cb8da5cdafc2b58c0a5e24036c578de2a433c828ff7d60408301527f3b8ec09e026fdc305365dfc94e189a81b38c7597b3d941c279f042e8206e0bd860608301527fecd50eee38e386bd62be9bedb990706951b65fe053bd9d8a521af753d139e2da60808301527fdefff6d330bb5403f63b14f33b578274160de3a50df4efecf0e0db73bcdd3da560a08301527f617bdd11f7c0a11f49db22f629387a12da7596f9d1704d7465177c63d88ec7d760c08301527f292c23a9aa1d8bea7e2435e555a4a60e379a5a35f3f452bae60121073fb6eead60e08301527fe1cea92ed99acdcb045a6726b2f87107e8a61620a232cf4d7d5b5766b3952e106101008301527f7ad66c0a68c72cb89e4fb4303841966e4062a76ab97451e3b9fb526a5ceb7f826101208301527fe026cc5a4aed3c22a58cbd3d2ac754c9352c5436f638042dca99034e836365166101408301527f3d04cffd8b46a874edf5cfae63077de85f849a660426697b06a829c70dd1409c6101608301527fad676aa337a485e4728a0b240d92b3ef7b3c372d06d189322bfd5f61f1e7203e6101808301527fa2fca4a49658f9fab7aa63289c91b7c7b6c832a6d0e69334ff5b0a3483d09dab6101a08301527f4ebfd9cd7bca2505f7bef59cc1c12ecc708fff26ae4af19abe852afe9e20c8626101c08301527f2def10d13dd169f550f578bda343d9717a138562e0093b380a1120789d53cf106101e08301527f776a31db34a1a0a7caaf862cffdfff1789297ffadc380bd3d39281d340abd3ad6102008301527fe2e7610b87a5fdf3a72ebe271287d923ab990eefac64b6e59d79f8b7e08c46e36102208301527f504364a5c6858bf98fff714ab5be9de19ed31a976860efbd0e772a2efe23e2e06102408301527f4f05f4acb83f5b65168d9fef89d56d4d77b8944015e6b1eed81b0238e2d0dba36102608301527f44a6d974c75b07423e1d6d33f481916fdd45830aea11b6347e700cd8b9f0767c6102808301527fedf260291f734ddac396a956127dde4c34c0cfb8d8052f88ac139658ccf2d5076102a08301527f6075c657a105351e7f0fce53bc320113324a522e8fd52dc878c762551e01a46e6102c08301527f6ca6a3f763a9395f7da16014725ca7ee17e4815c0ff8119bf33f273dee11833b6102e08301527f1c25ef10ffeb3c7d08aa707d17286e0b0d3cbcb50f1bd3b6523b63ba3b52dd0f6103008301527ffffc43bd08273ccf135fd3cacbeef055418e09eb728d727c4d5d5c556cdea7e36103208301527fc5ab8111456b1f28f3c7a0a604b4553ce905cb019c463ee159137af83c350b226103408301527f0ff273fcbf4ae0f2bd88d6cf319ff4004f8d7dca70d4ced4e74d2c74139739e66103608301527f7fa06ba11241ddd5efdc65d4e39c9f6991b74fd4b81b62230808216c876f827c6103808301527f7e275adf313a996c7e2950cac67caba02a5ff925ebf9906b58949f3e77aec5b96103a08301527f8f6162fa308d2b3a15dc33cffac85f13ab349173121645aedf00f471663108be6103c08301527f78ccaaab73373552f207a63599de54d7d8d0c1805f86ce7da15818d09f4cff626103e08301528190849081106200086657fe5b60200201519392505050565b69ffffffff000000000000603082901b16600160501b600160f01b03605084901b1660ff60f01b60f086901b167fff0000000000000000000000000000000000000000000000000000000000000060f888901b16171717949350505050565b6040518061040001604052806020906020820280368337509192915050565b60008060008060008060008060008060006101608c8e03121562000912578687fd5b8b516200091f8162000a32565b60208d0151909b50620009328162000a32565b60408d0151909a50620009458162000a32565b60608d0151909950620009588162000a32565b60808d01519098506200096b8162000a32565b60a08d01519097506200097e8162000a32565b8096505060c08c0151945060e08c015193506101008c015192506101208c015191506101408c015190509295989b509295989b9093969950565b600060208284031215620009ca578081fd5b5051919050565b918252602082015260400190565b9485526020850193909352604084019190915260608301526001600160a01b0316608082015260a00190565b83815260208101839052606081016005831062000a2457fe5b826040830152949350505050565b6001600160a01b038116811462000a4857600080fd5b50565b60805160a05160c05160e05160601c6101005160601c610120516101405161016051610180516101a0516101c05160601c6101e05160601c6102005160601c6102205160601c6102405161026051614b0562000b8360003980610ab65280611c9c5250806111b352806118985280611ea852806120815250806113a95280611860528061193a525080610c925280610dbb528061205252508061117b528061195e5280611e03525080610699528061156e5280611aaa5280611e60525080612df6525080612e38525080612e17525080612d9d525080612dcd5250806105685280610937528061144652806119825280611c29525080610afc52806118f45280612c15525080611e3c5280612ba7525080610ff1528061294f525080610f6d5280611e8452806128ed5280612ce35250614b056000f3fe6080604052600436106101ee5760003560e01c80638a4068dd1161010d578063b32c4d8d116100a0578063d1b243e21161006f578063d1b243e2146104e7578063e42f7554146104fc578063eb84e07614610511578063f698da2514610531578063fe5f960414610546576101ee565b8063b32c4d8d1461046f578063c08caaef1461049d578063ca858d9d146104bd578063d089e11a146104d2576101ee565b8063ac96f0cd116100dc578063ac96f0cd14610410578063acb8cc4914610425578063acf5b54d1461043a578063b02c43d01461044f576101ee565b80638a4068dd146103be5780639251597f146103d357806398d17621146103e6578063a72598ce146103fb576101ee565b80634f6de740116101855780635f6e91d5116101545780635f6e91d5146103525780636c7ac9d814610372578063796f077b146103875780637ae8c568146103a9576101ee565b80634f6de740146102dd5780635ac44282146102fd5780635b097d371461032a5780635f5b95b81461033f576101ee565b80632538507d116101c15780632538507d1461026857806325d5971f1461028857806331c2b7db146102a85780634e23e8c3146102bd576101ee565b8063035695e3146101f3578063069321b0146102085780630d64ac34146102335780630ed75b9c14610246575b600080fd5b6102066102013660046138a9565b610566565b005b34801561021457600080fd5b5061021d610923565b60405161022a91906141f4565b60405180910390f35b610206610241366004613ab0565b610935565b34801561025257600080fd5b5061025b610c90565b60405161022a919061430c565b34801561027457600080fd5b50610206610283366004613c71565b610cb4565b34801561029457600080fd5b506102066102a3366004613b4a565b610e8b565b3480156102b457600080fd5b5061021d610fef565b3480156102c957600080fd5b506102066102d8366004613c12565b611013565b3480156102e957600080fd5b506102066102f8366004613cda565b61128c565b34801561030957600080fd5b5061031d610318366004613b4a565b6113e0565b60405161022a91906149b8565b34801561033657600080fd5b5061021d61143e565b61020661034d366004613993565b611444565b34801561035e57600080fd5b5061020661036d366004613c12565b6116f8565b34801561037e57600080fd5b5061025b6118f2565b34801561039357600080fd5b5061039c611916565b60405161022a9190614320565b3480156103b557600080fd5b5061025b611938565b3480156103ca57600080fd5b5061025b61195c565b6102066103e1366004613993565b611980565b3480156103f257600080fd5b5061025b611c27565b34801561040757600080fd5b5061021d611c4b565b34801561041c57600080fd5b50610206611c51565b34801561043157600080fd5b5061039c611c7d565b34801561044657600080fd5b5061021d611c9a565b34801561045b57600080fd5b5061021d61046a366004613b4a565b611cbe565b34801561047b57600080fd5b5061048f61048a366004613b4a565b611cd0565b60405161022a9291906141cd565b3480156104a957600080fd5b506102066104b8366004613e22565b611ce9565b3480156104c957600080fd5b5061021d611e3a565b3480156104de57600080fd5b5061025b611e5e565b3480156104f357600080fd5b5061021d611e82565b34801561050857600080fd5b5061021d611ea6565b34801561051d57600080fd5b5061020661052c366004613b85565b611eca565b34801561053d57600080fd5b5061021d6120ac565b34801561055257600080fd5b5061021d610561366004613b4a565b6120bc565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b1580156105bf57600080fd5b505afa1580156105d3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906105f79190613882565b6001600160a01b0316336001600160a01b0316146106305760405162461bcd60e51b815260040161062790614722565b60405180910390fd5b600454156106505760405162461bcd60e51b8152600401610627906145f9565b6060896001600160401b038111801561066857600080fd5b50604051908082528060200260200182016040528015610692578160200160208202803683370190505b50905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663ebf0c7176040518163ffffffff1660e01b815260040160206040518083038186803b1580156106f057600080fd5b505afa158015610704573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107289190613a51565b905060005b8b8110156108ff5761073d612fbb565b6040518061010001604052808481526020018d8d8581811061075b57fe5b9050604002016002806020026040519081016040528092919082600260200280828437600092019190915250505081526020018b8b8581811061079a57fe5b9050608002016000600481106107ac57fe5b602002013581526020018989858181106107c257fe5b9050602002013581526020018b8b858181106107da57fe5b9050608002016001600481106107ec57fe5b602002013581526020018b8b8581811061080257fe5b90506080020160026004811061081457fe5b602002013581526020018b8b8581811061082a57fe5b90506080020160036004811061083c57fe5b6020020135815260200187878581811061085257fe5b90506020028101906108649190614a23565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250505091525090508d8d838181106108a857fe5b905060200201356108b8826125fe565b6040516020016108c99291906141cd565b604051602081830303815290604052805190602001208483815181106108eb57fe5b60209081029190910101525060010161072d565b5061091561090c83612655565b8c8360026128eb565b505050505050505050505050565b600080516020614ad983398151915281565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b15801561098e57600080fd5b505afa1580156109a2573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109c69190613882565b6001600160a01b0316336001600160a01b0316146109f65760405162461bcd60e51b815260040161062790614722565b60045415610a165760405162461bcd60e51b8152600401610627906145f9565b60018054600019016000818152602081815260409182902082518084019093528054835284015490820152909190610a4d906129e1565b03836020015114610a705760405162461bcd60e51b815260040161062790614566565b600081815260208190526040902054610a8990846129ee565b610aa55760405162461bcd60e51b8152600401610627906143c6565b82515182516020840151610adb92917f000000000000000000000000000000000000000000000000000000000000000091612a19565b610af75760405162461bcd60e51b8152600401610627906146df565b6000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663d86ee48d6040518163ffffffff1660e01b81526004016040805180830381600087803b158015610b5457600080fd5b505af1158015610b68573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b8c9190613b62565b600185016000818152600260205260409081902083905587519051939550919350917f1e6cd0ae9aa6f6e6ea5cf7f07396a96fc1ea5ae4ec7d5054f53e204c47f79bd891610bdd91869186916149f5565b60405180910390a16000610bfa8387600001518860200151612a31565b9050600081600080516020614ad983398151915260001b604051602001610c229291906141cd565b60408051601f198184030181529082905280516020918201209250600091610c5c918491600080516020614ad983398151915291016141cd565b604051602081830303815290604052805190602001209050610c858160016000801b60046128eb565b505050505050505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000838152602081815260409182902082518084019093528054835260010154908201528390610ce390612adc565b4310610d015760405162461bcd60e51b815260040161062790614630565b600454811080610d115750600454155b610d2d5760405162461bcd60e51b8152600401610627906147ed565b600084815260208190526040902054610d469084612aec565b610d625760405162461bcd60e51b815260040161062790614667565b610d6a61300a565b6040805160a081018252855160209081015181015182528651518183015286510151519181019190915260608101610da06120ac565b815260200185600001516020015160e00151815250905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663be42c7c48386886000015160200151604001516040518463ffffffff1660e01b8152600401610e1593929190614982565b60206040518083038186803b158015610e2d57600080fd5b505afa158015610e41573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e659190613a95565b9050600081600b811115610e7557fe5b14610e8357610e8386612aff565b505050505050565b600081815260208181526040918290208251808401909352805483526001015490820152610eb890612b0f565b6001600160a01b0316336001600160a01b031614610ee85760405162461bcd60e51b815260040161062790614450565b600081815260208181526040918290208251808401909352805483526001015490820152610f1590612adc565b4311610f335760405162461bcd60e51b8152600401610627906145a8565b610f3e816003612b25565b15610f5b5760405162461bcd60e51b815260040161062790614840565b610f66816003612b4c565b60405133907f000000000000000000000000000000000000000000000000000000000000000080156108fc02916000818181858888f19350505050158015610fb2573d6000803e3d6000fd5b507f1248d48e2de900a1010c7fce73506969ecec243600bfc08b641b158f26d857cd3382604051610fe49291906141db565b60405180910390a150565b7f000000000000000000000000000000000000000000000000000000000000000081565b600084815260208181526040918290208251808401909352805483526001015490820152849061104290612adc565b43106110605760405162461bcd60e51b815260040161062790614630565b6004548110806110705750600454155b61108c5760405162461bcd60e51b8152600401610627906147ed565b8484846020015160008082600014156110e0575060001984016000818152602081815260409182902082518084019093528054835260019081015491830191909152906110d8906129e1565b0391506110e9565b50506000198101835b8184602001511461110c5760405162461bcd60e51b815260040161062790614566565b60008181526020819052604090205461112590856129ee565b6111415760405162461bcd60e51b815260040161062790614374565b60008a81526020819052604090205461115a9089612b74565b6111765760405162461bcd60e51b8152600401610627906144d7565b6000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663929314928c60000151600001517f00000000000000000000000000000000000000000000000000000000000000008d6000015160200151604001518e6000015160200151606001518e6040518663ffffffff1660e01b815260040161120d9594939291906142c9565b604080518083038186803b15801561122457600080fd5b505afa158015611238573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061125c9190613a69565b9092509050600081600b81111561126f57fe5b14158061127e57508951518214155b15610915576109158c612aff565b60008381526020818152604091829020825180840190935280548352600101549082015283906112bb90612adc565b43106112d95760405162461bcd60e51b815260040161062790614630565b6004548110806112e95750600454155b6113055760405162461bcd60e51b8152600401610627906147ed565b60008481526020819052604090205461131e9084612b74565b61133a5760405162461bcd60e51b81526004016106279061451f565b61134261300a565b6040805160a0810182528551602090810151810151825286515181830152865101515191810191909152606081016113786120ac565b81528551602090810151606001519101526040516001627e8b6b60e01b031981529091506000906001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063ff81749590610e1590859088906004016148ae565b6113e8613041565b6000828152602081905260409020600101546114165760405162461bcd60e51b815260040161062790614759565b5060009081526020818152604091829020825180840190935280548352600101549082015290565b60045481565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b15801561149d57600080fd5b505afa1580156114b1573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906114d59190613882565b6001600160a01b0316336001600160a01b0316146115055760405162461bcd60e51b815260040161062790614722565b600454156115255760405162461bcd60e51b8152600401610627906145f9565b6060876001600160401b038111801561153d57600080fd5b50604051908082528060200260200182016040528015611567578160200160208202803683370190505b50905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663ebf0c7176040518163ffffffff1660e01b815260040160206040518083038186803b1580156115c557600080fd5b505afa1580156115d9573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906115fd9190613a51565b90506000805b8a8110156116d557828a8a8381811061161857fe5b90506040020189898481811061162a57fe5b9050602002013588888581811061163d57fe5b905060200281019061164f9190614a23565b60405160200161166395949392919061410e565b6040516020818303038152906040528051906020012091508b8b8281811061168757fe5b90506020020135826040516020016116a09291906141cd565b604051602081830303815290604052805190602001208482815181106116c257fe5b6020908102919091010152600101611603565b506116eb6116e284612655565b8b8460016128eb565b5050505050505050505050565b600084815260208181526040918290208251808401909352805483526001015490820152849061172790612adc565b43106117455760405162461bcd60e51b815260040161062790614630565b6004548110806117555750600454155b6117715760405162461bcd60e51b8152600401610627906147ed565b8484846020015160008082600014156117c5575060001984016000818152602081815260409182902082518084019093528054835260019081015491830191909152906117bd906129e1565b0391506117ce565b50506000198101835b818460200151146117f15760405162461bcd60e51b815260040161062790614566565b60008181526020819052604090205461180a90856129ee565b6118265760405162461bcd60e51b815260040161062790614374565b60008a81526020819052604090205461183f9089612b74565b61185b5760405162461bcd60e51b8152600401610627906144d7565b6000807f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663336920368c60000151600001517f00000000000000000000000000000000000000000000000000000000000000008d6000015160200151604001518e6000015160200151606001518e6040518663ffffffff1660e01b815260040161120d9594939291906142c9565b7f000000000000000000000000000000000000000000000000000000000000000081565b60405180604001604052806006815260200165487562626c6560d01b81525081565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663e9790d026040518163ffffffff1660e01b815260040160206040518083038186803b1580156119d957600080fd5b505afa1580156119ed573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611a119190613882565b6001600160a01b0316336001600160a01b031614611a415760405162461bcd60e51b815260040161062790614722565b60045415611a615760405162461bcd60e51b8152600401610627906145f9565b6060876001600160401b0381118015611a7957600080fd5b50604051908082528060200260200182016040528015611aa3578160200160208202803683370190505b50905060007f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031663ebf0c7176040518163ffffffff1660e01b815260040160206040518083038186803b158015611b0157600080fd5b505afa158015611b15573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611b399190613a51565b90506000805b8a811015611c1157828a8a83818110611b5457fe5b905060400201898984818110611b6657fe5b90506020020135888885818110611b7957fe5b9050602002810190611b8b9190614a23565b604051602001611b9f95949392919061410e565b6040516020818303038152906040528051906020012091508b8b82818110611bc357fe5b9050602002013582604051602001611bdc9291906141cd565b60405160208183030381529060405280519060200120848281518110611bfe57fe5b6020908102919091010152600101611b3f565b506116eb611c1e84612655565b8b8460036128eb565b7f000000000000000000000000000000000000000000000000000000000000000081565b60015481565b600060045411611c735760405162461bcd60e51b815260040161062790614333565b611c7b612b87565b565b604051806040016040528060018152602001603160f81b81525081565b7f000000000000000000000000000000000000000000000000000000000000000081565b60026020526000908152604090205481565b6000602081905290815260409020805460019091015482565b6000838152602081815260409182902082518084019093528054835260010154908201528390611d1890612adc565b4310611d365760405162461bcd60e51b815260040161062790614630565b600454811080611d465750600454155b611d625760405162461bcd60e51b8152600401610627906147ed565b600084815260208190526040902054611d7b9084612b74565b611d975760405162461bcd60e51b81526004016106279061451f565b611d9f61300a565b6040805160a081018252855160209081015181015182528651518183015286510151519181019190915260608101611dd56120ac565b81528551602090810151606001519101526040516344ec5a7760e01b81529091506000906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906344ec5a7790610e15908590889060040161495d565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000848152602081815260409182902082518084019093528054835260010154908201528490611ef990612adc565b4310611f175760405162461bcd60e51b815260040161062790614630565b600454811080611f275750600454155b611f435760405162461bcd60e51b8152600401610627906147ed565b848484602001516000808260001415611f9757506000198401600081815260208181526040918290208251808401909352805483526001908101549183019190915290611f8f906129e1565b039150611fa0565b50506000198101835b81846020015114611fc35760405162461bcd60e51b815260040161062790614566565b600081815260208190526040902054611fdc90856129ee565b611ff85760405162461bcd60e51b815260040161062790614374565b60008a8152602081905260409020546120119089612aec565b61202d5760405162461bcd60e51b8152600401610627906144d7565b88515188516020015160405163ab5a164f60e01b815260009283926001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169263ab5a164f9261120d9290917f000000000000000000000000000000000000000000000000000000000000000091908f90600401614229565b60006120b6612d99565b90505b90565b60036020526000908152604090205481565b60006120d8613058565b600080516020614ad983398151915281527f633dc4d7da7256660a892f8f1604a44b5432649cc8ec5cb3ced4c4e6ac94dd1d6020808301919091527f890740a8eb06ce9be422cb8da5cdafc2b58c0a5e24036c578de2a433c828ff7d60408301527f3b8ec09e026fdc305365dfc94e189a81b38c7597b3d941c279f042e8206e0bd860608301527fecd50eee38e386bd62be9bedb990706951b65fe053bd9d8a521af753d139e2da60808301527fdefff6d330bb5403f63b14f33b578274160de3a50df4efecf0e0db73bcdd3da560a08301527f617bdd11f7c0a11f49db22f629387a12da7596f9d1704d7465177c63d88ec7d760c08301527f292c23a9aa1d8bea7e2435e555a4a60e379a5a35f3f452bae60121073fb6eead60e08301527fe1cea92ed99acdcb045a6726b2f87107e8a61620a232cf4d7d5b5766b3952e106101008301527f7ad66c0a68c72cb89e4fb4303841966e4062a76ab97451e3b9fb526a5ceb7f826101208301527fe026cc5a4aed3c22a58cbd3d2ac754c9352c5436f638042dca99034e836365166101408301527f3d04cffd8b46a874edf5cfae63077de85f849a660426697b06a829c70dd1409c6101608301527fad676aa337a485e4728a0b240d92b3ef7b3c372d06d189322bfd5f61f1e7203e6101808301527fa2fca4a49658f9fab7aa63289c91b7c7b6c832a6d0e69334ff5b0a3483d09dab6101a08301527f4ebfd9cd7bca2505f7bef59cc1c12ecc708fff26ae4af19abe852afe9e20c8626101c08301527f2def10d13dd169f550f578bda343d9717a138562e0093b380a1120789d53cf106101e08301527f776a31db34a1a0a7caaf862cffdfff1789297ffadc380bd3d39281d340abd3ad6102008301527fe2e7610b87a5fdf3a72ebe271287d923ab990eefac64b6e59d79f8b7e08c46e36102208301527f504364a5c6858bf98fff714ab5be9de19ed31a976860efbd0e772a2efe23e2e06102408301527f4f05f4acb83f5b65168d9fef89d56d4d77b8944015e6b1eed81b0238e2d0dba36102608301527f44a6d974c75b07423e1d6d33f481916fdd45830aea11b6347e700cd8b9f0767c6102808301527fedf260291f734ddac396a956127dde4c34c0cfb8d8052f88ac139658ccf2d5076102a08301527f6075c657a105351e7f0fce53bc320113324a522e8fd52dc878c762551e01a46e6102c08301527f6ca6a3f763a9395f7da16014725ca7ee17e4815c0ff8119bf33f273dee11833b6102e08301527f1c25ef10ffeb3c7d08aa707d17286e0b0d3cbcb50f1bd3b6523b63ba3b52dd0f6103008301527ffffc43bd08273ccf135fd3cacbeef055418e09eb728d727c4d5d5c556cdea7e36103208301527fc5ab8111456b1f28f3c7a0a604b4553ce905cb019c463ee159137af83c350b226103408301527f0ff273fcbf4ae0f2bd88d6cf319ff4004f8d7dca70d4ced4e74d2c74139739e66103608301527f7fa06ba11241ddd5efdc65d4e39c9f6991b74fd4b81b62230808216c876f827c6103808301527f7e275adf313a996c7e2950cac67caba02a5ff925ebf9906b58949f3e77aec5b96103a08301527f8f6162fa308d2b3a15dc33cffac85f13ab349173121645aedf00f471663108be6103c08301527f78ccaaab73373552f207a63599de54d7d8d0c1805f86ce7da15818d09f4cff626103e08301528190849081106125ab57fe5b60200201519392505050565b69ffffffff000000000000603082901b16600160501b600160f01b03605084901b1660ff60f01b60f086901b166001600160f81b031960f888901b16171717949350505050565b80516020808301516040808501516060860151608087015160a088015160c089015160e08a01519551600099612638999098979101614139565b604051602081830303815290604052805190602001209050919050565b60006020825111156126795760405162461bcd60e51b8152600401610627906147b6565b612681613077565b600080516020614ad983398151915281527f633dc4d7da7256660a892f8f1604a44b5432649cc8ec5cb3ced4c4e6ac94dd1d60208201527f890740a8eb06ce9be422cb8da5cdafc2b58c0a5e24036c578de2a433c828ff7d60408201527f3b8ec09e026fdc305365dfc94e189a81b38c7597b3d941c279f042e8206e0bd86060808301919091527fecd50eee38e386bd62be9bedb990706951b65fe053bd9d8a521af753d139e2da60808301527fdefff6d330bb5403f63b14f33b578274160de3a50df4efecf0e0db73bcdd3da560a08301528351600181811692918101901c90816001600160401b038111801561277857600080fd5b506040519080825280602002602001820160405280156127a2578160200160208202803683370190505b5090508560005b60005b858503811015612831576000600182901b90508381815181106127cb57fe5b60200260200101518482600101815181106127e257fe5b60200260200101516040516020016127fb9291906141cd565b6040516020818303038152906040528051906020012085838151811061281d57fe5b6020908102919091010152506001016127ac565b85600114156128a35782600182901b8151811061284a57fe5b602002602001015187836006811061285e57fe5b60200201516040516020016128749291906141cd565b6040516020818303038152906040528051906020012084828151811061289657fe5b6020026020010181815250505b84600114156128b257506128c9565b5060018085169550938401841c93839250016127a9565b826000815181106128d657fe5b60200260200101519650505050505050919050565b7f000000000000000000000000000000000000000000000000000000000000000034101561292b5760405162461bcd60e51b815260040161062790614877565b604051806040016040528085815260200161297583600481111561294b57fe5b86337f000000000000000000000000000000000000000000000000000000000000000043016125b7565b9052600180546000908152602081815260409182902084518155930151928201929092555490517f3c367559be37298dc40fed468321ed4d44d99ab593f8e0bb7c82dbc84ab0bad6916129cb91859085906149cf565b60405180910390a1505060018054810190555050565b6020015160f01c60ff1690565b6000612a1083612a018460000151612e63565b84602001518560400151612a19565b90505b92915050565b600084612a27858585612a31565b1495945050505050565b600083815b8351811015612ad357600185821c16612a8c5781848281518110612a5657fe5b6020026020010151604051602001612a6f9291906141cd565b604051602081830303815290604052805190602001209150612acb565b838181518110612a9857fe5b602002602001015182604051602001612ab29291906141cd565b6040516020818303038152906040528051906020012091505b600101612a36565b50949350505050565b6020015160301c63ffffffff1690565b6000612a1083612a018460000151612e80565b6004819055612b0c612b87565b50565b602081015160501c6001600160a01b0316919050565b610100820460009081526020919091526040902054600160ff9092169190911b9081161490565b61010082046000908152602091909152604090208054600160ff9093169290921b9091179055565b6000612a1083612a018460000151612ea5565b6004546001546000918291908103908290600019015b82821015612c88577f00000000000000000000000000000000000000000000000000000000000000005a11612bd157612c88565b8181036000818152602081815260408083208381556001018390556002909152902054955093508415612c7d576040516356f0001360e11b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000169063ade0002690612c4a9088906004016141f4565b600060405180830381600087803b158015612c6457600080fd5b505af1158015612c78573d6000803e3d6000fd5b505050505b600190910190612b9d565b6001805483900390558282148015612ca05760006004555b7f595cb200f38fcdbf73429ffce272e53b2e923c231b2059ce20d1f4909cff1bec828483604051612cd393929190614a0b565b60405180910390a16000612d07847f0000000000000000000000000000000000000000000000000000000000000000612eb9565b90506000612d216003612d1b846002612eb9565b90612ef3565b90506000612d2f8383612f25565b604051909150339083156108fc029084906000818181858888f19350505050158015612d5f573d6000803e3d6000fd5b5060405160009082156108fc0290839083818181858288f19350505050158015612d8d573d6000803e3d6000fd5b50505050505050505050565b60007f0000000000000000000000000000000000000000000000000000000000000000612dc4612f4d565b1415612df157507f00000000000000000000000000000000000000000000000000000000000000006120b9565b612e5c7f00000000000000000000000000000000000000000000000000000000000000007f00000000000000000000000000000000000000000000000000000000000000007f0000000000000000000000000000000000000000000000000000000000000000612f51565b90506120b9565b6000816000015182602001516040516020016126389291906141cd565b60008160000151612e9483602001516125fe565b6040516020016126389291906141cd565b60008160000151612e948360200151612f92565b600082612ec857506000612a13565b82820282848281612ed557fe5b0414612a105760405162461bcd60e51b81526004016106279061469e565b6000808211612f145760405162461bcd60e51b8152600401610627906144a0565b818381612f1d57fe5b049392505050565b600082821115612f475760405162461bcd60e51b815260040161062790614419565b50900390565b4690565b6000838383612f5e612f4d565b30604051602001612f739594939291906141fd565b6040516020818303038152906040528051906020012090509392505050565b600081600001518260200151836040015184606001516040516020016126389493929190614192565b6040805161010081019091526000815260208101612fd7613095565b81526020016000815260200160008019168152602001600081526020016000815260200160008152602001606081525090565b6040518060a0016040528061301d613095565b81526000602082018190526040820181905260608083019190915260809091015290565b604080518082019091526000808252602082015290565b6040518061040001604052806020906020820280368337509192915050565b6040518060c001604052806006906020820280368337509192915050565b60405180604001604052806002906020820280368337509192915050565b600082601f8301126130c3578081fd5b81356130d66130d182614a8d565b614a67565b818152915060208083019084810160005b84811015613110576130fe888484358a010161329e565b845292820192908201906001016130e7565b505050505092915050565b60008083601f84011261312c578081fd5b5081356001600160401b03811115613142578182fd5b60208301915083602060408302850101111561315d57600080fd5b9250929050565b60008083601f840112613175578182fd5b5081356001600160401b0381111561318b578182fd5b60208301915083602060808302850101111561315d57600080fd5b6000601f83818401126131b7578182fd5b82356131c56130d182614a8d565b81815292506020808401908581016080808502880183018910156131e857600080fd5b60005b8581101561325057898784011261320157600080fd5b61320a82614a67565b80848486018d81111561321c57600080fd5b60005b600481101561323c5782358552938801939188019160010161321f565b5091885250958501959350506001016131eb565b5050505050505092915050565b60008083601f84011261326e578182fd5b5081356001600160401b03811115613284578182fd5b602083019150836020808302850101111561315d57600080fd5b600082601f8301126132ae578081fd5b81356132bc6130d182614a8d565b8181529150602080830190848101818402860182018710156132dd57600080fd5b60005b84811015613110578135845292820192908201906001016132e0565b600082601f83011261330c578081fd5b813561331a6130d182614a8d565b818152915060208083019084810160005b84811015613110578135870160a080601f19838c0301121561334c57600080fd5b6133566040614a67565b6133628b87850161383d565b815290820135906001600160401b0382111561337d57600080fd5b61338b8b878486010161329e565b8187015286525050928201929082019060010161332b565b600082601f8301126133b3578081fd5b81356133c16130d182614a8d565b81815291506020808301908481016080808502870183018810156133e457600080fd5b60005b8581101561340b576133f9898461383d565b855293830193918101916001016133e7565b50505050505092915050565b600082601f830112613427578081fd5b6134316040614a67565b905080828460408501111561344557600080fd5b60005b6002811015613467578135835260209283019290910190600101613448565b50505092915050565b600082601f830112613480578081fd5b81356001600160401b03811115613495578182fd5b6134a8601f8201601f1916602001614a67565b91508082528360208285010111156134bf57600080fd5b8060208401602084013760009082016020015292915050565b8051600c8110612a1357600080fd5b600081830360808112156134f9578182fd5b6135036060614a67565b9150604081121561351357600080fd5b5061351e6040614a67565b8235815260208301356020820152808252506040820135602082015260608201356001600160401b0381111561355357600080fd5b61355f8482850161329e565b60408301525092915050565b60006060828403121561357c578081fd5b6135866060614a67565b905081356001600160401b038082111561359f57600080fd5b90830190604082860312156135b357600080fd5b6135bd6040614a67565b823581526020830135828111156135d357600080fd5b6135df87828601613614565b6020830152508084525060208401356020840152604084013591508082111561360757600080fd5b5061355f8482850161329e565b60006101208284031215613626578081fd5b61010061363281614a67565b9150823582526136458460208501613417565b6020830152606083013560408301526080830135606083015260a0830135608083015260c083013560a083015260e083013560c08301528083013590506001600160401b0381111561369657600080fd5b6136a284828501613470565b60e08301525092915050565b6000608082840312156136bf578081fd5b6136c96080614a67565b905081356001600160401b03808211156136e257600080fd5b6136ee858386016133a3565b8352602084013591508082111561370457600080fd5b613710858386016130b3565b6020840152604084013591508082111561372957600080fd5b613735858386016131a6565b6040840152606084013591508082111561374e57600080fd5b5061375b848285016130b3565b60608301525092915050565b600060a08284031215613778578081fd5b6137826080614a67565b9050813581526137958360208401613417565b60208201526060820135604082015260808201356001600160401b038111156137bd57600080fd5b61375b84828501613470565b6000606082840312156137da578081fd5b6137e46060614a67565b905081356001600160401b03808211156137fd57600080fd5b908301906040828603121561381157600080fd5b61381b6040614a67565b8235815260208301358281111561383157600080fd5b6135df87828601613767565b60006080828403121561384e578081fd5b6138586080614a67565b90508135815260208201356020820152604082013560408201526060820135606082015292915050565b600060208284031215613893578081fd5b81516001600160a01b0381168114612a10578182fd5b60008060008060008060008060008060a08b8d0312156138c7578586fd5b8a356001600160401b03808211156138dd578788fd5b6138e98e838f0161325d565b909c509a5060208d0135915080821115613901578788fd5b61390d8e838f0161311b565b909a50985060408d0135915080821115613925578788fd5b6139318e838f01613164565b909850965060608d0135915080821115613949578586fd5b6139558e838f0161325d565b909650945060808d013591508082111561396d578384fd5b5061397a8d828e0161325d565b915080935050809150509295989b9194979a5092959850565b6000806000806000806000806080898b0312156139ae578182fd5b88356001600160401b03808211156139c4578384fd5b6139d08c838d0161325d565b909a50985060208b01359150808211156139e8578384fd5b6139f48c838d0161311b565b909850965060408b0135915080821115613a0c578384fd5b613a188c838d0161325d565b909650945060608b0135915080821115613a30578384fd5b50613a3d8b828c0161325d565b999c989b5096995094979396929594505050565b600060208284031215613a62578081fd5b5051919050565b60008060408385031215613a7b578182fd5b82519150613a8c84602085016134d8565b90509250929050565b600060208284031215613aa6578081fd5b612a1083836134d8565b60008060408385031215613ac2578182fd5b82356001600160401b0380821115613ad8578384fd5b613ae4868387016134e7565b93506020850135915080821115613af9578283fd5b9084019060408287031215613b0c578283fd5b613b166040614a67565b82358152602083013582811115613b2b578485fd5b613b378882860161329e565b6020830152508093505050509250929050565b600060208284031215613b5b578081fd5b5035919050565b60008060408385031215613b74578182fd5b505080516020909101519092909150565b60008060008060808587031215613b9a578182fd5b8435935060208501356001600160401b0380821115613bb7578384fd5b613bc3888389016134e7565b94506040870135915080821115613bd8578384fd5b613be48883890161356b565b93506060870135915080821115613bf9578283fd5b50613c06878288016132fc565b91505092959194509250565b60008060008060808587031215613c27578182fd5b8435935060208501356001600160401b0380821115613c44578384fd5b613c50888389016134e7565b94506040870135915080821115613c65578384fd5b613be4888389016137c9565b600080600060608486031215613c85578081fd5b8335925060208401356001600160401b0380821115613ca2578283fd5b613cae8783880161356b565b93506040860135915080821115613cc3578283fd5b50613cd0868287016136ae565b9150509250925092565b600080600060608486031215613cee578081fd5b8335925060208401356001600160401b0380821115613d0b578283fd5b613d17878388016137c9565b93506040860135915080821115613d2c578283fd5b9085019060c08288031215613d3f578283fd5b613d4960c0614a67565b823582811115613d57578485fd5b613d63898286016133a3565b825250602083013582811115613d77578485fd5b613d83898286016130b3565b602083015250604083013582811115613d9a578485fd5b613da6898286016131a6565b604083015250606083013582811115613dbd578485fd5b613dc9898286016130b3565b606083015250608083013582811115613de0578485fd5b613dec8982860161329e565b60808301525060a083013582811115613e03578485fd5b613e0f898286016130b3565b60a0830152508093505050509250925092565b600080600060608486031215613e36578081fd5b8335925060208401356001600160401b0380821115613e53578283fd5b613cae878388016137c9565b6000815180845260208085018081965082840281019150828601855b85811015613ea5578284038952613e93848351613f10565b98850198935090840190600101613e7b565b5091979650505050505050565b6000815180845260208085019450808401835b83811015613f0557815187865b6004811015613eef57825182529185019190850190600101613ed2565b5050506080969096019590820190600101613ec5565b509495945050505050565b6000815180845260208085019450808401835b83811015613f0557815187529582019590820190600101613f23565b6000815180845260208085018081965082840281019150828601855b85811015613ea5578284038952815160a0613f778683516140e8565b868201519150806080870152613f8f81870183613f10565b9a87019a9550505090840190600101613f5b565b6000815180845260208085019450808401835b83811015613f0557613fc98783516140e8565b6080969096019590820190600101613fb6565b8060005b6002811015613fff578151845260209384019390910190600101613fe0565b50505050565b6000815180845261401d816020860160208601614aac565b601f01601f19169290920160200192915050565b600061403e838351613fdc565b602082015160408401526040820151606084015260608201516080840152608082015160c060a085015261407560c0850182614005565b949350505050565b60008151608084526140926080850182613fa3565b9050602083015184820360208601526140ab8282613e5f565b915050604083015184820360408601526140c58282613eb2565b915050606083015184820360608601526140df8282613e5f565b95945050505050565b805182526020810151602083015260408101516040830152606081015160608301525050565b6000868252604086602084013760608201859052828460808401379101608001908152949350505050565b600089825261414b602083018a613fdc565b8760608301528660808301528560a08301528460c08301528360e0830152610100835161417e8183860160208801614aac565b929092019091019998505050505050505050565b60008582526141a46020830186613fdc565b83606083015282516141bd816080850160208701614aac565b9190910160800195945050505050565b918252602082015260400190565b6001600160a01b03929092168252602082015260400190565b90815260200190565b9485526020850193909352604084019190915260608301526001600160a01b0316608082015260a00190565b60008582528460208301526080604083015283516080830152602084015161425460a0840182613fdc565b50604084015160e083015260608401516101008301526080840151610120818185015260a086015161014085015260c086015161016085015260e0860151915080610180850152506142aa6101a0840182614005565b905082810360608401526142be8185613f3f565b979650505050505050565b600086825285602083015284604083015260a060608301526142ee60a0830185614005565b82810360808401526143008185613f3f565b98975050505050505050565b6001600160a01b0391909116815260200190565b600060208252612a106020830184614005565b60208082526021908201527f42617463684d616e616765723a204973206e6f7420726f6c6c696e67206261636040820152606b60f81b606082015260800190565b60208082526032908201527f70726576696f757320636f6d6d69746d656e7420697320616273656e7420696e604082015271040e8d0ca40c6eae4e4cadce840c4c2e8c6d60731b606082015260800190565b60208082526033908201527f70726576696f757320636f6d6d69746d656e7420697320616273656e7420696e604082015272040e8d0ca40e0e4caecd2deeae640c4c2e8c6d606b1b606082015260800190565b6020808252601e908201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604082015260600190565b60208082526030908201527f596f7520617265206e6f742074686520636f727265637420636f6d6d6974746560408201526f0e440ccdee440e8d0d2e640c4c2e8c6d60831b606082015260800190565b6020808252601a908201527f536166654d6174683a206469766973696f6e206279207a65726f000000000000604082015260600190565b60208082526028908201527f54617267657420636f6d6d69746d656e7420697320616273656e7420696e20746040820152670d0ca40c4c2e8c6d60c31b606082015260800190565b60208082526027908201527f526f6c6c75703a20436f6d6d69746d656e74206e6f742070726573656e7420696040820152660dc40c4c2e8c6d60cb1b606082015260800190565b60208082526022908201527f70726576696f757320636f6d6d69746d656e74206861732077726f6e672070616040820152610e8d60f31b606082015260800190565b60208082526031908201527f54686973206261746368206973206e6f74207965742066696e616c697365642c60408201527020636865636b206261636b20736f6f6e2160781b606082015260800190565b6020808252601d908201527f42617463684d616e616765723a20497320726f6c6c696e67206261636b000000604082015260600190565b60208082526017908201527f426174636820616c72656164792066696e616c69736564000000000000000000604082015260600190565b6020808252601f908201527f436f6d6d69746d656e74206e6f742070726573656e7420696e20626174636800604082015260600190565b60208082526021908201527f536166654d6174683a206d756c7469706c69636174696f6e206f766572666c6f6040820152607760f81b606082015260800190565b60208082526023908201527f526f6c6c75703a2053746174652073756274726565206973206e6f7420766163604082015262185b9d60ea1b606082015260800190565b60208082526018908201527f526f6c6c75703a20496e76616c69642070726f706f7365720000000000000000604082015260600190565b6020808252603f908201527f42617463682069642067726561746572207468616e20746f74616c206e756d6260408201527f6572206f6620626174636865732c20696e76616c696420626174636820696400606082015260800190565b6020808252601b908201527f4d65726b6c65547265653a20546f6f206d616e79206c65617665730000000000604082015260600190565b60208082526033908201527f416c7265616479207375636365737366756c6c792064697370757465642e20526040820152726f6c6c206261636b20696e2070726f6365737360681b606082015260800190565b60208082526019908201527f526f6c6c75703a20416c72656164792077697468647261776e00000000000000604082015260600190565b6020808252601a908201527f526f6c6c75703a2077726f6e67207374616b6520616d6f756e74000000000000604082015260600190565b6000604082526148c16040830185614031565b8281036020840152835160c082526148dc60c0830182613fa3565b9050602085015182820360208401526148f58282613e5f565b9150506040850151828203604084015261490f8282613eb2565b915050606085015182820360608401526149298282613e5f565b915050608085015182820360808401526149438282613f10565b91505060a085015182820360a08401526142be8282613e5f565b6000604082526149706040830185614031565b82810360208401526140df818561407d565b6000606082526149956060830186614031565b82810360208401526149a7818661407d565b915050826040830152949350505050565b815181526020918201519181019190915260400190565b8381526020810183905260608101600583106149e757fe5b826040830152949350505050565b9283526020830191909152604082015260600190565b92835260208301919091521515604082015260600190565b6000808335601e19843603018112614a39578283fd5b8301803591506001600160401b03821115614a52578283fd5b60200191503681900382131561315d57600080fd5b6040518181016001600160401b0381118282101715614a8557600080fd5b604052919050565b60006001600160401b03821115614aa2578081fd5b5060209081020190565b60005b83811015614ac7578181015183820152602001614aaf565b83811115613fff575050600091015256fe290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563a164736f6c634300060c000a290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563",
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

// SubmitDeposits is a paid mutator transaction binding the contract method 0x0d64ac34.
//
// Solidity: function submitDeposits(((bytes32,bytes32),uint256,bytes32[]) previous, (uint256,bytes32[]) vacant) payable returns()
func (_Rollup *RollupTransactor) SubmitDeposits(opts *bind.TransactOpts, previous TypesCommitmentInclusionProof, vacant TypesSubtreeVacancyProof) (*types.Transaction, error) {
	return _Rollup.contract.Transact(opts, "submitDeposits", previous, vacant)
}

// SubmitDeposits is a paid mutator transaction binding the contract method 0x0d64ac34.
//
// Solidity: function submitDeposits(((bytes32,bytes32),uint256,bytes32[]) previous, (uint256,bytes32[]) vacant) payable returns()
func (_Rollup *RollupSession) SubmitDeposits(previous TypesCommitmentInclusionProof, vacant TypesSubtreeVacancyProof) (*types.Transaction, error) {
	return _Rollup.Contract.SubmitDeposits(&_Rollup.TransactOpts, previous, vacant)
}

// SubmitDeposits is a paid mutator transaction binding the contract method 0x0d64ac34.
//
// Solidity: function submitDeposits(((bytes32,bytes32),uint256,bytes32[]) previous, (uint256,bytes32[]) vacant) payable returns()
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
