package applier

import "github.com/Worldcoin/hubble-commander/models"

type Proofs struct {
	SenderStateProof   models.StateMerkleProof
	ReceiverStateProof models.StateMerkleProof
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

func NewSenderPartialSyncedGenericTransaction(
	tx models.GenericTransaction,
	senderUserState *models.UserState,
) *SyncedGenericTransaction {
	return &SyncedGenericTransaction{
		Tx: tx,
		Proofs: Proofs{
			SenderStateProof: models.StateMerkleProof{
				UserState: senderUserState,
			},
		},
	}
}
