package main

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rbo13/go-api-assessment/src/db"

	_ "github.com/go-sql-driver/mysql"
)

type registerPayload struct {
	Teacher  string   `json:"teacher"`
	Students []string `json:"students"`
}

type commonStudentsReponse struct {
	Students []string `json:"students"`
}

func main() {
	ctx := context.Background()

	conn, err := db.CreateNewConnection(&db.Config{
		Ctx:      ctx,
		MaxConns: 16,
		DSN:      "root:password@tcp(localhost:3306)/api_db?parseTime=true&loc=Local",
	})
	if err != nil {
		log.Fatalf("Cannot start API due to: %v \n", err)
	}
	defer conn.Close()

	e := echo.New()

	e.GET("/api/commonstudents", func(c echo.Context) error {
		json := map[string]interface{}{}

		json["students"] = ""
		return c.JSON(http.StatusOK, json)
	})

	e.POST("/api/register", func(c echo.Context) error {
		var payload registerPayload
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		// teacherRepo := mysql.NewTeacherRepository(conn)
		// teacherSrvc := service.NewTeacher(teacherRepo)

		// newTeacher := domain.Teacher{
		// 	Email: payload.Teacher,
		// }

		// currentTeacher, err := teacherSrvc.RetrieveTeacherByEmail(c.Request().Context(), newTeacher.Email)
		// if err != nil {
		// 	return c.JSON(http.StatusBadRequest, err)
		// }

		// if err := teacherSrvc.AddTeacher(c.Request().Context(), newTeacher); err != nil {
		// 	return c.JSON(http.StatusBadRequest, err)
		// }

		return c.JSON(http.StatusNoContent, payload)
	})

	e.Logger.Fatal(e.Start(":3000"))
}
