package main

import (
	"fiber-template/pkg/configs"
	"fiber-template/pkg/middleware"
	"fiber-template/pkg/routes"
	"fiber-template/pkg/utils"
	"fiber-template/pkg/utils/logger"
	"fiber-template/platform/cache"
	"fiber-template/platform/database"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/gofiber/websocket/v2"

	_ "github.com/joho/godotenv/autoload" // load .env file automatically
	// "github.com/robfig/cron/v3" // for cron job
)

func main() {
	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Initialize Loggers
	logger.InitLogger()

	// Open & Close log file
	defer logger.FileClose()

	// // Database pool
	database.InitDBConnection(app)

	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware for app.

	// // Redis pool
	cache.InitRedisConnection(app)

	// Routes.
	routes.HttpRoutes(app)      // Register a public routes for app.
	routes.WebsocketRoutes(app) // Register a public websocket routes for app.

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
