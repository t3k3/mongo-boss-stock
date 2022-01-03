package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/t3k3/mongo-boss-stock/handlers"
)

func StoreRoute(route fiber.Router) {
	route.Get("/", handlers.GetAllStores)
	route.Get("/:id", handlers.GetStore)
	route.Post("/", handlers.AddStore)
	route.Put("/:id", handlers.UpdateStore)
	route.Delete("/:id", handlers.DeleteStore)
}
