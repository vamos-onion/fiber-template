package controllers

import (
	"context"
	"encoding/base64"
	"fiber-template/app/interfaces"
	"fiber-template/app/models"
	"fiber-template/app/queries"
	"fiber-template/pkg/middleware"
	"fiber-template/pkg/utils"
	log "fiber-template/pkg/utils/logger"
	"fiber-template/platform/cache"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type JwtController struct {
	ctx        *fiber.Ctx
	apiRequest *models.API
	log        *log.User
	account    *models.SsoUser
	rdbQuery   interfaces.RdbQuery
	redisQuery interfaces.RedisQuery
}

func newJwtClient(c *fiber.Ctx) *JwtController {
	l := log.NewUser()
	s := &queries.SqlQuery{
		Log: l,
		DB:  c.UserContext().Value(middleware.Tx("RdbConnection")).(*gorm.DB),
	}
	return &JwtController{
		ctx:        c,
		apiRequest: &models.API{},
		log:        l,
		account:    &models.SsoUser{},
		rdbQuery:   s,
		redisQuery: &queries.RedisQuery{Log: l},
	}
}

// SSO is verifying user and set token as key, value in Redis.
//	- Single Sign On
//	- Some customers would not use OAuth, and this method is for them.
//	- SSO := timestamp + base64(organization + " " + username + " " + sha256(saltkey + username + timestamp))
//
func SSO(c *fiber.Ctx) error {
	// Create a new client
	j := newJwtClient(c)

	// Create a new user struct.
	requestBody := &models.Login{
		// SSO string `json:"sso"`
	}

	// Checking received data from JSON body.
	if err := c.BodyParser(j.apiRequest); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:      fiber.StatusBadRequest,
			Transaction: j.apiRequest.Transaction,
			Response:    err.Error(),
		})
	}

	// Parsing received data from JSON body
	if err := zeroValid(j.apiRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:      fiber.StatusBadRequest,
			Transaction: j.apiRequest.Transaction,
			Response:    err.Error(),
		})
	}

	// Parsing received data from JSON body's body field.
	jsonParse(j.apiRequest.Body, &requestBody)
	if err := zeroValid(requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:      fiber.StatusBadRequest,
			Transaction: j.apiRequest.Transaction,
			Response:    err.Error(),
		})
	}

	// Verifying the length of the bytes of the requested payload
	if len(requestBody.SSO) < 10 {
		return c.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:      fiber.StatusBadRequest,
			Transaction: j.apiRequest.Transaction,
			Response:    "failed; your request contains invalid data",
		})
	}

	// Verifying timestamp field
	ssoTimestamp := requestBody.SSO[:10]
	minAccessTime := time.Now().Add(time.Minute * 10).Unix()
	clientAccessTime, _ := strconv.ParseInt(ssoTimestamp, 10, 64)
	if clientAccessTime > minAccessTime {
		strconv.FormatInt(minAccessTime, 10)
		return c.Status(fiber.StatusUnauthorized).JSON(&models.R{
			Status:      fiber.StatusUnauthorized,
			Transaction: j.apiRequest.Transaction,
			Response:    "unauthorized, your request contains invalid data",
		})
	}

	// Verifying base64 encoded field
	ssoBase64 := strings.Trim(requestBody.SSO[10:], "\n")
	b64d, _ := base64.StdEncoding.DecodeString(ssoBase64)
	separated := strings.Split(string(b64d), " ")
	if len(separated) < 3 {
		// failed to verify sso secret key
		return c.Status(fiber.StatusUnauthorized).JSON(&models.R{
			Status:      fiber.StatusUnauthorized,
			Transaction: j.apiRequest.Transaction,
			Response:    "unauthorized, your request format is wrong",
		})
	}

	// SSO verifying logic
	ssoSecret := os.Getenv("SSO_SECRET") + separated[1] + ssoTimestamp
	utils.HashPassword(&ssoSecret)
	if separated[2] != ssoSecret {
		// failed to verify sso secret key
		return c.Status(fiber.StatusUnauthorized).JSON(&models.R{
			Status:      fiber.StatusUnauthorized,
			Transaction: j.apiRequest.Transaction,
			Response:    "unauthorized, your request contains unauthorized data",
		})
	}
	j.log.Infoln("verified user")

	// Set credentials as token claim
	credentials := []string{}
	credentials = append(credentials, "next_to_us_sso:is_valid")

	/***
	* @ Parse from DB
	*	- set credential
	**/
	credentials = append(credentials, "next_to_us_sso:user")
	v, _ := utils.MakeUintUniqueID()

	// Find the user info & upsert info
	m := &models.SsoUser{
		Organization: separated[0],
		Username:     separated[1],
		UserUuid:     utils.MakeUID(),
		UserIndex:    v,
		Status:       true,
	}
	var err error
	err = j.rdbQuery.SSO(m)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&models.R{
			Status:      fiber.StatusInternalServerError,
			Transaction: j.apiRequest.Transaction,
			Response:    err.Error(),
		})
	}

	// Generate a new pair of access and refresh tokens.
	tokens, err := utils.GenerateNewTokens(j.account.UserUuid.String(), credentials)
	if err != nil {
		// Return status 500 and token generation error.
		return c.Status(fiber.StatusInternalServerError).JSON(&models.R{
			Status:      fiber.StatusInternalServerError,
			Transaction: j.apiRequest.Transaction,
			Response:    err.Error(),
		})
	}

	// Define user ID.
	userID := j.account.UserUuid

	// Create a new Redis connection.
	connRedis := cache.Redis.UserConn

	// Set expires hours count for refresh key from .env file.
	hoursCount, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT"))

	// Save refresh token to Redis.
	errSaveToRedis := connRedis.Set(context.Background(), userID.String(), tokens.Refresh, time.Hour*time.Duration(hoursCount)).Err()
	if errSaveToRedis != nil {
		// Return status 500 and Redis connection error.
		return c.Status(fiber.StatusInternalServerError).JSON(&models.R{
			Status:      fiber.StatusInternalServerError,
			Transaction: j.apiRequest.Transaction,
			Response:    errSaveToRedis.Error(),
		})
	}

	accessExpiry, _ := strconv.Atoi(os.Getenv("HTTP_COOKIE_ACCESS_EXPIRY"))
	refreshExpiry, _ := strconv.Atoi(os.Getenv("HTTP_COOKIE_REFRESH_EXPIRY"))

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access",
		Value:    tokens.Access,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * time.Duration(accessExpiry)),
		Domain:   "your.domain.com",
		Secure:   true,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteNoneMode,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh",
		Value:    tokens.Refresh,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * time.Duration(refreshExpiry)),
		Domain:   "your.domain.com",
		Secure:   true,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteNoneMode,
	})

	// Return status 200 OK.
	return c.JSON(&models.R{
		Status:      fiber.StatusOK,
		Transaction: j.apiRequest.Transaction,
		Response:    "success",
	})
}

// UserSignOut method to de-authorize user and delete refresh token from Redis.
//
func UserSignOut(c *fiber.Ctx) error {
	// Create a new client
	j := newJwtClient(c)

	// Checking received data from JSON body.
	if err := c.BodyParser(j.apiRequest); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:      fiber.StatusBadRequest,
			Transaction: j.apiRequest.Transaction,
			Response:    err.Error(),
		})
	}

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(&models.R{
			Status:      fiber.StatusInternalServerError,
			Transaction: j.apiRequest.Transaction,
			Response:    err.Error(),
		})
	}

	// Define user ID.
	userID := claims.UserID.String()

	// Create a new Redis connection.
	connRedis := cache.Redis.UserConn

	// Save refresh token to Redis.
	errDelFromRedis := connRedis.Del(context.Background(), userID).Err()
	if errDelFromRedis != nil {
		// Return status 500 and Redis deletion error.
		return c.Status(fiber.StatusInternalServerError).JSON(&models.R{
			Status:      fiber.StatusInternalServerError,
			Transaction: j.apiRequest.Transaction,
			Response:    errDelFromRedis.Error(),
		})
	}

	// Return status 204 no content.
	return c.SendStatus(fiber.StatusNoContent)
}

// before access token expired
//	- both access & refresh token must be contained
//
func RenewTokens(c *fiber.Ctx) error {
	// Create a new client
	j := newJwtClient(c)

	// Checking received data from JSON body.
	if err := c.BodyParser(j.apiRequest); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:      fiber.StatusBadRequest,
			Transaction: j.apiRequest.Transaction,
			Response:    err.Error(),
		})
	}

	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(&models.R{
			Status:      fiber.StatusInternalServerError,
			Transaction: j.apiRequest.Transaction,
			Response:    err.Error(),
		})
	}

	// Set expiration time from JWT data of current user.
	expiresAccessToken := claims.Expires

	// Checking, if now time greather than Access token expiration time.
	if now > expiresAccessToken {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(&models.R{
			Status:      fiber.StatusUnauthorized,
			Transaction: j.apiRequest.Transaction,
			Response:    "unauthorized, check expiration time of your token",
		})
	}

	// Create a new renew refresh token struct.
	renew := &models.Renew{}
	renew.RefreshToken = c.Cookies("refresh")
	if strings.Count(renew.RefreshToken, "") == 1 {
		return c.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:      fiber.StatusBadRequest,
			Transaction: j.apiRequest.Transaction,
			Response:    "unauthorized, please check your refresh token",
		})
	}

	// Set expiration time from Refresh token of current user.
	expiresRefreshToken, err := utils.ParseRefreshToken(renew.RefreshToken)
	if err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:      fiber.StatusBadRequest,
			Transaction: j.apiRequest.Transaction,
			Response:    err.Error(),
		})
	}

	// Checking, if now time greather than Refresh token expiration time.
	if now > expiresRefreshToken {
		return c.Status(fiber.StatusUnauthorized).JSON(&models.R{
			Status:      fiber.StatusUnauthorized,
			Transaction: j.apiRequest.Transaction,
			Response:    "unauthorized, your session was ended earlier",
		})
	}

	// Define user ID.
	userID := claims.UserID

	// Set credentials list
	var credentials []string
	for i, v := range claims.Credentials {
		if v {
			credentials = append(credentials, i)
		}
	}

	// Generate JWT Access & Refresh tokens.
	tokens, err := utils.GenerateNewTokens(userID.String(), credentials)
	if err != nil {
		// Return status 500 and token generation error.
		return c.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:      fiber.StatusBadRequest,
			Transaction: j.apiRequest.Transaction,
			Response:    err.Error(),
		})
	}

	// Create a new Redis connection.
	connRedis := cache.Redis.UserConn

	// Set expires hours count for refresh key from .env file.
	hoursCount, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT"))

	// Save refresh token to Redis.
	errRedis := connRedis.Set(context.Background(), userID.String(), tokens.Refresh, time.Hour*time.Duration(hoursCount)).Err()
	if errRedis != nil {
		// Return status 500 and Redis connection error.
		return c.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:      fiber.StatusBadRequest,
			Transaction: j.apiRequest.Transaction,
			Response:    errRedis.Error(),
		})
	}

	accessExpiry, _ := strconv.Atoi(os.Getenv("HTTP_COOKIE_ACCESS_EXPIRY"))
	refreshExpiry, _ := strconv.Atoi(os.Getenv("HTTP_COOKIE_REFRESH_EXPIRY"))

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access",
		Value:    tokens.Access,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * time.Duration(accessExpiry)),
		Domain:   "*",
		Secure:   true,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteNoneMode,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh",
		Value:    tokens.Refresh,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * time.Duration(refreshExpiry)),
		Domain:   "*",
		Secure:   true,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteNoneMode,
	})

	// Return status 200 OK.
	return c.JSON(&models.R{
		Status:      fiber.StatusOK,
		Transaction: j.apiRequest.Transaction,
		Response:    "success",
	})
}

// After access token expired
//	- expired token must be contained
//	- both access & refresh token must be contained
//
func RecreateTokens(c *fiber.Ctx) error {
	// Create a new client
	j := newJwtClient(c)

	// Checking received data from JSON body.
	if err := c.BodyParser(j.apiRequest); err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:      fiber.StatusBadRequest,
			Transaction: j.apiRequest.Transaction,
			Response:    err.Error(),
		})
	}

	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return c.Status(fiber.StatusInternalServerError).JSON(&models.R{
			Status:      fiber.StatusInternalServerError,
			Transaction: j.apiRequest.Transaction,
			Response:    err.Error(),
		})
	}

	// Set expiration time from JWT data of current user.
	expiresAccessToken := claims.Expires

	// // Checking, if now time greather than Access token expiration time.
	if now < expiresAccessToken {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(&models.R{
			Status:      fiber.StatusUnauthorized,
			Transaction: j.apiRequest.Transaction,
			Response:    "unauthorized, check expiration time of your token",
		})
	}

	// Create a new renew refresh token struct.
	renew := &models.Renew{}
	renew.RefreshToken = c.Cookies("refresh")
	if strings.Count(renew.RefreshToken, "") == 1 {
		return c.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:      fiber.StatusBadRequest,
			Transaction: j.apiRequest.Transaction,
			Response:    "unauthorized, please check your refresh token",
		})
	}

	// Set expiration time from Refresh token of current user.
	expiresRefreshToken, err := utils.ParseRefreshToken(renew.RefreshToken)
	if err != nil {
		// Return status 400 and error message.
		return c.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:      fiber.StatusBadRequest,
			Transaction: j.apiRequest.Transaction,
			Response:    err.Error(),
		})
	}

	// Checking, if now time greather than Refresh token expiration time.
	if now > expiresRefreshToken {
		// Return status 401 and unauthorized error message.
		return c.Status(fiber.StatusUnauthorized).JSON(&models.R{
			Status:      fiber.StatusUnauthorized,
			Transaction: j.apiRequest.Transaction,
			Response:    "unauthorized, your session was ended earlier",
		})
	}

	// Define user ID.
	userID := claims.UserID

	// Generate JWT Access & Refresh tokens.
	tokens, err := utils.GenerateNewTokens(userID.String(), nil) // credentials)
	if err != nil {
		// Return status 500 and token generation error.
		return c.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:      fiber.StatusBadRequest,
			Transaction: j.apiRequest.Transaction,
			Response:    err.Error(),
		})
	}

	// Create a new Redis connection.
	connRedis := cache.Redis.UserConn

	// Set expires hours count for refresh key from .env file.
	hoursCount, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_KEY_EXPIRE_HOURS_COUNT"))

	// Save refresh token to Redis.
	errRedis := connRedis.Set(context.Background(), userID.String(), tokens.Refresh, time.Hour*time.Duration(hoursCount)).Err()
	if errRedis != nil {
		// Return status 500 and Redis connection error.
		return c.Status(fiber.StatusBadRequest).JSON(&models.R{
			Status:      fiber.StatusBadRequest,
			Transaction: j.apiRequest.Transaction,
			Response:    errRedis.Error(),
		})
	}

	accessExpiry, _ := strconv.Atoi(os.Getenv("HTTP_COOKIE_ACCESS_EXPIRY"))
	refreshExpiry, _ := strconv.Atoi(os.Getenv("HTTP_COOKIE_REFRESH_EXPIRY"))

	// Set cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access",
		Value:    tokens.Access,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * time.Duration(accessExpiry)),
		Domain:   "*",
		Secure:   true,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteNoneMode,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh",
		Value:    tokens.Refresh,
		Path:     "/",
		Expires:  time.Now().Add(time.Minute * time.Duration(refreshExpiry)),
		Domain:   "*",
		Secure:   true,
		HTTPOnly: true,
		SameSite: fiber.CookieSameSiteNoneMode,
	})

	// Return status 200 OK.
	return c.JSON(&models.R{
		Status:      fiber.StatusOK,
		Transaction: j.apiRequest.Transaction,
		Response:    "success",
	})
}

func tokenValidCheck(c *fiber.Ctx) (int, error) {
	// Get now time.
	now := time.Now().Unix()

	// Get claims from JWT.
	claims, err := utils.ExtractTokenMetadata(c)
	if err != nil {
		// Return status 500 and JWT parse error.
		return fiber.StatusInternalServerError,
			fmt.Errorf("%s", err.Error())
	}

	// Set expiration time from JWT data of current book.
	expires := claims.Expires

	// Checking, if now time greather than expiration from JWT.
	if now > expires {
		// Return status 401 and unauthorized error message.
		return fiber.StatusUnauthorized,
			fmt.Errorf("unauthorized, check expiration time of your token")
	}

	// Set credential `next_to_us_sso:user` from JWT data of current request.
	credential := claims.Credentials["next_to_us_sso:user"]

	// Only user with `next_to_us_sso:user` credential can access.
	if !credential {
		// Return status 403 and permission denied error message.
		return fiber.StatusForbidden,
			fmt.Errorf("permission denied, check credentials of your token")
	}
	c.Locals("id", claims.UserID)
	return fiber.StatusOK, nil
}
