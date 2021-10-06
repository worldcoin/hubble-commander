package executor

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
		Transfer: generic.Transaction.(*models.Transfer),
		Proofs:   generic.Proofs,
	}
}

type SyncedCreate2Transfer struct {
	Transfer *models.Create2Transfer
	Proofs
}

func NewSyncedCreate2TransferFromGeneric(generic *SyncedGenericTransaction) *SyncedCreate2Transfer {
	return &SyncedCreate2Transfer{
		Transfer: generic.Transaction.(*models.Create2Transfer),
		Proofs:   generic.Proofs,
	}
}

type SyncedGenericTransaction struct {
	Transaction models.GenericTransaction
	Proofs
}

func NewPartialSyncedGenericTransaction(
	tx models.GenericTransaction,
	senderUserState, receiverUserState *models.UserState,
) *SyncedGenericTransaction {
	return &SyncedGenericTransaction{
		Transaction: tx,
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