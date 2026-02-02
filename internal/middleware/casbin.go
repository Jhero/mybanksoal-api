package middleware

import (
	"net/http"

	"github.com/casbin/casbin/v3"
	"github.com/labstack/echo/v4"
)

func CasbinMiddleware(enforcer *casbin.Enforcer) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole := c.Get("role").(string)
			path := c.Path()
			method := c.Request().Method

			// Check permission
			allowed, err := enforcer.Enforce(userRole, path, method)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error checking permissions"})
			}

			if !allowed {
				return c.JSON(http.StatusForbidden, map[string]string{"error": "Access forbidden"})
			}

			return next(c)
		}
	}
}
