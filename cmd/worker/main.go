package main

import (
	"flag"

	"github.com/MahikaJaguste/distributed-task-queue/pkg/worker"
)

func main() {
	flag.IntVar(&worker.CONCURRENCY, "concurrency", 5, "concurrency limit for  worker")
	flag.Parse()
	worker.StartWorkerServer()
}
