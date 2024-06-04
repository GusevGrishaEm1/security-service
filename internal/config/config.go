package config

import (
	"errors"
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string     `yaml:"env" env-default:"local"`
	StoragePath string     `yaml:"storage_path" env-required:"true"`
	TokenTTL    int        `yaml:"token_ttl" env-default:"1h"`
	SecretKey   string     `yaml:"secret_key" env-required:"true"`
	GRPC        GRPCConfig `yaml:"grpc"`
}

type GRPCConfig struct {
	Host    string `yaml:"host" env-default:"localhost"`
	Port    string `yaml:"port" env-default:"50051"`
	Timeout int    `yaml:"timeout" env-default:"5"`
}

func Load() (*Config, error) {
	var res string

	flag.StringVar(&res, "config", "", "config file path")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_FILE")
	}

	if res == "" {
		return nil, errors.New("config file path is empty")
	}

	var config Config

	if err := cleanenv.ReadConfig(res, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
