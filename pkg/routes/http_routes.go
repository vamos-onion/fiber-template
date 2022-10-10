package routes

import (
	"fiber-template/app/controllers"

	"github.com/gofiber/fiber/v2"
)

// See : https://docs.gofiber.io/guide/routing
func HttpRoutes(a *fiber.App) {
	rest_route := a.Group("/rest")
	rest_route.Get("/get/:param", controllers.Get)
	rest_route.Post("/post", controllers.Post)
	rest_route.Post("/sso", controllers.SSO)
}
