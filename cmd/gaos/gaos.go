// Copyright 2016 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "github.com/alecthomas/kingpin"

var (
	cli   = kingpin.New("gaos", "An opensource.apple.com resource management cli tool.")
	debug = cli.Flag("debug", "Enable debug mode.").Bool()

	cmdList     = cli.Command("list", "List all package available to opensource.apple.com.")
	cmdVersions = cli.Command("versions", "List all <package> versions available to opensource.apple.com.")
)

var (
	Version   = "0.0.1"
	GitCommit = "HEAD"
)

func init() {
	cli.Version(Version)
}

func main() {
	kingpin.Parse()
}
