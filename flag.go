package jsonnext

import (
	"flag"
	"fmt"
)

// ConfigFlags defines a set of flags in the given FlagSet for a Config struct
// to populate its fields from the command line. The return value is a pointer
// to the Config struct that stores the values of the flags.
// The fields are populated with the following flags:
//  Config.ImportPath:
//   -J, -jpath
//  Config.ExtVars:
//   -V, -ext-str: ext var as string literal
//   -ext-code: ext var as code literal
//   -ext-str-file: ext var as string from file
//   -ext-code-file: ext var as code from file
//  Config.TLAVars:
//   -A, -tla-str: top-level arg as string literal
//   -tla-code: top-level arg as code literal
//   -tla-str-file: top-level arg as string from file
//   -tla-code-file: top-level arg as code from file
func ConfigFlags(fs *flag.FlagSet) *Config {
	c := NewConfig()
	ConfigFlagsVar(fs, c)
	return c
}

// ConfigFlagsVar defines a set of flags in the given FlagSet for a Config
// struct to populate the fields from the command line. The argument c points
// to the Config struct to populate. The set of flags defined is described in
// the ConfigFlags function description.
func ConfigFlagsVar(fs *flag.FlagSet, c *Config) {
	StringSliceVar(fs, &c.ImportPath, "jpath", "Add a library search `dir`")
	ExtStrVar(fs, c.ExtVars, "ext-str", "Add extVar `var[=str]` (from environment if <str> is omitted)")
	ExtCodeVar(fs, c.ExtVars, "ext-code", "Add extVar `var[=code]` (from environment if <code> is omitted)")
	ExtStrFileVar(fs, c.ExtVars, "ext-str-file", "Add extVar `var=file` string from a file")
	ExtCodeFileVar(fs, c.ExtVars, "ext-code-file", "Add extVar `var=file` code from a file")

	TLAStrVar(fs, c.TLAVars, "tla-str", "Add top-level arg `var=[=str]` (from environment if <str> is omitted)")
	TLACodeVar(fs, c.TLAVars, "tla-code", "Add top-level arg `var[=code]` (from environment if <code> is omitted)")
	TLAStrFileVar(fs, c.TLAVars, "tla-str-file", "Add top-level arg `var=file` string from a file")
	TLACodeFileVar(fs, c.TLAVars, "tla-code-file", "Add top-level arg `var=file` code from a file")

	// Add short flags. TODO(camh): consider making these optional.
	StringSliceVar(fs, &c.ImportPath, "J", "Add a library search `dir`")
	ExtStrVar(fs, c.ExtVars, "V", "Add extVar `var[=str]` (from environment if <str> is omitted)")
	TLAStrVar(fs, c.TLAVars, "A", "Add top-level arg `var[=str]` (from environment if <str> is omitted)")
}

// StringSliceVar defines a flag in the given FlagSet with the given name and
// usage string. The argument p is a pointer to a []string variable in which to
// store the value of the flag. The value given to each instance of the flag is
// appended to the slice.
func StringSliceVar(fs *flag.FlagSet, p *[]string, name, usage string) {
	fs.Var((*stringSliceValue)(p), name, usage)
}

// StringSlice defines a flag in the given FlagSet with the given name and
// usage string. The return value is the address of a []string variable that
// stores the value of the flag. The value given to each instance of the flag
// is appended to the slice.
func StringSlice(fs *flag.FlagSet, name, usage string) *[]string {
	p := &[]string{}
	StringSliceVar(fs, p, name, usage)
	return p
}

type stringSliceValue []string

func (s *stringSliceValue) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func (s *stringSliceValue) Get() interface{} { return []string(*s) }
func (s *stringSliceValue) String() string   { return fmt.Sprint(*s) }

// ExtStrVar defines flag with the given name and usage string in the given
// FlagSet to set a VMVar in the given VMVarMap. The VMVar sets an extVar
// string literal in a jsonnet VM.
//
// The flag value on the command line is parsed as "key[=string]". If "=string"
// is omitted then key is looked up in the environment for the value. It is an
// error if "=string" is not provided and the environment variable does not
// exist.
//
// The flag can be repeated multiple times on the command line. Each instance
// adds a new value to the given VMVarMap.
func ExtStrVar(fs *flag.FlagSet, m VMVarMap, name, usage string) {
	fs.Var(vmVarMapValue{m, NewExtStr}, name, usage)
}

// ExtStrFileVar defines flag with the given name and usage string in the given
// FlagSet to set a VMVar in the given VMVarMap. The VMVar sets an extVar
// string import in a jsonnet VM.
//
// The flag value on the command line is parsed as "key[=filename]". If
// "=filename" is omitted then key is looked up in the environment for the
// value. It is an error if "=filename" is not provided and the environment
// variable does not exist.
//
// The flag can be repeated multiple times on the command line. Each instance
// adds a new value to the given VMVarMap.
func ExtStrFileVar(fs *flag.FlagSet, m VMVarMap, name, usage string) {
	fs.Var(vmVarMapValue{m, NewExtStrFile}, name, usage)
}

// ExtCodeVar defines flag with the given name and usage string in the given
// FlagSet to set a VMVar in the given VMVarMap. The VMVar sets an extVar code
// literal in a jsonnet VM.
//
// The flag value on the command line is parsed as "key[=code]". If "=code" is
// omitted then key is looked up in the environment for the value. It is an
// error if "=code" is not provided and the environment variable does not
// exist.
//
// The flag can be repeated multiple times on the command line. Each instance
// adds a new value to the given VMVarMap.
func ExtCodeVar(fs *flag.FlagSet, m VMVarMap, name, usage string) {
	fs.Var(vmVarMapValue{m, NewExtCode}, name, usage)
}

// ExtCodeFileVar defines flag with the given name and usage string in the
// given FlagSet to set a VMVar in the given VMVarMap. The VMVar sets an extVar
// code import in a jsonnet VM.
//
// The flag value on the command line is parsed as "key[=filename]". If
// "=filename" is omitted then key is looked up in the environment for the
// value. It is an error if "=filename" is not provided and the environment
// variable does not exist.
//
// The flag can be repeated multiple times on the command line. Each instance
// adds a new value to the given VMVarMap.
func ExtCodeFileVar(fs *flag.FlagSet, m VMVarMap, name, usage string) {
	fs.Var(vmVarMapValue{m, NewExtCodeFile}, name, usage)
}

// TLAStrVar defines flag with the given name and usage string in the given
// FlagSet to set a VMVar in the given VMVarMap. The VMVar sets a top-level arg
// string literal in a jsonnet VM.
//
// The flag value on the command line is parsed as "key[=string]". If "=string"
// is omitted then key is looked up in the environment for the value. It is an
// error if "=string" is not provided and the environment variable does not
// exist.
//
// The flag can be repeated multiple times on the command line. Each instance
// adds a new value to the given VMVarMap.
func TLAStrVar(fs *flag.FlagSet, m VMVarMap, name, usage string) {
	fs.Var(vmVarMapValue{m, NewTLAStr}, name, usage)
}

// TLAStrFileVar defines flag with the given name and usage string in the given
// FlagSet to set a VMVar in the given VMVarMap. The VMVar sets a top-level arg
// string import in a jsonnet VM.
//
// The flag value on the command line is parsed as "key[=filename]". If
// "=filename" is omitted then key is looked up in the environment for the
// value. It is an error if "=filename" is not provided and the environment
// variable does not exist.
//
// The flag can be repeated multiple times on the command line. Each instance
// adds a new value to the given VMVarMap.
func TLAStrFileVar(fs *flag.FlagSet, m VMVarMap, name, usage string) {
	fs.Var(vmVarMapValue{m, NewTLAStrFile}, name, usage)
}

// TLACodeVar defines flag with the given name and usage string in the given
// FlagSet to set a VMVar in the given VMVarMap. The VMVar sets a top-level arg
// code literal in a jsonnet VM.
//
// The flag value on the command line is parsed as "key[=code]". If "=code" is
// omitted then key is looked up in the environment for the value. It is an
// error if "=code" is not provided and the environment variable does not
// exist.
//
// The flag can be repeated multiple times on the command line. Each instance
// adds a new value to the given VMVarMap.
func TLACodeVar(fs *flag.FlagSet, m VMVarMap, name, usage string) {
	fs.Var(vmVarMapValue{m, NewTLACode}, name, usage)
}

// TLACodeFileVar defines flag with the given name and usage string in the
// given FlagSet to set a VMVar in the given VMVarMap. The VMVar sets a
// top-level arg code import in a jsonnet VM.
//
// The flag value on the command line is parsed as "key[=filename]". If
// "=filename" is omitted then key is looked up in the environment for the
// value. It is an error if "=filename" is not provided and the environment
// variable does not exist.
//
// The flag can be repeated multiple times on the command line. Each instance
// adds a new value to the given VMVarMap.
func TLACodeFileVar(fs *flag.FlagSet, m VMVarMap, name, usage string) {
	fs.Var(vmVarMapValue{m, NewTLACodeFile}, name, usage)
}

// vmVarMapValue wraps a VMVarMap with the additional information it needs to
// parse the various types of VMVars, using a makevar function to construct the
// values in the map.
type vmVarMapValue struct {
	m       VMVarMap
	makevar func(string) VMVar
}

// Set sets the flag var value from v. It implements the flag.Value interface.
func (mv vmVarMapValue) Set(v string) error {
	return mv.m.SetVar(v, mv.makevar)
}

// String returns a string representation of the value. It implements the
// flag.Value interface.
func (mv vmVarMapValue) String() string {
	return fmt.Sprint(mv.m)
}

// Get returns the underlying VMVarMap. It implements the flag.Getter interface.
func (mv vmVarMapValue) Get() interface{} {
	return mv.m
}
