// Copyright 2020 The appleopensource Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/spf13/pflag"
)

func addGlobalFlags(flags *pflag.FlagSet, opts *Options) {
	flags.BoolVarP(&opts.debug, "debug", "d", false, "Use debug output")
	flags.StringVarP(&opts.configPath, "config", "c", "", "config file path")

	addProfilingFlags(flags)
}
