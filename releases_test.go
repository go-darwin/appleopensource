// Copyright 2016 Koichi Shiraishi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package appleopensource

import "testing"

func TestPlatform_String(t *testing.T) {
	tests := []struct {
		name string
		p    Platform
		want string
	}{
		{
			name: "Unknown",
			p:    Unknown,
			want: "",
		},
		{
			name: "macos",
			p:    MacOS,
			want: "macos",
		},
		{
			name: "xcode",
			p:    Xcode,
			want: "xcode",
		},
		{
			name: "ios",
			p:    IOS,
			want: "ios",
		},
		{
			name: "server",
			p:    Server,
			want: "server",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("Platform.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
