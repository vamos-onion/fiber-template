package middleware

import (
	"fiber-template/app/models"
	"os"

	"github.com/gofiber/fiber/v2"

	jwtMiddleware "github.com/gofiber/jwt/v2"
)

// JWTProtected func for specify routes group with JWT authentication.
// See: https://github.com/gofiber/jwt
func JwtInCookie() func(*fiber.Ctx) error {
	// Create config for JWT authentication middleware.
	config := jwtMiddleware.Config{
		SigningKey:   []byte(os.Getenv("JWT_SECRET_KEY")),
		ContextKey:   "jwt", // used in private routes
		TokenLookup:  "cookie:access,refresh",
		ErrorHandler: jwtError,
	}

	return jwtMiddleware.New(config)
}

func JwtInBody() func(*fiber.Ctx) error {
	// Create config for JWT authentication middleware.
	config := jwtMiddleware.Config{
		SigningKey:   []byte(os.Getenv("JWT_SECRET_KEY")),
		ContextKey:   "jwt", // used in private routes
		ErrorHandler: jwtError,
	}

	return jwtMiddleware.New(config)
}

func jwtError(c *fiber.Ctx, err error) error {
	// Return status 401 and failed authentication error.
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:   fiber.StatusBadRequest,
			Response: err.Error(),
		})
	}

	// Return status 401 and failed authentication error.
	return c.Status(fiber.StatusUnauthorized).JSON(&models.R{
		Status:   fiber.StatusUnauthorized,
		Response: err.Error(),
	})
}
