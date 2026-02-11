package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type AuthMiddlewareConfig struct {
	keyfunc  keyfunc.Keyfunc
	audience string
	issuer   string
}

func NewAuthMiddlewareConfig(keyfunc keyfunc.Keyfunc, audience, issuer string) AuthMiddlewareConfig {
	return AuthMiddlewareConfig{keyfunc: keyfunc, audience: audience, issuer: issuer}
}

func isValidToken(config AuthMiddlewareConfig, tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, config.keyfunc.Keyfunc, jwt.WithAudience(config.audience), jwt.WithIssuer(config.issuer))
	if err != nil {
		return jwt.MapClaims{}, err
	}

	if !token.Valid {
		return jwt.MapClaims{}, fmt.Errorf("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return jwt.MapClaims{}, fmt.Errorf("token claims are not of type jwt.MapClaims")
	}

	return claims, nil
}

func AuthMiddleware(config AuthMiddlewareConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header required")
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == authHeader {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization format")
			}

			_, err := isValidToken(config, token)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			return next(c)
		}
	}
}
