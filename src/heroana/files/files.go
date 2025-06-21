package files

import (
	"fmt"
	"io/fs"
	"log"
	"strings"
)

type ExportedPackage struct {
	GoPath      string
	PythonPath  string
	PackageName string
}

const (
	GoExportedDirective = "// melo:"
)

func ScanModule(fileSystem fs.FS, path, moduleName string) ([]ExportedPackage, error) {
	exportedPackages := []ExportedPackage{}
	err := fs.WalkDir(fileSystem, path, getScanModuleWalker(fileSystem, &exportedPackages, moduleName))

	if err != nil {
		log.Printf("Error scanning module: %v", err)
		return nil, err
	}

	return exportedPackages, nil
}

func getScanModuleWalker(fileSystem fs.FS, exportedPackages *[]ExportedPackage, moduleName string) fs.WalkDirFunc {
	return func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error reading %s: %w", path, err)
		}

		if dirEntry.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		packageName, pythonPath, exported, err := parseGoFile(fileSystem, path)
		if err != nil {
			return fmt.Errorf("error parsing %s: %w", path, err)
		}

		if !exported {
			return nil
		}

		exportedPackage := genExportedPackage(path, moduleName, packageName, pythonPath)

		*exportedPackages = append(*exportedPackages, exportedPackage)
		return nil
	}
}

func goFormatPath(path, moduleName string) string {
	splittedPath := strings.Split(path, "/")
	splittedPath = splittedPath[1 : len(splittedPath)-1]
	log.Println(splittedPath)

	return strings.Join(append([]string{moduleName}, splittedPath...), "/")
}

func parseGoFile(fileSystem fs.FS, filePath string) (packageName string, pythonPath string, exported bool, err error) {
	content, err := fs.ReadFile(fileSystem, filePath)
	if err != nil {
		return
	}

	for _, line := range strings.Split(string(content), "\n") {
		if strings.HasPrefix(line, GoExportedDirective) {
			pythonPath = strings.TrimPrefix(line, GoExportedDirective)
			exported = true
		} else if strings.HasPrefix(line, "package ") {
			packageName = strings.TrimPrefix(line, "package ")
			break
		}
	}

	if pythonPath != "" {
		exported = true
	}
	return
}

func genExportedPackage(path, moduleName, packageName, pythonPath string) ExportedPackage {
	exportedPackege := ExportedPackage{
		GoPath:     goFormatPath(path, moduleName),
		PythonPath: pythonPath,
	}

	if packageName == "main" {
		return exportedPackege
	}

	splittedGoPath := strings.Split(exportedPackege.GoPath, "/")

	if packageName != splittedGoPath[len(splittedGoPath)-1] {
		exportedPackege.PackageName = packageName
	}

	return exportedPackege
}
