package api

import (
	"github.com/Worldcoin/hubble-commander/commander/executor"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var getMassMigrationCommitmentProofAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(50004, "mass migration commitment inclusion proof could not be generated"),
}

func (a *API) GetMassMigrationCommitmentProof(commitmentID models.CommitmentID) (*dto.MassMigrationCommitmentProof, error) {
	if !a.cfg.EnableProofMethods {
		return nil, APIErrProofMethodsDisabled
	}
	commitmentInclusionProof, err := a.unsafeGetMassMigrationCommitmentProof(commitmentID)
	if err != nil {
		return nil, sanitizeError(err, getMassMigrationCommitmentProofAPIErrors)
	}
	return commitmentInclusionProof, nil
}

func (a *API) unsafeGetMassMigrationCommitmentProof(commitmentID models.CommitmentID) (*dto.MassMigrationCommitmentProof, error) {
	batch, err := a.storage.GetBatch(commitmentID.BatchID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	commitments, err := a.storage.GetCommitmentsByBatchID(commitmentID.BatchID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	unsortedMassMigrations, err := a.storage.GetMassMigrationsByCommitmentID(commitmentID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// TODO remove when new primary key for transactions with transaction index is implement
	txQueue := executor.NewTxQueue(models.MassMigrationArray(unsortedMassMigrations))
	massMigrations := txQueue.PickTxsForCommitment().ToMassMigrationArray()

	serializedMassMigrations, err := encoder.SerializeMassMigrations(massMigrations)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	commitment := commitments[commitmentID.IndexInBatch].(*models.TxCommitment)

	withdrawTree, meta, err := a.generateWithdrawTreeAndMetaForMassMigrationCommitmentProof(commitment, massMigrations)
	if err != nil {
		return nil, err
	}

	leafHashes := make([]common.Hash, 0, len(commitments))
	for i := range commitments {
		leafHashes = append(leafHashes, commitments[i].LeafHash())
	}
	batchLeafTree, err := merkletree.NewMerkleTree(leafHashes)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	proofBase := &dto.CommitmentInclusionProofBase{
		StateRoot: commitments[commitmentID.IndexInBatch].GetPostStateRoot(),
		Path: &dto.MerklePath{
			Path:  uint32(commitmentID.IndexInBatch),
			Depth: batchLeafTree.Depth(),
		},
		Witness: batchLeafTree.GetWitness(uint32(commitmentID.IndexInBatch)),
	}

	return &dto.MassMigrationCommitmentProof{
		CommitmentInclusionProofBase: *proofBase,
		Body: &dto.MassMigrationBody{
			AccountRoot:  *batch.AccountTreeRoot,
			Signature:    commitment.CombinedSignature,
			Meta:         meta,
			WithdrawRoot: withdrawTree.Root(),
			Transactions: serializedMassMigrations,
		},
	}, nil
}

func (a *API) generateWithdrawTreeAndMetaForMassMigrationCommitmentProof(
	commitment *models.TxCommitment,
	massMigrations []models.MassMigration,
) (*merkletree.MerkleTree, *dto.MassMigrationMeta, error) {
	meta := &dto.MassMigrationMeta{
		Amount:      models.MakeUint256(0),
		FeeReceiver: commitment.FeeReceiver,
	}

	hashes := make([]common.Hash, 0, len(massMigrations))

	for i := range massMigrations {
		meta.Amount = *meta.Amount.Add(&massMigrations[i].Amount)

		senderLeaf, err := a.storage.StateTree.Leaf(massMigrations[i].FromStateID)
		if err != nil {
			return nil, nil, err
		}
		if i == 0 {
			meta.TokenID = senderLeaf.TokenID
			meta.SpokeID = massMigrations[0].SpokeID
		}

		hash, err := encoder.HashUserState(&models.UserState{
			PubKeyID: senderLeaf.PubKeyID,
			TokenID:  meta.TokenID,
			Balance:  massMigrations[i].Amount,
			Nonce:    models.MakeUint256(0),
		})
		if err != nil {
			return nil, nil, errors.WithStack(err)
		}
		hashes = append(hashes, *hash)
	}

	withdrawTree, err := merkletree.NewMerkleTree(hashes)
	if err != nil {
		return nil, nil, errors.WithStack(err)
	}

	return withdrawTree, meta, nil
}
