package applier

import "github.com/Worldcoin/hubble-commander/models"

type Proofs struct {
	SenderStateProof   models.StateMerkleProof
	ReceiverStateProof models.StateMerkleProof
}

type SyncedTransfer struct {
	Transfer *models.Transfer
	Proofs
}

func NewSyncedTransferFromGeneric(generic *SyncedGenericTransaction) *SyncedTransfer {
	return &SyncedTransfer{
		Transfer: generic.Tx.ToTransfer(),
		Proofs:   generic.Proofs,
	}
}

type SyncedCreate2Transfer struct {
	Transfer *models.Create2Transfer
	Proofs
}

func NewSyncedCreate2TransferFromGeneric(generic *SyncedGenericTransaction) *SyncedCreate2Transfer {
	return &SyncedCreate2Transfer{
		Transfer: generic.Tx.ToCreate2Transfer(),
		Proofs:   generic.Proofs,
	}
}

type SyncedGenericTransaction struct {
	Tx models.GenericTransaction
	Proofs
}

func NewPartialSyncedGenericTransaction(
	tx models.GenericTransaction,
	senderUserState, receiverUserState *models.UserState,
) *SyncedGenericTransaction {
	return &SyncedGenericTransaction{
		Tx: tx,
		Proofs: Proofs{
			SenderStateProof: models.StateMerkleProof{
				UserState: senderUserState,
			},
			ReceiverStateProof: models.StateMerkleProof{
				UserState: receiverUserState,
			},
		},
	}
}

type SyncedTx interface {
	SyncedTx() models.GenericTransaction
	SyncedPubKeyID() uint32
}

type SyncedTransfer2 struct {
	tx *models.Transfer
}

func (a *SyncedTransfer2) SyncedTx() models.GenericTransaction {
	return a.tx
}

func (a *SyncedTransfer2) SyncedPubKeyID() uint32 {
	panic("SyncedPubKeyID cannot be invoked on SyncedTransfer2")
}

type SyncedC2T2 struct {
	Tx       *models.Create2Transfer
	PubKeyID uint32
}

func (a *SyncedC2T2) SyncedTx() models.GenericTransaction {
	return a.Tx
}

func (a *SyncedC2T2) SyncedPubKeyID() uint32 {
	return a.PubKeyID
}
