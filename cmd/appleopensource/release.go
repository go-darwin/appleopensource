// Copyright 2016 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	cli "github.com/alecthomas/kingpin"
	"github.com/zchee/appleopensource"
)

var (
	releaseMacOS        = cmdRelease.Command("macos", "macOS release.")
	releaseMacOSVersion = releaseMacOS.Arg("version", "Release version. Default is latest release.").Default(appleopensource.Release[appleopensource.MacOS][0]).String()
	releaseMacOSList    = releaseMacOS.Flag("list", "List available macOS releases.").Short('l').Bool()

	releaseXcode        = cmdRelease.Command("xcode", "Xcode(Developer Tool) release.")
	releaseXcodeVersion = releaseXcode.Arg("version", "Release version. Default is latest release.").Default(appleopensource.Release[appleopensource.Xcode][0]).String()
	releaseXcodeList    = releaseXcode.Flag("list", "List available Xcode(Developer Tool) releases.").Short('l').Bool()

	releaseIOS        = cmdRelease.Command("ios", "iOS release.")
	releaseIOSVersion = releaseIOS.Arg("version", "Release version. Default is latest release.").Default(appleopensource.Release[appleopensource.IOS][0]).String()
	releaseIOSOSList  = releaseIOS.Flag("list", "List available iOS releases.").Short('l').Bool()

	releaseServer        = cmdRelease.Command("server", "macOS Server release.")
	releaseServerVersion = releaseServer.Arg("version", "Release version. Default is latest release.").Default(appleopensource.Release[appleopensource.Server][0]).String()
	releaseServerList    = releaseServer.Flag("list", "List available Server releases.").Short('l').Bool()

	releaseNocache = cmdRelease.Flag("no-cache", "Disable the cache.").Short('n').Bool()
	releaseQuite   = cmdRelease.Flag("quiet", "suppress some output").Short('q').Bool()
)

func init() {
	releaseMacOS.Action(runReleaseMacOS)
	releaseXcode.Action(runReleaseXcode)
	releaseIOS.Action(runReleaseIOS)
	releaseServer.Action(runReleaseServer)
}

func indexRelease(platform appleopensource.Platform, version string) ([]byte, error) {
	cachedir := filepath.Join(cacheDir(), "release")

	fname := filepath.Join(cachedir, fmt.Sprintf("%s_%s.html", platform, strings.Replace(version, ".", "", -1)))
	if isExist(fname) && !*releaseNocache {
		return ioutil.ReadFile(fname)
	}

	if err := os.MkdirAll(cachedir, 0775); err != nil {
		return nil, err
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

func runRelease(platform appleopensource.Platform, version string) error {
	if !*releaseQuite {
		fmt.Printf("Release version: %s\n", version)
	}

	release, err := indexRelease(platform, version)
	if err != nil {
		cli.Fatalf(err.Error())
	}

	list, err := appleopensource.ListRelease(release)
	if err != nil {
		cli.Fatalf(err.Error())
	}

	var buf bytes.Buffer
	tbuf := tabwriter.NewWriter(&buf, 2, 1, 2, ' ', 0)

	for _, b := range list {
		if !*releaseQuite {
			if b.Updated {
				tbuf.Write([]byte("\u2022 ")) // u2022: •
			} else {
				tbuf.Write([]byte("  "))
			}
		}
		tbuf.Write([]byte(fmt.Sprintf("%s\t%s", b.Name, b.Version)))
		if !*releaseQuite {
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

func runReleaseMacOS(ctx *cli.ParseContext) error {
	return runRelease(appleopensource.MacOS, *releaseMacOSVersion)
}

func runReleaseXcode(ctx *cli.ParseContext) error {
	return runRelease(appleopensource.Xcode, *releaseXcodeVersion)
}

func runReleaseIOS(ctx *cli.ParseContext) error {
	return runRelease(appleopensource.IOS, *releaseIOSVersion)
}

func runReleaseServer(ctx *cli.ParseContext) error {
	return runRelease(appleopensource.Server, *releaseServerVersion)
}
