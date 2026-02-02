package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/jovan/mybanksoal-api/config"
	"github.com/jovan/mybanksoal-api/internal/repository"
	"github.com/jovan/mybanksoal-api/pkg/utils"
)

func AuthMiddleware(cfg *config.Config, userRepo repository.UserRepository) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 1. Check for API Key
			apiKey := c.Request().Header.Get("X-API-KEY")
			if apiKey != "" {
				user, err := userRepo.FindByAPIKey(apiKey)
				if err != nil {
					return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid API Key"})
				}
				c.Set("user_id", user.ID)
				c.Set("role", user.Role)
				return next(c)
			}

			// 2. Check for Bearer Token
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing authorization header"})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid authorization header"})
			}

			claims, err := utils.ParseToken(parts[1], cfg)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			}

			c.Set("user_id", claims.UserID)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}
