package routes

import (
	"fiber-template/app/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func WebsocketRoutes(a *fiber.App) {
	controllers.InitController()
	websocket_route := a.Group("/admin/v1")
	// websocket_route.Get("/wsconn/", websocket.New(controllers.WsConn)) // app/controllers/account_controller.go

	// @ ACL
	var allowed []string
	// allowed = append(allowed, "origin.example.com")
	allowed = append(allowed, "*")
	websocket_route.Get("/wsconn/", websocket.New(controllers.WsConn, websocket.Config{
		Origins: allowed,
	})) // app/controllers/account_controller.go

	test_websocket := a.Group("/wstest")
	test_websocket.Get("/tx", websocket.New(controllers.WsTxTest, websocket.Config{
		Origins: allowed,
	}))
	test_websocket.Get("/supertx", websocket.New(controllers.WsSuperTxTest, websocket.Config{
		Origins: allowed,
	}))
}
