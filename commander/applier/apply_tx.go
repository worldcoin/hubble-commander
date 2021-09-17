package applier

import (
	"errors"

	"github.com/Worldcoin/hubble-commander/models"
)

var (
	ErrNonceTooLow         = errors.New("nonce too low")
	ErrNonceTooHigh        = errors.New("nonce too high")
	ErrInvalidSlicesLength = errors.New("invalid slices length")
	ErrNilReceiverStateID  = errors.New("transfer receiver state id cannot be nil")

	ErrBalanceTooLow              = errors.New("not enough balance")
	ErrInvalidSenderTokenID       = errors.New("invalid sender token ID")
	ErrInvalidReceiverTokenID     = errors.New("invalid receiver token ID")
	ErrInvalidTokenAmount         = errors.New("amount cannot be equal to 0")
	ErrInvalidCommitmentStateRoot = errors.New("invalid commitment post state root")
)

func (a *Applier) ApplyTx(
	tx models.GenericTransaction,
	receiverLeaf *models.StateLeaf,
	commitmentTokenID models.Uint256,
) (transferError, appError error) {
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

	if tErr := validateTransferNonce(&senderLeaf.UserState, tx.GetNonce()); tErr != nil {
		return tErr, nil
	}

	newSenderState, newReceiverState, tErr := calculateStateAfterTransfer(senderLeaf.UserState, receiverLeaf.UserState, tx)
	if tErr != nil {
		return tErr, nil
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

func (a *Applier) ApplyTransferForSync(transfer models.GenericTransaction, commitmentTokenID models.Uint256) (
	synced *SyncedTransfer,
	transferError, appError error,
) {
	receiverLeaf, err := a.storage.StateTree.LeafOrEmpty(*transfer.GetToStateID())
	if err != nil {
		return nil, nil, err
	}

	genericSynced, transferError, appError := a.applyGenericTransactionForSync(transfer, receiverLeaf, commitmentTokenID)
	if appError != nil {
		return nil, nil, appError
	}
	return NewSyncedTransferFromGeneric(genericSynced), transferError, nil
}

func (a *Applier) applyGenericTransactionForSync(
	tx models.GenericTransaction,
	receiverLeaf *models.StateLeaf,
	commitmentTokenID models.Uint256,
) (synced *SyncedGenericTransaction, transferError, appError error) {
	senderLeaf, err := a.storage.StateTree.LeafOrEmpty(tx.GetFromStateID())
	if err != nil {
		return nil, nil, err
	}

	synced = NewPartialSyncedGenericTransaction(tx.Copy(), &senderLeaf.UserState, &receiverLeaf.UserState)

	newSenderState, newReceiverState, tErr := calculateStateAfterTransfer(senderLeaf.UserState, receiverLeaf.UserState, tx)
	if tErr != nil {
		return a.fillSenderWitness(synced, tErr)
	}

	senderWitness, appError := a.storage.StateTree.Set(senderLeaf.StateID, newSenderState)
	if appError != nil {
		return nil, nil, appError
	}
	synced.SenderStateProof.Witness = senderWitness

	if tErr := a.validateSenderTokenID(senderLeaf, commitmentTokenID); tErr != nil {
		return synced, tErr, nil
	}

	receiverWitness, appError := a.storage.StateTree.Set(receiverLeaf.StateID, newReceiverState)
	if appError != nil {
		return nil, nil, appError
	}
	synced.ReceiverStateProof.Witness = receiverWitness

	if tErr := a.validateReceiverTokenID(receiverLeaf, commitmentTokenID); tErr != nil {
		return synced, tErr, nil
	}

	synced.Transaction.SetNonce(senderLeaf.Nonce)

	return synced, nil, nil
}

func (a *Applier) fillSenderWitness(synced *SyncedGenericTransaction, tErr error) (*SyncedGenericTransaction, error, error) {
	witness, appError := a.storage.StateTree.GetLeafWitness(synced.Transaction.GetFromStateID())
	if appError != nil {
		return nil, nil, appError
	}
	synced.SenderStateProof.Witness = witness

	return synced, tErr, nil
}

func (a *Applier) validateSenderTokenID(senderState *models.StateLeaf, commitmentTokenID models.Uint256) error {
	if senderState.TokenID.Cmp(&commitmentTokenID) != 0 {
		return ErrInvalidSenderTokenID
	}
	return nil
}

func (a *Applier) validateReceiverTokenID(receiverState *models.StateLeaf, commitmentTokenID models.Uint256) error {
	if receiverState.TokenID.Cmp(&commitmentTokenID) != 0 {
		return ErrInvalidReceiverTokenID
	}
	return nil
}

func validateTransferNonce(senderState *models.UserState, transferNonce models.Uint256) error {
	comparison := transferNonce.Cmp(&senderState.Nonce)
	if comparison > 0 {
		return ErrNonceTooHigh
	} else if comparison < 0 {
		return ErrNonceTooLow
	}
	return nil
}

func calculateStateAfterTransfer(
	senderState, receiverState models.UserState, // nolint:gocritic
	transfer models.GenericTransaction,
) (
	newSenderState, newReceiverState *models.UserState,
	err error,
) {
	amount := transfer.GetAmount()

	if amount.CmpN(0) <= 0 {
		return nil, nil, ErrInvalidTokenAmount
	}

	totalAmount := amount.Add(transfer.GetFee())
	if senderState.Balance.Cmp(totalAmount) < 0 {
		return nil, nil, ErrBalanceTooLow
	}

	newSenderState = &senderState
	newReceiverState = &receiverState

	newSenderState.Nonce = *newSenderState.Nonce.AddN(1)
	newSenderState.Balance = *newSenderState.Balance.Sub(totalAmount)
	newReceiverState.Balance = *newReceiverState.Balance.Add(&amount)

	return newSenderState, newReceiverState, nil
}
