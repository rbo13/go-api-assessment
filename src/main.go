package main

import (
	"context"
	"database/sql"
	"log"

	database "github.com/rbo13/go-api-assessment/generated/db"
	"github.com/rbo13/go-api-assessment/src/db"

	_ "github.com/go-sql-driver/mysql"
)

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

	query := database.New(conn)

	res, err := query.CreateTeacher(ctx, database.CreateTeacherParams{
		TeacherName: sql.NullString{
			String: "Teacher Richard",
			Valid:  true,
		},
		Email: sql.NullString{
			String: "teacherrichard@gmail.com",
			Valid:  true,
		},
	})
	if err != nil {
		log.Fatalf("Error inserting teacher: %v \n", err)
	}

	if res != nil {
		log.Println("Successfully inserted! ")
	}
}
