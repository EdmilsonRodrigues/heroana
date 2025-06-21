package generator

import (
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"strings"

	"golang.org/x/tools/go/packages"
)

func InspectPackage(packagePath string) (exportedObjects ExportedObjects, err error) {
	pkg, err := getPackage(packagePath)
	if err != nil {
		err = fmt.Errorf("error inspecting package: %w", err)
		return
	}

	for _, file := range pkg.Syntax {
		inspectAbstractSyntaxTree(file, pkg, &exportedObjects)
	}

	return

	// // We're interested in the main package provided, assuming one is found.
	// // For packages with multiple Go files, they are combined into one *packages.Package object.
	// pkg := pkgs[0]

	// fmt.Println("\n--- Exported Functions ---")
	// if len(exportedFunctions) == 0 {
	// 	fmt.Println("(none)")
	// } else {
	// 	for _, f := range exportedFunctions {
	// 		fmt.Println(f)
	// 	}
	// }

	// fmt.Println("\n--- Exported Variables ---")
	// if len(exportedVars) == 0 {
	// 	fmt.Println("(none)")
	// } else {
	// 	for _, v := range exportedVars {
	// 		fmt.Println(v)
	// 	}
	// }

	// fmt.Println("\n--- Exported Constants ---")
	// if len(exportedConsts) == 0 {
	// 	fmt.Println("(none)")
	// } else {
	// 	for _, c := range exportedConsts {
	// 		fmt.Println(c)
	// 	}
	// }

	// fmt.Println("\n--- Exported Types ---")
	// if len(exportedTypes) == 0 {
	// 	fmt.Println("(none)")
	// } else {
	// 	for _, t := range exportedTypes {
	// 		fmt.Println(t)
	// 	}
	// }

	// return nil
}

func getPackage(packagePath string) (*packages.Package, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedDeps | packages.NeedExportFile |
			packages.NeedTypes | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedModule,
		Tests: false,
		Logf:  log.Printf,
	}

	pkgs, err := packages.Load(cfg, packagePath)
	if err != nil {
		return nil, fmt.Errorf("loading packages: %w", err)
	}

	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("packages contained errors")
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no packages found for %q", packagePath)
	}

	pkg := pkgs[0]
	log.Printf("\n\n--- Inspecting Package: %s (%s)  ---\n", pkg.Name, pkg.PkgPath)
	log.Printf("\n\n--- %+v ----\n\n", pkg.Types)

	return pkg, nil
}

func inspectAbstractSyntaxTree(file *ast.File, pkg *packages.Package, exportedObjects *ExportedObjects) {
	methods := make(map[string][]*ExportedRoutine)
	ast.Inspect(file, func(node ast.Node) bool {
		switch declaration := node.(type) {
		case *ast.FuncDecl: // Function or Method
			if !ast.IsExported(declaration.Name.Name) {
				return true
			}

			if exportedRoutine, receiver := parseRoutineDeclaration(pkg, declaration); receiver == "" {
				exportedObjects.ExportedFunctions = append(exportedObjects.ExportedFunctions, exportedRoutine)
			} else {
				methods[receiver] = append(methods[receiver], &exportedRoutine)
			}

		case *ast.GenDecl: // General declaration (const, var, type)
			for _, spec := range declaration.Specs {
				switch specification := spec.(type) {
				case *ast.ValueSpec: // Var or Const
					for index, name := range specification.Names {
						if !ast.IsExported(name.Name) {
							continue
						}

						variable := pkg.TypesInfo.Defs[name]
						if variable == nil {
							continue
						}
						typeName := strings.TrimLeft(variable.Type().String(), "untyped ")

						switch name.Obj.Kind {
						case ast.Con:
							exportedObjects.ExportedConstants = append(exportedObjects.ExportedConstants, ExportedConstant{
								Name:  name.Name,
								Type:  typeName,
								Value: parseVariableValue(specification.Values[index].(*ast.BasicLit), typeName),
								Doc:   specification.Doc.Text(), // TODO: Find a way to get the doc
							})
						case ast.Var:
							exportedObjects.ExportedVariables = append(exportedObjects.ExportedVariables, ExportedVariable{
								Name:  name.Name,
								Type:  typeName,
								Value: parseVariableValue(specification.Values[index].(*ast.BasicLit), typeName),
								Doc:   specification.Doc.Text(), // TODO: Find a way to get the doc
							})
						}
					}
				case *ast.TypeSpec: // Type declaration
					if !ast.IsExported(specification.Name.Name) {
						continue
					}
					object := pkg.TypesInfo.Defs[specification.Name]
					if object == nil {
						continue
					}

					switch underlying := object.Type().Underlying().(type) {
					case *types.Struct:
						exportedStruct := parseExportedStruct(underlying, specification)
						exportedObjects.ExportedStructs = append(exportedObjects.ExportedStructs, exportedStruct)
					case *types.Interface:
						exportedInterface := parseExportedInterface(underlying, specification)
						exportedObjects.ExportedInterfaces = append(exportedObjects.ExportedInterfaces, exportedInterface)
					default:
						exportedObjects.ExportedTypes = append(exportedObjects.ExportedTypes, ExportedType{
							Name: specification.Name.Name,
							Type: underlying.String(),
							Doc:  specification.Doc.Text(), // TODO: Find a way to get the doc
						})
					}

				}
			}
		}
		return true // Continue inspecting child nodes
	})

	for structName, structMethods := range methods {
		exportedStruct := findStructByName(exportedObjects.ExportedStructs, structName)
		if exportedStruct == nil {
			continue
		}
		for _, method := range structMethods {
			exportedStruct.Methods = append(exportedStruct.Methods, *method)
		}
	}

}

func parseRoutineDeclaration(pkg *packages.Package, declaration *ast.FuncDecl) (exportedRoutine ExportedRoutine, receiver string) {
	signature := pkg.TypesInfo.Defs[declaration.Name].Type().(*types.Signature)

	if signature.Recv() != nil { // Method
		splittedSignature := strings.Split(signature.Recv().Type().String(), ".")
		receiver = splittedSignature[len(splittedSignature)-1]
	}

	return ExportedRoutine{
		Name:        declaration.Name.Name,
		Arguments:   parseArguments(signature.Params()),
		ReturnTypes: parseReturnTypes(signature.Results()),
		Doc:         strings.TrimRight(declaration.Doc.Text(), "\n"),
	}, receiver

}

func parseArguments(arguments *types.Tuple) []ExportedArgument {
	exportedArguments := make([]ExportedArgument, 0, arguments.Len())
	for index := range arguments.Len() {
		argument := arguments.At(index)
		exportedArguments = append(exportedArguments, ExportedArgument{
			Name: argument.Name(),
			Type: argument.Type().String(),
		})
	}
	return exportedArguments
}

func parseReturnTypes(results *types.Tuple) []string {
	returnTypes := make([]string, 0, results.Len())
	for index := range results.Len() {
		result := results.At(index)
		returnTypes = append(returnTypes, result.Type().String())
	}
	return returnTypes
}

func parseVariableValue(value *ast.BasicLit, typeName string) string {
	if typeName == "string" {
		return strings.Trim(value.Value, "\"")
	}
	return value.Value
}

func parseExportedStruct(structType *types.Struct, specification *ast.TypeSpec) ExportedStruct {
	return ExportedStruct{
		Name:   specification.Name.Name,
		Fields: parseStructFields(structType),
		Doc:    strings.TrimRight(specification.Doc.Text(), "\n"),
	}
}

func parseStructFields(structType *types.Struct) []ExportedField {
	exportedFields := make([]ExportedField, 0, structType.NumFields())
	for field := range structType.Fields() {
		exportedFields = append(exportedFields, ExportedField{
			Name: field.Name(),
			Type: field.Type().String(),
			// TODO: Find a way to get the doc
		})
	}
	return exportedFields
}

func parseExportedInterface(interfaceType *types.Interface, specification *ast.TypeSpec) ExportedInterface {
	return ExportedInterface{
		Name:    specification.Name.Name,
		Methods: parseInterfaceMethods(interfaceType),
		Doc:     strings.TrimRight(specification.Doc.Text(), "\n"),
	}
}

func parseInterfaceMethods(interfaceType *types.Interface) []ExportedRoutine {
	exportedMethods := make([]ExportedRoutine, 0, interfaceType.NumExplicitMethods())
	for method := range interfaceType.ExplicitMethods() {
		exportedMethods = append(exportedMethods, parseInterfaceMethod(method))
	}
	return exportedMethods
}

func parseInterfaceMethod(method *types.Func) ExportedRoutine {
	return ExportedRoutine{
		Name:        method.Name(),
		Arguments:   parseArguments(method.Signature().Params()),
		ReturnTypes: parseReturnTypes(method.Signature().Results()),
	}
}

func findStructByName(structs []ExportedStruct, name string) *ExportedStruct {
	for _, exportedStruct := range structs {
		if exportedStruct.Name == name {
			return &exportedStruct
		}
	}
	return nil
}
