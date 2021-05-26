package config

import (
	"os"
	"testing"

	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestGetEnvOrDefault(t *testing.T) {
	_ = os.Setenv("HUBBLE_FOO", "foo")

	require.Equal(t, "foo", *getEnvOrDefault("HUBBLE_FOO", ref.String("bar")))
	require.Equal(t, "bar", *getEnvOrDefault("HUBBLE_BAR", ref.String("bar")))
}

func TestGetFromViperOrDefault(t *testing.T) {
	_ = os.Setenv("HUBBLE_VERSION", "env")

	viper.AutomaticEnv()
	require.Equal(t, "env", *getFromViperOrDefault("hubble_version", ref.String("default")))
	require.Equal(t, "default", *getFromViperOrDefault("hubble_default", ref.String("default")))
}
