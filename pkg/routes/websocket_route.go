package routes

import (
	"fiber-template/app/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func WebsocketRoutes(a *fiber.App) {
	controllers.InitController()
	websocket_route := a.Group("/websocket")

	// @ ACL
	var allowed []string
	// allowed = append(allowed, "origin.example.com")
	allowed = append(allowed, "*")
	websocket_route.Get("/wsconn/", websocket.New(controllers.WsConn, websocket.Config{
		Origins: allowed,
	}))
}
