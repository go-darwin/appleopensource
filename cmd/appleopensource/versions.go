// Copyright 2016 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	cli "github.com/alecthomas/kingpin"
	"github.com/zchee/appleopensource"
)

var (
	versionsProject = cmdVersions.Arg("project", "List version of the project name. Available project name is \"list\" command result.").Required().String()

	versionsTarballs = cmdVersions.Flag("tarballs", "List the tarballs resource versions.").Short('t').Bool()
	versionsSource   = cmdVersions.Flag("source", "List the source resource versions.").Short('s').Bool()
	versionsNoCache  = cmdVersions.Flag("no-cache", "Disable the cache.").Short('n').Bool()
)

func runVersions(ctx *cli.ParseContext) error {
	mode := appleopensource.TypeTarballs
	switch {
	case *versionsSource:
		mode = appleopensource.TypeSource
	case *versionsTarballs:
		// nothing to do
	}

	buf, err := indexVersion(*versionsProject, mode.String())
	if err != nil {
		log.Fatal(err)
	}

	list, err := appleopensource.ListVersions(buf)

	fmt.Println(strings.Join(list, "\n"))

	return nil
}

// index return the opensource.apple.com project index, and caches the HTML DOM tree into cacheDir.
func indexVersion(project, typ string) ([]byte, error) {
	cachedir := filepath.Join(cacheDir(), typ)
	fname := filepath.Join(cachedir, fmt.Sprintf("%s.html", project))
	if isExist(fname) && !*versionsNoCache {
		return ioutil.ReadFile(fname)
	}

	if err := os.MkdirAll(cachedir, 0755); err != nil {
		return nil, err
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
