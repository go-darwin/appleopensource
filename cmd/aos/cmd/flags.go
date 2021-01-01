// Copyright 2020 The appleopensource Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/spf13/pflag"
)

func addGlobalFlags(flags *pflag.FlagSet, a *aos) {
	flags.BoolVar(&a.noCache, "no-cache", false, "Do not use cache")
	flags.BoolVarP(&a.debug, "debug", "d", false, "Use debug output")
	flags.StringVarP(&a.configPath, "config", "c", "", "config file path")

	addProfilingFlags(flags)
}
