package worker

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	forms "github.com/albrow/forms"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/MahikaJaguste/distributed-task-queue/pkg/common/db"
)

type WorkerResponse struct {
	Id     int
	Body   string
	Status string
}

func StartWorkerServer() {
	err := godotenv.Load(db.ENV_FILE_PATH)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /execute", handleTaskExecution)

	db.SetupDb()

	port := 8001
	fmt.Printf("Server listening on %d!\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Fatal(err)
	}
}

func handleTaskExecution(w http.ResponseWriter, req *http.Request) {
	data, err := forms.Parse(req)
	if err != nil {
		// in case of any error
		// TODO
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	taskId := data.GetInt("taskId")

	go execute(taskId)

	fmt.Printf("Task with id = %d is picked\n", taskId)

	response := WorkerResponse{
		Id:     taskId,
		Body:   fmt.Sprintf("Task with id = %d is picked", taskId),
		Status: "PICKED",
	}

	resBytes, _ := json.Marshal(response)
	w.Write(resBytes)
}

func execute(taskId int) {
	row := db.DBCon.QueryRow("select id, name from tasks where id=?", taskId)
	var id int
	var name string
	if err := row.Scan(&id, &name); err != nil {
		if err == sql.ErrNoRows {
			// TODO
			fmt.Printf("No such task with taskId = %d\n", taskId)
			return
		}
	}

	fmt.Printf("Task %d: %s\n", taskId, name)
}
