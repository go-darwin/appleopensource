// Copyright 2017 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/urfave/cli"

	"github.com/zchee/appleopensource"
	"github.com/zchee/appleopensource/pkg/fs"
)

var listCommand = cli.Command{
	Name:  "list",
	Usage: "List all project available to opensource.apple.com.",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "tarballs, t",
			Usage: "List the tarballs resources",
		},
		cli.BoolFlag{
			Name:  "source, s",
			Usage: "List the source resources",
		},
	},
	Before: initList,
	Action: runList,
}

var (
	listSource   bool
	listTarballs bool
)

var listCachedir = filepath.Join(cacheDir(), "list")

func initList(ctx *cli.Context) error {
	listSource = ctx.Bool("source")
	listTarballs = ctx.Bool("tarballs")

	if err := fs.MkdirAll(listCachedir, 0700); err != nil {
		return err
	}
	return nil
}

// index return the opensource.apple.com project index, and caches the HTML DOM tree into cacheDir.
func indexList(ctx *cli.Context, typ appleopensource.ResourceType) ([]byte, error) {
	fname := filepath.Join(listCachedir, fmt.Sprintf("%s.html", typ))
	if fs.IsExist(fname) && !noCache {
		return ioutil.ReadFile(fname)
	}

	buf, err := appleopensource.IndexProject(typ)
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(fname, buf, 0664); err != nil {
		return nil, err
	}

	return buf, nil
}

func runList(ctx *cli.Context) error {
	mode := appleopensource.TarballsResource
	switch {
	case listTarballs:
		// nothing to do
	case listSource:
		mode = appleopensource.SourceResource
	}

	index, err := indexList(ctx, mode)
	if err != nil {
		return err
	}

	list, err := appleopensource.ListProject(index)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	for _, b := range list {
		buf.WriteString(b.Name + "\n")
	}

	// TODO(zchee): another way of trim last new line
	buf.Truncate(buf.Len() - 1)

	fmt.Printf(buf.String())

	return nil
}
