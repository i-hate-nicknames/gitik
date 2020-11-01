package data

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const gitDir = ".gitik"

// Init initializes a new repository
func Init() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	err = os.Mkdir(gitDir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(dir, gitDir)
}

// HashObject calculates sha1 sum of given data, and puts it
// in the git directory using the hash as the name
// Basically, it's a store mechanism for a content-based database
func HashObject(fileName string) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	return hashAndStore(data)
}

func hashAndStore(data []byte) error {
	hash := sha1.Sum(data)
	buf := bytes.NewBuffer(hash[:])
	oid := fmt.Sprintf("%x", buf)
	file, err := os.Create(filepath.Join(gitDir, oid))
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	return err
}
