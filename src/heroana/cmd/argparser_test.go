package cmd_test

import (
	"reflect"
	"testing"

	"github.com/EdmilsonRodrigues/melo-project/src/melo/cmd"
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
				Arguments:             []string{cmd.BuildFlag, ".", "--" + cmd.OutputFlag, "output"},
				ExpectedCallsHelpNum:  0,
				ExpectedCallBuildArgs: [][]string{{".", "output"}},
			},
			{
				Name: "Successfully build project with default output path",
				Arguments:             []string{cmd.BuildFlag, "."},
				ExpectedCallsHelpNum:  0,
				ExpectedCallBuildArgs: [][]string{{".", cmd.DefaultOutputPath}},
			},
			{
				Name: "Successfully print help",
				Arguments:             []string{cmd.HelpFlag},
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
				Arguments:             []string{cmd.BuildFlag},
				ExpectedCallsHelpNum:  1,
				ExpectedCallBuildArgs: nil,
			},
			{
				Name: "Print help when passing build flag without input path",
				Arguments:             []string{cmd.BuildFlag, "--" + cmd.OutputFlag, "output"},
				ExpectedCallsHelpNum:  1,
				ExpectedCallBuildArgs: nil,
			},
			{
				Name: "Print help when passing build flag without output path",
				Arguments:             []string{cmd.BuildFlag, ".", "--" + cmd.OutputFlag},
				ExpectedCallsHelpNum:  1,
				ExpectedCallBuildArgs: nil,
			},
			{
				Name: "Print help when passing build flag with output flag",
				Arguments:             []string{cmd.BuildFlag, ".", "-o", "output"},
				ExpectedCallsHelpNum:  1,
				ExpectedCallBuildArgs: nil,
			},
			{
				Name: "Print help even if passing random arguments after help command",
				Arguments:             []string{cmd.HelpFlag, "--" + cmd.OutputFlag, "output"},
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

			cmd.HelpFunction = spyHelpFunc(&helpCalls)
			cmd.BuildFunction = spyBuildFunc(&buildCalls)

			cmd.ParseArguments(arguments)

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

func spyHelpFunc(num *int) cmd.HelpFunctionType {
	return func() {
		*num++
	}
}

func spyBuildFunc(callsArgs *[][]string) cmd.BuildFunctionType {
	return func(inputPath string, outputPath string) {
		*callsArgs = append(*callsArgs, []string{inputPath, outputPath})
	}
}
