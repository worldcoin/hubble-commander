package commander

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/bls"
	"github.com/Worldcoin/hubble-commander/config"
	"github.com/Worldcoin/hubble-commander/db/postgres"
	"github.com/Worldcoin/hubble-commander/eth"
	"github.com/Worldcoin/hubble-commander/models"
	"github.com/Worldcoin/hubble-commander/models/enums/txtype"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
)

func TestSyncBatches(t *testing.T) {
	db, err := postgres.NewTestDB()
	require.NoError(t, err)

	storage := st.NewTestStorage(db.DB)
	tree := st.NewStateTree(storage)
	cfg := &config.RollupConfig{
		MinCommitmentsPerBatch: 1,
		MaxCommitmentsPerBatch: 32,
	}

	client, err := eth.NewTestClient()
	require.NoError(t, err)

	seedDB(t, storage, tree)

	tx := models.Transfer{
		TransactionBase: models.TransactionBase{
			TxType:      txtype.Transfer,
			FromStateID: 0,
			Amount:      models.MakeUint256(400),
			Fee:         models.MakeUint256(0),
			Nonce:       models.MakeUint256(0),
			Signature:   *bls.MockSignature().ModelsSignature(),
		},
		ToStateID: 1,
	}
	err = storage.AddTransfer(&tx)
	require.NoError(t, err)

	commitments, err := createTransferCommitments([]models.Transfer{tx}, storage, cfg)
	require.NoError(t, err)
	require.Len(t, commitments, 1)

	err = submitBatch(txtype.Transfer, commitments, storage, client.Client, cfg)
	require.NoError(t, err)

	// Recreate database
	db, err = postgres.NewTestDB()
	require.NoError(t, err)
	storage = st.NewTestStorage(db.DB)
	tree = st.NewStateTree(storage)

	seedDB(t, storage, tree)

	err = SyncBatches(storage, client.Client, cfg)
	require.NoError(t, err)

	state0, err := tree.Leaf(0)
	require.NoError(t, err)
	require.Equal(t, models.MakeUint256(600), state0.Balance)

	state1, err := tree.Leaf(1)
	require.NoError(t, err)
	require.Equal(t, models.MakeUint256(400), state1.Balance)

	batches, err := storage.GetBatchesInRange(nil, nil)
	require.NoError(t, err)
	require.Len(t, batches, 1)
}

func seedDB(t *testing.T, storage *st.Storage, tree *st.StateTree) {
	err := storage.AddAccountIfNotExists(&models.Account{
		PubKeyID:  0,
		PublicKey: models.PublicKey{},
	})
	require.NoError(t, err)

	err = tree.Set(0, &models.UserState{
		PubKeyID:   0,
		TokenIndex: models.MakeUint256(0),
		Balance:    models.MakeUint256(1000),
		Nonce:      models.MakeUint256(0),
	})
	require.NoError(t, err)

	err = tree.Set(1, &models.UserState{
		PubKeyID:   0,
		TokenIndex: models.MakeUint256(0),
		Balance:    models.MakeUint256(0),
		Nonce:      models.MakeUint256(0),
	})
	require.NoError(t, err)
}
