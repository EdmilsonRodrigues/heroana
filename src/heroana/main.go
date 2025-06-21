package main

import (
	"os"

	"github.com/EdmilsonRodrigues/melo-project/src/melo/cmd"
)

func main()  {
	cmd.ParseArguments(os.Args[1:])
}
