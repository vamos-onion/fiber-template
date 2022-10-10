package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// Create config for http authentication middleware.

func Acl() func(*fiber.Ctx) error {
	config := cors.Config{
		AllowOrigins: "localhost",
	}
	return cors.New(config)
}

func BasicAuth() func(*fiber.Ctx) error {
	config := basicauth.Config{
		Users: map[string]string{
			"username": "password",
		},
	}
	return basicauth.New(config)
}
