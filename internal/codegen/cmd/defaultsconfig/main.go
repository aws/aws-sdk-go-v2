package main

import (
	"flag"
	"log"
)

func main() {
	jsonFile := flag.String("json", "", "The JSON defaults configuration file.")
	outputFile := flag.String("output", "defaults.go", "The output filename.")
	packageName := flag.String("p", "defaults", "The Go package name.")
	resolverName := flag.String("r", "GetModeConfiguration", "The configuration resolver function name.")
	flag.Parse()

	validateStringArgument(jsonFile, "-json file name is required")
	validateStringArgument(outputFile, "-output file name is required")
	validateStringArgument(packageName, "-p package name is required")
	validateStringArgument(resolverName, "-r resolver name is required")

	if err := generateConfigPackage(*jsonFile, *outputFile, *packageName, *resolverName); err != nil {
		log.Fatal(err)
	}
}
