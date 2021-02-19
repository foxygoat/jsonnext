package kong_test

import (
	"fmt"
	"os"

	jxkong "foxygo.at/jsonnext/kong"
	"github.com/alecthomas/kong"
)

func Example() {
	// Define kong CLI struct embedding jxkong.Config, adding your own
	// application-specific flags and args.
	cli := struct {
		jxkong.Config
		Verbose  bool
		Filename string `arg:""`
	}{
		Config: *jxkong.NewConfig(), // foxygo.at/jsonnext/kong imported as jxkong
	}

	// Simulate command line arguments
	os.Args = []string{
		"prog",
		"--ext-str=extvar=hello", // string external var
		`--tla-code=tlavar=1+1`,  // code top-level arg
		"--verbose",              // application-specific flag
		"filename",               // application-specific arg
	}

	// Use kong to parse command line into our CLI struct
	kong.Parse(&cli)

	// Create and configure jsonnet VM with command line args, and JPATH env var
	vm := cli.Config.MakeVM("JPATH")

	// Evaluate jsonnet snippet that references defined vars
	code := "function(tlavar) std.repeat(std.extVar('extvar'), tlavar)"
	result, _ := vm.EvaluateAnonymousSnippet("<literal>", code)

	fmt.Println(cli.Filename)
	fmt.Printf("verbose: %v\n", cli.Verbose)
	fmt.Println(result)
	// Output:
	// filename
	// verbose: true
	// "hellohello"
}
