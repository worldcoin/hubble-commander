package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateConfig(t *testing.T) {
	cfg := CreateConfig(
		"dev-0.1.0",
		8080,
		"hubble_test",
		"hubble_admin",
		"root_password",
	)

	assert.Equal(
		t,
		&Config{
			Version:    "dev-0.1.0",
			Port:       8080,
			DBName:     "hubble_test",
			DBUser:     "hubble_admin",
			DBPassword: "root_password",
		},
		cfg,
	)
}

func TestGetConfig(t *testing.T) {
	os.Setenv("HUBBLE_VERSION", "dev-0.1.0")
	os.Setenv("HUBBLE_PORT", "8080")
	os.Setenv("HUBBLE_DBNAME", "hubble_test")
	os.Setenv("HUBBLE_DBUSER", "hubble_admin")
	os.Setenv("HUBBLE_DBPASSWORD", "root")

	cfg := GetConfig()

	assert.Equal(
		t,
		&Config{
			Version:    "dev-0.1.0",
			Port:       8080,
			DBName:     "hubble_test",
			DBUser:     "hubble_admin",
			DBPassword: "root",
		},
		cfg,
	)
}
