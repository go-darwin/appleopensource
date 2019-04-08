// Copyright 2019 The appleopensource Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

var (
	// version will be increased when upgrading release version.
	tag = "dev"
	// gitCommit will be the hash that the binary was built from and will be populated by the Makefile
	gitCommit = ""

	version = tag + "@" + gitCommit
)
