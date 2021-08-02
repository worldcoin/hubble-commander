package eth

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (c *Client) DisputeSignatureTransfer(
	batchID *models.Uint256,
	targetProof *models.TransferCommitmentInclusionProof,
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
		*TransferProofToCalldata(targetProof),
		*signatureProofToCalldata(signatureProof),
	)
	if err != nil {
		return err
	}
	return c.waitForRollbackToFinish(sink, subscription, transaction.Hash())
}

func (c *Client) DisputeSignatureCreate2Transfer(
	batchID *models.Uint256,
	targetProof *models.TransferCommitmentInclusionProof,
	signatureProof *models.SignatureProofWithReceiver,
) error {
	sink := make(chan *rollup.RollupRollbackStatus)
	subscription, err := c.Rollup.WatchRollbackStatus(&bind.WatchOpts{}, sink)
	if err != nil {
		return err
	}
	defer subscription.Unsubscribe()

	transaction, err := c.rollup().DisputeSignatureCreate2Transfer(
		batchID.ToBig(),
		*TransferProofToCalldata(targetProof),
		*signatureProofWithReceiverToCalldata(signatureProof),
	)
	if err != nil {
		return err
	}
	return c.waitForRollbackToFinish(sink, subscription, transaction.Hash())
}

func signatureProofToCalldata(proof *models.SignatureProof) *rollup.TypesSignatureProof {
	states := make([]rollup.TypesUserState, 0, len(proof.UserStates))
	stateWitnesses := make([][][32]byte, 0, len(proof.UserStates))
	pubkeys := make([][4]*big.Int, 0, len(proof.PublicKeys))
	pubkeyWitnesses := make([][][32]byte, 0, len(proof.PublicKeys))

	for i := range proof.UserStates {
		stateProof := stateMerkleProofToCalldata(&proof.UserStates[i])
		states = append(states, stateProof.State)
		stateWitnesses = append(stateWitnesses, stateProof.Witness)

		pubkeys = append(pubkeys, proof.PublicKeys[i].PublicKey.BigInts())
		pubkeyWitnesses = append(pubkeyWitnesses, proof.PublicKeys[i].Witness.Bytes())
	}
	return &rollup.TypesSignatureProof{
		States:          states,
		StateWitnesses:  stateWitnesses,
		Pubkeys:         pubkeys,
		PubkeyWitnesses: pubkeyWitnesses,
	}
}

func signatureProofWithReceiverToCalldata(proof *models.SignatureProofWithReceiver) *rollup.TypesSignatureProofWithReceiver {
	states := make([]rollup.TypesUserState, 0, len(proof.UserStates))
	stateWitnesses := make([][][32]byte, 0, len(proof.UserStates))
	senderPubkeys := make([][4]*big.Int, 0, len(proof.SenderPublicKeys))
	senderPubkeyWitnesses := make([][][32]byte, 0, len(proof.SenderPublicKeys))
	receiverPubkeyHashes := make([][32]byte, 0, len(proof.ReceiverPublicKeys))
	receiverPubkeyWitnesses := make([][][32]byte, 0, len(proof.ReceiverPublicKeys))

	for i := range proof.UserStates {
		stateProof := stateMerkleProofToCalldata(&proof.UserStates[i])
		states = append(states, stateProof.State)
		stateWitnesses = append(stateWitnesses, stateProof.Witness)

		senderPubkeys = append(senderPubkeys, proof.SenderPublicKeys[i].PublicKey.BigInts())
		senderPubkeyWitnesses = append(senderPubkeyWitnesses, proof.SenderPublicKeys[i].Witness.Bytes())

		receiverPubkeyHashes = append(receiverPubkeyHashes, proof.ReceiverPublicKeys[i].PublicKeyHash)
		receiverPubkeyWitnesses = append(receiverPubkeyWitnesses, proof.ReceiverPublicKeys[i].Witness.Bytes())
	}
	return &rollup.TypesSignatureProofWithReceiver{
		States:                  states,
		StateWitnesses:          stateWitnesses,
		PubkeysSender:           senderPubkeys,
		PubkeyWitnessesSender:   senderPubkeyWitnesses,
		PubkeyHashesReceiver:    receiverPubkeyHashes,
		PubkeyWitnessesReceiver: receiverPubkeyWitnesses,
	}
}
