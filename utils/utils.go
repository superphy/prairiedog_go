package utils

import (
	"log"
	"os"
	"path/filepath"
)

func Walk(root string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return files, err
}
