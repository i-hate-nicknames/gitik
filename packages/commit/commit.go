package commit

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/i-hate-nicknames/gitik/packages/constants"
	"github.com/i-hate-nicknames/gitik/packages/plumbing"
	"github.com/i-hate-nicknames/gitik/packages/storage"
)

type Commit struct {
	OID     storage.OID
	Tree    storage.OID
	Parent  storage.OID
	Message string
}

func CommitCurrentTree(message string) (storage.OID, error) {
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
	commitOID, err := storage.StoreObject(c.Encode(), constants.TypeCommit)
	if err != nil {
		return "", err
	}
	err = setHead(commitOID)
	if err != nil {
		return "", fmt.Errorf("make commit: cannot write commit to head: %w", err)
	}
	return commitOID, nil
}

func Log() ([]Commit, error) {
	head, err := getHead()
	if err != nil {
		return nil, err
	}
	return LogFrom(head)
}

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

func (c Commit) Encode() []byte {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("tree %s\n", string(c.Tree)))
	if c.Parent != "" {
		buf.WriteString(fmt.Sprintf("parent %s\n", string(c.Parent)))
	}
	buf.WriteString("\n" + c.Message + "\n")
	return buf.Bytes()
}

var ErrInvalidEncoding = errors.New("invalid encoding")

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
