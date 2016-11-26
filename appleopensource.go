// Copyright 2016 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package appleopensource

import (
	"io/ioutil"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ListMode represets a list mode.
type ListMode string

const (
	// ModeTarballs tarballs mode.
	ModeTarballs ListMode = "tarballs"
	// ModeSource source mode.
	ModeSource = "source"
)

// ListTarballs return the lists of opensource.apple.com/tarballs
func ListTarballs() ([]string, error) {
	return ListPackage(ModeTarballs)
}

// ListSource return the lists of opensource.apple.com/source.
func ListSource() ([]string, error) {
	return ListPackage(ModeSource)
}

// ListPackage base function that return the opensource.apple.com package list.
func ListPackage(typ ListMode) ([]string, error) {
	cacheFile, err := cacheFile(string(typ))
	if err != nil {
		return nil, err
	}

	f, err := ioutil.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}

	sr := strings.NewReader(string(f))
	dom, err := goquery.NewDocumentFromReader(sr)
	if err != nil {
		return nil, err
	}

	packages := dom.Find("#content > div.column > table > tbody")

	// Subtracts the number of <th>, <hr> and "Parent Directory"
	list := make([]string, packages.Children().Length()-4)

	packages.Children().Each(func(_ int, s *goquery.Selection) {
		if name := s.Find("td a").Text(); name != "" && name[len(name)-1] == byte('/') {
			// Subtracts the count of header's <th> and <hr>, and slice start count is 0.
			// Also trims the "/" at the end of the name.
			list[s.Index()-3] = name[:len(name)-1]
		}
	})

	return list, nil
}
