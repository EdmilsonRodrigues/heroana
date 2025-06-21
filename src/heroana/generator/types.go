package generator

type ExportedObjects struct {
	ExportedConstants  []ExportedConstant
	ExportedVariables  []ExportedVariable
	ExportedTypes      []ExportedType
	ExportedStructs    []ExportedStruct
	ExportedInterfaces []ExportedInterface
	ExportedFunctions  []ExportedRoutine
}

type ExportedConstant struct {
	Name  string
	Type  string
	Value any
	Doc   string
}

type ExportedVariable struct {
	Name  string
	Type  string
	Value any
	Doc   string
}

type ExportedType struct {
	Name string
	Type string
	Doc  string
}

type ExportedField struct {
	Name string
	Type string
	Doc  string
}

type ExportedArgument struct {
	Name string
	Type string
}

type ExportedRoutine struct {
	Name        string
	Arguments   []ExportedArgument
	ReturnTypes []string
	Doc         string
}

type ExportedInterface struct {
	Name    string
	Methods []ExportedRoutine
	Doc     string
}

type ExportedStruct struct {
	Name    string
	Fields  []ExportedField
	Methods []ExportedRoutine
	Doc     string
}
