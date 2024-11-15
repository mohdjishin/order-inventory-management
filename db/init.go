package db

import (
	"fmt"
	"sync"

	"github.com/mohdjishin/order-inventory-management/config"
	// "github.com/mohdjishin/order-inventory-management/internal/models"
	log "github.com/mohdjishin/order-inventory-management/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBManagerInterface interface {
	Connect()
	GetDB() *gorm.DB
}

type DBManager struct {
	db *gorm.DB
}

var (
	dbManagerInstance DBManagerInterface
	once              sync.Once
)

func SetDbManager(manager DBManagerInterface) {
	dbManagerInstance = manager
}

func GetDbManagerInstance() DBManagerInterface {
	once.Do(func() {
		if dbManagerInstance == nil {
			dbManagerInstance = &DBManager{}
			dbManagerInstance.Connect()
		}
	})
	return dbManagerInstance
}

func init() {
	_ = GetDbManagerInstance()
}

func (m *DBManager) Connect() {
	log.Info().Msg("Connecting to database")
	var err error
	fmt.Println("DSN: ", config.Get().DSN)
	m.db, err = gorm.Open(postgres.Open(config.Get().DSN), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

}

func (m *DBManager) GetDB() *gorm.DB {
	return m.db
}

func GetDb() *gorm.DB {
	return GetDbManagerInstance().GetDB()
}
