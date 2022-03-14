package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/mtbuzato/go-challenge/internal/api"
	"github.com/mtbuzato/go-challenge/internal/apivanilla"
	"github.com/mtbuzato/go-challenge/internal/orm"
	"github.com/mtbuzato/go-challenge/internal/repository"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load .env.")
	}

	cfg := mysql.Config{
		User:   os.Getenv("MYSQL_USER"),
		Passwd: os.Getenv("MYSQL_PASSWORD"),
		Addr:   os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT"),
		DBName: os.Getenv("MYSQL_DATABASE"),
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	var repo api.TaskRepository
	if os.Getenv("DB_IMPL") == "orm" {
		repo, err = orm.NewTaskRepository(db)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		repo = repository.NewTaskRepository(db)
	}

	var server http.Handler
	if os.Getenv("SERVER_IMPL") == "vanilla" {
		server = apivanilla.NewAPIServer(repo)
	} else {
		server = api.NewAPIServer(repo)
	}

	http.ListenAndServe(":8080", server)
}
