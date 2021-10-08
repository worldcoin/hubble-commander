package applier

import (
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/pkg/errors"
)

func (a *Applier) ApplyCreate2Transfer(
	create2Transfer *models.Create2Transfer,
	commitmentTokenID models.Uint256,
) (applyResult *ApplySingleC2TResult, transferError, appError error) {
	pubKeyID, isPending, appError := a.getPubKeyID(&create2Transfer.ToPublicKey, commitmentTokenID)
	if appError != nil {
		return nil, nil, appError
	}

	nextAvailableStateID, appError := a.storage.StateTree.NextAvailableStateID()
	if appError != nil {
		return nil, nil, appError
	}

	applyResult = &ApplySingleC2TResult{
		tx:       create2Transfer.Clone(),
		pubKeyID: *pubKeyID,
	}
	applyResult.tx.ToStateID = nextAvailableStateID

	if isPending {
		applyResult.pendingAccount = &models.AccountLeaf{
			PubKeyID:  *pubKeyID,
			PublicKey: create2Transfer.ToPublicKey,
		}
	}

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
) (synced *SyncedGenericTransaction, transferError, appError error) {
	if create2Transfer.ToStateID == nil {
		return nil, nil, errors.WithStack(ErrNilReceiverStateID)
	}

	receiverLeaf, appError := newUserLeaf(*create2Transfer.ToStateID, pubKeyID, commitmentTokenID)
	if appError != nil {
		return nil, nil, appError
	}
	genericSynced, transferError, appError := a.applyTxForSync(create2Transfer, receiverLeaf, commitmentTokenID)
	if appError != nil {
		return nil, nil, appError
	}
	return genericSynced, transferError, nil
}

func (a *Applier) getPubKeyID(
	publicKey *models.PublicKey,
	tokenID models.Uint256,
) (pubKeyID *uint32, isPending bool, err error) {
	pubKeyID, err = a.storage.GetUnusedPubKeyID(publicKey, &tokenID)
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
