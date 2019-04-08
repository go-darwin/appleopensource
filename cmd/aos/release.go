// Copyright 2017 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/urfave/cli"

	"github.com/zchee/appleopensource"
	"github.com/zchee/appleopensource/pkg/fs"
)

var releaseCommand = cli.Command{
	Name:  "release",
	Usage: "List all projects included to the releases available to opensource.apple.com.",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "suppress some output",
		},
	},
	Subcommands: []cli.Command{
		{
			Name:   "macos",
			Usage:  "macOS release",
			Action: runReleaseMacOS,
		},
		{
			Name:   "xcode",
			Usage:  "Developer Tool(Xcode) release",
			Action: runReleaseXcode,
		},
		{
			Name:   "ios",
			Usage:  "iOS release",
			Action: runReleaseMacOS,
		},
		{
			Name:   "server",
			Usage:  "macOS Server release",
			Action: runReleaseMacOS,
		},
	},
	Before: initRelease,
}

var (
	releaseVersion string
	releaseList    string
	releaseQuiet   bool
)

var releaseCachedir = filepath.Join(cacheDir(), "release")

func initRelease(ctx *cli.Context) error {
	releaseVersion = ctx.Args().Get(1)
	releaseList = ctx.String("list")
	releaseQuiet = ctx.Bool("quiet")

	releaseCachedir = ctx.Args().Get(1)
	if err := fs.MkdirAll(releaseCachedir, 0700); err != nil {
		return err
	}
	return nil
}

func indexRelease(ctx *cli.Context, platform appleopensource.Platform, version string) ([]byte, error) {
	fname := filepath.Join(releaseCachedir, fmt.Sprintf("%s_%s.html", platform, strings.Replace(version, ".", "", -1)))
	if fs.IsExist(fname) && !noCache {
		return ioutil.ReadFile(fname)
	}

	buf, err := appleopensource.IndexRelease(platform, version)
	if err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(fname, buf, 0664); err != nil {
		return nil, err
	}

	return buf, nil
}

func runRelease(ctx *cli.Context, platform appleopensource.Platform, version string) error {
	if !releaseQuiet {
		fmt.Printf("Release version: %s\n", version)
	}

	release, err := indexRelease(ctx, platform, version)
	if err != nil {
		log.Fatal(err.Error())
	}

	list, err := appleopensource.ListRelease(release)
	if err != nil {
		log.Fatal(err.Error())
	}

	var buf bytes.Buffer
	tbuf := tabwriter.NewWriter(&buf, 2, 1, 2, ' ', 0)

	for _, b := range list {
		if !releaseQuiet {
			if b.Updated {
				tbuf.Write([]byte("\u2022 ")) // u2022: •
			} else {
				tbuf.Write([]byte("  "))
			}
		}
		tbuf.Write([]byte(fmt.Sprintf("%s\t%s", b.Name, b.Version)))
		if !releaseQuiet {
			tbuf.Write([]byte("\t"))
			if b.ComingSoon {
				tbuf.Write([]byte(appleopensource.ComingSoon))
			}
		}
		tbuf.Write([]byte("\n"))
	}
	tbuf.Flush()

	fmt.Printf(buf.String())

	return nil
}

func runReleaseMacOS(ctx *cli.Context) error {
	return runRelease(ctx, appleopensource.MacOS, releaseVersion)
}

func runReleaseXcode(ctx *cli.Context) error {
	return runRelease(ctx, appleopensource.Xcode, releaseVersion)
}

func runReleaseIOS(ctx *cli.Context) error {
	return runRelease(ctx, appleopensource.IOS, releaseVersion)
}

func runReleaseServer(ctx *cli.Context) error {
	return runRelease(ctx, appleopensource.Server, releaseVersion)
}
