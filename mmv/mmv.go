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
	"path/filepath"
	"regexp"
	"strings"
)

type MultiMover struct {
	editor   FileEditor
	dryRun   bool
	includes *regexp.Regexp
	excludes *regexp.Regexp
}

func NewMultiMover(options ...func(*MultiMover)) *MultiMover {
	mover := &MultiMover{
		editor: &DefaultEditor{},
	}
	for _, option := range options {
		option(mover)
	}
	return mover
}

func (mm *MultiMover) MultiMoveDir(path string) error {
	paths, err := pathsInDir(path)
	if err != nil {
		return err
	}

	paths = filterStringSlice(paths, func(s string) bool {
		if mm.includes == nil {
			return true
		}
		return mm.includes.MatchString(s)
	})
	paths = filterStringSlice(paths, func(s string) bool {
		if mm.excludes == nil {
			return true
		}
		return !mm.excludes.MatchString(s)
	})

	newPaths, err := multiEdit(paths, mm.editor)
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
		if err := mm.rename(oldPath, newPath); err != nil {
			return err
		}
	}

	return nil
}

func (mm *MultiMover) rename(old, new string) error {
	if old == new {
		return nil
	}
	if mm.dryRun {
		fmt.Printf("mv %s %s\n", old, new)
		return nil
	}
	return os.Rename(old, new)
}

// multiEdit takes a string slice, puts each string on a new line in a file,
// opens an editor, prompts the user to edit them, then returns the edited
// strings
func multiEdit(items []string, editor FileEditor) ([]string, error) {
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
	if err := editor.Edit(f.Name()); err != nil {
		return nil, err
	}

	// Read the new contents back
	bytes, err := ioutil.ReadFile(f.Name())
	s := string(bytes)
	// Remove any trailing newlines - some editors insert them
	s = strings.TrimSpace(s)
	return strings.Split(s, "\n"), nil
}

func filterStringSlice(ss []string, include func(string) bool) []string {
	var filtered []string
	for _, s := range ss {
		if include(s) {
			filtered = append(filtered, s)
		}
	}
	return filtered
}
