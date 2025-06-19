package cli_test

import (
	"reflect"
	"testing"

	"github.com/EdmilsonRodrigues/melo-project/src/melo/cli"
)

func TestParseArguments(t *testing.T) {
	t.Run("should parse option", func(t *testing.T) {
		allArguments := []struct {
			Name                  string
			Arguments             []string
			ExpectedCallsHelpNum  int
			ExpectedCallBuildArgs [][]string
		}{
			{
				Name: "Successfully build project",
				Arguments:             []string{cli.BuildFlag, ".", "--" + cli.OutputFlag, "output"},
				ExpectedCallsHelpNum:  0,
				ExpectedCallBuildArgs: [][]string{{".", "output"}},
			},
			{
				Name: "Successfully build project with default output path",
				Arguments:             []string{cli.BuildFlag, "."},
				ExpectedCallsHelpNum:  0,
				ExpectedCallBuildArgs: [][]string{{".", cli.DefaultOutputPath}},
			},
			{
				Name: "Successfully print help",
				Arguments:             []string{cli.HelpFlag},
				ExpectedCallsHelpNum:  1,
				ExpectedCallBuildArgs: nil,
			},
			{
				Name: "Print help when no arguments",
				Arguments:             []string{},
				ExpectedCallsHelpNum:  1,
				ExpectedCallBuildArgs: nil,
			},
			{
				Name: "Print help when passing build flag without arguments",
				Arguments:             []string{cli.BuildFlag},
				ExpectedCallsHelpNum:  1,
				ExpectedCallBuildArgs: nil,
			},
			{
				Name: "Print help when passing build flag without input path",
				Arguments:             []string{cli.BuildFlag, "--" + cli.OutputFlag, "output"},
				ExpectedCallsHelpNum:  1,
				ExpectedCallBuildArgs: nil,
			},
			{
				Name: "Print help when passing build flag without output path",
				Arguments:             []string{cli.BuildFlag, ".", "--" + cli.OutputFlag},
				ExpectedCallsHelpNum:  1,
				ExpectedCallBuildArgs: nil,
			},
			{
				Name: "Print help when passing build flag with output flag",
				Arguments:             []string{cli.BuildFlag, ".", "-o", "output"},
				ExpectedCallsHelpNum:  1,
				ExpectedCallBuildArgs: nil,
			},
			{
				Name: "Print help even if passing random arguments after help command",
				Arguments:             []string{cli.HelpFlag, "--" + cli.OutputFlag, "output"},
				ExpectedCallsHelpNum:  1,
				ExpectedCallBuildArgs: nil,
			},
			{
				Name: "Print help when no known command is sent",
				Arguments:             []string{"randomcommand"},
				ExpectedCallsHelpNum:  1,
				ExpectedCallBuildArgs: nil,
			},
		}
		assertion := func(arguments []string, expectedCallsHelpNum int, expectedCallBuildArgs [][]string) {
			t.Helper()
			helpCalls := 0
			buildCalls := [][]string{}

			cli.HelpFunction = spyHelpFunc(&helpCalls)
			cli.BuildFunction = spyBuildFunc(&buildCalls)

			cli.ParseArguments(arguments)

			if helpCalls != expectedCallsHelpNum {
				t.Errorf("Expected %d help calls, got %d", expectedCallsHelpNum, helpCalls)
			}

			if len(buildCalls) != len(expectedCallBuildArgs) {
				t.Errorf("Expected %d build calls, got %d", len(expectedCallBuildArgs), len(buildCalls))
			}

			if len(buildCalls) == 0 {
				return
			}

			if !reflect.DeepEqual(buildCalls, expectedCallBuildArgs) {
				t.Errorf("Expected %+v arguments on build calls, got %+v", expectedCallBuildArgs, buildCalls)
			}
		}
		for _, argument := range allArguments {
			t.Run(argument.Name, func(t *testing.T) {
				assertion(argument.Arguments, argument.ExpectedCallsHelpNum, argument.ExpectedCallBuildArgs)
			})
		}
	})
}

func spyHelpFunc(num *int) cli.HelpFunctionType {
	return func() {
		*num++
	}
}

func spyBuildFunc(callsArgs *[][]string) cli.BuildFunctionType {
	return func(inputPath string, outputPath string) {
		*callsArgs = append(*callsArgs, []string{inputPath, outputPath})
	}
}
