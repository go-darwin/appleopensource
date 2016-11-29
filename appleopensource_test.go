// Copyright 2016 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package appleopensource

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

// testListTimestamp date of created the testdata from https://opensource.apple.com
const timestamp = "2016-11-29"

var (
	wantTarballsIndex []byte
	wantSourceIndex   []byte

	wantIndexVersionCsu []byte
	wantIndexVersionXnu []byte

	wantIndexReleaseMacOS []byte
	wantIndexReleaseXcode []byte
)

func TestMain(m *testing.M) {
	wantTarballsIndex = readTestFile(fmt.Sprintf("testdata/index_tarballs_%s.golden", timestamp))
	wantSourceIndex = readTestFile(fmt.Sprintf("testdata/index_source_%s.golden", timestamp))

	wantIndexVersionCsu = readTestFile(fmt.Sprintf("testdata/version_Csu_%s.golden", timestamp))
	wantIndexVersionXnu = readTestFile(fmt.Sprintf("testdata/version_xnu_%s.golden", timestamp))

	wantIndexReleaseMacOS = readTestFile(fmt.Sprintf("testdata/release_macos_1012_%s.golden", timestamp))
	wantIndexReleaseXcode = readTestFile(fmt.Sprintf("testdata/release_xcode_731_%s.golden", timestamp))

	os.Exit(m.Run())
}

func readTestFile(filename string) []byte {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return buf
}

func TestResourceType_String(t *testing.T) {
	tests := []struct {
		name string
		r    ResourceType
		want string
	}{
		{
			name: "tarballs",
			r:    TypeTarballs,
			want: "tarballs",
		},
		{
			name: "source",
			r:    TypeSource,
			want: "source",
		},
		{
			name: "empty",
			r:    0,
			want: "",
		},
		{
			name: "unknown",
			r:    3,
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.String(); got != tt.want {
				t.Errorf("ResourceType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProject_Tarball(t *testing.T) {
	type fields struct {
		Name       string
		Version    string
		Updated    bool
		ComingSoon bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "xnu",
			fields: fields{
				Name:    "xnu",
				Version: "3789.1.32",
			},
			want: "https://opensource.apple.com/tarballs/xnu/xnu-3789.1.32.tar.gz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Project{
				Name:       tt.fields.Name,
				Version:    tt.fields.Version,
				Updated:    tt.fields.Updated,
				ComingSoon: tt.fields.ComingSoon,
			}
			if got := p.Tarball(); got != tt.want {
				t.Errorf("Project.Tarball() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProject_Source(t *testing.T) {
	type fields struct {
		Name       string
		Version    string
		Updated    bool
		ComingSoon bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "xnu",
			fields: fields{
				Name:    "xnu",
				Version: "3789.1.32",
			},
			want: "https://opensource.apple.com/source/xnu/xnu-3789.1.32/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Project{
				Name:       tt.fields.Name,
				Version:    tt.fields.Version,
				Updated:    tt.fields.Updated,
				ComingSoon: tt.fields.ComingSoon,
			}
			if got := p.Source(); got != tt.want {
				t.Errorf("Project.Source() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexProject(t *testing.T) {
	type args struct {
		typ string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "tarballs",
			args:    args{typ: "tarballs"},
			want:    wantTarballsIndex,
			wantErr: false,
		},
		{
			name:    "source",
			args:    args{typ: "source"},
			want:    wantSourceIndex,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IndexProject(tt.args.typ)
			if (err != nil) != tt.wantErr {
				t.Errorf("IndexProject(%v) error = %v, wantErr %v", tt.args.typ, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IndexProject(%v) = %v, want %v", tt.args.typ, got, tt.want)
			}
		})
	}
}

func TestIndexVersion(t *testing.T) {
	type args struct {
		project string
		typ     string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "Csu (tarball)",
			args:    args{project: "Csu", typ: "tarballs"},
			want:    wantIndexVersionCsu,
			wantErr: false,
		},
		{
			name:    "xnu (source)",
			args:    args{project: "xnu", typ: "source"},
			want:    wantIndexVersionXnu,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IndexVersion(tt.args.project, tt.args.typ)
			if (err != nil) != tt.wantErr {
				t.Errorf("IndexVersion(%v, %v) error = %v, wantErr %v", tt.args.project, tt.args.typ, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IndexVersion(%v, %v) = %v, want %v", tt.args.project, tt.args.typ, got, tt.want)
			}
		})
	}
}

func TestIndexRelease(t *testing.T) {
	type args struct {
		platform Platform
		version  string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "macos 10.12",
			args:    args{platform: MacOS, version: "10.12"},
			want:    wantIndexReleaseMacOS,
			wantErr: false,
		},
		{
			name:    "Xcode 7.3.1",
			args:    args{platform: Xcode, version: "7.3.1"},
			want:    wantIndexReleaseXcode,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IndexRelease(tt.args.platform, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("IndexRelease(%v, %v) error = %v, wantErr %v", tt.args.platform, tt.args.version, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IndexRelease(%v, %v) = %v, want %v", tt.args.platform, tt.args.version, got, tt.want)
			}
		})
	}
}
