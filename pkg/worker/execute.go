package worker

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	forms "github.com/albrow/forms"

	"github.com/MahikaJaguste/distributed-task-queue/pkg/common/db"
)

type WorkerResponse struct {
	Id     int
	Body   string
	Status string
}

func HandleTaskExecution(w http.ResponseWriter, req *http.Request) {
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
	row := db.DBCon.QueryRow("select * from tasks where id=?", taskId)
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
