package scripts

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
)

func Test_exportAccounts(t *testing.T) {
	storage, err := st.NewTestStorage()
	require.NoError(t, err)

	defer func() {
		err = storage.Teardown()
		require.NoError(t, err)
	}()

	expectedLeaves, err := addAccounts(storage.Storage)
	require.NoError(t, err)

	file, err := os.CreateTemp("", "export_account_leaves_test")
	require.NoError(t, err)
	defer func() {
		err = os.Remove(file.Name())
		require.NoError(t, err)
	}()

	exportedLeaves, err := exportData(storage.Storage, file, exportAndCountAccounts)
	require.NoError(t, err)
	require.Len(t, expectedLeaves, exportedLeaves)

	err = file.Close()
	require.NoError(t, err)

	leaves := readAccountsFromFile(t, file.Name())
	require.Equal(t, expectedLeaves, leaves)
}

func readAccountsFromFile(t *testing.T, fileName string) []models.AccountLeaf {
	bytes, err := os.ReadFile(fileName)
	require.NoError(t, err)

	leaves := make([]models.AccountLeaf, 0, 4)
	err = json.Unmarshal(bytes, &leaves)
	require.NoError(t, err)

	return leaves
}

func addAccounts(storage *st.Storage) ([]models.AccountLeaf, error) {
	leaves := make([]models.AccountLeaf, 0, 4)
	for i := uint32(0); i < 4; i++ {
		leaf := &models.AccountLeaf{
			PubKeyID:  i,
			PublicKey: models.PublicKey{1, 2, byte(i)},
		}
		err := storage.AccountTree.SetSingle(leaf)
		if err != nil {
			return nil, err
		}

		leaves = append(leaves, *leaf)
	}
	return leaves, nil
}
