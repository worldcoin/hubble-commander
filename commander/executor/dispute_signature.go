package executor

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/batchtype"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

func (c *DisputeContext) DisputeSignature(
	batch *eth.DecodedBatch,
	commitmentIndex int,
	stateProofs []models.StateMerkleProof,
) error {
	switch batch.Type {
	case batchtype.Transfer:
		return c.disputeTransferSignature(batch, commitmentIndex, stateProofs)
	case batchtype.Create2Transfer:
		return c.disputeCreate2TransferSignature(batch, commitmentIndex, stateProofs)
	case batchtype.Genesis, batchtype.MassMigration:
		return errors.New("unsupported batch type")
	}
	return nil
}

func (c *DisputeContext) disputeTransferSignature(
	batch *eth.DecodedBatch,
	commitmentIndex int,
	stateProofs []models.StateMerkleProof,
) error {
	signatureProof, err := c.signatureProof(stateProofs)
	if err != nil {
		return err
	}

	targetCommitmentProof, err := targetCommitmentInclusionProof(batch, uint32(commitmentIndex))
	if err != nil {
		return err
	}

	return c.client.DisputeSignatureTransfer(&batch.ID, targetCommitmentProof, signatureProof)
}

func (c *DisputeContext) disputeCreate2TransferSignature(
	batch *eth.DecodedBatch,
	commitmentIndex int,
	stateProofs []models.StateMerkleProof,
) error {
	signatureProof, err := c.signatureProofWithReceiver(&batch.Commitments[commitmentIndex], stateProofs)
	if err != nil {
		return err
	}

	targetCommitmentProof, err := targetCommitmentInclusionProof(batch, uint32(commitmentIndex))
	if err != nil {
		return err
	}

	return c.client.DisputeSignatureCreate2Transfer(&batch.ID, targetCommitmentProof, signatureProof)
}

func (c *DisputeContext) signatureProof(stateProofs []models.StateMerkleProof) (*models.SignatureProof, error) {
	proof := &models.SignatureProof{
		UserStates: stateProofs,
		PublicKeys: make([]models.PublicKeyProof, 0, len(stateProofs)),
	}

	for i := range stateProofs {
		publicKeyProof, err := c.publicKeyProof(stateProofs[i].UserState.PubKeyID)
		if err != nil {
			return nil, err
		}
		proof.PublicKeys = append(proof.PublicKeys, *publicKeyProof)
	}
	return proof, nil
}

func (c *DisputeContext) signatureProofWithReceiver(
	commitment *encoder.DecodedCommitment,
	stateProofs []models.StateMerkleProof,
) (*models.SignatureProofWithReceiver, error) {
	pubKeyIDs := encoder.DeserializeCreate2TransferPubKeyIDs(commitment.Transactions)

	proof := &models.SignatureProofWithReceiver{
		UserStates:         stateProofs,
		SenderPublicKeys:   make([]models.PublicKeyProof, 0, len(stateProofs)),
		ReceiverPublicKeys: make([]models.ReceiverPublicKeyProof, 0, len(stateProofs)),
	}
	for i := range stateProofs {
		publicKeyProof, err := c.publicKeyProof(stateProofs[i].UserState.PubKeyID)
		if err != nil {
			return nil, err
		}
		receiverPublicKeyProof, err := c.receiverPublicKeyProof(pubKeyIDs[i])
		if err != nil {
			return nil, err
		}

		proof.SenderPublicKeys = append(proof.SenderPublicKeys, *publicKeyProof)
		proof.ReceiverPublicKeys = append(proof.ReceiverPublicKeys, *receiverPublicKeyProof)
	}
	return proof, nil
}

func (c *DisputeContext) publicKeyProof(pubKeyID uint32) (*models.PublicKeyProof, error) {
	account, err := c.storage.AccountTree.Leaf(pubKeyID)
	if err != nil {
		return nil, err
	}
	witness, err := c.storage.AccountTree.GetWitness(pubKeyID)
	if err != nil {
		return nil, err
	}

	return &models.PublicKeyProof{
		PublicKey: &account.PublicKey,
		Witness:   witness,
	}, nil
}

func (c *DisputeContext) receiverPublicKeyProof(pubKeyID uint32) (*models.ReceiverPublicKeyProof, error) {
	account, err := c.storage.AccountTree.Leaf(pubKeyID)
	if err != nil {
		return nil, err
	}
	witness, err := c.storage.AccountTree.GetWitness(pubKeyID)
	if err != nil {
		return nil, err
	}

	return &models.ReceiverPublicKeyProof{
		PublicKeyHash: crypto.Keccak256Hash(account.PublicKey.Bytes()),
		Witness:       witness,
	}, nil
}
