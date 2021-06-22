package executor

import (
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrNotEnoughTransfers = NewRollupError("not enough transfers")
)

func (t *TransactionExecutor) CreateTransferCommitments(
	domain *bls.Domain,
) (commitments []models.Commitment, err error) {
	pendingTransfers, err := t.storage.GetPendingTransfers(pendingTxsCountMultiplier * t.cfg.TxsPerCommitment)
	if err != nil {
		return nil, err
	}

	if len(pendingTransfers) < int(t.cfg.TxsPerCommitment) {
		return []models.Commitment{}, nil
	}

	commitments = make([]models.Commitment, 0, t.cfg.MaxCommitmentsPerBatch)

	for len(commitments) != int(t.cfg.MaxCommitmentsPerBatch) {
		var commitment *models.Commitment

		pendingTransfers, commitment, err = t.createTransferCommitment(pendingTransfers, domain)
		if err != nil {
			return nil, err
		}
		if commitment == nil {
			break
		}

		commitments = append(commitments, *commitment)
	}

	return commitments, nil
}

func (t *TransactionExecutor) createTransferCommitment(
	pendingTransfers []models.Transfer,
	domain *bls.Domain,
) (
	[]models.Transfer,
	*models.Commitment,
	error,
) {
	startTime := time.Now()

	preparedTransfers, err := t.prepareTransfers(pendingTransfers)
	if err != nil {
		return nil, nil, err
	}

	commitment, err := t.prepareTransferCommitment(preparedTransfers, domain)
	if err != nil {
		return nil, nil, err
	}

	err = t.markTransfersAsIncluded(preparedTransfers.appliedTransfers, commitment.ID)
	if err != nil {
		return nil, nil, err
	}

	log.Printf(
		"Created a %s commitment from %d transactions in %s",
		txtype.Transfer,
		len(preparedTransfers.appliedTransfers),
		time.Since(startTime).Round(time.Millisecond).String(),
	)

	return preparedTransfers.newPendingTransfers, commitment, nil
}

func removeTransfer(transferList, toRemove []models.Transfer) []models.Transfer {
	outputIndex := 0
	for i := range transferList {
		transfer := &transferList[i]
		if !transferExists(toRemove, transfer) {
			transferList[outputIndex] = *transfer
			outputIndex++
		}
	}

	return transferList[:outputIndex]
}

func transferExists(transferList []models.Transfer, tx *models.Transfer) bool {
	for i := range transferList {
		if transferList[i].Hash == tx.Hash {
			return true
		}
	}
	return false
}

func combineTransferSignatures(transfers []models.Transfer, domain *bls.Domain) (*models.Signature, error) {
	signatures := make([]*bls.Signature, 0, len(transfers))
	for i := range transfers {
		sig, err := bls.NewSignatureFromBytes(transfers[i].Signature.Bytes(), *domain)
		if err != nil {
			return nil, err
		}
		signatures = append(signatures, sig)
	}
	return bls.NewAggregatedSignature(signatures).ModelsSignature(), nil
}

func (t *TransactionExecutor) markTransfersAsIncluded(transfers []models.Transfer, commitmentID int32) error {
	hashes := make([]common.Hash, 0, len(transfers))
	for i := range transfers {
		hashes = append(hashes, transfers[i].Hash)
	}
	return t.storage.BatchMarkTransactionAsIncluded(hashes, &commitmentID)
}
