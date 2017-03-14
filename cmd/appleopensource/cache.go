// Copyright 2017 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/pkgutil/osutil"
	"github.com/urfave/cli"
	"github.com/zchee/appleopensource"
)

var cacheCommand = cli.Command{
	Name:  "cache",
	Usage: "Manage the cache",
	Subcommands: []cli.Command{
		{
			Name:  "list",
			Usage: "List cache project",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "tarballs, t",
					Usage: "List the tarballs resources type cache",
				},
				cli.BoolFlag{
					Name:  "source, s",
					Usage: "List the source resources type cache",
				},
			},
			Before: initCache,
			Action: runCacheList,
		},
		{
			Name:   "delete",
			Usage:  "Delete cache",
			Action: runCacheDelete,
		},
	},
}

var (
	cacheListSource   bool
	cacheListTarballs bool
)

func initCache(ctx *cli.Context) error {
	cacheListSource = ctx.Bool("source")
	cacheListTarballs = ctx.Bool("tarballs")
	return nil
}

func runCacheList(ctx *cli.Context) error {
	dir := cacheDir()
	if osutil.IsNotExist(dir) {
		return fmt.Errorf("Not exists cache")
	}

	typ := appleopensource.TarballsResource
	switch {
	case cacheListTarballs:
		// nothing to do
	case cacheListSource:
		typ = appleopensource.SourceResource
	}

	files, err := ioutil.ReadDir(filepath.Join(dir, typ.String()))
	if err != nil {
		return errors.Wrapf(err, "Not exists the %s type cache", typ.String())
	}

	var buf bytes.Buffer
	for _, f := range files {
		buf.WriteString(strings.TrimSuffix(f.Name(), ".html"))
		buf.WriteString("\n")
	}

	fmt.Printf(buf.String())

	return nil
}

func runCacheDelete(ctx *cli.Context) error {
	if dir := cacheDir(); osutil.IsExist(dir) {
		log.Printf("Delete %s cache", dir)
		return os.RemoveAll(dir)
	}

	return fmt.Errorf("Not exists cache")
}
