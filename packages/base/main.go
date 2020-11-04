package base

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/i-hate-nicknames/gitik/packages/data"
)

type treeEntry struct {
	name, oid string
	otype     data.ObjectType
}

func (te treeEntry) String() string {
	return fmt.Sprintf("%s %s %s", te.otype.String(), te.oid, te.name)
}

func WriteTree(directory string) (string, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return "", nil
	}
	var entries []treeEntry
	for _, f := range files {
		if isIgnored(f.Name()) {
			continue
		}
		var entry treeEntry
		fullPath := filepath.Join(directory, f.Name())
		if f.IsDir() {
			oid, err := WriteTree(fullPath)
			if err != nil {
				return "", err
			}
			entry = treeEntry{name: f.Name(), oid: oid, otype: data.TypeTree}
		} else if f.Mode().IsRegular() {
			oid, err := data.HashObject(fullPath, data.TypeBlob)
			if err != nil {
				return "", err
			}
			entry = treeEntry{name: f.Name(), oid: oid, otype: data.TypeBlob}
		}
		entries = append(entries, entry)
	}
	var lines []string
	for _, entry := range entries {
		lines = append(lines, entry.String())
	}

	return data.HashAndStore([]byte(strings.Join(lines, "\n")), data.TypeTree)
}

func isIgnored(path string) bool {
	return path == data.GitDir
}
