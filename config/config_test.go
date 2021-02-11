package config

import (
	"os"
	"testing"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGetEnvOrDefault(t *testing.T) {
	viper.SetEnvPrefix("hubble")
	os.Setenv("HUBBLE_FOO", "foo")
	
	assert.Equal(t, "foo", *getEnvOrDefault("foo", utils.MakeStringPointer("bar")))
	assert.Equal(t, "bar", *getEnvOrDefault("bar", utils.MakeStringPointer("bar")))
}
