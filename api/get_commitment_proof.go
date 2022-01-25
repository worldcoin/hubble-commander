package api

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/dto"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

var (
	ErrUnsupportedCommitmentTypeForProofing = fmt.Errorf(
		"commitment inclusion proof can only be generated for Transfer/Create2Transfer commitments",
	)

	APIErrProofMethodsDisabled              = NewAPIError(50000, "proof methods disabled")
	APIErrCannotGenerateCommitmentProof     = NewAPIError(50001, "commitment inclusion proof could not be generated")
	APIErrUnsupportedCommitmentTypeForProof = NewAPIError(
		50008,
		"commitment inclusion proof can only be generated for Transfer/Create2Transfer commitments",
	)

	getCommitmentProofAPIErrors = map[error]*APIError{
		storage.AnyNotFoundError:                APIErrCannotGenerateCommitmentProof,
		ErrUnsupportedCommitmentTypeForProofing: APIErrUnsupportedCommitmentTypeForProof,
	}
)

func (a *API) GetCommitmentProof(commitmentID models.CommitmentID) (*dto.CommitmentInclusionProof, error) {
	if !a.cfg.EnableProofMethods {
		return nil, APIErrProofMethodsDisabled
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
	if batch.Type != batchtype.Transfer && batch.Type != batchtype.Create2Transfer {
		//TODO: support MassMigration and Deposit types
		return nil, errors.WithStack(ErrUnsupportedCommitmentTypeForProofing)
	}

	commitments, err := a.storage.GetCommitmentsByBatchID(commitmentID.BatchID)
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

	commitment, err := a.storage.GetCommitment(&commitmentID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	transactions, err := a.getTransactionsForCommitment(commitment)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	body := &dto.CommitmentProofBody{
		AccountRoot:  *batch.AccountTreeRoot,
		Transactions: transactions,
	}

	if commitment.GetCommitmentBase().Type == batchtype.MassMigration {
		body.Signature = commitment.ToMMCommitment().CombinedSignature
		body.FeeReceiver = commitment.ToMMCommitment().Meta.FeeReceiver
	} else {
		body.Signature = commitment.ToTxCommitment().CombinedSignature
		body.FeeReceiver = commitment.ToTxCommitment().FeeReceiver
	}

	return &dto.CommitmentInclusionProof{
		CommitmentInclusionProofBase: dto.CommitmentInclusionProofBase{
			StateRoot: commitment.GetCommitmentBase().PostStateRoot,
			Path: &dto.MerklePath{
				Path:  path.Path,
				Depth: path.Depth,
			},
			Witness: tree.GetWitness(uint32(commitmentID.IndexInBatch)),
		},
		Body: body,
	}, nil
}
