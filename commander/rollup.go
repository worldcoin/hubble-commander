package commander

import (
	"log"
	"time"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/encoder"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/Worldcoin/hubble-commander/testutils/deployer"
	simulator2 "github.com/Worldcoin/hubble-commander/testutils/simulator"
	"github.com/ethereum/go-ethereum/common"
)

func RollupLoop(cfg *config.Config) {
	storage, err := st.NewStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}
	stateTree := st.NewStateTree(storage)

	err = PopulateGenesisAccounts(stateTree, []GenesisAccount{
		{
			accountIndex: 0,
			balance:      models.MakeUint256(1000),
		},
		{
			accountIndex: 1,
			balance:      models.MakeUint256(1000),
		},
		{
			accountIndex: 2,
			balance:      models.MakeUint256(1000),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	simulator, err := simulator2.NewAutominingSimulator()
	if err != nil {
		log.Fatal(err)
	}

	stateRoot, err := stateTree.Root()
	if err != nil {
		log.Fatal(err)
	}

	contracts, err := deployer.DeployConfiguredRollup(simulator, deployer.DeploymentConfig{
		GenesisStateRoot: stateRoot,
	})
	if err != nil {
		log.Fatal(err)
	}

	client := eth.NewTestClient(simulator.Account, contracts.Rollup)

	for {
		transactions, err := storage.GetPendingTransactions()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("%d transactions in the pool", len(transactions))

		if len(transactions) < 2 { // TODO: change to 32 transactions
			time.Sleep(500 * time.Millisecond)
			continue
		}

		feeReceiver := models.MakeUint256(0) // TODO: Get from config

		log.Printf("Applying %d transactions", len(transactions))

		includedTransactions, err := ApplyTransactions(stateTree, transactions, uint32(feeReceiver.Uint64()))
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Creating a commitment from %d transactions", len(includedTransactions))

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

		_, err = client.SubmitTransfersBatch([]*models.Commitment{commitment})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Sumbmited commitment %s on chain", commitment.LeafHash.Hex())

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
