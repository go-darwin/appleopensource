// Copyright 2016 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"os/user"
	"path/filepath"
)

// isExist returns whether the filename is exists.
func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// xdgCacheHome return the $XDG_CACHE_HOME env path or "~/.cache".
func xdgCacheHome() string {
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		u, _ := user.Current()
		cacheHome = filepath.Join(u.HomeDir, ".cache")
	}

	return cacheHome
}
