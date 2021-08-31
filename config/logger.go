package config

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func GetConfigAndSetupLogger() *Config {
	cfg := GetConfig()
	setupLogger(cfg)
	logConfig(cfg)
	return cfg
}

func setupLogger(cfg *Config) {
	if cfg.Log.Format == LogFormatJSON {
		log.SetFormatter(&log.JSONFormatter{})
	}
	log.SetOutput(os.Stdout)
	log.SetLevel(cfg.Log.Level)
}

func logConfig(cfg *Config) {
	jsonCfg, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}
	log.Debugf("Loaded config: %s", string(jsonCfg))
}
