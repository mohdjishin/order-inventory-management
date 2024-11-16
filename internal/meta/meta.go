package meta

import (
	"log"

	"github.com/mohdjishin/order-inventory-management/db"
)

var (
	Version    string
	CommitHash string
	BuildTime  string
)

func GetVersion() string {
	if Version == "" {
		Version = "development"
	}
	return Version
}

func GetCommitHash() string {
	if CommitHash == "" {
		CommitHash = "unknown"
	}
	return CommitHash
}

func GetBuildTime() string {
	if BuildTime == "" {
		BuildTime = "unknown"
	}
	return BuildTime
}

func GetDatabaseStats() string {

	sqlDB, err := db.GetDb().DB() // Get the underlying *sql.DB
	if err != nil {
		log.Fatalf("Failed to get database handle: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database connection is not active: %v", err)
		return "Connection failed"
	}
	return "Ok"
}
