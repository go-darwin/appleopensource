// Copyright 2016 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package appleopensource

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/blang/semver"
)

const rootURL = "https://opensource.apple.com"

// ResourceType represents a list mode.
type ResourceType int

const (
	// TypeTarballs tarballs type.
	TypeTarballs ResourceType = 1 << iota

	// TypeSource source type.
	TypeSource
)

// String implements a fmt.Stringer interface.
func (r ResourceType) String() string {
	switch r {
	case TypeTarballs:
		return "tarballs"
	case TypeSource:
		return "source"
	default:
		return ""
	}
}

// Project represents a apple open source project.
type Project struct {
	Name       string
	Version    string
	Updated    bool // for release only
	ComingSoon bool // for release only
}

func (p *Project) Tarball() string {
	return fmt.Sprintf("%s/%s/%s/%s-%s.tar.gz", rootURL, TypeTarballs, p.Name, p.Name, p.Version)
}

func (p *Project) Source() string {
	return fmt.Sprintf("%s/%s/%s/%s-%s/", rootURL, TypeSource, p.Name, p.Name, p.Version)
}

func index(baseURL *url.URL) ([]byte, error) {
	dom, err := goquery.NewDocument(baseURL.String())
	if err != nil {
		return nil, err
	}

	table, err := dom.Find("body #content > div.column").Html()
	if err != nil {
		return nil, err
	}

	// 6 is 404 not found
	if len(strings.TrimSpace(table)) == 0 {
		return nil, fmt.Errorf("Not found %s page", baseURL.String())
	}

	return bytes.TrimSpace([]byte(table)), nil
}

// IndexProject return the index of opensource.apple.com/<typ> HTML DOM tree.
func IndexProject(typ string) ([]byte, error) {
	var rtype string

	switch typ {
	case TypeTarballs.String():
		rtype = TypeTarballs.String()
	case TypeSource.String():
		rtype = TypeSource.String()
	default:
		return nil, errors.New("unknown resource type")
	}

	baseURL, err := url.Parse(rootURL)
	if err != nil {
		return nil, err
	}
	baseURL.Path = path.Join(baseURL.Path, fmt.Sprint(rtype))

	return index(baseURL)
}

// ListProject parses the buf HTML DOM tree, and return the project list.
func ListProject(buf []byte) ([]Project, error) {
	r := bytes.NewReader(buf)
	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	projects := dom.Find("table > tbody > tr")

	// Subtracts the number of <th>, <hr> and "Parent Directory"
	list := make([]Project, projects.Length()-4)

	projects.Each(func(i int, s *goquery.Selection) {
		if name := s.Find("td > a").Text(); name != "" && name[len(name)-1] == byte('/') {
			// Subtracts the count of header's <th> and <hr>, and slice start count is 0.
			// Also trims the "/" at the end of the name.
			list[i-3].Name = name[:len(name)-1]
		}
	})

	return list, nil
}

// IndexVersion return the index of all versions of the project HTML DOM tree.
func IndexVersion(project, typ string) ([]byte, error) {
	var rtype string

	switch typ {
	case TypeTarballs.String():
		rtype = TypeTarballs.String()
	case TypeSource.String():
		rtype = TypeSource.String()
	default:
		return nil, errors.New("unknown resource type")
	}

	baseURL, err := url.Parse(rootURL)
	if err != nil {
		return nil, err
	}
	baseURL.Path = path.Join(baseURL.Path, fmt.Sprint(rtype), project)

	return index(baseURL)
}

var zero = ".0.0"

// ListVersions parses the buf HTML DOM tree, and return the available versions of the project.
func ListVersions(buf []byte) ([]string, error) {
	r := bytes.NewReader(buf)
	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	versions := dom.Find("table > tbody > tr")

	// Subtracts the number of <th>, <hr> and "Parent Directory"
	vlist := make([]semver.Version, versions.Length()-4)

	var start int
	versions.Each(func(i int, s *goquery.Selection) {
		if name := s.Find("td > a").Text(); name != "" && name != "Parent Directory" {
			if start == 0 {
				start = strings.Index(name, "-") + 1
			}

			var end int
			switch {
			case name[len(name)-1] == '/':
				// Trims the "/" at the end of the name.
				end = len(name) - 1
			case name[len(name)-7] == '.':
				// Trims the ".tar.gz" at the end of the name.
				end = len(name) - 7
			}

			// Subtracts the count of header's <th> and <hr>, and slice start count is 0.
			// Ignored parse miss
			vlist[i-3], _ = semver.ParseTolerant(name[start:end])
		}
	})
	semver.Sort(vlist)

	list := make([]string, len(vlist))
	for i, v := range vlist {
		// Try trims the ".0.0"
		list[i] = strings.TrimSuffix(v.String(), zero)
	}
	return list, nil
}

// IndexRelease return the index of all releases of the platform.
func IndexRelease(platform Platform, version string) ([]byte, error) {
	var prefix string

	switch platform {
	case MacOS:
		v, err := semver.ParseTolerant(version)
		if err != nil {
			return nil, err
		}
		// wtf why does not use unified url?
		threshold, err := semver.ParseTolerant("10.9")
		if err != nil {
			return nil, err
		}
		switch v.Compare(threshold) {
		case 0, 1: // 0 is equal, 1 is greater than threshold
			prefix = "os-x"
		case -1: // -1 is less than threshold
			prefix = "mac-os-x"
		}
	case Xcode:
		prefix = "developer-tools"
	case IOS:
		prefix = "ios"
	case Server:
		prefix = "os-x-server"
	default:
		return nil, errors.New("unknown platform")
	}

	baseURL, err := url.Parse(rootURL)
	if err != nil {
		return nil, err
	}
	baseURL.Path = path.Join(baseURL.Path, "release", fmt.Sprintf("%s-%s.html", prefix, strings.Replace(version, ".", "", -1)))

	return index(baseURL)
}

const ComingSoon = "(coming soon!)"

// ListRelease parses the release buf HTML DOM, and return the Project slice.
func ListRelease(buf []byte) ([]Project, error) {
	r := bytes.NewReader(buf)
	dom, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	release := dom.Find("table > tbody > tr")

	projects := make([]Project, release.Length())

	release.Each(func(i int, s *goquery.Selection) {
		// td.project-name is e.g. "xnu-3789.1.32"
		if data := strings.TrimSpace(s.Find("td.project-name").Text()); data != "" {
			p := strings.Split(data, "-")

			projects[i].Name = p[0]
			version := p[1]

			if strings.Contains(version, ComingSoon) {
				version = strings.TrimSpace(strings.Replace(version, ComingSoon, "", -1))
				projects[i].ComingSoon = true
			}

			projects[i].Version = version
		}
		if updated := strings.TrimSpace(s.Find("td.project-updated").Text()); updated != "" {
			projects[i].Updated = true
		}
	})

	return projects, nil
}
