package main

import (
	"flag"
	"os"
)

func main() {
	flag.Parse()
	values := flag.Args()
	if values == nil {
		os.Exit(111)
	}
	dir := values[0]
	envData, err := ReadDir(dir)
	if err != nil {
		os.Exit(111)
	}
	returnCode := RunCmd(values[1:], envData)
	os.Exit(returnCode)
}
