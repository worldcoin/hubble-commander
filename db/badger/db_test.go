package badger

import (
	"path"
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/stretchr/testify/require"
)

func TestDatabase_Clone(t *testing.T) {
	cfg := config.GetTestConfig().Badger

	primary, err := NewDatabase(cfg)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, primary.Prune())
	}()

	err = primary.Insert(testStruct.Name, testStruct)
	require.NoError(t, err)

	clonedConfig := cfg
	clonedConfig.Path = path.Join(cfg.Path, "cloned")
	cloned, err := primary.Clone(clonedConfig)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, cloned.Prune())
	}()

	var value someStruct
	err = cloned.Get(testStruct.Name, &value)
	require.NoError(t, err)
	require.Equal(t, testStruct, value)
}
