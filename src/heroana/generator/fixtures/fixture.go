// melo:package

package fixtures


// Go doc for my constant
const MyConst = "hello"

// Go doc for my type
type MyType string

// Go doc for my variable
var MyVar = "world"

// Go doc for my struct
type MyStruct struct {
	// Go doc for my field
	Name string
}

// Go doc for my interface
type MyInterface interface {
	SayHello(name string) string
}

// Go doc for my function
func SumTwoNumbers(a, b int) int {
	return a + b
}

// Go doc for my second function
func CanSumTwoNumbers(a, b int) (sum int, err error) {
	return a + b, nil
}

// Go doc for my method
func (s MyStruct) CanSumTwoNumbers2(a, b int) (sum int, err error) {
	return a + b, nil
}

