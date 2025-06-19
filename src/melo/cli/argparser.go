package cli

import (
	"fmt"
	"log"
	"os"
)

const (
	DefaultOutputPath = "build"
	BuildFlag         = "build"
	HelpFlag          = "help"

	OutputFlag = "output"
)

type (
	HelpFunctionType  func()
	BuildFunctionType func(inputPath string, outputPath string)
)

var (
	HelpFunction  HelpFunctionType
	BuildFunction BuildFunctionType
)

func ParseArguments(arguments []string) {
	if len(arguments) < 1 {
		fmt.Println("Error: Missing command")
		HelpFunction()
		return
	}

	command, arguments := arguments[0], arguments[1:]

	switch command {
	case BuildFlag:
		usage := func() {
			fmt.Fprintf(os.Stderr, "Usage: %s %s <inputPath> [-o <outputPath>]\n", os.Args[0], BuildFlag)
			fmt.Fprintf(os.Stderr, "Options for '%s' command:\n", BuildFlag)
			fmt.Fprintf(os.Stderr, "  --%s <outputPath>	Output folder path\n", OutputFlag)
		}

		if len(arguments) < 1 {
			fmt.Println("Error: Missing input file path")
			usage()
			return
		}

		inputPath, arguments := arguments[0], arguments[1:]
		outputPath := parseBuildArguments(arguments, usage)

		log.Printf("Building %s to %s\n", inputPath, outputPath)
		BuildFunction(inputPath, outputPath)

	case HelpFlag:
		HelpFunction()

	default:
		fmt.Println("Error: Unknown command")
		HelpFunction()
	}
}

func parseBuildArguments(arguments []string, usage HelpFunctionType) (outputPath string) {
	outputPath = DefaultOutputPath
	for len(arguments) > 0 {
		usedArguments := 0
		switch arguments[0] {
		case fmt.Sprintf("--%s", OutputFlag):
			if len(arguments) < 2 {
				fmt.Println("Error: Missing output folder path")
				usage()
				return
			}
			outputPath = arguments[1]
			usedArguments = 2
		default:
			fmt.Println("Error: Unknown option")
			usage()
			return
		}
		arguments = arguments[usedArguments:]
	}
	return
}
