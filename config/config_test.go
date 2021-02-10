package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Setenv("Version", "dev-0.1.0")
	os.Setenv("Port", "8080")
	os.Setenv("DBName", "hubble_test")
	os.Setenv("DBUser", "hubble")
	os.Setenv("DBPassword", "root")
	os.Exit(m.Run())
}

func TestGetConfig(t *testing.T) {
	cfg, err := GetConfig()
	fmt.Println(cfg)
	assert.NoError(t, err)
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
