// Copyright 2016 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package appleopensource

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

func Test_isExist(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "own file",
			args: args{filename: "./cache_test.go"},
			want: true,
		},
		{
			name: "not exist",
			args: args{filename: "./not_exist.go"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isExist(tt.args.filename); got != tt.want {
				t.Errorf("isExist(%v) = %v, want %v", tt.args.filename, got, tt.want)
			}
		})
	}
}

func Test_xdgCacheHome(t *testing.T) {
	u, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		want string
		env  string
	}{
		{
			name: "Default cache dir",
			want: filepath.Join(u.HomeDir, ".cache"),
			env:  "",
		},
		{
			name: "Use $XDG_CACHE_HOME env",
			want: "/home/testuser/.cache",
			env:  "/home/testuser/.cache",
		},
	}
	for _, tt := range tests {
		os.Setenv("XDG_CACHE_HOME", tt.env)

		t.Run(tt.name, func(t *testing.T) {
			if got := xdgCacheHome(); got != tt.want {
				t.Errorf("xdgCacheHome() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cacheDir(t *testing.T) {
	tmpdir := os.TempDir()

	tests := []struct {
		name        string
		want        string
		wantErr     bool
		env         string
		envXDGCache string
	}{
		{
			name:        "Use $APPLEOPENSOURCE_CACHE_DIR",
			want:        filepath.Join(tmpdir, ".cache", "appleopensource"),
			wantErr:     false,
			env:         filepath.Join(tmpdir, ".cache", "appleopensource"),
			envXDGCache: "",
		},
		{
			name:        "Use $XDG_CACHE_HOME/.cache/appleopensource",
			want:        filepath.Join(tmpdir, "cachehome", "appleopensource"),
			wantErr:     false,
			env:         "",
			envXDGCache: filepath.Join(tmpdir, "cachehome"),
		},
	}
	for _, tt := range tests {
		os.Setenv("APPLEOPENSOURCE_CACHE_DIR", tt.env)
		os.Setenv("XDG_CACHE_HOME", tt.envXDGCache)

		t.Run(tt.name, func(t *testing.T) {
			got, err := cacheDir()
			if (err != nil) != tt.wantErr {
				t.Errorf("cacheDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("cacheDir() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cacheFile(t *testing.T) {
	type args struct {
		typ string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
	// TODO(zchee): avoid always fetch HTML
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cacheFile(tt.args.typ)
			if (err != nil) != tt.wantErr {
				t.Errorf("cacheFile(%v) error = %v, wantErr %v", tt.args.typ, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("cacheFile(%v) = %v, want %v", tt.args.typ, got, tt.want)
			}
		})
	}
}
