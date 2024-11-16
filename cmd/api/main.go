package main

import (
	"github.com/mohdjishin/order-inventory-management/config"
	_ "github.com/mohdjishin/order-inventory-management/config"
	"github.com/mohdjishin/order-inventory-management/db/migrations"
	"github.com/mohdjishin/order-inventory-management/internal/router"
	"go.uber.org/zap"

	log "github.com/mohdjishin/order-inventory-management/logger"
)

func main() {

	if err := run(); err != nil {
		log.Fatal("Failed to run", zap.Error(err))

	}

}

func run() error {
	if err := migrations.Run(); err != nil {
		log.Fatal("Failed to run migrations", zap.Error(err))
		return err
	}

	r := router.New()
	log.Info("Starting server")
	log.Info("Server started on port " + config.Get().Port)
	log.Fatal("Listen ", zap.Error(r.Listen(config.Get().Port)))
	return nil
}
