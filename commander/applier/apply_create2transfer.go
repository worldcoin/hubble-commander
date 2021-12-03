package applier

import (
	"fmt"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
)

var ErrNilReceiverStateID = fmt.Errorf("transfer receiver state id cannot be nil")

func (a *Applier) ApplyCreate2Transfer(
	create2Transfer *models.Create2Transfer,
	commitmentTokenID models.Uint256,
) (applyResult *ApplySingleC2TResult, txError, appError error) {
	pubKeyID, isPending, appError := a.getPubKeyID(&create2Transfer.ToPublicKey)
	if appError != nil {
		return nil, nil, appError
	}

	nextAvailableStateID, appError := a.storage.StateTree.NextAvailableStateID()
	if appError != nil {
		return nil, nil, appError
	}

	receiverLeaf, appError := newUserLeaf(*nextAvailableStateID, *pubKeyID, commitmentTokenID)
	if appError != nil {
		return nil, nil, appError
	}

	updatedCreate2Transfer := create2Transfer.Clone()
	updatedCreate2Transfer.ToStateID = nextAvailableStateID

	txError, appError = a.ApplyTx(updatedCreate2Transfer, receiverLeaf, commitmentTokenID)
	if txError != nil || appError != nil {
		return nil, txError, appError
	}

	applyResult = &ApplySingleC2TResult{
		tx:       updatedCreate2Transfer,
		pubKeyID: *pubKeyID,
	}

	if isPending {
		applyResult.pendingAccount = &models.AccountLeaf{
			PubKeyID:  *pubKeyID,
			PublicKey: updatedCreate2Transfer.ToPublicKey,
		}
	}

	return applyResult, nil, nil
}

func (a *Applier) ApplyCreate2TransferForSync(
	create2Transfer *models.Create2Transfer,
	pubKeyID uint32,
	commitmentTokenID models.Uint256,
) (synced *SyncedGenericTransaction, txError, appError error) {
	if create2Transfer.ToStateID == nil {
		return nil, nil, errors.WithStack(ErrNilReceiverStateID)
	}

	receiverLeaf, appError := newUserLeaf(*create2Transfer.ToStateID, pubKeyID, commitmentTokenID)
	if appError != nil {
		return nil, nil, appError
	}
	genericSynced, txError, appError := a.applyTxForSync(create2Transfer, receiverLeaf, commitmentTokenID)
	if appError != nil {
		return nil, nil, appError
	}
	return genericSynced, txError, nil
}

func (a *Applier) getPubKeyID(publicKey *models.PublicKey) (pubKeyID *uint32, isPending bool, err error) {
	pubKeyID, err = a.storage.GetFirstPubKeyID(publicKey)
	if err == nil {
		return pubKeyID, false, nil
	}
	if st.IsNotFoundError(err) {
		pubKeyID, err = a.storage.AccountTree.NextBatchAccountPubKeyID()
		if err != nil {
			return nil, false, err
		}
		return pubKeyID, true, err
	}
	return nil, false, err
}

func newUserLeaf(stateID, pubKeyID uint32, tokenID models.Uint256) (*models.StateLeaf, error) {
	return st.NewStateLeaf(stateID, &models.UserState{
		PubKeyID: pubKeyID,
		TokenID:  tokenID,
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	})
}
