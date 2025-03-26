package worker

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/MahikaJaguste/distributed-task-queue/pkg/common/db"
)

var CONCURRENCY int
var WORKER_ID int
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
		fmt.Println("Error in registering worker")
		log.Fatal(err)
	}

	handleTaskExecution()
}

func setWorkerId() error {
	ctx := context.Background()

	tx, err := db.DBCon.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	_, err = tx.Exec("update workerCount set count = count+1 WHERE id = 1")

	if err != nil {
		return err
	}

	row := tx.QueryRow("select count from workerCount where id = 1")

	if err := row.Scan(&WORKER_ID); err != nil {
		return err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return err
	}

	fmt.Printf("Worker ID: %d\n", WORKER_ID)
	return nil
}

func handleTaskExecution() {
	ticker := time.NewTicker(SCAN_INTERVAL)
	defer ticker.Stop()

	parallelWorkers := make(chan struct{}, CONCURRENCY)

	for range ticker.C {

		select {
		case parallelWorkers <- struct{}{}: // Only proceed if there's a free slot
		default:
			// CONCURRENCY workers already present
			fmt.Println("Max parallelism reached, skipping this cycle")
			continue // Skip iteration if all worker slots are occupied
		}

		taskId, err := scanTasks()
		if err != nil {
			fmt.Println("Error in scanning tasks")
			fmt.Println(err)
			<-parallelWorkers
			continue
		} else if taskId == -1 {
			fmt.Println("No more tasks currently")
			<-parallelWorkers
			continue
		}

		fmt.Println("TaskId:", taskId)
		fmt.Printf("Task with id = %d is picked\n", taskId)

		go func() {
			execute(taskId)
			<-parallelWorkers // removes a struct from parallelWorkers, allowing another to proceed
		}()
	}
}

func scanTasks() (int, error) {
	fmt.Println("Inside scanTasks")
	var taskId int

	ctx := context.Background()

	tx, err := db.DBCon.BeginTx(ctx, nil)
	if err != nil {
		return taskId, err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	row := tx.QueryRow("select id from tasks where status = ? limit 1 for update skip locked", db.Pending)

	if err := row.Scan(&taskId); err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		}
	}

	_, err = tx.Exec("update tasks set pickedAt = now(), processedAt = now(), workerId = ?, status = ? WHERE id = ?", WORKER_ID, db.Processing, taskId)
	if err != nil {
		fmt.Println("Error in updating pickedAt")
		return -1, err
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return -1, err
	}

	return taskId, nil
}

func execute(taskId int) {
	ctx, cancelHeartbeat := context.WithCancel(context.Background())
	defer cancelHeartbeat()

	go func(ctx context.Context, taskId int) {
		sendHeartbeat(ctx, taskId)
	}(ctx, taskId)

	row := db.DBCon.QueryRow("select id, name from tasks where id=?", taskId)
	var id int
	var name string
	if err := row.Scan(&id, &name); err != nil {
		if err == sql.ErrNoRows {
			// TODO
			fmt.Printf("No such task with taskId = %d\n", taskId)
		} else {
			fmt.Println("Error in scanning row")
			fmt.Println(err)
		}
		return
	}

	fmt.Printf("Starting sleep for taskId = %d at %d\n", taskId, time.Now().Unix())
	fmt.Printf("Task %d: %s\n", taskId, name)
	time.Sleep(time.Minute)
	fmt.Printf("Sleep done for taskId = %d at %d\n", taskId, time.Now().Unix())

	cancelHeartbeat()

	fmt.Printf("Updating completed at for taskId = %d at %d\n", taskId, time.Now().Unix())
	_, err := db.DBCon.Exec("update tasks set completedAt = now(), status = ? WHERE id = ? and workerId = ?", db.Completed, taskId, WORKER_ID)
	if err != nil {
		fmt.Println("Error in updating completedAt")
		return
	}

}

func sendHeartbeat(ctx context.Context, taskId int) {
	heartbeatTicker := time.NewTicker(HEARTBEAT_DURATION)
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-ctx.Done(): // if cancelHeartbeat() execute
			fmt.Printf("Shutting down hearbeat for taskId = %d at %d\n", taskId, time.Now().Unix())
			return
		case <-heartbeatTicker.C:
			fmt.Printf("Sending hearbeat for taskId = %d at %d\n", taskId, time.Now().Unix())
			_, err := db.DBCon.ExecContext(ctx, "update tasks set processedAt = now() WHERE id = ? and workerId = ?", taskId, WORKER_ID)
			if err != nil {
				if err == context.Canceled {
					fmt.Printf("Shutting down hearbeat for taskId = %d at %d\n", taskId, time.Now().Unix())
					return
				}
				fmt.Printf("Error in sending heartbeat for taskId = %d at %d\n", taskId, time.Now().Unix())
				fmt.Println(err)
			}
		}
	}
}
