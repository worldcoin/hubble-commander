package applier

import (
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
)

func (a *Applier) ApplyCreate2Transfer(
	create2Transfer *models.Create2Transfer,
	commitmentTokenID models.Uint256,
) (applyResult *ApplySingleC2TResult, transferError, appError error) {
	pubKeyID, appError := a.getOrRegisterPubKeyID(&create2Transfer.ToPublicKey, commitmentTokenID)
	if appError != nil {
		return nil, nil, appError
	}

	nextAvailableStateID, appError := a.storage.StateTree.NextAvailableStateID()
	if appError != nil {
		return nil, nil, appError
	}

	applyResult = &ApplySingleC2TResult{
		tx:            create2Transfer.Clone(),
		addedPubKeyID: *pubKeyID,
	}
	applyResult.tx.ToStateID = nextAvailableStateID

	receiverLeaf, appError := newUserLeaf(*nextAvailableStateID, *pubKeyID, commitmentTokenID)
	if appError != nil {
		return nil, nil, appError
	}
	transferError, appError = a.ApplyTx(applyResult.tx, receiverLeaf, commitmentTokenID)
	return applyResult, transferError, appError
}

func (a *Applier) ApplyCreate2TransferForSync(
	create2Transfer *models.Create2Transfer,
	pubKeyID uint32,
	commitmentTokenID models.Uint256,
) (synced *SyncedCreate2Transfer, transferError, appError error) {
	if create2Transfer.ToStateID == nil {
		return nil, nil, ErrNilReceiverStateID
	}

	receiverLeaf, appError := newUserLeaf(*create2Transfer.ToStateID, pubKeyID, commitmentTokenID)
	if appError != nil {
		return nil, nil, appError
	}
	genericSynced, transferError, appError := a.applyTxForSync(create2Transfer, receiverLeaf, commitmentTokenID)
	if appError != nil {
		return nil, nil, appError
	}
	return NewSyncedCreate2TransferFromGeneric(genericSynced), transferError, nil
}

func (a *Applier) getOrRegisterPubKeyID(
	publicKey *models.PublicKey,
	tokenID models.Uint256,
) (*uint32, error) {
	pubKeyID, err := a.storage.GetUnusedPubKeyID(publicKey, &tokenID)
	if err != nil && !st.IsNotFoundError(err) {
		return nil, err
	} else if st.IsNotFoundError(err) {
		return a.client.RegisterAccountAndWait(publicKey)
	}
	return pubKeyID, nil
}

func newUserLeaf(stateID, pubKeyID uint32, tokenID models.Uint256) (*models.StateLeaf, error) {
	return st.NewStateLeaf(stateID, &models.UserState{
		PubKeyID: pubKeyID,
		TokenID:  tokenID,
		Balance:  models.MakeUint256(0),
		Nonce:    models.MakeUint256(0),
	})
}
