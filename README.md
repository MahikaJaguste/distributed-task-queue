# distributed-task-queue

A task queue built using MySQL and Go to experiment with locking and concurrency. 

### Code Structure

1. `pkg/submission`: Starts a server which exposes endpoint to submit tasks.

```
curl --location '127.0.0.1:8000/submit' \
--header 'Content-Type: application/json' \
--data '{
    "name": "Task 1"
}'
```

2. `pkg/worker`: Starts a worker service which picks pending task from queue every `SCAN_INTERVAL`, assigns it to itself, processing it while sending heartbeats to the database and marks it as complete. Simultaneously, the worker can spawn `CONCURRENCY` number of goroutines, each handling a task. 


3. `pkg/retry`: Retry service to detect failed tasks (tasks where 3 heartbeats were missed) and puts them back into pending status for available workers to pick.


4. `common/db`: Contains SQL statements used to setup the database tables and experiment with locking features by pausing transactions.


### Running the Code

Running the task submission service inside `/cmd/submission`: `go build -o ../build && ../build/submission -port=8080`

Running the worker service inside `/cmd/worker`: `go build -o ../build && ../build/worker`

Running the retry service inside `/cmd/retry`: `go build -o ../build && ../build/retry`