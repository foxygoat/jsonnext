// Package kong parses command line flags into jsonnext.Config using kong.
//
//   import "github.com/alecthomas/kong"
//
// Kong is a flag type library supporting complex command-line structures with
// minimal developer effort: CLI args are expressed as Go types and field tags.
//
// This package provides a Config struct that extends jsonnext.Config, adding
// fields and tags for kong to build a CLI parser. Some fields in
// jsonnext.Config already have kong field tags, but others need special
// handling as the default kong features do not support the style of CLI
// parsing that is compatible with the standard jsonnet CLI.
//
// This functionality is split into a separate sub-package so users of jsonnext
// do not need to depend on kong if the do not use it.
package kong

import (
	"foxygo.at/jsonnext"
	"github.com/alecthomas/kong"
)

// Config embeds jsonnext.Config to add kong command line parsing for external
// variables and top-level arguments. The flags for these options are parsed
// into the VMVarMaps in the embedded jsonnext.Config field.
//
// Since Config embeds jsonnext.Config, the method set of jsonnext.Config is
// also present on Config.
type Config struct {
	*jsonnext.Config

	// vmVarMaps is embedded and unexported, which prevents users of this
	// struct from accessing the vmVarMaps directly, but still allows kong
	// to access it to make flags from its contents.
	vmVarMaps
}

type vmVarMaps struct {
	ExtStr      vmVarMap `placeholder:"var[=str]" help:"Set extVar string (str from env if omitted)" short:"V"`
	ExtStrFile  vmVarMap `placeholder:"var[=filename]" help:"Set extVar string from a file (filename from env if omitted)"`
	ExtCode     vmVarMap `placeholder:"var[=code]" help:"Set extVar code (code from env if omitted)"`
	ExtCodeFile vmVarMap `placeholder:"var[=filename]" help:"Set extVar code from a file (filename from env if omitted)"`
	TLAStr      vmVarMap `placeholder:"var[=str]" help:" Set top-level arg string (str from env if omitted)" short:"A"`
	TLAStrFile  vmVarMap `placeholder:"var[=filename]" help:"Set top-level arg string from a file (filename from env if omitted)"`
	TLACode     vmVarMap `placeholder:"var[=code]" help:"Set top-level arg code (code from env if omitted)"`
	TLACodeFile vmVarMap `placeholder:"var[=filename]" help:"Set top-level arg code from a file (filename from env if omitted)"`
}

// NewConfig returns an initialised Config struct embedding a jsonnext.Config.
func NewConfig() *Config {
	c := jsonnext.NewConfig()
	return &Config{
		Config: c,
		vmVarMaps: vmVarMaps{
			ExtStr:      vmVarMap{c.ExtVars, jsonnext.NewExtStr, "string"},
			ExtStrFile:  vmVarMap{c.ExtVars, jsonnext.NewExtStrFile, "filename"},
			ExtCode:     vmVarMap{c.ExtVars, jsonnext.NewExtCode, "code"},
			ExtCodeFile: vmVarMap{c.ExtVars, jsonnext.NewExtCodeFile, "filename"},
			TLAStr:      vmVarMap{c.TLAVars, jsonnext.NewTLAStr, "string"},
			TLAStrFile:  vmVarMap{c.TLAVars, jsonnext.NewTLAStrFile, "filename"},
			TLACode:     vmVarMap{c.TLAVars, jsonnext.NewTLACode, "code"},
			TLACodeFile: vmVarMap{c.TLAVars, jsonnext.NewTLACodeFile, "filename"},
		},
	}
}

type vmVarMap struct {
	m       jsonnext.VMVarMap
	makevar func(string) jsonnext.VMVar
	valtype string
}

func (v *vmVarMap) Decode(ctx *kong.DecodeContext) error {
	// v will be the zero value the first time Decode is called on it as
	// kong creates the values fresh and does not use the values in the
	// Config struct. ctx.Value.Target contains that value from the Config
	// struct, so if we are the zero value, initialise ourself from
	// ctx.Value.Target.
	if v.m == nil {
		*v = ctx.Value.Target.Interface().(vmVarMap)
	}
	var valstr string
	if err := ctx.Scan.PopValueInto(v.valtype, &valstr); err != nil {
		return err
	}
	// Literals (i.e. not from files) can come from the environment
	return v.m.SetVar(valstr, v.makevar)
}
