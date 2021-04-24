package commit

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/i-hate-nicknames/gitik/pkg/constants"
	"github.com/i-hate-nicknames/gitik/pkg/plumbing"
	"github.com/i-hate-nicknames/gitik/pkg/storage"
)

// Commit represents a version control commit. It's a snapshot
// of repository together with a message and link to previous commit
type Commit struct {
	OID     storage.OID
	Tree    storage.OID
	Parent  storage.OID
	Message string
}

// SaveCurrentTree saves current working tree to the datastore, and creates a
// commit object that points to that tree. Additionally, it advances HEAD of
// the repository and point it to the fresly created commit
// Return new commit's storage ID
func SaveCurrentTree(message string) (storage.OID, error) {
	oid, err := plumbing.WriteTree(".")
	if err != nil {
		return "", err
	}
	c := Commit{Tree: oid, Message: message}
	headOID, err := getHead()
	if err != nil && !errors.Is(err, ErrNoHead) {
		return "", err
	}
	if err == nil {
		c.Parent = headOID
	}
	commitOID, err := storage.StoreObject(c.Encode(), storage.TypeCommit)
	if err != nil {
		return "", err
	}
	err = setHead(commitOID)
	if err != nil {
		return "", fmt.Errorf("make commit: cannot write commit to head: %w", err)
	}
	return commitOID, nil
}

// Log returns all commits that were made starting from HEAD
// and until the first commit, following the parent chain
func Log() ([]Commit, error) {
	head, err := getHead()
	if err != nil {
		return nil, err
	}
	return LogFrom(head)
}

// LogFrom returns all commits that were made starting from given commit
// and until the first commit, following the parent chain
func LogFrom(startFrom storage.OID) ([]Commit, error) {
	var log []Commit
	for currentOID := storage.OID(startFrom); currentOID != ""; {
		commit, err := getCommit(currentOID)
		if err != nil {
			return nil, err
		}
		log = append(log, commit)
		currentOID = commit.Parent
	}
	return log, nil
}

// Encode commit to byte sequence. This data can be later be used with
// Decode method to retrieve commit back
func (c Commit) Encode() []byte {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("tree %s\n", string(c.Tree)))
	if c.Parent != "" {
		buf.WriteString(fmt.Sprintf("parent %s\n", string(c.Parent)))
	}
	buf.WriteString("\n" + c.Message + "\n")
	return buf.Bytes()
}

// ErrInvalidEncoding signifies problems with commit encoding
var ErrInvalidEncoding = errors.New("invalid encoding")

// Decode data into a Commit
func Decode(data []byte) (Commit, error) {
	rawParts := bytes.Split(data, []byte("\n\n"))
	if len(rawParts) != 2 {
		return Commit{}, ErrInvalidEncoding
	}
	header, message := rawParts[0], rawParts[1]
	headerParts := bytes.Split(header, []byte("\n"))
	result := Commit{Message: string(message)}
	for _, line := range headerParts {
		parts := bytes.Split(line, []byte(" "))
		if len(parts) != 2 {
			return Commit{}, ErrInvalidEncoding
		}
		switch string(parts[0]) {
		case "tree":
			result.Tree = storage.OID(parts[1])
		case "parent":
			result.Parent = storage.OID(parts[1])
		default:
			return Commit{}, ErrInvalidEncoding
		}
	}
	return result, nil
}

func getCommit(oid storage.OID) (Commit, error) {
	obj, err := storage.GetObject(oid)
	if err != nil {
		return Commit{}, err
	}
	commit, err := Decode(obj.Data)
	if err != nil {
		return Commit{}, err
	}
	commit.OID = oid
	return commit, nil
}

func setHead(oid storage.OID) error {
	path := path.Join(constants.GitDir, constants.HeadName)
	return storage.WriteFile(path, []byte(oid))
}

// ErrNoHead is returned when repository has no HEAD
var ErrNoHead = errors.New("head not found or empty")

func getHead() (storage.OID, error) {
	path := path.Join(constants.GitDir, constants.HeadName)
	file, err := os.Open(path)
	defer file.Close()
	if errors.Is(err, os.ErrNotExist) {
		return "", ErrNoHead
	}
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	result := storage.OID(data)
	if result == "" {
		return "", ErrNoHead
	}
	return result, nil
}
