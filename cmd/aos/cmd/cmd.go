// Copyright 2020 The appleopensource Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	// version will be increased when upgrading release version.
	tag = "dev"
	// gitCommit will be the hash that the binary was built from and will be populated by the Makefile
	gitCommit = ""

	version = tag + "@" + gitCommit
)

var (
	// AppName returns the application name.
	AppName = filepath.Base(os.Args[0])
)

// aos represents a root command options.
type aos struct {
	noCache    bool
	debug      bool
	configPath string

	ioStreams *IOStreams
}

// NewCommand creates the aos root command.
func NewCommand(ctx context.Context, args []string) *cobra.Command {
	a := &aos{}
	cmd := &cobra.Command{
		Use:                AppName,
		Short:              "An opensource.apple.com resource management tool.",
		SilenceUsage:       false,
		PersistentPreRunE:  func(*cobra.Command, []string) error { return initProfiling() },
		PersistentPostRunE: func(*cobra.Command, []string) error { return flushProfiling() },
		Version:            version,
	}
	cmd.Flags().BoolP("version", "v", false, "Show "+AppName+" version.") // version flag is root only

	f := cmd.PersistentFlags()
	addGlobalFlags(f, a)
	f.Parse(args)

	a.ioStreams = &IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	cmd.SetIn(a.ioStreams.In)
	cmd.SetOut(a.ioStreams.Out)
	cmd.SetErr(a.ioStreams.ErrOut)

	cmd.AddCommand(a.newCmdCache(ctx, a.ioStreams))
	cmd.AddCommand(a.newCmdFetch(ctx, a.ioStreams))
	cmd.AddCommand(a.newCmdList(ctx, a.ioStreams))
	cmd.AddCommand(a.newCmdRelease(ctx, a.ioStreams))
	cmd.AddCommand(a.newCmdVersions(ctx, a.ioStreams))
	cmd.AddCommand(a.newCompletion(ctx, a.ioStreams))

	return cmd
}

const (
	exactArgs = iota
	minArgs
	maxArgs
)

func checkArgs(cmdName string, flags *pflag.FlagSet, expected, checkType int, args ...string) error {
	switch checkType {
	case exactArgs:
		if flags.NArg() != expected {
			return fmt.Errorf("%s: %q requires exactly %d argument(s), args: <%s>\n", AppName, cmdName, expected, strings.Join(args, " "))
		}
	case minArgs:
		if flags.NArg() < expected {
			return fmt.Errorf("%s: %q requires a minimum of %d argument(s), args: <%s>\n", AppName, cmdName, expected, strings.Join(args, " "))
		}
	case maxArgs:
		if flags.NArg() > expected {
			return fmt.Errorf("%s: %q requires a maximum of %d argument(s), args: <%s>\n", AppName, cmdName, expected, strings.Join(args, " "))
		}
	}

	return nil
}
