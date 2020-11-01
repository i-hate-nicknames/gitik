package data

import (
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
