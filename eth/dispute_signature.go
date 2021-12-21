package eth

import (
	"math/big"

	"github.com/Worldcoin/hubble-commander/contracts/rollup"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
)

func (c *Client) DisputeSignatureTransfer(
	batchID *models.Uint256,
	batchHash *common.Hash,
	targetProof *models.TransferCommitmentInclusionProof,
	signatureProof *models.SignatureProof,
) error {
	transaction, err := c.rollup().
		WithGasLimit(*c.config.SignatureDisputeGasLimit).
		DisputeSignatureTransfer(
			batchID.ToBig(),
			*transferProofToCalldata(targetProof),
			*signatureProofToCalldata(signatureProof),
		)
	if err != nil {
		return handleDisputeSignatureError(err)
	}

	err = c.waitForDispute(batchID, batchHash, transaction)
	if err == ErrBatchAlreadyDisputed || err == ErrRollbackInProcess {
		log.Info(err)
		return nil
	}
	return err
}

func (c *Client) DisputeSignatureCreate2Transfer(
	batchID *models.Uint256,
	batchHash *common.Hash,
	targetProof *models.TransferCommitmentInclusionProof,
	signatureProof *models.SignatureProofWithReceiver,
) error {
	transaction, err := c.rollup().
		WithGasLimit(*c.config.SignatureDisputeGasLimit).
		DisputeSignatureCreate2Transfer(
			batchID.ToBig(),
			*transferProofToCalldata(targetProof),
			*signatureProofWithReceiverToCalldata(signatureProof),
		)
	if err != nil {
		return handleDisputeSignatureError(err)
	}

	err = c.waitForDispute(batchID, batchHash, transaction)
	if err == ErrBatchAlreadyDisputed || err == ErrRollbackInProcess {
		log.Info(err)
		return nil
	}
	return err
}

func (c *Client) DisputeSignatureMassMigration(
	batchID *models.Uint256,
	batchHash *common.Hash,
	targetProof *models.MMCommitmentInclusionProof,
	signatureProof *models.SignatureProof,
) error {
	transaction, err := c.rollup().
		WithGasLimit(*c.config.SignatureDisputeGasLimit).
		DisputeSignatureMassMigration(
			batchID.ToBig(),
			*massMigrationProofToCalldata(targetProof),
			*signatureProofToCalldata(signatureProof),
		)
	if err != nil {
		return handleDisputeSignatureError(err)
	}

	err = c.waitForDispute(batchID, batchHash, transaction)
	if err == ErrBatchAlreadyDisputed || err == ErrRollbackInProcess {
		log.Info(err)
		return nil
	}
	return err
}

func handleDisputeSignatureError(err error) error {
	errMsg := getGasEstimateErrorMessage(err)
	if errMsg == msgSignatureMissingBatch || errMsg == msgBatchAlreadyDisputed {
		log.Info(err.Error())
		return nil
	}
	return err
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
