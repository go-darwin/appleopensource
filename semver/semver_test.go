// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package semver_test

import (
	"math/rand"
	"sort"
	"strings"
	"testing"

	"go-darwin.dev/appleopensource/semver"
)

var tests = []struct {
	in  string
	out string
}{
	{"bad", ""},
	{"1-alpha.beta.gamma", ""},
	{"1-pre", ""},
	{"1+meta", ""},
	{"1-pre+meta", ""},
	{"1.2-pre", ""},
	{"1.2+meta", ""},
	{"1.2-pre+meta", ""},
	{"1.0.0-alpha", "1.0.0-alpha"},
	{"1.0.0-alpha.1", "1.0.0-alpha.1"},
	{"1.0.0-alpha.beta", "1.0.0-alpha.beta"},
	{"1.0.0-beta", "1.0.0-beta"},
	{"1.0.0-beta.2", "1.0.0-beta.2"},
	{"1.0.0-beta.11", "1.0.0-beta.11"},
	{"1.0.0-rc.1", "1.0.0-rc.1"},
	{"1", "1.0.0"},
	{"1.0", "1.0.0"},
	{"1.0.0", "1.0.0"},
	{"1.2", "1.2.0"},
	{"1.2.0", "1.2.0"},
	{"1.2.3-456", "1.2.3-456"},
	{"1.2.3-456.789", "1.2.3-456.789"},
	{"1.2.3-456-789", "1.2.3-456-789"},
	{"1.2.3-456a", "1.2.3-456a"},
	{"1.2.3-pre", "1.2.3-pre"},
	{"1.2.3-pre+meta", "1.2.3-pre"},
	{"1.2.3-pre.1", "1.2.3-pre.1"},
	{"1.2.3-zzz", "1.2.3-zzz"},
	{"1.2.3", "1.2.3"},
	{"1.2.3+meta", "1.2.3"},
	{"1.2.3+meta-pre", "1.2.3"},
	{"1.2.3+meta-pre.sha.256a", "1.2.3"},
}

func TestIsValid(t *testing.T) {
	for _, tt := range tests {
		ok := semver.IsValid(tt.in)
		if ok != (tt.out != "") {
			t.Errorf("IsValid(%q) = %v, want %v", tt.in, ok, !ok)
		}
	}
}

func TestCanonical(t *testing.T) {
	for _, tt := range tests {
		out := semver.Canonical(tt.in)
		if out != tt.out {
			t.Errorf("Canonical(%q) = %q, want %q", tt.in, out, tt.out)
		}
	}
}

func TestMajor(t *testing.T) {
	for _, tt := range tests {
		out := semver.Major(tt.in)
		want := ""
		if i := strings.Index(tt.out, "."); i >= 0 {
			want = tt.out[:i]
		}
		if out != want {
			t.Errorf("Major(%q) = %q, want %q", tt.in, out, want)
		}
	}
}

func TestMajorMinor(t *testing.T) {
	for _, tt := range tests {
		out := semver.MajorMinor(tt.in)
		var want string
		if tt.out != "" {
			want = tt.in
			if i := strings.Index(want, "+"); i >= 0 {
				want = want[:i]
			}
			if i := strings.Index(want, "-"); i >= 0 {
				want = want[:i]
			}
			switch strings.Count(want, ".") {
			case 0:
				want += ".0"
			case 1:
				// ok
			case 2:
				want = want[:strings.LastIndex(want, ".")]
			}
		}
		if out != want {
			t.Errorf("MajorMinor(%q) = %q, want %q", tt.in, out, want)
		}
	}
}

func TestPrerelease(t *testing.T) {
	for _, tt := range tests {
		pre := semver.Prerelease(tt.in)
		var want string
		if tt.out != "" {
			if i := strings.Index(tt.out, "-"); i >= 0 {
				want = tt.out[i:]
			}
		}
		if pre != want {
			t.Errorf("Prerelease(%q) = %q, want %q", tt.in, pre, want)
		}
	}
}

func TestBuild(t *testing.T) {
	for _, tt := range tests {
		build := semver.Build(tt.in)
		var want string
		if tt.out != "" {
			if i := strings.Index(tt.in, "+"); i >= 0 {
				want = tt.in[i:]
			}
		}
		if build != want {
			t.Errorf("Build(%q) = %q, want %q", tt.in, build, want)
		}
	}
}

func TestCompare(t *testing.T) {
	for i, ti := range tests {
		for j, tj := range tests {
			cmp := semver.Compare(ti.in, tj.in)
			var want int
			if ti.out == tj.out {
				want = 0
			} else if i < j {
				want = -1
			} else {
				want = +1
			}
			if cmp != want {
				t.Errorf("Compare(%q, %q) = %d, want %d", ti.in, tj.in, cmp, want)
			}
		}
	}
}

func TestSort(t *testing.T) {
	versions := make([]string, len(tests))
	for i, test := range tests {
		versions[i] = test.in
	}
	rand.Shuffle(len(versions), func(i, j int) { versions[i], versions[j] = versions[j], versions[i] })
	semver.Sort(versions)
	if !sort.IsSorted(semver.ByVersion(versions)) {
		t.Errorf("list is not sorted:\n%s", strings.Join(versions, "\n"))
	}
}

var (
	v1 = "1.0.0+metadata-dash"
	v2 = "1.0.0+metadata-dash1"
)

func BenchmarkCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if semver.Compare(v1, v2) != 0 {
			b.Fatalf("bad compare")
		}
	}
}
