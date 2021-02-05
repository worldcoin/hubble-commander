package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	cfg, err := GetConfig("../config.template.yaml")
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, &Config{Version: "dev-0.1.0", Port: 8080}, cfg)
}
