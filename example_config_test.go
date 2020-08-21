package jsonnext_test

import (
	"fmt"
	"os"

	jsonnet "github.com/google/go-jsonnet"

	"foxygo.at/jsonnext"
)

func ExampleConfig_ConfigureVM() {
	vm := jsonnet.MakeVM()

	cfg := jsonnext.NewConfig()
	cfg.ExtVars["extvar"] = jsonnext.NewExtStr("hello")
	cfg.TLAVars["tlavar"] = jsonnext.NewTLACode(`" world"`)
	cfg.ConfigureVM(vm)

	// "<literal>" is the filename for error reporting.
	str, _ := vm.EvaluateSnippet("<literal>", "function(tlavar) std.extVar('extvar') + tlavar")
	fmt.Println(str)
	// Output: "hello world"
}

func ExampleConfig_ConfigureImporter() {
	vm := jsonnet.MakeVM()
	i := &jsonnext.Importer{}
	vm.Importer(i)

	// testdata/importer/hello.txt contains: hello world\n
	// testdata/importer/mellow.txt contains: mellow world\n
	_ = os.Setenv("EXAMPLE_PATH", "testdata/importer")
	cfg := jsonnext.NewConfig()
	cfg.ImportPath = []string{"testdata"}
	cfg.ConfigureImporter(i, "EXAMPLE_PATH")

	// "<literal>" is the filename for error reporting.
	str, _ := vm.EvaluateSnippet("<literal>", `[importstr "hello.txt", importstr "importer/mellow.txt"]`)
	fmt.Println(str)
	// Output:
	// [
	//    "hello world\n",
	//    "mellow world\n"
	// ]
}
