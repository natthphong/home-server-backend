package middleware

import (
	"strings"

	"github.com/natthphong/home-server-backend/api"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// JWTMiddlewareWithObjects validates JWT and checks for required objects
func JWTMiddlewareWithObjects(jwtSecret string, requireObjects []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(requireObjects) == 0 {
			return c.Next()
		}
		tokenString := c.Get("Authorization")
		if len(tokenString) == 0 {
			return api.JwtError(c, "Token Not Found")
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			return api.JwtError(c, "Token Expired")
		}
		if !token.Valid {
			return api.JwtError(c, "Invalid or expired token")
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return api.JwtError(c, "Failed to parse token claims")
		}
		roles, ok := claims["roles"].([]interface{})
		if !ok {
			return api.Forbidden(c)
		}

		for _, role := range roles {
			roleMap, isMap := role.(map[string]interface{})
			if !isMap {
				continue
			}

			objects, hasObjects := roleMap["objects"].([]interface{})
			if !hasObjects {
				continue
			}

			for _, object := range objects {
				objectStr, isString := object.(string)
				if !isString {
					continue
				}

				for _, requiredObject := range requireObjects {
					if objectStr == requiredObject {
						return c.Next()
					}
				}
			}
		}
		return api.Forbidden(c)
	}
}
