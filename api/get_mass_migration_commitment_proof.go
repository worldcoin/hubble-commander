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
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var (
	ErrOnlyMassMigrationCommitmentsForProofing = fmt.Errorf(
		"mass migration commitment inclusion proof cannot be generated for different type of commitments",
	)

	APIErrOnlyMassMigrationCommitmentsForProofing = NewAPIError(
		50009,
		"mass migration commitment inclusion proof cannot be generated for different type of commitments",
	)

	APIErrCannotGenerateMMCommitmentProof = NewAPIError(50004, "mass migration commitment inclusion proof could not be generated")

	getMassMigrationCommitmentProofAPIErrors = map[error]*APIError{
		storage.AnyNotFoundError:                   APIErrCannotGenerateMMCommitmentProof,
		ErrOnlyMassMigrationCommitmentsForProofing: APIErrOnlyMassMigrationCommitmentsForProofing,
	}
)

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

	if batch.Type != batchtype.MassMigration {
		return nil, ErrOnlyMassMigrationCommitmentsForProofing
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

	commitment := commitments[commitmentID.IndexInBatch].ToMMCommitment()

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
			Meta:         dto.NewMassMigrationMeta(commitment.Meta),
			WithdrawRoot: commitment.WithdrawRoot,
			Transactions: serializedMassMigrations,
		},
	}, nil
}
