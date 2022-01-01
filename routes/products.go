package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/t3k3/mongo-boss-stock/handlers"
)

func ProductRoute(route fiber.Router) {
	route.Get("/", handlers.GetAllProducts)
	route.Get("/:id", handlers.GetProduct)
	route.Post("/", handlers.AddProduct)
	route.Put("/:id", handlers.UpdateProduct)
	route.Delete("/:id", handlers.DeleteProduct)
}
