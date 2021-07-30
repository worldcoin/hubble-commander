package executor

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
)

func (t *TransactionExecutor) verifyTransferSignature(commitment *encoder.DecodedCommitment, transfers []models.Transfer) error {
	domain, err := t.client.GetDomain()
	if err != nil {
		return err
	}

	messages := make([][]byte, len(transfers))
	publicKeys := make([]*models.PublicKey, len(transfers))
	for i := range transfers {
		publicKeys[i], err = t.storage.GetPublicKeyByStateID(transfers[i].FromStateID)
		if err != nil {
			return err
		}
		messages[i], err = encoder.EncodeTransferForSigning(&transfers[i])
		if err != nil {
			return err
		}
	}

	genericTxs := models.TransferArray(transfers)
	return t.verifyCommitmentSignature(&commitment.CombinedSignature, domain, messages, publicKeys, genericTxs)
}

func (t *TransactionExecutor) verifyCreate2TransferSignature(
	commitment *encoder.DecodedCommitment,
	transfers []models.Create2Transfer,
) error {
	domain, err := t.client.GetDomain()
	if err != nil {
		return err
	}

	messages := make([][]byte, len(transfers))
	publicKeys := make([]*models.PublicKey, len(transfers))
	for i := range transfers {
		publicKeys[i], err = t.storage.GetPublicKeyByStateID(transfers[i].FromStateID)
		if err != nil {
			return err
		}
		messages[i], err = encoder.EncodeCreate2TransferForSigning(&transfers[i])
		if err != nil {
			return err
		}
	}

	genericTxs := models.Create2TransferArray(transfers)
	return t.verifyCommitmentSignature(&commitment.CombinedSignature, domain, messages, publicKeys, genericTxs)
}

func (t *TransactionExecutor) verifyCommitmentSignature(
	signature *models.Signature,
	domain *bls.Domain,
	messages [][]byte,
	publicKeys []*models.PublicKey,
	transfers models.GenericTransactionArray,
) error {
	sig, err := bls.NewSignatureFromBytes(signature.Bytes(), *domain)
	if err != nil {
		return t.createDisputableSignatureError(err.Error(), transfers)
	}
	aggregatedSignature := bls.AggregatedSignature{Signature: sig}
	isValid, err := aggregatedSignature.Verify(messages, publicKeys)
	if err != nil {
		return err
	}
	if !isValid {
		return t.createDisputableSignatureError(InvalidSignature, transfers)
	}
	return nil
}

func (t *TransactionExecutor) createDisputableSignatureError(reason string, transfers models.GenericTransactionArray) error {
	proofs, proofErr := t.stateMerkleProofs(transfers)
	if proofErr != nil {
		return proofErr
	}
	return NewDisputableErrorWithProofs(Signature, reason, proofs)
}
