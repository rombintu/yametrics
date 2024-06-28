package config

import (
	"flag"
)

type ServerConfig struct {
	Address       string `yaml:"address" env-default:"localhost:8080"`
	Environment   string `yaml:"Environment" env-default:"local"`
	StorageDriver string `yaml:"StorageDriver" env-default:"memory"`
}

type AgentConfig struct {
	ServerUrl      string `yaml:"serverUrl" env-default:"http://localhost:8080"`
	PollInterval   int64  `yaml:"pollInterval" env-default:"2"`
	ReportInterval int64  `yaml:"reportInterval" env-default:"10"`
	Mode           string `yaml:"mode" env-default:"debug"`
}

// type Config struct {
// 	Environment   string `yaml:"Environment" env-default:"local"`
// 	StorageDriver string `yaml:"StorageDriver" env-default:"memory"`
// 	Server        ServerConfig
// 	Agent         AgentConfig
// }

// func MustLoad() Config {
// path := fetchConfigPath()
// if path == "" {
// 	panic("config path is empty")
// }

// if _, err := os.Stat(path); os.IsNotExist(err) {
// 	panic("config file does not exist: " + path)
// }
// var config Config
// if err := cleanenv.ReadConfig(path, &config); err != nil {
// 	panic("failed to read config: " + err.Error())
// }
// return config
// 	config := tryLoadFromFlags()
// 	return config
// }

// func fetchConfigPath() string {
// 	var res string
// 	flag.StringVar(&res, "config", "./config/local.yaml", "Path to config file")
// 	flag.Parse()
// 	if res == "" {
// 		res = os.Getenv("CONFIG_PATH")
// 	}
// 	return res
// }

// func getenv(name string, defaultValue string) string {
// 	value, err := os.LookupEnv(name)
// 	if !err {
// 		return ""
// 	}
// 	return value
// }

// func tryLoadFromEnv() Config {
// 	var config Config
// 	address := getenv("ADDRESS", "localhost:8080")

// 	return config
// }

// Try load Server Config from flags
func LoadServerConfigFromFlags() ServerConfig {
	var config ServerConfig
	e := flag.String("environment", "local", "Environment")
	s := flag.String("storageDriver", "memory", "Storage driver")
	a := flag.String("a", "localhost:8080", "Server address")
	flag.Parse()

	config.Address = *a
	config.StorageDriver = *s
	config.Environment = *e
	return config
}

// Try load Server Config from flags
func LoadAgentConfigFromFlags() AgentConfig {
	var config AgentConfig
	a := flag.String("a", "localhost:8080", "Server address")
	r := flag.Int64("r", 10, "Report interval")
	p := flag.Int64("p", 2, "Poll interval")
	flag.Parse()

	config.ServerUrl = *a
	config.ReportInterval = *r
	config.PollInterval = *p

	return config
}
