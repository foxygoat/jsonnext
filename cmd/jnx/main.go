package main

import (
	"fmt"
	"os"

	jnxkong "foxygo.at/jsonnext/kong"
	"github.com/alecthomas/kong"
	jsonnet "github.com/google/go-jsonnet"
)

type config struct {
	jnxkong.Config
	Filename string `arg:"" optional:"" help:"File to evaluate. stdin is used if omitted or \"-\""`
}

func main() {
	c := &config{Config: *jnxkong.NewConfig()}
	kong.Parse(c)
	vm := c.Config.MakeVM("JNXPATH")

	out, err := run(vm, c.Filename)
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
