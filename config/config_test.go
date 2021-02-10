package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("HUBBLE_VERSION", "dev-0.1.0")
	os.Setenv("HUBBLE_PORT", "8080")
	os.Setenv("HUBBLE_DBNAME", "hubble_test")
	os.Setenv("HUBBLE_DBUSER", "hubble")
	os.Setenv("HUBBLE_DBPASSWORD", "root")
	os.Exit(m.Run())
}

func TestGetConfig(t *testing.T) {
	cfg := GetConfig()
	assert.Equal(
		t,
		&Config{
			Version:    "dev-0.1.0",
			Port:       8080,
			DBName:     "hubble_test",
			DBUser:     "hubble",
			DBPassword: "root",
		},
		cfg,
	)
}
