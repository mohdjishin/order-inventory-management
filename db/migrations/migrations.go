package migrations

import (
	"github.com/mohdjishin/order-inventory-management/db"
	"github.com/mohdjishin/order-inventory-management/internal/models"
	log "github.com/mohdjishin/order-inventory-management/logger"
)

func Run() error {
	log.Info().Msg("Running database migrations")

	err := db.GetDb().AutoMigrate(&models.Product{}, &models.Order{}, &models.Inventory{}, &models.PricingHistory{}, &models.User{})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to migrate database")
		return err
	}
	log.Info().Msg("Database migration successful")
	return nil
}
