package mmv

import (
	"fmt"
	"os"
)

func validateIsDir(path string) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", path)
	}
	return nil
}

func validateNewPaths(old, new []string) error {
	if len(new) > len(old) {
		return fmt.Errorf("There are more new paths than old. mmv can't create files or directories")
	}
	if len(old) > len(new) {
		return fmt.Errorf("There are more new paths than old. mmv can't delete files or directories")
	}
	// TODO
	//	- Check number of files is the same
	//	- Check all directories exist
	return nil
}
