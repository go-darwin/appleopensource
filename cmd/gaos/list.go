// Copyright 2016 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/alecthomas/kingpin"
	"github.com/zchee/appleopensource"
)

var (
	listArg = cmdList.Arg("mode", "resource type. (default: tarballs) [tarballs source]").Default("tarballs").String()
)

func init() {
	cmdList.Action(runList)
}

func runList(ctx *kingpin.ParseContext) error {
	mode := appleopensource.ListMode(*listArg)
	list, err := appleopensource.ListPackage(mode)
	if err != nil {
		return err
	}

	fmt.Println(list)
	fmt.Println(len(list))

	return nil
}
