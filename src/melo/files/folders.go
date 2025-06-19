package files

import (
	"io/fs"
	"log"
	"os"
	"strings"
)

const (
	goModStart = "module "
)

func CheckInputFolder(fileSystem fs.FS, path string) bool {
	log.Println("Checking input folder...", path)
	content, err := fs.ReadFile(fileSystem, strings.TrimRight(path, "/") + "/go.mod")
	if err != nil {
		log.Println("Error:", err)
		return false
	}

	if !strings.HasPrefix(string(content), goModStart) {
		log.Println("Input folder doesn't have a valid go.mod")
		return false
	}

	return true
}

func CreateOutputFolder(path string) error {
	log.Println("Creating output folder...", path)
	return os.Mkdir(path, fs.ModePerm)
}