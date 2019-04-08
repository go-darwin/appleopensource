// Copyright 2017 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"github.com/urfave/cli"
	"github.com/zchee/appleopensource/pkg/appleopensource"
)

var fetchCommand = cli.Command{
	Name:      "fetch",
	Usage:     "Fetch the tarballs",
	Before:    initFetch,
	Action:    runFetch,
	ArgsUsage: "<project> <version> <dist>",
}

var (
	fetchProject string
	fetchVersion []string
	fetchDistDir string
)

func initFetch(ctx *cli.Context) error {
	if err := checkArgs(ctx, 3, minArgs, "project version dist"); err != nil {
		return err
	}

	fetchProject = ctx.Args().First()
	args := ctx.Args().Tail()
	fetchVersion = args[:len(args)-1]
	fetchDistDir = args[len(args)-1]

	return nil
}

func runFetch(ctx *cli.Context) error {
	dlList := []string{}
	for _, v := range fetchVersion {
		p := appleopensource.Project{
			Name:    fetchProject,
			Version: v,
		}
		dlList = append(dlList, p.Tarball())
	}
	if err := appleopensource.Fetch(fetchDistDir, dlList...); err != nil {
		return err
	}

	return nil
}
