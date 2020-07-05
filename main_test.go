package main

import (
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
