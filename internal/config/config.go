package config

import (
	"errors"
	"os"
	"wishmill/internal/logger"

	"gopkg.in/yaml.v2"
)

var Config ServerConfig

type ServerConfig struct {
	Postgres_uri  string          `yaml:"postgres_uri"`
	Loglevel      string          `yaml:"loglevel"`
	DevMode       bool            `yaml:"dev_mode"`
	Token_secret  string          `yaml:"token_secret"`
	Oidc_provider []Oidc_provider `yaml:"oidc_providers" json:"oidc_providers"`
}

type Oidc_provider struct {
	Name           string `yaml:"name" json:"name"`
	Url            string `yaml:"url" json:"url"`
	Client_id      string `yaml:"client_id" json:"client_id"`
	Client_secret  string `yaml:"client_secret" json:"-"`
	Username_claim string `yaml:"username_claim" json:"-"`
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

	env_token_secret, _ := os.LookupEnv("TOKEN_SECRET")
	if env_postgres_uri != "" {
		yamlConfig.Token_secret = env_token_secret
	}

	if len(yamlConfig.Oidc_provider) == 0 {
		logger.FatalLogger.Panicln(errors.New("no oidc providers configured"))
	}

	for _, p := range yamlConfig.Oidc_provider {
		if (p.Username_claim != "email") && (p.Username_claim != "preferred_username") && (p.Username_claim != "username") && (p.Username_claim != "sub") && (p.Username_claim != "name") {
			logger.FatalLogger.Panicln(errors.New("username_claim has to be email, preferred_username, username, sub or name"))
		}
	}

	Config = yamlConfig
	return yamlConfig

}
