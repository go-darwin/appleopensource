package appleopensource_test

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"go-darwin.dev/appleopensource"
	"go-darwin.dev/appleopensource/semver"
)

func TestListVersion(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		args    []byte
		want    []semver.Version
		wantErr bool
	}{
		"objc4": {
			args: readTestdata(t, filepath.Join("testdata", "version", "objc4.html")),
			want: []semver.Version{
				"208",
				"217",
				"222",
				"235",
				"237",
				"266",
				"267",
				"267.1",
				"274",
				"371",
				"371.1",
				"371.2",
				"437",
				"437.1",
				"437.3",
				"493.9",
				"493.11",
				"532",
				"532.2",
				"551.1",
				"646",
				"647",
				"680",
				"706",
				"709",
				"709.1",
				"723",
				"750",
				"750.1",
				"756.2",
				"779.1",
				"781",
			},
			wantErr: false,
		},
		"empty": {
			args:    readTestdata(t, filepath.Join("testdata", "version", "empty.html")),
			want:    nil,
			wantErr: true,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := appleopensource.ListVersion(tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ListVersions(\n%s\n) error = %v, wantErr %v", string(tt.args), err, tt.wantErr)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("ListVersions(\n%s\n): (-want +got):\n%s", string(tt.args), diff)
			}
		})
	}
}

func Test_trimZeros(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		args semver.Version
		want semver.Version
	}{
		"trimed": {
			args: semver.Version("532.0.0"),
			want: semver.Version("532"),
		},
		"no trim": {
			args: semver.Version("532.1.3"),
			want: semver.Version("532.1.3"),
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := appleopensource.TrimZeros(tt.args)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("trimZeros(%s): (-want +got):\n%s", tt.args, diff)
			}
		})
	}
}
