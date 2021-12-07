package applier

import "github.com/Worldcoin/hubble-commander/models"

type Proofs struct {
	SenderStateProof   models.StateMerkleProof
	ReceiverStateProof *models.StateMerkleProof
}

type SyncedTxWithProofs struct {
	Tx models.GenericTransaction
	Proofs
}

func NewSyncedTxWithProofs(tx models.GenericTransaction, senderUserState, receiverUserState *models.UserState) *SyncedTxWithProofs {
	return &SyncedTxWithProofs{
		Tx: tx,
		Proofs: Proofs{
			SenderStateProof: models.StateMerkleProof{
				UserState: senderUserState,
			},
			ReceiverStateProof: &models.StateMerkleProof{
				UserState: receiverUserState,
			},
		},
	}
}

func NewSyncedTxWithSenderProof(tx models.GenericTransaction, senderUserState *models.UserState) *SyncedTxWithProofs {
	return &SyncedTxWithProofs{
		Tx: tx,
		Proofs: Proofs{
			SenderStateProof: models.StateMerkleProof{
				UserState: senderUserState,
			},
		},
	}
}
