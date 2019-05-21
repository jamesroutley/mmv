// Copyright © 2019 James Routley <jroutley@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mmv

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func MultiMoveDir(path string) error {
	paths, err := pathsInDir(path)
	if err != nil {
		return err
	}

	newPaths, err := multiEdit(paths)
	if err != nil {
		return err
	}

	newPaths = cleanPaths(newPaths)
	if err := validateNewPaths(paths, newPaths); err != nil {
		return err
	}

	// Individually move each old path to each new path
	// TODO: there's probably a 'correct' order to do this in
	for i, oldPath := range paths {
		oldPath = filepath.Join(path, oldPath)
		newPath := filepath.Join(path, newPaths[i])
		if err := rename(oldPath, newPath, true); err != nil {
			return err
		}
	}

	return nil
}

// BatchEdit takes a string slice, puts each string on a new line in a file,
// opens an editor, prompts the user to edit them, then returns the edited
// strings
func multiEdit(items []string) ([]string, error) {
	// Create a temporary file and write the items to it
	f, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		return nil, err
	}
	defer os.Remove(f.Name())
	for _, item := range items {
		fmt.Fprintf(f, "%s\n", item)
	}

	// Open edit the in a text editor
	if err := editFile(f.Name()); err != nil {
		return nil, err
	}

	// Read the new contents back
	bytes, err := ioutil.ReadFile(f.Name())
	s := string(bytes)
	// Remove any trailing newlines - some editors insert them
	s = strings.TrimSpace(s)
	return strings.Split(s, "\n"), nil
}

// editFile opens the file pointed to by `filename` in the user's chosen editor
func editFile(filename string) error {
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
