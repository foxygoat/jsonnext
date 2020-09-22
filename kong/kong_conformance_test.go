package kong_test

import (
	"testing"

	"github.com/alecthomas/kong"

	"foxygo.at/jsonnext"
	"foxygo.at/jsonnext/conformance"
	jxkong "foxygo.at/jsonnext/kong"
)

type suite struct{}

func (s *suite) Parse(t *testing.T, args []string) (*jsonnext.Config, error) {
	kcfg := jxkong.NewConfig()
	parser, err := kong.New(kcfg)
	if err != nil {
		return nil, err
	}
	_, err = parser.Parse(args[1:])
	return kcfg.Config, err
}

func TestConformance(t *testing.T) {
	s := conformance.NewSuite(&suite{})
	s.Run(t)
}
