package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/t3k3/mongo-boss-stock/handlers"
)

func CategoryRoute(route fiber.Router) {
	route.Get("/", handlers.GetAllCategories)
	route.Get("/:id", handlers.GetCategory)
	route.Post("/", handlers.AddCategory)
	route.Put("/:id", handlers.UpdateCategory)
	route.Delete("/:id", handlers.DeleteCategory)
}
