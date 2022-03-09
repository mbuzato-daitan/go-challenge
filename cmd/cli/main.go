package main

import (
	"database/sql"
	"log"
	"os"
	"reflect"

	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
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

	repo := repository.NewTaskRepository(db)

	testVanilla(repo)
}

func testVanilla(repo *repository.TaskRepository) {
	task, err := repo.Create("Test 1")
	if err != nil {
		log.Fatal(err)
	}

	tasks, err := repo.ListAll()
	if err != nil {
		log.Fatal(err)
	}

	if len(tasks) != 1 || !reflect.DeepEqual(task, tasks[0]) {
		log.Fatal("Expected task to have been created.")
	}

	gottenTask, err := repo.GetByID(task.ID)
	if err != nil {
		log.Fatal(err)
	}

	if !reflect.DeepEqual(task, gottenTask) {
		log.Fatal("Expected task to have been returned.")
	}

	tasks, err = repo.ListByCompletion(true)
	if err != nil {
		log.Fatal(err)
	}

	if len(tasks) != 0 {
		log.Fatal("Expected task to not have been found.")
	}

	task.Completed = true
	err = repo.Update(task)
	if err != nil {
		log.Fatal(err)
	}

	tasks, err = repo.ListByCompletion(true)
	if err != nil {
		log.Fatal(err)
	}

	if len(tasks) != 1 || !reflect.DeepEqual(task, tasks[0]) {
		log.Fatal("Expected task to have been found.")
	}
}
