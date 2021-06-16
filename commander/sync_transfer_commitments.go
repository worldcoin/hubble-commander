package commander

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
)

func (t *transactionExecutor) syncTransferCommitments(batch *eth.DecodedBatch) error {
	for i := range batch.Commitments {
		err := t.syncTransferCommitment(batch, &batch.Commitments[i])
		if err != nil {
			if err == ErrInvalidSignature { //nolint: staticcheck
				//TODO: dispute fraudulent commitment
			}
			return err
		}
	}
	return nil
}

func (t *transactionExecutor) syncTransferCommitment(
	batch *eth.DecodedBatch,
	commitment *encoder.DecodedCommitment,
) error {
	deserializedTransfers, err := encoder.DeserializeTransfers(commitment.Transactions)
	if err != nil {
		return err
	}

	transfers, err := t.ApplyTransfers(deserializedTransfers)
	if err != nil {
		return err
	}

	if len(transfers.invalidTransfers) > 0 {
		return ErrFraudulentTransfer
	}
	if len(transfers.appliedTransfers) != len(deserializedTransfers) {
		return ErrTransfersNotApplied
	}

	isValid, err := t.verifyTransferSignature(commitment, transfers.appliedTransfers)
	if err != nil {
		return err
	}
	if !isValid {
		return ErrInvalidSignature
	}

	commitmentID, err := t.storage.AddCommitment(&models.Commitment{
		Type:              batch.Type,
		Transactions:      commitment.Transactions,
		FeeReceiver:       commitment.FeeReceiver,
		CombinedSignature: commitment.CombinedSignature,
		PostStateRoot:     commitment.StateRoot,
		IncludedInBatch:   &batch.ID,
	})
	if err != nil {
		return err
	}
	for i := range transfers.appliedTransfers {
		transfers.appliedTransfers[i].IncludedInCommitment = commitmentID
	}

	for i := range transfers.appliedTransfers {
		hashTransfer, err := encoder.HashTransfer(&transfers.appliedTransfers[i])
		if err != nil {
			return err
		}
		transfers.appliedTransfers[i].Hash = *hashTransfer
	}

	return t.storage.BatchAddTransfer(transfers.appliedTransfers)
}

func (t *transactionExecutor) verifyTransferSignature(commitment *encoder.DecodedCommitment, transfers []models.Transfer) (bool, error) {
	domain, err := t.storage.GetDomain(t.client.ChainState.ChainID)
	if err != nil {
		return false, err
	}
	blsDomain, err := bls.DomainFromBytes(domain[:])
	if err != nil {
		return false, err
	}

	messages := make([][]byte, len(transfers))
	publicKeys := make([]*models.PublicKey, len(transfers))
	for i := range transfers {
		publicKeys[i], err = t.storage.GetPublicKeyByStateID(transfers[i].FromStateID)
		if err != nil {
			return false, err
		}
		messages[i], err = encoder.EncodeTransferForSigning(&transfers[i])
		if err != nil {
			return false, err
		}
	}

	sig, err := bls.NewSignatureFromBytes(commitment.CombinedSignature[:], *blsDomain)
	if err != nil {
		return false, err
	}
	signature := bls.AggregatedSignature{Signature: sig}
	return signature.Verify(messages, publicKeys)
}
