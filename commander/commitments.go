package commander

import (
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
)

func CommitmentsEndlessLoop(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) error {
	done := make(chan bool)
	return CommitmentsLoop(storage, client, cfg, done)
}

func CommitmentsLoop(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig, done <-chan bool) error {
	ticker := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		case <-done:
			ticker.Stop()
			return nil
		case <-ticker.C:
			err := CommitTransactions(storage, client, cfg)
			if err != nil {
				return err
			}
		}
	}
}

func CommitTransactions(storage *st.Storage, client *eth.Client, cfg *config.RollupConfig) error {
	stateTree := st.NewStateTree(storage)

	// TODO: wrap in a db transaction
	txs, err := storage.GetPendingTransactions()
	if err != nil {
		return err
	}
	txsCount := uint32(len(txs))
	log.Printf("%d transactions in the pool", txsCount)

	if txsCount < cfg.TxsPerCommitment {
		return nil
	}

	log.Printf("Applying %d transactions", txsCount)
	includedTxs, err := ApplyTransactions(storage, txs, cfg)
	if err != nil {
		return err
	}
	// TODO: if len(includedTxs) != txCountPerCommitment then fail and rollback

	log.Printf("Creating a commitment from %d transactions", len(includedTxs))
	commitment, err := CreateCommitment(stateTree, includedTxs, cfg.FeeReceiverIndex)
	if err != nil {
		return err
	}

	commitmentID, err := storage.AddCommitment(commitment)
	if err != nil {
		return err
	}

	_, err = client.SubmitTransfersBatch([]*models.Commitment{commitment})
	if err != nil {
		return err
	}
	log.Printf("Sumbmited commitment %s on chain", commitment.LeafHash().Hex())

	for i := range includedTxs {
		tx := includedTxs[i]
		err = storage.MarkTransactionAsIncluded(tx.Hash, *commitmentID)
		if err != nil {
			return err
		}
	}
	// tx.Commit()
	return nil
}

// TODO: Test me
func CreateCommitment(stateTree *st.StateTree, txs []models.Transaction, feeReceiver uint32) (*models.Commitment, error) {
	combinedSignature := models.Signature{models.MakeUint256(1), models.MakeUint256(2)} // TODO: Actually combine signatures

	serializedTxs, err := serializeTransactions(txs)
	if err != nil {
		return nil, err
	}

	accountRoot := &common.Hash{} // TODO: Read from account tree

	stateRoot, err := stateTree.Root()
	if err != nil {
		return nil, err
	}

	return &models.Commitment{
		Transactions:      serializedTxs,
		FeeReceiver:       feeReceiver,
		CombinedSignature: combinedSignature,
		PostStateRoot:     *stateRoot,
		AccountTreeRoot:   accountRoot,
	}, nil
}

func serializeTransactions(txs []models.Transaction) ([]byte, error) {
	buf := make([]byte, 0, len(txs)*12)

	for i := range txs {
		encoded, err := encoder.EncodeTransaction(&txs[i])
		if err != nil {
			return nil, err
		}
		buf = append(buf, encoded...)
	}

	return buf, nil
}
