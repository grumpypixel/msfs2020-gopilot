package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ConnectionName      string `yaml:"connection_name" env:"CONNECTION_NAME" env-default:"MyConnection"`
	ConnectionTimeout   int64  `yaml:"connection_timeout" env:"CONNECTION_TIMEOUT" env-default:"600"`
	DLLSearchPath       string `yaml:"dll_search_path" env:"DLL_SEARCH_PATH" env-default:"."`
	ServerAddress       string `yaml:"server_address" env:"SERVER_ADDRESS" env-default:"0.0.0.0:8888"`
	DataRequestInterval int64  `yaml:"data_request_interval" env:"DATA_REQUEST_INTERVAL" env-default:"200"`
	Verbose             bool   `yaml:"verbose" env:"VERBOSE" env-default:"false"`
}

func NewConfig(path string) (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, err
	}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
