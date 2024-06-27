package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Port   int    `yaml:"port" env-default:"8080"`
	Listen string `yaml:"listen" env-default:"0.0.0.0"`
}

type AgentConfig struct {
	ServerUrl      string `yaml:"serverUrl" env-default:"http://localhost:8080"`
	PollInterval   int    `yaml:"pollInterval" env-default:"2"`
	ReportInterval int    `yaml:"reportInterval" env-default:"10"`
}

type Config struct {
	Environment   string `yaml:"Environment" env-default:"local"`
	StorageDriver string `yaml:"StorageDriver" env-default:"memory"`
	Server        ServerConfig
	Agent         AgentConfig
}

func MustLoad() Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}
	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic("failed to read config: " + err.Error())
	}
	return config
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "./config/local.yaml", "Path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}
