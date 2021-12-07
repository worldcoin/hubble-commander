package applier

import (
	"github.com/Worldcoin/hubble-commander/models"
)

func (a *Applier) ApplyMassMigration(tx models.GenericTransaction, commitmentTokenID models.Uint256) (
	applyResult ApplySingleTxResult, txError, appError error,
) {
	senderLeaf, appErr := a.storage.StateTree.Leaf(tx.GetFromStateID())
	if appErr != nil {
		return nil, nil, appErr
	}
	appErr = a.validateSenderTokenID(senderLeaf, commitmentTokenID)
	if appErr != nil {
		return nil, nil, appErr
	}
	if txErr := validateTxNonce(&senderLeaf.UserState, tx.GetNonce()); txErr != nil {
		return nil, txErr, nil
	}

	newSenderState, txErr := calculateSenderStateAfterTx(senderLeaf.UserState, tx)
	if txErr != nil {
		return nil, txErr, nil
	}

	_, appErr = a.storage.StateTree.Set(tx.GetFromStateID(), newSenderState)
	if appErr != nil {
		return nil, nil, appErr
	}

	return &ApplySingleMassMigrationResult{tx: tx.ToMassMigration()}, nil, nil
}

func (a *Applier) ApplyMassMigrationForSync(
	tx models.GenericTransaction,
	commitmentTokenID models.Uint256,
) (synced *SyncedTxWithProofs, txErr, appErr error) {
	senderLeaf, appErr := a.storage.StateTree.Leaf(tx.GetFromStateID())
	if appErr != nil {
		return nil, nil, appErr
	}

	synced = NewSyncedTxWithSenderProof(tx.Copy(), &senderLeaf.UserState)

	newSenderState, txErr := calculateSenderStateAfterTx(senderLeaf.UserState, tx)
	if txErr != nil {
		return a.fillSenderWitness(synced, txErr)
	}

	senderWitness, appErr := a.storage.StateTree.Set(tx.GetFromStateID(), newSenderState)
	if appErr != nil {
		return nil, nil, appErr
	}
	synced.SenderStateProof.Witness = senderWitness

	txErr = a.validateSenderTokenID(senderLeaf, commitmentTokenID)
	if txErr != nil {
		return synced, txErr, nil
	}

	synced.Tx.SetNonce(senderLeaf.Nonce)

	return synced, nil, nil
}
