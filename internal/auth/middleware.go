package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/redhajuanda/gorengan/internal/httperror"
)

// IsLoggedIn is a JWT middleware
// JWTWithConfig provides a JSON Web Token (JWT) authentication middleware.
// - For valid token, it sets the user in context and calls next handler.
// - For invalid token, it sends “401 - Unauthorized” response.
// - For missing or invalid Authorization header, it sends “400 - Bad Request”.
var IsLoggedIn = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte("test"),
})

// IsAdmin checks wether user is an admin or not
func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		var role string
		if val, ok := claims["role"].(string); ok {
			role = val
		}

		if role != "admin" {
			return httperror.Unauthorized("")
		}
		return next(c)
	}
}
