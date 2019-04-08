// Copyright 2017 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"os"
	"path/filepath"

	xdgbasedir "github.com/zchee/go-xdgbasedir"

	"github.com/zchee/appleopensource/pkg/fs"
)

var rootCacheDir = os.Getenv("APPLEOPENSOURCE_CACHE_DIR")

// cacheDir create appleopensource cache directory into cacheHome, and return the cache directory path.
func cacheDir() string {
	if rootCacheDir == "" {
		rootCacheDir = filepath.Join(xdgbasedir.CacheHome(), "appleopensource")
	}
	if err := fs.MkdirAll(rootCacheDir, 0700); err != nil {
		log.Fatal(err)
	}

	return rootCacheDir
}
