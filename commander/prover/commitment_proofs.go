package prover

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/utils/merkletree"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func (c *Context) PreviousCommitmentInclusionProof(
	batch *eth.DecodedTxBatch,
	previousCommitmentIndex int,
) (*models.CommitmentInclusionProof, error) {
	if previousCommitmentIndex == -1 {
		return c.PreviousBatchCommitmentInclusionProof(batch.ID)
	}

	leafHashes := make([]common.Hash, 0, len(batch.Commitments))
	for i := range batch.Commitments {
		leafHashes = append(leafHashes, batch.Commitments[i].LeafHash(batch.AccountTreeRoot))
	}

	return createCommitmentInclusionProof(
		leafHashes,
		uint32(previousCommitmentIndex),
		batch.Commitments[previousCommitmentIndex].GetStateRoot(),
		*batch.Commitments[previousCommitmentIndex].BodyHash(batch.AccountTreeRoot),
	)
}

func (c *Context) PreviousBatchCommitmentInclusionProof(
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

	commitments, err := c.storage.GetCommitmentsByBatchID(previousBatch.ID)
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
		commitments[previousCommitmentIndex].GetPostStateRoot(),
		*commitments[previousCommitmentIndex].GetBodyHash(),
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
		CommitmentInclusionProofBase: models.CommitmentInclusionProofBase{
			StateRoot: commitmentStateRoot,
		},
		BodyRoot: commitmentBodyRoot,
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

func (c *Context) TargetTransferCommitmentInclusionProof(
	batch *eth.DecodedTxBatch,
	commitmentIndex uint32,
) (*models.TransferCommitmentInclusionProof, error) {
	proofBase, err := c.targetCommitmentInclusionProofBase(batch, commitmentIndex)
	if err != nil {
		return nil, err
	}

	commitment := batch.Commitments[commitmentIndex].ToDecodedCommitment()
	return &models.TransferCommitmentInclusionProof{
		CommitmentInclusionProofBase: *proofBase,
		Body: &models.TransferBody{
			AccountRoot:  batch.AccountTreeRoot,
			Signature:    commitment.CombinedSignature,
			FeeReceiver:  commitment.FeeReceiver,
			Transactions: commitment.Transactions,
		},
	}, nil
}

func (c *Context) TargetMMCommitmentInclusionProof(
	batch *eth.DecodedTxBatch,
	commitmentIndex uint32,
) (*models.MMCommitmentInclusionProof, error) {
	proofBase, err := c.targetCommitmentInclusionProofBase(batch, commitmentIndex)
	if err != nil {
		return nil, err
	}

	commitment := batch.Commitments[commitmentIndex].(*encoder.DecodedMMCommitment)
	return &models.MMCommitmentInclusionProof{
		CommitmentInclusionProofBase: *proofBase,
		Body: &models.MMBody{
			AccountRoot:  batch.AccountTreeRoot,
			Signature:    commitment.CombinedSignature,
			Meta:         commitment.Meta,
			WithdrawRoot: commitment.WithdrawRoot,
			Transactions: commitment.Transactions,
		},
	}, nil
}

func (c *Context) targetCommitmentInclusionProofBase(
	batch *eth.DecodedTxBatch,
	commitmentIndex uint32,
) (*models.CommitmentInclusionProofBase, error) {
	leafHashes := make([]common.Hash, 0, len(batch.Commitments))
	for i := range batch.Commitments {
		leafHashes = append(leafHashes, batch.Commitments[i].LeafHash(batch.AccountTreeRoot))
	}
	tree, err := merkletree.NewMerkleTree(leafHashes)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &models.CommitmentInclusionProofBase{
		StateRoot: batch.Commitments[commitmentIndex].GetStateRoot(),
		Path: &models.MerklePath{
			Path:  commitmentIndex,
			Depth: tree.Depth(),
		},
		Witness: tree.GetWitness(commitmentIndex),
	}, nil
}
