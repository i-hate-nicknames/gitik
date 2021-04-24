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

	"github.com/i-hate-nicknames/gitik/pkg/constants"
)

// OID is object id, a sha1 checksum of object contents, used to identify object for later retrieval
type OID string

// StoredObject represents an object that is retrieved from the storage
type StoredObject struct {
	ObjType ObjectType
	Data    []byte
}

// Init initializes a new repository
func Init() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	err = os.Mkdir(constants.GitDir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(dir, constants.GitDir)
}

// ErrInvalidObject is returned when object format is invalid
var ErrInvalidObject = errors.New("invalid object format")

// GetObject retrieves an object stored by HashObject under its object ID (oid)
// This is the retrieve process of the data stored by HashObject
func GetObject(oid OID) (StoredObject, error) {
	var obj StoredObject
	data, err := ioutil.ReadFile(filepath.Join(constants.GitDir, string(oid)))
	if err != nil {
		return obj, err
	}
	split := bytes.SplitN(data, []byte{0}, 2)
	if len(split) != 2 {
		return obj, ErrInvalidObject
	}
	objType, err := Decode(split[0])
	if err != nil {
		return obj, err
	}
	obj.Data = split[1]
	obj.ObjType = objType
	return obj, nil
}

// StoreObject calculates sha1 sum of given data, and puts it
// in the git directory using the hash as the name
// Basically, it's a store mechanism for a content-based database
func StoreObject(data []byte, objType ObjectType) (OID, error) {
	header := objType.Encode()
	header = append(header, byte(0))
	data = append(header, data...)
	hash := sha1.Sum(data)
	buf := bytes.NewBuffer(hash[:])
	oid := fmt.Sprintf("%x", buf)
	err := WriteFile(filepath.Join(constants.GitDir, oid), data)
	if err != nil {
		return "", err
	}
	return OID(oid), nil
}

// WriteFile writes data to a regular file under given path
// return error on any i/o error, or if a file with this name already exists
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
