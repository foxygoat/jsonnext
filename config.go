package jsonnext

import (
	"errors"
	"os"
	"strings"

	"foxygo.at/s/errs"
	jsonnet "github.com/google/go-jsonnet"
)

const (
	maxStackDepth       = 500 // default from jsonnet
	maxStackTraceOutput = 20  // default from jsonnet
)

// Sentinel errors returned by Config functions. Callers can use errors.Is
// with these sentinels to handle the specific types of errors.
var (
	ErrMissingKey   = errors.New("missing key")
	ErrMissingValue = errors.New("missing value")
)

// Config holds configuration for a jsonnet VM and the Importer defined in this
// jsonnext package. An application can populate this struct directly from
// whatever source of configuration it uses and use it to configure the jsonnet
// VM. This package provides two options for populating it from the command line
// (Go flags or Kong).
type Config struct {
	ImportPath []string `name:"jpath" sep:"none" short:"J" placeholder:"dir" help:"Add a library search dir"`
	ExtVars    VMVarMap `kong:"-"`
	TLAVars    VMVarMap `kong:"-"`
	MaxStack   int      `default:"500" help:"Number of allowed stack frames of jsonnet VM"`
	MaxTrace   int      `default:"20" help:"Maximum number of stack frames output on error"`
}

// NewConfig returns a new initialised but empty Config struct.
func NewConfig() *Config {
	return &Config{
		ExtVars:  VMVarMap{},
		TLAVars:  VMVarMap{},
		MaxStack: maxStackDepth,
		MaxTrace: maxStackTraceOutput,
	}
}

// MakeVM returns a jsonnet.VM configured with the external vars and top-level
// args in Config, and sets its Importer to be a jsonnext.Importer also
// configured from Config. The importer search path is extended with elements
// from the environment variable pathEnvVar. pathEnvVar may be the empty string
// if no environment variable should be used.
func (c *Config) MakeVM(pathEnvVar string) *jsonnet.VM {
	vm := jsonnet.MakeVM()
	i := &Importer{}
	vm.Importer(i)
	c.ConfigureImporter(i, pathEnvVar)
	c.ConfigureVM(vm)
	return vm
}

// ConfigureImporter sets up a jsonnext.Importer with the import path from
// the config and from a PATH-style environment variable. If envvar is the
// empty string, no paths are taken from the environment.
func (c *Config) ConfigureImporter(i *Importer, envvar string) {
	i.SearchPath = c.ImportPath
	if envvar != "" {
		i.AppendSearchFromEnv(envvar)
	}
}

// ConfigureVM sets the ExtVars and TLAVars in the jsonnet VM.
func (c *Config) ConfigureVM(vm *jsonnet.VM) {
	c.ExtVars.ConfigureVM(vm)
	c.TLAVars.ConfigureVM(vm)
	vm.MaxStack = c.MaxStack
	vm.ErrorFormatter.SetMaxStackTraceSize(c.MaxTrace)
}

// VMVarMap is a map of VMVars that contains a common namespace for variable
// names. The values of type VMVar know how to set themselves in a given VM.
type VMVarMap map[string]VMVar

// ConfigureVM sets the vars from m in the givem jsonnet VM.
func (m VMVarMap) ConfigureVM(vm *jsonnet.VM) {
	for key, vmvar := range m {
		vmvar.Set(key, vm)
	}
}

// SetVar sets a variable in m parsing the key and value from the given string
// v, using makevar to construct the VMVar value. If the value is omitted from
// the string, the value is taken from an environment variable of the same name
// as the key. An error is returned if the environment variable does not exist,
// or if the value string cannot be parsed due to a missing key or value (when
// required).
//
// makevar will typically be one of the VMVar constructor functions in this
// package - New{Ext,TLA}{Str,Code}{,File}.
func (m VMVarMap) SetVar(v string, makevar func(string) VMVar) error {
	parts := strings.SplitN(v, "=", 2)
	if parts[0] == "" {
		return errs.Errorf(`%v in "%s"`, ErrMissingKey, v)
	}
	if len(parts) == 1 {
		val, ok := os.LookupEnv(parts[0])
		if !ok {
			return errs.Errorf("%v: %s", ErrMissingValue, parts[0])
		}
		parts = append(parts, val)
	}
	m[parts[0]] = makevar(parts[1])
	return nil
}

// VMVar is a variable that can be set in a jsonnet VM, either as an external
// variable (extVar) or a top-level arg (TLA), as a string or code. Variants of
// VMVars that take the string or code from a file are turned into jsonnet
// imports (jsonnet "importstr" or "import" statements) and added as a code
// VMVar. It maps to the --ext-* and --tla-* command line arguments to the
// standard jsonnet binary.
type VMVar interface {
	Set(key string, vm *jsonnet.VM)
}

type (
	extStr      string
	extCode     string
	extStrFile  string
	extCodeFile string
	tlaStr      string
	tlaCode     string
	tlaStrFile  string
	tlaCodeFile string
)

// NewExtStr constructs a VMVar as a string external variable.
func NewExtStr(s string) VMVar                  { return extStr(s) }
func (v extStr) Set(key string, vm *jsonnet.VM) { vm.ExtVar(key, string(v)) }

// NewExtCode constructs a VMVar as a code external variable.
func NewExtCode(s string) VMVar                  { return extCode(s) }
func (v extCode) Set(key string, vm *jsonnet.VM) { vm.ExtCode(key, string(v)) }

// NewExtStrFile constructs a VMVar as a string external variable to be read
// from a file.
func NewExtStrFile(s string) VMVar                  { return extStrFile(s) }
func (v extStrFile) Set(key string, vm *jsonnet.VM) { vm.ExtCode(key, mkImportStr(string(v))) }

// NewExtCodeFile constructs a VMVar as a code external variable to be read
// from a file.
func NewExtCodeFile(s string) VMVar                  { return extCodeFile(s) }
func (v extCodeFile) Set(key string, vm *jsonnet.VM) { vm.ExtCode(key, mkImport(string(v))) }

// NewTLAStr constructs a VMVar as a string top-level arg.
func NewTLAStr(s string) VMVar                  { return tlaStr(s) }
func (v tlaStr) Set(key string, vm *jsonnet.VM) { vm.TLAVar(key, string(v)) }

// NewTLACode constructs a VMVar as a code top-level arg.
func NewTLACode(s string) VMVar                  { return tlaCode(s) }
func (v tlaCode) Set(key string, vm *jsonnet.VM) { vm.TLACode(key, string(v)) }

// NewTLAStrFile constructs a VMVar as a string top-level arg to be read
// from a file.
func NewTLAStrFile(s string) VMVar                  { return tlaStrFile(s) }
func (v tlaStrFile) Set(key string, vm *jsonnet.VM) { vm.TLACode(key, mkImportStr(string(v))) }

// NewTLACodeFile constructs a VMVar as a code top-level arg to be read
// from a file.
func NewTLACodeFile(s string) VMVar                  { return tlaCodeFile(s) }
func (v tlaCodeFile) Set(key string, vm *jsonnet.VM) { vm.TLACode(key, mkImport(string(v))) }

// Quote string using verbatim string: @'...'.
func quoteStr(s string) string    { return "@'" + strings.ReplaceAll(s, "'", "''") + "'" }
func mkImport(f string) string    { return "import " + quoteStr(f) }
func mkImportStr(f string) string { return "importstr " + quoteStr(f) }
