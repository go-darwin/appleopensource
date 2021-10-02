// Copyright 2021 The Go Darwin Authors
// SPDX-License-Identifier: BSD-3-Clause

package appleopensource

import "errors"

var (
	// ErrNotFoundVersion is the not found any product version error.
	ErrNotFoundVersion = errors.New("not found any version")

	// ErrNotFoundProduct is the not found any product error.
	ErrNotFoundProduct = errors.New("not found any product")
)
