package plumbing

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/i-hate-nicknames/gitik/packages/storage"
)

type ObjectType string

const (
	TypeBlob ObjectType = "blob"
	TypeTree ObjectType = "tree"
)

func (t ObjectType) String() string {
	switch t {
	case TypeBlob:
		return "blob"
	case TypeTree:
		return "tree"
	default:
		return "_unknown"
	}
}

type treeEntry struct {
	name  string
	oid   storage.OID
	otype ObjectType
}

func (te treeEntry) String() string {
	return fmt.Sprintf("%s %s %s", te.otype.String(), te.oid, te.name)
}

func WriteFile(fileName string) (storage.OID, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return storage.HashObject(data, []byte(TypeBlob))
}

func ReadFile(objectID storage.OID) (string, error) {
	data, err := storage.GetObject(objectID, []byte(TypeBlob))
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(data)
	return fmt.Sprintf(buf.String()), nil
}

func WriteTree(directory string) (storage.OID, error) {
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
			entry = treeEntry{name: f.Name(), oid: oid, otype: TypeTree}
		} else if f.Mode().IsRegular() {
			oid, err := WriteFile(fullPath)
			if err != nil {
				return "", err
			}
			entry = treeEntry{name: f.Name(), oid: oid, otype: TypeBlob}
		}
		entries = append(entries, entry)
	}
	var lines []string
	for _, entry := range entries {
		lines = append(lines, entry.String())
	}

	return storage.HashObject([]byte(strings.Join(lines, "\n")), []byte(TypeTree))
}

func isIgnored(path string) bool {
	return path == storage.GitDir
}
