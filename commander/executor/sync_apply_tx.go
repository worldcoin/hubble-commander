package executor

import "github.com/Worldcoin/hubble-commander/models"

func (c *SyncContext) ApplyTxs(txs SyncedTxs, feeReceiverStateID uint32) (
	models.GenericTransactionArray,
	[]models.StateMerkleProof,
	error,
) {
	appliedTxs := c.Syncer.NewTxArray(0, uint32(txs.Txs().Len()))
	stateChangeProofs := make([]models.StateMerkleProof, 0, 2*txs.Txs().Len()+1)
	combinedFee := models.NewUint256(0)

	tokenID, err := c.getCommitmentTokenID(txs.Txs(), feeReceiverStateID)
	if err != nil {
		return nil, nil, err
	}

	for i := 0; i < txs.Txs().Len(); i++ {
		synced, transferError, appError := c.Syncer.ApplyTx(txs.TxAt(i), *tokenID)
		if appError != nil {
			return nil, nil, appError
		}
		stateChangeProofs = append(
			stateChangeProofs,
			synced.SenderStateProof,
			synced.ReceiverStateProof,
		)
		if transferError != nil {
			return nil, nil, NewDisputableErrorWithProofs(Transition, transferError.Error(), stateChangeProofs)
		}

		appliedTxs = appliedTxs.AppendOne(synced.Tx)
		fee := synced.Tx.GetFee()
		*combinedFee = *combinedFee.Add(&fee)
	}

	stateProof, commitmentError, appError := c.ApplyFeeForSync(feeReceiverStateID, tokenID, combinedFee)
	if appError != nil {
		return nil, nil, appError
	}
	stateChangeProofs = append(stateChangeProofs, *stateProof)
	if commitmentError != nil {
		return nil, nil, NewDisputableErrorWithProofs(Transition, commitmentError.Error(), stateChangeProofs)
	}

	return appliedTxs, stateChangeProofs, nil
}

func (c *SyncContext) getCommitmentTokenID(txs models.GenericTransactionArray, feeReceiverStateID uint32) (
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
