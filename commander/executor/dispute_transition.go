package executor

import (
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func (t *TransactionExecutor) previousCommitmentInclusionProof(
	batch *eth.DecodedBatch,
	previousCommitmentIndex int,
) (*models.CommitmentInclusionProof, error) {
	if previousCommitmentIndex == -1 {
		return t.previousBatchCommitmentInclusionProof(batch.ID)
	}

	bodyHashes := make([]common.Hash, 0, len(batch.Commitments))
	for i := range batch.Commitments {
		bodyHashes = append(bodyHashes, batch.Commitments[i].BodyHash(*batch.AccountTreeRoot))
	}

	return createCommitmentInclusionProof(bodyHashes, uint32(previousCommitmentIndex), batch.Commitments[previousCommitmentIndex].StateRoot)
}

func (t *TransactionExecutor) previousBatchCommitmentInclusionProof(
	currentBatchID models.Uint256,
) (*models.CommitmentInclusionProof, error) {
	previousBatch, err := t.storage.GetBatch(*currentBatchID.SubN(1))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	commitments, err := t.storage.GetCommitmentsByBatchID(previousBatch.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	bodyHashes := make([]common.Hash, 0, len(commitments))
	for i := range commitments {
		bodyHashes = append(bodyHashes, commitments[i].BodyHash(*previousBatch.AccountTreeRoot))
	}

	previousCommitmentIndex := len(commitments) - 1
	return createCommitmentInclusionProof(bodyHashes, uint32(previousCommitmentIndex), commitments[previousCommitmentIndex].PostStateRoot)
}

func createCommitmentInclusionProof(
	bodyHashes []common.Hash,
	commitmentIndex uint32,
	commitmentStateRoot common.Hash,
) (*models.CommitmentInclusionProof, error) {
	proof := models.CommitmentInclusionProof{
		StateRoot: commitmentStateRoot,
		BodyRoot:  bodyHashes[commitmentIndex],
	}

	tree, err := merkletree.NewMerkleTree(bodyHashes)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	proof.Path = &models.MerklePath{
		Path:  commitmentIndex,
		Depth: tree.Depth(),
	}

	proof.Witness = tree.GetWitness(commitmentIndex)
	return &proof, nil
}

func targetCommitmentInclusionProof(
	batch *eth.DecodedBatch,
	commitmentIndex uint32,
) (*models.TransferCommitmentInclusionProof, error) {
	bodyHashes := make([]common.Hash, 0, len(batch.Commitments))
	for i := range batch.Commitments {
		bodyHashes = append(bodyHashes, batch.Commitments[i].BodyHash(*batch.AccountTreeRoot))
	}
	tree, err := merkletree.NewMerkleTree(bodyHashes)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	path := &models.MerklePath{
		Path:  commitmentIndex,
		Depth: tree.Depth(),
	}

	commitment := batch.Commitments[commitmentIndex]
	return &models.TransferCommitmentInclusionProof{
		StateRoot: commitment.StateRoot,
		Body: &models.TransferBody{
			AccountRoot:  *batch.AccountTreeRoot,
			Signature:    commitment.CombinedSignature,
			FeeReceiver:  commitment.FeeReceiver,
			Transactions: commitment.Transactions,
		},
		Path:    path,
		Witness: tree.GetWitness(commitmentIndex),
	}, nil
}
