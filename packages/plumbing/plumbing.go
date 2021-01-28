package plumbing

import (
	"bytes"
	"errors"
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

func parseEntry(data []byte) (treeEntry, error) {
	parts := bytes.Split(data, []byte(" "))
	if len(parts) != 3 {
		return treeEntry{}, fmt.Errorf("parseEntry: wrong length (%d), should be 3", len(parts))
	}
	var otype ObjectType
	switch header := string(parts[0]); header {
	case "blob":
		otype = TypeBlob
	case "tree":
		otype = TypeTree
	default:
		return treeEntry{}, fmt.Errorf("parseEntry: wrong entry type %s", header)
	}
	oid := storage.OID(parts[1])
	name := string(parts[2])
	return treeEntry{name, oid, otype}, nil
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

func ReadTree(oid storage.OID) error {
	entries, err := readTreeEntries(oid)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		fmt.Println(entry)
	}
	return nil
}

var errEmptyTree = errors.New("empty tree")

func readTreeEntries(oid storage.OID) ([]treeEntry, error) {
	data, err := storage.GetObject(oid, []byte(TypeTree))
	if err != nil {
		return nil, err
	}
	var entries []treeEntry
	if len(data) == 0 {
		return nil, errEmptyTree
	}
	entriesRaw := bytes.Split(data, []byte("\n"))
	for _, eraw := range entriesRaw {
		entry, err := parseEntry(eraw)
		if err != nil {
			return nil, err
		}
		switch entry.otype {
		case TypeBlob:
			entries = append(entries, entry)
		case TypeTree:
			children, err := readTreeEntries(entry.oid)
			if err == errEmptyTree {
				continue
			}
			if err != nil {
				return nil, err
			}
			entries = append(entries, children...)
		default:
			return nil, fmt.Errorf("readTreeEntries: unknown object type")
		}
	}
	return entries, nil
}

func isIgnored(path string) bool {
	return path == storage.GitDir
}
