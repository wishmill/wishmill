package db

import (
	"database/sql"
	"time"
	"wishmill/internal/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

const healthCheckIntervalDb int = 5

var dbHealth bool = false

var Db *sql.DB

func Init(postgres_uri string) {
	//Get postgres uri from environment variable
	if postgres_uri == "" {
		logger.FatalLogger.Panicln("No postgres uri configured")
	}

	//Try database connection and ping
	var err error
	Db, err = sql.Open("postgres", postgres_uri)
	if err != nil {
		logger.FatalLogger.Println("Connection to database failed")
		panic(err)
	}
	err = Db.Ping()
	if err != nil {
		logger.FatalLogger.Println("Initial ping to database failed")
		panic(err)
	}
	go checkDbHealth()
	logger.InfoLogger.Println("Successfully connected to database")
}

// Ping DB. Return true if succeeded
func PingDB() bool {
	if err := Db.Ping(); err != nil {
		logger.WarningLogger.Println("db: Could not ping database")
		return false
	}
	logger.DebugLogger.Println("db: Successfully pinged database")
	return true
}

func checkDbHealth() {
	for {
		dbHealth = PingDB()
		time.Sleep(time.Duration(healthCheckIntervalDb) * time.Second)
	}
}

func GetDbHealth() bool {
	return dbHealth
}

func migrateDb() {
	driver, err := postgres.WithInstance(Db, &postgres.Config{})
	if err != nil {
		logger.FatalLogger.Println(err)
	}

	m, err := migrate.NewWithDatabaseInstance("file", "postgres", driver)
	if err != nil {
		logger.FatalLogger.Println(err)
	}

	m.Up()
}
