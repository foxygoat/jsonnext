package jsonnext

import (
	"errors"
	"strings"
	"testing"

	"foxygo.at/s/test"
	jsonnet "github.com/google/go-jsonnet"
	"github.com/stretchr/testify/require"
)

func TestVars(t *testing.T) {
	const (
		extVar = iota
		tlaVar
	)
	tests := map[string]struct {
		varType int
		value   VMVar
	}{
		"ext-str":       {extVar, NewExtStr("val")},
		"ext-code":      {extVar, NewExtCode(`"val"`)},
		"tla-str":       {tlaVar, NewTLAStr("val")},
		"tla-code":      {tlaVar, NewTLACode(`"val"`)},
		"ext-str-file":  {extVar, NewExtStrFile("testdata/config/str")},
		"ext-code-file": {extVar, NewExtCodeFile("testdata/config/code")},
		"tla-str-file":  {tlaVar, NewTLAStrFile("testdata/config/str")},
		"tla-code-file": {tlaVar, NewTLACodeFile("testdata/config/code")},
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			c := NewConfig()
			var snippet string
			if tc.varType == extVar {
				c.ExtVars["var"] = tc.value
				snippet = "std.extVar('var')"
			} else {
				c.TLAVars["var"] = tc.value
				snippet = "function(var) var"
			}
			vm := jsonnet.MakeVM()
			c.ConfigureVM(vm)
			got, err := vm.EvaluateSnippet("<literal>", snippet)
			require.NoError(t, err)
			require.Equal(t, `"val"`, strings.TrimSuffix(got, "\n"))
		})
	}
}

func TestSetVar(t *testing.T) {
	m := VMVarMap{}
	err := m.SetVar("var=val", NewExtStr)
	require.NoError(t, err)
	require.Contains(t, m, "var")
	require.Equal(t, NewExtStr("val"), m["var"])
}

func TestSetVarFromEnv(t *testing.T) {
	test.Env.Set("var", "val")
	defer test.Env.Restore()
	m := VMVarMap{}
	err := m.SetVar("var", NewExtStr)
	require.NoError(t, err)
	require.Contains(t, m, "var")
	require.Equal(t, NewExtStr("val"), m["var"])
}

func TestSetVarNoKey(t *testing.T) {
	m := VMVarMap{}
	err := m.SetVar("=hello", NewExtStr)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrMissingKey), "error should be ErrMissingKey")
}

func TestSetVarNoValue(t *testing.T) {
	test.Env.Unset("var")
	defer test.Env.Restore()
	m := VMVarMap{}
	err := m.SetVar("var", NewExtStr)
	require.Error(t, err)
	require.True(t, errors.Is(err, ErrMissingValue), "error should be ErrMissingValue")
}

func TestMakeVM(t *testing.T) {
	c := NewConfig()

	c.ImportPath = []string{"testdata"}
	c.ExtVars["var"] = NewExtCode("importstr 'config/str'")

	test.Env.Set("JPATH", "testdata/config")
	defer test.Env.Restore()
	c.TLAVars["var"] = NewTLACode("importstr 'str'")

	vm := c.MakeVM("JPATH")

	got, err := vm.EvaluateSnippet("<literal>", "std.extVar('var')")
	require.NoError(t, err)
	require.Equal(t, `"val"`, strings.TrimSuffix(got, "\n"))

	got, err = vm.EvaluateSnippet("<literal>", "function(var) var")
	require.NoError(t, err)
	require.Equal(t, `"val"`, strings.TrimSuffix(got, "\n"))
}

func TestConfigureImporter(t *testing.T) {
	i := Importer{}
	c := NewConfig()
	c.ImportPath = []string{"a", "b"}
	c.ConfigureImporter(&i, "")
	require.Equal(t, []string{"a", "b"}, i.SearchPath)
}

func TestConfigureImporterFromEnv(t *testing.T) {
	test.Env.Set("JPATH", "c:d")
	defer test.Env.Restore()
	i := Importer{}
	c := NewConfig()
	c.ConfigureImporter(&i, "JPATH")
	require.Equal(t, []string{"c", "d"}, i.SearchPath)
}

func TestConfigureImporterFromConfigAndEnv(t *testing.T) {
	test.Env.Set("JPATH", "c:d")
	defer test.Env.Restore()
	i := Importer{}
	c := NewConfig()
	c.ImportPath = []string{"a", "b"}
	c.ConfigureImporter(&i, "JPATH")
	require.Equal(t, []string{"a", "b", "c", "d"}, i.SearchPath)
}
