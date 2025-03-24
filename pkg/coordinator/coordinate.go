package coordinator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/MahikaJaguste/distributed-task-queue/pkg/common/db"
)

var FETCHING_LIMIT = 2

type WorkerInfo struct {
	Addr string
	Port int
	// nextHeartbeatDeadline time
}

type CoordinatorServer struct {
	Port       int
	WorkerList []WorkerInfo
}

var coordinatorServer CoordinatorServer
var nextWorkerId = 0

func StartCoordinatorServer(port int) {
	err := godotenv.Load(db.ENV_FILE_PATH)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	mux := http.NewServeMux()
	// mux.HandleFunc("POST /heartbeat", handleReceiveHeartBeat)

	db.SetupDb()

	coordinatorServer.Port = port
	coordinatorServer.WorkerList = make([]WorkerInfo, 0)
	coordinatorServer.WorkerList = append(coordinatorServer.WorkerList, WorkerInfo{
		Addr: "localhost",
		Port: 8001,
	})
	// coordinatorServer.WorkerList = append(coordinatorServer.WorkerList, WorkerInfo{
	// 	Addr: "localhost",
	// 	Port: 8002,
	// })

	go distributeTasksJob()

	fmt.Printf("Coordinator server listening on %d!\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Fatal(err)
	}
}

// func handleReceiveHeartBeat(w http.ResponseWriter, req *http.Request) {
// 	// TODO
// 	return
// }

func distributeTasksJob() {
	for range time.Tick(time.Second * 10) {

		if !hasAvailableWorkers() {
			fmt.Println("no available workers")
			continue
		}

		taskIds, err := scanTasks()
		if err != nil {
			fmt.Println("Error in scanning tasks")
			fmt.Println(err)
			continue
		}
		fmt.Println("taskIds:", taskIds)

		distributeTasks(taskIds)

	}
}
func scanTasks() ([]int, error) {
	fmt.Println("Inside scanTasks")
	var taskIds []int

	ctx := context.Background()

	tx, err := db.DBCon.BeginTx(ctx, nil)
	if err != nil {
		return taskIds, err
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, "select id from tasks where pickedAt is null limit ? for update skip locked", FETCHING_LIMIT)
	if err != nil {
		return taskIds, err
	}
	defer rows.Close()

	for rows.Next() {
		var taskId int
		if err := rows.Scan(&taskId); err != nil {
			fmt.Println("Error in scanning row")
			continue
		}
		fmt.Println(taskId)
		taskIds = append(taskIds, taskId)
	}

	if len(taskIds) != 0 {
		idStr := "(" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(taskIds)), ","), "[]") + ")"
		updateQuery := fmt.Sprintf("update tasks set pickedAt = now() WHERE id in %s", idStr)

		_, err = tx.ExecContext(ctx, updateQuery)
		if err != nil {
			fmt.Println("Error in updating pickedAt")
			return []int{}, err
		}
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return []int{}, err
	}

	return taskIds, nil
}

func distributeTasks(taskIds []int) {
	for _, taskId := range taskIds {
		worker, err := getNextWorker()
		if err != nil {
			fmt.Println("Error in getting next worker")
			break
		}
		go assignTask(taskId, worker)
	}
}

func getNextWorker() (WorkerInfo, error) {
	if len(coordinatorServer.WorkerList) == 0 {
		return WorkerInfo{}, errors.New("no workers available")
	}
	idx := nextWorkerId % len(coordinatorServer.WorkerList)
	nextWorkerId = idx + 1
	return coordinatorServer.WorkerList[idx], nil
}

func hasAvailableWorkers() bool {
	return len(coordinatorServer.WorkerList) != 0
}

func assignTask(taskId int, worker WorkerInfo) {
	endpoint := fmt.Sprintf("http://%s:%d/execute", worker.Addr, worker.Port)
	v := url.Values{}
	v.Add("taskId", fmt.Sprintf("%d", taskId))
	_, err := http.PostForm(endpoint, v)
	if err != nil {
		fmt.Println(err)
		return
	}
	// fmt.Println(resp.Body)

}
