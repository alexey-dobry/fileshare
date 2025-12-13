package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/alexey-dobry/fileshare/pkg/logger/zap"
	"github.com/alexey-dobry/fileshare/pkg/validator"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/domain/jwt"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/server/grpc"
	"github.com/alexey-dobry/fileshare/services/auth_service/internal/store/authdb"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Logger zap.Config    `yaml:"logger"`
	GRPC   grpc.Config   `yaml:"grpc"`
	Store  authdb.Config `yaml:"store"`
	JWT    jwt.Config    `yaml:"jwt"`
}

func MustLoad() Config {
	var cfg Config
	configPath := ParseFlag(cfg)

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		errMsg := fmt.Sprintf("Failed to read config on path(%s): %s", configPath, err)
		panic(errMsg)
	}

	if err := validator.V.Struct(&cfg); err != nil {
		errMsg := fmt.Sprintf("Failed to validate config: %s", err)
		panic(errMsg)
	}

	return cfg
}

func ParseFlag(cfg Config) string {
	configPath := flag.String("config", "./configs/config.yaml", "config file path")
	configHelp := flag.Bool("help", false, "show configuration help")

	if *configHelp {
		headerText := "Configuration options:"
		help, err := cleanenv.GetDescription(&cfg, &headerText)
		if err != nil {
			errMsg := fmt.Sprintf("error getting configuration description: %s", err)
			panic(errMsg)
		}
		fmt.Println(help)
		os.Exit(0)
	}

	return *configPath
}
