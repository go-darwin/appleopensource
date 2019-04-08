// Copyright 2017 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/pkgutil/osutil"
	"github.com/urfave/cli"
	"github.com/zchee/appleopensource"
)

var versionsCommand = cli.Command{
	Name:  "versions",
	Usage: "List all versions of the project available to opensource.apple.com.",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "tarballs, t",
			Usage: "List the tarballs resource versions",
		},
		cli.BoolFlag{
			Name:  "source, s",
			Usage: "List the source resource versions",
		},
	},
	ArgsUsage: "<project name>",
	Before:    initVersions,
	Action:    runVersions,
}

var (
	versionsProject  string
	versionsSource   bool
	versionsTarballs bool
)

func initVersions(ctx *cli.Context) error {
	versionsProject = ctx.Args().First()
	versionsSource = ctx.Bool("source")
	versionsTarballs = ctx.Bool("tarballs")
	return nil
}

// index return the opensource.apple.com project index, and caches the HTML DOM tree into cacheDir.
func indexVersion(project string, typ appleopensource.ResourceType) ([]byte, error) {
	versionsCachedir := filepath.Join(cacheDir(), typ.String())
	if err := osutil.MkdirAll(versionsCachedir, 0700); err != nil {
		return nil, err
	}

	fname := filepath.Join(versionsCachedir, fmt.Sprintf("%s.html", project))
	if osutil.IsExist(fname) && !noCache {
		return ioutil.ReadFile(fname)
	}

	buf, err := appleopensource.IndexVersion(project, typ)
	if err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(fname, buf, 0664); err != nil {
		return nil, err
	}

	return buf, nil
}

func runVersions(ctx *cli.Context) error {
	if err := checkArgs(ctx, 1, minArgs, "project name"); err != nil {
		return err
	}

	mode := appleopensource.TarballsResource
	switch {
	case versionsTarballs:
		// nothing to do
	case versionsSource:
		mode = appleopensource.SourceResource
	}

	buf, err := indexVersion(versionsProject, mode)
	if err != nil {
		return err
	}

	list, err := appleopensource.ListVersions(buf)
	if err != nil {
		return err
	}

	fmt.Println(strings.Join(list, "\n"))

	return nil
}
