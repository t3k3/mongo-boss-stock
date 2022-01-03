package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/t3k3/mongo-boss-stock/handlers"
)

func OrderRoute(route fiber.Router) {
	route.Get("/", handlers.GetAllOrders)
	route.Get("/:id", handlers.GetOrder)
	route.Post("/", handlers.AddOrder)
	route.Put("/:id", handlers.UpdateOrder)
	route.Delete("/:id", handlers.DeleteOrder)
}
