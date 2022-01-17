package scripts

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/Worldcoin/hubble-commander/models"
	st "github.com/Worldcoin/hubble-commander/storage"
	"github.com/stretchr/testify/require"
)

func Test_exportStateLeaves(t *testing.T) {
	storage, err := st.NewTestStorage()
	require.NoError(t, err)

	defer func() {
		err = storage.Teardown()
		require.NoError(t, err)
	}()

	expectedLeaves, err := addStateLeaves(storage.Storage)
	require.NoError(t, err)

	file, err := os.CreateTemp("", "export_state_leaves_test")
	require.NoError(t, err)
	defer func() {
		err = os.Remove(file.Name())
		require.NoError(t, err)
	}()

	exportedLeaves, err := exportStateLeaves(storage.Storage, file)
	require.NoError(t, err)
	require.Len(t, expectedLeaves, exportedLeaves)

	err = file.Close()
	require.NoError(t, err)

	leaves := getStateLeavesFromFile(t, file.Name())
	require.Equal(t, expectedLeaves, leaves)
}

func getStateLeavesFromFile(t *testing.T, fileName string) []models.StateLeaf {
	bytes, err := os.ReadFile(fileName)
	require.NoError(t, err)

	leaves := make([]models.StateLeaf, 0, 4)
	err = json.Unmarshal(bytes, &leaves)
	require.NoError(t, err)

	return leaves
}

func addStateLeaves(storage *st.Storage) ([]models.StateLeaf, error) {
	userState := &models.UserState{
		PubKeyID: 0,
		TokenID:  models.MakeUint256(0),
		Balance:  models.MakeUint256(1000),
		Nonce:    models.MakeUint256(0),
	}

	leaves := make([]models.StateLeaf, 0, 4)
	for i := uint32(0); i < 4; i++ {
		userState.PubKeyID = i
		_, err := storage.StateTree.Set(i, userState)
		if err != nil {
			return nil, err
		}

		leaves = append(leaves, models.StateLeaf{
			StateID:   i,
			UserState: *userState,
		})
	}
	return leaves, nil
}
