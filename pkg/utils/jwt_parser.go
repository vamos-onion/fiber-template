package utils

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// TokenMetadata struct to describe metadata in JWT.
type TokenMetadata struct {
	UserID      uuid.UUID
	Credentials map[string]bool
	Expires     int64
}

// ExtractTokenMetadata func to extract metadata from JWT.
func ExtractTokenMetadata(c *fiber.Ctx) (*TokenMetadata, error) {
	token, err := verifyToken(c)
	if err != nil {
		return nil, err
	}

	// Setting and checking token and credentials.
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// User ID.
		userID, err := uuid.Parse(claims["id"].(string))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		// Expires time.
		expires := int64(claims["expires"].(float64))

		// User credentials.
		credentials := map[string]bool{
			"next_to_us_sso:is_valid": claims["next_to_us_sso:is_valid"].(bool),
			"next_to_us_sso:master":   claims["next_to_us_sso:master"].(bool),
			"next_to_us_sso:user":     claims["next_to_us_sso:user"].(bool),
		}

		return &TokenMetadata{
			UserID:      userID,
			Credentials: credentials,
			Expires:     expires,
		}, nil
	}

	return nil, err
}

func extractToken(c *fiber.Ctx) string {
	// Another method
	// JWT in bearer token
	//
	// bearToken := c.Get("Authorization")
	// // Normally Authorization HTTP header.
	// onlyToken := strings.Split(bearToken, " ")
	// if len(onlyToken) == 2 {
	// 	return onlyToken[1]
	// }
	// return ""

	// JWT in cookie
	cookieToken := c.Cookies("access")
	return string(cookieToken)
}

func verifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := extractToken(c)

	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
}
