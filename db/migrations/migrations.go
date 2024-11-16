package migrations

import (
	"github.com/mohdjishin/order-inventory-management/db"
	"github.com/mohdjishin/order-inventory-management/internal/models"
	log "github.com/mohdjishin/order-inventory-management/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Run() error {
	log.Info("Running database migrations")

	err := db.GetDb().AutoMigrate(&models.Product{}, &models.Order{}, &models.Inventory{}, &models.PricingHistory{}, &models.User{})
	if err != nil {
		log.Fatal("failed to migrate database", zap.Error(err))
		return err
	}
	if err := createTrigger(db.GetDb()); err != nil {

	}
	log.Info("Database migration successful")
	return nil
}

func createTrigger(db *gorm.DB) error {
	triggerFunction := `
		CREATE OR REPLACE FUNCTION adjust_product_price_on_stock_change()
		RETURNS TRIGGER AS $$
		BEGIN
			-- Check if the stock is being reduced (if new stock is less than old stock)
			IF (NEW.stock < OLD.stock) THEN
				-- Calculate the percentage of remaining stock
				DECLARE
					remaining_stock_percentage FLOAT;
				BEGIN
					-- Calculate the remaining stock percentage
					remaining_stock_percentage := (NEW.stock::FLOAT / OLD.stock) * 100;

					-- Update the product price based on the remaining stock percentage
					UPDATE products
					SET price = price * (1 + (100 - remaining_stock_percentage) / 100)
					WHERE id = NEW.product_id;
				END;
			END IF;
			-- Return the new inventory row
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`

	createTrigger := `
		CREATE TRIGGER adjust_price_on_inventory_update
		AFTER UPDATE ON inventories
		FOR EACH ROW
		WHEN (NEW.stock < OLD.stock)
		EXECUTE FUNCTION adjust_product_price_on_stock_change();
	`

	if err := db.Exec(triggerFunction).Error; err != nil {
		log.Error("Error creating trigger function:", zap.Error(err))
		return err
	}

	if err := db.Exec(createTrigger).Error; err != nil {
		log.Error("Error creating trigger:", zap.Error(err))
		return err
	}

	return nil
}
