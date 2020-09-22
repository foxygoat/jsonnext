// +build !flag

package main

import (
	jxkong "foxygo.at/jsonnext/kong"
	"github.com/alecthomas/kong"
)

// Parse CLI using Kong.
func parseCLI() *config {
	c := &config{Config: *jxkong.NewConfig()}
	kong.Parse(c)
	return c
}
