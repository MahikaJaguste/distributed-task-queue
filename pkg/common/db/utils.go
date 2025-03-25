package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

var (
	DBCon *sql.DB
)

var ENV_FILE_PATH = "../../.env"

type TaskStatus int
type Task struct {
	Id          int
	Name        string
	PickedAt    time.Time
	ProcessedAt time.Time
	CompletedAt time.Time
	WorkerId    int
	Status      TaskStatus
}

const (
	Pending TaskStatus = iota
	Processing
	Completed
)

func SetupDb() {

	cfg := mysql.Config{
		User:      os.Getenv("DB_USER"),
		Passwd:    os.Getenv("DB_PASSWORD"),
		Net:       "tcp",
		Addr:      os.Getenv("DB_HOST"),
		DBName:    os.Getenv("DB_NAME"),
		ParseTime: true,
	}

	var err error
	DBCon, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}
	DBCon.SetConnMaxLifetime(time.Minute * 3)
	DBCon.SetMaxOpenConns(10)
	DBCon.SetMaxIdleConns(10)

	pingErr := DBCon.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("DB Connected!")
}
