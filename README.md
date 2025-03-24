# distributed-task-queue

Running the coordinator service inside `/cmd/coordinator`: `go build -o ../build && ../build/coordinator -port=8000`
Running the worker service inside `/cmd/worker`: `go build -o ../build && ../build/worker -port=8001`
Running the task submission service inside `/cmd/submission`: `go build -o ../build && ../build/submission -port=8080`