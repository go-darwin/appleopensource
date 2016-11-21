// Copyright 2016 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package appleopensource

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"reflect"
	"testing"
)

// testListTimestamp date of created the package list from https://opensource.apple.com
const timestamp = "2016-11-20"

var (
	tarballsList []string
	sourceList   []string

	tarballsListGolden = fmt.Sprintf("testdata/packagelist_tarballs_%s.golden", timestamp)
	sourceListGolden   = fmt.Sprintf("testdata/packagelist_source_%s.golden", timestamp)
)

func TestMain(m *testing.M) {
	flag.Parse()

	tarballsList = readTestFile(tarballsListGolden)
	sourceList = readTestFile(sourceListGolden)

	os.Exit(m.Run())
}

func readTestFile(filename string) []string {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	list := []string{}
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		list = append(list, scan.Text())
	}

	return list
}

func TestListTarballs(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{
			name:    "match all sources list",
			want:    tarballsList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListTarballs()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListTarballs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListTarballs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListSource(t *testing.T) {
	tests := []struct {
		name    string
		want    []string
		wantErr bool
	}{
		{
			name:    "match all sources list",
			want:    sourceList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListSource()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listProject(t *testing.T) {
	type args struct {
		typ ListMode
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "tarballs",
			args:    args{typ: "tarballs"},
			want:    tarballsList,
			wantErr: false,
		},
		{
			name:    "source",
			args:    args{typ: "source"},
			want:    sourceList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListPackage(tt.args.typ)
			if (err != nil) != tt.wantErr {
				t.Errorf("listProject(%v) error = %v, wantErr %v", tt.args.typ, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listProject(%v) = %v, want %v", tt.args.typ, got, tt.want)
			}
		})
	}
}
