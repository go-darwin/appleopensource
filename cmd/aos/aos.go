// Copyright 2020 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"os"

	"github.com/zchee/appleopensource/cmd/aos/cmd"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := cmd.NewCommand(ctx, os.Args[1:])
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
