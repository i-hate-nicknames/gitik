package storage

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const GitDir = ".gitik"

type OID string

// Init initializes a new repository
func Init() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	err = os.Mkdir(GitDir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(dir, GitDir)
}

var UnexpectedTypeErr = errors.New("unexpected object type")

var InvalidObjectErr = errors.New("invalid object format")

// GetObject retrieves an object stored by HashObject under its object ID (oid)
// This is the retrieve process of the data stored by HashObject
func GetObject(oid OID, expectedHeader []byte) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath.Join(GitDir, string(oid)))
	if err != nil {
		return nil, err
	}
	split := bytes.SplitN(data, []byte{0}, 2)
	if len(split) != 2 {
		return nil, InvalidObjectErr
	}
	if string(split[0]) != string(expectedHeader) {
		return nil, UnexpectedTypeErr
	}
	return split[1], nil
}

// HashObject calculates sha1 sum of given data, and puts it
// in the git directory using the hash as the name
// Basically, it's a store mechanism for a content-based database
func HashObject(data []byte, header []byte) (OID, error) {
	header = append(header, byte(0))
	data = append(header, data...)
	hash := sha1.Sum(data)
	buf := bytes.NewBuffer(hash[:])
	oid := fmt.Sprintf("%x", buf)
	err := WriteFile(filepath.Join(GitDir, oid), data)
	if err != nil {
		return "", err
	}
	return OID(oid), nil
}

func WriteFile(path string, data []byte) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()
	_, err = file.Write(data)
	return
}
