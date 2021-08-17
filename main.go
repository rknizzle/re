package main

import (
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
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

func run(commandInput []string, done chan bool) error {
	// check for no command input
	if len(commandInput) == 0 {
		return errors.New("No arguments supplied")
	}

	// run the input command initially
	clearScreen()
	cmd := initializeCmd(commandInput)
	err := cmd.Start()
	if err != nil {
		return err
	}

	// start watching for file changes
	watcher, err = deployWatchers()
	if err != nil {
		return err
	}
	defer watcher.Close()

	// refresh the command any time a file is changed
	err = monitor(cmd, commandInput, done)
	if err != nil {
		return err
	}
	return nil
}

//Attempt to clear the terminal screen depending on the OS. Print a warning
//message if the system is not supported
//NOTE: only tested on macOS
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

// recursively add watchers to directories to capture file change events
func deployWatchers() (*fsnotify.Watcher, error) {
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// creates a new file watcher which watches all files for changes in the
	// directory that it is placed in
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	// starting at the root of the project, walk each file/directory searching
	// for directories to add a watcher to
	if err := filepath.Walk(path, watchDir); err != nil {
		return nil, err
	}
	return watcher, nil
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

// convert the command input -> an executable command (exec.Command) that
// will output its stdout and stderr to the console
var initializeCmd = func(commandInput []string) *exec.Cmd {
	cmd := exec.Command(commandInput[0], commandInput[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// reruns the command when a file change is detected
func monitor(cmd *exec.Cmd, commandInput []string, done chan bool) error {
	var readyForCommand = true
	// watch for events to refresh the command
	for {
		select {
		// watch for events
		case event := <-watcher.Events:
			if !(event.Op&fsnotify.Write == fsnotify.Write) {
				break
			}

			if readyForCommand {
				go pauseMonitoring(&readyForCommand)
				// kill the previous commands process before restarting
				err := cmd.Process.Kill()
				if err != nil {
					return err
				}

				// run the command
				clearScreen()
				cmd = initializeCmd(commandInput)
				err = cmd.Start()
				if err != nil {
					return err
				}
			}

		// watch for errors
		case err := <-watcher.Errors:
			return err

		// exit when 'done'
		case _ = <-done:
			return nil
		}
	}
}

// After an event is triggered, have a waiting period before
// another event can be triggered. Its common for a file save
// to trigger multiple events so creating a short artificial
// delay between triggers makes the program feel smoother
func pauseMonitoring(readyForCommand *bool) {
	*readyForCommand = false
	time.Sleep(2 * time.Second)
	*readyForCommand = true
}
