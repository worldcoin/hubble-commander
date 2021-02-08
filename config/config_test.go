package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	cfg, err := GetConfig("../config.template.yaml")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, &Config{Version: "dev-0.1.0", Port: 8080}, cfg)
}
