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
	ErrNotEnoughC2Transfers = NewRollupError("not enough create2transfers")
)

func (t *TransactionExecutor) CreateCreate2TransferCommitments(
	domain *bls.Domain,
) (commitments []models.Commitment, err error) {
	pendingTransfers, err := t.storage.GetPendingCreate2Transfers(pendingTxsCountMultiplier * t.cfg.TxsPerCommitment)
	if err != nil {
		return nil, err
	}

	if len(pendingTransfers) < int(t.cfg.TxsPerCommitment) {
		return []models.Commitment{}, nil
	}

	commitments = make([]models.Commitment, 0, t.cfg.MaxCommitmentsPerBatch)

	for len(commitments) != int(t.cfg.MaxCommitmentsPerBatch) {
		var commitment *models.Commitment

		pendingTransfers, commitment, err = t.createC2TCommitment(pendingTransfers, domain)
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

func (t *TransactionExecutor) createC2TCommitment(
	pendingTransfers []models.Create2Transfer,
	domain *bls.Domain,
) (
	[]models.Create2Transfer,
	*models.Commitment,
	error,
) {
	startTime := time.Now()

	preparedTransfers, err := t.prepareCreate2Transfers(pendingTransfers)
	if err != nil {
		return nil, nil, err
	}
	if preparedTransfers == nil {
		return nil, nil, nil
	}

	commitment, err := t.prepareC2TCommitment(preparedTransfers, domain)
	if err != nil {
		return nil, nil, err
	}

	err = t.markCreate2TransfersAsIncluded(preparedTransfers.appliedTransfers, commitment.ID)
	if err != nil {
		return nil, nil, err
	}

	log.Printf(
		"Created a %s commitment from %d transactions in %s",
		txtype.Create2Transfer,
		len(preparedTransfers.appliedTransfers),
		time.Since(startTime).Round(time.Millisecond).String(),
	)

	return preparedTransfers.newPendingTransfers, commitment, nil
}

func removeCreate2Transfer(transferList, toRemove []models.Create2Transfer) []models.Create2Transfer {
	outputIndex := 0
	for i := range transferList {
		transfer := &transferList[i]
		if !create2TransferExists(toRemove, transfer) {
			transferList[outputIndex] = *transfer
			outputIndex++
		}
	}

	return transferList[:outputIndex]
}

func create2TransferExists(transferList []models.Create2Transfer, tx *models.Create2Transfer) bool {
	for i := range transferList {
		if transferList[i].Hash == tx.Hash {
			return true
		}
	}
	return false
}

func combineCreate2TransferSignatures(transfers []models.Create2Transfer, domain *bls.Domain) (*models.Signature, error) {
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

func (t *TransactionExecutor) markCreate2TransfersAsIncluded(transfers []models.Create2Transfer, commitmentID int32) error {
	hashes := make([]common.Hash, 0, len(transfers))
	for i := range transfers {
		hashes = append(hashes, transfers[i].Hash)

		err := t.storage.SetCreate2TransferToStateID(transfers[i].Hash, *transfers[i].ToStateID)
		if err != nil {
			return err
		}
	}
	return t.storage.BatchMarkTransactionAsIncluded(hashes, &commitmentID)
}
