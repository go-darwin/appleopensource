// Copyright 2020 The appleopensource Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/zchee/appleopensource/pkg/appleopensource"
)

type fetch struct {
	ioStreams *IOStreams

	product  string
	versions []string
	dist     string
}

// newCmdList creates the list command.
func (a *aos) newCmdFetch(ctx context.Context, ioStreams *IOStreams) *cobra.Command {
	fetch := &fetch{
		ioStreams: ioStreams,
	}

	cmd := &cobra.Command{
		Use:   "fetch product [version...] dist",
		Short: "Fetch the tarballs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := checkArgs(cmd.Name(), cmd.Flags(), 3, minArgs, args...); err != nil {
				return err
			}

			fetch.product = args[0]
			fetch.versions = args[1 : len(args)-1]
			fetch.dist = args[len(args)-1]
			return fetch.run(ctx)
		},
	}

	return cmd
}

func (f *fetch) run(ctx context.Context) error {
	list := make([]string, len(f.versions))
	for i, v := range f.versions {
		p := appleopensource.Product{
			Name:    f.product,
			Version: v,
		}
		list[i] = p.Tarball()
	}

	return appleopensource.Fetch(ctx, f.dist, list...)
}
