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

func RollupLoop(storage *st.Storage, client *eth.Client, cfg *config.Config) error {
	stateTree := st.NewStateTree(storage)

	for {
		// TODO: wrap in a db transaction
		transactions, err := storage.GetPendingTransactions()
		if err != nil {
			return err
		}
		log.Printf("%d transactions in the pool", len(transactions))

		if len(transactions) < 2 { // TODO: change to 32 transactions
			time.Sleep(500 * time.Millisecond)
			continue
		}

		feeReceiver := cfg.FeeReceiverIndex

		log.Printf("Applying %d transactions", len(transactions))
		includedTransactions, err := ApplyTransactions(storage, transactions, feeReceiver)
		if err != nil {
			return err
		}
		// TODO: if len(includedTransactions) != txCountPerCommitment then fail and rollback

		log.Printf("Creating a commitment from %d transactions", len(includedTransactions))
		commitment, err := CreateCommitment(stateTree, includedTransactions, feeReceiver)
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

		time.Sleep(500)
	}
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
