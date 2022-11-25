package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

var Config ServerConfig

type ServerConfig struct {
	Postgres_uri string `yaml:"postgres_uri"`
	Loglevel     string `yaml:"loglevel"`
	DevMode      bool   `yaml:"dev_mode"`
}

func Init(yamlFilePath string) ServerConfig {
	yamlConfig := ServerConfig{}
	//Read and parse config file if exists
	if yamlFilePath != "" {
		yamlFile, err := os.ReadFile(yamlFilePath)
		if err != nil {
			panic(err)
		}

		err = yaml.UnmarshalStrict(yamlFile, &yamlConfig)
		if err != nil {
			panic(err)
		}
	}

	env_postgres_uri, _ := os.LookupEnv("POSTGRES_URI")
	if env_postgres_uri != "" {
		yamlConfig.Postgres_uri = env_postgres_uri
	}

	env_loglevel, _ := os.LookupEnv("LOGLEVEL")
	if env_postgres_uri != "" {
		yamlConfig.Loglevel = env_loglevel
	}

	Config = yamlConfig
	return yamlConfig

}
