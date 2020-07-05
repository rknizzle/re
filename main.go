package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var watcher *fsnotify.Watcher

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

// watchDir gets run as a walk func, searching for directories to add watchers
func watchDir(path string, fi os.FileInfo, err error) error {
	// ignore node_modules/ and .git/ because they cause the watchers to trigger
	if strings.Contains(path, "node_modules") || strings.Contains(path, ".git") {
		return nil
	}

	// since fsnotify can watch all the files in a directory, watchers only need
	// to be added to each nested directory
	if fi.Mode().IsDir() {
		return watcher.Add(path)
	}

	return nil
}
