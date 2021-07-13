// jnxflag evaluates a jsonnet file and outputs it as JSON.
//
// Usage of ./jnxflag:
//   -A var[=str]
//         Add top-level arg var[=str] (from environment if <str> is omitted)
//   -J dir
//         Add a library search dir
//   -V var[=str]
//         Add extVar var[=str] (from environment if <str> is omitted)
//   -ext-code var[=code]
//         Add extVar var[=code] (from environment if <code> is omitted)
//   -ext-code-file var=file
//         Add extVar var=file code from a file
//   -ext-str var[=str]
//         Add extVar var[=str] (from environment if <str> is omitted)
//   -ext-str-file var=file
//         Add extVar var=file string from a file
//   -jpath dir
//         Add a library search dir
//   -tla-code var[=code]
//         Add top-level arg var[=code] (from environment if <code> is omitted)
//   -tla-code-file var=file
//         Add top-level arg var=file code from a file
//   -tla-str var=[=str]
//         Add top-level arg var=[=str] (from environment if <str> is omitted)
//   -tla-str-file var=file
//         Add top-level arg var=file string from a file
//
// This program exists just to implement the standard Go flag package parsing.
// The full jnx program uses the kong library and has more features.

package main

import (
	"flag"
	"fmt"
	"os"

	"foxygo.at/jsonnext"
	jsonnet "github.com/google/go-jsonnet"
)

type config struct {
	jsonnext.Config
	Filename string `arg:"" optional:"" help:"File to evaluate. stdin is used if omitted or \"-\""`
}

func main() {
	cli := parseCLI()
	vm := cli.Config.MakeVM("JNXPATH")

	out, err := run(vm, cli.Filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Print(out)
}

// Parse CLI using Go's flag package and the helpers in jsonnext.
func parseCLI() *config {
	c := &config{}
	c.Config = *jsonnext.ConfigFlags(flag.CommandLine)

	flag.Parse()
	if flag.NArg() > 1 {
		flag.Usage()
		os.Exit(1)
	} else if flag.NArg() == 1 {
		c.Filename = flag.Args()[0]
	}

	return c
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
