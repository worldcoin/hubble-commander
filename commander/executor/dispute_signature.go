package executor

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

func (t *TransactionExecutor) DisputeSignature(
	batch *eth.DecodedBatch,
	commitmentIndex int,
	stateProofs []models.StateMerkleProof,
) error {
	switch batch.Type {
	case txtype.Transfer:
		return t.disputeTransferSignature(batch, commitmentIndex, stateProofs)
	case txtype.Create2Transfer:
		return t.disputeCreate2TransferSignature(batch, commitmentIndex, stateProofs)
	case txtype.Genesis, txtype.MassMigration:
		return errors.New("unsupported batch type")
	}
	return nil
}

func (t *TransactionExecutor) disputeTransferSignature(
	batch *eth.DecodedBatch,
	commitmentIndex int,
	stateProofs []models.StateMerkleProof,
) error {
	signatureProof, err := t.signatureProof(stateProofs)
	if err != nil {
		return err
	}

	targetCommitmentProof, err := targetCommitmentInclusionProof(batch, uint32(commitmentIndex))
	if err != nil {
		return err
	}

	return t.client.DisputeSignatureTransfer(&batch.ID, targetCommitmentProof, signatureProof)
}

func (t *TransactionExecutor) disputeCreate2TransferSignature(
	batch *eth.DecodedBatch,
	commitmentIndex int,
	stateProofs []models.StateMerkleProof,
) error {
	signatureProof, err := t.signatureProofWithReceiver(&batch.Commitments[commitmentIndex], stateProofs)
	if err != nil {
		return err
	}

	targetCommitmentProof, err := targetCommitmentInclusionProof(batch, uint32(commitmentIndex))
	if err != nil {
		return err
	}

	return t.client.DisputeSignatureCreate2Transfer(&batch.ID, targetCommitmentProof, signatureProof)
}

func (t *TransactionExecutor) signatureProof(stateProofs []models.StateMerkleProof) (*models.SignatureProof, error) {
	proof := &models.SignatureProof{
		UserStates: stateProofs,
		PublicKeys: make([]models.PublicKeyProof, 0, len(stateProofs)),
	}

	for i := range stateProofs {
		publicKeyProof, err := t.publicKeyProof(stateProofs[i].UserState.PubKeyID)
		if err != nil {
			return nil, err
		}
		proof.PublicKeys = append(proof.PublicKeys, *publicKeyProof)
	}
	return proof, nil
}

func (t *TransactionExecutor) signatureProofWithReceiver(
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
		publicKeyProof, err := t.publicKeyProof(stateProofs[i].UserState.PubKeyID)
		if err != nil {
			return nil, err
		}
		receiverPublicKeyProof, err := t.receiverPublicKeyProof(pubKeyIDs[i])
		if err != nil {
			return nil, err
		}

		proof.SenderPublicKeys = append(proof.SenderPublicKeys, *publicKeyProof)
		proof.ReceiverPublicKeys = append(proof.ReceiverPublicKeys, *receiverPublicKeyProof)
	}
	return proof, nil
}

func (t *TransactionExecutor) userStateProof(stateID uint32) (*models.StateMerkleProof, error) {
	leaf, err := t.storage.StateTree.Leaf(stateID)
	if err != nil {
		return nil, err
	}
	witness, err := t.storage.StateTree.GetUserStateWitness(leaf.StateID)
	if err != nil {
		return nil, err
	}
	return &models.StateMerkleProof{
		UserState: &leaf.UserState,
		Witness:   witness,
	}, nil
}

func (t *TransactionExecutor) publicKeyProof(pubKeyID uint32) (*models.PublicKeyProof, error) {
	account, err := t.storage.AccountTree.Leaf(pubKeyID)
	if err != nil {
		return nil, err
	}
	witness, err := t.storage.AccountTree.GetWitness(pubKeyID)
	if err != nil {
		return nil, err
	}

	return &models.PublicKeyProof{
		PublicKey: &account.PublicKey,
		Witness:   witness,
	}, nil
}

func (t *TransactionExecutor) receiverPublicKeyProof(pubKeyID uint32) (*models.ReceiverPublicKeyProof, error) {
	account, err := t.storage.AccountTree.Leaf(pubKeyID)
	if err != nil {
		return nil, err
	}
	witness, err := t.storage.AccountTree.GetWitness(pubKeyID)
	if err != nil {
		return nil, err
	}

	return &models.ReceiverPublicKeyProof{
		PublicKeyHash: crypto.Keccak256Hash(account.PublicKey.Bytes()),
		Witness:       witness,
	}, nil
}

func (t *TransactionExecutor) stateMerkleProofs(transfers models.GenericTransactionArray) ([]models.StateMerkleProof, error) {
	proofs := make([]models.StateMerkleProof, 0, transfers.Len())
	for i := 0; i < transfers.Len(); i++ {
		stateProof, err := t.userStateProof(transfers.At(i).GetFromStateID())
		if err != nil {
			return nil, err
		}
		proofs = append(proofs, *stateProof)
	}
	return proofs, nil
}
