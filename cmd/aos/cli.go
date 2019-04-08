// Copyright 2017 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"
)

const (
	exactArgs = iota
	minArgs
	maxArgs
)

func checkArgs(context *cli.Context, expected, checkType int, args ...string) (err error) {
	cmdName := context.Command.Name
	switch checkType {
	case exactArgs:
		if context.NArg() != expected {
			err = fmt.Errorf("%s: %q requires exactly <%s> %d argument(s)", os.Args[0], cmdName, strings.Join(args, " "), expected)
		}
	case minArgs:
		if context.NArg() < expected {
			err = fmt.Errorf("%s: %q requires a minimum of <%s> %d argument(s)", os.Args[0], cmdName, strings.Join(args, " "), expected)
		}
	case maxArgs:
		if context.NArg() > expected {
			err = fmt.Errorf("%s: %q requires a maximum of <%s> %d argument(s)", os.Args[0], cmdName, strings.Join(args, " "), expected)
		}
	}

	if err != nil {
		fmt.Printf("Incorrect Usage.\n\n")
		return err
	}
	return nil
}
