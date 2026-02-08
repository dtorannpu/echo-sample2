package handler

import (
	"echo-sample2/internal/httpclient"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authClient *httpclient.AuthHttpClient
}

func NewAuthHandler(authClient *httpclient.AuthHttpClient) *AuthHandler {
	return &AuthHandler{authClient: authClient}
}

func (a *AuthHandler) RegisterAuthRoutes(e *echo.Echo) {
	e.GET("/auth-well-known-config", func(c echo.Context) error {
		resp, err := a.authClient.GetAuthWellKnownConfig(c.Request().Context())
		if err != nil {
			c.Logger().Error("auth-well-known-configエラー", err)
			return err
		}
		defer resp.Body.Close()

		c.Response().Header().Set("Content-Type", resp.Header.Get("Content-Type"))

		return c.Stream(resp.StatusCode, resp.Header.Get("Content-Type"), resp.Body)
	})
}
