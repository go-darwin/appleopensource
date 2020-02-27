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
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/zchee/appleopensource/pkg/appleopensource"
)

type release struct {
	ioStreams *IOStreams

	version string
	quiet   bool
}

var releaseCachedir = filepath.Join(cacheDir(), "release")

// newCmdList creates the release command.
func newCmdRelease(ctx context.Context, ioStreams *IOStreams) *cobra.Command {
	release := &release{
		ioStreams: ioStreams,
	}

	cmd := &cobra.Command{
		Use:   "release",
		Short: "List all projects included to the releases available to opensource.apple.com.",
	}
	f := cmd.Flags()
	release.version = f.Arg(0)
	f.BoolVarP(&release.quiet, "quiet", "q", false, "suppress some output")

	cmd.AddCommand(release.cmdMacOS(ctx))
	cmd.AddCommand(release.cmdXCode(ctx))
	cmd.AddCommand(release.cmdIOS(ctx))
	cmd.AddCommand(release.cmdServer(ctx))

	return cmd
}

func (r *release) cmdMacOS(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "macos",
		Short: "macOS release",
		RunE: func(*cobra.Command, []string) error {
			return r.runRelease(ctx, appleopensource.MacOS, r.version)
		},
	}
}

func (r *release) cmdXCode(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "xcode",
		Short: "Developer Tool(Xcode) release",
		RunE: func(*cobra.Command, []string) error {
			return r.runRelease(ctx, appleopensource.Xcode, r.version)
		},
	}
}

func (r *release) cmdIOS(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "ios",
		Short: "iOS release",
		RunE: func(*cobra.Command, []string) error {
			return r.runRelease(ctx, appleopensource.IOS, r.version)
		},
	}
}

func (r *release) cmdServer(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "macOS Server release",
		RunE: func(*cobra.Command, []string) error {
			return r.runRelease(ctx, appleopensource.Server, r.version)
		},
	}
}

func (r *release) indexRelease(ctx context.Context, platform appleopensource.Platform, version string) ([]byte, error) {
	fname := filepath.Join(releaseCachedir, fmt.Sprintf("%s_%s.html", platform, strings.Replace(version, ".", "", -1)))
	// if _, err := os.Stat(fname); err == nil && !noCache {
	if _, err := os.Stat(fname); err == nil {
		return ioutil.ReadFile(fname)
	}

	buf, err := appleopensource.IndexRelease(platform, version)
	if err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(fname, buf, 0664); err != nil {
		return nil, err
	}

	return buf, nil
}

func (r *release) runRelease(ctx context.Context, platform appleopensource.Platform, version string) error {
	if !r.quiet {
		fmt.Printf("Release version: %s\n", version)
	}

	release, err := r.indexRelease(ctx, platform, version)
	if err != nil {
		log.Fatal(err.Error())
	}

	list, err := appleopensource.ListRelease(release)
	if err != nil {
		log.Fatal(err.Error())
	}

	var buf bytes.Buffer
	tbuf := tabwriter.NewWriter(&buf, 2, 1, 2, ' ', 0)

	for _, b := range list {
		if !r.quiet {
			if b.Updated {
				tbuf.Write([]byte("\u2022 ")) // u2022: •
			} else {
				tbuf.Write([]byte("  "))
			}
		}
		tbuf.Write([]byte(fmt.Sprintf("%s\t%s", b.Name, b.Version)))
		if !r.quiet {
			tbuf.Write([]byte("\t"))
			if b.ComingSoon {
				tbuf.Write([]byte(appleopensource.ComingSoon))
			}
		}
		tbuf.Write([]byte("\n"))
	}
	tbuf.Flush()

	fmt.Printf(buf.String())

	return nil
}
