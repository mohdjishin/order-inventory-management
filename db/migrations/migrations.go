package migrations

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/mohdjishin/order-inventory-management/db"
	"github.com/mohdjishin/order-inventory-management/internal/models"
	log "github.com/mohdjishin/order-inventory-management/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Run() error {
	log.Info("Running database migrations")

	err := db.GetDb().AutoMigrate(&models.Product{}, &models.Order{}, &models.Inventory{}, &models.PricingHistory{}, &models.User{})
	if err != nil {
		log.Fatal("failed to migrate database", zap.Error(err))
		return err
	}
	if err := initialSetup(); err != nil {
		log.Fatal("failed to run initial setup", zap.Error(err))
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

	if err := db.Exec(triggerFunction).Error; err != nil {
		log.Error("Error creating trigger function:", zap.Error(err))
		return err
	}

	if err := db.Exec(createTrigger).Error; err != nil {
		if strings.Contains(err.Error(), "already exists") {
			log.Warn("Trigger already exists, skipping creation.")
		} else {
			log.Error("Error creating trigger:", zap.Error(err))
			return err
		}
	}

	return nil
}

func initialSetup() error {
	log.Info("Creating initial admin user")

	var count int64
	err := db.GetDb().
		Model(&models.User{}).
		Where("role = ?", models.SuperAdmin).
		Count(&count).Error
	if err != nil {
		log.Error("Failed to check if super admin user exists", zap.Error(err))
		return err
	}

	if count > 0 {
		log.Info("Super admin user already exists, skipping creation")
		return nil
	}

	decordedP, err := base64.StdEncoding.DecodeString("cGFzc3dPcmRAMTIz")
	if err != nil {
		log.Error("Failed to decode the password", zap.Error(err))
		return errors.New("unable to decode password")
	}

	hp, err := bcrypt.GenerateFromPassword(decordedP, bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash the password", zap.Error(err))
		return errors.New("unable to generate password hash")
	}

	adminUser := models.User{
		Email:     "admin@oim.com",
		Password:  string(hp),
		Role:      models.SuperAdmin,
		FirstName: "Super",
		LastName:  "Admin",
		Phone:     "9605384376",
	}

	if err := db.GetDb().Create(&adminUser).Error; err != nil {
		log.Error("Failed to create initial admin user", zap.Error(err))
		return err
	}

	log.Info("Initial super admin user created successfully")

	log.Info("Creating trigger for adjusting product price on stock change")
	//  creating trigger only once per db
	if err := createTrigger(db.GetDb()); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			log.Fatal("failed to create trigger", zap.Error(err))
			return err

		} else {
			log.Warn("Trigger already exists, skipping creation.")
		}
	}
	return nil
}
