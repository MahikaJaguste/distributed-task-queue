package submission

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MahikaJaguste/distributed-task-queue/pkg/common/db"
	"github.com/albrow/forms"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var STREAM_KEY string

func StartSubmissionServer(port int) {
	err := godotenv.Load(db.ENV_FILE_PATH)
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /submit", handleTaskSubmission)

	db.SetupDb()
	STREAM_KEY = os.Getenv("STREAM_KEY")

	log.Printf("Submission server listening on %d!\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		log.Fatal(err)
	}
}

func handleTaskSubmission(w http.ResponseWriter, req *http.Request) {
	data, err := forms.Parse(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := data.Get("name")

	args := redis.XAddArgs{
		Stream: STREAM_KEY,
		ID:     "*",
		Values: map[string]interface{}{
			"name": name,
		},
	}
	result := db.RDB.XAdd(context.Background(), &args)
	id, err := result.Result()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	log.Printf("Task created with id = %s\n", id)
	w.Write(([]byte)(fmt.Sprintf("Task created with id = %s\n", id)))
}
