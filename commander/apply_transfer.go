package commander

import (
	"fmt"
	"math/big"

	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/storage"
)

func ApplyTransfer(stateTree *storage.StateTree, tx *models.Transaction) (txError error, appError error) {
	if stateTree == nil {
		return nil, fmt.Errorf("state tree cannot be nil")
	}
	if tx == nil {
		return nil, fmt.Errorf("transaction cannot be nil")
	}

	senderIndex := uint32(tx.FromIndex.Uint64())
	senderLeaf, err := stateTree.Leaf(senderIndex)
	if err != nil {
		return nil, err
	}
	receiverIndex := uint32(tx.ToIndex.Uint64())
	receiverLeaf, err := stateTree.Leaf(receiverIndex)
	if err != nil {
		return nil, err
	}

	if senderLeaf == nil || receiverLeaf == nil {
		return fmt.Errorf("sender/receiver cannot be nil"), nil
	}

	senderState := senderLeaf.UserState
	receiverState := receiverLeaf.UserState

	newSenderState, newReceiverState, err := CalculateStateAfterTransfer(&senderState, &receiverState, tx)
	if err != nil {
		return err, nil
	}

	err = stateTree.Set(senderIndex, &newSenderState)
	if err != nil {
		return nil, err
	}
	err = stateTree.Set(receiverIndex, &newReceiverState)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func CalculateStateAfterTransfer(
	senderState,
	receiverState *models.UserState,
	tx *models.Transaction,
) (
	newSenderState models.UserState,
	newReceiverState models.UserState,
	err error,
) {
	// TODO: Signature validation

	if senderState.Nonce.Cmp(&tx.Nonce.Int) != 0 {
		err = fmt.Errorf("incorrect nonce")
		return
	}

	totalAmount := big.NewInt(0)
	totalAmount.Add(&tx.Amount.Int, &tx.Fee.Int)

	if senderState.Balance.Cmp(totalAmount) < 0 {
		err = fmt.Errorf("amount exceeds balance")
		return
	}

	newSenderState = *senderState
	newReceiverState = *receiverState

	newSenderState.Nonce.Add(&senderState.Nonce.Int, big.NewInt(1))

	newSenderState.Balance.Sub(&senderState.Balance.Int, totalAmount)

	newReceiverState.Balance.Add(&receiverState.Balance.Int, &tx.Amount.Int)

	return newSenderState, newReceiverState, nil
}
