package submission

import (
	"fmt"
	"net/http"

	"github.com/MahikaJaguste/distributed-task-queue/pkg/common/db"

	forms "github.com/albrow/forms"
)

func HandleTaskSubmission(w http.ResponseWriter, req *http.Request) {
	data, err := forms.Parse(req)
	if err != nil {
		// in case of any error
		// TODO
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := data.Get("name")

	result, err := db.DBCon.Exec("INSERT INTO tasks (name) VALUES (?)", name)
	if err != nil {
		// TODO
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		// TODO
		return
	}
	fmt.Printf("Task created with id = %d\n", id)
	w.Write(([]byte)(fmt.Sprintf("Task created with id = %d\n", id)))
}
