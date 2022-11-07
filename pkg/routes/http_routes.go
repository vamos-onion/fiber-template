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
}

func AccountRoutes(a *fiber.App) {
	account_route := a.Group("/account")
	account_route.Post("/sso", controllers.SSO)
	account_route.Post("/register", controllers.Register)
	account_route.Post("/login", controllers.Login)
	account_route.Post("/modify", controllers.Modify)
	account_route.Post("/withdrawal", controllers.Withdrawal)
	// account_route.Get("/logout", controllers.Logout)
}
