package commander

import (
	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
)

func (t *transactionExecutor) syncTransferCommitments(batch *eth.DecodedBatch) error {
	for i := range batch.Commitments {
		if err := t.syncTransferCommitment(batch, &batch.Commitments[i]); err != nil {
			return err
		}
	}
	return nil
}

func (t *transactionExecutor) syncTransferCommitment(
	batch *eth.DecodedBatch,
	commitment *encoder.DecodedCommitment,
) error {
	deserializedTransfers, transferMessages, err := encoder.DeserializeTransfers(commitment.Transactions)
	if err != nil {
		return err
	}
	_, err = t.verifySignature(commitment, deserializedTransfers, transferMessages)
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

func (t *transactionExecutor) verifySignature(commitment *encoder.DecodedCommitment, transfers []models.Transfer, messages [][]byte) (bool, error) {
	domain, err := t.storage.GetDomain(t.client.ChainState.ChainID)
	if err != nil {
		return false, err
	}
	blsDomain, err := bls.DomainFromBytes(domain[:])
	if err != nil {
		return false, err
	}

	publicKeys := make([]*models.PublicKey, 0, len(transfers))
	for i := range transfers {
		publicKey, err := t.storage.GetPublicKeyByStateID(transfers[i].FromStateID)
		if err != nil {
			return false, err
		}
		publicKeys = append(publicKeys, publicKey)
	}

	sig, err := bls.NewSignatureFromBytes(commitment.CombinedSignature[:], *blsDomain)
	if err != nil {
		return false, err
	}
	signature := bls.AggregatedSignature{Signature: sig}
	return signature.Verify(messages, publicKeys)
}
