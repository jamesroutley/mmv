package mmv

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

type FileEditor interface {
	Edit(filename string) error
}

type DefaultEditor struct{}

func (e *DefaultEditor) Edit(filename string) error {
	editor := os.Getenv("VISUAL")
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		editor = "vim"
	}
	cmd := exec.Command("sh", "-c", fmt.Sprintf("%s %s", editor, filename))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

type MockEditor struct {
	// This message is written to the file when Edit is called.
	Message string
}

func (e *MockEditor) Edit(filename string) error {
	return ioutil.WriteFile(filename, []byte(e.Message), 0644)
}
