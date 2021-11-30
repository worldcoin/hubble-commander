package applier

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/pkg/errors"
)

var (
	ErrNonceTooLow            = fmt.Errorf("nonce too low")
	ErrNonceTooHigh           = fmt.Errorf("nonce too high")
	ErrBalanceTooLow          = fmt.Errorf("not enough balance")
	ErrInvalidSenderTokenID   = fmt.Errorf("invalid sender token ID")
	ErrInvalidReceiverTokenID = fmt.Errorf("invalid receiver token ID")
	ErrInvalidTokenAmount     = fmt.Errorf("amount cannot be equal to 0")
)

func (a *Applier) ApplyTx(
	tx models.GenericTransaction,
	receiverLeaf *models.StateLeaf,
	commitmentTokenID models.Uint256,
) (txError, appError error) {
	newSenderState, newReceiverState, txError, appError := a.validateAndCalculateStateAfterTx(tx, receiverLeaf, commitmentTokenID)
	if txError != nil || appError != nil {
		return txError, appError
	}

	_, appError = a.storage.StateTree.Set(tx.GetFromStateID(), newSenderState)
	if appError != nil {
		return nil, appError
	}
	_, appError = a.storage.StateTree.Set(receiverLeaf.StateID, newReceiverState)
	if appError != nil {
		return nil, appError
	}

	return nil, nil
}

func (a *Applier) validateAndCalculateStateAfterTx(
	tx models.GenericTransaction,
	receiverLeaf *models.StateLeaf,
	commitmentTokenID models.Uint256,
) (newSenderState, newReceiverState *models.UserState, txError, appError error) {
	senderLeaf, appError := a.storage.StateTree.Leaf(tx.GetFromStateID())
	if appError != nil {
		return nil, nil, nil, appError
	}

	appError = a.validateSenderTokenID(senderLeaf, commitmentTokenID)
	if appError != nil {
		return nil, nil, nil, appError
	}

	appError = a.validateReceiverTokenID(receiverLeaf, commitmentTokenID)
	if appError != nil {
		return nil, nil, nil, appError
	}

	if txErr := validateTxNonce(&senderLeaf.UserState, tx.GetNonce()); txErr != nil {
		return nil, nil, txErr, nil
	}

	newSenderState, newReceiverState, txErr := calculateStateAfterTx(senderLeaf.UserState, receiverLeaf.UserState, tx)
	if txErr != nil {
		return nil, nil, txErr, nil
	}

	return newSenderState, newReceiverState, nil, nil
}

func (a *Applier) applyTxForSync(
	tx models.GenericTransaction,
	receiverLeaf *models.StateLeaf,
	commitmentTokenID models.Uint256,
) (synced *SyncedGenericTransaction, txError, appError error) {
	senderLeaf, appError := a.storage.StateTree.LeafOrEmpty(tx.GetFromStateID())
	if appError != nil {
		return nil, nil, appError
	}

	synced = NewPartialSyncedGenericTransaction(tx.Copy(), &senderLeaf.UserState, &receiverLeaf.UserState)

	newSenderState, newReceiverState, txErr := calculateStateAfterTx(senderLeaf.UserState, receiverLeaf.UserState, tx)
	if txErr != nil {
		return a.fillSenderWitness(synced, txErr)
	}

	senderWitness, appError := a.storage.StateTree.Set(senderLeaf.StateID, newSenderState)
	if appError != nil {
		return nil, nil, appError
	}
	synced.SenderStateProof.Witness = senderWitness

	if txErr = a.validateSenderTokenID(senderLeaf, commitmentTokenID); txErr != nil {
		return synced, txErr, nil
	}

	receiverWitness, appError := a.storage.StateTree.Set(receiverLeaf.StateID, newReceiverState)
	if appError != nil {
		return nil, nil, appError
	}
	if senderLeaf.StateID == receiverLeaf.StateID {
		synced.ReceiverStateProof.UserState = newSenderState
	}
	synced.ReceiverStateProof.Witness = receiverWitness

	if txErr = a.validateReceiverTokenID(receiverLeaf, commitmentTokenID); txErr != nil {
		return synced, txErr, nil
	}

	synced.Tx.SetNonce(senderLeaf.Nonce)

	return synced, nil, nil
}

func (a *Applier) fillSenderWitness(synced *SyncedGenericTransaction, tErr error) (*SyncedGenericTransaction, error, error) {
	witness, appError := a.storage.StateTree.GetLeafWitness(synced.Tx.GetFromStateID())
	if appError != nil {
		return nil, nil, appError
	}
	synced.SenderStateProof.Witness = witness

	return synced, tErr, nil
}

func (a *Applier) validateSenderTokenID(senderState *models.StateLeaf, commitmentTokenID models.Uint256) error {
	if senderState.TokenID.Cmp(&commitmentTokenID) != 0 {
		return errors.WithStack(ErrInvalidSenderTokenID)
	}
	return nil
}

func (a *Applier) validateReceiverTokenID(receiverState *models.StateLeaf, commitmentTokenID models.Uint256) error {
	if receiverState.TokenID.Cmp(&commitmentTokenID) != 0 {
		return errors.WithStack(ErrInvalidReceiverTokenID)
	}
	return nil
}

func validateTxNonce(senderState *models.UserState, txNonce models.Uint256) error {
	comparison := txNonce.Cmp(&senderState.Nonce)
	if comparison > 0 {
		return errors.WithStack(ErrNonceTooHigh)
	} else if comparison < 0 {
		return errors.WithStack(ErrNonceTooLow)
	}
	return nil
}

func calculateStateAfterTx(
	senderState, receiverState models.UserState, // nolint:gocritic
	tx models.GenericTransaction,
) (
	newSenderState, newReceiverState *models.UserState,
	err error,
) {
	if tx.GetToStateID() == nil {
		panic("transaction ToStateID is nil")
	}

	fee := tx.GetFee()
	amount := tx.GetAmount()

	if amount.CmpN(0) <= 0 {
		return nil, nil, errors.WithStack(ErrInvalidTokenAmount)
	}

	totalAmount := amount.Add(&fee)
	if senderState.Balance.Cmp(totalAmount) < 0 {
		return nil, nil, errors.WithStack(ErrBalanceTooLow)
	}

	newSenderState = &senderState
	newSenderState.Nonce = *newSenderState.Nonce.AddN(1)
	newSenderState.Balance = *newSenderState.Balance.Sub(totalAmount)

	if tx.GetFromStateID() == *tx.GetToStateID() {
		newReceiverState = newSenderState.Copy()
	} else {
		newReceiverState = &receiverState
	}
	newReceiverState.Balance = *newReceiverState.Balance.Add(&amount)

	return newSenderState, newReceiverState, nil
}
