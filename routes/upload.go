package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/t3k3/mongo-boss-stock/handlers"
)

func UploadRoute(route fiber.Router) {
	route.Post("/", handlers.UploadFile)

}
