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

var ErrWaitForRollbackTimeout = errors.New("waitForRollbackToFinish: timeout")

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
	subscription, err := c.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink) // TODO-dis query receipts instead of subscribing events
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	transaction, err := c.rollup().
		DisputeTransitionCreate2Transfer(
			batchID.ToBig(),
			*CommitmentProofToCalldata(previous),
			*TransferProofToCalldata(target),
			StateMerkleProofsToCalldata(proofs),
		)
	// TODO-dis handle "Already successfully disputed. Roll back in process" error
	// TODO-dis handle error caused by reverted transaction (someone else already disputed, check against nextBatchID)
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
				transactionHash, err = c.KeepRollingBack()
				if err != nil {
					return err
				}
			}
		case <-time.After(*c.config.TxTimeout):
			return ErrWaitForRollbackTimeout
		}
	}
}

func (c *Client) KeepRollingBack() (common.Hash, error) {
	transaction, err := c.rollup().KeepRollingBack()
	// TODO-dis handle "BatchManager: Is not rolling back" error
	// TODO-dis handle error caused by reverted transaction (already rolled back, check against nextBatchID)
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

func (c *Client) GetInvalidBatchID() (*models.Uint256, error) {
	batchMarker, err := c.Rollup.InvalidBatchMarker(nil)
	if err != nil {
		return nil, err
	}
	return models.NewUint256FromBig(*batchMarker), err
}
