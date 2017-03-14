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

const root = "https://opensource.apple.com/"

var rootURL, _ = url.Parse(root)

// ResourceType represents a resource type.
type ResourceType int

const (
	// TarballsResource is a tarballs resource type.
	TarballsResource ResourceType = 1 << iota

	// SourceResource is a source resource type.
	SourceResource
)

// String implements a fmt.Stringer interface.
func (r ResourceType) String() string {
	switch r {
	case TarballsResource:
		return "tarballs"
	case SourceResource:
		return "source"
	default:
		return ""
	}
}

func index(u *url.URL) ([]byte, error) {
	dom, err := goquery.NewDocument(u.String())
	if err != nil {
		return nil, err
	}

	table, err := dom.Find("body #content > div.column").Html()
	if err != nil {
		return nil, err
	}

	// 0 is 404 not found
	if len(strings.TrimSpace(table)) == 0 {
		return nil, fmt.Errorf("Not found %s project", u.String())
	}

	return bytes.TrimSpace([]byte(table)), nil
}

// IndexProject return the index of opensource.apple.com/<typ> HTML DOM tree.
func IndexProject(typ ResourceType) ([]byte, error) {
	u := *rootURL // copy
	u.Path = path.Join(u.Path, typ.String())

	return index(&u)
}

// IndexVersion return the index of all versions of the project HTML DOM tree.
func IndexVersion(project string, typ ResourceType) ([]byte, error) {
	u := *rootURL // copy
	u.Path = path.Join(u.Path, typ.String(), project)

	return index(&u)
}

const (
	macOSPrefix  = "macos"
	osxPrefix    = "os-x"
	macOSXPrefix = "mac-os-x"
)

// IndexRelease return the index of projects of the specified platforms release version.
func IndexRelease(platform Platform, version string) ([]byte, error) {
	var prefix string

	switch platform {
	case MacOS:
		v, err := semver.ParseTolerant(version)
		if err != nil {
			return nil, err
		}
		// wtf why does not use unified url?
		// 10.12 ~ newer:  Uses 'macos'
		// 10.11.6 ~ 10.9: Uses 'os-x'
		// 10.9 ~ older:   Uses 'mac-os-x'
		threshold, err := semver.ParseTolerant("10.11.6")
		if err != nil {
			return nil, err
		}
		switch v.Compare(threshold) {
		case 1: // 1 is greater than threshold, 10.12 or newer
			prefix = macOSPrefix
		case 0: // 0 is equal, 10.11.6
			prefix = osxPrefix
		case -1: // -1 is less than threshold, 10.11.5 or older
			secondThreshold, err := semver.ParseTolerant("10.9")
			if err != nil {
				return nil, err
			}
			switch v.Compare(secondThreshold) {
			case 0, 1: // 0 is equal, 1 is greater than threshold
				prefix = osxPrefix
			case -1:
				prefix = macOSXPrefix
			}
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

	u := *rootURL // copy
	u.Path = path.Join(u.Path, "release", fmt.Sprintf("%s-%s.html", prefix, strings.Replace(version, ".", "", -1)))

	return index(&u)
}

// Project represents a Apple open source project.
type Project struct {
	Name       string
	Version    string
	Updated    bool // for release only
	ComingSoon bool // for release only
}

// Tarball return the tarballs resource download uri.
func (p *Project) Tarball() string {
	return root + path.Join(TarballsResource.String(), p.Name, fmt.Sprintf("%s-%s.tar.gz", p.Name, p.Version))
}

// Source return the source resource page uri.
func (p *Project) Source() string {
	return root + path.Join(SourceResource.String(), p.Name, fmt.Sprintf("%s-%s", p.Name, p.Version))
}

// ListProject parses the project list HTML DOM, and return the project list.
func ListProject(buf []byte) ([]Project, error) {
	dom, err := goquery.NewDocumentFromReader(bytes.NewReader(buf))
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

// ListVersions parses the project version index page HTML DOM, and return the available versions of the project.
func ListVersions(buf []byte) ([]string, error) {
	dom, err := goquery.NewDocumentFromReader(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	versions := dom.Find("table > tbody > tr")

	// Subtracts the number of <th>, <hr> and "Parent Directory"
	vlist := make([]semver.Version, versions.Length()-4)

	versions.Each(func(i int, s *goquery.Selection) {
		if name := s.Find("td > a").Text(); name != "" && name != "Parent Directory" {
			start := strings.Index(name, "-") + 1

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
		// Try trims the ".0" suffix
		list[i] = trimZeros(v.String())
	}
	return list, nil
}

// trimZeros trims ".0" as much as possible.
func trimZeros(version string) string {
	for {
		if strings.HasSuffix(version, ".0") {
			version = version[:len(version)-2]
			continue
		}
		return version
	}
}

// ComingSoon is a Apple's comming soon message.
const ComingSoon = "(coming soon!)"

// ListRelease parses the release page HTML DOM, and return the Project slice.
func ListRelease(buf []byte) ([]Project, error) {
	dom, err := goquery.NewDocumentFromReader(bytes.NewReader(buf))
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
