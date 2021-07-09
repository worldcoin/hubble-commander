package eth

import (
	"math/big"
	"time"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/eth/deployer"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"
	"github.com/pkg/errors"
)

func (c *Client) DisputeTransitionTransfer(
	batchID *models.Uint256,
	previous *models.CommitmentInclusionProof,
	target *models.TransferCommitmentInclusionProof,
	proofs []models.StateMerkleProof,
) error {
	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := c.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	transaction, err := c.rollup().
		WithGasLimit(0).
		DisputeTransitionTransfer(
			batchID.ToBig(),
			*CommitmentProofToCalldata(previous),
			*TransferProofToCalldata(target),
			StateMerkleProofsToCalldata(proofs),
		)
	if err != nil {
		return err
	}
	return c.waitForRollbackToFinish(sink, subscription, transaction.Hash())
}

func (c *Client) DisputeTransitionCreate2Transfer(
	batchID *models.Uint256,
	previous *models.CommitmentInclusionProof,
	target *models.TransferCommitmentInclusionProof,
	proofs []models.StateMerkleProof,
) error {
	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := c.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	transaction, err := c.rollup().
		WithGasLimit(0).
		DisputeTransitionCreate2Transfer(
			batchID.ToBig(),
			*CommitmentProofToCalldata(previous),
			*TransferProofToCalldata(target),
			StateMerkleProofsToCalldata(proofs),
		)
	if err != nil {
		return err
	}
	return c.waitForRollbackToFinish(sink, subscription, transaction.Hash())
}

func (c *Client) waitForRollbackToFinish(
	sink chan *rollup.RollupRollbackStatus,
	subscription event.Subscription,
	transactionHash common.Hash,
) (err error) {
	for {
		select {
		case err = <-subscription.Err():
			return errors.WithStack(err)
		case rollbackStatus := <-sink:
			if rollbackStatus.Raw.TxHash == transactionHash {
				if rollbackStatus.Completed {
					return nil
				}
				transactionHash, err = c.keepRollingBack()
				if err != nil {
					return err
				}
			}
		case <-time.After(*c.config.txTimeout):
			return errors.New("waitForRollbackToFinish: timeout")
		}
	}
}

func (c *Client) keepRollingBack() (common.Hash, error) {
	transaction, err := c.Rollup.KeepRollingBack(c.transactOpts(nil, 8_000_000))
	if err != nil {
		return common.Hash{}, err
	}
	receipt, err := deployer.WaitToBeMined(c.ChainConnection.GetBackend(), transaction)
	if err != nil {
		return common.Hash{}, err
	}
	return receipt.TxHash, nil
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
