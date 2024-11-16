package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mohdjishin/order-inventory-management/internal/meta"
)

func GetVersion(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"version":   meta.GetVersion(),
		"buildTime": meta.GetBuildTime(),
		"hash/tag":  meta.GetCommitHash(),
	})
}

func HealthCheck(c fiber.Ctx) error {

	return c.JSON(fiber.Map{
		"server":   "Ok",
		"database": meta.GetDatabaseStats(),
	})

}
