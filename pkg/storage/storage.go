package storage

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/i-hate-nicknames/gitik/pkg/constants"
)

// OID is object id, a sha1 checksum of object contents, used to identify object for later retrieval
type OID [sha1.Size]byte

// ZeroOID is OID of no object, used as empty value
var ZeroOID OID

// MakeOID makes Object ID from string encoding
func MakeOID(hexStr []byte) (OID, error) {
	var oid OID
	// this is twice as much length as we need,
	// because hexStr in string encoding takes twice as much space as raw bytes
	// but let it be this way for readability
	decoded := make([]byte, len(hexStr))
	n, err := hex.Decode(decoded, hexStr)
	if err != nil {
		return ZeroOID, fmt.Errorf("makeOID: %w", err)
	}
	if n != sha1.Size {
		return ZeroOID, fmt.Errorf("makeOID: invalid length (%d), expected %d, oid: %s", n, sha1.Size, hexStr)
	}
	copy(oid[:], decoded)
	return oid, nil
}

func (oid OID) String() string {
	return fmt.Sprintf("%x", oid[:])
}

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
	data, err := ioutil.ReadFile(getObjectPath(oid))
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
	oid := OID(hash)
	err := WriteFile(getObjectPath(oid), data)
	if err != nil {
		return ZeroOID, err
	}
	return oid, nil
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

// get file path for given object id
func getObjectPath(oid OID) string {
	// in future we might want to use first couple of bytes as a directory
	// like git does
	return filepath.Join(constants.GitDir, oid.String())
}
