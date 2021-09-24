package syncer

import (
	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrInvalidDataLength = NewDisputableError(Transition, "invalid data length")
	ErrTooManyTxs        = NewDisputableError(Transition, "too many transactions in a commitment")
)

func (c *Context) syncTxCommitment(commitment *encoder.DecodedCommitment) (models.GenericTransactionArray, error) {
	if len(commitment.Transactions)%c.Syncer.TxLength() != 0 {
		return nil, ErrInvalidDataLength
	}

	syncedTxs, err := c.Syncer.DeserializeTxs(commitment.Transactions)
	if err != nil {
		return nil, err
	}

	if uint32(syncedTxs.Txs().Len()) > c.cfg.MaxTxsPerCommitment {
		return nil, ErrTooManyTxs
	}

	appliedTxs, stateProofs, err := c.SyncTxs(syncedTxs, commitment.FeeReceiver)
	if err != nil {
		return nil, err
	}
	syncedTxs.SetTxs(appliedTxs)

	err = c.verifyStateRoot(commitment.StateRoot, stateProofs)
	if err != nil {
		return nil, err
	}

	err = c.Syncer.SetPublicKeys(syncedTxs)
	if err != nil {
		return nil, err
	}
	if !c.cfg.DisableSignatures {
		err = c.verifyTxSignature(commitment, syncedTxs.Txs())
		if err != nil {
			return nil, err
		}
	}

	return syncedTxs.Txs(), nil
}

func (c *Context) verifyStateRoot(commitmentPostState common.Hash, proofs []models.StateMerkleProof) error {
	postStateRoot, err := c.storage.StateTree.Root()
	if err != nil {
		return err
	}
	if *postStateRoot != commitmentPostState {
		return NewDisputableErrorWithProofs(Transition, applier.ErrInvalidCommitmentStateRoot.Error(), proofs)
	}
	return nil
}
