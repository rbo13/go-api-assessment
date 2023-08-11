package api

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rbo13/go-api-assessment/src/service"
)

func (a *api) Handlers() *echo.Echo {
	engine := echo.New()

	engine.Use(
		middleware.Recover(),
		middleware.Gzip(),
		middleware.RequestID(),
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Output: os.Stdout,
		}),
	)

	engine.GET("/api/commonstudents", func(c echo.Context) error {
		json := map[string]interface{}{}

		teacherSrvc := service.NewTeacher(a.teacherRepo)

		params := c.QueryParams()["teacher"]
		if params == nil {
			json["message"] = "Please add a query parameters"
			return c.JSON(http.StatusBadRequest, json)
		}

		res, err := teacherSrvc.RetrieveCommonStudents(c.Request().Context(), params)
		if err != nil {
			return c.JSON(http.StatusNotFound, err)
		}

		json["students"] = res

		return c.JSON(http.StatusOK, json)
	})

	return engine
}
