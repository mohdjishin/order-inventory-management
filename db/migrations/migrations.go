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
					old_price FLOAT;
					new_price FLOAT;
				BEGIN
					-- Retrieve the old price
					SELECT price INTO old_price FROM products WHERE id = NEW.product_id;

					-- Calculate the remaining stock percentage
					remaining_stock_percentage := (NEW.stock::FLOAT / OLD.stock) * 100;

					-- Calculate the new price
					new_price := old_price * (1 + (100 - remaining_stock_percentage) / 100);

					-- Update the product price
					UPDATE products
					SET price = new_price
					WHERE id = NEW.product_id;

					-- Insert into pricing_histories table
					INSERT INTO pricing_histories (product_id, old_price, new_price, created_at, updated_at)
					VALUES (NEW.product_id, old_price, new_price, NOW(), NOW());
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

	// Execute the function creation script
	if err := db.Exec(triggerFunction).Error; err != nil {
		log.Error("Error creating trigger function:", zap.Error(err))
		return err
	}

	// Execute the trigger creation script
	if err := db.Exec(createTrigger).Error; err != nil {
		log.Error("Error creating trigger:", zap.Error(err))
		return err
	}

	return nil
}
