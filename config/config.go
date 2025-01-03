package config

import (
	// "flag"
	// "os"
	"time"

	"github.com/spf13/viper"
	// "cleanenv"
)

const configPath = "config/config.yaml"

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	DB       DBConfig      `yaml:"db"`
	SSO      SSOConfig     `yaml:"sso"`
	Logger   LoggerConfig  `yaml:"logger"`
	TokenTTL time.Duration `yaml:"tokenTtl" env-default:"1h"`
}

type SSOConfig struct {
	GRPC GRPCConfig `yaml:"grpc"`
}

type DBConfig struct {
	StoragePath    string `yaml:"storagePath" env-required:"true"`
	MigrationPath  string `yaml:"migrationPath"`
	MigrationTable string `yaml:"migrationTable"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type LoggerConfig struct {
	Level string `yaml:"level"`
	File  string `yaml:"file"`
}

func New() (config *Config, err error) {
	viper.SetConfigFile(configPath)

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return config, err
}
