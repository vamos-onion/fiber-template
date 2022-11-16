package middleware

import (
	"context"
	"fiber-template/platform/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Tx string

// DB transaction settings
func AutoTx() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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
	}
}
