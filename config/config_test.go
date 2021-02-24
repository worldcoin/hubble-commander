package config

import (
	"os"
	"testing"

	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/stretchr/testify/require"
)

func TestGetEnvOrDefault(t *testing.T) {
	_ = os.Setenv("HUBBLE_FOO", "foo")

	require.Equal(t, "foo", *getEnvOrDefault("HUBBLE_FOO", ref.String("bar")))
	require.Equal(t, "bar", *getEnvOrDefault("HUBBLE_BAR", ref.String("bar")))
}
