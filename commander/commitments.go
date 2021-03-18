package commander

import (
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/utils"
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
	transactions, err := storage.GetPendingTransactions()
	if err != nil {
		return err
	}
	txsCount := uint32(len(transactions))
	log.Printf("%d transactions in the pool", txsCount)

	if txsCount < cfg.TxsPerCommitment {
		return nil
	}

	log.Printf("Applying %d transactions", txsCount)
	includedTransactions, err := ApplyTransactions(storage, transactions, cfg)
	if err != nil {
		return err
	}
	// TODO: if len(includedTransactions) != txCountPerCommitment then fail and rollback

	log.Printf("Creating a commitment from %d transactions", len(includedTransactions))
	commitment, err := CreateCommitment(stateTree, includedTransactions, cfg.FeeReceiverIndex)
	if err != nil {
		return err
	}

	err = storage.AddCommitment(commitment)
	if err != nil {
		return err
	}

	_, err = client.SubmitTransfersBatch([]*models.Commitment{commitment})
	if err != nil {
		return err
	}
	log.Printf("Sumbmited commitment %s on chain", commitment.LeafHash.Hex())

	for i := range includedTransactions {
		tx := includedTransactions[i]
		err = storage.MarkTransactionAsIncluded(tx.Hash, commitment.LeafHash)
		if err != nil {
			return err
		}
	}
	// tx.Commit()
	return nil
}

// TODO: Test me
func CreateCommitment(stateTree *st.StateTree, transactions []models.Transaction, feeReceiver uint32) (*models.Commitment, error) {
	combinedSignature := models.Signature{models.MakeUint256(1), models.MakeUint256(2)} // TODO: Actually combine signatures

	transactionsSerialized, err := serializeTransactions(transactions)
	if err != nil {
		return nil, err
	}

	accountRoot := common.Hash{} // TODO: Read from account tree

	bodyHash, err := encoder.GetCommitmentBodyHash(accountRoot, combinedSignature, feeReceiver, transactionsSerialized)
	if err != nil {
		return nil, err
	}

	stateRoot, err := stateTree.Root()
	if err != nil {
		return nil, err
	}

	leafHash := utils.HashTwo(*stateRoot, *bodyHash)

	commitment := models.Commitment{
		LeafHash:          leafHash,
		PostStateRoot:     *stateRoot,
		BodyHash:          *bodyHash,
		AccountTreeRoot:   accountRoot,
		CombinedSignature: combinedSignature,
		FeeReceiver:       feeReceiver,
		Transactions:      transactionsSerialized,
	}

	return &commitment, err
}

func serializeTransactions(transactions []models.Transaction) ([]byte, error) {
	buf := make([]byte, 0, len(transactions)*12)

	for i := range transactions {
		encoded, err := encoder.EncodeTransaction(&transactions[i])
		if err != nil {
			return nil, err
		}

		buf = append(buf, encoded...)
	}

	return buf, nil
}
