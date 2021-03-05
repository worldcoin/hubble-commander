package main

import (
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/api"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
)

func main() {
	cfg := config.GetConfig()

	go RollupLoop(&cfg)

	log.Fatal(api.StartAPIServer(&cfg))
}

func RollupLoop(cfg *config.Config) {
	storage, err := st.NewStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	for {
		transactions, err := storage.GetPendingTransactions()
		if err != nil {
			log.Fatal(err)
		}

		println("txs = %d", len(transactions))

		if len(transactions) < 2 { // TODO: change to 32 transactions
			time.Sleep(500 * time.Millisecond)
			continue
		}
		includedTransactions := transactions[0:2]

		feeReceiver := models.MakeUint256(0) // TODO: Get from config

		combinedSignature := models.Signature{models.MakeUint256(1), models.MakeUint256(2)} // TODO: Actually combine signatures

		transactionsSerialized, err := serializeTransactions(includedTransactions)
		if err != nil {
			log.Fatal(err)
		}

		accountRoot := common.Hash{} // TODO: Read from account tree

		bodyHash, err := encoder.GetCommitmentBodyHash(accountRoot, combinedSignature, feeReceiver, transactionsSerialized)
		if err != nil {
			log.Fatal(err)
		}

		stateRoot := common.Hash{} // TODO: Perform state transition and get tree root

		leafHash := st.HashTwo(stateRoot, *bodyHash)

		commitment := models.Commitment{
			LeafHash:          leafHash,
			BodyHash:          *bodyHash,
			AccountTreeRoot:   accountRoot,
			CombinedSignature: combinedSignature,
			FeeReceiver:       feeReceiver,
			Transactions:      transactionsSerialized,
		}

		err = storage.AddCommitment(&commitment)
		if err != nil {
			log.Fatal(err)
		}

		for _, tx := range includedTransactions {
			err := storage.MarkTransactionAsIncluded(tx.Hash, commitment.LeafHash)
			if err != nil {
				log.Fatal(err)
			}

		}

		time.Sleep(500)
	}
}

func serializeTransactions(transactions []models.Transaction) ([]byte, error) {
	buf := make([]byte, 0, len(transactions)*12)

	for _, tx := range transactions {
		encoded, err := encoder.EncodeTransaction(&tx)
		if err != nil {
			return nil, err
		}

		buf = append(buf, encoded...)
	}

	return buf, nil
}
