package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/t3k3/mongo-boss-stock/config"
	"github.com/t3k3/mongo-boss-stock/routes"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success":     true,
			"message":     "You are at the root endpoint 😉",
			"github_repo": "https://github.com/t3k3/mongo-boss-stock",
		})
	})

	api := app.Group("/api")

	routes.ProductRoute(api.Group("/products"))
	routes.CategoryRoute(api.Group("/categories"))
	routes.RepairRoute(api.Group("/repairs"))
	routes.StoreRoute(api.Group("/stores"))
	routes.CustomerRoute(api.Group("/customers"))
	routes.OrderRoute(api.Group("/orders"))
	routes.UploadRoute(api.Group("/upload"))
}

func main() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	app := fiber.New()

	app.Static("/static/thumbs", "./uploads/thumbs", fiber.Static{
		Compress: true,
	})

	app.Use(cors.New())
	app.Use(logger.New())

	//Eğer /dashboard endpointe istek gelirse web sayfasından ram, istekler, kullanımlar vs monitör edilebilir.
	//app.Get("/dashboard", monitor.New())

	config.ConnectDB()

	setupRoutes(app)

	port := os.Getenv("PORT")
	err := app.Listen(":" + port)

	if err != nil {
		log.Fatal("Error app failed to start")
		panic(err)
	}
}
