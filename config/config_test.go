package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	cfg, err := GetConfig("../config.template.yaml")
	assert.NoError(t, err)
	assert.Equal(
		t,
		&Config{
			Version:  "dev-0.1.0",
			Port:     8080,
			DBName:   "hubble_test",
			DBUser:   "hubble",
			DBPassword: "root",
		},
		cfg,
	)
}
