module github.com/MahikaJaguste/distributed-task-queue/task-worker

go 1.22.5

require (
	github.com/MahikaJaguste/distributed-task-queue/task-submission v0.0.0-20250322153200-7cc5ee922ad4
	github.com/go-sql-driver/mysql v1.9.1
	github.com/joho/godotenv v1.5.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/albrow/forms v0.3.3 // indirect
)

replace github.com/MahikaJaguste/distributed-task-queue/task-submission v0.0.0-20250322153200-7cc5ee922ad4 => ../task-submission
