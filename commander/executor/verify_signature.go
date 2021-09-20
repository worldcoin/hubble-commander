package executor

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
)

const InvalidSignatureMessage = "invalid commitment signature"

func (c *SyncContext) verifyTransferSignature(commitment *encoder.DecodedCommitment, txs models.GenericTransactionArray) error {
	domain, err := c.client.GetDomain()
	if err != nil {
		return err
	}

	messages := make([][]byte, txs.Len())
	publicKeys := make([]*models.PublicKey, txs.Len())
	for i := 0; i < txs.Len(); i++ {
		publicKeys[i], err = c.storage.GetPublicKeyByStateID(txs.At(i).GetFromStateID())
		if err != nil {
			return err
		}
		messages[i], err = c.Syncer.EncodeTxForSigning(txs.At(i))
		if err != nil {
			return err
		}
	}

	return c.verifyCommitmentSignature(&commitment.CombinedSignature, domain, messages, publicKeys, txs)
}

func (c *ExecutionContext) verifyCreate2TransferSignature(
	commitment *encoder.DecodedCommitment,
	transfers []models.Create2Transfer,
) error {
	domain, err := c.client.GetDomain()
	if err != nil {
		return err
	}

	messages := make([][]byte, len(transfers))
	publicKeys := make([]*models.PublicKey, len(transfers))
	for i := range transfers {
		publicKeys[i], err = c.storage.GetPublicKeyByStateID(transfers[i].FromStateID)
		if err != nil {
			return err
		}
		messages[i], err = encoder.EncodeCreate2TransferForSigning(&transfers[i])
		if err != nil {
			return err
		}
	}

	genericTxs := models.Create2TransferArray(transfers)
	return c.verifyCommitmentSignature(&commitment.CombinedSignature, domain, messages, publicKeys, genericTxs)
}

func (c *ExecutionContext) verifyCommitmentSignature(
	signature *models.Signature,
	domain *bls.Domain,
	messages [][]byte,
	publicKeys []*models.PublicKey,
	transfers models.GenericTransactionArray,
) error {
	if len(messages) == 0 {
		return nil
	}
	sig, err := bls.NewSignatureFromBytes(signature.Bytes(), *domain)
	if err != nil {
		return c.createDisputableSignatureError(err.Error(), transfers)
	}
	aggregatedSignature := bls.AggregatedSignature{Signature: sig}
	isValid, err := aggregatedSignature.Verify(messages, publicKeys)
	if err != nil {
		return err
	}
	if !isValid {
		return c.createDisputableSignatureError(InvalidSignatureMessage, transfers)
	}
	return nil
}

func (c *ExecutionContext) createDisputableSignatureError(reason string, transfers models.GenericTransactionArray) error {
	proofs, proofErr := c.stateMerkleProofs(transfers)
	if proofErr != nil {
		return proofErr
	}
	return NewDisputableErrorWithProofs(Signature, reason, proofs)
}

func (c *ExecutionContext) stateMerkleProofs(transfers models.GenericTransactionArray) ([]models.StateMerkleProof, error) {
	proofs := make([]models.StateMerkleProof, 0, transfers.Len())
	for i := 0; i < transfers.Len(); i++ {
		stateProof, err := c.userStateProof(transfers.At(i).GetFromStateID())
		if err != nil {
			return nil, err
		}
		proofs = append(proofs, *stateProof)
	}
	return proofs, nil
}

func (c *ExecutionContext) userStateProof(stateID uint32) (*models.StateMerkleProof, error) {
	leaf, err := c.storage.StateTree.Leaf(stateID)
	if err != nil {
		return nil, err
	}
	witness, err := c.storage.StateTree.GetLeafWitness(leaf.StateID)
	if err != nil {
		return nil, err
	}
	return &models.StateMerkleProof{
		UserState: &leaf.UserState,
		Witness:   witness,
	}, nil
}
