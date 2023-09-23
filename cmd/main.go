package main

import (
	"fmt"
	"log"

	"github.com/raja-dettex/dap/api"
)

func main() {
	// test cases and instructions to sprin up a database instance
	defer handlePanic()
	dbServer, err := api.New(":3000", "default")
	if err != nil {
		log.Fatal(err)
	}
	dbServer.RegisterHandlers()
	err = dbServer.Start()
	if err != nil {
		panic(err)
	}

}

func handlePanic() {
	if err := recover(); err != nil {
		fmt.Printf("panic raised %d", err)
	}
}
