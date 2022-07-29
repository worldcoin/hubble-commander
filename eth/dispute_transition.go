package eth

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
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
	batchHash *common.Hash,
	previous *models.CommitmentInclusionProof,
	target *models.TransferCommitmentInclusionProof,
	proofs []models.StateMerkleProof,
) error {
	transaction, err := c.rollup().
		WithGasLimit(*c.config.TransitionDisputeGasLimit).
		DisputeTransitionTransfer(
			batchID.ToBig(),
			*commitmentProofToCalldata(previous),
			*transferProofToCalldata(target),
			stateMerkleProofsToCalldata(proofs),
		)
	if err != nil {
		return handleDisputeTransitionError(err)
	}

	err = c.waitForDispute(batchID, batchHash, transaction)
	if err == ErrBatchAlreadyDisputed || err == ErrRollbackInProcess {
		log.Info(err)
		return nil
	}
	return err
}

func (c *Client) DisputeTransitionCreate2Transfer(
	batchID *models.Uint256,
	batchHash *common.Hash,
	previous *models.CommitmentInclusionProof,
	target *models.TransferCommitmentInclusionProof,
	proofs []models.StateMerkleProof,
) error {
	transaction, err := c.rollup().
		WithGasLimit(*c.config.TransitionDisputeGasLimit).
		DisputeTransitionCreate2Transfer(
			batchID.ToBig(),
			*commitmentProofToCalldata(previous),
			*transferProofToCalldata(target),
			stateMerkleProofsToCalldata(proofs),
		)
	if err != nil {
		return handleDisputeTransitionError(err)
	}

	err = c.waitForDispute(batchID, batchHash, transaction)
	if err == ErrBatchAlreadyDisputed || err == ErrRollbackInProcess {
		log.Info(err)
		return nil
	}
	return err
}

func (c *Client) DisputeTransitionMassMigration(
	batchID *models.Uint256,
	batchHash *common.Hash,
	previous *models.CommitmentInclusionProof,
	target *models.MMCommitmentInclusionProof,
	proofs []models.StateMerkleProof,
) error {
	transaction, err := c.rollup().
		WithGasLimit(*c.config.TransitionDisputeGasLimit).
		DisputeTransitionMassMigration(
			batchID.ToBig(),
			*commitmentProofToCalldata(previous),
			*massMigrationProofToCalldata(target),
			stateMerkleProofsToCalldata(proofs),
		)
	if err != nil {
		return handleDisputeTransitionError(err)
	}

	err = c.waitForDispute(batchID, batchHash, transaction)
	if err == ErrBatchAlreadyDisputed || err == ErrRollbackInProcess {
		log.Info(err)
		return nil
	}
	return err
}

func (c *Client) waitForDispute(batchID *models.Uint256, batchHash *common.Hash, tx *types.Transaction) error {
	receipt, err := c.WaitToBeMined(tx)
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
	err = c.isBatchAlreadyDisputed(batchID, batchHash)
	if err != nil {
		return err
	}
	err = c.GetRevertMessage(tx, receipt)
	if err != nil {
		// one day in the far future it would be nice to stitch the solidity
		// stack/call trace onto the other end of this one
		result := NewDisputeTxRevertedError(batchID.Uint64(), err.Error())
		return errors.WithStack(result)
	}
	return errors.WithStack(NewUnknownDisputeTxRevertedError(batchID.Uint64()))
}

func (c *Client) isBatchAlreadyDisputed(batchID *models.Uint256, batchHash *common.Hash) error {
	contractBatch, err := c.GetContractBatch(batchID)
	if err != nil {
		return err
	}
	if contractBatch.Hash != *batchHash {
		return ErrBatchAlreadyDisputed
	}
	return nil
}

func (c *Client) isBatchDuringDispute(batchID *models.Uint256) error {
	invalidBatchID, err := c.GetInvalidBatchID()
	if err != nil {
		return err
	}
	if invalidBatchID != nil && batchID.Cmp(invalidBatchID) >= 0 {
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

func commitmentProofToCalldata(proof *models.CommitmentInclusionProof) *rollup.TypesCommitmentInclusionProof {
	return &rollup.TypesCommitmentInclusionProof{
		Commitment: rollup.TypesCommitment{
			StateRoot: proof.StateRoot,
			BodyRoot:  proof.BodyRoot,
		},
		Path:    new(big.Int).SetUint64(uint64(proof.Path.Path)),
		Witness: proof.Witness.Bytes(),
	}
}

func transferProofToCalldata(proof *models.TransferCommitmentInclusionProof) *rollup.TypesTransferCommitmentInclusionProof {
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

func massMigrationProofToCalldata(proof *models.MMCommitmentInclusionProof) *rollup.TypesMMCommitmentInclusionProof {
	return &rollup.TypesMMCommitmentInclusionProof{
		Commitment: rollup.TypesMassMigrationCommitment{
			StateRoot: proof.StateRoot,
			Body: rollup.TypesMassMigrationBody{
				AccountRoot:  proof.Body.AccountRoot,
				Signature:    proof.Body.Signature.BigInts(),
				SpokeID:      new(big.Int).SetUint64(uint64(proof.Body.Meta.SpokeID)),
				WithdrawRoot: proof.Body.WithdrawRoot,
				TokenID:      proof.Body.Meta.TokenID.ToBig(),
				Amount:       proof.Body.Meta.Amount.ToBig(),
				FeeReceiver:  new(big.Int).SetUint64(uint64(proof.Body.Meta.FeeReceiver)),
				Txs:          proof.Body.Transactions,
			},
		},
		Path:    new(big.Int).SetUint64(uint64(proof.Path.Path)),
		Witness: proof.Witness.Bytes(),
	}
}

func stateMerkleProofsToCalldata(proofs []models.StateMerkleProof) []rollup.TypesStateMerkleProof {
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
	invalidBatchID := models.NewUint256FromBig(*batchMarker)
	if invalidBatchID.IsZero() {
		return nil, nil
	}
	return invalidBatchID, nil
}
