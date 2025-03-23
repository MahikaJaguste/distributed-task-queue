package coordinator

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/MahikaJaguste/distributed-task-queue/pkg/common/db"
)

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

func StartCoordinatorServer() {
	err := godotenv.Load(db.ENV_FILE_PATH)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	mux := http.NewServeMux()
	// mux.HandleFunc("POST /heartbeat", handleReceiveHeartBeat)

	db.SetupDb()

	port := 8000
	coordinatorServer.Port = port
	coordinatorServer.WorkerList = make([]WorkerInfo, 0)
	// coordinatorServer.WorkerList = append(coordinatorServer.WorkerList, WorkerInfo{
	// 	Addr: "localhost",
	// 	Port: 8001,
	// })
	// coordinatorServer.WorkerList = append(coordinatorServer.WorkerList, WorkerInfo{
	// 	Addr: "localhost",
	// 	Port: 8002,
	// })

	go distributeTasksJob()

	fmt.Printf("Server listening on %d!\n", port)
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

		tasks, err := scanTasks()
		if err != nil {
			fmt.Println("Error in scanning tasks")
			continue
		}
		fmt.Println("tasks:", tasks)

		distributeTasks(tasks)

	}
}
func scanTasks() ([]db.Task, error) {
	fmt.Println("Inside scanTasks")
	var tasks []db.Task

	rows, err := db.DBCon.Query("SELECT id, name FROM tasks WHERE pickedAt is NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task db.Task
		if err := rows.Scan(&task.Id, &task.Name); err != nil {
			// todo
			fmt.Println("Error in scanning row")
			return nil, err
		}
		fmt.Println(task)
		_, err := db.DBCon.Exec("UPDATE tasks SET pickedAt = ? WHERE id=? AND pickedAt is NULL", time.Now(), task.Id)

		if err != nil {
			fmt.Println("Error in updating pickedAt")
			continue
		}
		tasks = append(tasks, task)

	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error in rows")
		return nil, err
	}
	return tasks, nil
}

func distributeTasks(tasks []db.Task) {
	for _, task := range tasks {
		worker, err := getNextWorker()
		if err != nil {
			fmt.Println("Error in getting next worker")
			break
		}
		go assignTask(task, worker)
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

func assignTask(task db.Task, worker WorkerInfo) {
	endpoint := fmt.Sprintf("http://%s:%d/execute", worker.Addr, worker.Port)
	v := url.Values{}
	v.Add("taskId", fmt.Sprintf("%d", task.Id))
	resp, err := http.PostForm(endpoint, v)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.Body)

}
