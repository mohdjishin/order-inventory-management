package meta

import (
	"log"

	"github.com/mohdjishin/order-inventory-management/db"
)

var (
	CommitHash string
	BuildTime  string
)

func GetCommitHash() string {
	if CommitHash == "" {
		CommitHash = "development"

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
