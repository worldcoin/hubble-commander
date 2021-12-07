package syncer

import (
	"github.com/Worldcoin/hubble-commander/commander/applier"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
)

func (c *TxsContext) SyncTxs(txs SyncedTxs, feeReceiverStateID uint32) (
	models.GenericTransactionArray,
	[]models.StateMerkleProof,
	error,
) {
	appliedTxs := models.NewGenericTransactionArray(txtype.FromBatchType(c.BatchType), 0, txs.Txs().Len())
	stateChangeProofs := c.Syncer.NewStateChangeProofs(txs.Txs().Len())
	combinedFee := models.NewUint256(0)

	tokenID, err := c.getCommitmentTokenID(txs.Txs(), feeReceiverStateID)
	if err != nil {
		return nil, nil, err
	}

	for i := 0; i < txs.Txs().Len(); i++ {
		synced, txError, appError := c.Syncer.ApplyTx(txs.TxAt(i), *tokenID)
		if appError != nil {
			return nil, nil, appError
		}
		stateChangeProofs = appendStateChangeProofs(stateChangeProofs, synced)
		if txError != nil {
			return nil, nil, NewDisputableErrorWithProofs(Transition, txError.Error(), stateChangeProofs)
		}

		appliedTxs = appliedTxs.AppendOne(synced.Tx)
		fee := synced.Tx.GetFee()
		*combinedFee = *combinedFee.Add(&fee)
	}

	stateProof, commitmentError, appError := c.Syncer.ApplyFee(feeReceiverStateID, tokenID, combinedFee)
	if appError != nil {
		return nil, nil, appError
	}
	stateChangeProofs = append(stateChangeProofs, *stateProof)
	if commitmentError != nil {
		return nil, nil, NewDisputableErrorWithProofs(Transition, commitmentError.Error(), stateChangeProofs)
	}

	return appliedTxs, stateChangeProofs, nil
}

func (c *TxsContext) getCommitmentTokenID(txs models.GenericTransactionArray, feeReceiverStateID uint32) (
	tokenID *models.Uint256,
	err error,
) {
	var leaf *models.StateLeaf

	if txs.Len() > 0 {
		leaf, err = c.storage.StateTree.LeafOrEmpty(txs.At(0).GetFromStateID())
	} else {
		leaf, err = c.storage.StateTree.LeafOrEmpty(feeReceiverStateID)
	}
	if err != nil {
		return nil, err
	}

	return &leaf.TokenID, nil
}

func appendStateChangeProofs(proofs []models.StateMerkleProof, syncedTx *applier.SyncedTxWithProofs) []models.StateMerkleProof {
	if syncedTx.ReceiverStateProof != nil {
		return append(proofs, syncedTx.SenderStateProof, *syncedTx.ReceiverStateProof)
	}
	return append(proofs, syncedTx.SenderStateProof)
}
