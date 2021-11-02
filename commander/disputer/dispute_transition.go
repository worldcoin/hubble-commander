package disputer

import (
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func (c *Context) previousCommitmentInclusionProof(
	batch *eth.DecodedBatch,
	previousCommitmentIndex int,
) (*models.CommitmentInclusionProof, error) {
	if previousCommitmentIndex == -1 {
		return c.previousBatchCommitmentInclusionProof(batch.ID)
	}

	leafHashes := make([]common.Hash, 0, len(batch.Commitments))
	for i := range batch.Commitments {
		leafHashes = append(leafHashes, batch.Commitments[i].LeafHash(*batch.AccountTreeRoot))
	}

	return createCommitmentInclusionProof(
		leafHashes,
		uint32(previousCommitmentIndex),
		batch.Commitments[previousCommitmentIndex].StateRoot,
		batch.Commitments[previousCommitmentIndex].BodyHash(*batch.AccountTreeRoot),
	)
}

func (c *Context) previousBatchCommitmentInclusionProof(
	currentBatchID models.Uint256,
) (*models.CommitmentInclusionProof, error) {
	previousBatchID := currentBatchID.SubN(1)
	if previousBatchID.IsZero() {
		return c.genesisBatchCommitmentInclusionProof()
	}

	previousBatch, err := c.storage.GetBatch(*previousBatchID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	commitments, err := c.storage.GetTxCommitmentsByBatchID(previousBatch.ID)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	leafHashes := make([]common.Hash, 0, len(commitments))
	for i := range commitments {
		leafHashes = append(leafHashes, commitments[i].LeafHash())
	}

	previousCommitmentIndex := len(commitments) - 1
	return createCommitmentInclusionProof(
		leafHashes,
		uint32(previousCommitmentIndex),
		commitments[previousCommitmentIndex].PostStateRoot,
		*commitments[previousCommitmentIndex].BodyHash,
	)
}

func (c *Context) genesisBatchCommitmentInclusionProof() (*models.CommitmentInclusionProof, error) {
	previousBatch, err := c.storage.GetBatch(models.MakeUint256(0))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return createCommitmentInclusionProof(
		[]common.Hash{*previousBatch.PrevStateRoot},
		0,
		*previousBatch.PrevStateRoot,
		merkletree.GetZeroHash(0),
	)
}

func createCommitmentInclusionProof(
	leafHashes []common.Hash,
	commitmentIndex uint32,
	commitmentStateRoot, commitmentBodyRoot common.Hash,
) (*models.CommitmentInclusionProof, error) {
	proof := models.CommitmentInclusionProof{
		StateRoot: commitmentStateRoot,
		BodyRoot:  commitmentBodyRoot,
	}

	tree, err := merkletree.NewMerkleTree(leafHashes)
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
	leafHashes := make([]common.Hash, 0, len(batch.Commitments))
	for i := range batch.Commitments {
		leafHashes = append(leafHashes, batch.Commitments[i].LeafHash(*batch.AccountTreeRoot))
	}
	tree, err := merkletree.NewMerkleTree(leafHashes)
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

func (c *Context) DisputeTransition(
	batch *eth.DecodedBatch,
	commitmentIndex int,
	merkleProofs []models.StateMerkleProof,
) error {
	previousCommitmentProof, err := c.previousCommitmentInclusionProof(batch, commitmentIndex-1)
	if err != nil {
		return err
	}
	targetCommitmentProof, err := targetCommitmentInclusionProof(batch, uint32(commitmentIndex))
	if err != nil {
		return err
	}

	if batch.Type == batchtype.Transfer {
		err = c.client.DisputeTransitionTransfer(&batch.ID, batch.Hash, previousCommitmentProof, targetCommitmentProof, merkleProofs)
	} else {
		err = c.client.DisputeTransitionCreate2Transfer(&batch.ID, batch.Hash, previousCommitmentProof, targetCommitmentProof, merkleProofs)
	}
	return err
}
