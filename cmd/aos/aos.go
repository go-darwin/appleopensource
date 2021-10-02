// Copyright 2020 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"os"

	"go-darwin.dev/appleopensource/cmd/aos/cmd"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := cmd.NewCommand(ctx, os.Args[1:]).Execute(); err != nil {
		os.Exit(1)
	}
}
