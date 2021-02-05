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
	Tree    storage.OID
	Parent  storage.OID
	Message string
}

func (c Commit) Encode() []byte {
	var buf bytes.Buffer
	buf.WriteString(encodePair("tree", string(c.Tree)))
	if c.Parent != "" {
		buf.WriteString(encodePair("parent", string(c.Parent)))
	}
	buf.WriteString("\n" + c.Message + "\n")
	return buf.Bytes()
}

func encodePair(a, b string) string {
	return fmt.Sprintf("%s %s\n", a, b)
}

func MakeCommit(message string) (storage.OID, error) {
	oid, err := plumbing.WriteTree(".")
	if err != nil {
		return "", err
	}
	c := Commit{Tree: oid, Message: message}
	headOID, err := getHead()
	if err != nil && !errors.Is(err, errNoHead) {
		return "", err
	}
	if err == nil {
		c.Parent = headOID
	}
	commitOID, err := storage.HashObject(c.Encode(), constants.TypeCommit)
	if err != nil {
		return "", err
	}
	err = setHead(commitOID)
	if err != nil {
		return "", fmt.Errorf("make commit: cannot write commit to head: %w", err)
	}
	return commitOID, nil
}

func setHead(oid storage.OID) error {
	path := path.Join(constants.GitDir, constants.HeadName)
	return storage.WriteFile(path, []byte(oid))
}

var errNoHead = errors.New("head not found or empty")

func getHead() (storage.OID, error) {
	path := path.Join(constants.GitDir, constants.HeadName)
	file, err := os.Open(path)
	defer file.Close()
	if errors.Is(err, os.ErrNotExist) {
		return "", errNoHead
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
		return "", errNoHead
	}
	return result, nil
}
