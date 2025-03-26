package retry

import (
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/MahikaJaguste/distributed-task-queue/pkg/common/db"
	"github.com/MahikaJaguste/distributed-task-queue/pkg/worker"
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
	log.Println("Handling failed tasks")

	expiryDuration := worker.HEARTBEAT_DURATION * 3 // ie. failed 3 heartbeats
	expiry := time.Now().Add(time.Duration(-expiryDuration) * time.Second).Format(time.DateTime)
	result, err := db.DBCon.Exec("update tasks set pickedAt = null, processedAt = null, workerId = null, status = ? WHERE status = ? and processedAt < ?", db.Pending, db.Processing, expiry)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(rowsAffected), nil

}

func handleRetry() {
	for range time.Tick(worker.HEARTBEAT_DURATION * 3) {

		rowsAffected, err := handleFailedTasks()
		if err != nil {
			log.Println("Error in handling failed tasks: ", err)
			continue
		}

		log.Printf("Retry service affected %d rows.\n", rowsAffected)

	}
}
