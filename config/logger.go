package config

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func GetCommanderConfigAndSetupLogger() *Config {
	cfg := GetConfig()
	setupCommanderLogger(cfg)
	logConfig(cfg)
	return cfg
}

func GetDeployerConfigAndSetupLogger() *DeployerConfig {
	cfg := GetDeployerConfig()
	setupDeployerLogger()
	logConfig(cfg)
	return cfg
}

func setupCommanderLogger(cfg *Config) {
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

func logConfig(cfg interface{}) {
	jsonCfg, err := json.Marshal(cfg)
	if err != nil {
		log.Panicf("%+v", errors.WithStack(err))
	}
	log.Debugf("Loaded config: %s", string(jsonCfg))
}
