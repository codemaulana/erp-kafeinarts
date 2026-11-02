package routes

import (
	"erp-kafeinart/handler"

	"github.com/gofiber/fiber/v3"
)

// RegisterUserRoutes mendaftarkan rute CRUD pengguna.
func RegisterUserRoutes(api fiber.Router, userHandler *handler.UserHandler, jwtAuth fiber.Handler) {
	users := api.Group("/users", jwtAuth)
	users.Get("/", userHandler.GetUsers)
	users.Post("/", userHandler.CreateUser)
	users.Get("/:id", userHandler.GetUserByID)
	users.Put("/:id", userHandler.UpdateUser)
	users.Delete("/:id", userHandler.DeleteUser)
}
