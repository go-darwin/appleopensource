// Copyright 2016 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import cli "github.com/alecthomas/kingpin"

var (
	Version = "0.0.1"

	debug = cli.Flag("debug", "Enable debug mode.").Bool()

	cmdList     = cli.Command("list", "List all project available to opensource.apple.com.")
	cmdVersions = cli.Command("versions", "List all versions of the <project> available to opensource.apple.com.")
	cmdRelease  = cli.Command("release", "List all projects included to the releases available to opensource.apple.com.")
)

func init() {
	cli.CommandLine.Name = "gaos"
	cli.CommandLine.Help = "An opensource.apple.com resource management tool."
	cli.CommandLine.HelpFlag.Short('h')
	cli.CommandLine.UsageTemplate(cli.CompactUsageTemplate)

	cmdList.Action(runList)
	cmdVersions.Action(runVersions)
}

func main() {
	cli.Parse()
}
