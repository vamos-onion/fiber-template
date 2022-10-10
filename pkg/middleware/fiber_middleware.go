package middleware

import (
	"context"
	customLogger "fiber-template/pkg/utils/logger"
	"fiber-template/platform/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Tx string

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
	// DB transaction settings
	a.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) { // Returns true if the client requested upgrade to the WebSocket protocol
			db := database.DB.MariaDB
			c.Locals("RdbConnection", db)
			return c.Next()
		}
		var isTx bool
		tx := database.DB.MariaDB
		switch c.Method() {
		case "POST", "PUT", "DELETE":
			ctx := context.WithValue(c.UserContext(), Tx("RDBRollBack"), false)
			c.SetUserContext(ctx)
			isTx = true
			tx = tx.Begin()
			ctx = context.WithValue(c.UserContext(), Tx("RdbConnection"), tx)
			c.SetUserContext(ctx)
		default:
			isTx = false
			ctx := context.WithValue(c.UserContext(), Tx("RdbConnection"), tx)
			c.SetUserContext(ctx)
		}
		defer func(isTx bool) {
			if isTx {
				if c.Response().StatusCode() >= 500 || c.UserContext().Value(Tx("RDBRollBack")).(bool) {
					tx.Rollback()
				} else {
					tx.Commit()
				}
			}
		}(isTx)
		return c.Next()
	})
}
