package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client
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

	RDB = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDRESS"),
		// Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0, // use default DB
	})

	resp := RDB.Ping(context.Background())
	if resp.Err() != nil {
		log.Fatal(resp.Err())
	}

	log.Println("Redis Connected!")
}
