package commander

import (
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
)

func createCreate2TransferCommitments(
	pendingTransfers []models.Create2Transfer,
	storage *st.Storage,
	client *eth.Client,
	cfg *config.RollupConfig,
	domain bls.Domain,
) ([]models.Commitment, error) {
	stateTree := st.NewStateTree(storage)
	commitments := make([]models.Commitment, 0, 32)

	for {
		if len(commitments) >= int(cfg.MaxCommitmentsPerBatch) {
			break
		}
		if len(pendingTransfers) < int(cfg.TxsPerCommitment) {
			break
		}
		startTime := time.Now()

		initialStateRoot, err := stateTree.Root()
		if err != nil {
			return nil, err
		}

		appliedTxs, invalidTxs, addedPubKeyIDs, feeReceiverStateID, err := ApplyCreate2Transfers(storage, client, pendingTransfers, cfg)
		if err != nil {
			return nil, err
		}

		if len(appliedTxs) < int(cfg.TxsPerCommitment) {
			err = stateTree.RevertTo(*initialStateRoot)
			if err != nil {
				return nil, err
			}
			break
		}

		pendingTransfers = removeCreate2Transfer(pendingTransfers, append(appliedTxs, invalidTxs...))

		serializedTxs, err := encoder.SerializeCreate2Transfers(appliedTxs, addedPubKeyIDs)
		if err != nil {
			return nil, err
		}

		combinedSignature, err := combineCreate2TransferSignatures(appliedTxs, domain)
		if err != nil {
			return nil, err
		}

		log.Printf("Creating a %s commitment from %d transactions", txtype.Create2Transfer.String(), len(appliedTxs))
		commitment, err := createAndStoreCommitment(storage, txtype.Create2Transfer, *feeReceiverStateID, serializedTxs, combinedSignature)
		if err != nil {
			return nil, err
		}

		err = markCreate2TransfersAsIncluded(storage, appliedTxs, commitment.ID)
		if err != nil {
			return nil, err
		}
		err = setCreate2TransferToStateID(storage, appliedTxs)
		if err != nil {
			return nil, err
		}

		commitments = append(commitments, *commitment)
		log.Printf(
			"Created a %s commitment from %d transactions in %d ms",
			txtype.Create2Transfer,
			len(appliedTxs),
			time.Since(startTime).Milliseconds(),
		)
	}

	return commitments, nil
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

func combineCreate2TransferSignatures(transfers []models.Create2Transfer, domain bls.Domain) (*models.Signature, error) {
	signatures := make([]*bls.Signature, 0, len(transfers))
	for i := range transfers {
		sig, err := bls.NewSignatureFromBytes(transfers[i].Signature[:], domain)
		if err != nil {
			return nil, err
		}
		signatures = append(signatures, sig)
	}
	return bls.NewAggregatedSignature(signatures).ModelsSignature(), nil
}

func markCreate2TransfersAsIncluded(storage *st.Storage, transfers []models.Create2Transfer, commitmentID int32) error {
	hashes := make([]common.Hash, 0, len(transfers))
	for i := range transfers {
		hashes = append(hashes, transfers[i].Hash)
	}
	return storage.BatchMarkTransactionAsIncluded(hashes, commitmentID)
}

func setCreate2TransferToStateID(storage *st.Storage, transfers []models.Create2Transfer) error {
	for i := range transfers {
		err := storage.SetCreate2TransferToStateID(transfers[i].Hash, transfers[i].ToStateID)
		if err != nil {
			return err
		}
	}
	return nil
}
