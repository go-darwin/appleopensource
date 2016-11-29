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

	cli "github.com/alecthomas/kingpin"
	"github.com/pkg/errors"
	"github.com/zchee/appleopensource"
)

var (
	cmdCache = cli.Command("cache", "Manage gaos cache.")

	// list
	cmdCacheList      = cmdCache.Command("list", "List cached project.")
	cacheListTarballs = cmdCacheList.Flag("tarballs", "List tarballs type cache.").Short('t').Bool()
	cacheListSource   = cmdCacheList.Flag("source", "List source type cache.").Short('s').Bool()

	// delete
	cmdCacheDelete = cmdCache.Command("delete", "Delete cache.")
)

func init() {
	cmdCacheList.Action(runCacheList)
	cmdCacheDelete.Action(runCacheDelete)
}

func runCacheList(ctx *cli.ParseContext) error {
	dir := cacheDir()
	if !isExist(dir) {
		return fmt.Errorf("Not exists cache")
	}

	typ := appleopensource.TypeTarballs
	switch {
	case *cacheListSource:
		typ = appleopensource.TypeSource
	case *cacheListTarballs:
		// nothing to do
	}

	files, err := ioutil.ReadDir(filepath.Join(dir, typ.String()))
	if err != nil {
		cli.Fatalf(errors.Wrapf(err, "Not exists the %s type cache", typ.String()).Error())
	}

	var buf bytes.Buffer
	for _, f := range files {
		buf.WriteString(strings.TrimSuffix(f.Name(), ".html"))
		buf.WriteString("\n")
	}

	fmt.Printf(buf.String())

	return nil
}

func runCacheDelete(ctx *cli.ParseContext) error {
	if dir := cacheDir(); isExist(dir) {
		return os.RemoveAll(dir)
	}

	return fmt.Errorf("Not exists cache")
}

// cacheDir create appleopensource cache directory into cacheHome, and
// return the cache directory path.
func cacheDir() string {
	dir := os.Getenv("APPLEOPENSOURCE_CACHE_DIR")
	if dir == "" {
		dir = filepath.Join(xdgCacheHome(), "appleopensource")
	}

	return dir
}
