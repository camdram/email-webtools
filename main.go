package main

import "github.com/camdram/email-webtools/cmd"

// Software version defaults to the value below but is overridden by the compiler in Makefile.
var version = "dev-edge"

func main() {
	cmd.Execute(version)
}
