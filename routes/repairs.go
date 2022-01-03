package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/t3k3/mongo-boss-stock/handlers"
)

func RepairRoute(route fiber.Router) {
	route.Get("/", handlers.GetAllRepairs)
	route.Get("/:id", handlers.GetRepair)
	route.Post("/", handlers.AddRepair)
	route.Put("/:id", handlers.UpdateRepair)
	route.Delete("/:id", handlers.DeleteRepair)
}
