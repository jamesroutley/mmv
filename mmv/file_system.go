package mmv

import (
	"io/ioutil"
	"path/filepath"
)

func pathsInDir(path string) ([]string, error) {
	if err := validateIsDir(path); err != nil {
		return nil, err
	}
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	paths := make([]string, len(fileInfos))
	for i, info := range fileInfos {
		paths[i] = info.Name()
	}
	return paths, nil
}

func cleanPaths(paths []string) []string {
	clean := make([]string, len(paths))
	for i, path := range paths {
		clean[i] = filepath.Clean(path)
	}
	return clean
}
