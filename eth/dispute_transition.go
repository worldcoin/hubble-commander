package eth

import (
	"math/big"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Client) DisputeTransitionTransfer(
	batchID *models.Uint256,
	previous *models.CommitmentInclusionProof,
	target *models.TransferCommitmentInclusionProof,
	proofs []models.StateMerkleProof,
) (*types.Transaction, error) {
	return c.Rollup.DisputeTransitionTransfer(
		c.transactOpts(c.config.stakeAmount.ToBig(), 0),
		batchID.ToBig(),
		*CommitmentProofToCalldata(previous),
		*TransferProofToCalldata(target),
		StateMerkleProofsToCalldata(proofs),
	)
}

func (c *Client) DisputeTransitionCreate2Transfer(
	batchID *models.Uint256,
	previous *models.CommitmentInclusionProof,
	target *models.TransferCommitmentInclusionProof,
	proofs []models.StateMerkleProof,
) (*types.Transaction, error) {
	return c.Rollup.DisputeTransitionCreate2Transfer(
		c.transactOpts(c.config.stakeAmount.ToBig(), 0),
		batchID.ToBig(),
		*CommitmentProofToCalldata(previous),
		*TransferProofToCalldata(target),
		StateMerkleProofsToCalldata(proofs),
	)
}

func (c *Client) GetRollbackStatus(
	transactionHash common.Hash,
) (*rollup.RollupRollbackStatus, error) {
	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := c.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	if err != nil {
		return nil, err
	}
	defer subscription.Unsubscribe()

	for {
		select {
		case err = <-subscription.Err():
			return nil, errors.WithStack(err)
		case rollbackStatus := <-sink:
			if rollbackStatus.Raw.TxHash == transactionHash {
				return rollbackStatus, nil
			}
		case <-time.After(*c.config.txTimeout):
			return nil, errors.New("getRollbackStatus: timeout")
		}
	}
}

func CommitmentProofToCalldata(proof *models.CommitmentInclusionProof) *rollup.TypesCommitmentInclusionProof {
	return &rollup.TypesCommitmentInclusionProof{
		Commitment: rollup.TypesCommitment{
			StateRoot: proof.StateRoot,
			BodyRoot:  proof.BodyRoot,
		},
		Path:    new(big.Int).SetUint64(uint64(proof.Path.Path)),
		Witness: proof.Witnesses.Bytes(),
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
		Witness: proof.Witnesses.Bytes(),
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
			TokenID:  proof.UserState.TokenIndex.ToBig(),
			Balance:  proof.UserState.Balance.ToBig(),
			Nonce:    proof.UserState.Nonce.ToBig(),
		},
		Witness: proof.Witnesses.Bytes(),
	}
}
