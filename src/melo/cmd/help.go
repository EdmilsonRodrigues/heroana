package cmd

import "fmt"

const ProjectName = "Melo"

func Help() {
	fmt.Printf("%s is a project created to make a bridge between Go and Python.\n", ProjectName)
	fmt.Printf("By running the %s command, you can build your project.\n", BuildFlag)
	fmt.Print("This will generate a python module ready to be exported, where your go code will be run transparently.\n\n")
	fmt.Println("Commands:")
	fmt.Printf("  %s <inputPath> [--%s <outputPath>] \tBuild your project\n", BuildFlag, OutputFlag)
	fmt.Printf("  %s \t\t\t\t\t\tPrints this message\n\n", HelpFlag)
}