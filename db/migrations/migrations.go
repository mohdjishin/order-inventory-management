package migrations

import (
	"github.com/mohdjishin/order-inventory-management/db"
	"github.com/mohdjishin/order-inventory-management/internal/models"
	log "github.com/mohdjishin/order-inventory-management/logger"
	"go.uber.org/zap"
)

func Run() error {
	log.Info("Running database migrations")

	err := db.GetDb().AutoMigrate(&models.Product{}, &models.Order{}, &models.Inventory{}, &models.PricingHistory{}, &models.User{})
	if err != nil {
		log.Fatal("failed to migrate database", zap.Error(err))
		return err
	}
	log.Info("Database migration successful")
	return nil
}
