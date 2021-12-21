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
	senderLeaf, err := a.storage.StateTree.Leaf(tx.GetFromStateID())
	if err != nil {
		return nil, err
	}

	appError = a.validateSenderTokenID(senderLeaf, commitmentTokenID)
	if appError != nil {
		return nil, appError
	}

	appError = a.validateReceiverTokenID(receiverLeaf, commitmentTokenID)
	if appError != nil {
		return nil, appError
	}

	if tErr := validateTxNonce(&senderLeaf.UserState, tx.GetNonce()); tErr != nil {
		return tErr, nil
	}

	newSenderState, newReceiverState, txErr := calculateStateAfterTx(senderLeaf.UserState, receiverLeaf.UserState, tx)
	if txErr != nil {
		return txErr, nil
	}

	_, appError = a.storage.StateTree.Set(senderLeaf.StateID, newSenderState)
	if appError != nil {
		return nil, appError
	}
	_, appError = a.storage.StateTree.Set(receiverLeaf.StateID, newReceiverState)
	if appError != nil {
		return nil, appError
	}

	return nil, nil
}

func (a *Applier) applyTxForSync(
	tx models.GenericTransaction,
	receiverLeaf *models.StateLeaf,
	commitmentTokenID models.Uint256,
) (synced *SyncedTxWithProofs, txError, appError error) {
	senderLeaf, appError := a.storage.StateTree.LeafOrEmpty(tx.GetFromStateID())
	if appError != nil {
		return nil, nil, appError
	}

	synced = NewSyncedTxWithProofs(tx.Copy(), &senderLeaf.UserState, &receiverLeaf.UserState)

	newSenderState, newReceiverState, txErr := calculateStateAfterTx(senderLeaf.UserState, receiverLeaf.UserState, tx)
	if txErr != nil {
		return a.fillSenderWitness(synced, txErr)
	}

	senderWitness, appError := a.storage.StateTree.Set(senderLeaf.StateID, newSenderState)
	if appError != nil {
		return nil, nil, appError
	}
	synced.SenderStateProof.Witness = senderWitness

	txErr = a.validateSenderTokenID(senderLeaf, commitmentTokenID)
	if txErr != nil {
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

	txErr = a.validateReceiverTokenID(receiverLeaf, commitmentTokenID)
	if txErr != nil {
		return synced, txErr, nil
	}

	synced.Tx.SetNonce(senderLeaf.Nonce)

	return synced, nil, nil
}

func (a *Applier) fillSenderWitness(synced *SyncedTxWithProofs, tErr error) (*SyncedTxWithProofs, error, error) {
	witness, appError := a.storage.StateTree.GetLeafWitness(synced.Tx.GetFromStateID())
	if appError != nil {
		return nil, nil, appError
	}
	synced.SenderStateProof.Witness = witness

	return synced, tErr, nil
}

func (a *Applier) validateSenderTokenID(senderState *models.StateLeaf, commitmentTokenID models.Uint256) error {
	if !senderState.TokenID.Eq(&commitmentTokenID) {
		return errors.WithStack(ErrInvalidSenderTokenID)
	}
	return nil
}

func (a *Applier) validateReceiverTokenID(receiverState *models.StateLeaf, commitmentTokenID models.Uint256) error {
	if !receiverState.TokenID.Eq(&commitmentTokenID) {
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

	newSenderState, err = calculateSenderStateAfterTx(senderState, tx)
	if err != nil {
		return nil, nil, err
	}

	if tx.GetFromStateID() == *tx.GetToStateID() {
		newReceiverState = calculateReceiverStateAfterTx(*newSenderState.Copy(), tx)
	} else {
		newReceiverState = calculateReceiverStateAfterTx(receiverState, tx)
	}

	return newSenderState, newReceiverState, nil
}

func calculateSenderStateAfterTx(
	senderState models.UserState, // nolint:gocritic
	tx models.GenericTransaction,
) (newSenderState *models.UserState, err error) {
	fee := tx.GetFee()
	amount := tx.GetAmount()

	if amount.CmpN(0) <= 0 {
		return nil, errors.WithStack(ErrInvalidTokenAmount)
	}

	totalAmount := amount.Add(&fee)
	if senderState.Balance.Cmp(totalAmount) < 0 {
		return nil, errors.WithStack(ErrBalanceTooLow)
	}

	newSenderState = &senderState

	newSenderState.Nonce = *newSenderState.Nonce.AddN(1)
	newSenderState.Balance = *newSenderState.Balance.Sub(totalAmount)

	return newSenderState, nil
}

func calculateReceiverStateAfterTx(
	receiverState models.UserState, // nolint:gocritic
	tx models.GenericTransaction,
) *models.UserState {
	amount := tx.GetAmount()
	newReceiverState := &receiverState
	newReceiverState.Balance = *newReceiverState.Balance.Add(&amount)
	return newReceiverState
}
