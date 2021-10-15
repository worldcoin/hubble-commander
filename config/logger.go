package config

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func GetCommanderConfigAndSetupLogger() *Config {
	cfg := GetConfig()
	setupLogger(cfg)
	logCommanderConfig(cfg)
	return cfg
}

func GetDeployerConfigAndSetupLogger() *DeployerConfig {
	cfg := GetDeployerConfig()
	setupDeployerLogger()
	logDeployerConfig(cfg)
	return cfg
}

func setupLogger(cfg *Config) {
	if cfg.Log.Format == LogFormatJSON {
		log.SetFormatter(&log.JSONFormatter{})
	}
	log.SetOutput(os.Stdout)
	log.SetLevel(cfg.Log.Level)
}

func setupDeployerLogger() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func logCommanderConfig(cfg *Config) {
	jsonCfg, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}
	log.Debugf("Loaded config: %s", string(jsonCfg))
}

func logDeployerConfig(cfg *DeployerConfig) {
	jsonCfg, err := json.Marshal(cfg)
	if err != nil {
		log.Fatalf("%+v", errors.WithStack(err))
	}
	log.Debugf("Loaded config: %s", string(jsonCfg))
}
