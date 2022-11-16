package middleware

import (
	customLogger "fiber-template/pkg/utils/logger"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// type Tx string

// FiberMiddleware provide Fiber's built-in middlewares.
// See: https://docs.gofiber.io/api/middleware
func FiberMiddleware(a *fiber.App) { //, file *os.File) {
	a.Use(
		// Add CORS to each route.
		cors.New(cors.Config{
			// AllowOrigins:     "*", // "https://gofiber.io, https://gofiber.net",
			AllowMethods:     "GET,POST",
			AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
			AllowCredentials: true,
		}),
		// Add simple logger.
		logger.New(logger.Config{
			Format: "${time} | ${ip}:${port} | ${status} | ${method}\t| ${path} | ${latency}\n",
			// Next:         nil,
			TimeFormat: "02 Jan 06 15:04 MST", // "15:04:05",
			TimeZone:   "Local",
			// TimeInterval: 500 * time.Millisecond,
			// Output:       os.Stdout,
		}),
		logger.New(logger.Config{
			Format:     "${time} | ${ip}:${port} | ${status} | ${method}\t| ${path} | ${latency}\n",
			TimeFormat: "02 Jan 06 15:04 MST", // "2006-01-02T15:04:05",
			TimeZone:   "Local",
			Output:     customLogger.LogFiberFile, // as file,
		}),
	)
}
