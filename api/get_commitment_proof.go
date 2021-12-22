package api

import (
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var (
	errProofMethodsDisabled     = NewAPIError(50000, "proof methods disabled")
	getCommitmentProofAPIErrors = map[error]*APIError{
		storage.AnyNotFoundError: NewAPIError(50001, "commitment proof not found"),
	}
)

func (a *API) GetCommitmentProof(commitmentID models.CommitmentID) (*dto.CommitmentInclusionProof, error) {
	if !a.cfg.EnableProofMethods {
		return nil, errProofMethodsDisabled
	}
	commitmentProof, err := a.unsafeGetCommitmentProof(commitmentID)
	if err != nil {
		return nil, sanitizeError(err, getCommitmentProofAPIErrors)
	}
	return commitmentProof, nil
}

func (a *API) unsafeGetCommitmentProof(commitmentID models.CommitmentID) (*dto.CommitmentInclusionProof, error) {
	batch, err := a.storage.GetMinedBatch(commitmentID.BatchID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	commitments, err := a.storage.GetTxCommitmentsByBatchID(commitmentID.BatchID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	leafHashes := make([]common.Hash, 0, len(commitments))
	for i := range commitments {
		leafHashes = append(leafHashes, commitments[i].LeafHash())
	}
	tree, err := merkletree.NewMerkleTree(leafHashes)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	path := &dto.MerklePath{
		Path:  uint32(commitmentID.IndexInBatch),
		Depth: tree.Depth(),
	}

	commitment, err := a.storage.GetTxCommitment(&commitmentID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	transactions, err := a.getTransactionsForCommitment(commitment)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &dto.CommitmentInclusionProof{
		StateRoot: commitment.PostStateRoot,
		Body: &dto.CommitmentProofBody{
			AccountRoot:  *batch.AccountTreeRoot,
			Signature:    commitment.CombinedSignature,
			FeeReceiver:  commitment.FeeReceiver,
			Transactions: transactions,
		},
		Path:    path,
		Witness: tree.GetWitness(uint32(commitmentID.IndexInBatch)),
	}, nil
}
