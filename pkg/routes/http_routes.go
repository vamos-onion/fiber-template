package routes

import (
	"fiber-template/app/controllers"

	"github.com/gofiber/fiber/v2"
)

// See : https://docs.gofiber.io/guide/routing
func HttpRoutes(a *fiber.App) {
	rest_route := a.Group("/rest")
	rest_route.Get("/get/:param", controllers.Get)
	rest_route.Post("", controllers.Post)

	publisher_route := a.Group("/pub/v1")
	publisher_route.Post("/sso", controllers.SSO)

	test_route := a.Group("/test")
	test_route.Post("/ctx", controllers.TxTest)
}
