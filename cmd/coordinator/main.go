package main

import (
	"flag"

	"github.com/MahikaJaguste/distributed-task-queue/pkg/coordinator"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8000, "port for container")
	flag.Parse()
	coordinator.StartCoordinatorServer(port)
}
