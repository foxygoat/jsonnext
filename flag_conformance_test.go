package jsonnext_test

import (
	"flag"
	"testing"

	"foxygo.at/jsonnext"
	"foxygo.at/jsonnext/conformance"
)

type suite struct{}

func (s *suite) Parse(t *testing.T, args []string) (*jsonnext.Config, error) {
	fs := new(flag.FlagSet)
	cfg := jsonnext.ConfigFlags(fs)
	return cfg, fs.Parse(args[1:])
}

func TestConformance(t *testing.T) {
	s := conformance.NewSuite(&suite{})
	s.Run(t)
}
