package config

import (
	"os"
	"testing"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/require"
)

func TestGetEnvOrDefault(t *testing.T) {
	_ = os.Setenv("HUBBLE_FOO", "foo")

	require.Equal(t, "foo", *getEnvOrDefault("HUBBLE_FOO", utils.String("bar")))
	require.Equal(t, "bar", *getEnvOrDefault("HUBBLE_BAR", utils.String("bar")))
}
