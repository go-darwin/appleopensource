// Copyright 2021 The Go Darwin Authors
// SPDX-License-Identifier: BSD-3-Clause

package appleopensource

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"go-darwin.dev/appleopensource/semver"
)

// ListVersion parses the buf HTML DOM of the list of product versions page and returns the available product versions.
func ListVersion(buf []byte) ([]semver.Version, error) {
	dom, err := goquery.NewDocumentFromReader(bytes.NewReader(buf))
	if err != nil {
		return nil, fmt.Errorf("parse list of product versions HTML DOM: %w", err)
	}

	versions := dom.Find("table > tbody > tr")

	length := versions.Length()
	if length <= 4 {
		return nil, ErrNotFoundVersion
	}

	vlist := make([]semver.Version, length-4) // -4 is subtracts number of '<th>', "<hr>" x2 and "Parent Directory" section
	versions.Each(func(i int, s *goquery.Selection) {
		name := s.Find("td > a").Text()
		if name == "" || name == "Parent Directory" {
			return
		}

		var end int
		switch {
		case strings.HasSuffix(name, "/"):
			end = len(name) - 1
		case name[len(name)-7] == '.': // TODO(zchee): comment why -7
			end = len(name) - 7
		}

		start := strings.Index(name, "-") + 1 // TODO(zchee): comment why +1
		idx := i - 3                          // TODO(zchee): comment why -3
		vlist[idx] = trimZeros(semver.Canonical(semver.Version(name[start:end])))
	})
	semver.Sort(vlist)

	return vlist, nil
}

// trimZeros trims ".0" as much as possible.
func trimZeros(version semver.Version) semver.Version {
	for {
		if strings.HasSuffix(string(version), ".0") {
			version = version[:len(version)-2]
			continue
		}
		return version
	}
}
