package config

import "github.com/ilyakaznacheev/cleanenv"

import (
	"flag"
	"os"
	"time"
)


type Config struct {
	Env 		string 		  `yaml:"env"`
	StoragePath string 		  `yaml:"storage_path"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC 		GRPCConfig    `yaml:"grpc" env-required:"true"`
}


type GRPCConfig struct {
	Port 	int           `yaml:"port"`
	TimeOut time.Duration `yaml:"timeout"`
}


func MustLoad() *Config {

	path := "/home/anton/sso/config/local.yaml" //fetchConfigPath()

	// if path == "" {
	// 	panic("config path is empty")
	// }

	// if _, err := os.Stat(path); os.IsNotExist(err) {
	// 	panic("config file does not exist: " + path)
	// }

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func MustLoadPath(configPath string) *Config {
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}


func fetchConfigPath() string { 
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}