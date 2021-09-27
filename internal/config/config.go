package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	ConnectionName      string `yaml:"connection_name" env:"CONNECTION_NAME" env-default:"MyConnection"`
	ConnectionTimeout   int64  `yaml:"connection_timeout" env:"CONNECTION_TIMEOUT" env-default:"600"`
	SimConnectDLLPath   string `yaml:"simconnect_dll_path" env:"SIMCONNECT_DLL_PATH" env-default:"."`
	ServerAddress       string `yaml:"server_address" env:"SERVER_ADDRESS" env-default:"0.0.0.0:8888"`
	DataRequestInterval int64  `yaml:"data_request_interval" env:"DATA_REQUEST_INTERVAL" env-default:"200"`
	LogLevel            string `yaml:"log_level" env:"LOG_LEVEL" env-default:"info"`
}

func NewConfigFromFile(path string) (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err == nil {
		return cfg, nil
	} else {
		log.Errorf("Config error: %s", err.Error())
	}
	log.Info("Creating config from environment variables")
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
