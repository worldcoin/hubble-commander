package syncer

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
)

const (
	InvalidSignatureMessage    = "invalid commitment signature"
	NotExistingReceiverMessage = "nonexistent receiver"
)

func (c *Context) verifyTxSignature(commitment *encoder.DecodedCommitment, txs models.GenericTransactionArray) error {
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

func (c *Context) verifyCommitmentSignature(
	signature *models.Signature,
	domain *bls.Domain,
	messages [][]byte,
	publicKeys []*models.PublicKey,
	txs models.GenericTransactionArray,
) error {
	if len(messages) == 0 {
		return nil
	}
	sig, err := bls.NewSignatureFromBytes(signature.Bytes(), *domain)
	if err != nil {
		return c.createDisputableSignatureError(err.Error(), txs)
	}
	aggregatedSignature := bls.AggregatedSignature{Signature: sig}
	isValid, err := aggregatedSignature.Verify(messages, publicKeys)
	if err != nil {
		return err
	}
	if !isValid {
		return c.createDisputableSignatureError(InvalidSignatureMessage, txs)
	}
	return nil
}

func (c *Context) createDisputableSignatureError(reason string, txs models.GenericTransactionArray) error {
	proofs, proofErr := c.StateMerkleProofs(txs)
	if proofErr != nil {
		return proofErr
	}
	return NewDisputableErrorWithProofs(Signature, reason, proofs)
}

func (c *Context) StateMerkleProofs(txs models.GenericTransactionArray) ([]models.StateMerkleProof, error) {
	proofs := make([]models.StateMerkleProof, 0, txs.Len())
	for i := 0; i < txs.Len(); i++ {
		stateProof, err := c.UserStateProof(txs.At(i).GetFromStateID())
		if err != nil {
			return nil, err
		}
		proofs = append(proofs, *stateProof)
	}
	return proofs, nil
}

func (c *Context) UserStateProof(stateID uint32) (*models.StateMerkleProof, error) {
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
