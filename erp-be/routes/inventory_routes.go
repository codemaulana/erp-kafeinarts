package routes

import (
	"erp-kafeinart/handler"

	"github.com/gofiber/fiber/v3"
)

// RegisterInventoryRoutes mendaftarkan endpoint modul inventaris produk.
func RegisterInventoryRoutes(api fiber.Router, productHandler *handler.ProductHandler, jwtAuth fiber.Handler) {
	inventory := api.Group("/inventory", jwtAuth)
	inventory.Get("/products", productHandler.GetAll)
	inventory.Get("/products/:id", productHandler.GetByID)
	inventory.Post("/products", productHandler.Create)
	inventory.Put("/products/:id", productHandler.Update)
	inventory.Delete("/products/:id", productHandler.Delete)
	inventory.Post("/products/:id/adjust", productHandler.AdjustStock)
}
