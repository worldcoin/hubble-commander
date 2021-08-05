package eth

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	msgBatchAlreadyDisputed   = "Already successfully disputed. Roll back in process"
	msgTransitionMissingBatch = "Target commitment is absent in the batch"
	msgSignatureMissingBatch  = "Commitment not present in batch"
)

var (
	ErrBatchAlreadyDisputed = errors.New("batch already disputed")
	ErrRollbackInProcess    = errors.New("rollback in process")
)

func (c *Client) DisputeTransitionTransfer(
	batchID *models.Uint256,
	previous *models.CommitmentInclusionProof,
	target *models.TransferCommitmentInclusionProof,
	proofs []models.StateMerkleProof,
) error {
	transaction, err := c.rollup().
		DisputeTransitionTransfer(
			batchID.ToBig(),
			*CommitmentProofToCalldata(previous),
			*TransferProofToCalldata(target),
			StateMerkleProofsToCalldata(proofs),
		)
	if err != nil {
		return handleDisputeTransitionError(err)
	}

	err = c.waitForDispute(batchID, transaction)
	if err == ErrBatchAlreadyDisputed || err == ErrRollbackInProcess {
		log.Info(err)
		return nil
	}
	return err
}

func (c *Client) DisputeTransitionCreate2Transfer(
	batchID *models.Uint256,
	previous *models.CommitmentInclusionProof,
	target *models.TransferCommitmentInclusionProof,
	proofs []models.StateMerkleProof,
) error {
	transaction, err := c.rollup().
		DisputeTransitionCreate2Transfer(
			batchID.ToBig(),
			*CommitmentProofToCalldata(previous),
			*TransferProofToCalldata(target),
			StateMerkleProofsToCalldata(proofs),
		)
	if err != nil {
		return handleDisputeTransitionError(err)
	}

	err = c.waitForDispute(batchID, transaction)
	if err == ErrBatchAlreadyDisputed || err == ErrRollbackInProcess {
		log.Info(err)
		return nil
	}
	return err
}

func (c *Client) waitForDispute(batchID *models.Uint256, tx *types.Transaction) error {
	receipt, err := deployer.WaitToBeMined(c.ChainConnection.GetBackend(), tx)
	if err != nil {
		return err
	}
	if receipt.Status == types.ReceiptStatusSuccessful {
		return nil
	}

	err = c.isBatchDuringDispute(batchID)
	if err != nil {
		return err
	}
	err = c.isBatchAlreadyDisputed(batchID)
	if err != nil {
		return err
	}
	return errors.Errorf("dispute of batch #%d failed", batchID.Uint64())
}

func (c *Client) isBatchAlreadyDisputed(batchID *models.Uint256) error {
	nextBatchID, err := c.Rollup.NextBatchID(nil)
	if err != nil {
		return err
	}

	if batchID.CmpN(nextBatchID.Uint64()) < 0 {
		return ErrBatchAlreadyDisputed
	}
	return nil
}

func (c *Client) isBatchDuringDispute(batchID *models.Uint256) error {
	invalidBatchID, err := c.GetInvalidBatchID()
	if err != nil {
		return err
	}
	if !invalidBatchID.IsZero() && batchID.Cmp(invalidBatchID) >= 0 {
		return ErrRollbackInProcess
	}
	return nil
}

func handleDisputeTransitionError(err error) error {
	errMsg := getGasEstimateErrorMessage(err)
	if errMsg == msgTransitionMissingBatch || errMsg == msgBatchAlreadyDisputed {
		log.Info(err.Error())
		return nil
	}
	return err
}

func CommitmentProofToCalldata(proof *models.CommitmentInclusionProof) *rollup.TypesCommitmentInclusionProof {
	return &rollup.TypesCommitmentInclusionProof{
		Commitment: rollup.TypesCommitment{
			StateRoot: proof.StateRoot,
			BodyRoot:  proof.BodyRoot,
		},
		Path:    new(big.Int).SetUint64(uint64(proof.Path.Path)),
		Witness: proof.Witness.Bytes(),
	}
}

func TransferProofToCalldata(proof *models.TransferCommitmentInclusionProof) *rollup.TypesTransferCommitmentInclusionProof {
	return &rollup.TypesTransferCommitmentInclusionProof{
		Commitment: rollup.TypesTransferCommitment{
			StateRoot: proof.StateRoot,
			Body: rollup.TypesTransferBody{
				AccountRoot: proof.Body.AccountRoot,
				Signature:   proof.Body.Signature.BigInts(),
				FeeReceiver: new(big.Int).SetUint64(uint64(proof.Body.FeeReceiver)),
				Txs:         proof.Body.Transactions,
			},
		},
		Path:    new(big.Int).SetUint64(uint64(proof.Path.Path)),
		Witness: proof.Witness.Bytes(),
	}
}

func StateMerkleProofsToCalldata(proofs []models.StateMerkleProof) []rollup.TypesStateMerkleProof {
	result := make([]rollup.TypesStateMerkleProof, 0, len(proofs))
	for i := range proofs {
		result = append(result, *stateMerkleProofToCalldata(&proofs[i]))
	}
	return result
}

func stateMerkleProofToCalldata(proof *models.StateMerkleProof) *rollup.TypesStateMerkleProof {
	return &rollup.TypesStateMerkleProof{
		State: rollup.TypesUserState{
			PubkeyID: new(big.Int).SetUint64(uint64(proof.UserState.PubKeyID)),
			TokenID:  proof.UserState.TokenID.ToBig(),
			Balance:  proof.UserState.Balance.ToBig(),
			Nonce:    proof.UserState.Nonce.ToBig(),
		},
		Witness: proof.Witness.Bytes(),
	}
}

func (c *Client) GetInvalidBatchID() (*models.Uint256, error) {
	batchMarker, err := c.Rollup.InvalidBatchMarker(nil)
	if err != nil {
		return nil, err
	}
	return models.NewUint256FromBig(*batchMarker), err
}
