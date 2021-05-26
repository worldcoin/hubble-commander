package config

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Worldcoin/hubble-commander/utils/ref"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestGetEnvOrDefault(t *testing.T) {
	_ = os.Setenv("HUBBLE_FOO", "foo")

	require.Equal(t, "foo", *getEnvOrDefault("HUBBLE_FOO", ref.String("bar")))
	require.Equal(t, "bar", *getEnvOrDefault("HUBBLE_BAR", ref.String("bar")))
}

func TestGetTestConfig(t *testing.T) {
	config := GetViperConfig()
	fmt.Printf("%+v\n", config.Ethereum)
	fmt.Println(viper.GetBool("rollup.sync_batches"))
	WatchConfig(config)

	oldConfig := GetConfig()
	time.Sleep(10 * time.Second)
	fmt.Println(config.API.Version)

	require.Equal(t, oldConfig, *config)
}
