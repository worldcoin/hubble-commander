package api

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var (
	ErrMassMigrationWithSenderNotFound = fmt.Errorf(
		"mass migration with given sender " +
			"was not found in a commitment with given commitment index",
	)

	APIWithdrawProofCouldNotBeCalculated = NewAPIError(
		50005,
		"withdraw proof could not be calculated for a given batch",
	)
	APIErrMassMigrationWithSenderNotFound = NewAPIError(
		50006,
		"mass migration with given transaction hash was not found in a commitment with given commitment index",
	)
)

var getWithdrawProofAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError:           APIWithdrawProofCouldNotBeCalculated,
	ErrMassMigrationWithSenderNotFound: APIErrMassMigrationWithSenderNotFound,
}

func (a *API) GetWithdrawProof(batchID models.Uint256, commitmentIndex uint8, transactionHash common.Hash) (*dto.WithdrawProof, error) {
	if !a.cfg.EnableProofMethods {
		return nil, errProofMethodsDisabled
	}
	withdrawTreeProofAndRoot, err := a.unsafeGetWithdrawProof(batchID, commitmentIndex, transactionHash)
	if err != nil {
		return nil, sanitizeError(err, getWithdrawProofAPIErrors)
	}
	return withdrawTreeProofAndRoot, nil
}

func (a *API) unsafeGetWithdrawProof(
	batchID models.Uint256,
	commitmentIndex uint8,
	transactionHash common.Hash,
) (*dto.WithdrawProof, error) {
	// Verifies that batch exists
	_, err := a.storage.GetBatch(batchID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	commitmentID := models.CommitmentID{
		BatchID:      batchID,
		IndexInBatch: commitmentIndex,
	}

	unsortedMassMigrations, err := a.storage.GetMassMigrationsByCommitmentID(commitmentID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// TODO remove when new primary key for transactions with transaction index is implement
	txQueue := executor.NewTxQueue(models.MassMigrationArray(unsortedMassMigrations))
	massMigrations := txQueue.PickTxsForCommitment().ToMassMigrationArray()

	withdrawTree, targetUserState, massMigrationIndex, err := a.generateWithdrawTreeForWithdrawProof(massMigrations, transactionHash)
	if err != nil {
		return nil, err
	}

	return &dto.WithdrawProof{
		UserState: targetUserState,
		Path: dto.MerklePath{
			Path:  *massMigrationIndex,
			Depth: withdrawTree.Depth(),
		},
		Witness: withdrawTree.GetWitness(*massMigrationIndex),
		Root:    withdrawTree.Root(),
	}, nil
}

func (a *API) generateWithdrawTreeForWithdrawProof(
	massMigrations []models.MassMigration,
	transactionHash common.Hash,
) (
	*merkletree.MerkleTree,
	*dto.UserState,
	*uint32,
	error,
) {
	tokenID := models.MakeUint256(0)
	hashes := make([]common.Hash, 0, len(massMigrations))

	var (
		targetUserState    *dto.UserState
		massMigrationIndex int
	)

	for i := range massMigrations {
		senderLeaf, err := a.storage.StateTree.Leaf(massMigrations[i].FromStateID)
		if err != nil {
			return nil, nil, nil, err
		}
		if i == 0 {
			tokenID = senderLeaf.TokenID
		}

		massMigrationUserState := &models.UserState{
			PubKeyID: senderLeaf.PubKeyID,
			TokenID:  tokenID,
			Balance:  massMigrations[i].Amount,
			Nonce:    models.MakeUint256(0),
		}

		hash, err := encoder.HashUserState(massMigrationUserState)
		if err != nil {
			return nil, nil, nil, errors.WithStack(err)
		}
		hashes = append(hashes, *hash)

		if massMigrations[i].Hash == transactionHash {
			dtoMassMigrationUserState := dto.MakeUserState(massMigrationUserState)
			targetUserState = &dtoMassMigrationUserState
			massMigrationIndex = i
		}
	}
	if targetUserState == nil {
		return nil, nil, nil, errors.WithStack(ErrMassMigrationWithSenderNotFound)
	}

	withdrawTree, err := merkletree.NewMerkleTree(hashes)
	if err != nil {
		return nil, nil, nil, errors.WithStack(err)
	}

	return withdrawTree, targetUserState, ref.Uint32(uint32(massMigrationIndex)), nil
}
