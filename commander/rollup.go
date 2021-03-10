package commander

import (
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/ethereum/go-ethereum/common"
)

func RollupLoop(cfg *config.Config) {
	storage, err := st.NewStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}
	stateTree := st.NewStateTree(storage)

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

		feeReceiver := models.MakeUint256(0) // TODO: Get from config

		includedTransactions, err := applyTransactions(stateTree, transactions, uint32(feeReceiver.Uint64()))
		if err != nil {
			log.Fatal(err)
		}

		commitment, err := CreateCommitment(stateTree, includedTransactions, &feeReceiver)
		if err != nil {
			log.Fatal(err)
		}

		err = storage.AddCommitment(commitment)
		if err != nil {
			log.Fatal(err)
		}

		for i := range includedTransactions {
			tx := includedTransactions[i]
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

	for i := range transactions {
		encoded, err := encoder.EncodeTransaction(&transactions[i])
		if err != nil {
			return nil, err
		}

		buf = append(buf, encoded...)
	}

	return buf, nil
}

// TODO: add tests
func applyTransactions(
	stateTree *st.StateTree,
	transactions []models.Transaction,
	feeReceiverIndex uint32,
) (
	[]models.Transaction,
	error,
) {
	validTxs := make([]models.Transaction, 0, 32)

	// TODO: handle too many transaction
	for i := range transactions {
		tx := transactions[i]
		txError, appError := ApplyTransfer(stateTree, &tx, feeReceiverIndex)
		if appError != nil {
			return nil, appError
		}
		if txError == nil {
			validTxs = append(validTxs, tx)
		}
	}

	return validTxs, nil
}

// TODO: Test me
func CreateCommitment(stateTree *st.StateTree, transactions []models.Transaction, feeReceiver *models.Uint256) (*models.Commitment, error) {
	combinedSignature := models.Signature{models.MakeUint256(1), models.MakeUint256(2)} // TODO: Actually combine signatures

	transactionsSerialized, err := serializeTransactions(transactions)
	if err != nil {
		return nil, err
	}

	accountRoot := common.Hash{} // TODO: Read from account tree

	bodyHash, err := encoder.GetCommitmentBodyHash(accountRoot, combinedSignature, *feeReceiver, transactionsSerialized)
	if err != nil {
		return nil, err
	}

	stateRoot, err := stateTree.Root()
	if err != nil {
		return nil, err
	}

	leafHash := st.HashTwo(*stateRoot, *bodyHash)

	commitment := models.Commitment{
		LeafHash:          leafHash,
		PostStateRoot:     *stateRoot,
		BodyHash:          *bodyHash,
		AccountTreeRoot:   accountRoot,
		CombinedSignature: combinedSignature,
		FeeReceiver:       *feeReceiver,
		Transactions:      transactionsSerialized,
	}

	return &commitment, err
}
