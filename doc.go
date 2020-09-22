// Package jsonnext makes it simpler to use jsonnet in a program.
//
// jsonnet is a data templating language that can be considered an extension of
// JSON that makes it composable.
//
//  import "github.com/google/go-jsonnet"
//
// jsonnext provides some go types and functions to make it easier to plug
// jsonnet into a program with minimal work and to extend some of the jsonnet
// abstractions to add capability.
//
// Importer
//
// Importer is an implementation of jsonnet.Importer that supports netpaths -
// paths that start with a double-slash - // - for network paths. Netpaths are
// retrieved via https using a URLFetcher. Other paths are read from the
// filesystem.
//
// This form of netpath makes it simpler to use in a PATH-style,
// colon-separated environment variable where the colon in a URL would need to
// be escaped.
//
// Config
//
// A type to encapsulate the configurable properties of a jsonnet VM and
// extensions provided by this package. It is agnotistic to the source of the
// configuration - it could come from the command line, config files or
// constucted in code by an application.
//
// Config uses the types VMVarMap and VMVar that provide a generic abstraction
// of the various variable types used on the jsonnet command line. These are
// the combination of:
//  • external variables (extVars) vs top-level args (TLAs)
//  • strings vs code
//  • literals vs files
//
// VMVarMap has a helper to parse string forms of entries with the option to
// take values from environment variables to support jsonnet command-line
// use cases.
//
// Flag
//
// Functions that wrap flag.Var to define the flags for the Config struct using
// types that implement flag.Value to suitably parse the CLI options.
//
// The flag names and use align with the standard jsonnet CLI options as much
// as possible, with minor variations due to how the Go flag package works,
// such as single-hyphen long flag prefixes.
//
// The simplest use of flags is to let it create the Config struct for you with
// flags bound to the fields of that instance:
//
//     cfg := jsonnext.ConfigFlags(flag.CommandLine)
//     flag.Parse()
//     cfg.ConfigureVM(vm)
//
// However the functions that create the individual flag types are exported so
// they can be used to define different flag names for an application.
package jsonnext
