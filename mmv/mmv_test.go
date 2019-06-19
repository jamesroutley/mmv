package mmv

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/xerrors"
)

func newMockMultiMover(s string, options ...func(*MultiMover)) *MultiMover {
	mover := &MultiMover{
		editor: &MockEditor{
			Message: s,
		},
	}
	for _, option := range options {
		option(mover)
	}
	return mover
}

func TestMultiMoveDir(t *testing.T) {
	cases := []struct {
		name         string
		initialFiles []string
		editorInput  string
		finalFiles   []string
	}{
		{
			name:         "Rename a file",
			initialFiles: []string{"a.txt", "b.txt", "d.txt"},
			editorInput:  "a.txt\nb.txt\nc.txt",
			finalFiles:   []string{"a.txt", "b.txt", "c.txt"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir, cleanup := tempDirWithFiles(t, tc.initialFiles)
			defer cleanup()

			mover := newMockMultiMover(tc.editorInput)
			err := mover.MultiMoveDir(dir)
			require.NoError(t, err)

			finalFiles, err := pathsInDir(dir)
			require.NoError(t, err)
			assert.Equal(t, finalFiles, tc.finalFiles)
		})
	}
}

func TestMultiMoveDirWithFile(t *testing.T) {
	file, err := ioutil.TempFile("", "mmv_test")
	require.NoError(t, err)
	defer os.Remove(file.Name())

	mover := newMockMultiMover("") // input doesn't matter
	err = mover.MultiMoveDir(file.Name())
	fmt.Println(err.Error())
	assert.Error(t, err)
	assert.True(t, xerrors.Is(err, ErrNotDirectory))
}

func TestCreateAndDeleteFile(t *testing.T) {
	cases := []struct {
		name         string
		initialFiles []string
		editorInput  string
		expectedErr  error
	}{
		{
			name:         "Try to create",
			initialFiles: []string{"a.txt"},
			editorInput:  "a.txt\nb.txt",
			expectedErr:  ErrCannotCreate,
		},
		{
			name:         "Try to delete",
			initialFiles: []string{"a.txt", "b.txt"},
			editorInput:  "a.txt",
			expectedErr:  ErrCannotDelete,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir, cleanup := tempDirWithFiles(t, tc.initialFiles)
			defer cleanup()

			mover := newMockMultiMover(tc.editorInput)
			err := mover.MultiMoveDir(dir)
			assert.Error(t, err)
			assert.True(t, xerrors.Is(err, tc.expectedErr))
		})
	}
}

func tempDirWithFiles(t *testing.T, filepaths []string) (name string, cleanup func()) {
	t.Helper()
	dir, err := ioutil.TempDir("", "mmv_test")
	require.NoError(t, err)

	for _, path := range filepaths {
		fullPath := filepath.Join(dir, path)

		// Create any intermediate directories
		fileDir := filepath.Dir(fullPath)
		err := os.MkdirAll(fileDir, os.ModePerm)
		require.NoError(t, err)

		_, err = os.Create(fullPath)
		require.NoError(t, err)
	}

	return dir, func() { os.RemoveAll(dir) }
}

// allFilesInDir recursively searches dir for all files. It returns their
// paths, relative to `dir`
func allFilesInDir(t *testing.T, dir string) []string {
	t.Helper()
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		require.NoError(t, err)
		if info.IsDir() {
			return nil
		}
		relpath, err := filepath.Rel(dir, path)
		require.NoError(t, err)
		paths = append(paths, relpath)
		return nil
	})
	require.NoError(t, err)
	return paths
}
