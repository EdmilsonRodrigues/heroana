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
	HelpFunction  HelpFunctionType = Help
	BuildFunction BuildFunctionType
)

func ParseArguments(arguments []string) {
	if len(arguments) < 1 {
		fmt.Fprintf(os.Stderr, "Error: Missing command.\n\n")
		HelpFunction()
		return
	}

	command, arguments := arguments[0], arguments[1:]

	switch command {
	case BuildFlag:
		usage := func() {
			fmt.Fprintf(os.Stderr, "Usage: %s %s <inputPath> [--%s <outputPath>]\n", os.Args[0], BuildFlag, OutputFlag)
			fmt.Fprintf(os.Stderr, "Options for '%s' command:\n", BuildFlag)
			fmt.Fprintf(os.Stderr, "  --%s <outputPath>	Output folder path\n", OutputFlag)
		}

		if len(arguments) < 1 {
			fmt.Fprint(os.Stderr, "Error: Missing input file path\n\n")
			usage()
			return
		}

		inputPath, arguments := arguments[0], arguments[1:]
		outputPath, err := parseBuildArguments(arguments)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n\n", err.Error())
			usage()
			return
		}

		log.Printf("Building %s to %s\n", inputPath, outputPath)
		BuildFunction(inputPath, outputPath)

	case HelpFlag:
		HelpFunction()

	default:
		fmt.Fprint(os.Stderr, "Error: Unknown command\n\n")
		HelpFunction()
	}
}

func parseBuildArguments(arguments []string) (outputPath string, err error) {
	outputPath = DefaultOutputPath
	for len(arguments) > 0 {
		usedArguments := 0
		switch arguments[0] {
		case fmt.Sprintf("--%s", OutputFlag):
			if len(arguments) < 2 {
				err = fmt.Errorf("error: Missing output folder path")
				return
			}
			outputPath = arguments[1]
			usedArguments = 2
		default:
			err = fmt.Errorf("error: Unknown option")
			return
		}
		arguments = arguments[usedArguments:]
	}
	return
}
