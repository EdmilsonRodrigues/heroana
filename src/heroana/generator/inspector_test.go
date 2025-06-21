package generator_test

import (
	"reflect"
	"testing"

	"github.com/EdmilsonRodrigues/melo-project/src/melo/generator"
)

func TestInspectPackage(t *testing.T) {
	fixturePath := "github.com/EdmilsonRodrigues/melo-project/src/melo/generator/fixtures"

	expectedContents := generator.ExportedObjects{
		ExportedConstants: []generator.ExportedConstant{
			{
				Name:  "MyConst",
				Type:  "string",
				Value: "hello", 
				// Doc:   "Go doc for my constant",  // TODO Uncomment after finding a way to get the doc
			},
		},
		ExportedVariables: []generator.ExportedVariable{
			{
				Name:  "MyVar",
				Type:  "string",
				Value: "world",
				// Doc:   "Go doc for my variable",  // TODO Uncomment after finding a way to get the doc
			},
		},
		ExportedTypes: []generator.ExportedType{
			{
				Name: "MyType",
				Type: "string",
				// Doc:  "Go doc for my type",  // TODO Uncomment after finding a way to get the doc
			},
		},
		ExportedStructs: []generator.ExportedStruct{
			{
				Name: "MyStruct",
				Fields: []generator.ExportedField{
					{
						Name: "Name",
						Type: "string",
						// Doc:  "Go doc for my field",  // TODO Uncomment after finding a way to get the doc
					},
				},
				Methods: []generator.ExportedRoutine{
					{
						Name: "CanSumTwoNumbers2",
						Arguments: []generator.ExportedArgument{
							{
								Name: "a",
								Type: "int",
							},
							{
								Name: "b",
								Type: "int",
							},
						},
						ReturnTypes: []string{"int", "error"},
						Doc:         "Go doc for my method",
					},
				},
				Doc: "Go doc for my struct",
			},
		},
		ExportedInterfaces: []generator.ExportedInterface{
			{
				Name: "MyInterface",
				Methods: []generator.ExportedRoutine{
					{
						Name: "SayHello",
						Arguments: []generator.ExportedArgument{
							{
								Name: "name",
								Type: "string",
							},
						},
						ReturnTypes: []string{"string"},
						Doc:         "Go doc for my method",
					},
				},
				Doc: "Go doc for my interface",
			},
		},
		ExportedFunctions: []generator.ExportedRoutine{
			{
				Name: "SumTwoNumbers",
				Arguments: []generator.ExportedArgument{
					{
						Name: "a",
						Type: "int",
					},
					{
						Name: "b",
						Type: "int",
					},
				},
				ReturnTypes: []string{"int"},
				Doc:         "Go doc for my function",
			},
			{
				Name: "CanSumTwoNumbers",
				Arguments: []generator.ExportedArgument{
					{
						Name: "a",
						Type: "int",
					},
					{
						Name: "b",
						Type: "int",
					},
				},
				ReturnTypes: []string{"int", "error"},
				Doc:         "Go doc for my second function",
			},
		},
	}

	t.Run("should inspect fixture package", func(t *testing.T) {
		inspectedContents, err := generator.InspectPackage(fixturePath)
		if err != nil {
			t.Errorf("InspectPackage should not return error, got %v", err)
		}

		if !reflect.DeepEqual(inspectedContents, expectedContents) {
			t.Errorf("InspectPackage should return %v, got %v", expectedContents, inspectedContents)
		}
	})
}
