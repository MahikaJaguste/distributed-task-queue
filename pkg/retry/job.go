package retry

import (
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/MahikaJaguste/distributed-task-queue/pkg/common/db"
)

var CONCURRENCY int
var WORKER_ID int

func StartRetryJobServer() {
	err := godotenv.Load(db.ENV_FILE_PATH)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db.SetupDb()

	handleRetry()
}

func handleFailedTasks() (int, error) {
	fmt.Println("Inside handleFailedTasks")

	expiryDuration := 5 // 5 minutes
	expiry := time.Now().Add(time.Duration(-expiryDuration) * time.Minute).Format(time.DateTime)
	result, err := db.DBCon.Exec("update tasks set pickedAt = null, processedAt = null, workerId = null, status = ? WHERE status = ? and processedAt < ?", db.Pending, db.Processing, expiry)
	if err != nil {
		fmt.Println("Error in updating failed tasks")
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Error in updating failed tasks")
		return 0, err
	}

	return int(rowsAffected), nil

}

func handleRetry() {
	for range time.Tick(time.Minute * 5) {

		rowsAffected, err := handleFailedTasks()
		if err != nil {
			fmt.Println("Error in handling failed tasks")
			fmt.Println(err)
			continue
		}

		fmt.Printf("Retry service affected %d rows.\n", rowsAffected)

	}
}
