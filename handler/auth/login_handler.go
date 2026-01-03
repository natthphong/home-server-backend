package auth

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natthphong/home-server-backend/api"
)

func GenerateJWTForUser(
	db *pgxpool.Pool,
	userID, password, appCode string,
	jwtSecret string,
	accessTokenDuration, refreshTokenDuration time.Duration,
	refreshTokenFlag bool,
) (map[string]interface{}, error) {

	// Access Token
	accessTokenClaims := jwt.MapClaims{
		"username": "test_001",
		"exp":      time.Now().Add(accessTokenDuration).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, errors.New("Failed to generate access token")
	}

	// Refresh Token
	refreshTokenClaims := jwt.MapClaims{
		"username": "test_001",
		"exp":      time.Now().Add(refreshTokenDuration).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, errors.New("Failed to generate refresh token")
	}

	// Response
	response := map[string]interface{}{
		"accessToken":  accessTokenString,
		"refreshToken": refreshTokenString,
		"jwtBody":      accessTokenClaims,
	}
	return response, nil
}

func LoginHandler(db *pgxpool.Pool, jwtSecret string, accessTokenDuration, refreshTokenDuration time.Duration) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return api.BadRequest(c, "Invalid input")
		}

		// Call GenerateJWTForUser
		response, err := GenerateJWTForUser(db, req.Username, req.Password, req.AppCode, jwtSecret, accessTokenDuration, refreshTokenDuration, false)
		if err != nil {
			return api.InternalError(c, err.Error())
		}

		return api.Ok(c, response)
	}
}
