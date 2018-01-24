package main

import (
	"os"
	"strconv"

	"github.com/mastoj/proto-actor-demo/go-node/hello"
	"github.com/mastoj/proto-actor-demo/go-node/supervision"
	"github.com/mastoj/proto-actor-demo/go-node/worker"
)

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 || argsWithoutProg[0] == "h" {
		hello.Hello()
	} else if argsWithoutProg[0] == "s" {
		supervision.Run()
	} else if argsWithoutProg[0] == "w" {
		workerCount, _ := strconv.Atoi(argsWithoutProg[1])
		worker.Run(workerCount)
	}
	select {}
}
