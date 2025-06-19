package main

import (
	"os"

	"github.com/EdmilsonRodrigues/melo-project/src/melo/cli"
)


func main()  {
	cli.ParseArguments(os.Args[1:])
}
