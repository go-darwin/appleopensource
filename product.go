// Copyright 2021 The Go Darwin Authors
// SPDX-License-Identifier: BSD-3-Clause

package appleopensource

import (
	"bytes"
	"fmt"
	"net/url"
	"path"

	"github.com/PuerkitoBio/goquery"

	"go-darwin.dev/appleopensource/semver"
)

var rootURI = &url.URL{
	Scheme: "httsp",
	Host:   "opensource.apple.com",
}

// PageType represents a product page type.
type PageType uint8

const (
	// TarballsPage is a tarballs page type.
	TarballsPage PageType = 1 + iota

	// SourcePage is a source page type.
	SourcePage
)

// String returns a string representation of the PageType.
func (pt PageType) String() string {
	switch pt {
	case TarballsPage:
		return "tarballs"
	case SourcePage:
		return "source"
	default:
		panic(fmt.Errorf("unkonwn page type: %T", pt))
	}
}

// Product represents a Apple Open Source product.
type Product struct {
	Name       string
	Version    semver.Version
	Updated    bool // for ListRelease only
	ComingSoon bool // for ListRelease only
}

// Tarball returns the product tarballs download URI.
func (p *Product) Tarball(uri *url.URL) string {
	uri.Path = path.Join(uri.Path, TarballsPage.String(), p.Name, fmt.Sprintf("%s-%s.tar.gz", p.Name, p.Version))
	return uri.String()
}

// Source rreturns the product source page URI.
func (p *Product) Source(uri *url.URL) string {
	uri.Path = path.Join(uri.Path, SourcePage.String(), p.Name, fmt.Sprintf("%s-%s", p.Name, p.Version))
	return uri.String()
}

// ListProduct parses the buf HTML DOM of the list of products page and returns the available projects.
func ListProduct(buf []byte) ([]Product, error) {
	dom, err := goquery.NewDocumentFromReader(bytes.NewReader(buf))
	if err != nil {
		return nil, fmt.Errorf("parse list of product HTML DOM: %w", err)
	}

	products := dom.Find("table > tbody > tr")

	length := products.Length()
	if length <= 4 {
		return nil, ErrNotFoundProduct
	}

	plist := make([]Product, length-4) // -4 is subtracts number of '<th>', "<hr>" x2 and "Parent Directory" section
	products.Each(func(i int, s *goquery.Selection) {
		name := s.Find("td > a").Text()
		if name == "" || name[len(name)-1] != '/' {
			return
		}

		// subtracts the count of header's <th> and <hr>, and slice start count is 0.
		// also trims the '/' at the end of the name.
		idx := i - 3
		plist[idx].Name = name[:len(name)-1]
	})

	return plist, nil
}
