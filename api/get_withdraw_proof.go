package api

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var (
	ErrMassMigrationWithTxHashNotFound = fmt.Errorf("mass migration with given transaction hash was not found in a given commitment")
	ErrOnlyMassMigrationBatches        = fmt.Errorf("invalid batch type, only mass migration batches are supported")

	APIWithdrawProofCouldNotBeCalculated = NewAPIError(
		50005,
		"withdraw proof could not be calculated for a given batch",
	)
	APIErrOnlyMassMigrationBatches = NewAPIError(
		50006,
		"invalid batch type, only mass migration batches are supported",
	)
	APIErrMassMigrationWithTxHashNotFound = NewAPIError(
		50007,
		"mass migration with given transaction hash was not found in a given commitment",
	)
)

var getWithdrawProofAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError:           APIWithdrawProofCouldNotBeCalculated,
	ErrOnlyMassMigrationBatches:        APIErrOnlyMassMigrationBatches,
	ErrMassMigrationWithTxHashNotFound: APIErrMassMigrationWithTxHashNotFound,
}

func (a *API) GetWithdrawProof(commitmentID models.CommitmentID, transactionHash common.Hash) (*dto.WithdrawProof, error) {
	if !a.cfg.EnableProofMethods {
		return nil, APIErrProofMethodsDisabled
	}
	withdrawTreeProofAndRoot, err := a.unsafeGetWithdrawProof(commitmentID, transactionHash)
	if err != nil {
		return nil, sanitizeError(err, getWithdrawProofAPIErrors)
	}
	return withdrawTreeProofAndRoot, nil
}

func (a *API) unsafeGetWithdrawProof(
	commitmentID models.CommitmentID,
	transactionHash common.Hash,
) (*dto.WithdrawProof, error) {
	batch, err := a.storage.GetBatch(commitmentID.BatchID)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if batch.Type != batchtype.MassMigration {
		return nil, errors.WithStack(ErrOnlyMassMigrationBatches)
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
	withdrawTree *merkletree.MerkleTree,
	targetUserState *dto.UserState,
	massMigrationIndex *uint32,
	err error,
) {
	tokenID := models.MakeUint256(0)
	hashes := make([]common.Hash, 0, len(massMigrations))

	for i := range massMigrations {
		var senderLeaf *models.StateLeaf
		senderLeaf, err = a.storage.StateTree.Leaf(massMigrations[i].FromStateID)
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

		var hash *common.Hash
		hash, err = encoder.HashUserState(massMigrationUserState)
		if err != nil {
			return nil, nil, nil, errors.WithStack(err)
		}
		hashes = append(hashes, *hash)

		if massMigrations[i].Hash == transactionHash {
			dtoMassMigrationUserState := dto.MakeUserState(massMigrationUserState)
			targetUserState = &dtoMassMigrationUserState
			massMigrationIndex = ref.Uint32(uint32(i))
		}
	}
	if targetUserState == nil {
		return nil, nil, nil, errors.WithStack(ErrMassMigrationWithTxHashNotFound)
	}

	withdrawTree, err = merkletree.NewMerkleTree(hashes)
	if err != nil {
		return nil, nil, nil, errors.WithStack(err)
	}

	return withdrawTree, targetUserState, massMigrationIndex, nil
}
