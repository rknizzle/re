package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestMainThrowsAnErrorWhenGivenNoCommandArguments(t *testing.T) {
	err := run([]string{}, make(chan bool))
	if err != nil {
		if err.Error() != "No arguments supplied" {
			t.Errorf("Was expecting the error 'No arguments supplied' but instead got '%s'", err.Error())
		}
	}
}

// test that when the program starts up it runs the input command
func TestCommandRunsOnStartup(t *testing.T) {
	clearScreen = func() {}
	before := initializeCmd
	var runCount int = 0
	// replace the intializeCmd fxn with a fxn that increments a called counter
	// so that the test can verify if initializeCmd was called
	initializeCmd = func([]string) *exec.Cmd {
		runCount++
		cmd := exec.Command("echo", "")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd
	}

	// un-mock initializeCmd after the test
	defer func() {
		initializeCmd = before
	}()

	done := make(chan bool)
	// start the run function
	go run([]string{"blank", "echo", ""}, done)

	// this will block until the run function has ran the input command and is
	// now monitoring the files to re-run the command on file change
	done <- true

	// verify that initializeCmd was called
	if runCount == 0 {
		t.Errorf("initializeCmd was not called")
	} else if runCount != 1 {
		t.Errorf("Expected initializeCmd to be called once but it was called %d times", runCount)
	}
}
