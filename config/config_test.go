package config

import (
	"os"
	"testing"

	"github.com/Worldcoin/hubble-commander/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetEnvOrDefault(t *testing.T) {
	os.Setenv("HUBBLE_FOO", "foo")

	assert.Equal(t, "foo", *getEnvOrDefault("HUBBLE_FOO", utils.String("bar")))
	assert.Equal(t, "bar", *getEnvOrDefault("HUBBLE_BAR", utils.String("bar")))
}
