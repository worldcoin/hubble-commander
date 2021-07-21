package eth

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (c *Client) DisputeSignatureTransfer(
	batchID *models.Uint256,
	target *models.TransferCommitmentInclusionProof,
	signatureProof *models.SignatureProof,
) error {
	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := c.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	transaction, err := c.rollup().DisputeSignatureTransfer(
		batchID.ToBig(),
		*TransferProofToCalldata(target),
		*signatureProofToCalldata(signatureProof),
	)
	if err != nil {
		return err
	}
	return c.waitForRollbackToFinish(sink, subscription, transaction.Hash())
}

func signatureProofToCalldata(proof *models.SignatureProof) *rollup.TypesSignatureProof {
	result := &rollup.TypesSignatureProof{
		States:          make([]rollup.TypesUserState, 0, len(proof.UserStates)),
		StateWitnesses:  make([][][32]byte, 0, len(proof.UserStates)),
		Pubkeys:         make([][4]*big.Int, 0, len(proof.PublicKeys)),
		PubkeyWitnesses: make([][][32]byte, 0, len(proof.PublicKeys)),
	}
	for i := range proof.UserStates {
		stateProof := stateMerkleProofToCalldata(&proof.UserStates[i])
		result.States = append(result.States, stateProof.State)
		result.StateWitnesses = append(result.StateWitnesses, stateProof.Witness)

		result.Pubkeys = append(result.Pubkeys, proof.PublicKeys[i].PublicKey.BigInts())
		result.PubkeyWitnesses = append(result.PubkeyWitnesses, proof.PublicKeys[i].Witness.Bytes())
	}
	return result
}

func signatureProofWithReceiverToCalldata(proof *models.SignatureProofWithReceiver) *rollup.TypesSignatureProofWithReceiver {
	result := &rollup.TypesSignatureProofWithReceiver{
		States:                  make([]rollup.TypesUserState, 0, len(proof.UserStates)),
		StateWitnesses:          make([][][32]byte, 0, len(proof.UserStates)),
		PubkeysSender:           make([][4]*big.Int, 0, len(proof.SenderPublicKeys)),
		PubkeyWitnessesSender:   make([][][32]byte, 0, len(proof.SenderPublicKeys)),
		PubkeyHashesReceiver:    make([][32]byte, 0, len(proof.UserStates)),
		PubkeyWitnessesReceiver: make([][][32]byte, 0, len(proof.SenderPublicKeys)),
	}
	for i := range proof.UserStates {
		stateProof := stateMerkleProofToCalldata(&proof.UserStates[i])
		result.States = append(result.States, stateProof.State)
		result.StateWitnesses = append(result.StateWitnesses, stateProof.Witness)

		result.PubkeysSender = append(result.PubkeysSender, proof.SenderPublicKeys[i].PublicKey.BigInts())
		result.PubkeyWitnessesSender = append(result.PubkeyWitnessesSender, proof.SenderPublicKeys[i].Witness.Bytes())

		result.PubkeyHashesReceiver = append(result.PubkeyHashesReceiver, proof.ReceiverPublicKeys[i].PublicKeyHash)
		result.PubkeyWitnessesReceiver = append(result.PubkeyWitnessesReceiver, proof.ReceiverPublicKeys[i].Witness.Bytes())
	}
	return result
}
