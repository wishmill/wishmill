package main

import (
	"flag"
	"wishmill/internal/api"
	"wishmill/internal/config"
	"wishmill/internal/db"
	"wishmill/internal/logger"
)

func main() {
	var configFile string
	logger.Init()
	logger.InfoLogger.Println("Starting wishmill")
	flag.StringVar(&configFile, "config", "", "Path to config file")
	flag.Parse()
	serverConfig := config.Init(configFile)
	if serverConfig.DevMode {
		logger.InfoLogger.Println("Wishmill is running in dev mode")

	}
	logger.Init2(serverConfig.Loglevel, serverConfig.DevMode)
	db.Init(serverConfig.Postgres_uri)
	api.Init()
	select {}
}
