// Package conformance is a test suite for testing CLI parsing into jsonnext
// data structures. It can be used to check conformance of a CLI parsing
// implementation.
package conformance

import (
	"strings"
	"testing"

	"foxygo.at/s/test"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"foxygo.at/jsonnext"
)

// CLIParser should be implemented by a jsonnext.Config CLI parsing
// implementation to test its conformance to the standard parsing expected.
// Parse() is called with a given command line for which the parse must return
// a *jsonnext.Config corresponding to that set of CLI flags or return an error
// where those flags fail to parse..
type CLIParser interface {
	Parse(t *testing.T, args []string) (*jsonnext.Config, error)
}

// Suite is a conformance test for command line parsers that parse into a
// jsonnext.Config struct. The test case for the parser under test must
// implement CLIParser to parse a given command line into a *jsonnext.Config
// struct.
type Suite struct {
	suite.Suite
	parser CLIParser
}

// NewSuite returns a new Suite to test the conformance of the given parser.
func NewSuite(s CLIParser) *Suite {
	return &Suite{parser: s}
}

// Run runs the conformance test in this package against the CLIParser
// configured in the Suite.
func (s *Suite) Run(t *testing.T) {
	suite.Run(t, s)
}

var flags = map[string]struct {
	flag    string
	makevar func(string) jsonnext.VMVar
}{
	"ext-str":       {"--ext-str", jsonnext.NewExtStr},
	"ext-str-short": {"-V", jsonnext.NewExtStr},
	"ext-code":      {"--ext-code", jsonnext.NewExtCode},
	"ext-str-file":  {"--ext-str-file", jsonnext.NewExtStrFile},
	"ext-code-file": {"--ext-code-file", jsonnext.NewExtCodeFile},
	"tla-str":       {"--tla-str", jsonnext.NewTLAStr},
	"tla-str-short": {"-A", jsonnext.NewTLAStr},
	"tla-code":      {"--tla-code", jsonnext.NewTLACode},
	"tla-str-file":  {"--tla-str-file", jsonnext.NewTLAStrFile},
	"tla-code-file": {"--tla-code-file", jsonnext.NewTLACodeFile},
}

// TestVMVars tests that all styles of VMVars are correctly parsed on the
// command line, either as literals on the command line or from the
// environment.
func (s *Suite) TestVMVars() {
	for name, tc := range flags {
		name, tc := name, tc
		s.T().Run(name, func(t *testing.T) {
			test.Env.Set("environ", "hello")
			defer test.Env.Restore()
			input := []string{t.Name(), tc.flag, "environ", tc.flag, "literal=world"}
			expectedVars := jsonnext.VMVarMap{
				"environ": tc.makevar("hello"),
				"literal": tc.makevar("world"),
			}
			expected := jsonnext.NewConfig()
			if strings.HasPrefix(name, "ext-") {
				expected.ExtVars = expectedVars
			} else {
				expected.TLAVars = expectedVars
			}

			cfg, err := s.parser.Parse(t, input)
			require.NoError(t, err)
			require.Equal(t, expected, cfg)
		})
	}
}

// TestVMVarErr tests that all styles of VMVar flags generate errors when both
// a value is not provided with the corresponding environment variable unset,
// and when no variable name is provided.
func (s *Suite) TestVMVarErr() {
	for name, tc := range flags {
		tc := tc
		s.T().Run("err-"+name, func(t *testing.T) {
			test.Env.Unset("missing")
			defer test.Env.Restore()
			_, err := s.parser.Parse(t, []string{t.Name(), tc.flag, "missing"})
			require.Error(t, err)
			_, err = s.parser.Parse(t, []string{t.Name(), tc.flag, "=foo"})
			require.Error(t, err)
		})
	}
}

// TestVMVarOverride tests that when the same variable name is used for two
// VMVars of the same type (ext or TLA), that the last one has precedence and
// is present in the output configuration.
func (s *Suite) TestVMVarOverride() {
	t := s.T()
	f1 := flags["ext-str"]
	f2 := flags["ext-code"]

	cfg, err := s.parser.Parse(t, []string{t.Name(), f1.flag, "var=value1", f2.flag, "var=value2"})
	require.NoError(t, err)

	expected := jsonnext.NewConfig()
	expected.ExtVars = jsonnext.VMVarMap{"var": f2.makevar("value2")}
	require.Equal(t, expected, cfg)

	f3 := flags["tla-str"]
	f4 := flags["tla-code"]

	cfg, err = s.parser.Parse(t, []string{t.Name(), f3.flag, "var=value1", f4.flag, "var=value2"})
	require.NoError(t, err)

	expected = jsonnext.NewConfig()
	expected.TLAVars = jsonnext.VMVarMap{"var": f4.makevar("value2")}
	require.Equal(t, expected, cfg)
}

// TestImportPath tests that the ImportPath field is set by the --jpath and
// -J flags.
func (s *Suite) TestImportPath() {
	t := s.T()

	cfg, err := s.parser.Parse(t, []string{t.Name(), "--jpath", "a", "-J", "b"})
	require.NoError(t, err)
	expected := jsonnext.NewConfig()
	expected.ImportPath = []string{"a", "b"}
	require.Equal(t, expected, cfg)
}

// TestMaxStack tests that the MaxStack field is set by the --max-stack flag,
// and that the default is set when the flag is not present.
func (s *Suite) TestMaxStack() {
	t := s.T()

	cfg, err := s.parser.Parse(t, []string{t.Name(), "--max-stack", "10"})
	require.NoError(t, err)
	require.Equal(t, 10, cfg.MaxStack)

	cfg, err = s.parser.Parse(t, []string{t.Name()})
	require.NoError(t, err)
	require.Equal(t, 500, cfg.MaxStack)
}

// TestMaxTrace tests that the MaxTrace field is set by the --max-trace flag,
// and that the default is set when the flag is not present.
func (s *Suite) TestMaxTrace() {
	t := s.T()

	cfg, err := s.parser.Parse(t, []string{t.Name(), "--max-trace", "10"})
	require.NoError(t, err)
	require.Equal(t, 10, cfg.MaxTrace)

	cfg, err = s.parser.Parse(t, []string{t.Name()})
	require.NoError(t, err)
	require.Equal(t, 20, cfg.MaxTrace)
}
