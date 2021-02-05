package api

import (
	"testing"

	"github.com/Worldcoin/hubble-commander/config"
	"github.com/stretchr/testify/assert"
)

func TestSayHello(t *testing.T) {
	api := Api{&config.Config{Version: "v0123"}}
	assert.Equal(t, "v0123", api.GetVersion())
}
