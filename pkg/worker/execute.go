package worker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/MahikaJaguste/distributed-task-queue/pkg/common/db"
)

var CONCURRENCY int
var WORKER_ID string
var HEARTBEAT_DURATION = time.Second * 10
var SCAN_INTERVAL = time.Second * 10

func StartWorkerServer() {
	err := godotenv.Load(db.ENV_FILE_PATH)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db.SetupDb()
	err = setWorkerId()
	if err != nil {
		log.Fatal(err)
	}

	// handleTaskExecution()
}

func setWorkerId() error {
	flag := false

	for range 5 {
		uuidBytes := uuid.New()
		uuidStr := uuidBytes.String()

		if len(uuidStr) == 0 {
			continue
		}

		key := fmt.Sprintf("WORKER_IDS:%s", uuidStr)

		result := db.RDB.SetNX(context.Background(), key, 1, 0)
		resp, err := result.Result()
		if !resp || err != nil {
			continue
		}

		WORKER_ID = uuidStr
		flag = true
		break
	}

	if !flag {
		return errors.New("error in registering worker")
	}

	log.Printf("Worker ID: %s\n", WORKER_ID)
	return nil

}

// func handleTaskExecution() {
// 	ticker := time.NewTicker(SCAN_INTERVAL)
// 	defer ticker.Stop()

// 	parallelWorkers := make(chan struct{}, CONCURRENCY)
// 	defer close(parallelWorkers)

// 	for range ticker.C {

// 		log.Printf("Number of Running Goroutines: %d\n", runtime.NumGoroutine())

// 		select {
// 		case parallelWorkers <- struct{}{}: // Only proceed if there's a free slot
// 		default:
// 			// CONCURRENCY workers already present
// 			log.Println("Max parallelism reached, skipping this cycle")
// 			continue // Skip iteration if all worker slots are occupied
// 		}

// 		taskId, err := scanTasks()
// 		if err != nil {
// 			log.Println("Error in scanning tasks: ", err)
// 			<-parallelWorkers
// 			continue
// 		} else if taskId == -1 {
// 			log.Println("No more tasks currently")
// 			<-parallelWorkers
// 			continue
// 		}

// 		log.Println("TaskId: ", taskId)
// 		log.Printf("Task with id = %d is picked\n", taskId)

// 		go func() {
// 			execute(taskId)
// 			<-parallelWorkers // removes a struct from parallelWorkers, allowing another to proceed
// 		}()
// 	}
// }

// func scanTasks() (int, error) {
// 	log.Println("Scanning tasks")
// 	var taskId int

// 	ctx := context.Background()

// 	tx, err := db.DBCon.BeginTx(ctx, nil)
// 	if err != nil {
// 		return taskId, err
// 	}
// 	// Defer a rollback in case anything fails.
// 	defer tx.Rollback()

// 	row := tx.QueryRow("select id from tasks where status = ? limit 1 for update of tasks skip locked", db.Pending)

// 	if err := row.Scan(&taskId); err != nil {
// 		if err == sql.ErrNoRows {
// 			return -1, nil
// 		}
// 	}

// 	_, err = tx.Exec("update tasks set pickedAt = now(), processedAt = now(), workerId = ?, status = ? WHERE id = ?", WORKER_ID, db.Processing, taskId)
// 	if err != nil {
// 		log.Println("Error in updating pickedAt")
// 		return -1, err
// 	}

// 	// Commit the transaction.
// 	if err = tx.Commit(); err != nil {
// 		return -1, err
// 	}

// 	return taskId, nil
// }

// func execute(taskId int) {
// 	ctx, cancelHeartbeat := context.WithCancel(context.Background())
// 	defer cancelHeartbeat()

// 	go func(ctx context.Context, taskId int) {
// 		sendHeartbeat(ctx, taskId)
// 	}(ctx, taskId)

// 	row := db.DBCon.QueryRow("select id, name from tasks where id=?", taskId)
// 	var id int
// 	var name string
// 	if err := row.Scan(&id, &name); err != nil {
// 		if err == sql.ErrNoRows {
// 			log.Printf("No such task with taskId = %d\n", taskId)
// 		} else {
// 			log.Println("Error in scanning row: ", err)
// 		}
// 		return
// 	}

// 	log.Printf("Starting work for taskId = %d\n", taskId)
// 	log.Printf("Task %d: %s\n", taskId, name)
// 	time.Sleep(time.Minute)
// 	log.Printf("Completed work for taskId = %d\n", taskId)

// 	cancelHeartbeat()

// 	log.Printf("Updating completed at for taskId = %d\n", taskId)
// 	_, err := db.DBCon.Exec("update tasks set completedAt = now(), status = ? WHERE id = ? and workerId = ?", db.Completed, taskId, WORKER_ID)
// 	if err != nil {
// 		log.Println("Error in updating completedAt: ", err)
// 		return
// 	}

// }

// func sendHeartbeat(ctx context.Context, taskId int) {
// 	heartbeatTicker := time.NewTicker(HEARTBEAT_DURATION)
// 	defer heartbeatTicker.Stop()

// 	for {
// 		select {
// 		case <-ctx.Done(): // if cancelHeartbeat() execute
// 			log.Printf("Shutting down hearbeat for taskId = %d\n", taskId)
// 			return
// 		case <-heartbeatTicker.C:
// 			log.Printf("Sending hearbeat for taskId = %d\n", taskId)
// 			_, err := db.DBCon.ExecContext(ctx, "update tasks set processedAt = now() WHERE id = ? and workerId = ?", taskId, WORKER_ID)
// 			if err != nil {
// 				if err == context.Canceled {
// 					log.Printf("Shutting down hearbeat for taskId = %d\n", taskId)
// 					return
// 				}
// 				log.Println(fmt.Sprintf("Error in sending heartbeat for taskId = %d", taskId), err)
// 			}
// 		}
// 	}
// }
