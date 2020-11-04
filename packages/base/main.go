package base

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/i-hate-nicknames/gitik/packages/data"
)

func WriteTree(directory string) (string, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return "", nil
	}
	for _, f := range files {
		if isIgnored(f.Name()) {
			continue
		}
		fullPath := filepath.Join(directory, f.Name())
		if f.IsDir() {
			WriteTree(fullPath)
		}
		fmt.Println(fullPath)
	}
	return "", nil
}

func isIgnored(path string) bool {
	return path == data.GitDir
}
