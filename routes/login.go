package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/t3k3/mongo-boss-stock/handlers"
)

func LoginRoute(route fiber.Router) {
	route.Get("/", handlers.GetAllCategories)
}
