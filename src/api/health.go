package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *api) healthCheckHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		json := map[string]interface{}{
			"status":  "OK!",
			"message": "API is up and running!",
		}
		return c.JSON(http.StatusOK, json)
	}
}
