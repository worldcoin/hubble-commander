package api

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var getMassMigrationCommitmentInclusionProofAPIErrors = map[error]*APIError{
	storage.AnyNotFoundError: NewAPIError(50004, "mass migration commitment inclusion proof not found"),
}

func (a *API) GetMassMigrationCommitmentInclusionProof(batchID models.Uint256, commitmentIndex uint8) (*dto.MMCommitmentInclusionProof, error) {
	if !a.cfg.EnableProofMethods {
		return nil, errProofMethodsDisabled
	}
	commitmentInclusionProof, err := a.unsafeGetMassMigrationCommitmentInclusionProof(batchID, commitmentIndex)
	if err != nil {
		return nil, sanitizeError(err, getMassMigrationCommitmentInclusionProofAPIErrors)
	}
	return commitmentInclusionProof, nil
}

func (a *API) unsafeGetMassMigrationCommitmentInclusionProof(batchID models.Uint256, commitmentIndex uint8) (*dto.MMCommitmentInclusionProof, error) {
	batch, err := a.storage.GetBatch(batchID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	commitmentId := models.CommitmentID{
		BatchID:      batchID,
		IndexInBatch: commitmentIndex,
	}

	commitments, err := a.storage.GetCommitmentsByBatchID(batchID, batchtype.MassMigration)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	massMigrations, err := a.storage.GetMassMigrationsByCommitmentID(commitmentId)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	serializedMassMigrations, err := encoder.SerializeMassMigrations(massMigrations)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	commitment, err := a.storage.GetTxCommitment(&commitmentId)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	meta := &models.MassMigrationMeta{
		SpokeID:     0,
		TokenID:     models.MakeUint256(0),
		Amount:      models.MakeUint256(0),
		FeeReceiver: commitment.FeeReceiver,
	}

	hashes := make([]common.Hash, 0, len(massMigrations))

	for i := range massMigrations {
		meta.Amount = *meta.Amount.Add(&massMigrations[i].Amount)

		senderLeaf, err := a.storage.StateTree.Leaf(massMigrations[i].FromStateID)
		if err != nil {
			return nil, errors.WithStack(err)
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
			return nil, errors.WithStack(err)
		}
		hashes = append(hashes, *hash)
	}

	withdrawTree, err := merkletree.NewMerkleTree(hashes)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	leafHashes := make([]common.Hash, 0, len(commitments))
	for i := range commitments {
		leafHashes = append(leafHashes, commitments[i].LeafHash())
	}
	batchLeafTree, err := merkletree.NewMerkleTree(leafHashes)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	proofBase := &models.CommitmentInclusionProofBase{
		StateRoot: commitments[commitmentIndex].GetPostStateRoot(),
		Path: &models.MerklePath{
			Path:  uint32(commitmentIndex),
			Depth: batchLeafTree.Depth(),
		},
		Witness: batchLeafTree.GetWitness(uint32(commitmentIndex)),
	}

	return &dto.MMCommitmentInclusionProof{
		MMCommitmentInclusionProof: models.MMCommitmentInclusionProof{
			CommitmentInclusionProofBase: *proofBase,
			Body: &models.MMBody{
				AccountRoot:  *batch.AccountTreeRoot,
				Signature:    commitment.CombinedSignature,
				Meta:         meta,
				WithdrawRoot: withdrawTree.Root(),
				Transactions: serializedMassMigrations,
			},
		},
	}, nil
}
