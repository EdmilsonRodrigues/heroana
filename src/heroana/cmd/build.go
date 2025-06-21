package cmd

import (
	"os"

	"github.com/EdmilsonRodrigues/melo-project/src/melo/files"
)


func Build(inputPath string, outputPath string) {
	if !files.CheckInputFolder(os.DirFS("."), inputPath) {
		os.Exit(1)
	}

	if err := files.CreateOutputFolder(outputPath); err != nil {
		os.Exit(1)
	}

	_, err := files.ScanModule(os.DirFS("."), inputPath, outputPath)
	if err != nil {
		os.Exit(1)
	}

}