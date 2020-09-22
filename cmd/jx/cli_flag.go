// +build flag

package main

import (
	"flag"
	"os"

	"foxygo.at/jsonnext"
)

// Parse CLI using Go's flag package and the helpers in jsonnext.
func parseCLI() *config {
	c := &config{}
	c.Config.Config = jsonnext.ConfigFlags(flag.CommandLine)

	flag.Parse()
	if flag.NArg() > 1 {
		flag.Usage()
		os.Exit(1)
	} else if flag.NArg() == 1 {
		c.Filename = flag.Args()[0]
	}

	return c
}
