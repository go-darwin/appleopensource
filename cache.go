// Copyright 2016 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package appleopensource

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	base         = "https://opensource.apple.com"
	cacheDirName = "appleopensource"
)

// isExist returns whether the filename is exists.
func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// xdgCacheHome return the XDG_CACHE_HOME env or "~/.cache".
func xdgCacheHome() string {
	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		u, _ := user.Current()
		cacheHome = filepath.Join(u.HomeDir, ".cache")
	}

	return cacheHome
}

// cacheDir create appleopensource cache directory into cacheHome, and
// return the cache directory path.
func cacheDir() (string, error) {
	cacheDir := os.Getenv("APPLEOPENSOURCE_CACHE_DIR")
	if cacheDir == "" {
		cacheDir = filepath.Join(xdgCacheHome(), cacheDirName)
	}

	if !isExist(cacheDir) {
		err := os.MkdirAll(cacheDir, 0775)
		if err != nil {
			return "", err
		}
	}

	return cacheDir, nil
}

// cacheFile caches the opensource.apple.com website HTML DOM into cacheDir.
func cacheFile(typ string) (string, error) {
	cacheDir, err := cacheDir()
	if err != nil {
		return "", err
	}

	baseURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	baseURL.Path = path.Join(baseURL.Path, typ)
	host := strings.Replace(baseURL.Host, ".", "_", -1)

	fname := filepath.Join(cacheDir, fmt.Sprintf("./%s_%s.html", host, typ))
	if isExist(fname) {
		return fname, nil
	}

	targetURL := baseURL.String()
	dom, err := goquery.NewDocument(targetURL)
	if err != nil {
		return "", err
	}

	buf, err := dom.Find("body").Html()
	if err != nil {
		return "", err
	}

	if err := ioutil.WriteFile(fname, []byte(buf), 0664); err != nil {
		return "", err
	}

	return fname, nil
}
