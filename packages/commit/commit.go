package commit

import (
	"fmt"

	"github.com/i-hate-nicknames/gitik/packages/constants"
	"github.com/i-hate-nicknames/gitik/packages/plumbing"
	"github.com/i-hate-nicknames/gitik/packages/storage"
)

type Commit struct {
	Tree    storage.OID
	Message string
}

func (c Commit) Encode() []byte {
	return []byte(fmt.Sprintf("tree %s\n\n%s\n", c.Tree, c.Message))
}

func MakeCommit(message string) (storage.OID, error) {
	oid, err := plumbing.WriteTree(".")
	if err != nil {
		return "", err
	}
	c := Commit{Tree: oid, Message: message}

	commitOID, err := storage.HashObject(c.Encode(), constants.TypeCommit)
	if err != nil {
		return "", err
	}
	return commitOID, nil
}
