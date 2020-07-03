package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
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

/*
	Attempt to clear the terminal screen depending on the OS. Print a warning
	message if the system is not supported
	NOTE: only tested on macOS
*/
var clearScreen = func() {
	switch system := runtime.GOOS; system {
	case "darwin":
		func() {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
		}()
	case "linux":
		func() {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
		}()
	case "windows":
		func() {
			cmd := exec.Command("cmd", "/c", "cls")
			cmd.Stdout = os.Stdout
			cmd.Run()
		}()
	default:
		fmt.Println("Platform unsupported! Could not clear the screen")
	}
}
