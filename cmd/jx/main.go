package main

import (
	"fmt"
	"os"

	jxkong "foxygo.at/jsonnext/kong"
	jsonnet "github.com/google/go-jsonnet"
)

type config struct {
	jxkong.Config
	Filename string `arg:"" optional:"" help:"File to evaluate. stdin is used if omitted or \"-\""`
}

func main() {
	cli := parseCLI()
	vm := cli.Config.MakeVM("JXPATH")

	out, err := run(vm, cli.Filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(out)
}

func run(vm *jsonnet.VM, filename string) (string, error) {
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
