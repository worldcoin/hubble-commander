package executor

import (
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

func (t *TransactionExecutor) DisputeSignature(batch *eth.DecodedBatch, commitmentIndex int) error {
	switch batch.Type {
	case txtype.Transfer:
		return t.disputeTransferSignature(batch, commitmentIndex)
	case txtype.Create2Transfer:
		return t.disputeCreate2TransferSignature(batch, commitmentIndex)
	case txtype.MassMigration:
		return errors.New("unsupported batch type")
	}
	return nil
}

func (t *TransactionExecutor) disputeTransferSignature(batch *eth.DecodedBatch, commitmentIndex int) error {
	proof, err := t.signatureProof(&batch.Commitments[commitmentIndex])
	if err != nil {
		return err
	}

	targetCommitmentProof, err := targetCommitmentInclusionProof(batch, uint32(commitmentIndex))
	if err != nil {
		return err
	}

	return t.client.DisputeSignatureTransfer(&batch.ID, targetCommitmentProof, proof)
}

func (t *TransactionExecutor) disputeCreate2TransferSignature(batch *eth.DecodedBatch, commitmentIndex int) error {
	proof, err := t.signatureProofWithReceiver(&batch.Commitments[commitmentIndex])
	if err != nil {
		return err
	}

	targetCommitmentProof, err := targetCommitmentInclusionProof(batch, uint32(commitmentIndex))
	if err != nil {
		return err
	}

	return t.client.DisputeSignatureCreate2Transfer(&batch.ID, targetCommitmentProof, proof)
}

func (t *TransactionExecutor) signatureProof(commitment *encoder.DecodedCommitment) (*models.SignatureProof, error) {
	txs, err := encoder.DeserializeTransfers(commitment.Transactions)
	if err != nil {
		return nil, err
	}

	proof := &models.SignatureProof{
		UserStates: make([]models.StateMerkleProof, 0, len(txs)),
		PublicKeys: make([]models.PublicKeyProof, 0, len(txs)),
	}
	for i := range txs {
		stateProof, err := t.userStateProof(txs[i].FromStateID)
		if err != nil {
			return nil, err
		}
		publicKeyProof, err := t.publicKeyProof(stateProof.UserState.PubKeyID)
		if err != nil {
			return nil, err
		}

		proof.UserStates = append(proof.UserStates, *stateProof)
		proof.PublicKeys = append(proof.PublicKeys, *publicKeyProof)
	}
	return proof, nil
}

func (t *TransactionExecutor) signatureProofWithReceiver(
	commitment *encoder.DecodedCommitment,
) (*models.SignatureProofWithReceiver, error) {
	txs, pubKeyIDs, err := encoder.DeserializeCreate2Transfers(commitment.Transactions)
	if err != nil {
		return nil, err
	}

	proof := &models.SignatureProofWithReceiver{
		UserStates:         make([]models.StateMerkleProof, 0, len(txs)),
		SenderPublicKeys:   make([]models.PublicKeyProof, 0, len(txs)),
		ReceiverPublicKeys: make([]models.ReceiverPublicKeyProof, 0, len(txs)),
	}
	for i := range txs {
		stateProof, err := t.userStateProof(txs[i].FromStateID)
		if err != nil {
			return nil, err
		}
		publicKeyProof, err := t.publicKeyProof(stateProof.UserState.PubKeyID)
		if err != nil {
			return nil, err
		}
		receiverPublicKeyProof, err := t.receiverPublicKeyProof(pubKeyIDs[i])
		if err != nil {
			return nil, err
		}

		proof.UserStates = append(proof.UserStates, *stateProof)
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
	witness, err := t.storage.StateTree.GetWitness(leaf.StateID)
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
