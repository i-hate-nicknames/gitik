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
		return storage.ZeroOID, err
	}
	c := Commit{Tree: oid, Message: message}
	headOID, err := getHeadOID()
	if err != nil && !errors.Is(err, ErrNoHead) {
		return storage.ZeroOID, err
	}
	if err == nil {
		c.Parent = headOID
	}
	commitOID, err := storage.StoreObject(c.Encode(), storage.TypeCommit)
	if err != nil {
		return storage.ZeroOID, err
	}
	err = SetHead(commitOID)
	if err != nil {
		return storage.ZeroOID, fmt.Errorf("make commit: cannot write commit to head: %w", err)
	}
	return commitOID, nil
}

// Log returns all commits that were made starting from HEAD
// and until the first commit, following the parent chain
func Log() ([]Commit, error) {
	head, err := getHeadOID()
	if err != nil {
		return nil, err
	}
	return LogFrom(head)
}

// LogFrom returns all commits that were made starting from given commit
// and until the first commit, following the parent chain
func LogFrom(startFrom storage.OID) ([]Commit, error) {
	var log []Commit
	for currentOID := storage.OID(startFrom); currentOID != storage.ZeroOID; {
		commit, err := GetCommit(currentOID)
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
	buf.WriteString(fmt.Sprintf("tree %s\n", string(c.Tree[:])))
	if c.Parent != storage.ZeroOID {
		buf.WriteString(fmt.Sprintf("parent %s\n", string(c.Parent[:])))
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
		oid, err := storage.MakeOID(parts[1])
		if err != nil {
			return Commit{}, err
		}
		switch string(parts[0]) {
		case "tree":
			result.Tree = oid
		case "parent":
			result.Parent = oid
		default:
			return Commit{}, ErrInvalidEncoding
		}
	}
	return result, nil
}

// GetCommit gets commit by its ID
func GetCommit(oid storage.OID) (Commit, error) {
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

// SetHead sets current HEAD of gitik to give oid
func SetHead(oid storage.OID) error {
	path := path.Join(constants.GitDir, constants.HeadName)
	return storage.WriteFile(path, []byte(oid.String()))
}

// ErrNoHead is returned when repository has no HEAD
var ErrNoHead = errors.New("head not found or empty")

// GetHead returns current HEAD (i.e. currently checked out tree)
func GetHead() (Commit, error) {
	OID, err := getHeadOID()
	if err != nil {
		return Commit{}, err
	}
	return GetCommit(OID)
}

func getHeadOID() (storage.OID, error) {
	path := path.Join(constants.GitDir, constants.HeadName)
	file, err := os.Open(path)
	defer file.Close()
	if errors.Is(err, os.ErrNotExist) {
		return storage.ZeroOID, ErrNoHead
	}
	if err != nil {
		return storage.ZeroOID, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return storage.ZeroOID, err
	}
	if len(data) == 0 {
		return storage.ZeroOID, ErrNoHead
	}
	return storage.MakeOID(data)
}

type CheckoutError struct {
	origError    error
	recoverError error
}

func (ce CheckoutError) Error() string {
	if ce.origError != nil {
		msg := fmt.Sprintf("failed to read tree: %s", ce.origError)
		if ce.recoverError != nil {
			msg = fmt.Sprintf("failed to recover: %s, original error: %s", ce.recoverError, msg)
		}
		return msg
	}
	return ""
}

func (c Commit) Checkout(recover bool) error {
	head, err := GetHead()
	if err != nil {
		return err
	}
	var finalError CheckoutError
	err = plumbing.ReadTree(c.Tree)
	if err != nil {
		if !recover {
			return err
		}
		finalError.origError = err
		recoverErr := head.Checkout(false)
		if recoverErr != nil {
			finalError.recoverError = recoverErr
		}
		return finalError
	}
	return SetHead(c.OID)
}
