// Copyright 2020 The appleopensource Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"go-darwin.dev/appleopensource/pkg/appleopensource"
)

type cache struct {
	ioStreams *IOStreams

	listSource   bool
	listTarballs bool
}

// newCmdCache creates the cache command.
func (a *aos) newCmdCache(ctx context.Context, ioStreams *IOStreams) *cobra.Command {
	cache := &cache{
		ioStreams: ioStreams,
	}

	cmd := &cobra.Command{
		Use:   "cache product [version...] dist",
		Short: "Manage the cache",
	}

	cmd.AddCommand(cache.cmdList(ctx))
	cmd.AddCommand(cache.cmdDelete(ctx))

	return cmd
}

func (c *cache) cmdList(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List cache project",
		RunE: func(*cobra.Command, []string) error {
			return c.runList(ctx)
		},
	}
	f := cmd.Flags()
	f.BoolVarP(&c.listSource, "source", "s", false, "List the source resources type cache")
	f.BoolVarP(&c.listTarballs, "tarballs", "t", false, "List the tarballs resources type cache")

	return cmd
}

func (c *cache) runList(ctx context.Context) error {
	dir := cacheDir()
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		return fmt.Errorf("Not exists cache")
	}

	var typ appleopensource.ResourceType
	switch {
	case c.listTarballs:
		typ = appleopensource.TarballsResource
	case c.listSource:
		typ = appleopensource.SourceResource
	case c.listTarballs && c.listSource:
		return errors.New("-source and -tarballs flags are must be one")
	}

	files, err := ioutil.ReadDir(filepath.Join(dir, typ.String()))
	if err != nil {
		return errors.Wrapf(err, "Not exists the %s type cache", typ.String())
	}

	var buf bytes.Buffer
	for _, f := range files {
		buf.WriteString(strings.TrimSuffix(f.Name(), ".html") + "\n")
	}

	fmt.Printf(buf.String())

	return nil
}

func (c *cache) cmdDelete(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "delete",
		Short: "Delete cache",
		RunE: func(*cobra.Command, []string) error {
			return c.runDelete(ctx)
		},
	}
}

func (c *cache) runDelete(ctx context.Context) error {
	dir := cacheDir()
	if _, err := os.Stat(dir); err == nil {
		log.Printf("Delete %s cache", dir)
		return os.RemoveAll(dir)
	}

	return fmt.Errorf("Not exists cache")
}

// cacheDir create appleopensource cache directory into cacheHome, and return the cache directory path.
func cacheDir() string {
	rootCacheDir := os.Getenv("APPLEOPENSOURCE_CACHE_DIR")
	if rootCacheDir == "" {
		cacheDir, _ := os.UserCacheDir()
		rootCacheDir = filepath.Join(cacheDir, "appleopensource")
	}

	if fi, err := os.Stat(rootCacheDir); err == nil && fi.IsDir() {
		if err := os.MkdirAll(rootCacheDir, 0700); err != nil {
			return ""
		}
	}

	return rootCacheDir
}
