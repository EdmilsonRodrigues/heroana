package cli_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/EdmilsonRodrigues/melo-project/src/melo/cli"
)

func TestParseArguments(t *testing.T) {
	t.Run("should parse option", func(t *testing.T) {
		allArguments := []struct {
			Arguments             []string
			ExpectedCallsHelpNum  int
			ExpectedCallBuildArgs [][]string
		}{
			{[]string{cli.BuildFlag, ".", "--" + cli.OutputFlag, "output"}, 0, [][]string{{".", "output"}}},
			{[]string{cli.BuildFlag, "."}, 0, [][]string{{".", cli.DefaultOutputPath}}},
			{[]string{cli.HelpFlag}, 1, nil},
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
			t.Run("parsing arguments "+strings.Join(argument.Arguments, " "), func(t *testing.T) {
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
