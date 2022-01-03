package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/t3k3/mongo-boss-stock/handlers"
)

func CustomerRoute(route fiber.Router) {
	route.Get("/", handlers.GetAllCustomers)
	route.Get("/:id", handlers.GetCustomer)
	route.Post("/", handlers.AddCustomer)
	route.Put("/:id", handlers.UpdateCustomer)
	route.Delete("/:id", handlers.DeleteCustomer)
}
