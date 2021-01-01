// Copyright 2020 The appleopensource Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/zchee/appleopensource/pkg/appleopensource"
)

type versions struct {
	*aos

	ioStreams *IOStreams

	product  string
	source   bool
	tarballs bool
}

// newCmdVersions creates the versions command.
func (a *aos) newCmdVersions(ctx context.Context, ioStreams *IOStreams) *cobra.Command {
	versions := &versions{
		aos:       a,
		ioStreams: ioStreams,
	}

	cmd := &cobra.Command{
		Use:   "versions product",
		Short: "List all versions of the product.",
		RunE: func(_ *cobra.Command, args []string) error {
			versions.product = args[0]
			return versions.runVersions(ctx)
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&versions.source, "source", "s", false, "List the source resources type cache")
	f.BoolVarP(&versions.tarballs, "tarballs", "t", false, "List the tarballs resources type cache")

	return cmd
}

// index return the opensource.apple.com project index, and caches the HTML DOM tree into cacheDir.
func (v *versions) indexVersion(project string, typ appleopensource.ResourceType) ([]byte, error) {
	versionsCachedir := filepath.Join(cacheDir(), typ.String())

	if _, err := os.Stat(versionsCachedir); err != nil && errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(versionsCachedir, 0700); err != nil {
			return nil, err
		}
	}

	fname := filepath.Join(versionsCachedir, fmt.Sprintf("%s.html", project))

	if _, err := os.Stat(fname); err == nil && !v.noCache {
		return ioutil.ReadFile(fname)
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

func (v *versions) runVersions(ctx context.Context) error {
	mode := appleopensource.TarballsResource
	switch {
	case v.tarballs:
		// nothing to do
	case v.source:
		mode = appleopensource.SourceResource
	case v.tarballs && v.source:
		return errors.New("-source and -tarballs flags are must be one")
	}

	buf, err := v.indexVersion(v.product, mode)
	if err != nil {
		return err
	}

	list, err := appleopensource.ListVersions(buf)
	if err != nil {
		return err
	}

	_, err = fmt.Println(strings.Join(list, "\n"))

	return err
}
