// Copyright 2020 The appleopensource Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"go-darwin.dev/appleopensource/pkg/appleopensource"
)

type list struct {
	*aos

	ioStreams *IOStreams

	cacheDir string
	source   bool
	tarballs bool
}

// newCmdList creates the list command.
func (a *aos) newCmdList(ctx context.Context, ioStreams *IOStreams) *cobra.Command {
	list := &list{
		aos:       a,
		ioStreams: ioStreams,
		cacheDir:  filepath.Join(cacheDir(), "list"),
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all project available to opensource.apple.com.",
		PreRunE: func(*cobra.Command, []string) error {
			if _, err := os.Stat(list.cacheDir); err != nil && errors.Is(err, os.ErrNotExist) {
				return os.MkdirAll(list.cacheDir, 0755)
			}
			return nil
		},
		RunE: func(*cobra.Command, []string) error { return list.run(ctx) },
	}

	f := cmd.Flags()
	f.BoolVar(&list.source, "source", false, "List the source resources")
	f.BoolVar(&list.tarballs, "tarballs", false, "List the tarballs resources")

	return cmd
}

// index return the opensource.apple.com project index, and caches the HTML DOM tree into cacheDir.
func (l *list) indexList(ctx context.Context, typ appleopensource.ResourceType) ([]byte, error) {
	fname := filepath.Join(l.cacheDir, fmt.Sprintf("%s.html", typ))

	if _, err := os.Stat(fname); err == nil && !l.noCache {
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

func (l *list) run(ctx context.Context) error {
	var mode appleopensource.ResourceType
	switch {
	case l.source && l.tarballs:
		return errors.New("--source and --tarballs flags are must be one")
	case l.source:
		mode = appleopensource.SourceResource
	case l.tarballs:
		mode = appleopensource.TarballsResource
	default:
		mode = appleopensource.TarballsResource // default is tarballs mode
	}

	index, err := l.indexList(ctx, mode)
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
	s := buf.String()[:buf.Len()-2]

	_, err = fmt.Println(s)

	return err
}
