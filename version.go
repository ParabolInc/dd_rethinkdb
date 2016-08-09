package main

import "fmt"

const (
	Name = "dd_rethinkdb"
)

var (
	// Version should be updated by hand at each release
	Version = "0.1.0-dev"

	// GitCommit will be overwritten automatically by the build system
	GitCommit = "HEAD"
)

func NameVersion() string {
	return fmt.Sprintf("%s version %s, build %s", Name, Version, GitCommit)
}
