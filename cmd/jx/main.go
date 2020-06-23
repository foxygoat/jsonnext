package main

import (
	"fmt"
	"os"

	"foxygo.at/jsonnext"
	"github.com/alecthomas/kong"
	"github.com/google/go-jsonnet"
)

var CLI struct {
	ImportPath []string `name:"jpath" sep:":" short:"J"`
	Filename   string   `arg:"" optional:"" help:"File to evaluate. stdin is used if omitted or \"-\""`
}

func main() {
	kong.Parse(&CLI)

	vm := jsonnet.MakeVM()
	i := newImporter(CLI.ImportPath)
	vm.Importer(i)

	out, err := run(vm, CLI.Filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(out)
}

func newImporter(importPath []string) *jsonnext.Importer {
	i := jsonnext.Importer{SearchPath: importPath}
	i.AppendSearchFromEnv("JXPATH")
	return &i
}

func run(vm *jsonnet.VM, filename string) (string, error) {
	if filename == "" || filename == "-" {
		filename = "/dev/stdin"
	}

	node, _, err := vm.ImportAST("", filename)
	if err != nil {
		return "", err
	}

	out, err := vm.Evaluate(node)
	if err != nil {
		return "", err
	}

	return out, nil
}
