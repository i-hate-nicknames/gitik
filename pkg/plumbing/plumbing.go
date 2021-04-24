package plumbing

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/i-hate-nicknames/gitik/pkg/constants"
	"github.com/i-hate-nicknames/gitik/pkg/storage"
)

type treeEntry struct {
	name  string
	oid   storage.OID
	otype storage.ObjectType
}

func (te treeEntry) String() string {
	return fmt.Sprintf("%s %s %s", te.otype.String(), te.oid, te.name)
}

func parseEntry(data []byte) (treeEntry, error) {
	parts := bytes.Split(data, []byte(" "))
	if len(parts) != 3 {
		return treeEntry{}, fmt.Errorf("parseEntry: wrong length (%d), should be 3", len(parts))
	}
	otype, err := storage.Decode(parts[0])
	if err != nil {
		return treeEntry{}, err
	}
	oid := storage.OID(parts[1])
	name := string(parts[2])
	return treeEntry{name, oid, otype}, nil
}

// WriteFile writes contents of the given file path (relative to the root of the repository)
// to the object database. Return object id of the stored object
func WriteFile(fileName string) (storage.OID, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return storage.StoreObject(data, storage.TypeBlob)
}

// WriteTree writes contents of the given directory (relative to the root of the repository)
// to the object database. Return object id of the stored directory.
// Recursively writes all files found in the directory
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
			// todo: if tree wasn't written because it's empty, do not add it
			// to the entries
			if err != nil {
				return "", err
			}
			entry = treeEntry{name: f.Name(), oid: oid, otype: storage.TypeTree}
		} else if f.Mode().IsRegular() {
			oid, err := WriteFile(fullPath)
			if err != nil {
				return "", err
			}
			entry = treeEntry{name: f.Name(), oid: oid, otype: storage.TypeBlob}
		}
		entries = append(entries, entry)
	}
	var lines []string
	for _, entry := range entries {
		lines = append(lines, entry.String())
	}
	// todo: add empty tree error, and return it here when lines is empty,
	// instead of writing an empty tree
	return storage.StoreObject([]byte(strings.Join(lines, "\n")), storage.TypeTree)
}

// ReadTree reads directory under given storage id and writes it in the root
// directory of repository. The contents of root directory is removed before
// the write happens, but the ignored files are omitted
func ReadTree(oid storage.OID) error {
	entries, err := readTreeEntries(oid, ".")
	if err != nil {
		return err
	}
	err = emptyDir(".")
	if err != nil {
		return err
	}
	dirPerm := os.ModeDir | 0755
	for _, entry := range entries {
		// todo: consider storing permissions along with the name
		dirPath, _ := path.Split(entry.name)
		err := os.MkdirAll(dirPath, dirPerm)
		if err != nil {
			return err
		}
		data, err := readObject(oid, storage.TypeBlob)
		if err != nil {
			return err
		}
		err = storage.WriteFile(entry.name, data)
		if err != nil {
			return err
		}
	}
	return nil
}

var errEmptyTree = errors.New("empty tree")

func readTreeEntries(oid storage.OID, path string) ([]treeEntry, error) {
	data, err := readObject(oid, storage.TypeTree)
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
		if entry.name == ".." || entry.name == "." {
			return nil, fmt.Errorf("readTreeEntries: malformed entry, path %s, name %s", path, entry.name)
		}
		switch entry.otype {
		case storage.TypeBlob:
			entry.name = path + "/" + entry.name
			entries = append(entries, entry)
		case storage.TypeTree:
			children, err := readTreeEntries(entry.oid, path+"/"+entry.name)
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

// for testing purposes in the same directory, remove when done
var blacklist = []string{"gitik", ".git"}

func isIgnored(path string) bool {
	for _, item := range blacklist {
		if strings.Contains(path, item) {
			return true
		}
	}
	return path == constants.GitDir
}

func emptyDir(directory string) error {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}
	for _, f := range files {
		fullPath := filepath.Join(directory, f.Name())
		if isIgnored(f.Name()) {
			continue
		}
		if f.IsDir() {
			err := emptyDir(fullPath)
			if err != nil {
				return err
			}
			err = os.Remove(fullPath)
			if errors.Is(err, syscall.ENOTEMPTY) {
				continue
			} else if err != nil {
				return err
			}
		}
		if f.Mode().IsRegular() {
			err = os.Remove(fullPath)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func readObject(oid storage.OID, expectedType storage.ObjectType) ([]byte, error) {
	obj, err := storage.GetObject(oid)
	if err != nil {
		return nil, err
	}
	if obj.ObjType != expectedType {
		return nil, fmt.Errorf("unexpected type: want %s, got %s", expectedType, obj.ObjType)
	}
	return obj.Data, nil
}
