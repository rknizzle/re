package main

import (
	"log"
	"os"
)

// CLI entrypoint
func main() {
	// get the input command
	commandInput := os.Args[1:]
	err := run(commandInput, make(chan bool))
	if err != nil {
		log.Fatalf(err.Error())
	}
	return
}
