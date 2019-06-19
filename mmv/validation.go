package mmv

import (
	"os"

	"golang.org/x/xerrors"
)

var (
	ErrNotDirectory = xerrors.New("is not a directory")
	ErrCannotCreate = xerrors.New("there are more new paths than old. mmv can't create files or directories")
	ErrCannotDelete = xerrors.New("there are more new paths than old. mmv can't delete files or directories")
)

func validateIsDir(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fileInfo.IsDir() {
		return xerrors.Errorf("%s: %w", path, ErrNotDirectory)
	}
	return nil
}

func validateNewPaths(old, new []string) error {
	if len(new) > len(old) {
		return ErrCannotCreate
	}
	if len(old) > len(new) {
		return ErrCannotDelete
	}
	// TODO
	//	- Check all directories exist
	return nil
}
